package main

import (
	"context"
	"embed"
	"log/slog"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

//go:embed database/migrations/*.sql
var dbMigrationsFS embed.FS

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "6942"
	}
	DSN := os.Getenv("DATABASE_URL")
	if DSN == "" {
		DSN = "postgres://postgres@localhost:5432/fuwa"
	}
	logLevel := slog.LevelInfo
	switch os.Getenv("FUWA_LOG_LEVEL") {
	case "debug":
		logLevel = slog.LevelDebug
	}

	logHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: logLevel,
	})
	ctx := context.Background()
	log := slog.New(logHandler)
	dbConn, err := pgx.Connect(ctx, DSN)
	if err != nil {
		log.Error("could not connect to database", "error", err.Error())
		panic(err)
	}
	defer func() {
		if closeError := dbConn.Close(ctx); closeError != nil {
			log.Error("error closing database", "error", closeError)
			if err == nil {
				err = closeError
			}
		}
	}()
	applyMigrations(dbConn, log)
	sv := &httpServer{
		ServeMux: http.NewServeMux(),
		log:      log,
	}
	sv.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("fuwa.chat"))
	})
	log.Info("listening on :"+port, "port", port)
	err = sv.ListenAndServe(":" + port)
	if err != nil {
		log.Error("could not listen and serve", "error", err.Error())
		panic(err)
	}
}

func applyMigrations(dbConn *pgx.Conn, log *slog.Logger) {
	goose.SetBaseFS(dbMigrationsFS)
	if err := goose.SetDialect("postgres"); err != nil {
		panic(err)
	}
	db := stdlib.OpenDB(*dbConn.Config())
	if err := goose.Up(db, "database/migrations"); err != nil {
		log.Error("could not apply migrations", "error", err.Error())
		panic(err)
	}
	db.Close()

}
