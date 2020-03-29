package logging

import (
	"io/ioutil"
	"log"
	"os"
)

var (
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
)

func init() {
	hasDir := false

	allFiles, err := ioutil.ReadDir("./")
	if err != nil {
		log.Fatalln("讀取目錄失敗：", err)
	}

	for _, fi := range allFiles {
		if fi.Name() == "logs" {
			hasDir = true
		}
	}

	if !hasDir {
		err := os.Mkdir("logs", os.ModeDir)
		if err != nil {
			log.Fatalln("創建logs目錄失敗：", err)
		}
	}

	errFile, err := os.OpenFile(".\\logs\\errors.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("開啟log文件失敗：", err)
	}
	Info = log.New(os.Stdout, "[INFO] ", log.LstdFlags|log.Lshortfile|log.LUTC)
	Warning = log.New(os.Stdout, "[WARNING] ", log.LstdFlags|log.Lshortfile|log.LUTC)
	Error = log.New(errFile, "[ERROR] ", log.LstdFlags|log.Lshortfile|log.LUTC)

}
