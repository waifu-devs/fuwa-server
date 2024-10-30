package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/waifu-devs/fuwa-server/database"
	"go.opentelemetry.io/otel"
)

func setMessageRoutes(mux *httpMux) {
	mux.HandleFunc("POST /{serverID}/channels/{channelID}/messages", createMessage(mux))
	mux.HandleFunc("GET /{serverID}/channels/{channelID}/messages/{messageID}", getMessage(mux))
	mux.HandleFunc("GET /{serverID}/channels/{channelID}/messages", listMessages(mux))
	mux.HandleFunc("GET /{serverID}/channels/{channelID}/subscribe", subscribeMessages(mux))
}

func createMessage(mux *httpMux) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		serverID := r.PathValue("serverID")
		channelID := r.PathValue("channelID")
		req := database.CreateMessageParams{}
		err := json.NewDecoder(r.Body).Decode(&req)
		defer r.Body.Close()
		if err != nil {
			writeJSONError(w, http.StatusBadRequest, fmt.Errorf("could not decode: %w", err))
			return
		}
		req.MessageID = ulid.Make()
		req.Timestamp = time.Now().UnixMilli()
		req.ChannelID, err = ulid.Parse(channelID)
		if err != nil {
			writeJSONError(w, http.StatusBadRequest, errors.New("invalid channel id"))
			return
		}
		writeDBI, ok := mux.serverDBs.Load(serverID)
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("server not found"))
			return
		}
		writeDB := writeDBI.(*server).writeDB

		tracer := otel.Tracer("createMessage")
		ctx, span := tracer.Start(r.Context(), "createMessage")
		defer span.End()

		err = writeDB.CreateMessage(ctx, req)
		if err != nil {
			mux.log.Error("could not create message", "error", err.Error())
			writeJSONError(w, http.StatusInternalServerError, err)
			return
		}

		writeJSON(w, http.StatusOK, map[string]any{
			"message": req,
		})
	}
}

func getMessage(mux *httpMux) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		serverID := r.PathValue("serverID")
		messageID := r.PathValue("messageID")
		validMessageID, err := ulid.Parse(messageID)
		if err != nil {
			writeJSONError(w, http.StatusBadRequest, errors.New("invalid message id"))
			return
		}
		readDBI, ok := mux.serverDBs.Load(serverID)
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("server not found"))
			return
		}
		readDB := readDBI.(*server).readDB

		tracer := otel.Tracer("getMessage")
		ctx, span := tracer.Start(r.Context(), "getMessage")
		defer span.End()

		message, err := readDB.GetMessage(ctx, validMessageID)
		if err != nil {
			mux.log.Error("could not get message", "error", err.Error())
			writeJSONError(w, http.StatusInternalServerError, fmt.Errorf("could not get message: %w", err))
			return
		}

		writeJSON(w, http.StatusOK, map[string]any{"message": message})
	}
}

func listMessages(mux *httpMux) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		serverID := r.PathValue("serverID")
		channelID := r.PathValue("channelID")
		validChannelID, err := ulid.Parse(channelID)
		if err != nil {
			writeJSONError(w, http.StatusBadRequest, errors.New("invalid channel id"))
			return
		}
		messageID := r.URL.Query().Get("messageID")
		var validMessageID ulid.ULID
		if messageID != "" {
			validMessageID, err = ulid.Parse(messageID)
			if err != nil {
				writeJSONError(w, http.StatusBadRequest, errors.New("invalid message id"))
				return
			}
		} else {
			validMessageID = ulid.ULID{}
		}
		readDBI, ok := mux.serverDBs.Load(serverID)
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("server not found"))
			return
		}
		readDB := readDBI.(*server).readDB

		tracer := otel.Tracer("listMessages")
		ctx, span := tracer.Start(r.Context(), "listMessages")
		defer span.End()

		messages, err := readDB.ListMessages(ctx, database.ListMessagesParams{
			ChannelID: validChannelID,
			MessageID: validMessageID,
			Limit:     10,
		})
		if err != nil {
			mux.log.Error("could not list messages", "error", err.Error())
			writeJSONError(w, http.StatusInternalServerError, fmt.Errorf("could not list messages: %w", err))
			return
		}

		writeJSON(w, http.StatusOK, map[string]any{"messages": messages})
	}
}

func subscribeMessages(mux *httpMux) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		serverID := r.PathValue("serverID")
		channelID := r.PathValue("channelID")
		validChannelID, err := ulid.Parse(channelID)
		if err != nil {
			writeJSONError(w, http.StatusBadRequest, errors.New("invalid channel id"))
			return
		}

		messageID := r.URL.Query().Get("messageID")
		var validMessageID ulid.ULID
		if messageID != "" {
			validMessageID, err = ulid.Parse(messageID)
			if err != nil {
				writeJSONError(w, http.StatusBadRequest, errors.New("invalid message id"))
				return
			}
		} else {
			validMessageID = ulid.ULID{}
		}
		readDBI, ok := mux.serverDBs.Load(serverID)
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("server not found"))
			return
		}
		readDB := readDBI.(*server).readDB

		tracer := otel.Tracer("subscribeMessages")
		ctx, span := tracer.Start(r.Context(), "subscribeMessages")
		defer span.End()

		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.(http.Flusher).Flush()

		lastSentMessageID := validMessageID

		for {
			messages, err := readDB.ListMessages(ctx, database.ListMessagesParams{
				ChannelID: validChannelID,
				MessageID: lastSentMessageID,
				Limit:     10,
			})
			if err != nil {
				mux.log.Error("could not list messages", "error", err.Error())
				writeJSONError(w, http.StatusInternalServerError, fmt.Errorf("could not list messages: %w", err))
				return
			}

			for _, message := range messages {
				jsonMsg, err := json.Marshal(message)
				if err != nil {
					mux.log.Error("could not marshal message", "error", err.Error())
					writeJSONError(w, http.StatusInternalServerError, fmt.Errorf("could not marshal message: %w", err))
					return
				}
				w.Write([]byte("id: " + message.MessageID.String() + "\n"))
				w.Write([]byte("event: message\n"))
				w.Write([]byte("data: " + string(jsonMsg) + "\n"))
				w.Write([]byte("retry: 2500\n\n"))
				w.(http.Flusher).Flush()
				lastSentMessageID = message.MessageID
			}

			if len(messages) == 0 {
				w.Write([]byte(": ping\n\n"))
				w.(http.Flusher).Flush()
			}

			time.Sleep(500 * time.Millisecond)
		}
	}
}
