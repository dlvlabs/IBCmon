package rpc

import (
	"github.com/dlvlabs/ibcmon/logger"
	"github.com/pkg/errors"
)

func (c *Client) Connect() error {
	// for websocket connection
	err := c.rpcClient.Start()
	if err != nil {
		return errors.Wrap(err, "failed to connect to rpc client")
	}

	logger.Info("RPC connected")

	return nil
}

func (c *Client) Terminate() error {
	// for websocket connection
	err := c.rpcClient.Stop()
	if err != nil {
		return errors.Wrap(err, "failed to terminate rpc client")
	}

	logger.Info("RPC connection terminated")

	return nil
}
