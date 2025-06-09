package grpc

import (
	"crypto/tls"

	"github.com/dlvlabs/ibcmon/logger"
	"github.com/pkg/errors"

	"github.com/cosmos/cosmos-sdk/client/grpc/cmtservice"
	clientTypes "github.com/cosmos/ibc-go/v10/modules/core/02-client/types"
	connectionTypes "github.com/cosmos/ibc-go/v10/modules/core/03-connection/types"
	channelTypes "github.com/cosmos/ibc-go/v10/modules/core/04-channel/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

func (c *Client) Connect() error {
	var opts []grpc.DialOption
	if c.tlsConn {
		opts = append(
			opts,
			grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{})),
		)
	} else {
		opts = append(
			opts,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
	}

	conn, err := grpc.NewClient(
		c.host,
		opts...,
	)
	if err != nil {
		return errors.Wrap(err, "failed to create grpc client")
	}

	c.conn = conn
	c.clientQueryClient = clientTypes.NewQueryClient(conn)
	c.connectionQueryClient = connectionTypes.NewQueryClient(conn)
	c.channelQueryClient = channelTypes.NewQueryClient(conn)
	c.cmtServiceClient = cmtservice.NewServiceClient(conn)

	logger.Info("GRPC client created")

	return nil
}

func (c *Client) Terminate() error {
	err := c.conn.Close()
	if err != nil {
		return errors.Wrap(err, "failed to close grpc connection")
	}

	logger.Info("GRPC connection terminated")

	return nil
}
