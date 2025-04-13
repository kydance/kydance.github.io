# Make 构建系统详解：从基础到实战


{{< admonition type=abstract title="导语" open=true >}}
在现代软件开发中，高效的构建系统是项目成功的关键要素。Make 作为经典的构建工具，以其简洁的语法和强大的功能闻名；本文将带你深入了解这个强大工具的使用方法，从基础概念到高级技巧，帮助你构建更加专业和高效的开发工作流。
{{< /admonition >}}

<!--more-->

## Makefile 特殊字符说明

- `$`: 主要用于变量引用，`$(CC)` 引用名为 `CC` 的变量
- `#`: 注释
- `:`: 分隔目标和依赖，`target: dependencies`
- `;`: 分隔命令，`target: dependencies; command1; command2`
- `=`: 变量赋值(递归展开赋值)，`CC = gcc`
- `:=`: 立即赋值，`VERSION := $(shell git describe)`
- `+=`: 追加赋值，`CFLAGS += -O3 -Wall`
- `?=`: 条件赋值（如果变量未定义），`CC ?= gcc`
- `\`: 行继续符

    ```Makefile
    OBJS = main.o \
            helper.o \
            utils.o
    ```

- `%`: 通配符，用于模式规则，`%.o`, `%.c`
- `@`: 禁止命令回显，`@echo "Building...`
- `$<`: 第一个依赖项
- `$@`: 目标
- `$^`: 所有依赖项

    ```Makefile
    target: dep1 dep2
        command $< $@ $^
    ```

- `-`: 忽略命令错误，`rm temp.txt`
- `*`: 通配符，匹配任意字符串，`*.o`
- `wildcard`、`patsubst` 等: 函数调用，`$(wildcard *.c)`


---

> Author: [kyden](https://github.com/kydance)  
> URL: http://kydance.github.io/posts/make/  

