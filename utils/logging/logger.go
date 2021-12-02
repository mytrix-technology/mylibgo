package log

import (
	"fmt"
	"io"
)

type Option func(Logger)

var errMissingValue = fmt.Errorf("(MISSING)")

type Logger interface {
	Log(keyvals ...interface{}) error
	SetOutput(io.Writer)
}

// WithOutput configures the logger output to w io.Writer
func WithOutput(w io.Writer) Option {
	return func(l Logger) {
		l.SetOutput(w)
	}
}

// // WithLevel configures a logrus logger to log at level for all events.
// func WithLevel(level apexlog.Level) Option {
// 	return func(c *cliLogger) {
// 		c.level = level
// 	}
// }
