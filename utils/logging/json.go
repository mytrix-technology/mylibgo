package log

import (
	"io"
	"os"

	"github.com/go-kit/kit/log"
)

type jsonLogger struct {
	log.Logger
}

// NewJsonLogger returns a Logger with plain cli format
func NewJsonLogger(options ...Option) Logger {
	logger := &jsonLogger{log.NewJSONLogger(os.Stdout)}
	for _, optFunc := range options {
		optFunc(logger)
	}

	return logger
}

func (j jsonLogger) Log(keyvals ...interface{}) error {
	return j.Logger.Log(keyvals...)
}

func (j *jsonLogger) SetOutput(w io.Writer) {
	logger := log.NewJSONLogger(w)
	j.Logger = logger
}
