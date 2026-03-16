package cmd

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"team_service/global"
	"time"
)

func Bootstrap() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	deps := global.NewGlobalDependency()
	if err := deps.Start(ctx); err != nil {
		panic(err)
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	<-sig
	cancel()

	shutdownCtx, stop := context.WithTimeout(context.Background(), 10*time.Second)
	defer stop()

	done := make(chan struct{})

	go func() {
		defer close(done)
		_ = deps.Stop(shutdownCtx)
	}()

	go func() {
		<-sig
		os.Exit(1)
	}()
}
