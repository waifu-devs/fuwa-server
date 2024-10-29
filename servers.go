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
		servers, err := mux.readDB.ListServers(r.Context(), database.ListServersParams{
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
		err = mux.writeDB.CreateServer(r.Context(), req)
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
