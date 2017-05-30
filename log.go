package wrapplog

import (
	"os"
	"path"

	"github.com/Sirupsen/logrus"
)

type WrappHook struct {
	ServiceName string
}

func (w *WrappHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (w *WrappHook) Fire(entry *logrus.Entry) error {
	entry.Data["service"] = w.ServiceName

	return nil
}

func init() {
	serviceName := os.Getenv("SERVICE_NAME")
	if serviceName == "" {
		// Fallback to executable name when
		// service name isn't specified
		exePath, err := os.Executable()
		if err != nil {
			serviceName = "unknown"
		} else {
			serviceName = path.Base(exePath)
		}
	}

	logrus.AddHook(&WrappHook{ServiceName: serviceName})
	logrus.SetFormatter(&logrus.JSONFormatter{
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyLevel: "level",
			logrus.FieldKeyMsg:   "msg",
			logrus.FieldKeyTime:  "timestamp",
		},
	})
	logrus.SetOutput(os.Stdout)
}
