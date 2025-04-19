package main

import (
	"context"
	"flag"

	"github.com/bovinxx/auth-service/internal/app"
	"github.com/bovinxx/auth-service/internal/logger"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config-path", ".env", "path to config file")
}

func main() {
	logger.InitDefaulLogger()

	flag.Parse()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	app, err := app.NewApp(ctx)
	if err != nil {
		panic(err)
	}

	if err := app.Start(ctx); err != nil {
		panic(err)
	}
}
