package wrapplog

import (
	"os"
	"strings"

	"github.com/Sirupsen/logrus"
)

var jsonFormatter = logrus.JSONFormatter{}

type WrappFormatter struct{}

func init() {
	logrus.SetFormatter(&WrappFormatter{})
	logrus.SetOutput(os.Stdout)
}

// Format logs according to WEP-007
func (f *WrappFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	jsonBytes, err := (&jsonFormatter).Format(entry)
	prefix := []byte(strings.ToUpper(entry.Level.String()) + " ")
	return append(prefix[:], jsonBytes[:]...), err
}
