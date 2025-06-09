package grpc

import (
	"context"

	"github.com/pkg/errors"

	"github.com/cosmos/cosmos-sdk/client/grpc/cmtservice"
	codecTypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	clientTypes "github.com/cosmos/ibc-go/v10/modules/core/02-client/types"
	connectionTypes "github.com/cosmos/ibc-go/v10/modules/core/03-connection/types"
	channelTypes "github.com/cosmos/ibc-go/v10/modules/core/04-channel/types"
	"github.com/cosmos/ibc-go/v10/modules/core/exported"
)

func (c *Client) GetChainId(ctx context.Context) (string, error) {
	resp, err := c.cmtServiceClient.GetNodeInfo(
		ctx,
		&cmtservice.GetNodeInfoRequest{},
	)
	if err != nil {
		return "", errors.Wrap(err, "failed to get chain id")
	}

	return resp.DefaultNodeInfo.Network, nil
}

func (c *Client) GetClientStates(ctx context.Context) (clientTypes.IdentifiedClientStates, error) {
	var clientStates clientTypes.IdentifiedClientStates

	page := &query.PageRequest{}
	for {
		resp, err := c.clientQueryClient.ClientStates(
			ctx,
			&clientTypes.QueryClientStatesRequest{
				Pagination: page,
			},
		)
		if err != nil {
			return nil, errors.Wrap(err, "failed to get client states")
		}

		clientStates = append(clientStates, resp.ClientStates...)

		if len(resp.Pagination.NextKey) == 0 {
			break
		}
		page.Key = resp.Pagination.NextKey
	}

	return clientStates, nil
}

func (c *Client) GetClientState(ctx context.Context, clientId string) (*codecTypes.Any, error) {
	resp, err := c.clientQueryClient.ClientState(
		ctx,
		&clientTypes.QueryClientStateRequest{
			ClientId: clientId,
		},
	)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get client state for client: %s", clientId)
	}

	return resp.ClientState, nil
}

func (c *Client) GetClientStatus(ctx context.Context, clientId string) (bool, error) {
	resp, err := c.clientQueryClient.ClientStatus(
		ctx,
		&clientTypes.QueryClientStatusRequest{
			ClientId: clientId,
		},
	)
	if err != nil {
		return false, errors.Wrapf(err, "failed to get client status for client: %s", clientId)
	}

	if exported.Active.String() != resp.Status {
		return false, nil
	}
	return true, nil
}

func (c *Client) GetClientConnections(ctx context.Context, clientId string) ([]string, error) {
	resp, err := c.connectionQueryClient.ClientConnections(
		ctx,
		&connectionTypes.QueryClientConnectionsRequest{
			ClientId: clientId,
		},
	)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get client connections for client: %s", clientId)
	}

	return resp.ConnectionPaths, nil
}

func (c *Client) GetConnection(ctx context.Context, connectionId string) (*connectionTypes.ConnectionEnd, error) {
	resp, err := c.connectionQueryClient.Connection(
		ctx,
		&connectionTypes.QueryConnectionRequest{
			ConnectionId: connectionId,
		},
	)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get connection: %s", connectionId)
	}

	return resp.Connection, nil
}

func (c *Client) GetConnectionChannels(ctx context.Context, connectionId string) ([]*channelTypes.IdentifiedChannel, error) {
	var channelStates []*channelTypes.IdentifiedChannel

	page := &query.PageRequest{}
	for {
		resp, err := c.channelQueryClient.ConnectionChannels(
			ctx,
			&channelTypes.QueryConnectionChannelsRequest{
				Connection: connectionId,
				Pagination: page,
			},
		)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to get connection channels for connection: %s", connectionId)
		}

		channelStates = append(channelStates, resp.Channels...)

		if len(resp.Pagination.NextKey) == 0 {
			break
		}
		page.Key = resp.Pagination.NextKey
	}

	return channelStates, nil
}

func (c *Client) GetConsensusState(
	ctx context.Context,
	clientId string, revisionNumber uint64,
	revisionHeight uint64,
) (*codecTypes.Any, error) {
	resp, err := c.clientQueryClient.ConsensusState(
		ctx,
		&clientTypes.QueryConsensusStateRequest{
			ClientId:       clientId,
			RevisionNumber: revisionNumber,
			RevisionHeight: revisionHeight,
		},
	)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get consensus state for client: %s", clientId)
	}

	return resp.ConsensusState, nil
}

func (c *Client) GetNextSequenceSend(ctx context.Context, channelId, portId string) (uint64, error) {
	resp, err := c.channelQueryClient.NextSequenceSend(
		ctx,
		&channelTypes.QueryNextSequenceSendRequest{
			ChannelId: channelId,
			PortId:    portId,
		},
	)
	if err != nil {
		if errors.Is(err, UNIMPLMENTED) {
			// Not support `NextSequenceSend` query
			return 0, errors.Cause(err)
		}

		return 0, errors.Wrapf(err, "failed to get next sequence send for channel: %s", channelId)
	}

	return resp.NextSequenceSend, nil
}
