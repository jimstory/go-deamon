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

var redisAddr = lib.GetConfigure("REDIS_HOST") +":"+ lib.GetConfigure("REDIS_PORT")


// 建立redis 链接
var  redisConn = redis.NewClient(&redis.Options{
	Addr:     redisAddr,
	Password: "", // no password set
	DB:       0,  // use default DB
})
// 是否结束程序
var  quit =false
var quitChan = make(chan bool)

func worker(i int,queueChan chan string)  {
	waitGroup.Add(1)
	go func(i int) {
		// 如果程序不结束，一直执行
		str := ""
		for !quit {
			str  = <- queueChan
			lib.Trace.Println("worker id："+ strconv.Itoa(i) + " ,queue value :" + str)
			// 每隔协程处理时长，此处用来测试

			time.Sleep((time.Second * time.Duration(10)))
		}
		if str != "" {
			cmd := redisConn.RPush("queue1",str)
			if cmd.Err() != nil {
				fmt.Printf("%s\n", cmd.Err().Error())
			}
			lib.Trace.Println("worker id："+ strconv.Itoa(i) + " ,set to redis  queue value :" + str)
		}

		waitGroup.Done()
	}(i)

	// 依次开启10 个worker
	time.Sleep(time.Second*1)
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
		quitChan <- true
		close(queueChan)
		lib.Trace.Println("accept 中断信号")
	}()


	// 生成worker 协程
	for i:= 1;i<11;i++ {
		worker(i,queueChan)
	}

	// 从redis 队列取值，放入chan
	go func() {
		for {
			select {
				// 收到中断信号，退出协程
				case <- quitChan :
					return
			default:
				// 每隔10 s ，从redis 取次数据
				val, err := redisConn.BLPop(time.Second*10, "queue1").Result()
				if err != nil {
					if err == redis.Nil {
						//fmt.Println("queue is empty")
					} else {
						panic(err.Error())
					}
				} else {
					if !quit {
						// 信号没关闭，放入信道
						queueChan <- val[1]
					} else {
						// 重新放入redis
						redisConn.LPush("queue1", val)
					}
				}
			}
		}
	}()

	waitGroup.Wait()
}