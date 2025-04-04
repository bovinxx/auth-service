package main

import (
	"context"
	"flag"

	_ "net/http/pprof"

	"github.com/bovinxx/auth-service/internal/app"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config-path", ".env", "path to config file")
}

func main() {
	// go func() {
	// 	log.Println(http.ListenAndServe("localhost:6060", nil))
	// }()

	flag.Parse()
	ctx := context.Background()

	app, err := app.NewApp(ctx)
	if err != nil {
		panic(err)
	}

	if err := app.Run(); err != nil {
		panic(err)
	}
}
