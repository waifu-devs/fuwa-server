package main

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/oklog/ulid/v2"
	"github.com/waifu-devs/fuwa-server/database"
)

func setServerRoutes(mux *httpMux) {
	mux.HandleFunc("GET /servers", listServers(mux))
	mux.HandleFunc("POST /servers", createServers(mux))
}

func listServers(mux *httpMux) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		serverID := r.PathValue("serverID")
		servers, err := mux.serverDBs[serverID].readDB.ListServers(r.Context(), database.ListServersParams{
			Limit: 10,
		})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("could not list servers: " + err.Error()))
			return
		}
		json.NewEncoder(w).Encode(map[string]any{
			"servers": servers,
		})
	}
}

func validateCreateServerParams(args database.CreateServerParams) error {
	if args.Name == "" {
		return errors.New("server name must not be empty")
	}
	return nil
}

func createServers(mux *httpMux) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := database.CreateServerParams{}
		err := json.NewDecoder(r.Body).Decode(&req)
		defer r.Body.Close()
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("could not decode: " + err.Error()))
			return
		}
		err = validateCreateServerParams(req)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("invalid params: " + err.Error()))
			return
		}
		req.ServerID = ulid.Make()

		// Create a new database file for the new server
		writeDB, err := createDatabaseConnection(mux.cfg, req.ServerID.String(), 1)
		if err != nil {
			mux.log.Error("could not create database for new server", "error", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("could not create database: " + err.Error()))
			return
		}
		readDB, err := createDatabaseConnection(mux.cfg, req.ServerID.String(), max(4, runtime.NumCPU()))
		if err != nil {
			mux.log.Error("could not create database for new server", "error", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("could not create database: " + err.Error()))
			return
		}
		defer writeDB.Close()
		defer readDB.Close()

		// Apply migrations to the new database
		applyMigrations(writeDB, mux.log)

		// Add the new database connection to the serverDBs map
		mux.serverDBs[req.ServerID.String()] = &server{
			readDB:  database.New(readDB),
			writeDB: database.New(writeDB),
		}

		err = mux.serverDBs[req.ServerID.String()].writeDB.CreateServer(r.Context(), req)
		if err != nil {
			mux.log.Error("could not create server", "error", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("could not create: " + err.Error()))
			return
		}

		json.NewEncoder(w).Encode(map[string]any{
			"server": req,
		})
	}
}
