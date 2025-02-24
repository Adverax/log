//go:build windows
// +build windows

package logFileRotator

import (
	"os"
)

func chown(_ string, _ os.FileInfo) error {
	return nil
}
