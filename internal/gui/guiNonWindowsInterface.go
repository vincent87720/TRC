// +build !windows

package gui

import "github.com/vincent87720/TRC/internal/logging"

func StartWindow(path string) {
	logging.Warning.Printf("%+v\n", "Can't start gui. Please make sure your operating system is Windows.")
}
