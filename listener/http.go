package listener

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/Masterminds/log-go"
)

type HTTP struct{}

func NewHTTP() *HTTP {
	return new(HTTP)
}

func (h HTTP) GetName() string {
	return "HTTP"
}

func (h HTTP) Serve(ctx context.Context) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Debugw("/", log.Fields{"path": "/", "status": http.StatusOK})
		fmt.Fprintf(w, "OK")
	})

	server := &http.Server{
		Addr:    ":3000",
		Handler: mux,
	}

	log.Info("HTTP server started on port 3000")

	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("HTTP server died: %s", err.Error())
		}
	}()

	<-ctx.Done()

	log.Debug("stopping HTTP server")

	ctxShutDown, cancel := context.WithTimeout(context.Background(), time.Second)
	defer func() {
		cancel()
	}()

	// false positive (parent already cancelled)
	//nolint:contextcheck
	if err := server.Shutdown(ctxShutDown); err != nil {
		return fmt.Errorf("error stopping HTTP server: %w", err)
	}

	log.Debug("HTTP server shut down")

	return nil
}
