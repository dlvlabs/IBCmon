package grpc

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var UNIMPLMENTED = status.Errorf(codes.Unimplemented, "unknown method NextSequenceSend for service ibc.core.channel.v1.Query")
