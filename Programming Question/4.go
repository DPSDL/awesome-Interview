package main

import (
	"fmt"
	"math/rand"
	"time"
)

// 写代码实现两个 goroutine，其中一个产生随机数并写入到 go channel 中，另外一个从 channel 中读取数字并打印到标准输出。最终输出五个随机数。

func main() {
	// 创建一个 channel，用于传递随机数
	randChannel := make(chan int)

	// 启动第一个 goroutine 产生随机数并写入 channel
	go func() {
		for i := 0; i < 5; i++ {
			// 产生随机数
			num := rand.Intn(100)              // 生成 0 到 99 之间的随机数
			randChannel <- num                 // 发送到 channel
			time.Sleep(time.Millisecond * 500) // 暂停一段时间以便产生不同的随机数
		}
		close(randChannel) // 关闭 channel，表示没有更多数据发送
	}()

	// 启动第二个 goroutine 从 channel 中读取数字并打印
	go func() {
		for num := range randChannel {
			fmt.Println(num) // 打印从 channel 中读取的数字
		}
	}()

	// 等待输入，以便可以看到结果
	fmt.Scanln()
}
