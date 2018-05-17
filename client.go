package main

import (
	"fmt"
	"sync"
	"github.com/go-redis/redis"
	"strconv"
)

func in_array(str string, strArr []string) bool {
	for _,val := range strArr {
		if str == val {
			return true
		}
	}
	return false
}

//defer redisConn.Close()
var waitGroup sync.WaitGroup

func worker(i int,queueChan chan string,quit bool)  {
	waitGroup.Add(1)
	go func(i int) {
		// 如果程序不结束，一直执行
		for !quit {
			fmt.Println("worker id："+ strconv.Itoa(i) + " queue value :" + <- queueChan)
		}
		waitGroup.Done()
	}(i)
}

func main()  {

	// 队列信道
	queueChan := make(chan string, 10)

	//// 是否结束程序
	quit := false

	// 建立redis 链接
	redisConn := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	// 生成worker 协程
	for i:= 0;i<10;i++ {
		worker(i,queueChan,quit)
	}

	// 从redis 队列取值，放入chan
	go func() {
		for {
			val, err := redisConn.LPop("queue").Result()
			if err != nil {
				if (err == redis.Nil) {
					//fmt.Println("queue is empty")
				} else {
					panic(err.Error())
				}
			} else {
				queueChan <- val
			}
		}
	}()
	waitGroup.Wait()
}