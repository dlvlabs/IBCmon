package app

import (
	"context"
	"fmt"

	"github.com/dlvlabs/ibcmon/client/grpc"
	"github.com/dlvlabs/ibcmon/logger"

	connectionTypes "github.com/cosmos/ibc-go/v10/modules/core/03-connection/types"
)

// set all of open connections
func (connections *Connections) setOpenConnections(ctx context.Context, grpc *grpc.Client, clientId string) error {
	clientConnections, err := grpc.GetClientConnections(ctx, clientId)
	if err != nil {
		return err
	}

	for _, connectionId := range clientConnections {
		connection, err := grpc.GetConnection(ctx, connectionId)
		if err != nil {
			return err
		}

		if connection.State != connectionTypes.OPEN {
			msg := fmt.Sprintf("skipping connection: %s, only open connections are supported", connectionId)
			logger.Debug(msg)

			continue
		}

		channels := make(Channels)
		err = channels.setOpenChannels(ctx, grpc, connectionId, &connection.Counterparty)
		if err != nil {
			return err
		}

		(*connections)[connectionId] = channels
	}

	return nil
}
