// +build !windows

package main

func startWindow() {
	Warning.Printf("%+v\n", "Can't start gui. Please make sure your operating system is Windows.")
}
