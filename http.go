package main

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
	"sync"

	"github.com/oklog/ulid/v2"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/trace/noop"
)

// httpMux is a HTTP server
type httpMux struct {
	*http.ServeMux

	log *slog.Logger

	cfg *config

	serverDBs *sync.Map
}

type httpResponse struct {
	w          http.ResponseWriter
	statusCode int
}

func (r *httpResponse) Header() http.Header {
	return r.w.Header()
}

func (r *httpResponse) Write(b []byte) (int, error) {
	return r.w.Write(b)
}

func (r *httpResponse) WriteHeader(statusCode int) {
	r.statusCode = statusCode
	r.w.WriteHeader(statusCode)
}

func (r *httpResponse) Flush() {
	r.w.(http.Flusher).Flush()
}

type requestIDContextKey string

const RequestIDContextKey requestIDContextKey = "RequestID"

// Add helper functions for JSON responses
func writeJSON(w http.ResponseWriter, statusCode int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

func writeJSONError(w http.ResponseWriter, statusCode int, err error) {
	writeJSON(w, statusCode, map[string]any{"error": err.Error()})
}

func (s *httpMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var tracer trace.Tracer
	if s.cfg.otelEndpoint != "" && s.cfg.otelServiceName != "" {
		tracer = otel.Tracer("httpMux")
	} else {
		tracer = noop.NewTracerProvider().Tracer("")
	}

	ctx, span := tracer.Start(r.Context(), "request-incoming")
	defer span.End()

	response := &httpResponse{
		w:          w,
		statusCode: 200,
	}

	requestID := ulid.Make()

	span.SetAttributes(
		attribute.String("http.method", r.Method),
		attribute.String("http.path", r.URL.Path),
		attribute.String("http.user_agent", r.Header.Get("User-Agent")),
		attribute.String("http.request_id", requestID.String()),
	)

	ctx = context.WithValue(ctx, RequestIDContextKey, requestID)
	r = r.WithContext(ctx)

	response.w.Header().Set("X-Fuwa-Request-Id", requestID.String())

	// Handle the request
	s.ServeMux.ServeHTTP(response, r)

	s.log.Info(
		r.Method+" ["+strconv.FormatInt(int64(response.statusCode), 10)+"] "+r.URL.Path,
		"method", r.Method,
		"path", r.URL.Path,
		"status", response.statusCode,
		"user_agent", r.Header.Get("User-Agent"),
		"request_id", requestID,
	)
}
