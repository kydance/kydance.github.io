# Go 语言编程之道：编写优雅高效的 Golang 代码


{{&lt; admonition type=abstract title=&#34;导语&#34; open=true &gt;}}
想要写出优雅且高效的 Go 代码，仅仅了解语法是远远不够的。本文将带你深入探索 Go 语言的设计哲学和最佳实践，从代码格式化、命名规范到控制结构的巧妙运用，帮助你掌握编写地道 Go 代码的精髓。无论你是 Go 新手还是有经验的开发者，都能从中获得实用的编程技巧和深刻的设计思想。
{{&lt; /admonition &gt;}}

&lt;!--more--&gt;

## Formatting 格式化

在 Golang 中，`gofmt` 以包未处理对象而非源文件，它将 Go 程序按照标准风格缩进、对齐，保留注释并在需要时重新格式化。

- Indentation 缩进：使用 `制表符 Tab` 缩进，`gofmt` 默认使用
- Line length 行长度：Go 对行的长度没有限制
- Parentheses 括号：控制结构（`if`, `for`, `switch`）在语法上并不需要圆括号

## Commentary 注释

Go 支持 C 风格的块注释 `/* */` 和 C&#43;&#43; 风格的单行注释 `//`，其中，`//` 注释更常用，而 `/* */` 则主要用于包的注释

`godoc` 即使一个程序，又是一个 Web 服务器，它对 Go 的源码进行处理，并提取包中的文档内容：
出现在顶级声明之前，且与该声明之间没有空行的注释，将与该声明一起被提出来，作为该条目的说明文档。

每个包都应包含一段包注释，即放置在包子句前的一个块注释。
对于包含多个文件的包，包注释只需出现在其中的任一文件中即可。
包注释应在整体上对该包进行介绍，并提供包的相关信息。
它将出现在 `godoc` 页面中的最上面，并为紧随其后的内容建立详细的文档。

```Go
/*
Package regex implements a simple library for regular expressions.

The syntax of the regular expressions accepted is:

  regexp:
    concatenation { &#39;|&#39; concatenation }
  concatenation:
    { closure }
  closure:
    term [ &#39;*&#39; | &#39;&#43;&#39; | &#39;?&#39; ]
  term:
    &#39;^&#39;
    &#39;$&#39;
    &#39;.&#39;
    character
    &#39;[&#39; [ &#39;^&#39; ] character-range &#39;]&#39;
    &#39;(&#39; regexp &#39;)&#39;
*/
package regex
```

## Names 命名

### Package names 包名

**当一个包被导入后，包名就会成为内容的访问器 `import &#34;bytes&#34;`，按照惯例，包应当以某个小写的单个单词命名，且不应使用下划线或驼峰记法**。
例如，`err` 的命名就是出于简短考虑。

包名是导入时所需的唯一默认名称，它并不需要在所有源码中保持唯一，即便在少数发生冲突的情况下，也可为导入的包选择一个别名来局部使用。

另一个约定：**包名应为其源码目录的基本名称**。例如，`src/pkg/encoding/base64` 中的包应作为 `encoding/basee64` 导入，其包名应为 `base64` 而不是 `encoding_base64` / `encodingBase64`.

*长命名并不会使包更具有可读性，反而一份有用的说明文档通常比额外的长名更具价值。*

### Getter / Setter

Go 并不对 getter 和 setter 提供自动支持。

如将 `Get` 放入 getter 的名字中，既不符合习惯，也没有必要，但**大写字母**作为字段导出的 getter 是一个不错的选择，另外 `Set` 放入 setter 是个不错的选择。

```Golang
type Object struct {
  ower string
}

func (o *Object) Ower() string { return o.ower }
func (o *Object) SetOwer(s string) { o.ower = s }
```

### Interface names 接口名

按照规定，只包含一个方法的接口应当以该方法的名称加上 `er` 后缀来命名，如 `Reader` / `Writer` / `Formater` 等。

字符串转换方法命名应为 `String` 而非 `ToString`

### MixedCaps 驼峰记法

Go 中约定使用驼峰记法

## 分号

和 C 一样，Go 的正式语法使用分号 `;` 来结束语句，但 Go 的分号不一定出现在源码中，而是词法分析器会使用一条简单的规则来自动插入分号

规则：**如在新行前的最后一个标记为标识符（`int`/`float64`等）、数值或字符串常量之类的基本字面或`break`、`continue`、`fallthrough`、`return`、`&#43;&#43;`、`--`、`)`、`}` 之一，则词法分析器将始终在该标记后面插入分号**，即**如果新行前的标记为语句的末尾，则插入分号`;`**。

通常，Go 程序只在诸如 `for` 循环子句这样的地方使用分号，来以此将初始化器、条件及增量元素分开；

## Control structures 控制结构

Go 不再使用 `do` / `while` 循环，只有一个更为通用的 `for`，

```go
// C: for
for init; condition; post { }
// C: while
for condition { }
// C: for(;;)
for { }

// [12]aT, []vT, map[sting]any mT
for key, value := range aT/vT/mT { }
for key := range aT/vT/mT { }
for _, value := range aT/vT/mT { }
```

&gt; Go 没有逗号操作符，且 `&#43;&#43;`/`--` 是语句而非表达式

```go
for i, j := 0, len(aT) - 1; i &lt; j; i, j = i &#43; 1, j - 1 { // Not: i&#43;&#43;, j--
  a[i], a[j] = a[j], a[i]
}
```

`switch` 更加灵活，其表达式无需为常量或整数，`case` 语句会自上而下逐一进行求值直至匹配为止，它不会自动下溯，但 `case` 可通过逗号分隔来列举相同的处理条件

`break` 语句会使 `switch` 提前终止

```go
func unhex(c byte) byte {
  switch {
  case &#39;0&#39; &lt;= c &amp;&amp; c &lt;= &#39;9&#39;:
    return c - &#39;0&#39;
  case &#39;a&#39; &lt;= c &amp;&amp; c &lt;= &#39;f&#39;:
    return c - &#39;a&#39; &#43; 10
  case &#39;A&#39; &lt;= c &amp;&amp; c &lt;= &#39;F&#39;:
    return c - &#39;A&#39; &#43; 10
  }
  return 0
}

func shouldEscape(c byte) bool {
  switch c {
    case &#39; &#39;, &#39;?&#39;, &#39;&amp;&#39;, &#39;=&#39;, &#39;#&#39;, &#39;&#43;&#39;, &#39;%&#39;:
      return true
  }
  return false
}
```

`if` 强制使用大括号，并且接受初始化语句

```go
if err := file.Chmod(0664); err != nil {
  return err
}
```

## Function 函数

Go 与众不同的特性之一，就是函数和方法可以返回多个值，返回值或结果“形参”可被命名，并作常规变量使用。

Go 的 `defer` 语句用于预设一个函数调用（即延迟执行函数），该函数会在执行 `defer` 的函数返回之前立即执行。
被推迟的多个函数，会按照后进先出（LIFO）的顺序执行。

```go
func Contents(filename string) (string, error) {
  f, err := os.Open(filename)
  if err != nil {
    return &#34;&#34;, err
  }
  defer f.Close()

  var result []byte
  buf := make([]byte, 100)
  for {
    n, err := f.Read(buf[0:])
    result = append(result, buf[0:n]...)
    if err != nil {
      if err == io.EOF {
        break
      }
      return &#34;&#34;, err
    }
  }
  return string(result), nil
}
```

## Data 数据

`new(T)` 会为类型为 `T` 的新项分配已置零的内存空间， 并返回它的地址，也就是一个类型为 `*T` 的值(返回一个指针， 该指针指向新分配的，类型为 `T` 的零值)。

内建函数 `make(T, args)` 的目的不同于 `new(T)`。它只用于创建切片、映射和信道，并返回类型为 `T`（而非 `*T` ）的一个已初始化 （而非置零）的值。 出现这种用差异的原因在于，这三种类型本质上为引用数据类型，它们在使用前必须初始化。

```go
// Allocates slice structure; *p == nil; rarely useful
var p *[]int = new([]int)
// The slice v now refers to a new array of 100 ints
var v []int = make([]int, 100)

// 惯用法
v := make([]int, 100)
```

### Array 数组

数组主要用作切片的构件，主要特点：

- 数组是值，讲一个数组赋值给另一个数组会复制其所有元素
- 如将数组作为参数传入某个函数，则会收到该数组的一份副本而非指针
- 数组的大小是其类型的一部分

```go
func Sum(a *[3]float64) (sum float64) {
  for _, v := range *a {
    sum &#43;= v
  }
  return
}

aV := [...]float64{1, 2, 0.7}
fmt.Println(Sum(&amp;aV))
```

{{&lt; admonition type=note title=&#34;Go array&#34; open=true &gt;}}
在 C 语言中，数组变量是指向第一个元素的指针，但 Go 语言中并不是。

Go 语言中，数组变量属于值类型（value type），因此当一个数组变量呗赋值或传递时，实际上会复制整个数组 -&gt;

为了避免复制数组，一般传递指向数组的指针: `func f(pa *[3]uint8) { ... }`

{{&lt; /admonition &gt;}}

{{&lt; admonition type=tip title=&#34;Go array&#34; open=true &gt;}}
Go 中的数组类型定义了长度和元素类型。
例如，`[2]int` 类型表示由 2 个 int 整型组成的数组。

数组以索引方式访问（a[i] 访问数组 a 的第 i 个元素）。
数组的长度固定，且是数组类型的一部分 -&gt; 长度不同的 2 个数组不可以相互赋值，因为它们属于不同的类型。
{{&lt; /admonition &gt;}}

### Slice 切片

切片通过对数组进行封装，为数据序列提供了更通用、强大而方便的接口，其本质是是在数组 Array 之上的抽象数据类型。

```Go
struct {
    ptr *[]T
    len int
    cap int
}

```

slice 保存了对底层数组的引用，如将某个 slice 赋值给另一个 slice，则他们会引用同一个数组。

**若某个函数将一个切片作为参数传入，则它对该切片元素的修改对调用者而言同样可见，
这可以理解为传递了底层数组的指针**

只要切片不超出底层数组的限制，它的长度就是可变的，只需将它赋予其自身的切片即可。
切片的容量可通过内建函数 `cap` 获得，它将给出该切片可取得的最大长度。

尽管 Append 可修改 slice 的元素，但切片自身（其运行时数据结构包含指针、长度和容量）是通过值传递的.

{{&lt; admonition type=tip title=&#34;Go slice&#34; open=true &gt;}}
Go 中 slice 容量指的是当前切片以及预分配的内存能够容纳的元素个数.

若数据超出其容量，则会重新分配该切片，返回值即为所得的切片。

为了减少内存分配、拷贝的次数，在容量较小时，一般是以 2 的倍数进行扩大（2 -&gt; 4 -&gt; 8 -&gt; 16），
当达到 2048 时，为避免申请的内存过大，从而浪费空间 =&gt; [Go 语言 1.20 实现如下](https://github.com/golang/go/blob/release-branch.go1.20/src/runtime/slice.go#L157)：

```Go
// growslice allocates new backing store for a slice.
//
// arguments:
//
// oldPtr = pointer to the slice&#39;s backing array
// newLen = new length (= oldLen &#43; num)
// oldCap = original slice&#39;s capacity.
//    num = number of elements being added
//     et = element type
//
// return values:
//
// newPtr = pointer to the new backing store
// newLen = same value as the argument
// newCap = capacity of the new backing store
//
// Requires that uint(newLen) &gt; uint(oldCap).
// Assumes the original slice length is newLen - num
//
// A new backing store is allocated with space for at least newLen elements.
// Existing entries [0, oldLen) are copied over to the new backing store.
// Added entries [oldLen, newLen) are not initialized by growslice
// (although for pointer-containing element types, they are zeroed). They
// must be initialized by the caller.
// Trailing entries [newLen, newCap) are zeroed.
//
// growslice&#39;s odd calling convention makes the generated code that calls
// this function simpler. In particular, it accepts and returns the
// new length so that the old length is not live (does not need to be
// spilled/restored) and the new length is returned (also does not need
// to be spilled/restored).
func growslice(oldPtr unsafe.Pointer, newLen, oldCap, num int, et *_type) slice {
 oldLen := newLen - num

 // ...

 newcap := oldCap
 doublecap := newcap &#43; newcap
 if newLen &gt; doublecap {
  newcap = newLen
 } else {
  const threshold = 256
  if oldCap &lt; threshold {
   newcap = doublecap
  } else {
   // Check 0 &lt; newcap to detect overflow
   // and prevent an infinite loop.
   for 0 &lt; newcap &amp;&amp; newcap &lt; newLen {
    // Transition from growing 2x for small slices
    // to growing 1.25x for large slices. This formula
    // gives a smooth-ish transition between the two.
    newcap &#43;= (newcap &#43; 3*threshold) / 4
   }
   // Set newcap to the requested cap when
   // the newcap calculation overflowed.
   if newcap &lt;= 0 {
    newcap = newLen
   }
  }
 }

 // ...

 return slice{p, newLen, newcap}
}
```

{{&lt; /admonition &gt;}}

#### 二维数组

一种是独立地分配每一个切片；而另一种就是只分配一个数组， 将各个切片都指向它

```go
// 独立地分配每一个切片
pic := make([][]uint8, YSize)
for i := range pic {
  // 一次一行
  pic[i] = make([]uint8, XSize)
}

// 顶层 slice
pic := make([][]uint8, YSize)
// 分配一个大的切片来保存所有像素
pixels := make([]uint8, XSize*YSize)
// 遍历行，从剩余像素切片的前面切出每一行
for i := range pic {
  pic[i], pixels = picxels[:XSize], pixels[XSize:]
}
```

#### 切片操作及性能

Go 语言在 Github 上的官方
wiki - [SliceTricks](https://github.com/golang/go/wiki/SliceTricks) 介绍了切片常见的操作技巧。
另一个项目
[Go Slice Tricks Cheat Sheet](https://ueokande.github.io/go-slice-tricks/) 将这些操作以图片的形式呈现了出来，非常直观。

{{&lt; figure src=&#34;/posts/effective-go/copy.png&#34; title=&#34;&#34; &gt;}}

---

{{&lt; figure src=&#34;/posts/effective-go/append.png&#34; title=&#34;&#34; &gt;}}

slice 有 3 个属性，指针（ptr）、长度（len）和容量（cap），因此当 append 时存在两种场景：

- append 后的长度小于等于 cap，将会直接使用原底层数组剩余的空间
- append 后的长度大于 cap，将会分配一块更大的区域来容纳新的底层数组

&gt; 为了避免内存发生拷贝，若能够知道最终的切片的大小，预先设置 cap 的值能够获得最好的性能

---

{{&lt; figure src=&#34;/posts/effective-go/delete.png&#34; title=&#34;&#34; &gt;}}

slice 的底层是数组，所以 delete 意味着后面的元素需要逐个向前移位
=&gt; delete 的复杂度为O(N)
=&gt; slice 不适合大量随机删除的场景（链表 list 更适合）

---

{{&lt; figure src=&#34;/posts/effective-go/delete_gc.png&#34; title=&#34;&#34; &gt;}}

删除后，将空余位置置空，有助于垃圾回收。

---

{{&lt; figure src=&#34;/posts/effective-go/insert.png&#34; title=&#34;&#34; &gt;}}

insert 和 append 类似，即在某个位置添加一个元素后，将该位置后面的元素再 append 回去，复杂度为 O(N) =&gt; 不适合大量随机插入的场景。

---

{{&lt; figure src=&#34;/posts/effective-go/filter_in_place.png&#34; title=&#34;&#34; &gt;}}

当原切片不会再被使用时，就地 filter 方式是比较推荐的，可以节省内存空间。

---

{{&lt; figure src=&#34;/posts/effective-go/push.png&#34; title=&#34;&#34; &gt;}}

在末尾追加元素，不考虑内存拷贝的情况，复杂度为 O(1)。

---

{{&lt; figure src=&#34;/posts/effective-go/push_front.png&#34; title=&#34;&#34; &gt;}}

在头部追加元素，时间和空间复杂度均为 O(N)，不推荐。

---

{{&lt; figure src=&#34;/posts/effective-go/pop.png&#34; title=&#34;&#34; &gt;}}

尾部删除元素，复杂度 O(1)

---

{{&lt; figure src=&#34;/posts/effective-go/pop_front.png&#34; title=&#34;&#34; &gt;}}

头部删除元素，如果使用切片方式，复杂度为 O(1)。

&gt; 需要注意的是，底层数组没有发生改变，第 0 个位置的内存仍旧没有释放。
&gt; 如果有大量这样的操作，头部的内存会一直被占用。

### Map

可以关联不同类型的值。其键可以是任何相等性操作符支持的类型， 如整数、浮点数、复数、字符串、指针、接口（只要其动态类型支持相等性判断）、结构以及数组。 切片不能用作映射键，因为它们的相等性还未定义。与切片一样，

映射也是引用类型。 若将映射传入函数中，并更改了该映射的内容，则此修改对调用者同样可见。
**若试图通过映射中不存在的键来取值，就会返回与该映射中项的类型对应的零值**

要删除 map 中的某项，可使用内建函数 `delete`，它以映射及要被删除的键为实参。
即便对应的键不在该 map 中，此操作也是安全的。

## Reference

- [Effective Go](https://go.dev/doc/effective_go)


---

> Author: [kyden](https://github.com/kydance)  
> URL: http://kydance.github.io/posts/effective-go/  

