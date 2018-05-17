package main

import (
	"fmt"
	"sync"
	"github.com/go-redis/redis"
	"strconv"
	"lib"
	"os"
	"os/signal"
	"time"
)

var waitGroup sync.WaitGroup

// 建立redis 链接
var  redisConn = redis.NewClient(&redis.Options{
Addr:     "localhost:6379",
Password: "", // no password set
DB:       0,  // use default DB
})
// 是否结束程序
var  quit bool =false

func worker(i int,queueChan chan string)  {
	waitGroup.Add(1)
	go func(i int) {
		// 如果程序不结束，一直执行
		str := ""
		for !quit {
			str  = <- queueChan
			fmt.Println("worker id："+ strconv.Itoa(i) + " queue value :" + str)
		}
		redisConn.LPush("queue1",str)
		lib.Trace.Println("worker id："+ strconv.Itoa(i) + " set to redis  queue value :" + str)
		waitGroup.Done()
	}(i)
}

func main()  {
	// 队列信道
	queueChan := make(chan string, 10)

	// 监听信号
	sigs := make(chan os.Signal, 1)
	//接受kill 和中断信号
	signal.Notify(sigs, os.Interrupt, os.Kill)

	// 监听信号，进行处理
	go func() {
		<-sigs
		// 收到信号，进行退出,关闭信道
		quit = true
		close(queueChan)
		lib.Trace.Println("accept 中断信号")
		os.Exit(155)
	}()


	// 生成worker 协程
	for i:= 0;i<10;i++ {
		//worker(i,queueChan)

		waitGroup.Add(1)
		go func(i int) {
			// 如果程序不结束，一直执行
			str := ""
			for !quit {
				str  = <- queueChan
				fmt.Println("worker id："+ strconv.Itoa(i) + " queue value :" + str)
			}
			redisConn.LPush("queue1",str)
			lib.Trace.Println("worker id："+ strconv.Itoa(i) + " set to redis  queue value :" + str)
			waitGroup.Done()
		}(i)
	}

	// 从redis 队列取值，放入chan
	go func() {
		for {
			val, err := redisConn.LPop("queue1").Result()
			if err != nil {
				if err == redis.Nil {
					//fmt.Println("queue is empty")
				} else {
					panic(err.Error())
				}
			} else {
				queueChan <- val
			}

			// 每隔10s 取一次数据
			time.Sleep(time.Second*10)
		}
	}()

	waitGroup.Wait()
}