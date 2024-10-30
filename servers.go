package main

import (
	"fmt"
	"net/http"
	"runtime"

	"github.com/oklog/ulid/v2"
	"github.com/waifu-devs/fuwa-server/database"
	"go.opentelemetry.io/otel"
)

func setServerRoutes(mux *httpMux) {
	mux.HandleFunc("POST /servers", createServers(mux))
}

func createServers(mux *httpMux) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tracer := otel.Tracer("createServers")
		_, span := tracer.Start(r.Context(), "createServers")
		defer span.End()

		serverID := ulid.Make()

		// Create a new database file for the new server
		writeDB, err := createDatabaseConnection(mux.cfg, serverID.String()+".db", 1)
		if err != nil {
			mux.log.Error("could not create database for new server", "error", err.Error())
			writeJSONError(w, http.StatusInternalServerError, fmt.Errorf("could not create database: %w", err))
			return
		}
		readDB, err := createDatabaseConnection(mux.cfg, serverID.String()+".db", max(4, runtime.NumCPU()))
		if err != nil {
			mux.log.Error("could not create database for new server", "error", err.Error())
			writeJSONError(w, http.StatusInternalServerError, fmt.Errorf("could not create database: %w", err))
			return
		}
		defer writeDB.Close()
		defer readDB.Close()

		// Apply migrations to the new database
		applyMigrations(writeDB, mux.log)

		// Add the new database connection to the serverDBs map
		mux.serverDBs.Store(serverID.String(), &server{
			readDB:  database.New(database.NewTracedDB(readDB)),
			writeDB: database.New(database.NewTracedDB(writeDB)),
		})

		writeJSON(w, http.StatusOK, map[string]any{
			"serverID": serverID.String(),
		})
	}
}
