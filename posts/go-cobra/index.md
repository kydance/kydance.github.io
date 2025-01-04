# 浅析现代化命令行框架 Cobra


{{&lt; admonition type=abstract title=&#34;导语&#34; open=true &gt;}}
**Cobra** 是一个可以创建强大的现代化 CLI 应用程序库，它还提供了一个可以生成应用和命令文件的程序的命令行工具：`cobra-cli`.
{{&lt; /admonition &gt;}}

&lt;!--more--&gt;

## I. Cobra 简介

**Cobra** 是一个可以创建强大的现代化 CLI 应用程序库，它还提供了一个可以生成应用和命令文件的程序的命令行工具：`cobra-cli`。
许多大型项目（e.g. kubernetes, Docker, Etcd, Rkt, Hugo etc.）都采用了 cobra 来构建他们的应用程序。

{{&lt; figure src=&#34;/posts/go-cobra/CobraMain.png&#34; title=&#34;&#34; &gt;}}

Cobra 具有很多特性，一些核心特性如下：

- 可以构建基于子命令的 CLI，并支持嵌套子命令：`app server`, `app fetch`
- 可以通过 `cobra-cli init appname &amp; cobra-cli add cmdname` 轻松生成应用和子命令
- 智能化命令建议：`app srver...did you mean app server`
- 自动生成命令和标志的 helpe 文本，并能自动识别 `-h`, `--help` 等标志
- 自动为应用程序生成 bash、zsh、fish、powershell 自动补全脚本
- 支持命令别名、自定义帮助、自定义用法等
- 可以与 viper、pflag 紧密集成，用于构建 12-factor 应用程序

Cobra 建立在 commands、arguments 和 flags 结构之上。Commands 代表命令，arguments 代表非选项参数，flags 代表选项参数（标志）。

{{&lt; admonition type=Tips title=&#34;CLI 模式&#34; open=true &gt;}}
一个好的应用程序应该是易懂的，用户可以清晰知道如何去使用这个应用程序，因此通常遵循如下模式：
`APPNAME VERB NOUN --ADJECTIVE` 或者 `APPNAME COMMAND ARG --FLAG`，例如：

```bash
# clone 是一个 Commands
# URL 是一个非选项参数
# bare 一个选项参数
git clone URL --bare
```

NOTE：`VERB` 代表动词，`NOUN` 代表名词，`ADJECTIVE` 代表形容词
{{&lt; /admonition &gt;}}

## II. `cobra-cli` 命令安装

Cobra 提供了 `cobra-cli` 命令，用来初始化一个应用程序并为其添加命令，方便开发基于 Cobra 的应用，可用以下方法进行安装：

```bash
$ go install github.com/spf13/cobra-cli@latest
# ...
```

`cobra-cli` 提供了 4 个子命令：

- `init`: 初始化一个 cobra 应用程序
- `add`: 给通过 cobra init 创建的应用程序添加子命令
- `completion`: 为指定的 shell 生成命令自动补全脚本
- `help`: 打印任意命令的帮助信息

`cobra-cli` 还提供了一些全局参数：

- `-a`, `--author`: 指定 Copyright 版权声明中的作者
- `--config`: 指定 cobra 配置文件的路径
- `-l`, `--license`: 指定生成的应用程序所使用的开源协议，内置的有：GPLv2, GPLv3, LGPL, AGPL, MIT, 2-Clause BSD or 3-Clause BSD；
- `--viper`: 使用 viper 作为命令行参数解析工具，默认为 true。

## III. Cobra 使用

在构建 cobra 应用时，可以自行组织代码目录结构，但 cobra 建议如下目录结构：

```bash
$ tree app_name
app_name
├── cmd
│   ├── add.go
│   ├── create.go
│   └── list.go
└── main.go
```

`main.go` 文件的目的只有一个：初始化 cobra 应用程序并注册子命令

```go
package main

import (
  &#34;{pathtToApp}/cmd&#34;
)

func main() {
  cmd.Execute()
}
```

### 使用 `cobra-cli` 命令生成应用程序并添加子命令

可以选择使用 `cobra-cli` 命令行工具快速生成一个应用程序，并添加子命令，然后基于生成的代码进行二次开发，提高开发效率，具体方法如下：

#### 1. 初始化应用程序

使用 `cobra-cli init` 命令初始化一个应用程序，然后就可以基于这个 Demo 进行二次开发，提高开发效率：

```bash
$ mkdir -p kyden-demo &amp;&amp; cd kydne-demo &amp;&amp; go mod init kyden-demo
$ cobra-cli init --license=MIT --viper
$ ls
cmd  go.mod  go.sum  LICENSE  main.go
```

#### 2. 添加子命令

当一个应用程序初始化完成之后，就可以使用 `cobra-cli add` 命令添加一些命令：

```bash
$ cobra-cli add serve
$ cobra-cli add config
$ cobra-cli add create -p &#39;configCmd&#39; # 此命令的父命令的变量名（默认为 &#39;rootCmd&#39;）

$ tree kyden-demo 
kyden-demo
├── LICENSE
├── cmd
│   ├── config.go
│   ├── create.go
│   ├── root.go
│   └── serve.go
├── go.mod
├── go.sum
└── main.go
```

执行 `cobra-cli add` 命令之后，会在 `cmd` 目录下生成命令源码文件。
`cobra-cli` 不仅可以添加命令，也可以添加子命令，例如 `cobra-cli add create -p &#39;configCmd&#39;` 给 `config` 命令添加了 `create` 子命令，`-p` 指定子命令的父命令：`&lt;父命令&gt;Cmd`.

#### 3. 编译运行

在生成命令后，可以直接执行 `go build` 命令编译应用程序：

```bash
$ go build -v .
$ ./kyden-demo -h
A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.

Usage:
  kyden-demo [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  config      A brief description of your command
  help        Help about any command
  serve       A brief description of your command

Flags:
      --config string   config file (default is $HOME/.kyden-demo.yaml)
  -h, --help            help for kyden-demo
  -t, --toggle          Help message for toggle

Use &#34;kyden-demo [command] --help&#34; for more information about a command.
```

#### 4. 配置 cobra

`cobra` 在生成应用程序时，也会在当前目录下生成 `LINCENSE` 文件，并且会在生成的 Go 源码文件中中，添加 LINCENSE Header。

LINCENSE 和 LINCENSE Header 的内容可以通过 cobra 配置文进行配置，默认配置文件 `~/.cobra.yaml`:

```bash
$ cat ~/.cobra.yaml
author: Kyden &lt;kytedance@gmail.com&gt;
year: 2024
license:
  header: This file is part of CLI application foo.
  text: |
    {{ .copyright }}

    This is my license. There are many like it, but this one is mine.
    My license is my best friend. It is my life. I must master it as I must
    master my life.

$ cobra-cli init
Copyright © 2024 Kyden &lt;kytedance@gmail.com&gt;

This is my license. There are many like it, but this one is mine.
My license is my best friend. It is my life. I must master it as I must
master my life.
```

`{{ .copyright }}` 的具体内容会根据 `author` 和 `year` 生成，根据此配置生成的 LICENSE 文件内容.

也可以使用内建的 licenses，内建的 licenses 有：GPLv2, GPLv3, LGPL, AGPL, MIT, 2-Clause BSD or 3-Clause BSD。

### 使用 `cobra` 库创建命令

当使用 cobra 库编码实现一个应用程序，需要首选创建一个空的 `main.go` 文件和一个 rootCmd 文件，然后根据需要添加其他命令。

具体步骤如下：

1. 创建 rootCmd

```bash
$ mkdir -p cobrademo &amp;&amp; cobrademo
$ go mod init cobrademo
go: creating new go.mod: module cobrademo
go: to add module requirements and sums:
        go mod tidy
$ cobra-cli init
Using config file: /Users/kyden/.cobra.yml
Your Cobra application is ready at
/tmp/cobrademo
```

通常情况下，会将 `rootCmd` 放在 `cmd/root.go` 文件中

```go
/*
Copyright © 2024 Kyden &lt;kytedance@gmail.com&gt;
This file is part of CLI application foo.
*/
package cmd

import (
	&#34;os&#34;

	&#34;github.com/spf13/cobra&#34;
)



// rootCmd represents the base command when called without any subcommands
var rootCmd = &amp;cobra.Command{
	Use:   &#34;cobrademo&#34;,
	Short: &#34;A brief description of your application&#34;,
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
    Run: func(cmd *cobra.Command, args []string) {
        // Do stuff here
    },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&amp;cfgFile, &#34;config&#34;, &#34;&#34;, &#34;config file (default is $HOME/.cobrademo.yaml)&#34;)

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP(&#34;toggle&#34;, &#34;t&#34;, false, &#34;Help message for toggle&#34;)
}
```

还可以在 `init()` 函数中定义标志和处理配置，例如：`cmd/helper.go`:

```go

```

2. 创建 `main.go`

还需要一个 main 函数来调用 rootCmd，通常会创建一个 `main.go` 文件，在 `main.go` 中调用 `rootCmd.Execute()` 来执行命令：

```go
package main

import (
  &#34;{pathToApp}/cmd&#34;
)

func main() {
  cmd.Execute()
}
```

在 `main.go` 中不建议放太多代码，通常只需要调用 `cmd.Execute()` 即可

3. 添加命令

除了 `rootCmd`，还可以调用 `AddCommand()` 来添加其他命令，通常情况下，会把其他命令的源码文件放在 `cmd` 目录下，例如添加一个 `version` 命令（`cmd/version.go`）：

```go
/*
Copyright © 2024 Kyden &lt;kytedance@gmail.com&gt;
This file is part of CLI application foo.
*/
package cmd

import (
	&#34;fmt&#34;

	&#34;github.com/spf13/cobra&#34;
)

// versionCmd represents the version command
var versionCmd = &amp;cobra.Command{
	Use:   &#34;version&#34;,
	Short: &#34;A brief description of your command&#34;,
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(&#34;version called&#34;)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// versionCmd.PersistentFlags().String(&#34;foo&#34;, &#34;&#34;, &#34;A help for foo&#34;)

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// versionCmd.Flags().BoolP(&#34;toggle&#34;, &#34;t&#34;, false, &#34;Help message for toggle&#34;)
}

```

4. 编译运行

```go
$ go build -v .
$ ./cobrademo -h
A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.

Usage:
  cobrademo [flags]
  cobrademo [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  helper      A brief description of your command
  version     A brief description of your command

Flags:
  -h, --help     help for cobrademo
  -t, --toggle   Help message for toggle

Use &#34;cobrademo [command] --help&#34; for more information about a command.
```

### 使用标志

cobra 可以跟 pflag 结合使用，实现强大的标志功能。
具体步骤如下：

1. 使用持久化的标志

标志是可以&#34;持久化&#34;的，即该标志可用于它所分配的命令以及该命令下的每个子命令。
例如，在 `rootCmd` 中定义持久化标志：

```go
rootCmd.PersistentFlags().BoolVarP(&amp;Verbose, &#34;verbose&#34;, &#34;v&#34;, false, &#34;verbose output&#34;)
```

2. 使用本地标志

本地标志，只能在其所绑定的命令上使用：

```go
rootCmd.Flags().StringVarP(&amp;Source, &#34;source&#34;, &#34;s&#34;, &#34;&#34;, &#34;Source directory to read from&#34;)
```

上面的 `--source` 标志智能在 `rootCmd` 命令上引用，而不能在 `rootCmd` 的子命令上引用。

3. 将标志绑定到 viper

可以讲标志绑定到 viper，这样就可以使用 `viper.Get()` 获取标志的值。

```go
var auther string

func init() {
	rootCmd.PersistentFlags().StringVar(
    &amp;auther, &#34;author&#34;, &#34;Your Name&#34;, &#34;Author name for copyright attribution&#34;)
	viper.BindPFlag(&#34;author&#34;, rootCmd.PersistentFlags().Lookup(&#34;auther&#34;))
}
```

4. 设置标志为必选

默认情况下，标志是可选的，也可以设置标志为必选。
当设置标志为必选时，若不提供标志时，cobra 会报错：

```go
rootCmd.Flags().StringVarP(&amp;Region, &#34;region&#34;, &#34;r&#34;, &#34;&#34;, &#34;AWS region (required)&#34;)
rootCmd.MarkFlagRequired(&#34;region&#34;)
```

## IV. Reference

- [cobra](https://github.com/spf13/cobra)


---

> : [kyden](https:github.com/kydance)  
> URL: http://kyden.us.kg/posts/go-cobra/  

