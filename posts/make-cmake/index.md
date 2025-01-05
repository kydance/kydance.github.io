# Make 与 CMake 实战指南：现代 C/C&#43;&#43; 构建系统精解


{{&lt; admonition type=abstract title=&#34;导语&#34; open=true &gt;}}
在现代软件开发中，高效的构建系统是项目成功的关键要素。Make 作为经典的构建工具，以其简洁的语法和强大的功能闻名；而 CMake 则通过其跨平台能力和灵活的配置，成为了现代 C/C&#43;&#43; 项目的标配工具。本文将带你深入了解这两个强大工具的使用方法，从基础概念到高级技巧，帮助你构建更加专业和高效的开发工作流。
{{&lt; /admonition &gt;}}

&lt;!--more--&gt;

## Makefile 特殊字符说明

- `$`: 主要用于变量引用，`$(CC)` 引用名为 `CC` 的变量
- `#`: 注释
- `:`: 分隔目标和依赖，`target: dependencies`
- `;`: 分隔命令，`target: dependencies; command1; command2`
- `=`: 变量赋值(递归展开赋值)，`CC = gcc`
- `:=`: 立即赋值，`VERSION := $(shell git describe)`
- `&#43;=`: 追加赋值，`CFLAGS &#43;= -O3 -Wall`
- `?=`: 条件赋值（如果变量未定义），`CC ?= gcc`
- `\`: 行继续符

    ```Makefile
    OBJS = main.o \
            helper.o \
            utils.o
    ```

- `%`: 通配符，用于模式规则，`%.o`, `%.c`
- `@`: 禁止命令回显，`@echo &#34;Building...`
- `$&lt;`: 第一个依赖项
- `$@`: 目标
- `$^`: 所有依赖项

    ```Makefile
    target: dep1 dep2
        command $&lt; $@ $^
    ```

- `-`: 忽略命令错误，`rm temp.txt`
- `*`: 通配符，匹配任意字符串，`*.o`
- `wildcard`、`patsubst` 等: 函数调用，`$(wildcard *.c)`


---

> Author: [kyden](https:github.com/kydance)  
> URL: http://kyden.us.kg/posts/make-cmake/  

