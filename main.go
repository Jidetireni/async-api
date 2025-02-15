// main.go
package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/Jidetireni/async-api.git/api"
	"github.com/Jidetireni/async-api.git/config"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	conf, err := config.NewConfig()
	if err != nil {
		return err
	}

	jsonHandler := slog.NewJSONHandler(os.Stdout, nil)
	logger := slog.New(jsonHandler)
	apiServer := api.NewApi(conf, logger)

	if err := apiServer.Start(ctx); err != nil {
		return err
	}
	return nil

}
