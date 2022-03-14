// +build windows

package gui

var (
	INITPATH string
)

func StartWindow(path string) {
	INITPATH = path
	exportAssets()
	runMainWindow()
}
