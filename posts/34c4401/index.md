# MCP


{{< admonition type=abstract title="导语" open=true >}}
**这是导语部分**
{{< /admonition >}}

<!--more-->

## 模型上下文协议 MCP

MCP（Model Context Protocol，模型上下文协议），
起源于2024年11月25日[Anthropic](https://www.anthropic.com/)发布的文章：
[Introducing the Model Context Protocol](https://www.anthropic.com/news/model-context-protocol)。

MCP 在AI领域中的作用，就犹如扩展坞在笔记本电脑中的作用一样，为AI模型的能力扩展提供了新的可能性。
例如，现代扩展坞可以连接显示器、键盘以及移动硬盘等多种硬件外设，使得笔记本电脑瞬间扩展了多种功能；
而 MCP Server 将会作为一个智能中枢平台，可以动态接入各类专业能力模块（如知识库、计算工具、领域模型等）。

当 LLM 需要完成特定任务时，它就可以调用这些 MCP Server，实时获得精准的上下文支持，从而实现能力的弹性扩展。
这种架构打破了传统 AI 模型的封闭性，让大语言模型像搭载了多功能扩展坞的超级工作站，随时都能获取最合适的专业工具。

### 关键概念

- MCP Server: Agent，自己实现、用来响应请求的“服务器”
- MCP Client: 会调用 MCP Server 的客户端
- Resources: 资源，MCP Server 消费的 API 响应、文件等，最终作为响应返回给 CMP Client
- Tools: 工具，MCP Client 可调用的工具/方法

### JSON-RPC

#### JSON-RPC 与 gRPC 的主要区别

- JSON-RPC 没有强制的 client stub 生成机制，采用更加动态的方式调用服务。
而 MCP 是透过 lift cycle 来动态获取 Capabilities，但这并非 JSON-RPC 规范本身的一部分;
- JSON-RPC 使用 JSON 进行编码（动态弱型别），而 gRPC 使用 Protocol Buffers（静态强型别）
- JSON-RPC 可以透过特定实现提供的机制来动态发现服务能力，而 gRPC 通常在编译时确定

#### JSON-RPC 的 Request 与 Response 格式

Request，与 HTTP method 不太一样的是，HTTP API 都会把 method 声明在 URL 中，
而 JSON-RPC 则是把 POST + body 中的 payload 会注明 method。

```json
{
    jsonrpc: "2.0";
    id: string | number;
    method: string;
    params?: {
        [key: string]: unknown;
    };
}
```

Response，会把 Request 中的 id 带到 Response 中（RPC 的 client 能批量发送请求，与 HTTP API 不同）

```json
{
    jsonrpc: "2.0";
    id: string | number;
    result?: {
        [key: string]: unknown;
    }
    error?: {
        code: number;
        message: string;
        data?: unknown;
    }
}
```

> 批量请求
>
> gRPC 是通过 Streaming 的形式来批量请求
>
> ```json
> Requests :
> [
>   {"jsonrpc": "2.0", "method": "sum", "params": [1,2], "id": 1},
>   {"jsonrpc": "2.0", "method": "subtract", "params": [42,23], "id": 2},
>   {"jsonrpc": "2.0", "method": "foo.get", "params": {"name": "myself"}, "id": 3}
> ]
> 
> Responses:
> [
>   {"jsonrpc": "2.0", "result": 3, "id": 1},
>   {"jsonrpc": "2.0", "error": {"code": -32601, "message": "Method not found"}, "id": 2},
>   {"jsonrpc": "2.0", "result": {"firstName": "John", "lastName": "Doe", "age": 30}, "id": 3}
> ]
> ```

### Tools

> Like resources, tools are identified by unique names and can include descriptions to guide their usage. However, unlike resources, tools represent dynamic operations that can modify state or interact with external systems.

### Resource Capability

> Each resource is identified by a unique URI and can contain either text or binary data.

Resource，在 MCP Server 中用来表示 Client 能够存取该 Server 哪些内部资源或提供哪些定制化的请求，
并且把这些内容当作跟 LLM 交互的上下文来使用，使得 AI 能够给我们提供更精准的应答。
因此，Resource 决定了不同的 MCP Client 在合适（When）以及如何（How）被使用。

常见的 Resource 如：`File contents`, `Database records`, `API response`,`Live system data`,`Screenshots and images`,`Log files` 等等。

**Tools 能做到修改状态（需要额外给予LLM R/W 权限），
而 Resource 不行（额外给予 LLM ReadOnly 权限），且 Resource 语意上主要是提供结构化资料呀以及所在的 URI 提供存取。**

> 可以把 Resource 理解成能回传哪些 database、tables 等，而 Tools 则是允许你去操作（CRUD、Calculate）

#### method

- `resource/list`
- `resource/read`
- `resource/subscribe`
- `notifications/resource/update`

{{< figure src="/posts/MCP/mcp-resource.png" title="" >}}

### Prompt Capability

Prompt 与 Tools 作用一样，就是预先设置好的工作流程跟 Prompt 模版，客制化一些场景下可提供给AI进行采纳的 Prompt，如果需要 AI 客制化回传某些格式化内容时，可以提供自定义的 Prompt。

Prompt 的 Metadata 如下：

```json
{
  name: string;              // Unique identifier for the prompt
  description?: string;      // Human-readable description
  arguments?: [              // Optional list of arguments
    {
      name: string;          // Argument identifier
      description?: string;  // Argument description
      required?: boolean;    // Whether argument is required
    }
  ]
}
```

- name: 在同一个 Server 内必须唯一
- arguments：每一个 argument 内部都是一个物件，
其中的 name 在每个 prompt 中也必须是唯一的，
可配合 required 属性表明是否 Client 端提供。最重要的还是 `description` 描述要清晰.

> 跟 Tools 和 Resource 类似，Prompt 一样需要提供 `prompts/list` 获取所有可用的 prompt；
> 使用 `prompts/get` 获取指定 prompt 的模版内容

在 Prompt 模板中，我們可以設定 2 种角色（Role）給LLM使用，分別是 user 与 assistant:

- **User**: 这个角色代表与系统交互的使用者，也就是提问的人所发出的信息；需要提供上下文或是问题给 LLM 来产生回应；
- **Assistant**: 则代表 LLM 本身，负责处理使用者的请求并产生响应；主要根据内部逻辑和 context 进行推理，提供建议或解决方案，能维持整串对话的连贯性与对上下文的理解。

### MCP Transport

#### stdio 标准输入输出

MCP 提供多种 Client 与 Server 交互的方式，其中最常见的就是 stdio，且交互数据格式仍满足 JSON-RPC 协议

{{< figure src="/posts/MCP/server-stdio.png" title="" >}}

#### HTTP with SSE

SSE(Server-Sent Events)，是通过 HTTP API 的形式来调用 JSON-RPC Server。

{{< figure src="/posts/MCP/server-sse.png" title="" >}}

#### Streamable HTTP

**TODO**

### MCP 链接的 Life Cycle

Client 从建立链接开始到结束会经历 3 个阶段，其中最重要的是 Initialize 阶段的准备，这个阶段会提供 MCP Server 具体功能清单

1. **Initialize**: 初始化阶段，能力协商和协议版本协商

2. **Operator**： 正常协议通信

3. **Shutdown**： 正常终止链接

{{< figure src="/posts/MCP/mcp-lifecycle.png" title="" >}}

## Server Inspector

MCP 官方提供 [Server Inspector](https://github.com/modelcontextprotocol/inspector)
来协助 MCP 开发者对 Server 进行简单的除错调用，其作用类似于 PostMan、BloomRPC 这类的工具。

> MCP Inspector 是专为 Model Context Protocol（MCP）服务器设计的交互式调试工具，支持开发者通过多种方式快速测试与优化服务端功能

### 安装、启动

Inspector 直接通过 npx 运行，无需安装

```bash
npx @modelcontextprotocol/inspector <command>

# 传递参数 arg1 arg2
npx @modelcontextprotocol/inspector <command> <arg1> <arg2>

# 传递环境变量 KEY=value KEY2=$VALUE2
npx @modelcontextprotocol/inspector -e KEY=value -e KEY2=$VALUE2 <command> <arg1> <arg2>
```

- `@` 表示这是一个 **作用域包**，用于明确包的归属组织或用途，`@modelcontextprotocol` 是组织名，`inspector` 是包名
- `<command>` 是指让 MCP Server 执行起来的命令
- `<arg...>` 是 MCP Server 执行时对应的命令行参数

### Server 运行

Inspector 运行起来后，会启动两个服务：

- 带 UI 的客户端，默认端口：5173
- MCP Proxy Server，默认端口：3000

## MCP Server 中的 `hollo world`

```Go
// main.go

package main

import (
 "context"
 "flag"
 "fmt"
 "log"
 "os"

 "github.com/mark3labs/mcp-go/mcp"
 "github.com/mark3labs/mcp-go/server"
)

func main() {
 var transport string
 flag.StringVar(&transport, "t", "stdio", "Transport type (stdio or sse)")
 flag.StringVar(
  &transport,
  "transport",
  "stdio",
  "Transport type (stdio or sse)",
 )
 addr := flag.String("sse-address", "localhost:5568", "The host and port to start the sse server on")
 flag.Parse()
 fmt.Println(*addr)

 if err := run(transport, *addr); err != nil {
  panic(err)
 }
}

func run(transport, addr string) error {
 // Create MCP server with explicit options
 s := server.NewMCPServer(
  "Demo 🚀",
  "1.0.0",
 )

 // Add tool with more explicit configuration
 tool := mcp.NewTool("hello_world",
  mcp.WithDescription("Say hello to someone"),
  mcp.WithString("name",
   mcp.Required(),
   mcp.Description("Name of the person to greet"),
  ),
 )

 // Add tool handler
 s.AddTool(tool, helloHandler)
 // s.AddTools(server.ServerTool{Tool: tool, Handler: helloHandler})

 // Debug information
 log.Printf("Registered tool: hello_world")

 switch transport {
 case "stdio":
  srv := server.NewStdioServer(s)
  return srv.Listen(context.Background(), os.Stdin, os.Stdout)
 case "sse":
  // Create the SSE server with explicit debugging
  srv := server.NewSSEServer(s)

  log.Printf("SSE server listening on %s", addr)
  if err := srv.Start(addr); err != nil {
   return fmt.Errorf("Server error: %v", err)
  }
  // This code is unreachable as Start() blocks until error
 default:
  return fmt.Errorf(
   "Invalid transport type: %s. Must be 'stdio' or 'sse'",
   transport,
  )
 }
 return nil
}

func helloHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
 name, ok := request.Params.Arguments["name"].(string)
 if !ok {
  return mcp.NewToolResultError("name must be a string"), nil
 }

 return mcp.NewToolResultText(fmt.Sprintf("Hello, %s!", name)), nil
}
```

1. 使用 SSE 启用服务，以便能用 HTTP 那样进行 Debug

```bash
> go run ./main.go -t sse
localhost:5568
2025/05/17 11:36:23 Registered tool: hello_world
2025/05/17 11:36:23 SSE server listening on localhost:5568
```

2. 建立 TCP 连接以获取 session id

```bash
> curl http://localhost:5568/sse
event: endpoint
data: /message?sessionId=05dfc3f8-bc10-42d4-b4fb-1c9de72ed444
```

3. initialize 阶段，获取 MCP Server 的基本信息

```bash
> curl -X POST --data '{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "initialize",
  "params": {
    "protocolVersion": "2024-11-05",
    "capabilities": {
      "roots": {
        "listChanged": true
      },
      "sampling": {}
    },
    "clientInfo": {
      "name": "ExampleClient",
      "version": "1.0.0"
    }
  }
}' http://localhost:5568/message?sessionId=05dfc3f8-bc10-42d4-b4fb-1c9de72ed444

event: message
data: {"jsonrpc":"2.0","id":1,"result":{"protocolVersion":"2024-11-05","capabilities":{"tools":{}},"serverInfo":{"name":"Demo 🚀","version":"1.0.0"}}}
```

4. `tool/list` 获取该 MCP Server 能调用的功能（需要传入 session id）

```bash
> curl -X POST --data '{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "tools/list",
  "params":{}
}' http://localhost:5568/message?sessionId=05dfc3f8-bc10-42d4-b4fb-1c9de72ed444

event: message
data: {"jsonrpc":"2.0","id":1,"result":{"tools":[{"annotations":{"destructiveHint":true,"openWorldHint":true},"description":"Say hello to someone","inputSchema":{"properties":{"name":{"description":"Name of the person to greet","type":"string"}},"required":["name"],"type":"object"},"name":"hello_world"}]}}
```

5. `tool/call` 来调用执行

```bash
> curl -X POST --data '{
"jsonrpc": "2.0",
"id": 1,
"method": "tools/call",
"params": {
  "name": "hello_world",
  "arguments": {
    "name": "kyden"
  }
}
}' http://localhost:5568/message?sessionId=05dfc3f8-bc10-42d4-b4fb-1c9de72ed444

event: message
data: {"jsonrpc":"2.0","id":1,"result":{"content":[{"type":"text","text":"Hello, kyden!"}]}}
```

### VSCode Cline 调用

VSCode 配置 Cline MCP Server：`Users/kyden/Library/Application Support/Code/User/globalStorage/saoudrizwan.claude-dev/settings/cline_mcp_settings.json`

```json
{
  "mcpServers": {
    "helloWorld": {
      "command": "/Users/kyden/git-space/mcp.test/mcptest",
      "args": [],
      "env": {}
    }
  }
}
```

Cline 调用相应 MCP Server

{{< figure src="/posts/MCP/vscode-cline.png" title="" >}}

## Reference

- [MCP Introduction](https://modelcontextprotocol.io/introduction)
- [Golang 中的内建 JSON-RPC](https://pkg.go.dev/net/rpc/jsonrpc)
- [Server Inspector](https://github.com/modelcontextprotocol/inspector)


---

> Author: [kyden](https://github.com/kydance)  
> URL: http://localhost:1313/posts/34c4401/  

