package app

import (
	"context"
	"fmt"
	"time"

	"github.com/dlvlabs/ibcmon/alert"
	"github.com/dlvlabs/ibcmon/client/grpc"
	"github.com/dlvlabs/ibcmon/logger"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/ibc-go/v10/modules/core/exported"
	tendermint "github.com/cosmos/ibc-go/v10/modules/light-clients/07-tendermint"
)

func (app *App) checkClientsHealth(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Minute)
	defer cancel()

	err := app.connectGRPCs()
	if err != nil {
		return err
	}
	defer func() {
		err = app.terminateGRPCs()
		if err != nil {
			logger.Error(err)
		}
	}()

	g, ctx := errgroup.WithContext(ctx)

	for chainId, clients := range app.Store.IBCInfo {
		for clientId, client := range clients {
			g.Go(func() error {
				err := client.checkHealth(
					ctx, app.grpcs[chainId], app.cdc, clientId,
					app.cfg.Rule.ClientExpiredWarningTime,
				)
				if err != nil {
					logger.Error(err)
					return err
				}

				return nil
			})
		}
	}

	return g.Wait()
}

func (client *Client) checkHealth(
	ctx context.Context,
	grpc *grpc.Client,
	cdc codectypes.InterfaceRegistry,
	clientId string,
	warningTime time.Duration,
) error {
	// Get client state for new RevisionNumber and RevisionHeight
	state, err := grpc.GetClientState(ctx, clientId)
	if err != nil {
		return err
	}

	if state.TypeUrl != "/ibc.lightclients.tendermint.v1.ClientState" {
		msg := fmt.Sprintf("skipping client: %s, only tendermint light clients are supported", clientId)
		logger.Debug(msg)
	}

	var iClientState exported.ClientState
	if err := cdc.UnpackAny(state, &iClientState); err != nil {
		return errors.Wrapf(err, "failed to unpack client state for client: %s", clientId)
	}
	clientState, ok := iClientState.(*tendermint.ClientState)
	if !ok {
		return errors.Wrapf(err, "invalid client state type: %T", iClientState)
	}

	// Get consensus state for client updated timestamp
	state, err = grpc.GetConsensusState(ctx, clientId, clientState.LatestHeight.RevisionNumber, clientState.LatestHeight.RevisionHeight)
	if err != nil {
		return err
	}

	var iConsensusState exported.ConsensusState
	if err := cdc.UnpackAny(state, &iConsensusState); err != nil {
		return errors.Wrapf(err, "failed to unpack consensus state for client: %s", clientId)
	}
	consensusState, ok := iConsensusState.(*tendermint.ConsensusState)
	if !ok {
		return errors.Wrapf(err, "invalid consensus state type: %T", iConsensusState)
	}

	// Update stored client info
	// Don't need mutex lock here, because update each client
	if client.warnExpiration(consensusState.Timestamp, warningTime) {
		client.Health = false

		msg := fmt.Sprintf(
			"client %s(%s) would be expired, consensus state timestamp: %s",
			clientId, client.ChainId, consensusState.Timestamp,
		)
		logger.Warn(msg)
		alert.SendTg(msg)

		return err
	}
	client.Health = true

	client.RevisionNumber = clientState.LatestHeight.RevisionNumber
	client.RevisionHeight = clientState.LatestHeight.RevisionHeight
	client.TrustingPeriod = clientState.TrustingPeriod
	client.ClientUpdated = consensusState.Timestamp

	logger.Info(fmt.Sprintf("client %s is healthy", clientId))

	return nil
}

func (client *Client) warnExpiration(latestTimestamp time.Time, warningTime time.Duration) bool {
	// "expired" means: timestamp + trusting period <= current time
	warnExpirationTime := latestTimestamp.Add(client.TrustingPeriod).Add(-warningTime)
	return !warnExpirationTime.After(time.Now().UTC())
}
