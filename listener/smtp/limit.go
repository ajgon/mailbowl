package smtp

import (
	"github.com/Masterminds/log-go"
	"github.com/ajgon/mailbowl/config"
)

type Limit struct {
	Connections int
	MessageSize int
	Recipients  int
}

func NewLimit(conf config.SMTPLimit) *Limit {
	var (
		oldValue    int
		connections = conf.Connections
		messageSize = conf.MessageSize
		recipients  = conf.Recipients
	)

	if connections == 0 { // less than zero disables limiter completely
		oldValue = connections
		connections = config.GetDefaultInt("smtp.limit.connections")

		log.Debugf(
			"invalid smtp.limit.connections value `%d`, setting default `%d`", oldValue, connections,
		)
	}

	if messageSize <= 0 {
		oldValue = messageSize
		messageSize = config.GetDefaultInt("smtp.limit.message_size")

		log.Debugf(
			"invalid smtp.limit.message_size value `%d`, setting default `%d`", oldValue, messageSize,
		)
	}

	if recipients <= 0 {
		oldValue = recipients
		recipients = config.GetDefaultInt("smtp.limit.recipients")

		log.Debugf(
			"invalid smtp.limit.recipients value `%d`, setting default `%d`", oldValue, messageSize,
		)
	}

	return &Limit{
		Connections: connections,
		MessageSize: messageSize,
		Recipients:  recipients,
	}
}
