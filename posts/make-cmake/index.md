# 浅析 Make 与 Cmake


{{&lt; admonition type=abstract title=&#34;导语&#34; open=true &gt;}}
导语内容
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

