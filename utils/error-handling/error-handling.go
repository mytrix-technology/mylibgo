package error_handling

import (
	"encoding/json"
	"fmt"
	"runtime"
)

type ErrorString struct {
	code       int
	stacktrace string
	messages   []string
	status     string
}

func (e ErrorString) Code() int {
	return e.code
}

func (e ErrorString) Error() string {
	return e.stacktrace
}

func (e ErrorString) Messages() []string {
	return e.messages
}

func (e ErrorString) Status() string {
	return e.status
}

func NewError(cause error, vals ...interface{}) error {
	if cause == nil {
		return nil
	}

	_, ok := cause.(*ErrorString)
	if ok {
		return cause
	}

	j, _ := json.Marshal(vals)
	stacktrace := fmt.Sprintf(" Message %s %s", cause.Error(), j)

	for i := 1; i <= 3; i++ {
		pc, file, line, _ := runtime.Caller(i)
		f := runtime.FuncForPC(pc)
		if f == nil || line == 0 {
			break
		}

		stacktrace += fmt.Sprintf(" --- at %s:%d ---", file, line)
	}

	return &ErrorString{500, stacktrace, []string{cause.Error()}, "failed"}
}

func NewValidationError(text string, vals ...interface{}) error {
	j, _ := json.Marshal(vals)
	stacktrace := fmt.Sprintf("\nMessage: %s %s", text, j)

	for i := 1; i <= 3; i++ {
		pc, file, line, _ := runtime.Caller(i)
		f := runtime.FuncForPC(pc)
		if f == nil || line == 0 {
			break
		}

		stacktrace += fmt.Sprintf("\n--- at %s:%d ---", file, line)
	}

	return &ErrorString{400, stacktrace, []string{text}, "failed"}
}

func NewValidationErrors(text []string, vals ...interface{}) error {
	j, _ := json.Marshal(vals)
	stacktrace := fmt.Sprintf(" Message: %s %s", text, j)

	for i := 1; i <= 3; i++ {
		pc, file, line, _ := runtime.Caller(i)
		f := runtime.FuncForPC(pc)
		if f == nil || line == 0 {
			break
		}

		stacktrace += fmt.Sprintf(" --- at %s:%d ---", file, line)
	}

	return &ErrorString{400, stacktrace, text, "failed"}
}

func NewNotFoundError(text string, vals ...interface{}) error {
	j, _ := json.Marshal(vals)
	stacktrace := fmt.Sprintf("\nMessage: %s %s", text, j)

	for i := 1; i <= 3; i++ {
		pc, file, line, _ := runtime.Caller(i)
		f := runtime.FuncForPC(pc)
		if f == nil || line == 0 {
			break
		}

		stacktrace += fmt.Sprintf("\n--- at %s:%d ---", file, line)
	}

	return &ErrorString{404, stacktrace, []string{text}, "failed"}
}
