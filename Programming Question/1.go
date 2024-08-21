package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

//使用go实现1000个并发控制并设置执行超时时间1秒

func main() {
	const numGoroutines = 1000

	var wg sync.WaitGroup

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			select {
			case <-time.After(500 * time.Millisecond):
				fmt.Printf("Goroutine %d completed\n", i)
			case <-ctx.Done():
				fmt.Printf("Goroutine %d timed out\n", i)
			}
		}(i)
	}
	wg.Wait()

}
