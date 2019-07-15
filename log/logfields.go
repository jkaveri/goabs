package log

import (
	"fmt"
	"strings"
)

const (
	fieldLevel      = "level"
	fieldMessage    = "message"
	fieldError      = "error"
	fieldFormatArgs = "@@fargs"
)

// Fields fields contains all logging data
type Fields map[string]interface{}

// Level extract level in fields
// return `LevelNone` if there are no level was set or level type
func (d Fields) Level() Level {
	if val, ok := d[fieldLevel]; ok {
		if l, ok := val.(Level); ok {
			return l
		}
	}
	return LevelNone
}

// Message extract log message in fields
// return empty string if log message haven't set yet.
// return formatted string (like `fmt.Sprintf`) if the message is format string
// and user has attached the format arguments by using `WithFormatArgs` method
func (d Fields) Message() string {
	val, ok := d[fieldMessage]
	if !ok {
		return ""
	}
	s, ok := val.(string)
	if !ok {
		return ""
	}
	val, ok = d[fieldFormatArgs]
	if !ok {
		return s
	}

	if args, ok := val.([]interface{}); ok {
		return fmt.Sprintf(s, args...)
	}

	return s
}

// Error extract error from Fields
// return nil if not exist
func (d Fields) Error() error {
	if val, ok := d[fieldError]; ok {
		if err, ok := val.(error); ok {
			return err
		}
	}
	return nil
}

// Rest return all log data in Fields except "known-fields"
// known fields are: Error, Message, Level, FormatArgs
func (d Fields) Rest() map[string]interface{} {
	rest := make(map[string]interface{})
	for key, val := range d {
		if d.IsKnowFieldKey(key) {
			continue
		}
		rest[key] = val
	}
	return rest
}

// String serialize the Fields to a string
// return a string in this format:
// ```text
// [level="%s"][ message="%s"][ error="%s"][...[ key="%s"]
// ```
// log message will use the `fields.Message()` method.
func (d Fields) String() string {
	var sb strings.Builder
	str := d.Level().String()
	first := true
	if str != "" {
		first = false
		writeLogItem(&sb, fieldLevel, str)
	}

	str = d.Message()
	if str != "" {
		if !first {
			// nolint
			sb.WriteRune(' ')
		}
		first = false
		writeLogItem(&sb, fieldMessage, str)
	}

	if err := d.Error(); err != nil {
		if !first {
			// nolint
			sb.WriteRune(' ')
		}
		first = false
		writeLogItem(&sb, fieldError, err.Error())
	}
	for key, val := range d {
		// ignore known fields
		if d.IsKnowFieldKey(key) {
			continue
		}
		if !first {
			// nolint
			sb.WriteRune(' ')
		}
		first = false
		writeLogItem(&sb, key, getString(val))
	}
	return sb.String()
}

// IsKnowFieldKey check if
func (d Fields) IsKnowFieldKey(key string) bool {
	return key == fieldFormatArgs ||
		key == fieldMessage ||
		key == fieldError ||
		key == fieldLevel
}
