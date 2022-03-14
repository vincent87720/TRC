package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/user"
	"strings"

	"github.com/fatih/color"
	"github.com/vincent87720/TRC/internal/command"
	file "github.com/vincent87720/TRC/internal/file"
	"github.com/vincent87720/TRC/internal/logging"
)

var (
	INITPATH string
)

//選擇使用模式，0: argsMode 1: manualMode
func selectMode() (mode int) {
	if len(os.Args) > 1 {
		return 0
	} else {
		return 1
	}
}

func autoMode(trcCmd *command.CommandSet, f *command.Flags) {
	flag.Parse()
	command.AnalyseCommand(trcCmd, f)
}

func manualMode(trcCmd *command.CommandSet, f *command.Flags) (err error) {

	usr, err := user.Current()
	if err != nil {
		return err
	}

	path, err := os.Getwd()
	if err != nil {
		return err
	}

	for {
		color.Set(color.FgHiRed)
		fmt.Print("\n", usr.Username)
		color.Set(color.FgHiCyan)
		fmt.Print(" TRC ")
		color.Set(color.FgWhite)
		fmt.Print(path, "\r\n")
		fmt.Print(">")

		reader := bufio.NewReader(os.Stdin)
		data, err := reader.ReadString('\n')
		if err != nil {
			return err
		}
		dataXi := strings.Fields(data)
		os.Args = os.Args[:0]
		for _, val := range dataXi {
			os.Args = append(os.Args, val)
		}
		err = f.InitFlag()
		if err != nil {
			return err
		}
		flag.Parse()
		command.AnalyseCommand(trcCmd, f)
	}
}

func init() {
	INITPATH = file.GetInitialPath()
	logging.InitLogging(INITPATH)
}

func main() {
	f := command.NewFlag()
	err := f.InitFlag()
	if err != nil {
		logging.Error.Printf("%+v\n", err)
	}

	var trcCmd command.CommandSet
	trcCmd.CommandInit(&f, INITPATH)

	mode := selectMode()
	if mode == 0 {
		autoMode(&trcCmd, &f)
	} else if mode == 1 {
		err = manualMode(&trcCmd, &f)
		if err != nil {
			logging.Error.Printf("%+v\n", err)
		}
	}
}
