# Go 内存优化实战：逃逸分析完全指南


{{< admonition type=abstract title="导语" open=true >}}
在 Go 语言中，编译器通过逃逸分析来决定变量的内存分配位置，这直接影响着程序的性能和内存使用效率。但很多开发者对此知之甚少，导致代码中潜藏着性能隐患。本文将带你深入理解 Go 的逃逸分析机制，通过丰富的示例解析各种逃逸场景，帮助你掌握内存优化的关键技巧。从理论到实践，让你的 Go 程序更快、更高效。
{{< /admonition >}}

<!--more-->

## I. Golang Escape Ananlysis

Golang 编译器会自动决定把一个变量放在堆栈还是栈上，即逃逸分析（Escape Analysis）.

Go 声称逃逸分析可以释放程序员关于内存的使用限制，更多地关注程序逻辑本身。

---

## II. 逃逸规则

众所周知，当变量需要使用堆（heap）空间时，那么变量就应该进行逃逸。

一般情况下，一个引用对象中的引用类成员进行赋值，可能出现逃逸现象：可以理解为访问一个引用对象实际上底层就是通过一个指针来间接的访问，但如果再访问里面的引用成员就会有第二次间接访问，这样操作这部分对象的话，极大可能会出现逃逸的现象。

Golang 中的引用类型有函数类型 `func()`，接口类型 `interface`，切片类型 `slice`，字典类型 map，管道类型 `channel`，指针类型 `*` 等。

### 逃逸场景I： `[]interface{}`

**使用 `[]` 赋值 `[]interface{}` 数据类型，必定逃逸**

```Go
package main

import "fmt"

func main() {
	data := []interface{}{1, 2}
	val := data[0]
	fmt.Printf("%v\n", val)
	data[0] = 3
}
```

```Bash
$ go build -v -gcflags='-m' ./main.go
command-line-arguments
# command-line-arguments
./main.go:8:12: inlining call to fmt.Printf
./main.go:6:23: []interface {}{...} does not escape
./main.go:6:24: 1 escapes to heap
./main.go:6:27: 2 escapes to heap
./main.go:8:12: ... argument does not escape
./main.go:9:12: 3 escapes to heap
```

---

### 逃逸场景II： `map[string]interface{}`

**使用 `[]` 赋值 `map[string]interface{}` 数据类型，必定逃逸**

```Go
package main

import "fmt"

func main() {
	dat := make(map[string]interface{})
	dat["BlogName"] = "Kyden's Blog"
	val := dat["BlogName"]
	fmt.Printf("%v\n", val)
}
```

```Bash
$ go build -v -gcflags='-m' ./main.go
command-line-arguments
# command-line-arguments
./main.go:9:12: inlining call to fmt.Printf
./main.go:6:13: make(map[string]interface {}) does not escape
./main.go:7:20: "Kyden's Blog" escapes to heap
./main.go:9:12: ... argument does not escape
```

---

### 逃逸场景 III： `map[interface{}]interface{}`

**使用 `[]` 赋值 `map[interface{}]interface{}` 数据类型，必定逃逸**

```Go
package main

import (
	"fmt"
)

func main() {
	dat := make(map[interface{}]interface{})
	dat["BlogName"] = "Kyden's Blog"
	val := dat["BlogName"]
	fmt.Printf("%v\n", val)
}
```

```Bash
$ go build -v -gcflags='-m' ./main.go
command-line-arguments
# command-line-arguments
./main.go:11:12: inlining call to fmt.Printf
./main.go:8:13: make(map[interface {}]interface {}) does not escape
./main.go:9:6: "BlogName" escapes to heap
./main.go:9:20: "Kyden's Blog" escapes to heap
./main.go:10:13: "BlogName" does not escape
./main.go:11:12: ... argument does not escape
```

---

### 逃逸场景 IV：`map[string][]string`

`map[string][]string` 数据类型，赋值会发生 `[]string` 逃逸

```Go
package main

import (
	"fmt"
)

func main() {
	dat := make(map[string][]string)
	dat["BlogName"] = []string{"Kyden's Blog"}
	val := dat["BlogName"]
	fmt.Printf("%v\n", val)
}
```

```Bash
$ go build -v -gcflags='-m' ./main.go
command-line-arguments
# command-line-arguments
./main.go:11:12: inlining call to fmt.Printf
./main.go:8:13: make(map[string][]string) does not escape
./main.go:9:28: []string{...} escapes to heap
./main.go:11:12: ... argument does not escape
./main.go:11:21: val escapes to heap
```

---

### 逃逸场景 V：`[]*int`

`[]*int` 数据类型，赋值的右值会发生逃逸

```Go
package main

import "fmt"

func main() {
	dat := []*int{nil}
	a := 10
	dat[0] = &a
	fmt.Printf("%v\r\n", *dat[0])
	fmt.Printf("%v\r\n", dat[0])
}
```

```Bash
$ go build -v -gcflags='-m' ./main.go
command-line-arguments
# command-line-arguments
./main.go:9:12: inlining call to fmt.Printf
./main.go:10:12: inlining call to fmt.Printf
./main.go:7:2: moved to heap: a
./main.go:6:15: []*int{...} does not escape
./main.go:9:12: ... argument does not escape
./main.go:9:23: *dat[0] escapes to heap
./main.go:10:12: ... argument does not escape
```

---

### 逃逸场景 VI：`func(*int)`

`func(*int)` 数据类型，进行函数赋值，会使传递的形参逃逸

```Go
package main

import "fmt"

func f(a *int) {
	fmt.Printf("%v\n", *a)
	return
}

func main() {
	a := 10
	fn := f
	fn(&a)
	fmt.Printf("a = %v\n", a)
}
```

```Bash
$ go build -v -gcflags='-m' ./main.go
# command-line-arguments
./main.go:6:12: inlining call to fmt.Printf
./main.go:14:12: inlining call to fmt.Printf
./main.go:5:8: a does not escape
./main.go:6:12: ... argument does not escape
./main.go:6:21: *a escapes to heap
./main.go:14:12: ... argument does not escape
./main.go:14:25: a escapes to heap
```


---

> Author: [kyden](https://github.com/kydance)  
> URL: http://kydance.github.io/posts/golang-escape-analysis/  

