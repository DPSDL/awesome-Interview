## new 和make的区别？
在 Go 语言中，`new` 和 `make` 是两个不同的内置函数，虽然它们的名字相似，但它们的用途和功能各有不同。以下是它们的主要区别：

### 1. `new`

- **用途**：用于分配内存并初始化为零值。
- **返回值**：返回一个指向类型的指针（例如 `*T`），其中 `T` 是你指定的类型。
- **用于**：一般用于创建值类型的指针，例如数组、结构体等。
- **示例**：

```go
type Person struct {
	Name string
	Age  int
}

p := new(Person) // p 是 *Person 类型的指针
fmt.Println(*p)  // 输出: { 0}
```

在这个示例中，`new(Person)` 返回一个指向 `Person` 结构体的指针。这个结构体内部的字段被初始化为零值。

### 2. `make`

- **用途**：用于初始化切片、映射（map）和通道（channel），并返回一个值，而不是指针。
- **返回值**：返回一个被初始化的数据结构的值。
- **用于**：只能用于切片、映射和通道，不能用于其他类型，如数组或结构体。
- **示例**：

```go
s := make([]int, 5) // 创建一个长度为 5 的切片
fmt.Println(s)      // 输出: [0 0 0 0 0]

m := make(map[string]int) // 创建一个空的映射
m["a"] = 1
fmt.Println(m) // 输出: map[a:1]

ch := make(chan int) // 创建一个 int 类型的通道
go func() { ch <- 42 }
fmt.Println(<-ch) // 输出: 42
```

### 总结

- **`new`**：
    - 用于分配内存，初始化为零值，返回类型的指针。
    - 适用于所有类型。

- **`make`**：
    - 用于初始化切片、映射和通道，返回初始化后的值。
    - 仅适用于这三种类型。

### 使用场景

- 当你需要创建一个结构体或基本数据类型的指针时，可以使用 `new`。
- 当你需要创建一个切片、映射或通道时，使用 `make`。

这两者在内存分配和初始化行为上的差异，使得 Go 的内存管理更为灵活且类型安全。如果你有其他问题或需要更深入的解释，请告知我！

## go协程上下文切换什么？
首先，我们需要知道几种比较常见的触发协程上下文切换的场景：

● 协作式抢占的时候，检查抢占标志位，从而触发上下文切换

● 锁阻塞/channel阻塞/io阻塞等等，这类最后都会调用runtime.gopark将当前协程挂起，然后调度新的协程

协程的上下文切换，涉及到挂起协程的上下文保存，和调度协程的上下文恢复

可以看到，需要保存的有：

● 与栈相关的SP和BP寄存器

● PC寄存器
● 用于保存函数闭包的上下文信息，也就是DX寄存器

而对于其他通用寄存器，因为go的函数调用规约，参数和返回值是通过栈进行传递的，并且总是在函数调用的时候触发协程切换，并不需要保存

##  init函数执行流程

在 Go 语言中，`init` 函数是一个特殊的初始化函数，它用于初始化包的状态。每个包可以有多个 `init` 函数，这些函数会在包被导入时自动执行。以下是 `init` 函数的执行流程的详细解释：

### 1. `init` 函数的定义

- **定义**：`init` 函数没有参数，也没有返回值，可以在包中定义多个 `init` 函数。
- **特点**：
    - 不需要手动调用，Go 运行时在包被导入时自动调用。
    - 如果在一个包中定义多个 `init` 函数，它们会按照出现的顺序执行。

### 2. 执行流程

当程序导入一个包时，执行流程如下：

#### 1. 包级变量的初始化

- 在执行任何 `init` 函数之前，Go 会首先初始化包级变量。
    - 如果变量有显式的初始值，则使用该值。
    - 如果没有初始值，则会设置为类型的零值。

#### 2. 执行 `init` 函数

- 所有包级变量初始化完成后，Go 会执行该包中的所有 `init` 函数，按照它们的出现顺序。

#### 3. 多个包的初始化

如果当前包依赖于其他包，初始化的过程会遵循下列顺序：

1. 初始化依赖的包。
2. 执行依赖包的所有 `init` 函数。
3. 初始化当前包。
4. 执行当前包的所有 `init` 函数。

### 3. 示例代码

以下是一个简单的示例，展示 `init` 函数的执行流程：

**包1 (`package1.go`)**：

```go
package package1

import "fmt"

// 网络级变量
var a = 10

// 初始化函数
func init() {
	fmt.Println("package1 init: a =", a)
}
```

**包2 (`package2.go`)**：

```go
package package2

import (
	"fmt"
	"yourmodule/package1" // 将此替换为适当的模块名
)

// 包级变量
var b = 20

// 初始化函数
func init() {
	fmt.Println("package2 init: b =", b)
}
```

**主程序 (`main.go`)**：

```go
package main

import (
	"fmt"
	"yourmodule/package2" // 将此替换为适当的模块名
)

// 包级变量
var c = 30

// 初始化函数
func init() {
	fmt.Println("main init: c =", c)
}

func main() {
	fmt.Println("Main function")
}
```

### 执行顺序

运行 `main.go` 时，输出将会是：

```
package1 init: a = 10
package2 init: b = 20
main init: c = 30
Main function
```

### 重要细节

1. **执行顺序**：所有 `init` 函数总是会在 `main` 函数执行之前运行。
2. **全局变量优先**：通过 `init` 函数初始化的全局变量确保在执行 `main` 函数之前完成。
3. **调试方便**：使用 `init`函数可以帮助调试和设置包的初始状态，而不需要在使用包时显式进行初始化操作。

### 总结

- `init` 函数在每个包被导入时自动执行。
- 如果一个包依赖于其他包，先初始化依赖包，后初始化当前包。
- 使用`init`函数可以自动设置包的初始状态。

如果你有进一步的问题或需要详细的信息，请告诉我！


## go GMP流程

在 Go 语言中，GMP 模型是 Go 运行时的调度模型，代表 Goroutines、M (操作系统线程) 和 P (处理器抽象)。这个模型让 Go 能够高效地管理并发和并行计算。下面是对 Go 中 GMP 流程的详细解释。

### 1. GMP 模型的组成

- **G (Goroutine)**: 是 Go 的轻量级线程，每个 Goroutine 的内存占用非常小，通常只有几 KB。
- **M (Machine)**: 是操作系统线程，每个 M 对应一个操作系统线程，用于执行 Goroutine。
- **P (Processor)**: 是 Go 的执行上下文，P 维护着一个可执行的 Goroutine 列表。每个 P 绑定一个 M，用于执行在该 P 上排队的 Goroutine。

### 2. GMP 的工作流程

#### 1. Goroutine 的创建

当使用 `go` 关键字创建 Goroutine 时，新的 Goroutine 会被添加到 P 的运行队列中，准备在 M 上运行。

```go
go myFunction() // 创建一个新的 Goroutine
```

#### 2. 调度

- 每个 M 在运行时都会尝试从其绑定的 P 的就绪队列中获取 Goroutine 进行执行。
- Go 的调度器会根据 Goroutine 的状态（就绪、等待、阻塞等）决定何时让某个 Goroutine 运行。

#### 3. 运行

- 如果 P 与 M 绑定，且 P 的就绪队列中有 Goroutine，M 将会执行一个 Goroutine。
- 当一个 Goroutine 执行完成或遇到阻塞（例如 I/O 操作）时，它可能会让出 P。这时，调度器会寻找其他就绪的 Goroutine。

#### 4. Goroutine 的阻塞和唤醒

- 当一个 Goroutine 产生阻塞（如等待 I/O、等待通道、系统调用等），相应的 M 会被释放并返回到空闲状态，准备去执行其他 Goroutine。
- 如果这个 Goroutine 有阻塞事件（比如 I/O 完成或通道有可用数据），调度器会尝试将其唤醒并放回 P 的就绪队列。

#### 5. 工作窃取和负载均衡

- Go 的调度器还实现了工作窃取机制。如果一个 P 的就绪队列为空，P 会试图从其他 P 的队列中窃取 Goroutine，以充分利用 CPU。
- 这个机制帮助 Go 达到负载均衡，并增加 Goroutine 的并发度。

### 3. 示例代码

以下是一段简单的代码，演示了如何利用 Goroutines：

```go
package main

import (
	"fmt"
	"sync"
	"time"
)

func count(id int, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 0; i < 5; i++ {
		fmt.Printf("Goroutine %d: %d\n", id, i)
		time.Sleep(100 * time.Millisecond)
	}
}

func main() {
	var wg sync.WaitGroup
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go count(i, &wg) // 创建 Goroutine
	}
	wg.Wait() // 等待所有 Goroutines 完成
}
```

### 4. 总结

- **G**：轻量级的 Goroutine，允许多任务并发。
- **M**：操作系统线程，用于执行 Goroutine。
- **P**：处理器抽象，管理 Goroutine 的执行。
- **调度器**：调度 Goroutine 的执行，将其分配到 M 上运行，以实现高效的并发。

通过 `GMP` 模型，Go 语言能够高效地管理数千到数百万的 Goroutine，从而在并发操作时保持高性能。这使得 Go 在处理网络、并发 I/O 等场景时非常有优势。如果你有更多问题或者希望深入某些方面，请告诉我！

## go 内存管理