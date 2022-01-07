package smtp

import (
	"time"

	"github.com/Masterminds/log-go"
	"github.com/ajgon/mailbowl/config"
)

type Timeout struct {
	Read  time.Duration
	Write time.Duration
	Data  time.Duration
}

func NewTimeout(read, write, data string) *Timeout {
	var err error

	defaults := config.GetDefaults()

	readDuration, err := time.ParseDuration(read)
	if err != nil {
		log.Debugf("invalid smtp.timeout.read value `%s`, setting default `%s`", read, defaults.SMTPTimeoutRead)

		readDuration, _ = time.ParseDuration(defaults.SMTPTimeoutRead)
	}

	writeDuration, err := time.ParseDuration(write)
	if err != nil {
		log.Debugf("invalid smtp.timeout.write value `%s`, setting default `%s`", write, defaults.SMTPTimeoutWrite)

		writeDuration, _ = time.ParseDuration(defaults.SMTPTimeoutWrite)
	}

	dataDuration, err := time.ParseDuration(data)
	if err != nil {
		log.Debugf("invalid smtp.timeout.data value `%s`, setting default `%s`", data, defaults.SMTPTimeoutData)

		dataDuration, _ = time.ParseDuration(defaults.SMTPTimeoutData)
	}

	return &Timeout{
		Read:  readDuration,
		Write: writeDuration,
		Data:  dataDuration,
	}
}
