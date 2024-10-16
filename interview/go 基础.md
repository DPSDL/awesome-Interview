
Golang

new 和make的区别？
● make和new都是golang用来分配内存的內建函数，且在堆上分配内存，make 即分配内存，也初始化内存。new只是将内存清零，并没有初始化内存。
● make返回的还是引用类型本身；而new返回的是指向类型的指针。
● make只能用来分配及初始化类型为slice，map，channel的数据；new可以分配任意类型的数据。
go协程上下文切换什么？
首先，我们需要知道几种比较常见的触发协程上下文切换的场景：
● 协作式抢占的时候，检查抢占标志位，从而触发上下文切换
● 锁阻塞/channel阻塞/io阻塞等等，这类最后都会调用runtime.gopark将当前协程挂起，然后调度新的协程
协程的上下文切换，涉及到挂起协程的上下文保存，和调度协程的上下文恢复

可以看到，需要保存的有：
● 与栈相关的SP和BP寄存器
● PC寄存器
● 用于保存函数闭包的上下文信息，也就是DX寄存器
而对于其他通用寄存器，因为go的函数调用规约，参数和返回值是通过栈进行传递的，并且总是在函数调用的时候触发协程切换，并不需要保存


以看到，mcall中就会保存当前g的上下文，然后切换到g0栈执行。
上下文本切换

你能说说go语言中的垃圾回么？

好的，一般我们在程序中会申请新的变量，当这些变量不在需要的时候，就需要我们手动或者程序自动的将这些内存释放掉，这一过程就叫做垃圾回收（GC）,其中go语言中有自动回收垃圾的机制。

go语言中的垃圾回收

当前goland中使用的的垃圾回收主要是三色标记法配合写屏障和辅助GC。三色标记法是标记-清除法的一种增强版本。

标记-清除法（Mark and sweep）

原始的标记法分为两个步骤：
1.标记。先STW(stop the world),暂停整个程序的全部运行线程，将被引用的对象打上标记。
2.清除没有被打上标记的对象。然后恢复所有线程。

三色标记法

三色标记法是对标记阶段的改进：
1.初始状态所有的对象都是白色；
2.从root出发进行扫描所有的根对象，将所有他们引用的对象标记为灰色。
3.分析灰色对象是否引用了其他对象，如果没有引用则将该灰色对象变为黑色；如果有引用则将该灰色对象变为黑色同时将他引用的对象变为灰色。
4.不断重复步骤3，知道所有的灰色对象为空。此时白色对象即为需要回收的对象。

GO的GC是如何工作的？

go的GC主要如下：

1. Mark包含两部分：

● Mark Prepare:初始化GC任务,包括开启写屏障以及辅助GC。整个过程需要STW。
● GC Drains:扫描所有根节点root对象，循环标记对象，直到灰色对象为空。整个过程是后台进行的。

2. Mark Termination:完成标记对象，重新扫描全局指针和堆栈，整个过程会有新的对象可能加进来，需要通过写屏障标记下来，再次扫描一遍。整个过程需要STW。
3.SWeep:按照标记结果，回收所有的白色对象。该过程在后台进行。

写屏障：主要就是记录标记前后对象的颜色状态，对标记过程中某些对象的标记进行修改。

辅助GC：主要就是有些时候GC的回收速度赶不上用户代码分配的速度，就需要进行STW，保证新把垃圾对象清理掉，然后用户可以分配对象。
defer 的作用和特点是什么？
defer 函数调用的顺序是后进先出，当产生 panic 的时候，会先执行 panic 前面的 defer 函数后才真的抛出异常。一般的，recover 会在 defer 函数里执行并捕获异常，防止程序崩溃。
defer 的作用是：
● 你只需要在调用普通函数或方法前加上关键字 defer，就完成了 defer 所需要的语法。当 defer 语句被执行时，跟在 defer 后面的函数会被延迟执行。直到包含该 defer 语句的函数执行完毕时，defer 后的函数才会被执行，不论包含defer 语句的函数是通过 return 正常结束，还是由于 panic 导致的异常结束。你可以在一个函数中执行多条 defer 语句，它们的执行顺序与声明顺序相反。
defer 的常用场景：
● defer 语句经常被用于处理成对的操作，如打开、关闭、连接、断开连接、加锁、释放锁。
● 通过 defer 机制，不论函数逻辑多复杂，都能保证在任何执行路径下，资源被释放。
● 释放资源的 defer 应该直接跟在请求资源的语句后。
数组和切片的区别是什么？
● 数组： 数组固定长度。数组长度是数组类型的一部分，所以[3]int 和[4]int 是两种不同的数组类型数组需要指定大小，不指定也会根据初始化，自动推算出大小，大小不可改变。数组是通过值传递的；
● 切片： 切片可以改变长度。切片是轻量级的数据结构，三个属性，指针，长度，容量。不需要指定大小切片是地址传递（引用传递）可以通过数组来初始化，也可以通过内置函数 make()来初始化，初始化的时候 len=cap，然后进行扩容。
 Slice 的底层实现：
切片是基于数组实现的，它的底层是数组，它自己本身非常小，可以理解为对底层数组的抽象。因为基于数组实现，所以它的底层的内存是连续分配的，效率非常高，还可以通过索引获得数据。
切片本身并不是动态数组或者数组指针。它内部实现的数据结构通过指针引用底层数组，设定相关属性将数据读写操作限定在指定的区域内。切片本身是一个只读对象，其工作机制类似数组指针的一种封装。
切片对象非常小，是因为它是只有 3 个字段的数据结构：
● 指向底层数组的指针
● 切片的长度
● 切片的容量
Slice 的扩容机制：
● 首先判断，如果新申请容量大于 2 倍的旧容量，最终容量就是新申请的容量
● 否则判断，如果旧切片的长度小于 1024，则最终容量就是旧容量的两倍
● 否则判断，如果旧切片长度大于等于 1024，则最终容量从旧容量开始循环增加原来的 1/4, 直到最终容量大于等于新申请的容量
● 如果最终容量计算值溢出，则最终容量就是新申请容量
介绍一下Context
context 从它的字面量就可以看出来，是用来传递上下文信息的。在 Go 里并没有直接为我们提供一个统一的 context 对象，而是设计了一个接口类型的 Context。然后在这些接口上来实现了几种具体类型的 context。我们来看看官方的 context 类型。主要有四种：
● emptyCtx：空的 context，实现了上面的 4 个接口，但都是直接 return 默认值，没有具体功能代码。
● cancelCtx：用来取消通知用的 context
● timerCtx：用来超时通知用的 context
● valueCtx：用来传值的 context
map 的底层原理是什么？
Golang 中 map 的底层实现是一个散列表，因此实现 map 的过程实际上就是实现散表的过程。在这个散列表中，主要出现的结构体有两个，一个叫 hmap(a header for a go map)，一个叫 bmap(a bucket for a Go map，通常叫其 bucket)。
map 查找过程：
● Go 语言中 map 采用的是哈希查找表，由一个 key 通过哈希函数得到哈希值，64位系统中就生成一个 64bit 的哈希值，由这个哈希值将 key 对应存到不同的桶（bucket）中，当有多个哈希映射到相同的的桶中时，使用链表解决哈希冲突。细节：key 经过 hash 后共 64 位，根据 hmap 中 B 的值，计算它到底要落在哪个桶时，桶的数量为 2^B，如 B=5，那么用 64 位最后 5 位表示第几号桶，在用 hash 值的高 8 位确定在 bucket 中的存储位置，当前 bmap 中的 bucket 未找到，则查询对应的 overflow bucket，对应位置有数据则对比完整的哈希值，确定是否是要查找的数据。如果当前 map 处于数据搬移状态，则优先从 oldbuckets 查找。
● 扩容：map 扩容的目的是为了避免增大哈希碰撞的概率，否则如果大量Key落在同一个bucket 里，会导致map退化成链表，那么访问的时间的复杂度也会降为O(n)。在 go map 中有一个指标来衡量map的元素数量和桶的比值（装载因子），扩容时机：① 装载因子超过阈值，源码里定义的阈值是 6.5。这种情况下map会进行翻倍扩容，创建一个新的buckets数组，其容量是之前的两倍，并将旧数据逐步迁移至新的buckets。② 由overflow指针所连接的溢出桶数量过多。出现了这种情况可能是因为哈希表里有过多的空键值对（可能是由于原桶中有太多的键值对被删除），这个时候就需要通过map扩容做内存整理。目的就是为了清除bmap桶中空闲的键值对。这种情况下map扩容步骤与情况一基本相同，只不过扩容后map容量还是原来的大小。Go会创建一个与原buckets数组容量相同的buckets数组，并将旧数组的数据逐步迁移至这个新数组。扩容过程：由于 map 扩容需要将原有的数据重新搬迁到新的内存地址，如果有大量数据需要搬迁，会非常影响性能。因此 Go map 的扩容采取了一种称为“渐进式”地方式，原有的数据并不会一次性搬迁完毕，每次最多只会搬迁 2 个 bucket。

● 当我们插入一个k-v对时，需要确定他应该插入到bucket数组的哪一个槽中。bucket数组的长度为2^B，即2的次幂数，而2^B-1转换成二进制后一定是低位全1，高位全0的形式，因此在进行按位与操作后，一定能求得一个在[0,2^B-1]区间上的任意一个数，也就是数组中的下标位置，相较之下，能获得比取模更加优秀的执行效率。
● 涉及到扩容，每一次bucket数组都会变为现在的两倍，方便我们进行hash迁移。
map触发扩容的条件有两种：
1. 负载因子大于6.5时（负载因子 = 键数量 / bucket数量）
2. overflow的数量达到2^min(15,B)
等量扩容 所谓等量扩容，并不是扩大容量，而是bucket数量不变，重新做一遍类似增量扩容的搬迁动作，把松散的键值对重新排列一次，以使bucket的使用率更高，从而保证更快的存取速度。
sync.Map 了解吗？
普通的map在并发读写下会造成panic，在 Go 1.9 引入的一种并发安全 map，适用于读多写少的场景。一般情况下解决并发读写 map 的思路是加一把大锁，或者把一个 map 分成若干个小 map，对 key 进行哈希，只操作相应的小 map。前者锁的粒度比较大，影响效率；后者实现起来比较复杂，容易出错。而使用 sync.map 之后，对 map 的读写，不需要加锁。并且它通过空间换时间的方式，使用 read 和 dirty 两个 map 来进行读写分离，降低锁时间来提高效率。
sync.map 适用于读多写少的场景。对于写多的场景，会导致 read map 缓存失效，需要加锁，导致冲突变多；而且由于未命中 read map 次数过多，导致 dirty map 提升为 read map，这是一个 O(N) 的操作，会进一步降低性能。

Go 的内建 map 是不支持并发写操作的，原因是 map 写操作不是并发安全的，当你尝试多个 Goroutine 操作同一个 map，会产生报错：fatal error: concurrent map writes。
sync.Map 类型的底层数据结构如下：
type Map struct {  mu Mutex  read atomic.Value // readOnly  dirty map[interface{}]*entry  misses int }
● mu：互斥锁，用于保护 read 和 dirty。
● read：只读数据，支持并发读取（atomic.Value 类型）。如果涉及到更新操作，则只需要加锁来保证数据安全。read 实际存储的是 readOnly 结构体，内部也是一个原生 map，amended 属性用于标记 read 和 dirty 的数据是否一致。
● dirty：读写数据，是一个原生 map，也就是非线程安全。操作 dirty 需要加锁来保证数据安全。
● misses：统计有多少次读取 read 没有命中。每次 read 中读取失败后，misses 的计数值都会加 1
channel 是什么？
Go 语言中，不要通过共享内存来通信，而要通过通信来实现内存共享。Go 的 CSP(Communicating Sequential Process)并发模型，中文可以叫做通信顺序进程，是通过 goroutine 和 channel 来实现的。channel 收发遵循先进先出 FIFO 的原则。分为有缓冲区和无缓冲区，channel 中包括 buffer、sendx 和 recvx 收发的位置(ring buffer 记录实现)、sendq、recv。当 channel 因为缓冲区不足而阻塞了队列，则使用双向链表存储。
● 给一个 nil channel 发送数据，造成永远阻塞
● 从一个 nil channel 接收数据，造成永远阻塞
● 给一个已经关闭的 channel 发送数据，引起 panic
● 从一个已经关闭的 channel 接收数据，如果缓冲区中为空，则返回一个零值
● 无缓冲的 channel 是同步的，而有缓冲的 channel 是非同步的
● 关闭一个 nil channel 将会发生 panic
channel 中使用了 ring buffer（环形缓冲区) 来缓存写入的数据。ring buffer 有很多好处，而且非常适合用来实现 FIFO 式的固定长度队列。

channel是golang中用来实现多个goroutine通信的管道，它的底层是一个叫做hchan的结构体。在go的runtime包下。

总结hchan结构体的主要组成部分有四个：
● 用来保存goroutine之间传递数据的循环数组。=====> buf。
● 用来记录此循环链表当前发送或接收数据的下标值。=====> sendx和recvx。
● 用于保存向该chan发送和从改chan接收数据的goroutine的队列。=====> sendq 和 recvq
● 保证channel写入和读取数据时线程安全的锁。 =====> lock
channel 分配在栈上还是堆上？哪些对象分配在堆上，哪些对象分配在栈上？
Channel 被设计用来实现协程间通信的组件，其作用域和生命周期不可能仅限于某个函数内部，所以 golang 直接将其分配在堆上。Golang 中的变量只要被引用就一直会存活，存储在堆上还是栈上由内部实现决定而和具体的语法没有关系。知道变量的存储位置确实和效率编程有关系。如果可能，Golang 编译器会将函数的局部变量分配到函数栈帧（stack frame）上。然而，如果编译器不能确保变量在函数 return 之后不再被引用，编译器就会将变量分配到堆上。而且，如果一个局部变量非常大，那么它也应该被分配到堆上而不是栈上。当前情况下，如果一个变量被取地址，那么它就有可能被分配到堆上,然而，还要对这些变量做逃逸分析，如果函数 return 之后，变量不再被引用，则将其分配到栈上。
Mutex 有几种状态？
● mutexLocked — 表示互斥锁的锁定状态；
● mutexWoken — 表示从正常模式被从唤醒；
● mutexStarving — 当前的互斥锁进入饥饿状态；
● waitersCount — 当前互斥锁上等待的 Goroutine 个数；
Mutex 正常模式和饥饿模式？
正常模式（非公平锁）
正常模式下，所有等待锁的 goroutine 按照 FIFO（先进先出）顺序等待。唤醒的 goroutine 不会直接拥有锁，而是会和新请求 goroutine 竞争锁。新请求的goroutine 更容易抢占：因为它正在 CPU 上执行，所以刚刚唤醒的 goroutine 有很大可能在锁竞争中失败。在这种情况下，这个被唤醒的 goroutine 会加入到等待队列的前面。
饥饿模式（公平锁）
为了解决了等待 goroutine 队列的长尾问题饥饿模式下，直接由 unlock 把锁交给等待队列中排在第一位的 goroutine (队头)，同时，饥饿模式下，新进来的 goroutine 不会参与抢锁也不会进入自旋状态，会直接进入等待队列的尾部。这样很好的解决了老的 goroutine 一直抢不到锁的场景。饥饿模式的触发条件：当一个 goroutine 等待锁时间超过 1 毫秒时，或者当前队列只剩下一个 goroutine 的时候，Mutex 切换到饥饿模式。
总结
对于两种模式，正常模式下的性能是最好的，goroutine 可以连续多次获取锁，饥饿模式解决了取锁公平的问题，但是性能会下降，这其实是性能和公平的一个平衡模式。
自旋
如果 Goroutine 占用锁资源的时间比较短，那么每次释放资源后，都调用信号量来唤起正在阻塞等候的goroutine，将会很浪费资源。因此在符合一定条件后，mutex 会让等候的 Goroutine 去空转 CPU，在空转完后再次调用 CAS 方法去尝试性的占有锁资源，直到不满足自旋条件，则最终才加入到等待队列里。
RWMutex 实现？
通过记录 readerCount 读锁的数量来进行控制，当有一个写锁的时候，会将读 锁数量设置为负数 1<<30。目的是让新进入的读锁等待之前的写锁释放通知读 锁。同样的当有写锁进行抢占时，也会等待之前的读锁都释放完毕，才会开始进行后续的操作。 而等写锁释放完之后，会将值重新加上 1<<30, 并通知刚才 新进入的读锁（rw.readerSem），两者互相限制。
● RWMutex 是单写多读锁，该锁可以加多个读锁或者一个写锁
● 读锁占用的情况下会阻止写，不会阻止读，多个 Goroutine 可以同时获取读锁
● 写锁会阻止其他 Goroutine（无论读和写）进来，整个锁由该 Goroutine 独占
● 适用于读多写少的场景
● RWMutex 类型变量的零值是一个未锁定状态的互斥锁
● RWMutex 在首次被使用之后就不能再被拷贝
● RWMutex 的读锁或写锁在未锁定状态，解锁操作都会引发 panic
● RWMutex 的一个写锁去锁定临界区的共享资源，如果临界区的共享资源已被（读锁或写锁）锁定，这个写锁操作的 goroutine 将被阻塞直到解锁
● RWMutex 的读锁不要用于递归调用，比较容易产生死锁
● RWMutex 的锁定状态与特定的 goroutine 没有关联。一个 goroutine 可以 RLock（Lock），另一个 goroutine 可以 RUnlock（Unlock）
● 写锁被解锁后，所有因操作锁定读锁而被阻塞的 goroutine 会被唤醒，并都可以成功锁定读锁
● 读锁被解锁后，在没有被其他读锁锁定的前提下，所有因操作锁定写锁而被阻塞的 Goroutine，其中等待时间最长的一个 Goroutine 会被唤醒
WaitGroup 的原理？
一个 WaitGroup 对象可以等待一组协程结束。WaitGroup 主要维护了 2 个计数器，一个是请求计数器 v，一个是等待计数器 w，二者组成一个 64bit 的值，请求计数器占高 32bit，等待计数器占低32bit。每次 Add 执行，请求计数器 v 加 1，Done 方法执行，等待计数器减 1，v 为 0 时通过信号量唤醒 Wait()。
GMP 指的是什么？
G 代表着 goroutine，P 代表着上下文处理器，M 代表 thread 线程，在 GPM 模型，有一个全局队列（Global Queue）：存放等待运行的 G，还有一个 P 的本地队列：也是存放等待运行的 G，但数量有限，不超过 256 个。GPM 的调度流程从 go func()开始创建一个 goroutine，新建的 goroutine 优先保存在 P 的本地队列中，如果 P 的本地队列已经满了，则会保存到全局队列中。M 会从 P 的队列中取一个可执行状态的 G 来执行，如果 P 的本地队列为空，就会从其他的 MP 组合偷取一个可执行的 G 来执行，当 M 执行某一个 G 时候发生系统调用或者阻塞，M 阻塞，如果这个时候 G 在执行，runtime 会把这个线程 M 从 P 中摘除，然后创建一个新的操作系统线程来服务于这个 P，当 M 系统调用结束时，这个 G 会尝试获取一个空闲的 P 来执行，并放入到这个 P 的本地队列，如果这个线程 M 变成休眠状态，加入到空闲线程中，然后整个 G 就会被放入到全局队列中。
关于 G,P,M 的个数问题，G 的个数理论上是无限制的，但是受内存限制，P 的数量一般建议是逻辑 CPU 数量的 2 倍，M 的数据默认启动的时候是 10000，内核很难支持这么多线程数，所以整个限制客户忽略，M 一般不做设置，设置好 P，M 一般都是要大于 P。
● G（Goroutine）：我们所说的协程，为用户级的轻量级线程，每个 Goroutine 对象中的 scheduer 保存着其上下文信息。
● M（Machine）：对内核级线程的封装，数量对应真实的 CPU 数（真正干活的对象）。它是真正的调度执行者，M 需要跟 P 绑定，并且会让 P 按下面的原则挑出个 goroutine 来执行：优先从 P 的本地队列获取 goroutine 来执行；如果本地队列没有，从全局队列获取，如果全局队列也没有，会从其他的 P 上偷取 goroutine。
● P（Processor）：即为 G 和 M 的调度对象，用来调度 G 和 M 之间的关联关系，其数量可通过 GOMAXPROCS()来设置，默认为核心数。每当有 goroutine 要创建时，会被添加到 P 上的 goroutine 本地队列上，如果 P 的本地队列已满，则会维护到全局队列里。
GC 机制 ？
Go 采用的是三色标记法，将内存里的对象分为了三种：
● 白色对象：未被使用的对象；
● 灰色对象：当前对象有引用对象，但是还没有对引用对象继续扫描过；
● 黑色对象，对上面提到的灰色对象的引用对象已经全部扫描过了，下次不用再扫描它了。
当垃圾回收开始时，Go 会把根对象标记为灰色，其他对象标记为白色，然后从根对象遍历搜索，按照上面的定义去不断的对灰色对象进行扫描标记。当没有灰色对象时，表示所有对象已扫描过，然后就可以开始清除白色对象了。
GC 是怎么实现的？
Go 的 GC 回收有三次演进过程：
● Go1.3 之前普通标记清除（mark and sweep）方法，整体过程需要启动 STW，效率极低。Go1.5 三色标记法，堆空间启动写屏障，栈空间不启动，全部扫描之后，需要重新扫描一次栈(需要 STW)，效率普通。Go1.8 三色标记法，混合写屏障机制：栈空间不启动（全部标记成黑色），堆空间启用写屏障，整个过程不要 STW，效率高。
● Go1.3 之前的版本所谓标记清除是先启动 STW 暂停，然后执行标记，再执行数据回收，最后停止 STW。Go1.3 版本标记清除做了点优化，流程是：先启动 STW 暂停，然后执行标记，停止 STW，最后再执行数据回收。
● Go1.5 三色标记主要是插入屏障和删除屏障，写入屏障的流程：程序开始，全部标记为白色，1）所有的对象放到白色集合，2）遍历一次根节点，得到灰色节点，3）遍历灰色节点，将可达的对象，从白色标记灰色，遍历之后的灰色标记成黑色，4）由于并发特性，此刻外界向在堆中的对象发生添加对象，以及在栈中的对象添加对象，在堆中的对象会触发插入屏障机制，栈中的对象不触发，5）由于堆中对象插入屏障，则会把堆中黑色对象添加的白色对象改成灰色，栈中的黑色对象添加的白色对象依然是白色，6）循环第 5 步，直到没有灰色节点，7）在准备回收白色前，重新遍历扫描一次栈空间，加上 STW 暂停保护栈，防止外界干扰（有新的白色会被添加成黑色）在 STW 中，将栈中的对象一次三色标记，直到没有灰色，8）停止 STW，清除白色。至于删除写屏障，则是遍历灰色节点的时候出现可达的节点被删除，这个时候触发删除写屏障，这个可达的被删除的节点也是灰色，等循环三色标记之后，直到没有灰色节点，然后清理白色，删除写屏障会造成一个对象即使被删除了最后一个指向它的指针也依旧可以活过这一轮，在下一轮 GC 中被清理掉。
● GoV1.8 混合写屏障规则是：
  ○ GC 开始将栈上的对象全部扫描并标记为黑色(之后不再进行第二次重复扫描，无需 STW)，
  ○ GC 期间，任何在栈上创建的新对象，均为黑色。
  ○ 被删除的对象标记为灰色
  ○ 被添加的对象标记为灰色。
GC 的触发时机？
分为系统触发和主动触发。
● gcTriggerHeap：当所分配的堆大小达到阈值（由控制器计算的触发堆的大小）时，将会触发。
● gcTriggerTime：当距离上一个 GC 周期的时间超过一定时间时，将会触发。时间周期以runtime.forcegcperiod 变量为准，默认 2 分钟。
● gcTriggerCycle：如果没有开启 GC，则启动 GC。
● 手动触发的 runtime.GC 方法。
GC 的流程是什么？
Go1.14 版本以 STW 为界限，可以将 GC 划分为五个阶段：
● GCMark 标记准备阶段，为并发标记做准备工作，启动写屏障
● STWGCMark 扫描标记阶段，与赋值器并发执行，写屏障开启并发
● GCMarkTermination 标记终止阶段，保证一个周期内标记任务完成，停止写屏障
● GCoff 内存清扫阶段，将需要回收的内存归还到堆中，写屏障关闭
● GCoff 内存归还阶段，将过多的内存归还给操作系统，写屏障关闭。
GC 如何调优？
通过 go tool pprof 和 go tool trace 等工具
● 控制内存分配的速度，限制 Goroutine 的数量，从而提高赋值器对 CPU 的利用率。
● 减少并复用内存，例如使用 sync.Pool 来复用需要频繁创建临时对象，例如提前分配足够的内存来降低多余的拷贝。
● 需要时，增大 GOGC 的值，降低 GC 的运行频率。
知道 golang 的内存逃逸吗？
本该分配到栈上的变量，跑到了堆上，这就导致了内存逃逸。2)栈是高地址到低地址，栈上的变量，函数结束后变量会跟着回收掉，不会有额外性能的开销。3)变量从栈逃逸到堆上，如果要回收掉，需要进行 gc，那么 gc 一定会带来额外的性能开销。编程语言不断优化 gc 算法，主要目的都是为了减少 gc 带来的额外性能开销，变量一旦逃逸会导致性能开销变大。
内存逃逸的情况如下：
● 方法内返回局部变量指针。
● 向 channel 发送指针数据。
● 在闭包中引用包外的值。
● 在 slice 或 map 中存储指针。
● 切片（扩容后）长度太大。
● 在 interface 类型上调用方法。
请简述 Go 是如何分配内存的？
Go 程序启动的时候申请一大块内存，并且划分 spans，bitmap，areana 区域；

● arena 区域按照页划分成一个个小块，span 管理一个或者多个页，mcentral 管理多个 span 供现场申请使用；mcache 作为线程私有资源，来源于 mcentral。
● bitmap区用于保存arena对应地址(指针大小为 8B，bitmap中一个byte大小的内存对应arena区域中的4个指针，因此大小为 512G/(4 * 8B))中是否保存了对象，以及对象是否被gc过，主要用于gc。
● spans区域存放了mspan的指针，用于表示arena区的某一页属于哪个mspan。大小为 512G / 8KB(页的大小) * 8B(指针大小)。在创建 mspan的时候，按页填充对应的spans区域，在回收object时，很容易找到他所属的mspan。
Go是内置运行时runtime的语言；像这种内置运行时的语言会抛弃传统的内存管理方式，改为自己管理；这样可以完成类似预分配，内存池等操作，以避开系统调用产生的性能问题。Go的内存分配可以分为以下几点：
1. 每次从操作系统申请一大块内存，由Go来对这块内存做分配，减少系统调用。
2. 内存分配算法采用 TCMalloc算法，核心思想是把内存分的非常细，进行分级管理，以降低锁的粒度。
3. 回收对象内存时，并没有将其真正释放掉，而是放回预先分配的大内存中，以便复用。只有内存闲置过多的时候，才会尝试归还部分内存给操作系统，降低整体开销。
内存分配由内存分配器完成，分配器由3种组件构成：mcache, mcentral, mheap：
● mcache：每个工作线程都会绑定一个mcache，本地缓存可用的mspan资源，这样就可以直接给Goroutine分配，因为不存在多个Goroutine竞争的情况，所以不会消耗锁资源。
● mcentral：为所有mcache提供切分好的mspan资源。每个central保存一种特定大小的全局mspan列表，包括已分配出去的和未分配出去的。 每个mcentral对应一种mspan，而mspan的种类导致它分割的object大小不同。当工作线程的mcache中没有合适（也就是特定大小的）的mspan时就会从mcentral获取。
● mheap：代表Go程序持有的所有堆空间，Go程序使用一个mheap的全局对象_mheap来管理堆内存。当mcentral没有空闲的mspan时，会向mheap申请。而mheap没有资源时，会向操作系统申请新内存。mheap主要用于大对象的内存分配，以及管理未切割的mspan，用于给mcentral切割成小对象。同时我们也看到，mheap中含有所有规格的mcentral，所以，当一个mcache从mcentral申请mspan时，只需要在独立的mcentral中使用锁，并不会影响申请其他规格的mspan。
分布式
zookeeper 分布式锁
https://blog.csdn.net/m0_67322837/article/details/126016310

CAP 是什么？
CAP原理指的是在一个分布式系统中，Consistency（一致性）、 Availability（可用性）、 Partition tolerance（分区容错性），最多只能同时三个特性中的两个，三者不可兼得。
● 一致性：在一致性的需求下，当一个系统在数据一致的状态下执行更新操作后，应该保证系统 的数据仍然处于一致的状态。
● 可用性：系统提供的服务必须一直处于可用的状态，对于用户的每一个操作请求总是能够在 有限的时间内返回结果。
● 分区容错性：分布式系统在遇到任何网络分区故障的时候，仍然需要能够保证对外提供满足一致性和可用性的服务，除非是整个网络环境都发生了故障 。在分布式环境中，由于网络的问题可能导致某个节点和其它节点失去联系，这时候就形成了P(partition)，也就是由于网络问题，将系统的成员隔离成了2个区域，互相无法知道对方的状态，这在分布式环境下是非常常见的。因为P是必须的，那么我们需要选择的就是A和C。
BASE 理论
核心思想：既是无法做到强一致性（Strong consistency），但每个应用都可以根据自身的业务特点，采用适当的方式来使系统达到最终一致性（Eventual consistency）。
Basically Available 基本可用：假设系统，出现了不可预知的故障，但还是能用，相比较正常的系统而言会有响应时间和功能上的损失：
● 响应时间上的损失：正常情况下的搜索引擎0.5秒即返回给用户结果，而基本可用的搜索引擎可以在2秒作用返回结果。
● 功能上的损失：在一个电商网站上，正常情况下，用户可以顺利完成每一笔订单。但是到了大促期间，为了保护购物系统的稳定性，部分消费者可能会被引导到一个降级页面。
Soft state(软状态)：相对于原子性而言，要求多个节点的数据副本都是一致的，这是一种“硬状态”。软状态指的是：允许系统中的数据存在中间状态，并认为该状态不影响系统的整体可用性，即允许系统在多个不同节点的数据副本存在数据延时。
Eventually consistent(最终一致性)：强一致性读操作要么处于阻塞状态，要么读到的是最新的数据；最终一致性通常是异步完成的，读到的数据刚开始可能是旧数据，但是过一段时间后读到的就是最新的数据。
最终一致性解决方案
对于分布式事务一般采用的都是最终一致性方案，而不是强一致性。而在使用最终一致性的方案时，一定要提到的一个概念是状态机。
前提：
● 可查询操作：业务方需要提供可查询接口，来查询数据信息和状态，供其他服务知道数据状态。
● 幂等操作：同样的参数执行同一个方法，返回的结果都一样。在分布式环境，难免会出现数据的不一致，很多时候为了保证数据的一致，我们都会进行重试。如果不保证幂等，即使重试成功了，也无法保证数据的一致性。我们可以通过业务本身实现实现幂等，比如数据库的唯一索引来约束；也可以缓存（记录）请求和操作结果，当检测到一样的请求时，返回之前的结果。
● 补偿操作：某些数据存在不正常的状态，需要通过额外的方式使数据达到最终一致性的操作。
最终一致性解决方案：
● TCC两阶段补偿方案：TCC是Try-Conﬁrm-Cancel， 比如在支付场景中，先冻结一笔资金，再去发起支付。如果支付成功，则将冻结资金进行实际扣除；如果支付失败，则取消资金冻结。
  ○ Try阶段：完成业务检查（一致性），预留业务资源（准隔离性）
  ○ Conﬁrm阶段：确认执行业务操作，不做任何业务检查，只使用Try阶段 预留的业务资源
  ○ Cancel阶段：取消Try阶段预留的业务资源。Try阶段出现异常时，取消所有业务资源预留请求
● 最大努力通知：最大努力通知事务本质是通过引入定期校验机制来对最终一致性做兜底，对业务侵入性较低；适合于对最终一致性敏感度比较低、业务链路较短的场景。
● 基于可靠性分布式消息队列：阿里的 RocketMq。

缓存与数据库一致性问题
三种缓存读写策略：
● 旁路缓存：比较适合读请求比较多的场景。
  ○ 写：先更新DB，然后删除缓存。
  ○ 读：读缓存，如果无缓存则读DB并更新缓存。
● 读写穿透：服务端把缓存视为主要数据存储，从中读取数据并将数据写入其中。
  ○ 写：先查缓存，存在则更新缓存，再同步更新DB；不存在则只更新DB。
  ○ 读：读缓存，不存在则读DB并更新缓存。
● 写回策略：
● 延迟双删：
  ○ 为什么要延时呢？因为 DB 和 缓存主从节点的数据不是实时同步的，同步数据需要时间。
一致性hash是什么？
一般情况下，我们会使用hash表的方式以key-value的方式来存储数据，但是当数据量比较大的时候，我们就会把数据存储到多个节点上，然后通过hash取模的方法来决定当前key存储到哪个节点上。这种方式有一个非常明显的问题，就是当存储节点增加或者减少的时候，原本的映射关系就会发生变化。也就是需要对所有数据按照新的节点数量重新映射一遍，这个涉及到大量的数据迁移和重新映射，迁移代价很大。而一致性hash就是用来优化这种动态变化场景的算法，它的具体工作原理也很简单。首先，一致性Hash是通过一个Hash环的数据结构来实现的，然后我们把存储节点作为key进行hash之后，会在Hash环上确定一个位置。然后这个目标key会按照顺时针的方向找到离自己最近的一个节点进行数据存储。假设现在需要新增一个节点node4，那数据的映射关系的影响范围只限于node3和node1，只有少部分的数据需要重新映射迁移就行了。
一致性hash算法的好处是扩展性很强，在增加或者减少服务器的时候，数据迁移范围比较小。
另外，在一致性Hash算范里面，为了避免hash倾斜导致数据分配不均匀的情况，我们可以使用虚拟节点的方式来解决。