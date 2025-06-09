package server

import (
	"encoding/json"
	"net/http"
)

func (server *Server) getIBCInfo(w http.ResponseWriter, r *http.Request) {
	resp := server.QueryIBCInfo()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(resp)

	return
}

func (server *Server) getClientHealth(w http.ResponseWriter, r *http.Request) {
	resp := server.QueryClientHealth()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(resp)

	return
}

func (server *Server) getIBCPacket(w http.ResponseWriter, r *http.Request) {
	resp := server.QueryIBCPacket()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(resp)

	return
}
