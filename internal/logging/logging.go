package logging

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

//Use"Info"and"Warning"to print a message on stdout, and use "Error" to log the message in errors.log
var (
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
)

func InitLogging(INITPATH string) {
	if _, err := os.Stat(INITPATH + "/assets/logs"); os.IsNotExist(err) {
		err := os.MkdirAll(filepath.FromSlash(INITPATH+"/assets/logs"), os.ModeDir)
		if err != nil {
			fmt.Println(err)
		}
	}

	errFile, err := os.OpenFile(filepath.FromSlash(INITPATH+"/assets/logs/errors.log"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("開啟log文件失敗：", err)
	}
	Info = log.New(os.Stdout, "[INFO] ", log.LstdFlags|log.Lshortfile|log.LUTC)
	Warning = log.New(os.Stdout, "[WARNING] ", log.LstdFlags|log.Lshortfile|log.LUTC)
	Error = log.New(errFile, "[ERROR] ", log.LstdFlags|log.Lshortfile|log.LUTC)

}
