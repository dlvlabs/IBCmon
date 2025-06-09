package app

import (
	"time"
)

func (app *App) connectGRPCs() error {
	app.grpcsMutex.Lock()

	for _, grpc := range app.grpcs {
		err := grpc.Connect()
		if err != nil {
			return err
		}
	}
	return nil
}

func (app *App) terminateGRPCs() error {
	for _, grpc := range app.grpcs {
		err := grpc.Terminate()
		if err != nil {
			return err
		}
	}

	app.grpcsMutex.Unlock()

	return nil
}

func (app *App) updateStore(update func()) {
	app.storeMutex.Lock()
	defer app.storeMutex.Unlock()

	update()
	app.Store.Updated = time.Now().UTC()
}
