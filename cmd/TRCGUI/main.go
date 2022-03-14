//go:generate goversioninfo -manifest=../../tools/goversioninfo/goversioninfo.exe.manifest

package main

import (
	file "github.com/vincent87720/TRC/internal/file"
	"github.com/vincent87720/TRC/internal/gui"
	"github.com/vincent87720/TRC/internal/logging"
)

var (
	INITPATH string
)

func init() {
	INITPATH = file.GetInitialPath()
	logging.InitLogging(INITPATH)
}

func main() {
	gui.StartWindow(INITPATH)
}
