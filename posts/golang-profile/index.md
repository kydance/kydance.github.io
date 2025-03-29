# Go 性能优化实战：从 Benchmark 到 Profile 的完整指南


{{&lt; admonition type=abstract title=&#34;导语&#34; open=true &gt;}}
在高并发的生产环境中，性能问题往往在最意想不到的时候出现：CPU 突然飙升、内存悄然泄露、Goroutine 数暴增、接口延迟陡升......如何在这些危机时刻快速定位和解决问题？本文将为你揭示 Go 语言性能优化的完整工具链和方法论，从基准测试的正确姿势，到性能分析工具的熟练应用，再到实战中的优化策略。无论你是在进行性能优化，还是在为未来的性能问题未雨绸缪，这都是一份不可或缺的实战指南。
{{&lt; /admonition &gt;}}

&lt;!--more--&gt;

{{&lt; admonition type=note title=&#34;测试环境的稳定性、一致性&#34; open=true &gt;}}
性能测试的结果在很大程度上受到测试环境的影响，因此，在进行性能测试时应尽可能保持测试环境的稳定和一致。

- 测试机器在测试时，不要执行其他任务，不要与其他人共享硬件资源，不要开启节能模式
- 避免使用虚拟机和云主机：一般情况下，为了尽可能地提高资源利用率，虚拟机和云主机 CPU 和内存一般进行超分配，会导致超分机器的性能表现不稳定

&gt; 超分配是针对硬件资源来说的，商业上对应的就是云主机的超卖。虚拟化技术带来的最大直接收益是服务器整合，通过 CPU、内存、存储、网络的超分配（Overcommitment）技术，最大化服务器的使用率。Linux 上专门有一个指标，Steal Time(st)，用来衡量被虚拟机监视器(Hypervisor)偷去给其它虚拟机使用的 CPU 时间所占的比例。
&gt;
&gt; 例如，虚拟化的技能之一就是随心所欲的操控 CPU，例如一台 32U(物理核心)的服务器可能会创建出 128 个 1U(虚拟核心)的虚拟机，当物理服务器资源闲置时，CPU 超分配一般不会对虚拟机上的业务产生明显影响，但如果大部分虚拟机都处于繁忙状态时，那么各个虚拟机为了获得物理服务器的资源就要相互竞争，相互等待。
{{&lt; /admonition &gt;}}

## I. Benchmark

Go 语言 `testing` 标准库内置支持 Benchmark 测试。

Benchmark 和普通单元测试用例一样，都位于 `_test.go` 文件，并且函数名以 `Benchmark` 开头，参数是 `b *testing.B`.

`go test` 命令默认不运行 Benchmark 测试，需要在命令中加上 `-bench` 参数来进行测试：

- `go test -bench .`: 运行当前 `packge` 内的用例
- `go test -bench &#39;In$&#39; .`: `-bench` 参数支持正则表达式，只有匹配到的用例才会运行
- `go test -bench ./&lt;package name&gt;`: 运行子 `package` 内的用例
- `go test -bench ./...`: 运行当前目录下所有的 `package` 内的用例

### 浅析 Benchmark 工作原理

Benchmark 用例参数 `b *testing.B` 中的 `N` 属性表示该用例需要运行的次数。
一般情况下，不同用例的 `b.N` 是不一样的。

`b.N` 从 1 开始，若该用例能够在 1s 内完成，`b.N` 的值则会增加，再次执行。`b.N` 的值大概以 1, 2, 3, 5, 10, 20, 30, 50, 100 ... 这样的序列增加，越到后面，增加越快。

```bash
➜  ziwi go test -bench=&#39;Fib$&#39; -cpu=2,4 -benchtime=50x -count=2 -benchmem .
goos: darwin
goarch: arm64
pkg: go-temp/ziwi
cpu: Apple M1 Pro
BenchmarkFib-2                50               205.0 ns/op             0 B/op          0 allocs/op
BenchmarkFib-2                50               386.7 ns/op             0 B/op          0 allocs/op
BenchmarkFib-4                50               383.3 ns/op             0 B/op          0 allocs/op
BenchmarkFib-4                50               385.0 ns/op             0 B/op          0 allocs/op
PASS
ok      go-temp/ziwi    0.196s
```

- `BenchmarkFib-10`: `-10` 即 `GOMAXPROCS`，默认等于 CPU 核数，可通过 `-cpu` 参数改变 `GOMAXPROCS`，`-cpu` 支持传入一个列表作为参数
- `6097884` 和 `183.7 ns/op`: 表示该用例执行了 `6097884` 次，每次执行需要花费的时间为 `183.7 ns/op`
- 为了提高性能测试的准确度，可以使用 `-benchtime` 和 `-count` 两个参数分别调整测试时长(默认 1s)和执行轮数。
其中，`-benchtime` 的值除了是时间外，还可以是具体次数：`go test -bench=&#39;Fib$&#39; -benchtime=300x .`
- `-benchmem` 参数可以度量内存分配的次数

### ResetTimer &amp; StopTimer &amp; StartTimer

- `b.ResetTimer()`: 用于将进行 Benchmark 开始前的准备工作所消耗的时间忽略掉
- `b.StopTimer()`: 暂停计时
- `b.StartTimer()`: 开始计时

---

## II. Profile

&gt; 当面对一个未知程序，如何分析这个程序的性能，并找到瓶颈点呢？
&gt;
&gt; **pprof 就是用来解决这个问题的**

### CPU 性能分析

CPU性能分析（CPU profiling）是最常见的性能分析类型，当启动 CPU 性能分析时，运行时（runtime）将每隔 10ms 中断一次，记录此时正在运行的协程（goroutines）的堆栈信息。
程序结束后，可以分析记录的数据找到最热代码路径（hosttest code paths）。

{{&lt; admonition note &#34;What’s the meaning of “hot codepath”&#34; true &gt;}}
Compiler hot paths are code execution paths in the compiler in which most of the execution time is spent, and which are potentially executed very often.

– [What’s the meaning of “hot codepath”](https://english.stackexchange.com/questions/402436/whats-the-meaning-of-hot-codepath-or-hot-code-path)
{{&lt; /admonition &gt;}}

一个函数在性能分析数据中出现的次数越多，说明执行该函数代码路径（code path）花费的时间占总运行时间的比重越大。

### 内存性能分析

内存性能分析（Memory profiling）记录堆内存分配时的堆栈信息，忽略栈内存分配信息，当启动 memory 性能分析时，默认每 1000 次采样 1 次（这个比例可调整）。

由于内存性能分析是基于采样的，因此**基于内存分析数据来判断程序所有的内存使用情况是很困难的**。

### 阻塞性能分析

阻塞性能分析（block profiling）是 Go 特有的，它用来记录一个协程等待一个共享资源花费的时间，因此在判断程序的并发瓶颈时会很有用。

阻塞场景：

- 在没有缓冲区的信道上发送或接收数据
- 从空的信道上接收数据，或发送数据到满的信道上
- 尝试获得一个已经被其他协程锁住的排它锁

{{&lt; admonition tip &#34;When using block profilling&#34; true &gt;}}
一般情况下，当所有的 CPU 和 memory 瓶颈解决后，才会考虑阻塞性能分析。
{{&lt; /admonition &gt;}}

### 实践场景

在进行 **API 压测**、**全链路压测**、**线上生产环境被高峰流量打爆**的过程中随时可能发生故障等问题，例如：

- CPU 占用过高，超过 90%；
- 内存爆掉，[OOM(Out of memory)](https://en.wikipedia.org/wiki/Out_of_memory)；
- Goroutine 数量过多，80W；
- 线程数超高；
- 延迟过高；

在发生以上故障时，一般需要结合 **pprof** 寻找故障原因，并根据不同的情况选择不同的方案；

&gt; 线上一定要具有开启 `pprof` 的能力，如果考虑安全性，也要具有通过配置开启的能力；

### 压测时需要关注的服务指标

- **Request rate**: The number of service requests per second.
- **Errors**: The number of request that failed.
- **Duration**: The time for requests to complete.
- **Goroutine / Thread 数量**: 如果 Goroutine 数量很多，需要关注这些 Goroutine 的执行情况.
- **GC 频率**
- **gctrace 的内容**:
- **GC 的 STW 时间**

还有一些其他 Memstats 相关的其他指标，可以参考 [Prometheus](https://github.com/prometheus/prometheus).

### 压测手段

- [wrk](https://github.com/wg/wrk): a HTTP benchmarking tool
- [wrk2](https://github.com/giltene/wrk2): a HTTP benchmarking tool based mostly on wrk
- [HEY](https://github.com/rakyll/hey): a tiny program that sends some load to a web application.
- [Vegate](https://github.com/tsenart/vegeta): a versatile HTTP load testing tool built out of a need to drill HTTP services with a constant request rate.
- [h2load](https://nghttp2.org/documentation/h2load-howto.html): HTTP/2 benchmarking tool
- [ghz](https://ghz.sh/): gRPC benchmarking and load testing tool

### pprof 应用实例

```Go
package main

import (
    &#34;net/http&#34;
    _ &#34;net/http/pprof&#34;
)

var quit chan struct{} = make(chan struct{})

func f() {
    &lt;- quit
}

func main() {
    go func() { http.ListenAndServe(&#34;:8080&#34;, nil) }()

    for i := 0; i &lt; 10000; i&#43;&#43; {
        go f()
    }

	for {} // Test
}
```

```Bash
go tool pprof -http=:9999 localhost:8080/debug/pprof/heap
```

&gt; 注意事项
&gt;
&gt; 1. 测试代码中引入 `net/http/pprof` 包： `_ &#34;net/http/pprof&#34;`
&gt; 2. 单独启动一个 Goroutine 开启监听(端口自定，例如这里是 8080)：`go func() { http.ListenAndServe(&#34;:8080&#34;, nil) }()`
&gt; 3. `$ go tool pprof -http=:9999 localhost:8080/debug/pprof/heap`

---

## III. Optimize

### 优化方向

{{&lt; figure src=&#34;/posts/golang-profile/优化范围.svg&#34; title=&#34;&#34; &gt;}}

在分析上图的应用程序运行过程，可以发现进行程序优化时，一般从可以从以下方面入手：

- 应用层优化: 主要指的是逻辑优化、内存使用优化、CPU 使用优化、阻塞优化等，并且本层优化效果可能优于底层优化；
- 底层优化：GC优化、Go 标准库优化、Go runtime 优化等

### 基本优化流程

1. **外部依赖**：在监控系统中查看是否存在问题，例如依赖的上游服务 (DB/redis/MQ) 延迟过高；
2. **CPU 占用**：通过查看 CPU profile 检查是否存在问题，优化占用 CPU 较多的部分逻辑；
3. **内存占用**：看 Prometheus，内存 RSS / Goroutine 数量 / Goroutine 栈占用 --&gt;&gt; 如果 Goroutine 数量不多，则重点关注 heap profile 中的 inuse --&gt;&gt; 定时任务类需要看 alloc
4. Goroutine 数量过多 --&gt;&gt; 从 profile 网页进去看看 Goroutine 的执行情况（在干什么？） --&gt;&gt; 检查死锁、阻塞等问题 --&gt;&gt; 个别不在意延迟的选择第三方库优化

### 常见优化场景

#### 字符串拼接

```Go
package main

import (
	&#34;fmt&#34;
	&#34;testing&#34;
)

func BenchmarkConcat0(b *testing.B) {
	var str string

	for i := 0; i &lt; b.N; i&#43;&#43; {
		str = &#34;&#34;
		str &#43;= &#34;userid : &#34; &#43; &#34;1&#34;
		str &#43;= &#34;localtion : &#34; &#43; &#34;ab&#34;
	}
}

func BenchmarkConcat1(b *testing.B) {
	var str string

	for i := 0; i &lt; b.N; i&#43;&#43; {
		str = &#34;&#34;
		str &#43;= fmt.Sprintf(&#34;userid : %v&#34;, &#34;1&#34;)
		str &#43;= fmt.Sprintf(&#34;localtion : %v&#34;, &#34;ab&#34;)
	}
}
```

```Bash
$ go test -bench=. -benchmem
goos: linux
goarch: amd64
pkg: github.com/lutianen/go-test/bench0
cpu: 11th Gen Intel(R) Core(TM) i7-11800H @ 2.30GHz
BenchmarkConcat0-16     35702518                32.86 ns/op           24 B/op          1 allocs/op
BenchmarkConcat1-16      8105732               140.9 ns/op            56 B/op          3 allocs/op
PASS
ok      github.com/lutianen/go-test/bench0      2.506s
```

#### 逃逸分析

用户声明的对象，被放在栈上还是堆上？
可以通过编译器的 escape analysis 来决定 `go build -gcflags=&#34;-m&#34; xxx.go`

```Go
package main

func main() {
	var sl = make([]int, 1024)
	println(sl[0])

	var sl0 = make([]int, 10240)
	println(sl0[0])
}
```

```Bash
$ go build -gcflags=&#34;-m&#34; main.go
# command-line-arguments
./main.go:3:6: can inline main
./main.go:4:15: make([]int, 1024) does not escape
./main.go:7:16: make([]int, 10240) escapes to heap
```

&gt; TODO: 各种逃逸分析的可能性有哪些？

#### Trasval 2-D Matrix

```Go
package bench1

import &#34;testing&#34;

func BenchmarkHorizontal(b *testing.B) {
	arrLen := 10000

	arr := make([][]int, arrLen, arrLen)

	for i := 0; i &lt; arrLen; i&#43;&#43; {
		arrInternal := make([]int, arrLen)
		for j := 0; j &lt; arrLen; j&#43;&#43; {
			arrInternal[j] = 0
		}
        arr[i] = arrInternal
	}

	for i := 0; i &lt; b.N; i&#43;&#43; {
		for x := 0; x &lt; len(arr); x&#43;&#43; {
			for y := 0; y &lt; len(arr); y&#43;&#43; {
				arr[x][y] = 1
			}
		}
	}
}

func BenchmarkVertical(b *testing.B) {
	arrLen := 10000

	arr := make([][]int, arrLen, arrLen)

	for i := 0; i &lt; arrLen; i&#43;&#43; {
		arrInternal := make([]int, arrLen)
		for j := 0; j &lt; arrLen; j&#43;&#43; {
			arrInternal[j] = 0
		}
        arr[i] = arrInternal
	}

	for i := 0; i &lt; b.N; i&#43;&#43; {
		for x := 0; x &lt; len(arr); x&#43;&#43; {
			for y := 0; y &lt; len(arr); y&#43;&#43; {
				arr[y][x] = 1
			}
		}
	}
}
```

```Bash
$ go test -bench=. -benchmem
goos: linux
goarch: amd64
pkg: github.com/lutianen/go-test/bench1
cpu: 11th Gen Intel(R) Core(TM) i7-11800H @ 2.30GHz
BenchmarkHorizontal-16                15          71020410 ns/op        54629717 B/op        666 allocs/op
BenchmarkVertical-16                   1        1059649022 ns/op        819445856 B/op     10002 allocs/op
PASS
ok      github.com/lutianen/go-test/bench1      3.676s
```

#### Zero Garbage / Allocation

Zero Grabage 一般指的是通过利用 `sync.Pool` 将堆分配完全消灭的优化技术。

例如，在 http router 框架 [fasthttp](https://github.com/valyala/fasthttp) 中应用较多.

{{&lt; figure src=&#34;/posts/golang-profile/Fasthttp-best-practices.png&#34; title=&#34;&#34; &gt;}}

#### False Sharing

{{&lt; figure src=&#34;/posts/golang-profile/False-Sharing.svg&#34; title=&#34;&#34; &gt;}}

CPU 运行过程中修改数据是一个 **cache line**为单位，当两个变量`A`/`B`满足以下条件：

- 在内存中相邻
- 并发修改频繁

那么，当 CPU0 修改变量 `A` 时，会导致 CPU1 中的变量 `B` 缓存失效。

解决方法，在定义数据结构中，填充一些 `padding` 用以满足该数据结构正好是 cache line 的整数倍；

```Go
type NoPad struct {
	x uint64
	y uint64
}

type WithPad struct {
	x uint64
	_ [6]uint64
	y uint64
}
```

&gt; 查看 cache line 大小：`cat /sys/devices/system/cpu/cpu&lt;core-num&gt;/cache/index0/coherency_line_size`

#### 降低外部命令调用频次

优化前：

```Go
func f(wr http.ResponseWriter, r *http.Request) {
	uuid, _ := exec.Command(&#34;uuidgen&#34;).Output() // Use exec.Command

	wr.Header()[&#34;Content-Type&#34;] = []string{&#34;application/text&#34;}
	io.WriteString(wr, string(uuid))
}
```

优化后：

```Go
import uuid &#34;github.com/satori/go.uuid&#34;

func f(wr http.ResponseWriter, r *http.Request) {
	uuid, _ := uuid.NewV4() // Replace exec.Command with existing library

	wr.Header()[&#34;Content-Type&#34;] = []string{&#34;application/text&#34;}
	io.WriteString(wr, uuid.String())
}
```

&gt; 总结：
&gt;
&gt; 1. 线上使用 `exec` 命令是非常危险的
&gt; 2. 采用第三方库代替外部命令

#### 阻塞导致高延迟

##### 锁阻塞

```Go
var mtx sync.Mutex
var data = map[string]string{
	&#34;hint&#34;: &#34;hello wold&#34;,
}

func f(wr http.ResponseWriter, r *http.Request) {
	mtx.Lock()
	defer mtx.Unlock()

	buf := data[&#34;hint&#34;]
	time.Sleep(time.Millisecond * 10) // 临界区内的慢操作
	wr.Header()[&#34;Content-Type&#34;] = []string{&#34;application/json&#34;}
	io.WriteString(wr, buf)
}
```

- **减小临界区 - 优化后**：

```Go
var mtx sync.Mutex
var data = map[string]string{
	&#34;hint&#34;: &#34;hello wold&#34;,
}

func f(wr http.ResponseWriter, r *http.Request) {
	mtx.Lock()
	buf := data[&#34;hint&#34;]
	mtx.Unlock()

	time.Sleep(time.Millisecond * 10) // 慢操作放置于临界区之外
	wr.Header()[&#34;Content-Type&#34;] = []string{&#34;application/json&#34;}
	io.WriteString(wr, buf)
}
```

在后端系统开发中，锁瓶颈是较常见的问题，例如文件锁
{{&lt; figure src=&#34;/posts/golang-profile/func-write-with-lock.png&#34; title=&#34;&#34; &gt;}}

- **双 Buffer 完全干掉锁阻塞**

	&gt; 使用双 Buffer / RCU 完全消除读阻塞：全量更新，直接替换原 config

	```Go
	func updateConfig() {
		var newConfig = &amp;MyConfig {
			WhiteList: make(map[int]struct{}),
		}

		// Do a lot of compulation
		for i :=0; i &lt; 1000; i&#43;&#43; {
			newConfig.WhiteList[i] = struct{}{}
		}

		config.Store(newConfig)
	}
	```

	&gt; 使用双 Buffer / RCU 完全消除读阻塞：部分更新，先拷贝原 config，然后更新 key，最后替换

	```Go
	// Partial update
	func updateConfig() {
		var oldConfig = getConfig()
		var newConfig = &amp;MyConfig{
			WhiteList: make(map[int]struct{})
		}

		// Copy from old
		for k,v := range oldConfig.WhiteList {
			newConfig.WhiteList[k] = v
		}

		// Modify some keys
		newConfig.WhiteList[123] = struct{}{}
		newConfig.WhiteList[124] = struct{}{}

		config.Store(newConfig)
	}
	```

	**NOTE: 当更新可能并发时，则需要在更新时加锁**

&gt; 优化锁阻塞瓶颈的手段总结:
&gt;
&gt; 1. 减小临界区：只锁必须锁的对象，临界区内尽量不放慢操作，如 `syscall`
&gt; 2. 降低锁粒度：全局锁 -&gt; 对象锁，全局锁 -&gt; 连接锁， 连接锁 -&gt; 请求锁，文件锁 -&gt; 多个文件各种锁
&gt; 3. 同步改异步：同步日志 -&gt; 异步日志，若队列满则丢弃，不阻塞业务逻辑

#### CPU 使用太高

##### 编解码使用 CPU 过高

通过更换 json 库，就可以提高系统的吞吐量：本质上是请求的 CPU 使用被优化了（可使用固定 QPS 压测来验证）

&gt; `encoding/json` --&gt;&gt; `json &#34;github.com/json-iterator/go&#34;`

##### GC 使用 CPU 过高

- 将变化较少的结构放在堆外，通过 cgo 来管理内存，让 GC 发现不了这些对象，也就不会扫描了
- [**offheap**](https://github.com/glycerine/offheap)，可以减少 Go 进程的内存占用和内存使用波动，但要用到 cgo

[Manual Memory Management in Go using jemalloc](https://dgraph.io/blog/post/manual-memory-management-golang-jemalloc/)

```Go
func BenchmarkMapWithoutPtrs(b *testing.B) {
	for i := 0; i &lt; b.N; i&#43;&#43; {
		var m = make(map[int]int)
		for i := 0; i &lt; 10; i&#43;&#43; {
			m[i] = i
		}
	}
}

func BenchmarkMapWithPtrs(b *testing.B) {
	for i := 0; i &lt; b.N; i&#43;&#43; {
		var m = make(map[int]*int)
		for i := 0; i &lt; 10; i&#43;&#43; {
			var v = i
			m[i] = &amp;v
		}
	}
}
```

```Bash
$ got -bench . -benchmem

BenchmarkMapWithoutPtrs-16       3362536               412.1 ns/op           292 B/op          1 allocs/op
BenchmarkMapWithPtrs-16          2580622               524.8 ns/op           371 B/op         11 allocs/op
```

结论: 当 map 中含有大量的指针 key 时，会给 GC 扫描造成压力

解决方案（**只适用于内存不紧张，且希望提高整体吞吐量的服务**）：

- 调大 GOGC
- 程序启动阶段 make 一个全局超大的 slice（如1GB）*TODO 如何解决的？*

#### 内存占用过高

##### 堆分配导致内存占用过高

```Go
const max = 1 &lt;&lt; 14
//go:noinline
func Steal() {
	var buf = make([]int, max)

	for j := 0; j &lt; max; j&#43;&#43; {
		buf = append(buf, make([]int, max)...)
	}
}

func BenchmarkSteal(b *testing.B) {
	for i := 0; i &lt; b.N; i&#43;&#43; {
		Steal()
	}
}
```

```Bash
$ go test -bench . -benchmem
BenchmarkSteal-16              1        1386661490 ns/op        10764864792 B/op              51 allocs/op
```

##### Goroutine 数量太多导致内存占用过高

**Goroutine 涉及到的占用内存可能如下**：

1. Goroutine 栈占用的内存(**难优化**，一条 TCP 连接至少对应一个 Goroutine)
2. TCP Read Buffer 占用的内存(**难优化**，因为大部分连接阻塞在 Read 上，Read Buffer 基本没有可以释放的时机)

	```Go
	func f() {
		var l net.Listener
		for {
			c, _ := l.Accept()
			go func() {
				var buf = make([]byte, 4096)
				for {
					c.Read(buf)
				}
			}()
		}
	}
	```

3. TCP Writer Buffer 占用的内存(**易优化**，因为活跃连接不多)

&gt; 原因：
&gt;
&gt; 1. `gopark(...)` 的 Goroutine， 占用内存
&gt; 2. 阻塞的 Read Buffer 很难找到时机释放，占用内存

**Solution**: 在一些不太重视延迟的场景中（例如推送系统），可以使用某些库进行优化：evio、gev、gnet、easygo、gaio、netpoll

&gt; NOTE: **一定要进行在真实业务场景中做压测**，不要相信某些库的 README 中的压测数据

#### 常见优化场景总结

1. CPU 使用太高
   - 应用逻辑导致
     - JSON 序列化
       - 使用一些优化的 JSON 库替代标准库
       - 使用二进制编码方式代替 JSON 编码
       - 同物理节点通信，使用共享内存 IPC，直接干掉序列化开销
     - MD5 计算 HASH 值成本太高 --&gt; 使用 [cityhash](https://github.com/google/cityhash), [murmurhash](https://zh.wikipedia.org/zh-cn/Murmur%E5%93%88%E5%B8%8C)
     - 其他应用逻辑：只能具体情况具体分析
   - GC 使用 CPU 过高
     - 减少堆上对象分配
       - `sync.Pool` 进行堆对象重用
       - `Map` -&gt; `slice`
       - 指针 -&gt; 非指针对象
       - 多个小对象 -&gt; 合并为一个大对象
     - offheap
     - 降低 GC 频率
       - 修改 GOGC
       - 在程序开始时 `make` 一个全局大 `slice`
   - 调度相关的函数使用 CPU 过高
     - 尝试使用 Goroutine Pool，减少 Goroutine 的创建与销毁
     - 控制最大 Goroutine 数量
2. 内存使用过高
   - 堆内存占用内存空间过高
     - `sync.Pool` 对象复用
     - 为不同大小的对象提供不同大小 level 的 `sync.Pool`
     - offheap
   - Goroutine 栈占用过多内存
     - 减少 Goroutine 数量
       - 如每个连接一读一写 --&gt;&gt; 合并为一个连接一个 goroutine
       - Goroutine pool 限制最大 goroutine 数量
       - 使用裸 epoll 库(evio, gev等)修改网络编程方式（只适用于对延迟不敏感的业务）
     - 通过修改代码，减少函数调用层级（难）
3. 阻塞问题
   - 上游系统阻塞
     - 让上游赶紧解决
   - 锁阻塞
     - 减少临界区范围
     - 降低锁粒度
       - Global Lock --&gt;&gt; Shareded Lock
       - Global Lock --&gt;&gt; Connection Level Lock
       - Connection Level Lock --&gt;&gt; Request Level Lock
     - 同步改异步
       - 日志场景：同步日志 --&gt;&gt; 异步日志
       - Metrics 上报场景：`select` --&gt;&gt; `select` &#43; `default`
     - 个别场景使用双 Buffer 完全消灭阻塞

---

## IV. Coutinuous Profiling

压测是一个蹲点行为，然而真实场景并不美好，它们通常是难以发现的偶发问题：

- 该到吃饭的时候，CPU 使用尖刺
- 凌晨四点半，系统发生 OOM
- 刚睡着的时候，Goroutine 数量爆炸
- 产品被部署到客户那里，想登陆客户的环境并不方便

此时 Coutinuout Profiling 就派上用场了.

{{&lt; figure src=&#34;/posts/golang-profile/Continuous-Profiling.svg&#34; title=&#34;&#34; &gt;}}

**自省式的 Profile Dumper**，可以根据 CPU 利用率、Memory 利用率、Goroutine 数量等多个指标检测系统，设置定时周期进行检测，当发现某个指标异常时，自动 Dump file.

---

## V. Summary

1. `_pad` 优化，针对**多个线程更新同一个结构体内不同的字段**场景有效，而针对**一个线程同时更新整个结构体**的场景意义不大；

2. 第三方接口出现问题，如何保护自己的服务？
    &gt; 对外部调用必须有超时 ==&gt; 熔断

3. goroutine 初始化栈空间为 2KB，最大 1GB，那么 heap 为什么不爆栈？
    &gt; 在 Go 语言中，goroutine 和 heap 使用单独的内存空间：Goroutine 有自己的堆栈空间，用于存储局部变量、函数帧和其他运行时信息；heap 则是一个共享内存空间，用于存储动态分配的对象，例如 slice、map 和 strings。
    &gt;
    &gt; 当 Goroutine 需要分配的内存多于起堆栈上的可用内存时，它将自动从 stack 中分配内存，采用的是 stack 分配机制完成，运行 goroutine 分配任何数量的内存，而不用担心 stack 空间耗尽；
    &gt; 除了堆分配之外，goroutine 还可以使用一种称为堆栈复制的技术来在它们之间共享数据，堆栈复制比堆分配更有效，但它只能用于共享足够小以适合堆栈的数据。

---

## VI. Reference

- [Benchmarks Game](https://benchmarksgame-team.pages.debian.net/benchmarksgame/fastest/go-gpp.html)
- [Go Web Frame Benchmarks](https://github.com/smallnest/go-web-framework-benchmark)
- [Go HTTP Router Benchmark](https://github.com/julienschmidt/go-http-routing-benchmark)
- [Web 场景跨语言性能对比](https://www.techempower.com/benchmarks/)
- 《Systems Performance》
- [Dave 分享的 High Performance Go Workshop](https://dave.cheney.net/high-performance-go-workshop/dotgo-paris.html)
- [go-perfbook: best practices for writing high-performance Go code](https://github.com/dgryski/go-perfbook)
- [Delve](https://github.com/go-delve/delve/tree/master/Documentation)
- [What is Continuous Profiling?](https://www.opsian.com/blog/what-is-continuous-profiling)
- [Google-Wide Profiling: A Continuous Profiling Infrastructure for Data Centers](https://research.google/pubs/google-wide-profiling-a-continuous-profiling-infrastructure-for-data-centers/)
- [Go 语言笔试面试题汇总](https://geektutu.com/post/qa-golang.html)
- [七天用Go从零实现系列](https://geektutu.com/post/gee.html)
- [How to write benchmarks in Go](https://dave.cheney.net/2013/06/30/how-to-write-benchmarks-in-go)


---

> Author: [kyden](https://github.com/kydance)  
> URL: http://kydance.github.io/posts/golang-profile/  

