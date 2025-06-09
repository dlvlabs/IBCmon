package server

import (
	"fmt"
	"net/http"

	"github.com/dlvlabs/ibcmon/logger"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func (server *Server) Run() error {
	r := prometheus.NewRegistry()
	r.MustRegister(newIBCInfoCollector(server))
	r.MustRegister(newClientHealthCollector(server))
	r.MustRegister(newIBCPacketCollector(server))

	server.mux.HandleFunc("/ibc-info", server.getIBCInfo)
	server.mux.HandleFunc("/client-health", server.getClientHealth)
	server.mux.HandleFunc("/ibc-packet", server.getIBCPacket)
	server.mux.Handle("/metrics", promhttp.HandlerFor(r, promhttp.HandlerOpts{}))

	msg := fmt.Sprintf("starting server on %s", server.port)
	logger.Info(msg)
	httpServer := &http.Server{
		Addr:    server.port,
		Handler: server.mux,
	}

	if err := httpServer.ListenAndServe(); err != nil {
		return err
	}

	return nil
}
