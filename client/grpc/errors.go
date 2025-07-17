package grpc

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var UNIMPLMENTED = status.Errorf(codes.Unimplemented, "unknown method NextSequenceSend for service ibc.core.channel.v1.Query")

func CONNECTION_NOT_FOUND(clientId string) error {
	return status.Errorf(codes.NotFound, "%s: light client connection paths not found", clientId)
}
