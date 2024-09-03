package main

import (
	"fmt"
	"sync"
)

// @desc:启动3个goroutine 循环100次顺序打印123

func main() {
	var wg sync.WaitGroup
	wg.Add(3) // 启动三个 goroutine

	// 创建通道用于控制打印顺序
	ch1 := make(chan struct{})
	ch2 := make(chan struct{})
	ch3 := make(chan struct{})

	// 启动第一个 goroutine 打印 1.md
	go func() {
		defer wg.Done() // 完成时通知 WaitGroup
		for i := 0; i < 100; i++ {
			fmt.Print("1.md")
			ch2 <- struct{}{} // 通知第二个 goroutine 打印 2
			<-ch1             // 等待被通知继续
		}
	}()

	// 启动第二个 goroutine 打印 2
	go func() {
		defer wg.Done() // 完成时通知 WaitGroup
		for i := 0; i < 100; i++ {
			<-ch2 // 等待第一个 goroutine 的通知
			fmt.Print("2")
			ch3 <- struct{}{} // 通知第三个 goroutine 打印 3
		}
	}()

	// 启动第三个 goroutine 打印 3
	go func() {
		defer wg.Done() // 完成时通知 WaitGroup
		for i := 0; i < 100; i++ {
			<-ch3 // 等待第二个 goroutine 的通知
			fmt.Print("3")
			ch1 <- struct{}{} // 通知第一个 goroutine 继续
		}
	}()

	// 首次开启第一个 goroutine，开始打印
	ch1 <- struct{}{} // 启动第一个 goroutine

	// 等待所有 goroutine 完成
	wg.Wait()
	fmt.Println() // 打印换行
}
