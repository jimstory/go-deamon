package lib

import (
"log"
"time"
"os"
"fmt"
)

var Trace * log.Logger

func init()  {
	date := time.Now().Format("2006-01-02")
	logPath := GetConfigure("LOG_PATH")
	debuggerFile, err := os.OpenFile("./"+logPath+"/" + date + ".log" , os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println(err)
		log.Fatalln("fail  to create debug log")
	}
	//defer debuggerFile.Close()
	Trace = log.New(debuggerFile,"TRACE: ",log.Ldate | log.Llongfile)
}
