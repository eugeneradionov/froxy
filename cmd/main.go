package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/eugeneradionov/froxy/config"
	"github.com/eugeneradionov/froxy/pkg/http"
	"github.com/eugeneradionov/froxy/pkg/logger"
	"github.com/eugeneradionov/froxy/pkg/validator"
	"github.com/eugeneradionov/froxy/services/proxy"
	prtr "github.com/eugeneradionov/froxy/services/proxy/transport/http"
	"github.com/eugeneradionov/froxy/store/inmemory"
	"go.uber.org/zap"
)

func main() {
	err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %s", err.Error())
	}

	err = logger.Load(config.Get().Logger.Preset)
	if err != nil {
		log.Fatalf("Failed to load logger: %s", err.Error())
	}

	err = validator.Load()
	if err != nil {
		logger.Get().Fatal("Failed to load validator", zap.Error(err))
	}

	srv := http.NewServer(config.Get().HTTPServer, logger.Get())

	store := inmemory.NewStore()

	proxySrv := proxy.New(config.Get().Proxy, store)
	proxyTr := prtr.New(logger.Get(), proxySrv)

	srv.MountRoutes(proxyTr)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-signalChannel
		logger.Get().Info("shutdown signal received! Bye!", zap.String("signal", sig.String()))
		cancel()
	}()

	logger.Get().Info("Starting HTTP server", zap.String("listen_url", config.Get().HTTPServer.ListenURL))
	err = srv.Serve(ctx)
	if err != nil {
		logger.Get().Error("Failed to initialize HTTP server", zap.Error(err))
		os.Exit(1) // nolint:gocritic
	}
}
