package wrapplog

// Rip-off of the logrus TextFormatter with Wrapp defaults
import (
	"bytes"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
)

func prefixFieldClashes(entry *logrus.Entry) {
	_, ok := entry.Data["time"]
	if ok {
		entry.Data["fields.time"] = entry.Data["time"]
	}
	_, ok = entry.Data["msg"]
	if ok {
		entry.Data["fields.msg"] = entry.Data["msg"]
	}
	_, ok = entry.Data["level"]
	if ok {
		entry.Data["fields.level"] = entry.Data["level"]
	}
}

const (
	nocolor = 0
	red     = 31
	green   = 32
	yellow  = 33
	blue    = 34
)

var (
	baseTimestamp time.Time
	isTerminal    bool
)

func init() {
	baseTimestamp = time.Now()
	isTerminal = logrus.IsTerminal()
}
func miniTS() int {
	return int(time.Since(baseTimestamp) / time.Second)
}

// Formatter formats logs according to the Wrapp standard
// Uses nice colored output if the TTY is a terminal
type Formatter struct {
	// Set to true to bypass checking for a TTY before outputting colors.
	ForceColors   bool
	DisableColors bool
}

func (f *Formatter) Format(entry *logrus.Entry) ([]byte, error) {
	var keys []string
	for k := range entry.Data {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	b := &bytes.Buffer{}
	prefixFieldClashes(entry)
	isColored := (f.ForceColors || isTerminal) && !f.DisableColors
	if isColored {
		printColored(b, entry, keys)
	} else {
		f.appendKeyValue(b, "time", entry.Time.Format(time.RFC3339))
		f.appendKeyValue(b, "level", entry.Level.String())
		f.appendKeyValue(b, "msg", entry.Message)
		for _, key := range keys {
			f.appendKeyValue(b, key, entry.Data[key])
		}
	}
	b.WriteByte('\n')
	return b.Bytes(), nil
}
func printColored(b *bytes.Buffer, entry *logrus.Entry, keys []string) {
	var levelColor int
	switch entry.Level {
	case logrus.WarnLevel:
		levelColor = yellow
	case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
		levelColor = red
	default:
		levelColor = blue
	}
	levelText := strings.ToUpper(entry.Level.String())
	fmt.Fprintf(b, "\x1b[%dm%-5s\x1b[0m [%04d] %-44s ", levelColor, levelText, miniTS(), entry.Message)
	for _, k := range keys {
		v := entry.Data[k]
		fmt.Fprintf(b, " \x1b[%dm%s\x1b[0m=%v", levelColor, k, v)
	}
}
func (f Formatter) appendKeyValue(b *bytes.Buffer, key, value interface{}) {
	switch value.(type) {
	case string, error:
		fmt.Fprintf(b, "%v=%q ", key, value)
	default:
		fmt.Fprintf(b, "%v=%v ", key, value)
	}
}
