package formatter

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"strings"
)

var emojisLevel = [7]string{"💀", "🤬", "💩", "🙈", "🙃", "🤷", "🤮"}
var colors = [7]string{"[44;1m", "[41;1m", "[31;1m", "[33m", "[36m", "[37;1m", "[35;1m"}

const Start = "\033"
const End = "\033[0m"

// Format building log message.
func (f *Config) Format(entry *logrus.Entry) ([]byte, error) {
	output := f.LogFormat
	if output == "" {
		output = defaultLogFormat
	}

	timestampFormat := f.TimestampFormat
	if timestampFormat == "" {
		timestampFormat = defaultTimestampFormat
	}

	output = strings.Replace(output, "%time%", entry.Time.Format(timestampFormat), 1)

	output = strings.Replace(output, "%msg%", entry.Message, 1)

	level := strings.ToUpper(entry.Level.String())

	i, _ := logrus.ParseLevel(level)
	l := level
	if f.Color {
		color := colors[i]
		l = Start + color + level + End
	}

	output = strings.Replace(output, "%lvl%", l, 1)

	emoji := emojisLevel[i]
	output = strings.Replace(output, "%emoji%", emoji, 1)

	for k, val := range entry.Data {
		switch val.(type) {
		case []string:
			v := strings.Join(val.([]string), "\n\t  - ")
			output += fmt.Sprintf("\n\t%s: \n\t  - %v", k, v)
		default:
			output += fmt.Sprintf("\n\t%s: %v", k, val)
		}
		//strings.Join(s, ", "
		//output += fmt.Sprintf("\n\t%s: %v", k, val)
		//output = strings.Replace(output, "%"+k+"%", s, 1)
	}
	output += "\n"
	return []byte(output), nil
}
