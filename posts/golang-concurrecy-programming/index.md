# Go 并发编程实战指南：从理论到性能优化


{{&lt; admonition type=abstract title=&#34;导语&#34; open=true &gt;}}
并发编程是 Go 语言最引以为豪的特性之一，但如何正确、高效地使用并发特性却是每个 Go 开发者必须面对的挑战。本文将带你深入探索 Go 并发编程的核心机制，从锁的选择到协程的生命周期管理，通过实战案例和性能数据，帮你掌握并发编程的精髓。无论是构建高并发服务还是优化性能瓶颈，这篇文章都能给你带来实用的指导。
{{&lt; /admonition &gt;}}

&lt;!--more--&gt;

## 读写锁与互斥锁

Go 语言标准库 `sync` 提供了 2 种锁，互斥锁(`sync.Mutex`)和读写锁(`sync.RWMutex`)。

### 互斥锁(`sync.Mutex`)

互斥锁(`sync.Mutex`)，不可同时被多个协程持有，当一个协程获取到互斥锁后，其他协程只能等待该协程释放锁后才能获取锁。
Go 语言标准库 `sync` 提供了 `sync.Mutex` 类型，它有两个方法：`Lock()` 和 `Unlock()`，分别用于获取和释放锁。
可以在代码前调用 `Lock()` 方法获取锁，在代码后调用 `Unlock()` 方法释放锁，也可以使用 `defer` 语句在函数退出时自动释放锁(可以保证互斥锁一定会被释放)。

{{&lt; admonition info &#34;互斥锁如何实现公平？&#34; true&gt;}}
互斥锁有两种状态：正常状态和饥饿状态。

在正常状态下，所有等待锁的 goroutine 按照FIFO顺序等待。
唤醒的 goroutine 不会直接拥有锁，而是会和新请求锁的 goroutine 竞争锁的拥有。
新请求锁的 goroutine 具有优势：它正在 CPU 上执行，而且可能有好几个，所以刚刚唤醒的 goroutine 有很大可能在锁竞争中失败。在这种情况下，这个被唤醒的 goroutine 会加入到等待队列的前面。 如果一个等待的 goroutine 超过 1ms 没有获取锁，那么它将会把锁转变为饥饿模式。

在饥饿模式下，锁的所有权将从 unlock 的 goroutine 直接交给交给等待队列中的第一个。新来的 goroutine 将不会尝试去获得锁，即使锁看起来是 unlock 状态, 也不会去尝试自旋操作，而是放在等待队列的尾部。

如果一个等待的 goroutine 获取了锁，并且满足一以下其中的任何一个条件：(1)它是队列中的最后一个；(2)它等待的时候小于1ms。它会将锁的状态转换为正常状态。

正常状态有很好的性能表现，饥饿模式也是非常重要的，因为它能阻止尾部延迟的现象。
{{&lt; /admonition &gt;}}

### 读写锁(`sync.RWMutex`)

为保证读操作的安全，只要保证并发读时没有写操作即可。
在这种场景下，允许同时有多个协程获取读锁，但是只能有一个协程获取写锁，写锁会阻塞其他读锁和写锁，因此也被称为 `多读单写锁(multiple readers, single writer lock)`，简称读写锁(`sync.RWMutex`)。

Go 语言标准库 `sync` 提供了 `sync.RWMutex` 类型及其四种方法：`RLock()`、`RUnlock()`、`Lock()`、`Unlock()`，分别用于获取和释放读锁和写锁。

&gt; **读写锁的存在是为了解决读多写少的性能问题**：读场景较多时，读写锁可有效减少锁阻塞的时间。

{{&lt; admonition type=tip title=&#34;`sync.Mutex` 与 `sync.RWMutex` 性能对比&#34; open=true &gt;}}
读写操作耗时 1 微秒：

- 读写比为 9:1 时，`sync.RWMutex` 性能约为 `sync.Mutex` 的 8 倍
- 读写比为 1:9 时，`sync.RWMutex` 与 `sync.Mutex` 性能相当
- 读写比为 1:1 时，`sync.RWMutex` 性能约为 `sync.Mutex` 的 2 倍

读写操作耗时 0.1 微秒：`sync.RWMutex` 性能优势下降到 3 倍

读写操作耗时 10 微秒：`sync.RWMutex` 的性能与 1 微秒时基本一致
{{&lt; /admonition &gt;}}

## 协程超时返回

超时控制在网络编程中时非常常见的，利用 `context.WithTimeout` 和 `time.After` 可以轻松实现超时返回。

### `time.After` 实现超时控制

```Go
// ziwi.go
func doBadThing(done chan struct{}) {
 time.Sleep(time.Second)

 done &lt;- struct{}{}
}

func timeout(f func(chan struct{})) error {
 done := make(chan struct{})
 go f(done)

 select {
 case &lt;-done:
  log.Println(&#34;done&#34;)
  return nil

 case &lt;-time.After(time.Millisecond):
  return fmt.Errorf(&#34;timeout&#34;)
 }
}
```

```Go
// ziwi_test.go
func test(t *testing.T, f func(chan struct{})) {
 t.Helper()

 for range 1000 {
  _ = timeout(f)
 }

 time.Sleep(2 * time.Second)
 t.Log(runtime.NumGoroutine())
}

func TestBadTimeout(t *testing.T) { test(t, doBadThing) }
```

在这个典型的 `time.After` 实现超时返回的例子中：

- 利用 `time.After` 启动一个异步的定时器，返回一个 channel：当超过指定时间后，该 channel 将会收到信号
- 启动子协程函数 `f`，函数执行结束后，将向 channel `done` 发送结束信号
- 使用 `select` 阻塞等待 `done` 或 `time.After` 的信息：若超时，则返回错误；若没超时，则返回 `nil`

如果 `f` 调用能在超时前正常退出，那么启动的子协程（goroutine）将能够正常退出。
然而在发生超时的场景下，测试程序输出如下：

```Shell
$ go test -run ^TestBadTimeout$ . -v
=== RUN   TestBadTimeout
    ziwi_test.go:20: 1002
--- PASS: TestBadTimeout (3.26s)
PASS
ok      ziwi    3.950s
```

不难发现，最终程序存在 1002 个协程，说明在主协程退出前，即使 1000 个子协程都执行完成，但子协程并没有正常退出，原因如下：
**当超时发生时，select 接收到 `time.After` 的超时信号，`done` 则没有了接收方（receiver），
由于没有接受者且无缓冲区，发送者（sender）`done` 会一直阻塞，导致协程不能退出，随着时间的积累，造成内存耗尽，程序崩溃**

#### 解决方案: 创建有缓冲区的 channel

将创建 channel `done` 时，缓冲区设置为 1 =&gt; 即使没有接收方，发送方也不会发生阻塞。

```Go
func timeoutWithBuffer(f func(chan struct{})) error {
 done := make(chan struct{}, 1)
 go f(done)

 select {
 case &lt;-done:
  log.Println(&#34;done&#34;)
  return nil

 case &lt;-time.After(time.Millisecond):
  return fmt.Errorf(&#34;timeout&#34;)
 }
}
```

```Shell
$ go test -run ^TestTimeout . -v
=== RUN   TestTimeoutWithBuffer
    ziwi_test.go:28: 2
--- PASS: TestTimeoutWithBuffer (3.29s)
PASS
ok      ziwi    3.966s
```

#### 解决方案: 使用 select 尝试发送

使用 select 尝试向 channel `done` 发送信号，如果失败，则说明缺少接收者，即超时了，那么直接退出即可。

```Go
func doGoodThing(done chan struct{}) {
 time.Sleep(time.Second)

 select {
 case done &lt;- struct{}{}:
 default:
  return
 }
}
```

```Shell
$ go test -run ^TestGood . -v
=== RUN   TestGoodTimeout
    ziwi_test.go:21: 2
--- PASS: TestGoodTimeout (3.25s)
PASS
ok      ziwi    3.924s
➜
```

---

## Channel 关闭原则

&gt; **一个常用的使用 Go channel 的原则是：不要在数据接收方或在有多个发送者的情况下关闭通道，也就是只应该让一个通道唯一的发送者关闭通道**

### 粗鲁关闭（非常不推荐）

如果 channel 已经关闭，再次关闭会产生 Panic，这时通过 `recover` 使程序恢复正常

```go
func SafeClose[T any](ch chan T)(justClosed bool) {
    defer func () {
        if recover() != nil {
            justClosed = false // 一个函数的返回结果可以在 defer 调用中修改
        }
    }()

    close(ch) // 如果 ch 已关闭，则将 Panic
    return true
}
```

---

## channel 忘记关闭

```Go
// ziwi.go
func do(taskCh chan int) {
 for {
  select {
  case t := &lt;-taskCh:
   time.Sleep(time.Millisecond)
   fmt.Printf(&#34;%d &#34;, t)
  }
 }
}

func sendTasks() {
 tashCh := make(chan int)
 go do(taskCh)

 for i := 0; i &lt; 1000; i&#43;&#43; {
  tashCh &lt;- i
 }
}

// ziwi_test.go
func TestDo(t *testing.T) {
 t.Log(runtime.NumGoroutine())
 sendTasks()
 time.Sleep(time.Second)
 t.Log(runtime.NumGoroutine())
}
```

```Shell
$ go test -run ^TestDo$ . -v
=== RUN   TestDo
    ziwi_test.go:33: 2
    ziwi_test.go:36: 3
--- PASS: TestDo (2.14s)
PASS
ok      ziwi    3.231s
```

根据测试结果，不难发现，子协程多了一个，即有一个协程没有得到释放。
显然，这个子协程是 `sendTasks` 中的 `go do(taskCh)`，它一直处于阻塞状态，等待接收任务，直到程序结束也没有释放。

### 解决方案

```Go
func doCheckClose(taskCh chan int) {
 for {
  select {
  case t, beforeClosed := &lt;-taskCh:
   if !beforeClosed {
    fmt.Println(&#34;closed&#34;)
    return
   }

   time.Sleep(time.Millisecond)
   fmt.Printf(&#34;%d &#34;, t)
  }
 }
}

func sendTasksCheckClose() {
 taskCh := make(chan int)
 go doCheckClose(taskCh)

 for i := 0; i &lt; 1000; i&#43;&#43; {
  taskCh &lt;- i
 }

 close(taskCh)
}
```

- `t, beforeClosed := &lt;-taskCh`：判断 channel 是否已经关闭，`beforeClosed` 为 false 表示 channel 已被关闭 =&gt; 不再阻塞等待，直接返回，协程退出
- `sendTasks` 函数中，任务发送结束之后，使用 `close(taskCh)` 将 channel taskCh 关闭

{{&lt; admonition type=tip title=&#34;关于 channel 与 Goroutine 的垃圾回收&#34; open=true &gt;}}
一个通道被其发送数据协程队列和接收数据协程队列中的所有协程引用着。
因此，如果一个通道的这两个队列只要有一个不为空，则此通道肯定不会被垃圾回收。
另一方面，如果一个协程处于一个通道的某个协程队列之中，则此协程也肯定不会被垃圾回收，即使此通道仅被此协程所引用。
事实上，一个协程只有在退出后才能被垃圾回收。
{{&lt; /admonition &gt;}}

---

## 常见问题

### 1. Kill goroutine 可能吗？

答案是：不能。
goroutine 只能自己退出，而不能被其他 goroutine 强制关闭或杀死。

&gt; goroutine 被设计为不可以从外部无条件地结束掉，只能通过 channel 来与它通信。
&gt; 也就是说，每一个 goroutine 都需要承担自己退出的责任。
&gt; (A goroutine cannot be programmatically killed.
&gt; It can only commit a cooperative suicide.)

Github 讨论：[question: is it possible to a goroutine immediately stop another goroutine?](https://github.com/golang/go/issues/32610)

由于 Goroutine 不能被强制关闭或杀死，在超时或其他类似场景下，为了 Goroutine 尽可能正常退出，建议如下：

- 尽量使用非阻塞 I/O（非阻塞 I/O 常用来实现高性能的网络库），阻塞 I/O 很可能导致 goroutine 在某个调用一直等待，而无法正确结束
- 业务逻辑总是考虑退出机制，避免死循环
- 任务分段执行，超时后即时退出，避免 goroutine 无用的执行过多，浪费资源


---

> Author: [kyden](https://github.com/kydance)  
> URL: http://kyden.us.kg/posts/golang-concurrecy-programming/  

