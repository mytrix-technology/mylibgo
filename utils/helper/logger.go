package helper

type DebugFieldLogger func(keyvals ...interface{}) error

type DebugLogger func(format string, args ...interface{}) error

func CreateNoopFieldLogger() DebugFieldLogger {
	return func(keyvals ...interface{}) error {
		return nil
	}
}
