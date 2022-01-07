package config

import (
	"os"

	"github.com/Masterminds/log-go"
)

func Reload() {
	err := ConfigureLogger(os.Stdout, os.Stderr)
	handleReloadError(err)
}

func handleReloadError(err error) {
	if err == nil {
		return
	}

	log.Errorf("problem reloading configuration: %s", err.Error())
}
