package app

import (
	"context"
	"fmt"

	"github.com/dlvlabs/ibcmon/client/grpc"
	"github.com/dlvlabs/ibcmon/logger"

	connectionTypes "github.com/cosmos/ibc-go/v10/modules/core/03-connection/types"
	channelTypes "github.com/cosmos/ibc-go/v10/modules/core/04-channel/types"
)

// set all of open channels
func (channels *Channels) setOpenChannels(
	ctx context.Context,
	grpc *grpc.Client,
	connectionId string,
	counterparty *connectionTypes.Counterparty,
) error {
	connectionChnnels, err := grpc.GetConnectionChannels(ctx, connectionId)
	if err != nil {
		return err
	}

	for _, channel := range connectionChnnels {
		if channel.State != channelTypes.OPEN {
			msg := fmt.Sprintf("skipping channel: %s, only open channels are supported", channel.ChannelId)
			logger.Debug(msg)

			continue
		}

		(*channels)[channel.ChannelId] = &Channel{
			PortId: channel.PortId,
			Counterparty: &Counterparty{
				ClientId:     counterparty.ClientId,
				ConnectionId: counterparty.ConnectionId,
				ChannelId:    channel.Counterparty.ChannelId,
				PortId:       channel.Counterparty.PortId,
			},
		}
	}

	return nil
}
