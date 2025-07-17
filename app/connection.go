package app

import (
	"context"
	"fmt"

	"github.com/dlvlabs/ibcmon/client/grpc"
	"github.com/dlvlabs/ibcmon/logger"
	"github.com/pkg/errors"

	connectionTypes "github.com/cosmos/ibc-go/v10/modules/core/03-connection/types"
)

// set all of open connections
func (connections *Connections) setOpenConnections(ctx context.Context, grpcClient *grpc.Client, clientId string) error {
	clientConnections, err := grpcClient.GetClientConnections(ctx, clientId)
	if err != nil {
		if errors.Is(errors.Cause(err), grpc.CONNECTION_NOT_FOUND(clientId)) {
			msg := fmt.Sprintf("no connection paths found for client %s", clientId)
			logger.Info(msg)

			return nil
		}
		return err
	}

	for _, connectionId := range clientConnections {
		connection, err := grpcClient.GetConnection(ctx, connectionId)
		if err != nil {
			return err
		}

		if connection.State != connectionTypes.OPEN {
			msg := fmt.Sprintf("skipping connection: %s, only open connections are supported", connectionId)
			logger.Debug(msg)

			continue
		}

		channels := make(Channels)
		err = channels.setOpenChannels(ctx, grpcClient, connectionId, &connection.Counterparty)
		if err != nil {
			return err
		}

		(*connections)[connectionId] = channels
	}

	return nil
}
