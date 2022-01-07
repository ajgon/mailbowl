package process

import (
	"context"

	"github.com/Masterminds/log-go"
)

func Main(run func(context.Context) error) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	reloadChan, interruptChan := AttachSignals()

	defer Cleanup(cancel, reloadChan, interruptChan)

	go HandleReload(ctx, reloadChan)
	go HandleInterrupt(ctx, cancel, interruptChan)

	if err := run(ctx); err != nil {
		log.Errorf("unprocessable error: %s", err.Error())
	}
}
