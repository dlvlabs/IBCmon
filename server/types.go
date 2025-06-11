package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/dlvlabs/ibcmon/app"
)

// response for "/ibc-info"
type (
	IBCInfos []IBCInfo

	IBCInfo struct {
		Updated time.Time `json:"updated"`

		Source      IBC `json:"source"`
		Destination IBC `json:"destination"`
	}
)

// response for "/client-health"
type (
	ClientHealths []ClientHealth
	ClientHealth  struct {
		Health        bool      `json:"health"`
		ClientUpdated time.Time `json:"client_updated"`

		Source      string `json:"source"`
		Destination string `json:"destination"`

		ClientId       string        `json:"client_id"`
		TrustingPeriod time.Duration `json:"trusting_period"`
	}
)

// response for "/ibc-packet"
type (
	IBCPackets []IBCPacket
	IBCPacket  struct {
		Updated time.Time `json:"updated"`

		Health bool `json:"health"`

		Source      IBC `json:"source"`
		Destination IBC `json:"destination"`

		Sequence             uint64         `json:"sequence"`
		ConsecutiveMissed    uint64         `json:"consecutive_missed"`
		LatestSucceedPackets SucceedPackets `json:"latest_succeed_packets"`
	}
	// PakcetType => SucceedPacket
	SucceedPackets map[string]SucceedPacket
	SucceedPacket  struct {
		Hash     string `json:"hash"`
		Sequence uint64 `json:"sequence"`
		Data     string `json:"data"`
	}
)

type IBC struct {
	Path string `json:"path"`

	ChainId      string
	ClientId     string
	ConnectionId string
	ChannelId    string
	PortId       string
}

func newIBC(chainId, clientId, connectionId, channelId, portId string) IBC {
	path := fmt.Sprintf(
		"%s(%s/%s/%s/%s)",
		chainId,
		clientId, connectionId, channelId, portId,
	)
	return IBC{
		Path: path,

		ChainId:      chainId,
		ClientId:     clientId,
		ConnectionId: connectionId,
		ChannelId:    channelId,
		PortId:       portId,
	}
}

type Server struct {
	Store        *app.Store
	mux          *http.ServeMux
	port         string
	MetricPrefix string
}

func NewServer(store *app.Store, port int, prefix string) *Server {
	server := Server{
		Store:        store,
		mux:          http.NewServeMux(),
		port:         fmt.Sprintf(":%d", port),
		MetricPrefix: prefix,
	}

	return &server
}
