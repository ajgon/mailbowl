package listener

import "context"

type Listener interface {
	GetName() string
	Serve(context.Context) error
}
