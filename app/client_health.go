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
				updatedTimestamp, err := client.checkHealth(
					ctx, app.grpcs[chainId], app.cdc, clientId,
					app.cfg.Rule.ClientExpiredWarningTime,
				)
				if err != nil {
					logger.Error(err)
					return err
				}

				app.updateStore(func() { client.ClientUpdated = updatedTimestamp })

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
) (time.Time, error) {
	cs, err := grpc.GetConsensusState(ctx, clientId, client.RevisionNumber, client.RevisionHeight)
	if err != nil {
		return time.Time{}, err
	}

	var iConsensusState exported.ConsensusState
	if err := cdc.UnpackAny(cs, &iConsensusState); err != nil {
		return time.Time{}, errors.Wrapf(err, "failed to unpack consensus state for client: %s", clientId)
	}
	consensusState, ok := iConsensusState.(*tendermint.ConsensusState)
	if !ok {
		return time.Time{}, errors.Wrapf(err, "invalid consensus state type: %T", iConsensusState)
	}

	if client.warnExpiration(consensusState.Timestamp, warningTime) {
		client.Health = false

		msg := fmt.Sprintf(
			"client %s(%s) would be expired, consensus state timestamp: %s",
			clientId, client.ChainId, consensusState.Timestamp,
		)
		logger.Warn(msg)
		alert.SendTg(msg)

		return time.Time{}, err
	}

	client.Health = true

	logger.Info(fmt.Sprintf("client %s is healthy", clientId))

	return consensusState.Timestamp, nil
}

func (client *Client) warnExpiration(latestTimestamp time.Time, warningTime time.Duration) bool {
	// "expired" means: timestamp + trusting period <= current time
	warnExpirationTime := latestTimestamp.Add(client.TrustingPeriod).Add(-warningTime)
	return !warnExpirationTime.After(time.Now().UTC())
}
