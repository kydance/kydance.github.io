# Golang Netpoll


{{&lt; admonition type=abstract title=&#34;导语&#34; open=true &gt;}}
深入剖析 Golang 网络编程之 Netpoll，主要涉及 Linux 环境下的 Epoll 初始化、 Go 网络编程基本流程（Listen、Accept、Read、Write）以及netpoll 执行流程
{{&lt; /admonition &gt;}}

&lt;!--more--&gt;

{{&lt; admonition type=tip title=&#34;Golang 源码版本&#34; open=true &gt;}}
本文所涉及的源码版本：[v1.22.3](https://github.com/golang/go/tree/release-branch.go1.22/src)
{{&lt; /admonition &gt;}}

## I. 基础概念

**网络编程，是允许不同计算机上的程序通过网络通信的开发过程，涉及多种协议（HTTP、TCP/IP等）以及不同编程语言的应用**。

### 同步、异步、并发模型

| IO 模型    | 读写操作和阻塞阶段                                           |
| ---------- | :----------------------------------------------------------- |
| 阻塞 IO    | 程序阻塞于读写函数                                           |
| IO 复用    | 程序阻塞于 IO 复用系统调用，但可同时监听多个 IO 事件；对 IO 本身的读写操作是非阻塞的 |
| SIGIO 信号 | 信号触发读写就绪事件，用户程序执行读写操作；程序本身没有阻塞阶段 |
| 异步 IO    | 内核执行读写操作并触发读写完成事件；程序没有阻塞阶段         |

&gt; **主要用于区分内核向应用程序通知的是何种 IO 事件（就绪事件 or 完成事件），以及由谁来完成 IO 读写（应用程序 or 内核）**

#### IO模型中的同步

- **同步** IO 模型，指的是应用程序发起 IO 操作后，必须等待 IO 操作完成后才能继续执行后续的操作，即 IO 操作的结果需要立即返回给应用程序；在此期间，应用程序处于阻塞状态，无法做其他操作。
- 优点：编程模型简单
- 缺点：效率较低（应用程序的执行速度被 IO 操作所限制）

&gt; **对于操作系统内核来说，同步 IO 操作是指在内核处理 IO 请求时需要等待**

#### IO 模型中的异步

- **异步** IO 模型，指的是应用程序发起 IO 操作后，无须等待 IO 操作完成，可以立即进行后续的操作；在此期间，操作系统负责把 IO 操作的结果返回给应用程序；
- 优点：可以充分利用系统资源，提高 IO 操作的效率
- 缺点：编程模型相对复杂

&gt; **对于操作系统内核来说，异步 IO 操作指的是，在内核处理 IO 请求时无需等待，立即返回**

#### 并发模式

&gt; **并发模式，指的是 I/O 处理单元和多个逻辑单元之间协调完成任务的方法**

### Linux Epoll

{{&lt; figure src=&#34;/posts/golang-netpoll/Epoll-Linux.svg&#34; title=&#34;&#34; &gt;}}

- epoll 在内核里使用**红黑树(Red-black tree)来跟踪进程所有待检测的文件描述字 `fd`**，把需要监控的 socket 通过 `epoll_ctl()` 函数加入内核中的红黑树里（红黑树是个高效的数据结构，增删改一般时间复杂度是 `O(logn)`）
- epoll 使用**事件驱动**的机制，在内核里**维护了一个链表(List)来记录就绪事件**。
当某个 socket 有事件发生时，内核通过**回调函数**将其加入到这个就绪事件列表中。
当用户调用 `epoll_wait()` 函数时，**只**会返回有事件发生的文件描述符的个数，不需要像 select/poll 那样轮询扫描整个 socket 集合，大大提高了检测的效率
- 两种触发模式

  - **Level trigger**：服务器端不断地从 epoll_wait 中苏醒，直到内核缓冲区数据被 read 函数读完才结束
  - **Edge trigger**：服务器端只会从 epoll_wait 中苏醒一次
- 事件宏

  - `EPOLLIN` 表示对应的文件描述符**可读（包括对端 socket 正常关闭）**
  - `EPOLLOUT` 表示对应的文件描述符**可写**
  - `EPOLLPRI` 表示对应的文件描述符**有紧急的数据可读（带外数据）**
  - `EPOLLERR` 表示对应的文件描述符**发生错误**
  - `EPOLLHUP` 表示对应的文件描述符**被挂断**
  - `EPOLLET` 将 EPOLL 设为**边缘触发模式**（默认电平触发）
  - `EPOLLONESHOT` **只监听一次事件**，当监听完这次事件之后，如果还需要继续监听这个 socket 的话，需要再次把这个 socket 加入到内核中的事件注册表中

## II. 应用示例

```Go
package main

import &#34;net&#34;

func main() {
	l, _ := net.Listen(&#34;tcp&#34;, &#34;127.0.0.1:2333&#34;)

	for {
		conn, _ := l.Accept()

		go func() {
			defer conn.Close()

			buf := make([]byte, 4096)
			_, _ = conn.Read(buf)

			conn.Write(buf)
		}()
	}
}
```

## III. 相关数据结构

```Go
// src/net/fd_fake.go
// Network file descriptor.
type netFD struct {
	pfd poll.FD

	// immutable until Close
	family      int
	sotype      int
	isConnected bool // handshake completed or use of association with peer
	net         string
	laddr       Addr
	raddr       Addr

	// The only networking available in WASI preview 1 is the ability to
	// sock_accept on a pre-opened socket, and then fd_read, fd_write,
	// fd_close, and sock_shutdown on the resulting connection. We
	// intercept applicable netFD calls on this instance, and then pass
	// the remainder of the netFD calls to fakeNetFD.
	*fakeNetFD
}

// poll.FD`: `src/internal/poll/fd_unix.go
// FD is a file descriptor. The net and os packages use this type as a
// field of a larger type representing a network connection or OS file.
type FD struct {
	// Lock sysfd and serialize access to Read and Write methods.
	fdmu fdMutex

	// System file descriptor. Immutable until Close.
	Sysfd int

	// Platform dependent state of the file descriptor.
	SysFile

	// I/O poller.
	pd pollDesc

	// Semaphore signaled when file is closed.
	csema uint32

	// Non-zero if this file has been set to blocking mode.
	isBlocking uint32

	// Whether this is a streaming descriptor, as opposed to a
	// packet-based descriptor like a UDP socket. Immutable.
	IsStream bool

	// Whether a zero byte read indicates EOF. This is false for a
	// message based socket connection.
	ZeroReadIsEOF bool

	// Whether this is a file rather than a network socket.
	isFile bool
}

// Addr represents a network end point address.
//
// The two methods [Addr.Network] and [Addr.String] conventionally return strings
// that can be passed as the arguments to [Dial], but the exact form
// and meaning of the strings is up to the implementation.
type Addr interface {
	Network() string // name of the network (for example, &#34;tcp&#34;, &#34;udp&#34;)
	String() string  // string form of address (for example, &#34;192.0.2.1:25&#34;, &#34;[2001:db8::1]:80&#34;)
}

// fdMutex is a specialized synchronization primitive that manages
// lifetime of an fd and serializes access to Read, Write and Close
// methods on FD.
type fdMutex struct {
	state uint64
	rsema uint32
	wsema uint32
}

type SysFile struct {
	// Writev cache.
	iovecs *[]syscall.Iovec
}

type pollDesc struct {
	runtimeCtx uintptr
}
```

通过源码可以看到，Golang 网络编程涉及到的 `netFD`, `poll.FD`, `Addr`, `SysFile` 以及 `pollDesc` 之间的关系如下：

- `fdmu` 是为了保证对同一个文件的读、写操作能分别被序列化
- `Sysfd` 就是操作系统中 `syscall` 返回的 fd 值
- `pd`，`pollDesc` I/O poller，是 Go 对 poll 过程的一个抽象，所有平台的抽象都是一样的
- `csema`，当文件被关闭时会被触发
- `isBlocking` 表明 FD 是否为 blocking 模式
- `IsStream` 标志该 FD 是否是流式，与流式相反的是基于 packet 的，即 UDP socket
- `ZeroReadIsEOF`，当连接读到 0 长度时，用来区分是否代表 EOF. 如果是基于 packet 的 socket 连接，则始终是 `false`
- `isFile` 标志该 FD 是否代表文件，还是网络连接
- **`netFD` 结构中包含一个 `poll.FD` 类型的成员 `pfd`，以及 `Addr` 接口类型的 `laddr` 和 `raddr`**
- **`poll.FD` 结构含有 `SysFile` 和 `pollDesc` 类型的成员，以及 `fdMutex` 类型的 `fdmu`**

## IV. TCP 网络编程基本流程

本部分涉及众多函数调用，为了描述清晰，采用了图的形式，其中，每一块第一行表示该块所表示的函数名称，其他部分表示该函数中关键函数调用。

### 创建 TCP socket 并监听: `net.Listen`

{{&lt; figure src=&#34;/posts/golang-netpoll/golang-netpoll-listen.svg&#34; title=&#34;&#34; &gt;}}

&gt; NOTE
&gt;
&gt; **需要注意的是，在执行 `net.(*netFD).listenStream` 之前，由于 `maxListenerBacklog` 函数调用了 `open(&#34;/proc/sys/net/core/somaxconn&#34;)`，则会导致 epoll 底层红黑树的提前创建: `runtime.netpollinit` -&gt; `syscall.EpollCreate1` -&gt; `Syscall6(SYS_EPOLL_CREATE1, uintptr(flags), 0, 0, 0, 0, 0)`**.
&gt;
&gt; 另外，当启用 Timer 时，也存在提前初始化 netpoll 的可能，原因：
&gt; **Timers rely on the network poller**
&gt;
&gt; `time.NewTimer` -&gt; `runtime.startTimer` -&gt; `runtime.addtimer` -&gt; `runtime.doaddtimer` -&gt; `netpollGenericInit()`
&gt;
&gt; ```Go
&gt; // doaddtimer adds t to the current P&#39;s heap.
&gt; // The caller must have locked the timers for pp.
&gt; func doaddtimer(pp *p, t *timer) {
&gt; 	// Timers rely on the network poller, so make sure the poller
&gt; 	// has started.
&gt; 	if netpollInited.Load() == 0 {
&gt; 		netpollGenericInit()
&gt; 	}
&gt;   ...
&gt; }
&gt; ```

### 获取 TCP 连接: `net.(*TCPListener).Accept`

{{&lt; figure src=&#34;/posts/golang-netpoll/golang-netpoll-accept.svg&#34; title=&#34;&#34; &gt;}}

### TCP 连接读数据: `net.(*TCPConn).Read`

{{&lt; figure src=&#34;/posts/golang-netpoll/golang-netpoll-read.svg&#34; title=&#34;&#34; &gt;}}

### TCP 连接写数据: `net.(*TCPConn).Write`

{{&lt; figure src=&#34;/posts/golang-netpoll/golang-netpoll-write.svg&#34; title=&#34;&#34; &gt;}}

## V. netpoll 执行流程: `netpoll`

在调度和 GC 的关键点上都会检查一次 netpoll，确定是否存在 ready 状态的 FD：

- `startTheWorldWithSema`

	```Go
	// reason is the same STW reason passed to stopTheWorld. start is the start
	// time returned by stopTheWorld.
	//
	// now is the current time; prefer to pass 0 to capture a fresh timestamp.
	//
	// stattTheWorldWithSema returns now.
	func startTheWorldWithSema(now int64, w worldStop) int64 {
		assertWorldStopped()

		mp := acquirem() // disable preemption because it can be holding p in a local var
		if netpollinited() {
			list, delta := netpoll(0) // non-blocking
			injectglist(&amp;list)
			netpollAdjustWaiters(delta)
		}
		lock(&amp;sched.lock)

		procs := gomaxprocs
		if newprocs != 0 {
			procs = newprocs
			newprocs = 0
		}
		p1 := procresize(procs)
		sched.gcwaiting.Store(false)
		if sched.sysmonwait.Load() {
			sched.sysmonwait.Store(false)
			notewakeup(&amp;sched.sysmonnote)
		}
		unlock(&amp;sched.lock)

		worldStarted()
		...
	}
	```

- `findrunnable`

	```Go
	// Finds a runnable goroutine to execute.
	// Tries to steal from other P&#39;s, get g from local or global queue, poll network.
	// tryWakeP indicates that the returned goroutine is not normal (GC worker, trace
	// reader) so the caller should try to wake a P.
	func findRunnable() (gp *g, inheritTime, tryWakeP bool) {
		...
		// Poll network until next timer.
		if netpollinited() &amp;&amp; (netpollAnyWaiters() || pollUntil != 0) &amp;&amp; sched.lastpoll.Swap(0) != 0 {
			sched.pollUntil.Store(pollUntil)
			if mp.p != 0 {
				throw(&#34;findrunnable: netpoll with p&#34;)
			}
			if mp.spinning {
				throw(&#34;findrunnable: netpoll with spinning&#34;)
			}
			delay := int64(-1)
			if pollUntil != 0 {
				if now == 0 {
					now = nanotime()
				}
				delay = pollUntil - now
				if delay &lt; 0 {
					delay = 0
				}
			}
			if faketime != 0 {
				// When using fake time, just poll.
				delay = 0
			}
			list, delta := netpoll(delay) // block until new work is available
			...
		}
		...
	}
	```

- `pollWork`

	```Go
	// pollWork reports whether there is non-background work this P could
	// be doing. This is a fairly lightweight check to be used for
	// background work loops, like idle GC. It checks a subset of the
	// conditions checked by the actual scheduler.
	func pollWork() bool {
		if sched.runqsize != 0 {
			return true
		}
		p := getg().m.p.ptr()
		if !runqempty(p) {
			return true
		}
		if netpollinited() &amp;&amp; netpollAnyWaiters() &amp;&amp; sched.lastpoll.Load() != 0 {
			if list, delta := netpoll(0); !list.empty() {
				injectglist(&amp;list)
				netpollAdjustWaiters(delta)
				return true
			}
		}
		return false
	}
	```

- `sysmon`

	```Go
	// Always runs without a P, so write barriers are not allowed.
	//
	//go:nowritebarrierrec
	func sysmon() {
		...
		lock(&amp;sched.sysmonlock)
		// Update now in case we blocked on sysmonnote or spent a long time
		// blocked on schedlock or sysmonlock above.
		now = nanotime()

		// trigger libc interceptors if needed
		if *cgo_yield != nil {
			asmcgocall(*cgo_yield, nil)
		}
		// poll network if not polled for more than 10ms
		lastpoll := sched.lastpoll.Load()
		if netpollinited() &amp;&amp; lastpoll != 0 &amp;&amp; lastpoll&#43;10*1000*1000 &lt; now {
			sched.lastpoll.CompareAndSwap(lastpoll, now)
			list, delta := netpoll(0) // non-blocking - returns list of goroutines
			if !list.empty() {
				// Need to decrement number of idle locked M&#39;s
				// (pretending that one more is running) before injectglist.
				// Otherwise it can lead to the following situation:
				// injectglist grabs all P&#39;s but before it starts M&#39;s to run the P&#39;s,
				// another M returns from syscall, finishes running its G,
				// observes that there is no work to do and no other running M&#39;s
				// and reports deadlock.
				incidlelocked(-1)
				injectglist(&amp;list)
				incidlelocked(1)
				netpollAdjustWaiters(delta)
			}
		}
		...
	}
	```

根据 ready 的事件时 Read 或 Write，分别从 poolDesc 的 rg、wg 上获取该唤醒的 goroutine.
然后将已经 ready 的 goroutine push 到 toRun 链表，并且 toRun 链表最终会从 `netpoll()` 返回，通过 `injectglist` 进入全局队列.

&gt; 相当于每次调度循环都要执行 netpoll，检查频率还是比较高的

```Go
// netpoll checks for ready network connections.
// Returns list of goroutines that become runnable.
// delay &lt; 0: blocks indefinitely
// delay == 0: does not block, just polls
// delay &gt; 0: block for up to that many nanoseconds
func netpoll(delay int64) (gList, int32) {
	if epfd == -1 {
		return gList{}, 0
	}
	var waitms int32
	if delay &lt; 0 {
		waitms = -1
	} else if delay == 0 {
		waitms = 0
	} else if delay &lt; 1e6 {
		waitms = 1
	} else if delay &lt; 1e15 {
		waitms = int32(delay / 1e6)
	} else {
		// An arbitrary cap on how long to wait for a timer.
		// 1e9 ms == ~11.5 days.
		waitms = 1e9
	}
	var events [128]syscall.EpollEvent
retry:
	n, errno := syscall.EpollWait(epfd, events[:], int32(len(events)), waitms)
	if errno != 0 {
		if errno != _EINTR {
			println(&#34;runtime: epollwait on fd&#34;, epfd, &#34;failed with&#34;, errno)
			throw(&#34;runtime: netpoll failed&#34;)
		}
		// If a timed sleep was interrupted, just return to
		// recalculate how long we should sleep now.
		if waitms &gt; 0 {
			return gList{}, 0
		}
		goto retry
	}
	var toRun gList
	delta := int32(0)
	for i := int32(0); i &lt; n; i&#43;&#43; {
		ev := events[i]
		if ev.Events == 0 {
			continue
		}

		if *(**uintptr)(unsafe.Pointer(&amp;ev.Data)) == &amp;netpollBreakRd {
			if ev.Events != syscall.EPOLLIN {
				println(&#34;runtime: netpoll: break fd ready for&#34;, ev.Events)
				throw(&#34;runtime: netpoll: break fd ready for something unexpected&#34;)
			}
			if delay != 0 {
				// netpollBreak could be picked up by a
				// nonblocking poll. Only read the byte
				// if blocking.
				var tmp [16]byte
				read(int32(netpollBreakRd), noescape(unsafe.Pointer(&amp;tmp[0])), int32(len(tmp)))
				netpollWakeSig.Store(0)
			}
			continue
		}

		var mode int32
		if ev.Events&amp;(syscall.EPOLLIN|syscall.EPOLLRDHUP|syscall.EPOLLHUP|syscall.EPOLLERR) != 0 {
			mode &#43;= &#39;r&#39;
		}
		if ev.Events&amp;(syscall.EPOLLOUT|syscall.EPOLLHUP|syscall.EPOLLERR) != 0 {
			mode &#43;= &#39;w&#39;
		}
		if mode != 0 {
			tp := *(*taggedPointer)(unsafe.Pointer(&amp;ev.Data))
			pd := (*pollDesc)(tp.pointer())
			tag := tp.tag()
			if pd.fdseq.Load() == tag {
				pd.setEventErr(ev.Events == syscall.EPOLLERR, tag)
				delta &#43;= netpollready(&amp;toRun, pd, mode)
			}
		}
	}
	return toRun, delta
}

// netpollready is called by the platform-specific netpoll function.
// It declares that the fd associated with pd is ready for I/O.
// The toRun argument is used to build a list of goroutines to return
// from netpoll. The mode argument is &#39;r&#39;, &#39;w&#39;, or &#39;r&#39;&#43;&#39;w&#39; to indicate
// whether the fd is ready for reading or writing or both.
//
// This returns a delta to apply to netpollWaiters.
//
// This may run while the world is stopped, so write barriers are not allowed.
//
//go:nowritebarrier
func netpollready(toRun *gList, pd *pollDesc, mode int32) int32 {
	delta := int32(0)
	var rg, wg *g
	if mode == &#39;r&#39; || mode == &#39;r&#39;&#43;&#39;w&#39; {
		rg = netpollunblock(pd, &#39;r&#39;, true, &amp;delta)
	}
	if mode == &#39;w&#39; || mode == &#39;r&#39;&#43;&#39;w&#39; {
		wg = netpollunblock(pd, &#39;w&#39;, true, &amp;delta)
	}
	if rg != nil {
		toRun.push(rg)
	}
	if wg != nil {
		toRun.push(wg)
	}
	return delta
}
```

## VI. 总结

1. Golang 通过对 Linux 内核提供的 `epoll` 实现进行封装，实现了**同步编程异步执行**的效果，其核心数据结构是 `netFD`，并将 `Sysfd` 与 `pollDesc` 结构绑定。
当某个 `netFD` 产生 `EAGAIN` 错误时，则当前 Goroutine 将会被存储到其对应的 `pollDesc` 中，同时 Goroutine 会 `gopark()`，直至这个 `netFD` 再次发生读写事件，会将此 Goroutine 设置为 ready 并放入 `toRun` 队列等待重新运行，而底层事件通知机制就是 epoll.

2. Golang 中 netpoll 的创建与初始化的可能来源：Timer、读文件、TCP Listen.

3. 如下的调度和 GC 关键函数 `startTheWorldWithSema`、`findrunnable`、`pollWork`、`sysmon` 都会进行 `netpoll` 执行流程，检查是否存在 ready 状态的 FD.

## VII. Reference

- [runtime/netpoll.go](https://github.com/golang/go/blob/release-branch.go1.22/src/runtime/netpoll.go)
- [runtime/netpoll_epoll.go](https://github.com/golang/go/blob/release-branch.go1.22/src/runtime/netpoll_epoll.go)
- [runtime/proc.go](https://github.com/golang/go/blob/release-branch.go1.10/src/runtime/proc.go)
- [net/fd_unix.go](https://github.com/golang/go/blob/release-branch.go1.10/src/net/fd_unix.go)
- [internal/poll/fd_poll_runtime.go](https://github.com/golang/go/blob/release-branch.go1.10/src/internal/poll/fd_poll_runtime.go)
- [internal/poll/fd_unix.go](https://github.com/golang/go/blob/release-branch.go1.10/src/internal/poll/fd_unix.go)


---

> : [kyden](https:github.com/kydance)  
> URL: http://kyden.us.kg/posts/golang-netpoll/  

