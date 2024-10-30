package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/waifu-devs/fuwa-server/database"
)

func setChannelRoutes(mux *httpMux) {
	mux.HandleFunc("GET /{serverID}/channels", listChannels(mux))
	mux.HandleFunc("POST /{serverID}/channels", createChannels(mux))
}

func listChannels(mux *httpMux) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		serverID := r.PathValue("serverID")
		readDBI, ok := mux.serverDBs.Load(serverID)
		if !ok {
			writeJSONError(w, http.StatusNotFound, errors.New("server not found"))
			return
		}
		readDB := readDBI.(*server).readDB
		channels, err := readDB.ListChannels(r.Context(), database.ListChannelsParams{
			Limit: 10,
		})
		if err != nil {
			mux.log.Error("could not list channels", "error", err.Error())
			writeJSONError(w, http.StatusInternalServerError, err)
			return
		}

		writeJSON(w, http.StatusOK, map[string]any{"channels": channels})
	}
}

func validateCreateChannelParams(args database.CreateChannelParams) error {
	switch args.Type {
	case "text":
		break
	default:
		return errors.New("invalid channel type")
	}
	return nil
}

func createChannels(mux *httpMux) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		serverID := r.PathValue("serverID")
		req := database.CreateChannelParams{}
		err := json.NewDecoder(r.Body).Decode(&req)
		defer r.Body.Close()
		if err != nil {
			writeJSONError(w, http.StatusBadRequest, fmt.Errorf("could not decode: %w", err))
			return
		}
		err = validateCreateChannelParams(req)
		if err != nil {
			writeJSONError(w, http.StatusBadRequest, fmt.Errorf("invalid params: %w", err))
			return
		}
		writeDBI, ok := mux.serverDBs.Load(serverID)
		if !ok {
			writeJSONError(w, http.StatusNotFound, errors.New("server not found"))
			return
		}
		writeDB := writeDBI.(*server).writeDB

		req.ChannelID = ulid.Make()
		req.CreatedAt = time.Now().UnixMilli()
		err = writeDB.CreateChannel(r.Context(), req)
		if err != nil {
			mux.log.Error("could not create channel", "error", err.Error())
			writeJSONError(w, http.StatusInternalServerError, err)
			return
		}

		writeJSON(w, http.StatusOK, map[string]any{
			"channel": req,
		})
	}
}
