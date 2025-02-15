# Redis 核心精讲：从入门到性能优化


{{&lt; admonition type=abstract title=&#34;导语&#34; open=true &gt;}}
在当今高并发的互联网应用中，Redis 作为内存数据库和缓存系统的标配，以其卓越的性能和灵活的数据结构赢得了开发者的青睐。本文将带你深入了解 Redis 的核心特性，从五大数据类型的实战应用到单线程架构的性能优势，全方位提升你的 Redis 开发技能。无论是构建高性能缓存系统，还是开发实时数据处理应用，这都是一份不可或缺的实战指南。
{{&lt; /admonition &gt;}}

&lt;!--more--&gt;

## I. 简介

**Redis**(Remote Dictionary Service) 全称远程字典服务，一种**NoSQL** (Not Only SQL).

Redis is an in-memory data structure store used as a database,cache, message broker, and streaming engine.

---

## II. 基本数据类型

Redis 数据库中的每个键值对（Key-Value pair）都是由对象（Object）组成，其中：

- 数据库键（**Key**）：总是一个字符串对象（String Object）
- 数据库键对应的值（**Value**）：可以是
**字符串对象（String Object）**、**列表对象（List Object）**、
**哈希对象（Hash Object）**、**集合对象（Set Object）**、
**有序集合对象（Sorted Set Object）**中的一种

### String 字符串类型

Value 可以是字符串、也可以是数字

使用场景：计数（点赞数、粉丝数）、缓存

### List 列表类型

在 Redis 中，可以把 List 搞成队列、栈、阻塞队列.

List 的 Key 的底层实现就是一个链表，其中链表的每一个节点都保存了一个整数值.

Redis 链表实现的特性：

- 双向：链表节点都有 `prev` 和 `next` 指针 -&gt; 获取某个节点的前继节点和后继节点的复杂度都是O(1)
- 无环：链表头节点的 `prev` 指针和表尾节点的 `next` 指针都指向 `NULL`
- 表头指针 / 表尾指针：List 结构中存在 `head` 和 `tail` 指针
- 长度计数器：List 结构中存在 `len` 属性
- 多态：List 节点使用 `void*` 指针来保存节点值，并可以通过 List 结构中的
`dup`、`free`、`match`、`sane` 属性为节点值设置类型特定函数 -&gt; List 可以存储各种不同类型的值

使用场景：列表（关注列表、粉丝列表、消息列表，...）

---

## III. 常见问题

1. Redis 为什么是单线程？

    官方表示：Redis 是基于内存操作的，CPU不是Redis的性能瓶颈，Redis的瓶颈是根据机器的内存和网络带宽，既然可以用单线程实现，就没必要使用多线程了。

    Redis 的提供数据为 100000&#43; (10W&#43;) 的QPS，非常快

2. Redis 为什么单线程还这么快？

    多线程（CPU上下文切换）不一定比单线程效率高！！！Redis 是将所有的数据全部存放在内存中的，所有说单线程去操作效率就是最高的，多线程（CPU上下文切换，是一个耗时操作）；

    对于内存系统来说，如果没有上下文切换效率就是最高的！多次读写都是在一个CPU上完成的，在内存情况下就是最佳方案！

---

### IV. Reference


---

> Author: [kyden](https://github.com/kydance)  
> URL: http://kyden.us.kg/posts/redis/  

