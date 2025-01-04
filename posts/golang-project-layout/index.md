# Golang Project Stardard Layout


{{&lt; admonition type=abstract title=&#34;导语&#34; open=true &gt;}}
Go 应用程序项目的基本布局介绍
{{&lt; /admonition &gt;}}

&lt;!--more--&gt;

## `Go Module`

从 Go 1.14 版本开始，除非存在特定不使用 `Go Modules` 的理由，否则请使用，并且一旦使用，就无需再担心 `$GOPATH` 以及项目的存放位置。

## Go 目录

### `/cmd`

`/cmd`，**本项目的主干**，其中每一个应用程序的目录名应该与你想要的可执行程序的名称相对应，例如 `/cmd/myApp`。

在 `/cmd` 目录下，不应该放置太多代码：

- 如果认为代码可以导入并可在其他项目中使用，那么它应该位于 `/pkg` 目录中.
- 如果代码不是可重用的，或者不希望其他人重用它，那么应该位于 `/internal` 目录中.

该目录下，通常有一个小的 `main` 函数，从 `/internal` 和 `pkg` 目录中导入和调用代码，除此之外没有别的东西.

微服务中的 app 服务类型分为4类：interface、service、job、admin

```bash
|---cmd
|   |---kydenapp-admin
|   |---kydenapp-interface
|   |---kydenapp-job
|   |---kydenapp-service
|   |---kydenapp-task
```

- interface: 对外的 BFF 服务，接受来自用户的请求，比如暴露了 HTTP/gRPC 接口
- service: 对内的微服务，仅接受来自内部其他服务或网关的请求，比如暴露了 gRPC 接口只对内服务
- admin: 区别于 service，更多是面向运营测的服务，通常数据权限更高，隔离带来更好的代码级别安全
- job: 流式任务处理的服务，上游一般依赖 message broker
- task: 定时任务，类似 cronjob，部署到 task 托管平台中

**`/cmd` 应用目录负责程序的: 启动、关闭、配置初始化等**

&gt; `DTO(Data Transfer Object)`，数据传输对象，这个概念来源于 J2EE 的设计模式，
&gt; 但这里泛指用于展示层/API层与服务层（业务逻辑层）之间的数据传输对象。

### `internal`

`internal`，**私有应用程序和库代码**，它是不希望其他人在其应用程序或库中导入的代码。
该目录由 Go 强制执行，确保私有包不可导入。

### `pkg`

`/pkg`，**外部应用程序可以使用的库代码（例如 `/pkg/mypubliclib`）**.

如果应用程序项目真的很小，并且额外的嵌套并不能增加多少价值(除非你真的想要:-)，那就不要使用它。
当它变得足够大时，根目录会变得非常繁琐时(尤其是当你有很多非 Go 应用组件时)，请考虑使用。

### `api`

`/api`，协议定义目录，(`xxapi.proto`) protobuf 文件，以及生成的 go 文件。
通常把 api 文档直接在 proto 文件中描述。

### `configs`

配置文件模版或默认配置

### `scripts`

执行各种构建、安装、分析等操作的脚本，是的根级别的 `Makefile` 变得小而简单

### `test`

额外的外部测试应用程序和测试数据

&gt; Go 会忽略以 `.` 或 `_` 开头的目录和文件

### `docs`

设计和用户文档（godoc 生成的文档除外）

### `tools`

项目的支持工具，可以从 `/pkg` 和 `/internal` 目录导入代码

### `examples`

应用程序和/或公共库的示例

### `third_party`

外部辅助工具，分叉代码和其他第三方工具（例如 `Swagger UI`）

### `assets`

与存储库一起使用的其他资源（图像、徽标等）

&gt; 按理来说我们不应该 `src` 目录，但有些 Go 项目拥有一个 `src` 文件夹，这通常发生在开发人员具有 Java 背景
&gt; `$GOPATH` 环境变量指向你的(当前)工作空间(默认情况下，它指向非 windows 系统上的 `$HOME/go`)，这个工作空间包括顶层 `/pkg`, `/bin` 和 `/src` 目录，而实际项目最终是 `/src` 下的一个子目录，即 `/xxx/workspace/src/proj/src/xxx.go`（Go 1.11 之后，项目 `proj` 可以放在 `GOPATH` 之外）.

## Reference

- [https://talks.golang.org/2014/names.slide](https://talks.golang.org/2014/names.slide)
- [https://golang.org/doc/effective_go.html#names](https://golang.org/doc/effective_go.html#names)
- [https://blog.golang.org/package-names](https://blog.golang.org/package-names)
- [https://go.dev/wiki/CodeReviewComments](https://go.dev/wiki/CodeReviewComments)
- [Style guideline for Go packages (rakyll/JBD)](https://rakyll.org/style-packages)


---

> : [kyden](https:github.com/kydance)  
> URL: http://kyden.us.kg/posts/golang-project-layout/  

