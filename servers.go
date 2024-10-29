package main

import (
	"net/http"
	"runtime"

	"github.com/oklog/ulid/v2"
	"github.com/waifu-devs/fuwa-server/database"
)

func setServerRoutes(mux *httpMux) {
	mux.HandleFunc("POST /servers", createServers(mux))
}

func createServers(mux *httpMux) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// err := json.NewDecoder(r.Body).Decode(&req)
		// defer r.Body.Close()
		// if err != nil {
		// 	w.WriteHeader(http.StatusBadRequest)
		// 	w.Write([]byte("could not decode: " + err.Error()))
		// 	return
		// }
		// err = validateCreateServerParams(req)
		// if err != nil {
		// 	w.WriteHeader(http.StatusBadRequest)
		// 	w.Write([]byte("invalid params: " + err.Error()))
		// 	return
		// }
		serverID := ulid.Make()

		// Create a new database file for the new server
		writeDB, err := createDatabaseConnection(mux.cfg, serverID.String(), 1)
		if err != nil {
			mux.log.Error("could not create database for new server", "error", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("could not create database: " + err.Error()))
			return
		}
		readDB, err := createDatabaseConnection(mux.cfg, serverID.String(), max(4, runtime.NumCPU()))
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
		mux.serverDBs[serverID.String()] = &server{
			readDB:  database.New(readDB),
			writeDB: database.New(writeDB),
		}

		// json.NewEncoder(w).Encode(map[string]any{
		// 	"server": req,
		// })
	}
}
