# 设计模式精讲：从理论到实战的最佳实践指南


{{&lt; admonition type=abstract title=&#34;导语&#34; open=true &gt;}}
设计模式是每个程序员的必修课，但如何正确理解和灵活运用却是一门艺术。本文将带你深入浅出地探索 Golang 项目中常用的 8 种经典设计模式，通过 Go 和 C&#43;&#43; 的实际代码示例，让你真正理解每种模式的精髓。从面向对象设计原则到具体实现，从模式分类到实战应用，助你构建更优雅、更可维护的代码架构。
{{&lt; /admonition &gt;}}

&lt;!--more--&gt;

## I. 前言

软件设计模式（Design Pattern），是一套被反复使用、多数人知晓的、经过分类编目的、代码设计经验的总结，使用设计模式是为了可重用代码、让代码更容易被他人理解并且保证代码可靠性。
简单来说，设计模式就是**在一定环境下，用固定套路解决问题**。

&gt; 设计模式的基础：**多态**

### 面向对象设计原则

    &gt; 目的：高内聚、低耦合

如何同时提⾼⼀个软件系统的**可维护性**和**可复⽤性**是⾯向对象设计需要解决的核⼼问题之⼀。
⾯向对象设计原则为⽀持可维护性复⽤⽽诞⽣，这些原则蕴含在很多设计模式中，它们是从许多设计⽅案中总结出的指导性原则。

- 单一职责原则: 类的职责单⼀，对外只提供⼀种功能，⽽引起类变化的原因都应该只有⼀个
- 开闭原则: **类的改动是通过增加代码进⾏的，⽽不是修改源代码**
- 里式代换原则: 任何抽象类出现的地⽅都可以⽤他的实现类进⾏替换，实际就是虚拟机制，语⾔级别实现⾯向对象功能
- 依赖倒转原则: **依赖于抽象(接⼝)，不要依赖具体的实现(类)，也就是针对接⼝编程**
- 接口隔离原则: 不应该强迫⽤户的程序依赖他们不需要的接⼝⽅法。⼀个接⼝应该只提供⼀种对外功能，不应该把所有操作都封装到⼀个接⼝中去
- 合成复用原则: 如果使⽤继承，会导致⽗类的任何变换都可能影响到⼦类的⾏为。如果使⽤对象组合，就降低了这种依赖关系。**对于继承和组合，优先使⽤组合**
- 迪米特法则: **⼀个对象应当对其他对象尽可能少的了解，从⽽降低各个对象之间的耦合，提⾼系统的可维护性**

## II. 分类

- 创建型（Creational）模式：如何创建对象

| 模式名称 | 用途 |
| :--- | :--- |
| **单例模式** &lt;br&gt; 🌟🌟🌟🌟 | 保证一个类仅有一个实例，并提供一个访问它的全局访问点 |
| **简单工厂方法** &lt;br&gt; 🌟🌟🌟 | 通过专门定义一个类来负责创建其他类的实例，被创建的实例通常都具有共同的基类 |
| **抽象工厂方法** &lt;br&gt; 🌟🌟🌟🌟🌟 | 提供一个创建一系列相关或相互依赖的接口，而无需指定它们具体的类 |
| 原型模式 | ⽤原型实例指定创建对象的种类，并且通过拷⻉这些原型创建新的对象 |
| 建造者模式 | 将⼀个复杂的构建与其表示相分离，使得同样的构建过程可以创建不同的表示 |

- 结构型（Structural）模式：如何实现类或对象的组合

| 模式名称 | 用途 |
| :--- | :--- |
| **适配器模式** &lt;br&gt; 🌟🌟🌟🌟 | 将一个类的接口转换成客户希望的另外一个接口，使得原本由于接口不兼容而不能一起工作的那些类可以一起工作 |
| 桥接模式 | 将抽象部分与实际部分分离，使它们可以独立的变化 |
| **组合模式** &lt;br&gt; 🌟🌟🌟🌟 | 将对象组合成树形结构以表示 “部分 - 整体” 的层次结构，使得用户对单个对象和组合对象的使用具有一致性 |
| **装饰模式** &lt;br&gt; 🌟🌟🌟 | 动态地给一个对象添加一些额外的职责：就增加功能来说，此模式比生成子类更加灵活 |
| **外观模式** &lt;br&gt; 🌟🌟🌟🌟🌟 | 为子系统的一组接口提供一个一致的界面，此模式定义了一个高层次接口，使得这一子系统更容易使用 |
| 享元模式 | 以共享的方式高效的支持大量的细粒度的对象 |
| 代理模式 | 为其他对象提供一种代理以控制这个对象的访问 |

- 行为型（Behavioral）模式：类或对象如何交互以及如何分配指责

## III. 创建型设计模式

### 1. 单例模式（Singleton Pattern）

意图：它是一种创建型设计模式，限制了实例化类的对象个数，确保一个类只有一个实例，并且提供一个全局访问点。

{{&lt; admonition type=warning title=&#34;warning&#34; open=true &gt;}}

Singleton Pattern 同时解决了两个问题，因此违法了**单一职责原则**:

1. 保证一个类只用一个实例。
2. 为该实例提供一个全局访问节点。

{{&lt; /admonition &gt;}}

#### 应用场景

- 配置管理器：在应用程序中，配置信息通常需要一个实例来管理，如此可以保证配置信息的一致性
- 连接池：数据库连接池需要限制数据库连接的数量，以避免过多的连接消耗资源
- 日志记录器：日志系统通常只需要一个实例来记录应用程序的日志信息，以避免日志信息的冗余和混乱
- 硬件管理器：对于某些硬件设备，如打印机 / 扫描仪等，可能只需要一个管理器来控制对它们的访问
- 应用状态管理：在某些应用中，需要全局的管理状态，如用户会话管理或权限验证状态

#### 解决方案

- 将默认构造函数设为私有，防止其他对象使用单例类的 `new` 运算符
- 新建一个静态构建方法作为构造函数：该函数会“偷偷”调用私有构造函数来创建对象，并将其保存到一个静态成员变量中，之后所有对于该函数的调用都将返回这一缓存对象。

#### 单例模式结构

{{&lt; figure src=&#34;/posts/design-pattern/FactoryMethod-Singleton.svg&#34; title=&#34;&#34; &gt;}}

#### 与其他模式的关系

- **外观模式**类通常可以转化为**单例模式**类，因为在大部分情况下一个外观对象就足够啦
- 如果能将对象的所有共享状态简化为一个享元对象，那么**享元模式**就和**单例**类似，但二者有两个根本性的不同：
	1. 单例只有一个单例实体，但享元类可以有多个实体，各实体的内在状态也可以不同
	2. 单例对象可以是可变的，享元对象不可变
- **抽象工厂模式**、**生成器模式**和**原型模式**都可以用**单例**来实现

#### 应用示例

    ```Go
    // Singleton.go
    package singleton

    import &#34;sync&#34;

    var instance *Singleton
    var once sync.Once

    type Singleton struct {
        str string
    }

    func GetInstance() *Singleton {
        if instance != nil {
            return instance
        }

        once.Do(func() {
            instance = &amp;Singleton{}
        })
        return instance
    }
    ```

### 2. 工厂模式（Factory Pattern）

亦称：虚拟构造函数、Virtual Constructor、Factory Method

意图：它是一种创建型设计模式，**其在父类中提供一个创建对象的方法，允许子类决定实例化对象的类型**

由于 Golang 中缺少类和继承等 OOP 特性，因此，无法使用 Go 来实现经典的工厂方法模式，但我们仍能实现基础版本，即简单工厂。

    ```Go
    // iGun.go
    package factory

    type Gun interface {
        setName(name string)
        setPower(power int)
        name() string
        power() int
    }

    // gun.go
    type gun struct {
        name string
        power int
    }

    func (g *gun) setName(name string) { g.name = name }
    func (g *gun) setPower(power int) { g.power = power }
    func (g *gun) name() string { return g.name }
    func (g *gun) power() int { return g.power }

    // ak47.go
    type ak47 struct {
        gun
    }

    func newAk47() Gun {
        return &amp;ak47{
            gun: gun{
                name: &#34;AK47&#34;,
                power: 10,
            }
        }
    }

    // m16.go
    type m16 struct {
        gun
    }

    func m16() Gun {
        return &amp;gun{
            name: &#34;M16&#34;,
            power: 17,
        }
    }

    // Factory.go
    func GunFactory(gunType string) (Gun, error) {
        switch gunType {
        case &#34;ak47&#34;:
            return newAk47(), nil
        case &#34;m16&#34;:
            return newM16(), nil
        default:
            return nil, errors.New(&#34;wrong gun type&#34;)
        }
    }
    ```

## IV. 行为设计模式

### 1. 策略模式（Strategy Pattern）

**策略模式**是一种行为设计模式，它能让你定义一系列算法，并将每种算法分别放入独立的类中，以使算法的对象能够相互替换。

在项目开发中，我们经常要根据不同的场景，采取不同的措施，也就是不同的策略。通过 `if ... else ...` 的形式来调用不同的策略，这种方式称之为**硬编码**。

#### 内存缓存示例

假设构建内存缓存的场景，由于数据存于内存中，其大小会受到限制。
在达到其大小上限后，一些数据就必须被移除以留出空间，而此类操作可通过多种算法实现，例如：

- 最少最近使用（LRU）算法：移除最近最少使用的数据
- 最近最少使用（LFU）算法：移除使用频率最少使用的数据
- 先进先出（FIFO）算法：移除最先进入的数据

问题在于如何将缓存类与这些算法解耦，以便在运行时更改算法。
另外，在添加新算法时，缓存类不应该改变。

这就是策略模式发挥作用的场景：创建一系列算法，每个算法都有自己的类，这些类中的每一个都遵循相同的接口，这使得这些算法可以相互替换。

```Go
// cache.go
type Cache struct {
    storage     map[string]any
    rmAlgo      RmAlgo
    capacity    int
    maxCapacity int
}

func initCache(algo RmAlgo) *Cache {
    return &amp;Cache{
        storage:     make(map[string]any),
        rmAlgo:      algo,
        capacity:    0,
        maxCapacity: 100,
    }
}

func (c *Cache) rm() {
    c.rmAlgo.Rm(c)
    c.capacity--
}
func (c *Cache) setRmAlgo(algo RmAlgo) { c.rmAlgo = algo }
func (c *Cache) get(key string) any { return c.storage[key] }

func (c *Cache) add(key string, value any) {
    if c.capacity &gt;= c.maxCapacity {
        c.rm()
    }
    c.storage[key] = value
    c.capacity&#43;&#43;
}

// iCache.go 策略接口
type RmAlgo interface {
    Rm(c *Cache)
}

// fifo.go
type Fifo struct{}
func (f *Fifo) Rm(c *Cache) { fmt.Println(&#34;rm by fifo strategy&#34;) }

// lru.go
type Lru struct{}
func (l *Lru) Rm(c *Cache) { fmt.Println(&#34;rm by lru strategy&#34;) }

// lfu.go
type Lfu struct{}
func (l *Lfu) Rm(c *Cache) { fmt.Println(&#34;rm by lfu strategy&#34;) }
```

### 2. 模板方法模式（Template Method Pattern）

**模板方法模式**是一种行为设计模式，它定义了一个操作中的算法的骨架，允许子类在不修改结构的情况下重写算法的特定步骤。

#### OTP 示例

假设在处理一个一次性密码（OTP）的场景，将 OTP 传递给用户的方式多种多样（短信、邮件等），但无论是短信还是邮件，整个 OTP 处理过程都是相同的：

1. 生成一个随机的 n 位 OTP 数字
2. 在缓存中保存这组数字以便进行后续验证
3. 准备内容
4. 发送通知

后续引入的任何新 OTP 类型都很有可能需要进行相同的步骤。

首先，定一个由固定数量的方法组成的基础模板算法，然后将实现每一个步骤方法，但不改变模版方法。

```Go
// iOtp.go
type IOtp interface {
    GenerateRandomOtp(length int) string
    CacheOtp(otp string)
    PrepareContent() string
    SendNotification(message string) error
}

type Otp struct {
    iOtp IOtp
}
func (o *Otp) GenAndSendOtp(length int) error {
    opt := o.iOtp.GenerateRandomOtp(length)
    o.iOtp.CacheOtp(opt)
    content := o.iOtp.PrepareContent()
    return o.iOtp.SendNotification(content)
}

// sms.go
type Sms struct {
    Otp
}
func (s *Sms) GenerateRandomOtp(length int) string {
    opt := &#34;&#34;
    for i := range length {
        opt &#43;= strconv.Itoa(rand.Intn(10))
    }
    fmt.Println(&#34;SMS: Generate otp %s&#34;, opt)
    return opt
}
func (s *Sms) CacheOtp(otp string) { fmt.Println(&#34;SMS: Cache otp %s&#34;, otp) }
func (s *Sms) PrepareContent() string { return fmt.Sprintf(&#34;Your OTP is %s&#34;, otp) }
func (s *Sms) SendNotification(message string) error {
    fmt.Println(&#34;SMS: Send message %s&#34;, message)
    return nil
}

// email.go
type Email struct {
    Otp
}
func (e *Email) GenerateRandomOtp(length int) string {
    opt := &#34;&#34;
    for i := range length {
        opt &#43;= strconv.Itoa(rand.Intn(10))
    }
    fmt.Println(&#34;Email: Generate otp %s&#34;, opt)
    return opt
}
func (e *Email) CacheOtp(otp string) { fmt.Println(&#34;Email: Cache otp %s&#34;, otp) }
func (e *Email) PrepareContent() string { return fmt.Sprintf(&#34;Your OTP is %s&#34;, otp) }
func (e *Email) SendNotification(message string) error {
    fmt.Println(&#34;Email: Send message %s&#34;, message)
    return nil
}
```

## V. 结构型设计模式

### 1. 代理模式（Proxy Pattern）

**代理模式**是一种结构设计模式，让你能够提供对象的替代品或其占位符。
代理控制着对于原对象的访问，并允许在将请求提交给对象前后进行一些处理(访问控制、缓存等)。

代理模式建议新建一个与原服务对象接口相同的代理类，然后更新应用以将代理对象传递给所有原始对象客户端。
代理类接收到客户端请求后会创建实际的服务对象，并将所有工作委派给它。

#### Nginx 代理示例

Nginx 这样的 web 服务器可充当应用程序服务器的代理：

- 提供了的应用程序服务器的受控访问权限
- 可限制速度
- 可缓存请求

```Go
// server.go
type Server interface {
    HandleRequest(string, string) (int, string)
}

// nginx.go
type Nginx struct {
    application       *Application
    maxAllowedRequest int
    rateLimiter       map[string]int
}

func NewNginx() *Nginx {
    return &amp;Nginx{
        application:       &amp;Application{},
        maxAllowedRequest: 10,
        rateLimiter:       make(map[string]int),
    }
}
func (n *Nginx) HandleRequest(url string, method string) (int, string) {
    allowed := n.checkRateLimit(url)
    if !allowed {
        return 403, &#34;Forbidden&#34;
    }
    return n.application.HandleRequest(url, method)
}

func (n *Nginx) checkRateLimit(url string) bool {
    if n.rateLimiter[url] == 0 {
        n.rateLimiter[url] = 1
    }
    if n.rateLimiter[url] &gt; n.maxAllowedRequest {
        return false
    }
    n.rateLimiter[url]&#43;&#43;
    return true
}

// application.go
type Application struct {}
func (a *Application) HandleRequest(url string, method string) (int, string) {
    if url == &#34;/app/status&#34; &amp;&amp; method == &#34;GET&#34; {
        return 200, &#34;OK&#34;
    }

    if url == &#34;/create/user&#34; &amp;&amp; method == &#34;POST&#34; {
        return 201, &#34;User Created&#34;
    }

    return 404, &#34;Not Found&#34;
}
```

#### 2. 选项模式

**选项模式**是一种结构设计模式，可以创建一个带有默认值的 struct 变量，并选择性地修改其中一些参数的值。

在 Python 中，创建一个对象时，可以给参数设置默认值，这样在不传入任何参数时，
可以返回携带默认值的对象，并在需要时修改对象的属性。
这种特性可以大大简化开发者创建一个对象的成本，尤其是在对象拥有众多属性时。

然而，在 Go 生态中，因为不支持给参数设置默认值，为了既能够创建带默认值的实例，又能够自定义参数的实例，开发者一般会通过以下两种方法实现：

1. 分别开发两个用来创建实例的函数，一个带有默认值，一个不带默认值：此时需要实现两个函数，实现方式很不优雅；

    ```Go
    package options

    const (
        defaultTimeout = 10
        defaultCaching = false
    )

    type Connection struct {
        addr string
        cache int
        timeout time.Duration
    }

    func NewConnection(addr string) (*Connection, error) {
        return &amp;Connection{
            addr:   addr,
            cache:  defaultCaching,
            timeout: defaultTimeout,
        }, nil
    }

    func NewConnectionWithOptions(addr string, cache bool, timeout time.Duration) (*Connection, error) {
        return &amp;Connection{
            addr:   addr,
            cache:  cache,
            timeout: timeout,
        }, nil
    }
    ```

2. 创建一个带有默认值的选项，并用该选项创建实例: 每次创建实例时，都需要创建 `Options`，操作起来比较麻烦；

```Go
package options

const (
    defaultTimeout = 10
    defaultCaching = false
)

type Connection struct {
    addr string
    cache int
    timeout time.Duration
}

type ConnectionOption struct {
    Cache   bool
    Timeout time.Duration
}

func NewDefaultConnectionOption() *ConnectionOption {
    return &amp;ConnectionOption{
        Cache:   defaultCaching,
        Timeout: defaultTimeout,
    }
}

func NewConnection(addr string, opts *ConnectionOption) (*Connection, error) {
    return &amp;Connection{
        addr:   addr,
        cache:  opt.Cache,
        timeout: opt.Timeout,
    }, nil
}
```

##### 解决方案

```Go
package options

import &#34;time&#34;

type Connection struct {
    addr string
    cache int
    timeout time.Duration
}

const (
    defaultTimeout = 10
    defaultCaching = false
)

type options struct {
    timeout time.Duration
    cache   bool
}

// Option overrides behavior of Connection
type Option interface {
    apply(*options)
}

type optionFunc func(*options)
func (f optionFunc) apply(o *options) { f(o) }

func WithTimeout(t time.Duration) Option {
    return optionFunc(func(o *options) {
        o.timeout = t
    })
}

func WithCaching(c bool) Option {
    return optionFunc(func(o *options) {
        o.cache = c
    })
}

func NewConnection(addr string, opts ...Option) (*Connection, error) {
    o := &amp;options{
        timeout: defaultTimeout,
        cache:   defaultCaching,
    }
    for _, opt := range opts {
        opt.apply(o)
    }
    return &amp;Connection{
        addr:   addr,
        cache:  o.cache,
        timeout: o.timeout,
    }, nil
}
```

`Option` 类型的选项参数需要实现 `apply(*options)` 函数，结合 `WithTimeout`、`WithCache` 函数的返回值和 `optionFunc` 的 `apply` 方法实现，可以知道 `o.apply(&amp;options)` 其实就是把 `WithTimeout`、`WithCache` 的返回值赋值给 `options` 结构体变量，以此动态地设置 `options` 结构体变量的字段值。

同时，我们还可以在 `apply` 函数中自定义赋值逻辑，例如 `o.timeout = 10 * t`，使得设置结构体属性的灵活性更大。

Options 模式的优点：

- 支持传递多个参数，并在参数发生变化时保持兼容性
- 支持任意顺行传递参数
- 支持默认值
- 方便扩展
- 通过 `WithXXX` 的函数命名，可以使参数意义更加明确

{{&lt; admonition type=tip title=&#34;&#34; open=true &gt;}}
当结构体参数较少时，需要慎重考虑是否需要采用 Options 模式
{{&lt; /admonition &gt;}}


---

> Author: [kyden](https:github.com/kydance)  
> URL: http://kyden.us.kg/posts/design-pattern/  

