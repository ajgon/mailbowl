package relay

import (
	"fmt"
)

type Relay struct {
	OutgoingServer *OutgoingServer
}

func NewRelay(outgoingServer *OutgoingServer) *Relay {
	return &Relay{
		OutgoingServer: outgoingServer,
	}
}

func (r *Relay) Handle(from string, recipients []string, message []byte) error {
	err := r.OutgoingServer.Send(from, recipients, message)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}
