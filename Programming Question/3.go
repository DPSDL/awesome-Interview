package main

import (
	"fmt"
	"sync"
	"time"
)

// 编写一个程序限制10个goroutine执行，每执行完一个goroutine就放一个新的goroutine进来
func worker(id int, wg *sync.WaitGroup, sem chan struct{}) {
	defer wg.Done() // 通知 WaitGroup 当前 goroutine 完成
	fmt.Printf("Worker %d is starting\n", id)
	time.Sleep(time.Second) // 模拟一些工作
	fmt.Printf("Worker %d is done\n", id)

	<-sem // 释放一个信号量，允许其他 goroutine 进入
}

func main() {
	var wg sync.WaitGroup
	const maxWorkers = 10 // 最大 goroutine 数

	// 创建信号量通道
	sem := make(chan struct{}, maxWorkers)

	// 启动 20 个 goroutine
	for i := 1; i <= 20; i++ {
		wg.Add(1)         // 添加到 WaitGroup
		sem <- struct{}{} // 进入信号量，满的情况下会阻塞

		go worker(i, &wg, sem) // 启动 goroutine
	}

	// 等待所有 goroutine 完成
	wg.Wait()
	fmt.Println("All workers are done.")
}
