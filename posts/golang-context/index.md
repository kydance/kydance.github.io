# 浅析 Golang Context


{{&lt; admonition type=abstract title=&#34;导语&#34; open=true &gt;}}
当需要在多个 Goroutine 中传递上下文信息时，可以使用 Context 实现
{{&lt; /admonition &gt;}}

&lt;!--more--&gt;

## Context

Context 除了用来传递上下文信息，还可以用来传递终结执行子任务的相关信号，终止多个执行子任务的 Goroutine。

`context.Context` 接口数据结构：

```Go
// A Context carries a deadline, a cancellation signal, and other values across
// API boundaries.
//
// Context&#39;s methods may be called by multiple goroutines simultaneously.
type Context interface {
	// Deadline returns the time when work done on behalf of this context
	// should be canceled. Deadline returns ok==false when no deadline is
	// set. Successive calls to Deadline return the same results.
	Deadline() (deadline time.Time, ok bool)

	// Done returns a channel that&#39;s closed when work done on behalf of this
	// context should be canceled. Done may return nil if this context can
	// never be canceled. Successive calls to Done return the same value.
	// The close of the Done channel may happen asynchronously,
	// after the cancel function returns.
	//
	// WithCancel arranges for Done to be closed when cancel is called;
	// WithDeadline arranges for Done to be closed when the deadline
	// expires; WithTimeout arranges for Done to be closed when the timeout
	// elapses.
	//
	// Done is provided for use in select statements:
	//
	//  // Stream generates values with DoSomething and sends them to out
	//  // until DoSomething returns an error or ctx.Done is closed.
	//  func Stream(ctx context.Context, out chan&lt;- Value) error {
	//  	for {
	//  		v, err := DoSomething(ctx)
	//  		if err != nil {
	//  			return err
	//  		}
	//  		select {
	//  		case &lt;-ctx.Done():
	//  			return ctx.Err()
	//  		case out &lt;- v:
	//  		}
	//  	}
	//  }
	//
	// See https://blog.golang.org/pipelines for more examples of how to use
	// a Done channel for cancellation.
	Done() &lt;-chan struct{}

	// If Done is not yet closed, Err returns nil.
	// If Done is closed, Err returns a non-nil error explaining why:
	// Canceled if the context was canceled
	// or DeadlineExceeded if the context&#39;s deadline passed.
	// After Err returns a non-nil error, successive calls to Err return the same error.
	Err() error

	// Value returns the value associated with this context for key, or nil
	// if no value is associated with key. Successive calls to Value with
	// the same key returns the same result.
	//
	// Use context values only for request-scoped data that transits
	// processes and API boundaries, not for passing optional parameters to
	// functions.
	//
	// A key identifies a specific value in a Context. Functions that wish
	// to store values in Context typically allocate a key in a global
	// variable then use that key as the argument to context.WithValue and
	// Context.Value. A key can be any type that supports equality;
	// packages should define keys as an unexported type to avoid
	// collisions.
	//
	// Packages that define a Context key should provide type-safe accessors
	// for the values stored using that key:
	//
	// 	// Package user defines a User type that&#39;s stored in Contexts.
	// 	package user
	//
	// 	import &#34;context&#34;
	//
	// 	// User is the type of value stored in the Contexts.
	// 	type User struct {...}
	//
	// 	// key is an unexported type for keys defined in this package.
	// 	// This prevents collisions with keys defined in other packages.
	// 	type key int
	//
	// 	// userKey is the key for user.User values in Contexts. It is
	// 	// unexported; clients use user.NewContext and user.FromContext
	// 	// instead of using this key directly.
	// 	var userKey key
	//
	// 	// NewContext returns a new Context that carries value u.
	// 	func NewContext(ctx context.Context, u *User) context.Context {
	// 		return context.WithValue(ctx, userKey, u)
	// 	}
	//
	// 	// FromContext returns the User value stored in ctx, if any.
	// 	func FromContext(ctx context.Context) (*User, bool) {
	// 		u, ok := ctx.Value(userKey).(*User)
	// 		return u, ok
	// 	}
	Value(key any) any
}

```

- `Deadline`：返回 Context 被取消的时间，也就是完成工作的截至日期；
- `Done`：返回一个 channel，这个 channel 会在当前工作完成或者上下文被取消之后关闭，多次调用 `Done` 方法会返回同一个 channel；
- `Err`：放回 Context 结束的原因，只会在 `Done` 返回的 channel 被关闭时才会返回非空的值，如果 Context 被取消，会返回 Canceled 错误；如果 Context 超时，会返回 DeadlineExceeded 错误；
- `Value`：可用于从 Context 中获取传递的键值信息；

---

## Example

在 Web 请求的处理过程中，一个请求可能启动多个 goroutine 协同工作，这些 goroutine 之间可能需要共享请求的信息，且当请求被取消或者执行超时时，该请求对应的所有 goroutine 都需要快速结束，释放资源，Context 就是为了解决上述场景而开发的。

```Go
package main

import (
	&#34;context&#34;
	&#34;fmt&#34;
	&#34;time&#34;
)

const DB_ADDRESS = &#34;db_address&#34;
const CALCULATE_VALUE = &#34;calculate_value&#34;

func readDB(ctx context.Context, cost time.Duration) {
	fmt.Println(&#34;DB address is &#34;, ctx.Value(DB_ADDRESS))

	select {
	case &lt;-time.After(cost):
		fmt.Println(&#34;read data from db&#34;)
	case &lt;-ctx.Done():
		fmt.Println(ctx.Err())
	}
}

func calculate(ctx context.Context, cost time.Duration) {
	fmt.Println(&#34;calculate value is&#34;, ctx.Value(CALCULATE_VALUE))
	select {
	case &lt;-time.After(cost): //  模拟数据计算
		fmt.Println(&#34;calculate finish&#34;)
	case &lt;-ctx.Done():
		fmt.Println(ctx.Err()) // 任务取消的原因
		// 一些清理工作
	}
}

func main() {
	ctx := context.Background()

	// Add Context info
	ctx = context.WithValue(ctx, DB_ADDRESS, &#34;localhost:3306&#34;)
	ctx = context.WithValue(ctx, CALCULATE_VALUE, &#34;123&#34;)

	ctx, cancel := context.WithTimeout(ctx, time.Second*2)
	defer cancel()

	go readDB(ctx, time.Second*4)
	go calculate(ctx, time.Second*4)

	time.Sleep(time.Second * 5)
}

```

&gt; 使用 Context，能够有效地在一组 goroutine 中传递共享值、取消信号、deadline 等信息，及时关闭不需要的 goroutine。

---

## Reference

- [Go Context](https://github.com/golang/go/blob/release-branch.go1.22/src/context/context.go)


---

> : [kyden](https:github.com/kydance)  
> URL: http://kyden.us.kg/posts/golang-context/  

