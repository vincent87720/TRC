package main

import (
	"log"
	"os"
	"path/filepath"
)

var (
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
)

func initLogging() {
	err := os.MkdirAll(filepath.FromSlash(INITPATH+"/assets/logs"), os.ModeDir)
	if err != nil {
		Error.Printf("%+v\n", err)
	}

	errFile, err := os.OpenFile(filepath.FromSlash(INITPATH+"/assets/logs/errors.log"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("開啟log文件失敗：", err)
	}
	Info = log.New(os.Stdout, "[INFO] ", log.LstdFlags|log.Lshortfile|log.LUTC)
	Warning = log.New(os.Stdout, "[WARNING] ", log.LstdFlags|log.Lshortfile|log.LUTC)
	Error = log.New(errFile, "[ERROR] ", log.LstdFlags|log.Lshortfile|log.LUTC)

}
