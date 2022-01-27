package relay

import (
	"fmt"

	"github.com/ajgon/mailbowl/config"
)

type Relay struct {
	OutgoingServer *OutgoingServer
}

func NewRelay(conf config.Relay) (*Relay, error) {
	outgoingServer, err := NewOutgoingServer(conf.OutgoingServer)
	if err != nil {
		return nil, fmt.Errorf("error configuring outgoing server: %w", err)
	}

	return &Relay{
		OutgoingServer: outgoingServer,
	}, nil
}

func (r *Relay) Handle(from string, recipients []string, message []byte) error {
	err := r.OutgoingServer.Send(from, recipients, message)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}
