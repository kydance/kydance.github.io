# 深入解读 Golang 常用 Linter 工具及最佳实践


{{&lt; admonition type=abstract title=&#34;导语&#34; open=true &gt;}}
在 Golang 生态系统中，Linter 工具是开发者提升代码质量的关键。
本文将深入介绍几款常用的 Linter 工具及其最佳实践，帮助您在开发中避免常见错误并提高代码的可维护性。
{{&lt; /admonition &gt;}}

&lt;!--more--&gt;

## 目录

- [说明](#o-说明)
- [何为 Linter？](#i-何为-linter)
- [Gocyclo](#ii-gocyclo)
- [bodyclose](#iii-bodyclose)
- [sqlrows](#iv-sqlrows)
- [funlen](#v-funlen)
- [goconst](#vi-goconst)
- [ineffassign](#vii-ineffassign)
- [lll](#viii-lll)
- [errcheck](#ix-errcheck)
- [whitespace](#x-whitespace)
- [**GolangCI-Lint**](#xi-golangci-lint)
- [reviewdog](#xii-reviewdog)
- [Summary](#xiii-summary)
- [Reference](#xiv-reference)

## O. 说明

- 如特殊说明，文中代码已在在 Mac 和 Linux 系统下进行测试

## I. 何为 Linter？

Linter 是一种静态代码分析工具，用于在编译前检查代码中的错误、风格问题及潜在的 Bug。
在 Golang 生态中，Linter 工具帮助开发者在早期阶段就发现问题，从而避免后期修复的高成本。

---

## II. Gocyclo

Gocyclo 是一款用于分析 Go 代码中函数圈复杂度的 Linter 工具，帮助开发者识别需要重构的复杂函数。
通过降低圈复杂度，代码变得更加简洁、易读且更易维护。

### 函数圈复杂度(cyclomatic complexities)

圈复杂度，是一种衡量代码复杂性的指标，通过计算代码中的决策点（如if语句、循环等）来评估函数的复杂度，具体计算方法如下：

- 一个函数的基本圈复杂度为 `1`
- 当函数中存在的每一个 `if`, `for`, `case`, `&amp;&amp;` or `||`，都会使得该函数的圈复杂度加 `1`

&gt; 1. 在 Go 语言中，由于 `if err != nil` 的特殊情况存在，因此，其圈复杂度阈值默认为 15，而其他编程语言中圈复杂度阈值一般默认为 10。
&gt; 2. 在 Go 语言中，`switch` 中的 `default` 并不会增加函数的圈复杂度；

Gocyclo 可以作为单独的命令行工具使用，也可以与其他 Linter 工具(如 golangci-lint)集成使用，提供更全面的代码质量检查。
同时，它也可以集成到 CI/CD 流程中，帮助团队持续改善代码质量。

### 安装

```Bash
go install github.com/fzipp/gocyclo/cmd/gocyclo@latest
```

### 如何使用 Gocyclo linter ？

```Bash
Calculate cyclomatic complexities of Go functions.
Usage:
    gocyclo [flags] &lt;Go file or directory&gt; ...

Flags:
    -over N               show functions with complexity &gt; N only and
                          return exit code 1 if the set is non-empty
    -top N                show the top N most complex functions only
    -avg, -avg-short      show the average complexity over all functions;
                          the short option prints the value without a label
    -ignore REGEX         exclude files matching the given regular expression

The output fields for each line are:
&lt;complexity&gt; &lt;package&gt; &lt;function&gt; &lt;file:line:column&gt;
```

#### 使用示例

```Go
// gocyclo-test/main.go 
package main

import (
 &#34;fmt&#34;
 &#34;strconv&#34;
)

func main() {
    var a = 10
    if a == 10 {
        f()
    } else {
        fmt.Printf(&#34;%s&#34;, strconv.Itoa(a))
    }

    switch a{
    case 10:
        fmt.Println(a)
    default:
        fmt.Println(&#34;default&#34;)
    }
}

func f() {
    a := 10
    b := 12

    if a != b {
        // do something
        fmt.Println(&#34;a != b&#34;)
    }
}
```

```Bash
$ gocyclo gocyclo-test/main.go 
3 main main gocyclo-test/main.go:8:1
2 main f gocyclo-test/main.go:24:1
```

---

## III. bodyclose

在 Go 中，即使读取了所有的响应内容，也需要显式关闭响应体以释放资源，否则可能导致资源泄漏、连接池耗尽，进而影响应用性能。

`bodyclose` 主要关注于 HTTP 响应体的正确关闭，通过检查 `resp.Body` 是否被正确关闭。
它既可以单独使用，也可以集成到其他 linter 工具（例如 golangci-lint）中。

### 安装

```Bash
go install github.com/timakin/bodyclose@latest
```

### 如何使用 bodyclose ?

```Bash
$ bodyclose
bodyclose is a tool for static analysis of Go programs.

Usage of bodyclose:
	bodyclose unit.cfg	# execute analysis specified by config file
	bodyclose help    	# general help, including listing analyzers and flags
	bodyclose help name	# help on specific analyzer and its flags
```

#### 使用示例

这里展示借助 `golangci-lint` 的方式使用 `bodyclose`.

```Go
// main.go
package kyden

import (
 &#34;fmt&#34;
 &#34;io&#34;
 &#34;net/http&#34;
)

func f() error{
    resp, err := http.Get(&#34;http://example.com/&#34;)
    if err != nil {
        return err
    }
    // defer resp.Body.Close() // &lt;&lt;&lt;

    body, err := io.ReadAll(resp.Body)
    fmt.Println(body)
    return nil
}
```

```Bash
$ golangci-lint run --disable-all -E bodyclose main.go
main.go:11:26: response body must be closed (bodyclose)
    resp, err := http.Get(&#34;http://example.com/&#34;)
```

&gt; 避免使用 `http` 库中 `body` 忘记 `close` 的更优方案是:
&gt;
&gt; **对 Go 官方提供的 `http` 进行封装，使调用方（Caller）不用显示调用 `close` 函数.**
&gt;
&gt; ```go
&gt; package httpclient
&gt; 
&gt; import (
&gt;     &#34;io/ioutil&#34;
&gt;     &#34;net/http&#34;
&gt; )
&gt; 
&gt; // Client 是一个自定义的 HTTP 客户端结构体
&gt; type Client struct {
&gt;     http.Client
&gt; }
&gt; 
&gt; // Get 封装了 http.Get 方法
&gt; func (c *Client) Get(url string) (string, error) {
&gt;     resp, err := c.Client.Get(url)
&gt;     if err != nil {
&gt;         return &#34;&#34;, err
&gt;     }
&gt;     
&gt;     // 确保在函数返回时关闭响应体
&gt;     defer resp.Body.Close()
&gt; 
&gt;     // 读取响应内容
&gt;     body, err := ioutil.ReadAll(resp.Body)
&gt;     if err != nil {
&gt;         return &#34;&#34;, err
&gt;     }
&gt; 
&gt;     return string(body), nil
&gt; }
&gt; ```

---

## IV. sqlrows

在 Go 的 `database/sql` 包中，`sql.Rows` 是一个 `struct`，用于表示从数据库查询中返回的多行结果。

它提供了一组方法，允许开发者逐行读取查询结果。

- 迭代结果：使用 `Next()` 方法逐行遍历结果集。
- 扫描数据：使用 `Scan()` 方法将当前行的列值复制到指定的变量中。
- 关闭结果集：使用 `Close()` 方法释放与结果集相关的资源。

`sqlrows` 的[官方介绍](https://github.com/gostaticanalysis/sqlrows)：
`sqlrows` is a static code analyzer which helps uncover bugs by reporting a diagnostic for mistakes of `sql.Rows` usage.

### 安装

```Bash
go install github.com/gostaticanalysis/sqlrows@latest
```

### 如何使用 sqlrows ?

```Bash
$ sqlrows
sqlrows is a tool for static analysis of Go programs.

Usage of sqlrows:
	sqlrows unit.cfg	# execute analysis specified by config file
	sqlrows help    	# general help
	sqlrows help name	# help on specific analyzer and its flags
```

Go 源码【注意 Not Good(NG) 处】

```Go
// main.go
package kyden

import (
 &#34;context&#34;
 &#34;database/sql&#34;
)

func f(ctx context.Context, db *sql.DB) (interface{}, error) {
    rows, err := db.QueryContext(ctx, &#34;SELECT * FROM users&#34;)
    defer rows.Close() // NG: using rows before checking for errors

    if err != nil {
        return nil, err
    }
    // defer rows.Close() // NG: this return will not release a connection.

    for rows.Next() {
        err = rows.Scan()
        if err != nil {
            return nil, err
        }
    }
    return nil, nil
}
```

针对两种 NG 的不同输出：

```Bash
go vet -vettool=$(which sqlrows) main.go
# command-line-arguments
./main.go:10:11: using rows before checking for errors
```

```Bash
go vet -vettool=$(which sqlrows) main.go
# command-line-arguments
./main.go:9:33: rows.Close must be called
```

---

## V. funlen

`funlen`，用于检查函数的长度，确保函数的可读性和可维护性。
默认情况下，funlen 将函数的最大行数限制(`lines`)为 60 行，最大语句数(`statements`)限制为 40 条。

通常，funlen 会结合 golangci-lint 使用， 并集成到开发工作流中，提升代码质量.

### 安装

funlen 可以通过 golangci-lint 安装: `go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest`

### 如何使用 funlen ？

```yml
linters:
  disable-all: true
  enable:
    - funlen

linters-settings:
  funlen:
    lines: 60
    statements: 40
```

#### 使用示例

```go
// main.go
package main

import (
	&#34;fmt&#34;
)

func main() {
    f()
}

func f () {
    fmt.Println(&#34;Test funlen&#34;)

    a := 1
    fmt.Println(a)

    b := 1
    fmt.Println(b)

    c := 1
    fmt.Println(c)
}
```

下面的 `.golangci.yml` 仅用于展示 funlen 的用法，具体参数请根据实际项目自行调整。

```yml
# .golangci.yml 
linters:
  disable-all: true
  enable:
    - funlen

linters-settings:
  funlen:
    lines: 6
    statements: 4
```

```Bash
$ golangci-lint run
main.go:12: Function &#39;f&#39; has too many statements (7 &gt; 4) (funlen)
```

---

## VI. goconst

goconst 会扫描代码，识别出在多个地方重复出现的字符串。
这些字符串通常是相同的文本，开发者通过将重复的字符串提取为常量，代码变得更加清晰，减少了硬编码的出现，降低了出错的可能性。
可以根据项目需求自定义 goconst 的行为，例如设置字符串的最小长度、最小出现次数等。

goconst 通常作为 golangci-lint 的一部分使用。

### 如何使用 goconst ?

```yml
linters:
  disable-all: true
  enable:
    - goconst

linters-settings:
  goconst:
    min-len: 3
    min-occurrences: 3
```

#### 使用示例

```go
// main.go
package main

import &#34;fmt&#34;

func f() {
    a := &#34;Hello&#34;
    fmt.Println(a)

    b := &#34;Hello&#34;
    fmt.Println(b)

    c := &#34;Hello&#34;
    fmt.Println(c)
}
```

下面的 `.golangci.yml` 仅用于展示 funlen 的用法，具体参数请根据实际项目自行调整。

```yml
# .golangci.yml
linters:
  disable-all: true
  enable:
    - goconst

linters-settings:
  goconst:
    min-len: 3
    min-occurrences: 3
```

```Bash
$ golangci-lint run
main.go:7:10: string `Hello` has 3 occurrences, make it a constant (goconst)
    a := &#34;Hello&#34;
         ^
```

---

## VII. ineffassign

ineffassign，主要用于检测代码中对现有变量的赋值操作是否未被使用。
这种未使用的赋值通常是代码中的潜在错误，可能导致逻辑上的混乱或资源的浪费。

### 如何使用 ineffassign ?

通常作为 golangci-lint 的一部分使用。

```yml
linters:
  disable-all: true
  enable:
    - ineffassign
```

#### 使用示例

```go
// main.go
package main

import &#34;fmt&#34;

func f() {
    a := &#34;Hello&#34;

    // ...
    // Not assign a value to `a`
    // ...

    a = &#34;kyden&#34;
    fmt.Println(a)
}
```

```Bash
$ golangci-lint run
main.go:7:5: ineffectual assignment to a (ineffassign)
    a := &#34;Hello&#34;
    ^
```

---

## VIII. lll

通过限制行的长度，lll 有助于确保代码在查看时不会横向滚动，提升代码的可读性。

lll，主要用于检查代码行的长度，检查每一行的长度是否超过指定的最大值。
默认情况下，lll 将最大行长度限制为 120 个字符。

### 如何使用 lll ?

lll 通常作为 golangci-lint 的一部分使用。

```yml
linters:
  disable-all: true
  enable:
    - lll

linters-settings:
  lll:
    line-length: 80
```

#### 使用示例

```go
// main.go
package kyden

func f() int {
    a := &#34;This is a very long line that exceeds the maximum line length set by the linter and should be broken up into smaller, more manageable lines.&#34;
    return len(a)
}
```

```Bash
golangci-lint run
main.go:5: the line is 151 characters long, which exceeds the maximum of 80 characters. (lll)
    a := &#34;This is a very long line that exceeds the maximum line length set by the linter and should be broken up into smaller, more manageable lines.&#34;
```

&gt; 解决方案
&gt;
&gt;使用反引号（`）定义多行字符串，允许字符串跨越多行而不需要使用连接符

---

## IX. errcheck

errcheck，专门检查未处理的错误，确保开发者在调用可能返回错误的函数时，正确地检查和处理这些错误，从而提高代码的健壮性和可靠性。

- `errcheck` 会扫描 Go 代码，查找未检查错误的地方
- 除了检查函数返回的错误,还可以检查类型断言是否被忽略
- 可以检查是否将错误赋值给了空白标识符

### 如何使用 ?

`errcheck` 通常作为 golangci-lint 的一部分使用

```yml
linters-settings:
  errcheck:
    check-type-assertions: true # 检查类型断言是否被忽略,默认为 false
    check-blank: true # 检查是否将错误赋值给空白标识符,默认为 false
    disable-default-exclusions: true # 禁用默认的忽略函数列表,默认为 false
    exclude-functions:  # 指定要忽略检查的函数列表
        # ...
```

#### 使用示例

```go
// main.go
package main

import (
	&#34;fmt&#34;
)

func main() {
    hello(&#34;Kyden&#34;) // err Not Check

    _ = hello(&#34;Kyden&#34;) // err assign to _

    err := hello(&#34;Go&#34;)
    if err != nil {
        return
    }
}

func hello(str string) error {
    fmt.Printf(&#34;Hello, %s&#34;, str)

    return nil
}

```

下面的 `.golangci.yml` 仅用于展示 errcheck 的用法，具体参数请根据实际项目自行调整。

```yml
# .golangci.yml
linters:
  disable-all: true
  enable:
    - errcheck

linters-settings:
  errcheck:
    check-type-assertions: true
    check-blank: true
```

```Bash
golangci-lint run
main.go:9:10: Error return value is not checked (errcheck)
    hello(&#34;Kyden&#34;) // err Not Check
         ^
main.go:11:5: Error return value is not checked (errcheck)
    _ = hello(&#34;Kyden&#34;) // err assign to _
    ^
```

---

## X. whitespace

`whitespace` 是一个 Go 语言的 linter，主要用于检查代码中不必要的空行，即检查函数、条件语句（如 `if`、`for`）等开头和结尾的多余空行。

### 如何使用 whitespace ?

`whitespace` 也包含在 golangci-lint 中，只需在配置中启用即可。

```yml
linters:
  disable-all: true
  enable:
    - whitespace
```

#### 使用示例

```go
// main.go
package main

import (
	&#34;fmt&#34;
)

func main() {
    err := hello(&#34;Kyden&#34;)
    if err != nil {
        return
    }
}

func hello(str string) error {

    if len(str) &lt;= 0 {

        return fmt.Errorf(&#34;str len &lt;= 0&#34;)
    }
    fmt.Printf(&#34;Hello, %s&#34;, str)

    return nil

}
```

```Bash
$ gosrc golangci-lint run
main.go:15:31: unnecessary leading newline (whitespace)

^
main.go:25:1: unnecessary trailing newline (whitespace)

^
main.go:17:23: unnecessary leading newline (whitespace)

^
```

---

## XI. GolangCI-Lint

&gt; **生产级静态分析工具**
&gt;
&gt; [`golangci-lint` is a fast Go linters runner. It runs linters in parallel, uses caching, supports YAML configuration, integrates with all major IDEs, and includes over a hundred linters.](https://golangci-lint.run/)

`golangci-lint` 是一款快速的 Go 语言 linter，它并行运行多个 linter 程序，使用缓存，支持 YAML 配置，与所有主流集成开发环境集成，并包含一百多个 linter 程序。

### 安装

```Bash
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Verify
golangci-lint --version
```

### 如何使用 golangci-lint ?

在不进行任何配置的情况下，GolangCI-Lint 将默认采用启动以下 Linters:
`errcheck`, `gosimple`, `govet`, `ineffassign`, `staticcheck`, `unused`.

也可以通过传递 `-E`(`--enable`) 参数来启动 Linter，传递 `-D`(`--disable`) 来禁用 Linter.

```Bash
golangci-lint run --disable-all -E errcheck
```

{{&lt; figure src=&#34;/posts/golang-linter/golangci-lint-default.png&#34; title=&#34;&#34; &gt;}}

### Visual Studio Code 集成

由于个人一直使用 VSCode 开发各种程序，这里只展示其如何集成 GolangCI-Lint。

Step 1. **`settings.json` 启用 golangci-lint**

```json
&#34;go.lintTool&#34;: &#34;golangci-lint&#34;,
&#34;go.lintFlags&#34;: [
  &#34;--fast&#34; // Using it in an editor without --fast can freeze your editor.
]
```

---

Step 2. 配置 `.golangci.yml`

当使用 Golangci-lint 时，它会自动在编辑的 Go 文件所在的目录或父目录中查找 `.golangci.yml` 配置文件。
如果找到了配置文件，Golangci-lint 就会根据该配置文件的设置来运行 linter。

因此，在 VS Code 的设置中，不需要专门配置 Golangci-lint。
**只需要在项目根目录或相应的目录下创建 `.golangci.yml` 配置文件，并在其中指定需要启用的 linter 和相关选项即可**。

---

Step 3. **Enjoy your coding time 🥂**

---

&gt; [Golangci-lint 同样支持 GoLang、NeoVim 等流行 IDE 集成.](https://golangci-lint.run/welcome/integrations/)

### `.golangci.yml` 参考配置

这里给出一个个人在用的 golangci-lint 完整配置文件，以供参考：

```yml
run:
  timeout: 5m
  go: 1.21

linters-settings:
  funlen:
    lines: 150
    statements: 100
  goconst:
    min-len: 3
    min-occurrences: 3
  lll:
    line-length: 80
  govet:            # 对于linter govet，这里手动开启了它的某些扫描规则
    shadow: true
    check-unreachable: true
    check-rangeloops: true
    check-copylocks: true
    # 启动nilness检测
    enable:
      - nilness

linters:
  disable-all: true
  enable:
    - bodyclose
    - errcheck
    - funlen
    - goconst
    - gocyclo
    - gofmt
    - goimports
    - gosimple
    - govet
    - ineffassign
    - lll
    - misspell # Go 静态分析工具，专注于检查代码中的拼写错误
    - nilerr
    - rowserrcheck
    - staticcheck
    - typecheck
    - unconvert
    - unparam
    - unused
    - whitespace

issues:
  skip-dirs:
    - test

  exclude-files:
    - _test.go
```

更多详细信息，请参考[官方文档](https://golangci-lint.run/)

## XII. reviewdog

A code review dog who keeps your codebase healthy.

`reviewdog` 是一个用于自动化代码审查的工具，旨在通过集成各种 linter 工具来简化代码质量检查。它能够将 lint 工具的输出结果作为评论发布到代码托管服务（如 GitHub、GitLab 等），从而提高代码审查的效率和准确性。

### 功能

- 自动发布评论：reviewdog 可以将 lint 工具的结果自动发布为评论，帮助开发者快速识别代码中的问题。
- 支持多种 linter：它支持多种静态分析工具，包括 golangci-lint、eslint、pylint 等，可以方便地集成到现有的开发流程中。
- 过滤输出：支持根据 diff 过滤 lint 工具的输出，只报告在当前变更中出现的问题。
- 多种报告模式：支持多种报告模式，如 GitHub PR 评论、GitHub Checks、GitLab 合并请求讨论等。
- 本地运行：除了在 CI/CD 环境中运行外，reviewdog 也可以在本地环境中使用，方便开发者在提交代码前进行检查。

### 安装

```bash
# Install the latest version. (Install it into ./bin/ by default).
$ curl -sfL https://raw.githubusercontent.com/reviewdog/reviewdog/master/install.sh | sh -s

# Specify installation directory ($(go env GOPATH)/bin/) and version.
$ curl -sfL https://raw.githubusercontent.com/reviewdog/reviewdog/master/install.sh | sh -s -- -b $(go env GOPATH)/bin [vX.Y.Z]

# In alpine linux (as it does not come with curl by default)
$ wget -O - -q https://raw.githubusercontent.com/reviewdog/reviewdog/master/install.sh | sh -s [vX.Y.Z]
```

推荐使用第二种安装方式 `curl -sfL https://raw.githubusercontent.com/reviewdog/reviewdog/master/install.sh | sh -s -- -b $(go env GOPATH)/bin`，具体安装实例如下：

```bash
$ curl -sfL https://raw.githubusercontent.com/reviewdog/reviewdog/master/install.sh | sh -s -- -b $(go env GOPATH)/bin
reviewdog/reviewdog info checking GitHub for latest tag
reviewdog/reviewdog info found version: 0.20.1 for v0.20.1/Darwin/arm64
reviewdog/reviewdog info installed /Users/kyden/go/bin/reviewdog
```

### 如何使用 reviewdog ?

#### 本地使用

```Bash
golangci-lint run ./... 2&gt;&amp;1 | reviewdog -f=golangci-lint -reporter=local
```

&gt; [官方示例](https://github.com/reviewdog/reviewdog?tab=readme-ov-file#reporter-local--reporterlocal-default)

### Github Action

#### 1. 创建 GitHub Actions 工作流

在项目根目录下创建一个 GitHub Actions 工作流文件，`.github/workflows/reviewdog.yml`

#### 2. 配置 .golangci.yml

在项目根目录下创建一个 `.golangci.yml` 配置文件，配置需要启用的 linter

#### 3. 提交代码

当你提交代码并创建拉取请求时，GitHub Actions 会自动运行 reviewdog，并根据 lint 工具的输出在拉取请求中添加评论，指出代码中的问题。

&gt; [更多内容请参考官方示例](https://github.com/reviewdog/reviewdog?tab=readme-ov-file#github-actions)

## XIII. Summary

**综上所述，Golang 生态中有众多优秀的 Linter 工具，它们能够有效地检查代码质量，提高项目的可维护性和可靠性。
开发者可以根据项目需求，选择合适的 Linter 工具，并将其集成到 CI/CD 流程中，持续改善代码质量。
未来，随着 Golang 社区的不断发展，相信会有更多优秀的 Linter 工具问世，为 Golang 开发者提供更加强大的代码分析能力。**

## XIV. Reference

- [Cyclomatic complexity](https://en.wikipedia.org/wiki/Cyclomatic_complexity)
- [Gocyclo](https://github.com/fzipp/gocyclo)
- [bodyclose](https://github.com/timakin/bodyclose)
- [sqlrows](https://github.com/gostaticanalysis/sqlrows)
- [GolangCI-Lint](https://github.com/golangci/golangci-lint)
- [static analysis](https://github.com/analysis-tools-dev/static-analysis)
- [reviewdog](https://github.com/reviewdog/reviewdog)


---

> : [kyden](https:github.com/kydance)  
> URL: http://kyden.us.kg/posts/golang-linter/  

