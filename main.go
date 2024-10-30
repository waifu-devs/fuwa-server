package main

import (
	"database/sql"
	"embed"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"path"
	"runtime"
	"sync"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pressly/goose/v3"
	"github.com/waifu-devs/fuwa-server/database"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/semconv/v1.4.0"
	"context"
)

//go:embed database/migrations/*.sql
var dbMigrationsFS embed.FS

type server struct {
	readDB  *database.Queries
	writeDB *database.Queries
}

func main() {
	// ctx := context.Background()
	cfg := loadConfigFromEnv()

	logHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: cfg.logLevel,
	})
	log := slog.New(logHandler)
	if len(cfg.encryptionKey) == 0 {
		log.Error("missing FUWA_ENCRYPTION_KEY")
		panic("missing FUWA_ENCRYPTION_KEY")
	}
	if len(cfg.signingKey) == 0 {
		log.Error("missing FUWA_SIGNING_KEY")
		panic("missing FUWA_SIGNING_KEY")
	}
	// Create directory for the database file if it does not exist
	err := os.MkdirAll(cfg.dataPath, 0755)
	if err != nil {
		log.Error("could not create data directory", "error", err.Error())
		panic(err)
	}

	// Initialize OpenTelemetry SDK if configuration values are present
	var tracerProvider *sdktrace.TracerProvider
	if cfg.otelEndpoint != "" && cfg.otelServiceName != "" {
		exporter, err := otlptracehttp.New(context.Background(), 
			otlptracehttp.WithEndpoint(cfg.otelEndpoint),
			otlptracehttp.WithHeaders(map[string]string{
				"Authorization": "Bearer " + cfg.otelAuthToken,
			}),
		)
		if err != nil {
			log.Error("could not create OTLP trace exporter", "error", err.Error())
			panic(err)
		}
		tracerProvider = sdktrace.NewTracerProvider(
			sdktrace.WithBatcher(exporter),
			sdktrace.WithResource(resource.NewWithAttributes(
				semconv.SchemaURL,
				semconv.ServiceNameKey.String(cfg.otelServiceName),
			)),
		)
		otel.SetTracerProvider(tracerProvider)
		defer func() {
			if err := tracerProvider.Shutdown(context.Background()); err != nil {
				log.Error("could not shutdown tracer provider", "error", err.Error())
			}
		}()
	}

	// Initialize the serverConnections map
	serverConnections := &sync.Map{}

	// List all files in the cfg.dataPath directory
	files, err := os.ReadDir(cfg.dataPath)
	if err != nil {
		log.Error("could not read data directory", "error", err.Error())
		panic(err)
	}

	// Create connections to all the databases
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		serverID := file.Name()
		if path.Ext(serverID) != ".db" {
			continue
		}
		writeDB, err := createDatabaseConnection(cfg, serverID, 1)
		if err != nil {
			log.Error("could not connect to database", "error", err.Error())
			panic(err)
		}
		readDB, err := createDatabaseConnection(cfg, serverID, max(4, runtime.NumCPU()))
		if err != nil {
			log.Error("could not connect to database", "error", err.Error())
			panic(err)
		}
		serverConnections.Store(serverID[:len(serverID)-3], &server{
			readDB:  database.New(database.NewTracedDB(readDB)),
			writeDB: database.New(database.NewTracedDB(writeDB)),
		})
		applyMigrations(writeDB, log.With("serverID", serverID[:len(serverID)-3]))
	}

	defer func() {
		serverConnections.Range(func(key, value any) bool {
			server := value.(*server)
			database.GetDBFromQueries(server.readDB).Close()
			database.GetDBFromQueries(server.writeDB).Close()
			return true
		})
	}()

	mux := &httpMux{
		ServeMux:  http.NewServeMux(),
		log:       log,
		cfg:       cfg,
		serverDBs: serverConnections,
	}
	server := &http.Server{
		Addr:              ":" + cfg.port,
		Handler:           mux,
		ReadHeaderTimeout: 10 * time.Second,
	}
	setServerRoutes(mux)
	setChannelRoutes(mux)
	setMessageRoutes(mux)
	go func() {
		log.Info("listening on :"+cfg.port, "port", cfg.port)
		err = server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Error("could not listen and serve", "error", err.Error())
			panic(err)
		}
	}()
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, os.Kill)
	<-sig
	log.Info("shutting down")
	err = server.Shutdown(nil)
	if err != nil {
		log.Error("could not shutdown server", "error", err.Error())
	}
}

func createConnectionString(cfg *config, serverID string) string {
	connectionUrlParams := make(url.Values)
	connectionUrlParams.Add("_txlock", "immediate")
	connectionUrlParams.Add("_journal_mode", "WAL")
	connectionUrlParams.Add("_busy_timeout", "5000")
	connectionUrlParams.Add("_synchronous", "NORMAL")
	connectionUrlParams.Add("_cache_size", "1000000000")
	connectionUrlParams.Add("_foreign_keys", "true")
	return fmt.Sprintf("file:%s?%s", path.Join(cfg.dataPath, fmt.Sprintf("%s", serverID)), connectionUrlParams.Encode())
}

func createDatabaseConnection(cfg *config, serverID string, maxOpenConns int) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", createConnectionString(cfg, serverID))
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(maxOpenConns)
	_, err = db.Exec("PRAGMA temp_store = memory")
	if err != nil {
		return nil, err
	}
	return db, nil
}

type gooseLogger struct {
	log *slog.Logger
}

func (l *gooseLogger) Printf(format string, v ...interface{}) {
	l.log.Info(fmt.Sprintf(format, v...))
}

func (l *gooseLogger) Fatalf(format string, v ...interface{}) {
	l.log.Error(fmt.Sprintf(format, v...))
}

func applyMigrations(db *sql.DB, log *slog.Logger) {
	goose.SetBaseFS(dbMigrationsFS)
	goose.SetLogger(&gooseLogger{log: log})
	if err := goose.SetDialect("sqlite"); err != nil {
		log.Error("could not set dialog for goose", "error", err.Error())
		panic(err)
	}
	if err := goose.Up(db, "database/migrations"); err != nil {
		log.Error("could not apply migrations", "error", err.Error())
		panic(err)
	}
}
