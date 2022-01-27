package smtp

import (
	"time"

	"github.com/ajgon/mailbowl/config"
)

type Timeout struct {
	Read  time.Duration
	Write time.Duration
	Data  time.Duration
}

func NewTimeout(conf config.SMTPTimeout) *Timeout {
	return &Timeout{
		Read:  conf.Read,
		Write: conf.Write,
		Data:  conf.Data,
	}
}
