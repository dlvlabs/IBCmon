package app

import (
	"context"
	"fmt"

	"github.com/dlvlabs/ibcmon/client/grpc"
	"github.com/dlvlabs/ibcmon/logger"
	"github.com/pkg/errors"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/ibc-go/v10/modules/core/exported"
	tendermint "github.com/cosmos/ibc-go/v10/modules/light-clients/07-tendermint"
)

// set all of active and tendermint clients
func (clients *Clients) setActiveClients(ctx context.Context, grpc *grpc.Client, cdc codectypes.InterfaceRegistry) error {
	clientStates, err := grpc.GetClientStates(ctx)
	if err != nil {
		return err
	}

	for _, cs := range clientStates {
		// only support tendermint light clients
		if cs.ClientState.TypeUrl != "/ibc.lightclients.tendermint.v1.ClientState" {
			msg := fmt.Sprintf("skipping client: %s, only tendermint light clients are supported", cs.ClientId)
			logger.Debug(msg)

			continue
		}

		status, err := grpc.GetClientStatus(ctx, cs.ClientId)
		if err != nil {
			return err
		}

		if !status {
			msg := fmt.Sprintf("client %s is not active", cs.ClientId)
			logger.Debug(msg)

			continue
		}

		var iClientState exported.ClientState
		if err := cdc.UnpackAny(cs.ClientState, &iClientState); err != nil {
			return errors.Wrapf(err, "failed to unpack client state for client: %s", cs.ClientId)
		}
		clientState, ok := iClientState.(*tendermint.ClientState)
		if !ok {
			return errors.Wrapf(err, "invalid client state type: %T", iClientState)
		}

		connections := make(Connections)
		err = connections.setOpenConnections(ctx, grpc, cs.ClientId)
		if err != nil {
			return err
		}

		(*clients)[cs.ClientId] = &Client{
			Health: true,

			ChainId:        clientState.ChainId,
			RevisionNumber: clientState.LatestHeight.RevisionNumber,
			RevisionHeight: clientState.LatestHeight.RevisionHeight,
			TrustingPeriod: clientState.TrustingPeriod,

			Connections: connections,
		}
	}

	return nil
}

// set active and tendermint client
func (clients *Clients) setActiveClient(
	ctx context.Context,
	grpc *grpc.Client,
	cdc codectypes.InterfaceRegistry,
	clientId string,
) error {
	clientState, err := grpc.GetClientState(ctx, clientId)
	if err != nil {
		return err
	}

	if clientState.TypeUrl != "/ibc.lightclients.tendermint.v1.ClientState" {
		msg := fmt.Sprintf("skipping client: %s, only tendermint light clients are supported", clientId)
		logger.Debug(msg)
	}

	status, err := grpc.GetClientStatus(ctx, clientId)
	if err != nil {
		return err
	}

	if !status {
		msg := fmt.Sprintf("client %s is not active", clientId)
		logger.Debug(msg)

		return nil
	}

	var iClientState exported.ClientState
	if err := cdc.UnpackAny(clientState, &iClientState); err != nil {
		return errors.Wrapf(err, "failed to unpack client state for client: %s", clientId)
	}
	cs, ok := iClientState.(*tendermint.ClientState)
	if !ok {
		return errors.Wrapf(err, "invalid client state type: %T", iClientState)
	}

	connections := make(Connections)
	err = connections.setOpenConnections(ctx, grpc, clientId)
	if err != nil {
		return err
	}

	(*clients)[clientId] = &Client{
		Health: true,

		ChainId:        cs.ChainId,
		RevisionNumber: cs.LatestHeight.RevisionNumber,
		RevisionHeight: cs.LatestHeight.RevisionHeight,
		TrustingPeriod: cs.TrustingPeriod,

		Connections: connections,
	}

	return nil
}
