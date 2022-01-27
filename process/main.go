package process

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/Masterminds/log-go"
	"github.com/ajgon/mailbowl/config"
	"github.com/ajgon/mailbowl/listener"
)

const (
	ExitCodeGeneral    = 1
	ExitCodeTerminated = 130
)

type Manager struct {
	listeners []listener.Listener

	interruptChan chan os.Signal
	reloadChan    chan os.Signal
	restarting    bool
}

func NewManager() *Manager {
	manager := new(Manager)
	manager.listeners = make([]listener.Listener, 0)
	manager.interruptChan = make(chan os.Signal, 1)
	manager.reloadChan = make(chan os.Signal, 1)
	manager.restarting = true

	return manager
}

func (m *Manager) AddListener(listener listener.Listener) {
	m.listeners = append(m.listeners, listener)
}

func (m *Manager) Start() {
	for m.restarting {
		m.restarting = false
		m.startAllListeners()
	}
}

func (m *Manager) Restart(cancelCtx func()) {
	m.restarting = true

	cancelCtx()
}

func (m *Manager) startAllListeners() {
	var waitGroup sync.WaitGroup

	ctx, cancelCtx := context.WithCancel(context.Background())

	m.attachSignals()

	defer m.cleanup(cancelCtx)

	go m.handleReload(ctx, cancelCtx)
	go m.handleInterrupt(ctx, cancelCtx)

	for _, listener := range m.listeners {
		waitGroup.Add(1)

		go m.startListener(ctx, cancelCtx, &waitGroup, listener)
	}

	waitGroup.Wait()
}

func (m *Manager) attachSignals() {
	signal.Notify(m.reloadChan, syscall.SIGHUP, syscall.SIGUSR1, syscall.SIGUSR2)
	signal.Notify(m.interruptChan, syscall.SIGINT, syscall.SIGQUIT)
}

func (m *Manager) cleanup(cancelCtx func()) {
	signal.Stop(m.reloadChan)
	signal.Stop(m.interruptChan)
	cancelCtx()
}

func (m *Manager) handleReload(ctx context.Context, cancelCtx func()) {
	for {
		select {
		case <-m.reloadChan:
			log.Info("reloading config")
			config.Reload()
			m.Restart(cancelCtx)
		case <-ctx.Done():
			return
		}
	}
}

func (m *Manager) handleInterrupt(ctx context.Context, cancelCtx func()) {
	select {
	case <-m.interruptChan:
		log.Info("gracefully shutting down")
		cancelCtx()
	case <-ctx.Done():
		return
	}
	<-m.interruptChan
	log.Warn("forcing quit")
	os.Exit(ExitCodeTerminated)
}

func (m *Manager) startListener(ctx context.Context, cancelCtx func(), wg *sync.WaitGroup, listener listener.Listener) {
	defer wg.Done()

	if err := listener.Serve(ctx); err != nil {
		log.Errorf("unprocessable %s error: %s", listener.GetName(), err.Error())
		cancelCtx()
	}
}
