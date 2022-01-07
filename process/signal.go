package process

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/Masterminds/log-go"
	"github.com/ajgon/mailbowl/config"
)

const (
	ExitCodeGeneral    = 1
	ExitCodeTerminated = 130
)

func AttachSignals() (chan os.Signal, chan os.Signal) {
	reloadChan := make(chan os.Signal, 1)
	interruptChan := make(chan os.Signal, 1)

	signal.Notify(reloadChan, syscall.SIGHUP, syscall.SIGUSR1, syscall.SIGUSR2)
	signal.Notify(interruptChan, syscall.SIGINT, syscall.SIGQUIT)

	return reloadChan, interruptChan
}

func HandleReload(ctx context.Context, reloadChan chan os.Signal) {
	for {
		select {
		case <-reloadChan:
			log.Info("reloading config")
			config.Reload()
		case <-ctx.Done():
			return
		}
	}
}

func HandleInterrupt(ctx context.Context, cancel func(), interruptChan chan os.Signal) {
	select {
	case <-interruptChan:
		log.Info("gracefully shutting down")
		cancel()
	case <-ctx.Done():
		return
	}
	<-interruptChan
	log.Warn("forcing quit")
	os.Exit(ExitCodeTerminated)
}

func Cleanup(cancel func(), reloadChan, interruptChan chan os.Signal) {
	signal.Stop(reloadChan)
	signal.Stop(interruptChan)
	cancel()
}
