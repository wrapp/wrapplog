package wrapplog

import (
	"strings"

	"github.com/Sirupsen/logrus"
)

var jsonFormatter = logrus.JSONFormatter{}

type WrappFormatter struct {
	// TimestampFormat sets the format used for marshaling timestamps.
	TimestampFormat string
}

func (f *WrappFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	jsonBytes, err := (&jsonFormatter).Format(entry)
	prefix := []byte(strings.ToUpper(entry.Level.String()) + " ")
	return append(prefix[:], jsonBytes[:]...), err
}
