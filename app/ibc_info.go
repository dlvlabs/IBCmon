package app

import (
	"context"
	"fmt"
	"time"

	"github.com/dlvlabs/ibcmon/logger"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

type (
	// chainid => Clients
	IBCInfo map[string]Clients

	// clientId => Client
	Clients map[string]*Client

	Client struct {
		Health bool

		ChainId        string
		RevisionNumber uint64
		RevisionHeight uint64
		TrustingPeriod time.Duration

		// this value updated by app.checkClientHealth
		ClientUpdated time.Time

		Connections Connections
	}

	// connectionId => Channels
	Connections map[string]Channels

	// channelId => Channel
	Channels map[string]*Channel
	Channel  struct {
		PortId       string
		Counterparty *Counterparty

		// this value updated by app.trackIBCPacket
		IBCPacketTracker *IBCPacketTracker
	}
	Counterparty struct {
		ClientId     string
		ConnectionId string
		ChannelId    string
		PortId       string
	}
)

func (app *App) initIBCInfo(ctx context.Context) error {
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

	err = app.setBaseChain(ctx)
	if err != nil {
		return err
	}

	err = app.setCounterparties(ctx)
	if err != nil {
		return err
	}

	logger.Debug(fmt.Sprintf("IBCInfo: %s", app.Store.IBCInfo))

	return nil
}

func (app *App) setBaseChain(ctx context.Context) error {
	msg := fmt.Sprintf("init ibc info for basechain(%s)", app.cfg.General.baseChainId)
	logger.Info(msg)

	clients := make(Clients)
	err := clients.setActiveClients(ctx, app.grpcs[app.cfg.General.baseChainId], app.cdc)
	if err != nil {
		return err
	}
	app.updateStore(func() { app.Store.IBCInfo[app.cfg.General.baseChainId] = clients })

	return nil
}

func (app *App) setCounterparties(ctx context.Context) error {
	g, ctx := errgroup.WithContext(ctx)

	for _, client := range app.Store.IBCInfo[app.cfg.General.baseChainId] {
		chainId := client.ChainId

		// check whether the endpoint is in the config file
		_, ok := app.cfg.Counterparties[chainId]
		if !ok {
			msg := fmt.Sprintf("missing counterparty endpoints in config file for %s", chainId)
			return errors.New(msg)
		}

		for _, channels := range client.Connections {
			for _, channel := range channels {
				g.Go(func() error {
					msg := fmt.Sprintf("init ibc info for counterparty(%s)", chainId)
					logger.Info(msg)

					clients := make(Clients)
					err := clients.setActiveClient(ctx, app.grpcs[chainId], app.cdc, channel.Counterparty.ClientId)
					if err != nil {
						logger.Error(err)

						return err
					}
					app.updateStore(func() { app.Store.IBCInfo[chainId] = clients })

					return nil
				})
			}
		}
	}

	return g.Wait()
}
