# Model Context Protocol (MCP): 扩展 AI 能力的新范式


{{< admonition type=abstract title="导语" open=true >}}
在 AI 技术日新月异的今天，Anthropic 提出的 Model Context Protocol (MCP) 正引领着 AI 能力扩展的新浪潮。MCP 就像是为 AI 模型设计的多功能扩展坞，让大语言模型能够像连接了各种外设的超级工作站一样，实时接入专业能力模块，实现能力的弹性扩展。

MCP 通过标准化的协议，将 AI 模型与各类专业工具、知识库和计算资源无缝连接，打破了传统 AI 系统的封闭性。无论是文件内容、数据库记录、API 响应，还是实时系统数据，MCP 都能将其转化为 AI 可理解的上下文，显著提升 AI 在特定领域的表现。

本文将带您深入探索 MCP 的核心概念、工作原理及实际应用，并通过一个完整的 "Hello World" 示例，展示如何快速搭建和调试您的第一个 MCP 服务。无论您是 AI 开发者还是技术爱好者，MCP 都为您打开了一扇通往更智能、更强大 AI 应用的大门。
{{< /admonition >}}

<!--more-->

## 模型上下文协议 MCP

MCP（Model Context Protocol，模型上下文协议），起源于2024年11月25日 [Anthropic](https://www.anthropic.com/) 发布的文章：[Introducing the Model Context Protocol](https://www.anthropic.com/news/model-context-protocol)。这一协议的推出标志着 AI 领域迈向了更加开放、可扩展的新时代。

MCP 在 AI 领域中的作用，就犹如扩展坞在笔记本电脑中的作用一样，为 AI 模型的能力扩展提供了新的可能性。例如，现代扩展坞可以连接显示器、键盘以及移动硬盘等多种硬件外设，使得笔记本电脑瞬间扩展了多种功能；而 MCP Server 将会作为一个智能中枢平台，可以动态接入各类专业能力模块（如知识库、计算工具、领域模型等）。

当 LLM 需要完成特定任务时，它就可以调用这些 MCP Server，实时获得精准的上下文支持，从而实现能力的弹性扩展。这种架构打破了传统 AI 模型的封闭性，让大语言模型像搭载了多功能扩展坞的超级工作站，随时都能获取最合适的专业工具。

### 为什么 MCP 如此重要？

MCP 的出现解决了 AI 领域长期面临的几个关键挑战：

1. **知识时效性问题**：大语言模型的知识在训练后就固定了，而 MCP 允许模型实时获取最新信息
2. **专业能力限制**：通过 MCP，AI 可以调用专门的工具和模型来处理特定领域任务
3. **计算资源优化**：复杂计算可以交给专门的服务处理，而不必全部依赖 LLM 本身
4. **隐私与安全**：敏感数据可以保留在本地或专用服务中，只提供必要的处理结果
5. **可扩展性**：随着新工具和服务的开发，AI 的能力可以不断扩展，无需重新训练

### 关键概念

MCP 协议的核心由以下几个关键组件构成，它们共同定义了 AI 模型如何与外部系统交互：

- **MCP Server**：实现特定功能的服务端，负责响应 AI 模型的请求。可以理解为一个专用的“能力代理”，为 AI 提供特定领域的服务。每个 MCP Server 都可以提供一组独特的能力，如代码分析、数据库查询、图像处理等。

- **MCP Client**：调用 MCP Server 的客户端，通常是集成了 LLM 的应用程序。Client 负责发现可用的 Server，协商能力，并在需要时调用相应的服务。

- **Resources**：资源，指 MCP Server 可以提供给 Client 的各类数据和信息。这些资源可以是文件内容、数据库记录、API 响应、实时系统数据等。每个资源都由唯一的 URI 标识，并且只能被读取，不能被修改。

- **Tools**：工具，是 MCP Server 提供的可执行操作。与 Resources 不同，Tools 可以修改状态或与外部系统交互。每个工具都有唯一的名称、描述和参数模式，使 AI 模型能够正确地调用它们。

- **Prompts**：提示模板，允许 MCP Server 定义特定的提示模式，指导 AI 模型如何处理特定类型的任务。

- **Transport**：传输层，定义了 Client 和 Server 之间的通信方式，如 stdio、HTTP+SSE 等。

这些概念共同构成了 MCP 的基础架构，使 AI 模型能够以标准化的方式与各种外部服务和资源交互。

### JSON-RPC 通信协议

MCP 采用 JSON-RPC 2.0 作为其基础通信协议。JSON-RPC 是一种轻量级的远程过程调用协议，使用 JSON 作为数据格式，非常适合跨语言、跨平台的分布式系统通信。

#### 为什么选择 JSON-RPC？

MCP 选择 JSON-RPC 而非其他流行的 RPC 协议（如 gRPC）有以下几个关键原因：

1. **灵活性与简易性**：JSON-RPC 使用简单的 JSON 格式，更易于实现和调试
2. **广泛的语言支持**：几乎所有编程语言都有完善的 JSON 处理库
3. **动态类型**：适合 AI 系统的动态特性，可以在运行时发现和调用服务
4. **低门槛**：不需要复杂的工具链和代码生成步骤

#### JSON-RPC 与 gRPC 的主要区别

| 特性 | JSON-RPC | gRPC |
|---------|----------|------|
| 数据格式 | JSON（文本型，人类可读） | Protocol Buffers（二进制，高效） |
| 类型系统 | 动态弱类型 | 静态强类型 |
| 接口定义 | 无需预先定义，可动态发现 | 需要 .proto 文件预先定义 |
| 代码生成 | 无需生成客户端代码 | 需要生成客户端和服务端代码 |
| 传输协议 | 不特定（通常是 HTTP） | HTTP/2 |
| 流式支持 | 有限支持 | 原生支持双向流 |
| 性能 | 中等 | 高 |
| 实现复杂度 | 低 | 中到高 |

MCP 利用 JSON-RPC 的灵活性，并通过生命周期管理（lifecycle management）来动态发现服务能力，这一点是 MCP 协议的扩展，而非 JSON-RPC 本身的功能。

#### JSON-RPC 的请求与响应格式

##### 请求（Request）格式

与传统 HTTP API 不同，JSON-RPC 不在 URL 中指定方法，而是将方法名称和参数都放在 JSON 请求体中：

```json
{
    "jsonrpc": "2.0",
    "id": "uuid-or-number",
    "method": "method-name",
    "params": {
        "paramName1": "value1",
        "paramName2": "value2"
    }
}
```

其中：

- `jsonrpc`: 固定值 "2.0"，指定协议版本
- `id`: 请求标识符，可以是字符串或数字，用于关联请求和响应
- `method`: 要调用的方法名称
- `params`: 可选参数，可以是对象或数组

##### 响应（Response）格式

JSON-RPC 响应会包含原始请求的 `id`，以允许客户端关联请求和响应（特别是在批量请求的情况下）：

```json
{
    "jsonrpc": "2.0",
    "id": "uuid-or-number",
    "result": {
        "key1": "value1",
        "key2": "value2"
    }
}
```

如果发生错误，响应将包含 `error` 字段而非 `result`：

```json
{
    "jsonrpc": "2.0",
    "id": "uuid-or-number",
    "error": {
        "code": -32601,
        "message": "Method not found",
        "data": {
            "additionalInfo": "Optional error details"
        }
    }
}
```

#### 批量请求处理

JSON-RPC 2.0 支持批量请求，允许客户端在单个 HTTP 请求中发送多个 RPC 调用，这在需要执行多个相关操作时非常有用：

```json
// 批量请求示例
[
  {"jsonrpc": "2.0", "method": "sum", "params": [1,2], "id": 1},
  {"jsonrpc": "2.0", "method": "subtract", "params": [42,23], "id": 2},
  {"jsonrpc": "2.0", "method": "foo.get", "params": {"name": "myself"}, "id": 3}
]

// 相应的批量响应
[
  {"jsonrpc": "2.0", "result": 3, "id": 1},
  {"jsonrpc": "2.0", "error": {"code": -32601, "message": "Method not found"}, "id": 2},
  {"jsonrpc": "2.0", "result": {"firstName": "John", "lastName": "Doe", "age": 30}, "id": 3}
]
```

值得注意的是，与 gRPC 使用流式处理批量请求不同，JSON-RPC 的批量请求是通过将多个请求打包到一个数组中实现的。

### Tools 功能

Tools（工具）是 MCP 协议中最强大的组件之一，它允许 AI 模型执行各种操作和交互。与只能被读取的 Resources 不同，Tools 可以修改状态或与外部系统交互，这使得 AI 模型能够执行复杂的任务。

#### 工具的特点

Tools 具有以下几个关键特点：

1. **唯一标识**：每个工具都有一个唯一的名称，用于在 MCP Server 中识别

2. **描述性元数据**：工具包含清晰的描述，使 AI 模型能够理解其用途和功能

3. **结构化参数**：工具定义了所需的输入参数及其类型，使 AI 可以正确地调用它们

4. **状态修改能力**：工具可以执行操作并修改系统状态，如创建文件、发送电子邮件或查询数据库

#### 工具的定义格式

在 MCP 中，工具的定义通常包含以下元数据：

```json
{
  "name": "tool_name",
  "description": "What this tool does and when to use it",
  "inputSchema": {
    "type": "object",
    "properties": {
      "param1": {
        "type": "string",
        "description": "Description of parameter 1"
      },
      "param2": {
        "type": "number",
        "description": "Description of parameter 2"
      }
    },
    "required": ["param1"]
  },
  "annotations": {
    "destructiveHint": true,
    "openWorldHint": false
  }
}
```

其中：

- `name`：工具的唯一标识符
- `description`：工具的用途和功能描述
- `inputSchema`：使用 JSON Schema 定义的输入参数结构
- `annotations`：可选的元数据，提供关于工具行为的额外信息

#### 工具调用流程

当 AI 模型需要使用特定工具时，调用流程通常如下：

1. AI 模型首先通过 `tools/list` 方法获取可用工具列表

2. 根据工具描述和当前任务，AI 模型选择适合的工具

3. AI 模型通过 `tools/call` 方法调用工具，并提供必要的参数

4. MCP Server 执行工具操作并返回结果

5. AI 模型处理返回的结果并继续对话

#### 工具的应用场景

MCP 工具可以应用于各种场景，例如：

- **数据检索与分析**：查询数据库、分析数据集、生成报表

- **内容生成**：创建文档、生成图表、编辑媒体文件

- **系统集成**：与外部 API 交互、发送通知、执行工作流程

- **专业计算**：运行复杂算法、执行模拟、进行科学计算

### 资源能力 (Resource Capability)

Resources 是 MCP 协议中的另一个关键组件，它为 AI 模型提供了访问各种数据和信息的能力。与 Tools 不同，Resources 仅提供只读访问，不能修改系统状态。

#### 资源的核心特点

资源具有以下几个关键特点：

1. **唯一标识**：每个资源都由唯一的 URI 标识，便于引用和访问

2. **数据多样性**：资源可以包含文本或二进制数据，支持各种类型的信息

3. **只读特性**：资源只能被读取，不能被修改，确保了数据的安全性

4. **上下文增强**：资源为 AI 模型提供了额外的上下文信息，提高了应答的精准性

#### 常见的资源类型

MCP 服务器可以提供各种类型的资源，包括：

- **文件内容**：文档、代码、配置文件等

- **数据库记录**：数据表、查询结果、统计信息

- **API 响应**：来自外部服务的 JSON、XML 等格式的响应

- **实时系统数据**：监控指标、日志、状态信息

- **多媒体资源**：图像、音频、视频的元数据或内容

#### 资源与工具的区别

Resources 和 Tools 之间的主要区别在于：

| 特性 | Resources | Tools |
|---------|-----------|-------|
| 访问权限 | 只读 (ReadOnly) | 读写 (Read/Write) |
| 状态修改 | 不能修改状态 | 可以修改系统状态 |
| 主要用途 | 提供上下文信息 | 执行操作和交互 |
| 标识方式 | URI | 名称 |
| 安全级别 | 较高 | 需要额外权限控制 |

> 可以将 Resources 理解为“知识源”，提供各种数据库、表格和文件的访问；而 Tools 则是“操作工具”，允许执行 CRUD 操作、计算和其他交互。

#### 资源的访问方法

MCP 协议定义了以下几个方法来访问资源：

#### 资源访问方法

MCP 协议定义了以下几个方法来访问和管理资源：

- `resource/list`：列出可用的资源，包括其 URI 和元数据

- `resource/read`：读取特定资源的内容

- `resource/subscribe`：订阅资源的变化通知

- `notifications/resource/update`：接收资源更新的通知

下图展示了 MCP 资源的基本工作流程：

{{< figure src="/posts/MCP/mcp-resource.png" title="MCP 资源工作流程" >}}

#### 资源的应用场景

资源功能在多种场景中非常有用，例如：

- **知识库集成**：允许 AI 模型访问最新的专业知识和文档

- **实时数据分析**：提供实时系统指标和数据流

- **上下文感知应用**：使 AI 模型能够基于当前环境和用户状态提供定制化响应

- **安全数据访问**：允许受控的访问敏感数据，而不需要将其完全暴露给 AI 模型

### 提示能力 (Prompt Capability)

Prompt 能力是 MCP 协议中的另一个关键组件，它允许 MCP Server 定义特定的提示模板，指导 AI 模型如何处理特定类型的任务或生成特定格式的输出。

#### 提示的作用与价值

Prompt 能力主要解决以下几个问题：

1. **标准化输出**：确保 AI 模型生成的响应符合特定的格式和结构

2. **专业领域指导**：为 AI 模型提供特定领域的指导和最佳实践

3. **一致性保证**：确保不同的 AI 交互保持一致的语调和响应风格

4. **工作流程优化**：定义结构化的对话流程，提高交互效率

#### 提示的元数据结构

在 MCP 中，提示的元数据定义如下：

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

其中：

- `name`：在同一个 Server 内必须唯一的提示标识符
- `description`：提示的用途和功能描述
- `arguments`：可选的参数列表，每个参数都有自己的名称、描述和是否必需

#### 提示的访问方法

与 Tools 和 Resources 类似，MCP 协议定义了以下方法来访问提示：

- `prompts/list`：获取所有可用的提示及其元数据

- `prompts/get`：获取指定提示的模板内容

#### 提示中的角色定义

在 Prompt 模板中，可以定义两种主要角色（Role）来指导 AI 模型的行为：

- **User**：代表与系统交互的用户，提供问题和上下文信息

- **Assistant**：代表 AI 模型本身，负责处理用户请求并生成响应

这种角色定义允许 MCP Server 创建结构化的对话流，指导 AI 模型如何处理特定类型的交互。

#### 提示的应用场景

Prompt 能力在多种场景中非常有用，例如：

- **结构化数据生成**：创建符合特定格式的 JSON、XML 或其他数据结构

- **专业领域应用**：为医疗、法律、金融等领域提供特定的指导

- **多语言支持**：定义不同语言的响应模板

- **响应一致性**：确保在不同交互中保持一致的响应风格和语调

### MCP 传输层 (Transport)

MCP 协议的传输层定义了 Client 和 Server 之间的通信方式。MCP 设计了多种灵活的传输方式，以适应不同的部署场景和应用需求。

#### 传输层的作用

传输层在 MCP 中承担以下重要职责：

1. **建立连接**：在 Client 和 Server 之间创建可靠的通信通道

2. **消息传递**：确保 JSON-RPC 消息的可靠传递

3. **会话管理**：维护 Client 和 Server 之间的会话状态

4. **错误处理**：处理通信过程中的异常情况

#### 支持的传输方式

MCP 当前支持以下几种主要的传输方式：

##### 1. stdio 标准输入输出

stdio 是最简单的传输方式，利用标准输入和输出流进行通信。这种方式特别适合本地集成和命令行工具。

**主要特点：**

- **简单性**：无需网络设置或端口配置
- **低开销**：最小化的集成要求
- **本地性**：适合本地运行的应用
- **安全性**：无需网络暴露

{{< figure src="/posts/MCP/server-stdio.png" title="MCP stdio 传输模式" >}}

##### 2. HTTP with SSE

HTTP with SSE（Server-Sent Events）是一种基于 Web 的传输方式，允许通过 HTTP 协议进行双向通信，并使用 SSE 实现服务器到客户端的实时推送。

**主要特点：**

- **Web 兼容性**：可以在浏览器和 Web 环境中使用
- **分布式部署**：支持跨网络的分布式部署
- **实时推送**：服务器可以实时推送消息到客户端
- **扩展性**：可以集成到现有的 Web 服务中

{{< figure src="/posts/MCP/server-sse.png" title="MCP HTTP with SSE 传输模式" >}}

##### 3. Streamable HTTP

Streamable HTTP 是一种正在开发中的传输方式，旨在提供更高效的流式数据传输。这种方式特别适合需要处理大量数据或实时流的场景。

**预期特点：**

- **高效流式处理**：优化的数据流传输
- **低延迟**：减少端到端的响应时间
- **双向流**：支持客户端和服务器之间的双向数据流

> 注意：Streamable HTTP 传输方式目前仍在开发中，规范可能会发生变化。

#### 选择合适的传输方式

选择哪种传输方式取决于你的应用需求：

| 传输方式 | 适用场景 | 优势 | 劣势 |
|------------|------------|------|------|
| stdio | 本地应用、命令行工具 | 简单、快速集成 | 仅限本地使用 |
| HTTP with SSE | Web 应用、分布式系统 | 广泛兼容性、易于部署 | 需要网络配置 |
| Streamable HTTP | 大数据处理、实时流应用 | 高效、低延迟 | 仍在开发中 |

### MCP 生命周期 (Life Cycle)

MCP 协议定义了一个清晰的客户端和服务器之间的通信生命周期。这个生命周期管理确保了连接的可靠性和双方能力的正确协商。

#### 生命周期的三个阶段

MCP 生命周期包含以下三个清晰定义的阶段：

##### 1. 初始化阶段 (Initialize)

初始化阶段是 MCP 连接的第一步，也是最关键的一步。在这个阶段，客户端和服务器会进行以下操作：

- **协议版本协商**：确定双方支持的 MCP 协议版本

- **能力发现**：服务器向客户端公布其支持的能力（Tools、Resources、Prompts 等）

- **参数协商**：协商连接参数和限制

- **身份信息交换**：客户端和服务器交换基本信息（名称、版本等）

这个阶段使用 `initialize` 方法，客户端发送请求，服务器响应其能力。

##### 2. 操作阶段 (Operation)

初始化成功后，连接进入操作阶段。在这个阶段：

- **功能调用**：客户端可以调用服务器提供的各种功能（如 tools/call、resource/read 等）

- **状态维护**：服务器维护会话状态和上下文

- **双向通信**：客户端和服务器可以进行双向通信，服务器可以发送通知

- **错误处理**：处理操作过程中的各种错误情况

这个阶段是 MCP 连接的主要工作阶段，可能持续较长时间。

##### 3. 关闭阶段 (Shutdown)

当客户端或服务器需要终止连接时，进入关闭阶段：

- **优雅终止**：客户端发送 `shutdown` 请求，通知服务器即将终止连接

- **资源释放**：服务器释放与该连接相关的资源

- **状态清理**：清理会话状态和临时数据

关闭阶段确保了连接的正常终止，防止资源泄漏和状态不一致。

#### 生命周期的重要性

清晰定义的生命周期对 MCP 协议有以下重要作用：

1. **可靠性**：确保连接的建立、维护和终止都按照预期的方式进行

2. **兼容性**：通过协议版本协商，确保不同版本的客户端和服务器可以兼容工作

3. **动态发现**：允许客户端动态发现服务器的能力，而不需要预先编码

4. **资源管理**：确保资源在不再需要时被正确释放

下图展示了 MCP 生命周期的完整流程：

{{< figure src="/posts/MCP/mcp-lifecycle.png" title="MCP 生命周期流程" >}}

## Server Inspector 调试工具

MCP 生态系统的一个重要组成部分是开发者工具，其中 Server Inspector 是一个专为 MCP 服务器设计的交互式调试工具。它类似于 API 开发中的 Postman 或 gRPC 开发中的 BloomRPC，但专门针对 MCP 协议进行了优化。

### 为什么需要 Server Inspector

Server Inspector 解决了 MCP 开发中的几个关键问题：

1. **简化测试流程**：无需编写客户端代码就能测试 MCP 服务器

2. **可视化调试**：提供图形界面查看请求和响应

3. **快速原型制作**：快速验证 MCP 服务器的功能

4. **多传输方式支持**：支持 stdio 和 HTTP+SSE 传输方式

### 安装与使用

Server Inspector 的设计非常灵活，可以直接通过 npx 运行，无需安装：

```bash
# 基本用法
npx @modelcontextprotocol/inspector <command>

# 传递命令行参数
npx @modelcontextprotocol/inspector <command> <arg1> <arg2>

# 传递环境变量
npx @modelcontextprotocol/inspector -e KEY=value -e KEY2=$VALUE2 <command> <arg1> <arg2>
```

其中：

- `@modelcontextprotocol/inspector` 是包的完整名称，其中 `@modelcontextprotocol` 是组织作用域
- `<command>` 是要执行的 MCP 服务器命令，可以是可执行文件或脚本
- `<arg...>` 是传递给 MCP 服务器的命令行参数
- `-e KEY=VALUE` 参数允许设置环境变量

### 架构与组件

当启动 Server Inspector 时，它会创建两个主要组件：

1. **Web 界面**：一个基于浏览器的图形化界面，默认运行在端口 5173。该界面允许你：
   - 查看服务器能力（Tools、Resources 等）
   - 发送请求并查看响应
   - 查看详细的日志和调试信息
   - 保存和加载请求模板

2. **代理服务器**：一个运行在端口 3000 的代理服务，负责：
   - 与目标 MCP 服务器建立连接
   - 转发请求和响应
   - 提供日志和调试信息

### 使用场景

Server Inspector 在以下场景中特别有用：

1. **MCP 服务器开发**：在开发新的 MCP 服务器时进行快速测试

2. **API 调试**：调试现有 MCP 服务器的 API 和功能

3. **原型制作**：快速创建和测试 MCP 功能原型

4. **文档生成**：探索服务器能力并生成文档

### 高级功能

Server Inspector 还提供了一些高级功能：

- **会话管理**：维护多个 MCP 会话
- **请求历史**：记录和重放之前的请求
- **导出功能**：将请求和响应导出为代码或文档
- **批量测试**：执行批量请求进行压力测试

### 与其他工具的集成

Server Inspector 可以与其他开发工具集成：

- **CI/CD 系统**：在自动化测试中使用
- **IDE 插件**：与 VSCode 等 IDE 集成
- **文档系统**：生成 API 文档

要了解更多信息，可以访问 [Server Inspector 官方仓库](https://github.com/modelcontextprotocol/inspector)。

## MCP Server 中的 "Hello World" 示例

为了更好地理解 MCP 协议的实际应用，我们来创建一个简单的 "Hello World" 示例。这个示例将展示如何使用 Go 语言创建一个基本的 MCP 服务器，并实现一个简单的问候工具。

### 实现目标

我们将创建一个具有以下功能的 MCP 服务器：

1. 支持多种传输方式（stdio 和 SSE）
2. 提供一个简单的 `hello_world` 工具，可以问候用户
3. 实现完整的 MCP 生命周期

### 代码实现

以下是使用 Go 语言实现的完整代码：

```go
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
	// 定义命令行参数，支持多种传输方式
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

	// 启动服务器
	if err := run(transport, *addr); err != nil {
		panic(err)
	}
}

func run(transport, addr string) error {
 // Create MCP server with explicit options
	s := server.NewMCPServer(
		"Demo 🚀",  // 服务器名称
		"1.0.0",         // 版本号
	)

 // Add tool with more explicit configuration
	tool := mcp.NewTool("hello_world",
		mcp.WithDescription("Say hello to someone"),
		mcp.WithString("name",  // 定义字符串参数
			mcp.Required(),      // 设置为必需参数
			mcp.Description("Name of the person to greet"),  // 参数描述
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
	// 从请求中提取参数
	name, ok := request.Params.Arguments["name"].(string)
	if !ok {
		// 参数错误时返回错误
		return mcp.NewToolResultError("name must be a string"), nil
	}

	// 返回文本响应
	return mcp.NewToolResultText(fmt.Sprintf("Hello, %s!", name)), nil
}
```

### 代码解析

这个示例包含以下几个关键部分：

1. **服务器创建**：使用 `server.NewMCPServer()` 创建 MCP 服务器实例

2. **工具定义**：使用 `mcp.NewTool()` 定义工具，包括名称、描述和参数

3. **工具处理函数**：实现 `helloHandler` 函数处理工具调用

4. **多传输方式支持**：支持 stdio 和 HTTP+SSE 两种传输方式

### 运行与测试

让我们来测试这个 MCP 服务器。我们将使用 SSE 传输方式，以便于通过 HTTP 进行测试。

#### 1. 启动服务器

首先，我们使用 SSE 模式启动服务器：

```bash
> go run ./main.go -t sse
localhost:5568
2025/05/17 11:36:23 Registered tool: hello_world
2025/05/17 11:36:23 SSE server listening on localhost:5568
```

#### 2. 获取会话 ID

连接到 SSE 端点以获取会话 ID：

```bash
> curl http://localhost:5568/sse
event: endpoint
data: /message?sessionId=05dfc3f8-bc10-42d4-b4fb-1c9de72ed444
```

这个响应提供了一个会话 ID，我们将在后续请求中使用它。

#### 3. 初始化连接

接下来，我们发送 `initialize` 请求来启动 MCP 生命周期：

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

服务器响应了其协议版本、能力和基本信息。

#### 4. 获取可用工具

现在我们可以使用 `tools/list` 方法获取服务器提供的工具列表：

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

响应显示服务器提供了一个名为 `hello_world` 的工具，包含其描述和参数模式。

#### 5. 调用工具

最后，我们可以使用 `tools/call` 方法调用 `hello_world` 工具：

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

服务器成功响应并返回了文本“Hello, kyden!”。

### 关键要点

这个简单的示例演示了 MCP 协议的几个关键特性：

1. **多传输方式**：支持不同的通信方式（stdio 和 SSE）

2. **生命周期管理**：包含初始化、操作和关闭阶段

3. **能力发现**：客户端可以动态发现服务器能力

4. **结构化参数**：工具定义了清晰的参数模式

5. **JSON-RPC 通信**：使用标准的 JSON-RPC 2.0 协议进行通信

### 与 VSCode Cline 集成

MCP 协议的一个重要应用场景是与开发环境的集成。Claude 已经在 VSCode 中通过 Cline 插件提供了 MCP 支持，让开发者可以直接在 IDE 中使用自定义的 MCP 服务器。

#### 配置方法

VSCode Cline 插件的 MCP 服务器配置存储在以下路径：

```
Users/kyden/Library/Application Support/Code/User/globalStorage/saoudrizwan.claude-dev/settings/cline_mcp_settings.json
```

配置文件的格式如下：

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

其中：

- `helloWorld`：服务器的唯一标识符，将在 VSCode 中显示
- `command`：服务器可执行文件的绝对路径
- `args`：传递给服务器的命令行参数
- `env`：服务器运行时的环境变量

#### 使用方法

配置完成后，可以在 VSCode 中通过 Cline 插件直接调用配置的 MCP 服务器：

1. 打开 VSCode 的 Cline 面板
2. 在对话中使用 `@helloWorld` 语法调用服务器
3. Claude 将自动连接并使用该 MCP 服务器的功能

下图展示了 VSCode Cline 中调用 MCP 服务器的示例：

{{< figure src="/posts/MCP/vscode-cline.png" title="VSCode Cline 中的 MCP 服务器调用" >}}

#### 应用场景

在 VSCode 中集成 MCP 服务器有以下几个主要优势：

1. **开发者工具集成**：直接在 IDE 中使用自定义工具

2. **上下文感知**：让 AI 能够访问当前项目的上下文和状态

3. **专业能力增强**：为 AI 提供项目特定的工具和知识

4. **工作流程优化**：简化开发者与 AI 的协作流程

## 总结与展望

Model Context Protocol (MCP) 作为一种新兴的 AI 能力扩展协议，正在开启 AI 应用的新纪元。通过本文的探讨，我们可以看到 MCP 在以下几个方面的关键价值：

### MCP 的价值

1. **打破能力封闭**：通过标准化的接口，允许 AI 模型与各种外部服务和资源交互

2. **解决知识时效性**：让 AI 模型能够访问最新的信息和数据

3. **增强专业能力**：通过专用工具和资源，提升 AI 在特定领域的能力

4. **优化资源利用**：将复杂计算和处理任务分流到专用服务

5. **保障隐私与安全**：提供受控的数据访问机制

### 未来展望

MCP 协议仍处于早期发展阶段，但其潜力已经显现。随着生态系统的成熟，我们可以期待以下发展：

1. **更多标准化工具**：各领域将开发专用的 MCP 服务器和工具

2. **深度集成**：与开发环境、生产系统和云服务的更深入集成

3. **多模型协作**：促进不同 AI 模型之间的协作和专业能力共享

4. **社区生态系统**：形成开源的 MCP 工具和服务市场

随着 MCP 协议的发展，我们有理由相信，它将成为推动 AI 能力扩展和集成的关键标准，为构建更智能、更开放的 AI 生态系统铺平道路。

## 参考资料

- [MCP 官方介绍](https://modelcontextprotocol.io/introduction)
- [MCP 规范文档](https://modelcontextprotocol.io/specification)
- [Anthropic 关于 MCP 的公告](https://www.anthropic.com/news/model-context-protocol)
- [Golang 中的内建 JSON-RPC](https://pkg.go.dev/net/rpc/jsonrpc)
- [Server Inspector 工具](https://github.com/modelcontextprotocol/inspector)
- [MCP Go 客户端库](https://github.com/mark3labs/mcp-go)


---

> Author: [kyden](https://github.com/kydance)  
> URL: http://kydance.github.io/posts/model-context-protocol/  

