package grpc

import (
	"github.com/cosmos/cosmos-sdk/client/grpc/cmtservice"
	clientTypes "github.com/cosmos/ibc-go/v10/modules/core/02-client/types"
	connectionTypes "github.com/cosmos/ibc-go/v10/modules/core/03-connection/types"
	channelTypes "github.com/cosmos/ibc-go/v10/modules/core/04-channel/types"
	"google.golang.org/grpc"
)

type Client struct {
	host    string
	tlsConn bool
	conn    *grpc.ClientConn

	clientQueryClient     clientTypes.QueryClient
	connectionQueryClient connectionTypes.QueryClient
	channelQueryClient    channelTypes.QueryClient
	cmtServiceClient      cmtservice.ServiceClient
}

func New(host string, tlsConn bool) *Client {
	return &Client{
		host:    host,
		tlsConn: tlsConn,
	}
}
