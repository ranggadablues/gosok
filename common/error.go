package common

import (
	"fmt"
	"os"
	"runtime"
)

func IsLocalDev() bool {
	return os.Getenv("APP_ENV") == "local"
}

func Error(msg string, err error) error {
	if err == nil && msg == "" {
		return nil
	}

	if IsLocalDev() {
		// capture file + line number
		_, file, line, _ := runtime.Caller(1)
		if msg == "" {
			return fmt.Errorf("%w (at %s:%d)", err, file, line)
		}

		return fmt.Errorf("%s: %w (at %s:%d)", msg, err, file, line)
	}

	// production: hide internals
	return err
}

func ErrorInsert(err error) error {
	return Error("failed to insert", err)
}

func ErrorFind(err error) error {
	return Error("failed to find", err)
}

func ErrorUpdate(err error) error {
	return Error("failed to update", err)
}

func ErrorDelete(err error) error {
	return Error("failed to delete", err)
}
