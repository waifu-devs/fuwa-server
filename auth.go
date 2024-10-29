package main

import "net/http"

func setAuthRoutes(mux *httpMux) {
	mux.HandleFunc("GET /join/{serverID}", joinServer(mux))
}

func joinServer(mux *httpMux) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}
