package rpc

import (
	"github.com/dlvlabs/ibcmon/logger"
	"github.com/pkg/errors"

	cmthttp "github.com/cometbft/cometbft/rpc/client/http"
)

type Client struct {
	host string

	rpcClient *cmthttp.HTTP
}

func New(host string) (*Client, error) {
	result := &Client{
		host: host,
	}
	client, err := cmthttp.New(result.host, "/websocket")
	if err != nil {
		return nil, errors.Wrap(err, "failed to create rpc client")
	}
	result.rpcClient = client

	logger.Info("RPC client created")

	return result, nil
}
