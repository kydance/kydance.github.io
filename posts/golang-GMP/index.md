# 深入理解 Golang GMP


{{&lt; admonition type=abstract title=&#34;导语&#34; open=true &gt;}}
**Go 的调度流程本质上是一个生产-消费流程.**
{{&lt; /admonition &gt;}}

&lt;!--more--&gt;

## I. Process、Thread、Co-routine

### 进程 Process

在操作系统中，进程使用**进程控制块 (PCB, Process Control Block)** 数据结构 `task_struct` 来描述，PCB 是进程存在的唯一标识。

- 进程是指在系统中正在运行的一个应用程序，程序一旦运行就是进程；
- 进程可以认为是程序执行的一个实例，进程是系统进行资源分配的最小单位，且每个进程拥有独立的地址空间；
- 一个进程无法直接访问到另一个进程的变量和数据结构，如果希望一个进程去访问另一个进程的资源，需要使用进程间的通信，如`fifo`、`pipe`、`signal`、`socket` 等；
- 进程调度算法：先来先服务调度算法、短作业优先调度算法、最短剩余作业优先调度算法、最高响应比优先调度算法、最高优先级优先调度算法、时间片轮转算法（公平调度，$20 - 50 ms$）、多级反馈队列调度算法($最高优先级 &#43; 时间片轮转$)；

---

### 线程 Thread

**用户态线程**，是基于用户态的线程管理库来实现的，**线程控制块 (Thread Control Block)** 也是在库里实现，操作系统只能看到整个进程的PCB，即进程与用户线程属于**多对一**的模型。

**内核态线程(Thread)**，是由操作系统管理，对应的 TCB 存储在操作系统里，且其创建、销毁、调度都由操作系统完成；

**轻量级线程 LWP(Light-weight process)**，是由内核支持的用户线程，一个进程可以有一个或多个 LWP，每个 LWP 是跟内核线程一对一映射的，即 LWP 都是由一个内核线程支持，而且 LWP 是由内核管理并像普通进程一样被调度。
**在大多数系统中，LWP 和 普通进程的区别在于，LWP 只有一个最小的执行上下文和调度程序所需的统计信息。**

- 线程是进程的一个实体，是进程的一条执行路径；
- 线程是比进程更小的独立运行的基本单位
- **一个程序至少存在一个进程，一个进程可以有多个($&gt;=1$)线程**

&gt; **进程与线程的区别**
&gt;
&gt; - 进程是资源（包括内存、打开的文件等）分配的单位，线程是 CPU 调度的单位；
&gt; - 进程拥有一个完整的资源平台，而线程只独享必不可少的资源，如寄存器和栈；
&gt; - **同一进程的线程共享本进程的地址空间，而进程之间则是独立的地址空间**；
&gt; - **同一进程内的线程共享本地的资源，但是进程之间的资源是独立的**；
&gt; - **一个进程崩溃后，在保护模式下不会对其他进程产生影响，但是一个线程崩溃整个进程崩溃，即多进程比多线程健壮**；
&gt; - 进程切换，消耗的资源大（主要是虚拟地址空间的切换开销），线程同样具有就绪、阻塞、执行三种基本状态，同样具有状态之间的转换关系；
&gt; - 多进程、多线程都可以并发执行，线程能减少并发执行的时间和空间开销；
&gt; - 每个独立的进程有一个程序入口、程序出口；线程不能独立运行，必须依存于应用程序中，有应用程序提供多个线程执行控制；

---

### 协程 Co-routine

**协程**，又称 &#34;微线程&#34;，表现为一个可以 suspend 和 resume 的函数。

实现协程的关键点：**在于如何保存、恢复和切换上下文**，协程切换只涉及基本的CPU上下文切换（CPU寄存器）.

所有的协程共用的都是一个栈，即系统栈，也就也不必我们自行去给协程分配栈，因为是函数调用，我们当然也不必去显示的保存寄存器的值；

#### Co-routine 分类

**有栈 (stackful) 协程**：实现类似于内核态线程的实现，不同协程的切换还是要切换对应的栈上下文，只是不用陷入内核，例如 goroutine、libco

**无栈 (stackless) 协程**：无栈协程的上下文都会放到公共内存中，在协程切换时使用状态机来切换，而不用切换对应的上下文（都已经在堆中），相比有栈协程更轻量，例如 C&#43;&#43;20、Rust、JavaScript；**==本质就是一个状态机（state machine），即同一协程协程的切换本质不过是指令指针寄存器的改变==**

#### Co-routine 特点

- 一个线程可以有多个协程；协程不是被操作系统内核管理，而是完全由程序控制；
- 协程的开销远远小于线程；协程拥有自己的寄存器上下文和栈，在进行协程调度时，将寄存器上下文和栈保存到其他地方，在切换回来时恢复先前保存的寄存器上下文和栈；
- 每个协程表示一个执行单元，有自己的本地数据，与其他协程共享全局数据和其他资源；
- 跨平台、跨体系架构、无需线程上下文切换的开销、方便切换控制流，简化编程模型；
- 协程的执行效率极高，和多线程相比，线程数量越多，协程的性能优势越明显；

---

## II. GMP

Golang 为提供更加容易使用的并发工具，基于 GMP 模型实现了 goroutine 和 channel。

Goroutine 属于 Co-routine 的概念，非常轻量，一个 goroutine 初始空间只占几 KB 且可伸缩，使得在有限空间内支持大量 goroutine 并发。

Channel 可以独立创建和存取，在不同的 Goroutine 中传递使用，作为队列，遵循 FIFO 原则，同时保证同一时刻只能有一个 goroutine 访问。
channel 作为一种引用类型，声明时需要指定传输数据类型，声明形式如下(`T` 是 channel 可传输的数据类型)：

```go
// 声明 channel
var ch chan T	// 双向 channel
var ch chan&lt;- T	// 只能发送 msg 的 channel
var ch &lt;-chan T 	// 只能接收 msg 的 channel

// 创建 channel
ch := make(chan T, capicity)	// 双向 channel
ch := make(chan&lt;- T, capicity)	// 只能发送 msg 的 channel
ch := make(&lt;-chan T, capicity)	// 只能接收 msg 的 channel

// 访问 channel
ch &lt;- msg	// 发送 msg
msg := &lt;-ch	// 接收 msg
msg, ok := &lt;-ch // 接收 msg，同时判断 channel 是否接收成功
close(ch)	// 关闭 channel
```

---

### Golang 调度

#### 调度组件

- G：Goroutine，一个计算任务. 由需要执行的代码和其上下文组成，上下文包括：当前代码位置、栈空间(初始2K，可增长)、状态等。
- M：Machine，系统线程，执行实体。与 C 语言中的线程相同，通过 `clone` 创建。
- P: Processor，虚拟处理器，包含了 G 运行所需的资源，因此 M 必须获得 P 才能执行代码，否则必须陷入休眠（后台监控线程除外）。可理解为一种 token，有这个 token，才有在物理 CPU 核心上执行的权限。

相关数据结构定义如下：

`g` 的数据结构：

```Go
type g struct {
	// Stack parameters.
	// stack describes the actual stack memory: [stack.lo, stack.hi).
	// stackguard0 is the stack pointer compared in the Go stack growth prologue.
	// It is stack.lo&#43;StackGuard normally, but can be StackPreempt to trigger a preemption.
	// stackguard1 is the stack pointer compared in the //go:systemstack stack growth prologue.
	// It is stack.lo&#43;StackGuard on g0 and gsignal stacks.
	// It is ~0 on other goroutine stacks, to trigger a call to morestackc (and crash).
	stack       stack   // offset known to runtime/cgo
	stackguard0 uintptr // offset known to liblink
	stackguard1 uintptr // offset known to liblink

	_panic    *_panic // innermost panic - offset known to liblink
	_defer    *_defer // innermost defer
	m         *m      // current m; offset known to arm liblink
	sched     gobuf
	syscallsp uintptr // if status==Gsyscall, syscallsp = sched.sp to use during gc
	syscallpc uintptr // if status==Gsyscall, syscallpc = sched.pc to use during gc
	stktopsp  uintptr // expected sp at top of stack, to check in traceback
	// param is a generic pointer parameter field used to pass
	// values in particular contexts where other storage for the
	// parameter would be difficult to find. It is currently used
	// in four ways:
	// 1. When a channel operation wakes up a blocked goroutine, it sets param to
	//    point to the sudog of the completed blocking operation.
	// 2. By gcAssistAlloc1 to signal back to its caller that the goroutine completed
	//    the GC cycle. It is unsafe to do so in any other way, because the goroutine&#39;s
	//    stack may have moved in the meantime.
	// 3. By debugCallWrap to pass parameters to a new goroutine because allocating a
	//    closure in the runtime is forbidden.
	// 4. When a panic is recovered and control returns to the respective frame,
	//    param may point to a savedOpenDeferState.
	param        unsafe.Pointer
	atomicstatus atomic.Uint32
	stackLock    uint32 // sigprof/scang lock; TODO: fold in to atomicstatus
	goid         uint64
	schedlink    guintptr
	waitsince    int64      // approx time when the g become blocked
	waitreason   waitReason // if status==Gwaiting

	preempt       bool // preemption signal, duplicates stackguard0 = stackpreempt
	preemptStop   bool // transition to _Gpreempted on preemption; otherwise, just deschedule
	preemptShrink bool // shrink stack at synchronous safe point

	// asyncSafePoint is set if g is stopped at an asynchronous
	// safe point. This means there are frames on the stack
	// without precise pointer information.
	asyncSafePoint bool

	paniconfault bool // panic (instead of crash) on unexpected fault address
	gcscandone   bool // g has scanned stack; protected by _Gscan bit in status
	throwsplit   bool // must not split stack
	// activeStackChans indicates that there are unlocked channels
	// pointing into this goroutine&#39;s stack. If true, stack
	// copying needs to acquire channel locks to protect these
	// areas of the stack.
	activeStackChans bool
	// parkingOnChan indicates that the goroutine is about to
	// park on a chansend or chanrecv. Used to signal an unsafe point
	// for stack shrinking.
	parkingOnChan atomic.Bool
	// inMarkAssist indicates whether the goroutine is in mark assist.
	// Used by the execution tracer.
	inMarkAssist bool
	coroexit     bool // argument to coroswitch_m

	raceignore    int8  // ignore race detection events
	nocgocallback bool  // whether disable callback from C
	tracking      bool  // whether we&#39;re tracking this G for sched latency statistics
	trackingSeq   uint8 // used to decide whether to track this G
	trackingStamp int64 // timestamp of when the G last started being tracked
	runnableTime  int64 // the amount of time spent runnable, cleared when running, only used when tracking
	lockedm       muintptr
	sig           uint32
	writebuf      []byte
	sigcode0      uintptr
	sigcode1      uintptr
	sigpc         uintptr
	parentGoid    uint64          // goid of goroutine that created this goroutine
	gopc          uintptr         // pc of go statement that created this goroutine
	ancestors     *[]ancestorInfo // ancestor information goroutine(s) that created this goroutine (only used if debug.tracebackancestors)
	startpc       uintptr         // pc of goroutine function
	racectx       uintptr
	waiting       *sudog         // sudog structures this g is waiting on (that have a valid elem ptr); in lock order
	cgoCtxt       []uintptr      // cgo traceback context
	labels        unsafe.Pointer // profiler labels
	timer         *timer         // cached timer for time.Sleep
	selectDone    atomic.Uint32  // are we participating in a select and did someone win the race?

	coroarg *coro // argument during coroutine transfers

	// goroutineProfiled indicates the status of this goroutine&#39;s stack for the
	// current in-progress goroutine profile
	goroutineProfiled goroutineProfileStateHolder

	// Per-G tracer state.
	trace gTraceState

	// Per-G GC state

	// gcAssistBytes is this G&#39;s GC assist credit in terms of
	// bytes allocated. If this is positive, then the G has credit
	// to allocate gcAssistBytes bytes without assisting. If this
	// is negative, then the G must correct this by performing
	// scan work. We track this in bytes to make it fast to update
	// and check for debt in the malloc hot path. The assist ratio
	// determines how this corresponds to scan work debt.
	gcAssistBytes int64
}
```

`m` 的数据结构：

```Go
type m struct {
	g0      *g     // goroutine with scheduling stack
	morebuf gobuf  // gobuf arg to morestack
	divmod  uint32 // div/mod denominator for arm - known to liblink
	_       uint32 // align next field to 8 bytes

	// Fields not known to debuggers.
	procid        uint64            // for debuggers, but offset not hard-coded
	gsignal       *g                // signal-handling g
	goSigStack    gsignalStack      // Go-allocated signal handling stack
	sigmask       sigset            // storage for saved signal mask
	tls           [tlsSlots]uintptr // thread-local storage (for x86 extern register)
	mstartfn      func()
	curg          *g       // current running goroutine
	caughtsig     guintptr // goroutine running during fatal signal
	p             puintptr // attached p for executing go code (nil if not executing go code)
	nextp         puintptr
	oldp          puintptr // the p that was attached before executing a syscall
	id            int64
	mallocing     int32
	throwing      throwType
	preemptoff    string // if != &#34;&#34;, keep curg running on this m
	locks         int32
	dying         int32
	profilehz     int32
	spinning      bool // m is out of work and is actively looking for work
	blocked       bool // m is blocked on a note
	newSigstack   bool // minit on C thread called sigaltstack
	printlock     int8
	incgo         bool          // m is executing a cgo call
	isextra       bool          // m is an extra m
	isExtraInC    bool          // m is an extra m that is not executing Go code
	isExtraInSig  bool          // m is an extra m in a signal handler
	freeWait      atomic.Uint32 // Whether it is safe to free g0 and delete m (one of freeMRef, freeMStack, freeMWait)
	needextram    bool
	traceback     uint8
	ncgocall      uint64        // number of cgo calls in total
	ncgo          int32         // number of cgo calls currently in progress
	cgoCallersUse atomic.Uint32 // if non-zero, cgoCallers in use temporarily
	cgoCallers    *cgoCallers   // cgo traceback if crashing in cgo call
	park          note
	alllink       *m // on allm
	schedlink     muintptr
	lockedg       guintptr
	createstack   [32]uintptr // stack that created this thread, it&#39;s used for StackRecord.Stack0, so it must align with it.
	lockedExt     uint32      // tracking for external LockOSThread
	lockedInt     uint32      // tracking for internal lockOSThread
	nextwaitm     muintptr    // next m waiting for lock

	mLockProfile mLockProfile // fields relating to runtime.lock contention

	// wait* are used to carry arguments from gopark into park_m, because
	// there&#39;s no stack to put them on. That is their sole purpose.
	waitunlockf          func(*g, unsafe.Pointer) bool
	waitlock             unsafe.Pointer
	waitTraceBlockReason traceBlockReason
	waitTraceSkip        int

	syscalltick uint32
	freelink    *m // on sched.freem
	trace       mTraceState

	// these are here because they are too large to be on the stack
	// of low-level NOSPLIT functions.
	libcall   libcall
	libcallpc uintptr // for cpu profiler
	libcallsp uintptr
	libcallg  guintptr
	syscall   libcall // stores syscall parameters on windows

	vdsoSP uintptr // SP for traceback while in VDSO call (0 if not in call)
	vdsoPC uintptr // PC for traceback while in VDSO call

	// preemptGen counts the number of completed preemption
	// signals. This is used to detect when a preemption is
	// requested, but fails.
	preemptGen atomic.Uint32

	// Whether this is a pending preemption signal on this M.
	signalPending atomic.Uint32

	// pcvalue lookup cache
	pcvalueCache pcvalueCache

	dlogPerM

	mOS

	chacha8   chacha8rand.State
	cheaprand uint64

	// Up to 10 locks held by this m, maintained by the lock ranking code.
	locksHeldLen int
	locksHeld    [10]heldLockInfo
}
```

`p` 的数据结构：

```Go
type p struct {
	id          int32
	status      uint32 // one of pidle/prunning/...
	link        puintptr
	schedtick   uint32     // incremented on every scheduler call
	syscalltick uint32     // incremented on every system call
	sysmontick  sysmontick // last tick observed by sysmon
	m           muintptr   // back-link to associated m (nil if idle)
	mcache      *mcache
	pcache      pageCache
	raceprocctx uintptr

	deferpool    []*_defer // pool of available defer structs (see panic.go)
	deferpoolbuf [32]*_defer

	// Cache of goroutine ids, amortizes accesses to runtime·sched.goidgen.
	goidcache    uint64
	goidcacheend uint64

	// Queue of runnable goroutines. Accessed without lock.
	runqhead uint32
	runqtail uint32
	runq     [256]guintptr
	// runnext, if non-nil, is a runnable G that was ready&#39;d by
	// the current G and should be run next instead of what&#39;s in
	// runq if there&#39;s time remaining in the running G&#39;s time
	// slice. It will inherit the time left in the current time
	// slice. If a set of goroutines is locked in a
	// communicate-and-wait pattern, this schedules that set as a
	// unit and eliminates the (potentially large) scheduling
	// latency that otherwise arises from adding the ready&#39;d
	// goroutines to the end of the run queue.
	//
	// Note that while other P&#39;s may atomically CAS this to zero,
	// only the owner P can CAS it to a valid G.
	runnext guintptr

	// Available G&#39;s (status == Gdead)
	gFree struct {
		gList
		n int32
	}

	sudogcache []*sudog
	sudogbuf   [128]*sudog

	// Cache of mspan objects from the heap.
	mspancache struct {
		// We need an explicit length here because this field is used
		// in allocation codepaths where write barriers are not allowed,
		// and eliminating the write barrier/keeping it eliminated from
		// slice updates is tricky, more so than just managing the length
		// ourselves.
		len int
		buf [128]*mspan
	}

	// Cache of a single pinner object to reduce allocations from repeated
	// pinner creation.
	pinnerCache *pinner

	trace pTraceState

	palloc persistentAlloc // per-P to avoid mutex

	// The when field of the first entry on the timer heap.
	// This is 0 if the timer heap is empty.
	timer0When atomic.Int64

	// The earliest known nextwhen field of a timer with
	// timerModifiedEarlier status. Because the timer may have been
	// modified again, there need not be any timer with this value.
	// This is 0 if there are no timerModifiedEarlier timers.
	timerModifiedEarliest atomic.Int64

	// Per-P GC state
	gcAssistTime         int64 // Nanoseconds in assistAlloc
	gcFractionalMarkTime int64 // Nanoseconds in fractional mark worker (atomic)

	// limiterEvent tracks events for the GC CPU limiter.
	limiterEvent limiterEvent

	// gcMarkWorkerMode is the mode for the next mark worker to run in.
	// That is, this is used to communicate with the worker goroutine
	// selected for immediate execution by
	// gcController.findRunnableGCWorker. When scheduling other goroutines,
	// this field must be set to gcMarkWorkerNotWorker.
	gcMarkWorkerMode gcMarkWorkerMode
	// gcMarkWorkerStartTime is the nanotime() at which the most recent
	// mark worker started.
	gcMarkWorkerStartTime int64

	// gcw is this P&#39;s GC work buffer cache. The work buffer is
	// filled by write barriers, drained by mutator assists, and
	// disposed on certain GC state transitions.
	gcw gcWork

	// wbBuf is this P&#39;s GC write barrier buffer.
	//
	// TODO: Consider caching this in the running G.
	wbBuf wbBuf

	runSafePointFn uint32 // if 1, run sched.safePointFn at next safe point

	// statsSeq is a counter indicating whether this P is currently
	// writing any stats. Its value is even when not, odd when it is.
	statsSeq atomic.Uint32

	// Lock for timers. We normally access the timers while running
	// on this P, but the scheduler can also do it from a different P.
	timersLock mutex

	// Actions to take at some time. This is used to implement the
	// standard library&#39;s time package.
	// Must hold timersLock to access.
	timers []*timer

	// Number of timers in P&#39;s heap.
	numTimers atomic.Uint32

	// Number of timerDeleted timers in P&#39;s heap.
	deletedTimers atomic.Uint32

	// Race context used while executing timer functions.
	timerRaceCtx uintptr

	// maxStackScanDelta accumulates the amount of stack space held by
	// live goroutines (i.e. those eligible for stack scanning).
	// Flushed to gcController.maxStackScan once maxStackScanSlack
	// or -maxStackScanSlack is reached.
	maxStackScanDelta int64

	// gc-time statistics about current goroutines
	// Note that this differs from maxStackScan in that this
	// accumulates the actual stack observed to be used at GC time (hi - sp),
	// not an instantaneous measure of the total stack size that might need
	// to be scanned (hi - lo).
	scannedStackSize uint64 // stack size of goroutines scanned by this P
	scannedStacks    uint64 // number of goroutines scanned by this P

	// preempt is set to indicate that this P should be enter the
	// scheduler ASAP (regardless of what G is running on it).
	preempt bool

	// pageTraceBuf is a buffer for writing out page allocation/free/scavenge traces.
	//
	// Used only if GOEXPERIMENT=pagetrace.
	pageTraceBuf pageTraceBuf

	// Padding is no longer needed. False sharing is now not a worry because p is large enough
	// that its size class is an integer multiple of the cache line size (for any of our architectures).
}
```

- 在 `p` 的结构中，`runnext guintptr` 就是 run next，大小为 1，存放下一个将要运行的 G
- 在 `p` 的结构中，`runq [256]guintptr` 就是 local run queue，大小为 256 array，用于存放等待运行的 G

---

#### 调度流程

Go 的调度流程本质上是一个**生产-消费**流程：

{{&lt; figure src=&#34;/posts/golang-GMP/Go-调度本质.svg&#34; title=&#34;&#34; &gt;}}

为了实现简单、高效地调度 Goroutine，Golang 采用了 GMP 模型如下图所示：

{{&lt; figure src=&#34;/posts/golang-GMP/GMP.svg&#34; title=&#34;&#34; &gt;}}

- `global run queue`: 存放等待运行的 G
- `local run queue`: 256 大小的 array，用于存放等待运行的 G
- `runnext`: 存放下一个将要运行的 G

&gt; 由于将 Golang 的调度流程看作**生产者-消费者**流程，因此接下来将分别从生产者、消费者两个方面深入了解。

##### **Goroutine** 的生产端

Goroutine 生产流程：

{{&lt; figure src=&#34;/posts/golang-GMP/Goroutine-Producer.svg&#34; title=&#34;&#34; &gt;}}

##### **Goroutine** 的消费端

&gt; TODO
&gt;
&gt; 关于消费端函数调用链还需完善！！！

Goroutine 消费流程：

{{&lt; figure src=&#34;/posts/golang-GMP/Kyden-blog-Goroutine-Consumer.svg&#34; title=&#34;&#34; &gt;}}

---

### Goroutine 切换成本

`gobuf` 描述了一个 Goroutine 所有现场，从一个 `g` 切换到另一个 `g`，只要把这几个现场字段保存下来，再将 `g` 入队，M 就可以执行其他 `g` 了，无需进入内核态。

`gobuf` 数据结构如下

```Go
type gobuf struct {
	// The offsets of sp, pc, and g are known to (hard-coded in) libmach.
	//
	// ctxt is unusual with respect to GC: it may be a
	// heap-allocated funcval, so GC needs to track it, but it
	// needs to be set and cleared from assembly, where it&#39;s
	// difficult to have write barriers. However, ctxt is really a
	// saved, live register, and we only ever exchange it between
	// the real register and the gobuf. Hence, we treat it as a
	// root during stack scanning, which means assembly that saves
	// and restores it doesn&#39;t need write barriers. It&#39;s still
	// typed as a pointer so that any other writes from Go get
	// write barriers.
	sp   uintptr
	pc   uintptr
	g    guintptr
	ctxt unsafe.Pointer
	ret  uintptr
	lr   uintptr
	bp   uintptr // for framepointer-enabled architectures
}
```

---

### runtime 可拦截 goroutine 阻塞场景解析

Goroutine 属于协程的一种，因此存在运行态、阻塞态等各种状态。
那么 goroutine 什么情况下会发生阻塞？ 当 goroutine 发生阻塞时，GMP 模型如何应对？

显然，当 goroutine 发生可被 runtime 拦截的阻塞时，GMP 模型并不会阻塞调度循环，
而是把 goroutine 挂起，即让 `g` 先进某个数据结构，待 `ready` 后在继续执行，并不会占用线程，
同时线程会进入 `schedule`，继续消费队列，执行其他的 `g`.

#### 场景 I: 延迟

```Go
package main

import (
	&#34;fmt&#34;
	&#34;time&#34;
)

func main() {
	fmt.Println(&#34;Before: &#34;, time.Now())

	time.Sleep(30 * time.Minute)

	fmt.Println(&#34;After: &#34;, time.Now())
}
```

函数调用链如下：

```Go
time.Sleep -&gt;
  runtime.timeSleep {
    ...
    gp := getg()
    t := gp.timer
    ...
    t.arg = gp
    ...
  } -&gt;
    gopark(resetForSleep, unsafe.Pointer(t), waitReasonSleep, traceBlockSleep, 1)
```

显然，在 `runtime.timeSleep` 函数中，获取到的当前 `g` 被挂在 `runtime.timer.arg` 上，然后被挂起。

---

#### 场景 II: Channel send / recv (`chan` / `select`)

```Go
package main

import (
	&#34;fmt&#34;
	&#34;sync&#34;
	&#34;time&#34;
)

func main() {
	var ch = make(chan int)
	var wg = sync.WaitGroup{}
	wg.Add(2)

	go func(ch chan&lt;- int) {
		defer close(ch)
		defer wg.Done()

		time.Sleep(time.Second)
		ch &lt;- 1
	}(ch)

	go func(ch &lt;-chan int) {
		defer wg.Done()
		val := &lt;-ch
		fmt.Println(val)
	}(ch)

	wg.Wait()
}
```

函数 `ch&lt;-` 调用链如下：

```Go
ch&lt;- -&gt;
  runtime.chansend1 -&gt;
  runtime.chansend {
    ...
    gp := getg()
    mysg := acquireSudog()
    ...
    gp.waiting = mysg
    gp.param = nil
    c.sendq.enqueue(mysg)
    // Signal to anyone trying to shrink our stack that we&#39;re about
    // to park on a channel. The window between when this G&#39;s status
    // changes and when we set gp.activeStackChans is not safe for
    // stack shrinking.
    gp.parkingOnChan.Store(true)
    gopark(chanparkcommit, unsafe.Pointer(&amp;c.lock), waitReasonChanSend, traceBlockChanSend, 2)
    ...
  } -&gt;
  gopark
```

函数 `ch&lt;-` 调用链如下：

```Go
&lt;-ch -&gt;
  runtime.chanrecv1(c *hchan, elem unsafe.Pointer) -&gt;
    runtime.chanrecv(c *hchan, ep unsafe.Pointer, block bool) (selected, received bool) {
      ...
      // no sender available: block on this channel.
      gp := getg()
      mysg := acquireSudog()
      ...
      gp.waiting = mysg
      mysg.g = gp
      ...
      c.recvq.enqueue(mysg)
      // Signal to anyone trying to shrink our stack that we&#39;re about
      // to park on a channel. The window between when this G&#39;s status
      // changes and when we set gp.activeStackChans is not safe for
      // stack shrinking.
      gp.parkingOnChan.Store(true)
      gopark(chanparkcommit, unsafe.Pointer(&amp;c.lock), waitReasonChanReceive, traceBlockChanRecv, 2)
    } -&gt;
      runtime.gopark(unlockf func(*g, unsafe.Pointer) bool, lock unsafe.Pointer, reason waitReason, traceReason traceBlockReason, traceskip int)
```

根据调用链可知，`g` 被封装进 `sudog` 中，然后挂在了 `hchan.sendq` 链表上。

相关数据结构 `sudog`, `sendq` 如下：

`sudog` 的数据结构：

```Go
// sudog (pseudo-g) represents a g in a wait list, such as for sending/receiving
// on a channel.
//
// sudog is necessary because the g ↔ synchronization object relation
// is many-to-many. A g can be on many wait lists, so there may be
// many sudogs for one g; and many gs may be waiting on the same
// synchronization object, so there may be many sudogs for one object.
//
// sudogs are allocated from a special pool. Use acquireSudog and
// releaseSudog to allocate and free them.
type sudog struct {
	// The following fields are protected by the hchan.lock of the
	// channel this sudog is blocking on. shrinkstack depends on
	// this for sudogs involved in channel ops.

	g *g

	next *sudog
	prev *sudog
	elem unsafe.Pointer // data element (may point to stack)

	// The following fields are never accessed concurrently.
	// For channels, waitlink is only accessed by g.
	// For semaphores, all fields (including the ones above)
	// are only accessed when holding a semaRoot lock.

	acquiretime int64
	releasetime int64
	ticket      uint32

	// isSelect indicates g is participating in a select, so
	// g.selectDone must be CAS&#39;d to win the wake-up race.
	isSelect bool

	// success indicates whether communication over channel c
	// succeeded. It is true if the goroutine was awoken because a
	// value was delivered over channel c, and false if awoken
	// because c was closed.
	success bool

	// waiters is a count of semaRoot waiting list other than head of list,
	// clamped to a uint16 to fit in unused space.
	// Only meaningful at the head of the list.
	// (If we wanted to be overly clever, we could store a high 16 bits
	// in the second entry in the list.)
	waiters uint16

	parent   *sudog // semaRoot binary tree
	waitlink *sudog // g.waiting list or semaRoot
	waittail *sudog // semaRoot
	c        *hchan // channel
}
```

`hchan` / `waitq` 的数据结构：

```Go
type hchan struct {
	qcount   uint           // total data in the queue
	dataqsiz uint           // size of the circular queue
	buf      unsafe.Pointer // points to an array of dataqsiz elements
	elemsize uint16
	closed   uint32
	elemtype *_type // element type
	sendx    uint   // send index
	recvx    uint   // receive index
	recvq    waitq  // list of recv waiters
	sendq    waitq  // list of send waiters

	// lock protects all fields in hchan, as well as several
	// fields in sudogs blocked on this channel.
	//
	// Do not change another G&#39;s status while holding this lock
	// (in particular, do not ready a G), as this can deadlock
	// with stack shrinking.
	lock mutex
}

type waitq struct {
	first *sudog
	last  *sudog
}
```

---

#### 场景 III: Net read / recv

```Go
package main

import (
	&#34;net&#34;
)

func main() {
	l, _ := net.Listen(&#34;tcp&#34;, &#34;:6633&#34;)

	for {
		conn, _ := l.Accept()

		go func() {
			defer conn.Close()

			var buf = make([]byte, 512)
			n, _ := conn.Read(buf)
			n, _ = conn.Write(buf)
		}()
	}
}
```

函数 `conn.Read` 调用链如下：

```Go
conn.Read(buf) ---&gt;
  net.(*conn).Read(b []byte) (int, error) ---&gt;
    net.(*netFD).Read(p []byte) (n int, err error) ---&gt;
      poll.(*FD).Read(p []byte) (int, error) ---&gt;
        poll.(*pollDesc).waitRead(isFile bool) error ---&gt;
          poll.(*pollDesc).wait(mode int, isFile bool) error ---&gt;
            runtime.poll_runtime_pollWait(pd *pollDesc, mode int) int ---&gt;
              runtime.netpollblock(pd *pollDesc, mode int32, waitio bool) bool {
                gpp := &amp;pd.rg
                if mode == &#39;w&#39; {
                  gpp = &amp;pd.wg
                }
                ...
                gopark(netpollblockcommit, unsafe.Pointer(gpp), waitReasonIOWait, traceBlockNet, 5)
                ...
              } ---&gt;
                gopark(unlockf func(*g, unsafe.Pointer) bool, lock unsafe.Pointer, reason waitReason, traceReason traceBlockReason, traceskip int)
```

函数 `conn.Write` 调用链如下：

```Go
conn.Write(buf) ---&gt;
  net.(*conn).Write(b []byte) (int, error) ---&gt;
    net.(*netFD).Write(p []byte) (n int, err error) ---&gt;
      poll.(*FD).Write(p []byte) (int, error) ---&gt;
        poll.(*pollDesc).waitWrite(isFile bool) error ---&gt;
          poll.(*pollDesc).wait(mode int, isFile bool) error ---&gt;
            runtime.poll_runtime_pollWait(pd *pollDesc, mode int) int ---&gt;
              runtime.netpollblock(pd *pollDesc, mode int32, waitio bool) bool {
                gpp := &amp;pd.rg
                if mode == &#39;w&#39; {
                  gpp = &amp;pd.wg
                }
                ...
                gopark(netpollblockcommit, unsafe.Pointer(gpp), waitReasonIOWait, traceBlockNet, 5)
                ...
              } ---&gt;
                gopark(unlockf func(*g, unsafe.Pointer) bool, lock unsafe.Pointer, reason waitReason, traceReason traceBlockReason, traceskip int)
```

有关 `net.Conn` 读写详细内容，可参考[Netpoll of Network Program for Golang](https://lutianen.github.io/netpoll-of-network-program-for-golang/)

---

#### 场景 IV: 锁阻塞

```Go
package main

import (
	&#34;fmt&#34;
	&#34;sync&#34;
	&#34;time&#34;
)

var mtx sync.Mutex

func main() {
	go func() {
		mtx.Lock()
		defer mtx.Unlock()

		fmt.Printf(&#34;Start\n&#34;)
		time.Sleep(time.Second * 10)
		fmt.Printf(&#34;End\n&#34;)
	}()

	time.Sleep(time.Second) // Ensure child goroutine gets the mutex before main goroutine

	fmt.Printf(&#34;Try to acquire mutex\n&#34;)
	mtx.Lock()
	fmt.Printf(&#34;Main goroutine\n&#34;)
	mtx.Unlock()
}
```

函数 `mtx.Lock()` 调用链如下：

```Go
mtx.Lock() ---&gt;
	sync.(*Mutex).Lock() ---&gt;
		sync.(*Mutex) lockSlow() ---&gt;
			sync.runtime_SemacquireMutex(s *uint32, lifo bool, skipframes int) ---&gt;
				sync.sync_runtime_SemacquireMutex(addr *uint32, lifo bool, skipframes int) ---&gt;
					runtime.semacquire1(addr *uint32, lifo bool, profile semaProfileFlags, skipframes int, reason waitReason) {
						gp := getg()
						if gp != gp.m.curg {
							throw(&#34;semacquire not on the G stack&#34;)
						}

						// Easy case.
						if cansemacquire(addr) {
							return
						}

						// Harder case:
						//	increment waiter count
						//	try cansemacquire one more time, return if succeeded
						//	enqueue itself as a waiter
						//	sleep
						//	(waiter descriptor is dequeued by signaler)
						s := acquireSudog()
						root := semtable.rootFor(addr)
						...
							// Any semrelease after the cansemacquire knows we&#39;re waiting
							// (we set nwait above), so go to sleep.
							root.queue(addr, s, lifo)
							goparkunlock(&amp;root.lock, reason, traceBlockSync, 4&#43;skipframes)
							...
					} ---&gt;
						goparkunlock(lock *mutex, reason waitReason, traceReason traceBlockReason, traceskip int) ---&gt;
							gopark(unlockf func(*g, unsafe.Pointer) bool, lock unsafe.Pointer, reason waitReason, traceReason traceBlockReason, traceskip int)
```

相关数据结构: `semTable` 表现为大小为 251 的数组，其中 `semTable` 中的每一个元素都是一个具有不同地址的 sudog 平衡树.

这些 sudog 中的每一个都可以依次指向（通过 s.waitlink）等待同一地址的其他 sudog 的链表.

```Go
// Asynchronous semaphore for sync.Mutex.

// A semaRoot holds a balanced tree of sudog with distinct addresses (s.elem).
// Each of those sudog may in turn point (through s.waitlink) to a list
// of other sudogs waiting on the same address.
// The operations on the inner lists of sudogs with the same address
// are all O(1). The scanning of the top-level semaRoot list is O(log n),
// where n is the number of distinct addresses with goroutines blocked
// on them that hash to the given semaRoot.
// See golang.org/issue/17953 for a program that worked badly
// before we introduced the second level of list, and
// BenchmarkSemTable/OneAddrCollision/* for a benchmark that exercises this.
type semaRoot struct {
	lock  mutex
	treap *sudog        // root of balanced tree of unique waiters.
	nwait atomic.Uint32 // Number of waiters. Read w/o the lock.
}

// Prime to not correlate with any user patterns.
const semTabSize = 251

type semTable [semTabSize]struct {
	root semaRoot
	pad  [cpu.CacheLinePadSize - unsafe.Sizeof(semaRoot{})]byte
}
```

{{&lt; figure src=&#34;/posts/golang-GMP/Kyden-blog-semTable.svg&#34; title=&#34;&#34; &gt;}}

### runtime 不可拦截 goroutine 阻塞场景解析

`time.Sleep` / `channel send` / `channel recv` / `select` / `net read` / `net write` / `sync.Mutex` 等阻塞场景可被 runtime 拦截，然而仍存在一些阻塞情况是 runtime 无法拦截的，例如：**在执行 C 代码或阻塞在 syscall 上时，必须占用一个线程**。

---

## III. Sysmon

system monitor，高优先级，在专有线程中执行，不需要绑定 `p`.

---

## IV. Summary

- Runtime 构成：**Scheduler**、**Netpoll**、**内存管理**、**垃圾回收**
- GMP：M - 任务消费者；G - 计算任务；P - 可以使用 CPU 的 token
- GMP 中的队列抽象：P 的本地 runnext 字段 --&gt;&gt; P 的 local run queue --&gt;&gt; global run queue；采用多级队列减少锁竞争
- 调度循环：线程 M 在持有 P 的情况下不断消费运行队列中的 G 的过程
- 处理阻塞：
  - runtime 可以接管的阻塞：
    - channel send / recv，sync.Mutex，net read / write，select，time.Sleep
    - 所有 runtime 可接管的阻塞都是通过 `gopark` / `goparkunlock` 挂起，`goready` 恢复
  - runtime 不可接管的阻塞：syscall，cgo，长时间运行需要剥离 P 执行；
- sysmon：
  - 一个后台高级优先级循环，执行时不需要绑定任何的 P
  - 负责：
    - 检查是否已经没有活动线程，如果是则崩溃
    - 轮询 netpoll
    - 剥离在 syscall 上阻塞的 M 的 P
    - 发信号，抢占已经执行时间过长的 G

---

## V. Q &amp; A

1. 为什么阻塞等待的 goroutine，有时表现为 `g` 有时表现为 `sudog` ？

	- `sudog` (pseudo-g) 表示等待列表中的 `g`，例如用于在 channel 上的 `send`/`recv`.
	- `g` 与同步对象是多对多的关系: 一个 `g` 可以出现在多个等待列表中，因此一个 `g` 可能有多个 `sudog`；
	- 很多 `g` 可能在等待同一个同步对象，因此一个对象可能有很多 `sudog`
	- &gt; 一个 `g` 可能对应多个 `sudog`，比如一个 `g` 会同时 `select` 多个 channel

---

## VI. Reference

- [Golang的协程调度器原理及GMP设计思想](https://www.yuque.com/aceld/golang/srxd6d#0810e304)
- [Golang 生产-消费调度流程: Producer](https://www.figma.com/proto/gByIPDf4nRr6No4dNYjn3e/bootstrap?page-id=242%3A7&amp;node-id=242%3A215&amp;viewport=516%2C209%2C0.07501539587974548&amp;scaling=scale-down-width)
- [Golang 生产-消费调度流程: Consumer](https://www.figma.com/proto/gByIPDf4nRr6No4dNYjn3e/bootstrap?page-id=143%3A212&amp;node-id=143%3A213&amp;viewport=134%2C83%2C0.06213996931910515&amp;scaling=scale-down-width)
- [极端情况下收缩 Go 的线程数](https://xargin.com/shrink-go-threads/)
- [Go Scheduler 变更史](https://github.com/golang-design/history#scheduler)
- [internal/poll/fd_poll_runtime.go](https://github.com/golang/go/blob/release-branch.go1.10/src/internal/poll/fd_poll_runtime.go)
- [internal/poll/fd_unix.go](https://github.com/golang/go/blob/release-branch.go1.10/src/internal/poll/fd_unix.go)
- [net/fd_unix.go](https://github.com/golang/go/blob/release-branch.go1.10/src/net/fd_unix.go)
- [runtime/runtime2.go](https://github.com/golang/go/blob/release-branch.go1.22/src/runtime/runtime2.go)
- [runtime/time.go](https://github.com/golang/go/blob/release-branch.go1.22/src/runtime/time.go)
- [runtime/proc.go](https://github.com/golang/go/blob/release-branch.go1.22/src/runtime/proc.go)
- [runtime/netpoll.go](https://github.com/golang/go/blob/release-branch.go1.22/src/runtime/netpoll.go)
- [runtime/netpoll_epoll.go](https://github.com/golang/go/blob/release-branch.go1.22/src/runtime/netpoll_epoll.go)
- [runtime/sema.go](https://github.com/golang/go/blob/release-branch.go1.22/src/runtime/sema.go)
- [sync/mutex.go](https://github.com/golang/go/blob/release-branch.go1.22/src/sync/mutex.go)
- [time/sleep.go](https://github.com/golang/go/blob/release-branch.go1.22/src/time/sleep.go)


---

> : [kyden](https:github.com/kydance)  
> URL: http://kyden.us.kg/posts/golang-gmp/  

