package app

import (
	"context"
	"fmt"
	"time"

	"github.com/dlvlabs/ibcmon/logger"
	"github.com/pkg/errors"
)

func (app *App) Run(ctx context.Context) error {
	// app.initIBCInfo: run every cfg.General.IbcInfoUpdateInterval.

	// app.checkClientsHealth: run every cfg.General.ClientCheckInterval,
	// app.initIBCInfo should be done before this function.

	// app.trackIBCPacket: this function would be run continuously
	// app.initIBCInfo should be done before this function.

	for {
		appCtx, cancel := context.WithCancel(ctx)

		err := app.initIBCInfo(appCtx)
		if err != nil {
			return err
		}

		go func() {
			// Initialize ticker to fire immediately
			ticker := time.NewTicker(1 * time.Second)
			defer ticker.Stop()

			msg := fmt.Sprintf("check client health: %s", context.Canceled.Error())

			for {
				select {
				case <-ticker.C:
					err = app.checkClientsHealth(appCtx)
					if err != nil {
						if errors.Is(err, context.Canceled) {
							logger.Info(msg)
							return
						}

						logger.Error(err)
						panic(err)
					}

					// reset ticket
					ticker.Reset(app.cfg.General.ClientCheckInterval)
				case <-appCtx.Done():
					logger.Info(msg)
					return
				}
			}
		}()

		go func() {
			err = app.trackIBCPacket(appCtx)
			if err != nil {
				if errors.Is(err, context.Canceled) {
					msg := fmt.Sprintf("track ibc packet: %s", context.Canceled.Error())
					logger.Info(msg)
					return
				}

				logger.Error(err)
				panic(err)
			}
		}()

		time.Sleep(app.cfg.General.IbcInfoUpdateInterval)

		cancel()
	}
}
