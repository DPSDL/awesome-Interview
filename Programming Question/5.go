package main

import "fmt"

//channel 的使用方法以及常见问题case

func main() {
	var ch chan int // nil channel

	// 尝试向 nil channel 发送值，这将导致 panic
	ch <- 1 // 运行时将在此处引发 panic: send on closed channel
	fmt.Println("This will not print")
}
