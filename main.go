package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"runtime"

	"github.com/dlvlabs/ibcmon/alert"
	"github.com/dlvlabs/ibcmon/app"
	"github.com/dlvlabs/ibcmon/logger"
	"github.com/dlvlabs/ibcmon/server"

	"github.com/BurntSushi/toml"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	ctx := context.Background()

	cfgPath := flag.String("config", "", "Config file")
	flag.Parse()
	if *cfgPath == "" {
		panic("Error: Please input config file path with -config flag.")
	}

	f, err := os.ReadFile(*cfgPath)
	if err != nil {
		panic(err)
	}
	cfg := app.Config{}
	err = toml.Unmarshal(f, &cfg)
	if err != nil {
		panic(err)
	}

	if cfg.General.LogLevel == "production" {
		logger.InitLogger(true)
	} else {
		logger.InitLogger(false)
	}

	title := "ibcmon"
	tgTitle := fmt.Sprintf("🤖 %s 🤖", title)
	alert.SetTg(cfg.TG.Enable, tgTitle, cfg.TG.Token, cfg.TG.ChatID)

	app, error := app.NewApp(ctx, cfg)
	if error != nil {
		panic(error)
	}

	server := server.NewServer(&app.Store, cfg.General.ListenPort, title)
	go func() {
		if err := server.Run(); err != nil {
			panic(err)
		}
	}()

	err = app.Run(ctx)
	if err != nil {
		logger.Error(err)
		return
	}
}
