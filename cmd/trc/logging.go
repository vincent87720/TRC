package main

import (
	"log"
	"os"
)

var (
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
)

func initLogging() {

	errFile, err := os.OpenFile(".\\assets\\logs\\errors.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("開啟log文件失敗：", err)
	}
	Info = log.New(os.Stdout, "[INFO] ", log.LstdFlags|log.Lshortfile|log.LUTC)
	Warning = log.New(os.Stdout, "[WARNING] ", log.LstdFlags|log.Lshortfile|log.LUTC)
	Error = log.New(errFile, "[ERROR] ", log.LstdFlags|log.Lshortfile|log.LUTC)

}
