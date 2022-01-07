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

func NewLimit(connections, messageSize, recipients int) *Limit {
	defaults := config.GetDefaults()

	if connections == 0 { // less than zero disables limiter completely
		log.Debugf(
			"invalid smtp.limit.connections value `%d`, setting default `%d`", connections, defaults.SMTPLimitConnections,
		)

		connections = defaults.SMTPLimitConnections
	}

	if messageSize <= 0 {
		log.Debugf(
			"invalid smtp.limit.message_size value `%d`, setting default `%d`", messageSize, defaults.SMTPLimitMessageSize,
		)

		messageSize = defaults.SMTPLimitMessageSize
	}

	if recipients <= 0 {
		log.Debugf(
			"invalid smtp.limit.recipients value `%d`, setting default `%d`", recipients, defaults.SMTPLimitRecipients,
		)

		recipients = defaults.SMTPLimitRecipients
	}

	return &Limit{
		Connections: connections,
		MessageSize: messageSize,
		Recipients:  recipients,
	}
}
