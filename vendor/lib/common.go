package lib

import (
    "os"
    "bufio"
    "strings"
    "fmt"
    "time"
    "log"
)

// 获取配置文件信息
func GetConfigure(search string) string {
    confFile, err  := os.Open(".env")
    // 读取配置文件
    if  err != nil {
        fmt.Println(err)
    }
    defer confFile.Close()
    scanner := bufio.NewScanner(confFile)
    for  scanner.Scan() {
        txt := strings.TrimSpace(scanner.Text())
        if txt != "" && strings.Index(txt, "#") != 0 && strings.Index(txt,search) == 0 {
            pos := strings.Split(txt,"=")
            return strings.TrimSpace(pos[1])
        }
    }
    return ""
}

// 打印日志,返回全局log 对象，并在main 函数中关闭
func Logger() (* log.Logger,*os.File){
    date := time.Now().Format("2006-01-02")
    logPath := GetConfigure("LOG_PATH")
    debuggerFile, err := os.Create("./"+logPath+"/" + date + ".log")
    if err != nil {
        fmt.Println(err)
        log.Fatalln("fail  to create debug log")
    }
    //defer debuggerFile.Close()
    log := log.New(debuggerFile,"",log.Ldate | log.Llongfile)
    return log,debuggerFile
}
