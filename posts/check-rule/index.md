# 资格校验接口的微服务设计与实现


{{&lt; admonition type=abstract title=&#34;导语&#34; open=true &gt;}}
在设计和开发微服务接口的过程中，常常会遇到接口职责不够单一、功能混杂的问题。面对这种情况，该如何有效处理呢？
本文以资格校验服务为例，详细介绍如何通过工厂方法、流量镜像和流量回放等技术手段，来解决开发、测试和部署中遇到的此类问题。
{{&lt; /admonition &gt;}}

&lt;!--more--&gt;

## 资格服务的三层架构设计

在实现本接口的过程中，为了应对多种规则校验需求，我们设计了一个基于三层架构的系统。
该设计遵循简洁架构的原则，以确保系统的高内聚和低耦合。
具体划分为以下三层：**控制层（Controller）**、**业务层（Service / Biz）**、**存储层（Store）**。

其中，控制层负责处理外部请求和响应，业务层负责规则的具体逻辑校验，
另外，由于本接口不会存储用户数据，因此存储层的职责由**负责数据的持久化存储** 转变为了**负责外部数据进行交互**。

从架构上看，层与层之间的依赖关系自上而下递进，即控制层依赖业务层，业务层依赖存储层，具体代码架构如下图所示：

{{&lt; figure src=&#34;/posts/check-rule/Outline.svg&#34; title=&#34;&#34; &gt;}}

在各层之间的代码设计上，遵循了严格的依赖倒置原则（DIP）。
具体来说， **控制层（Controller）可以导入业务层（Service）和存储层（Store）** 的包，而非直接与存储层交互。
这样设计的好处是确保业务逻辑独立于数据存储，实现更强的扩展性和维护性。

需要特别注意的是，控制层不应直接导入存储层，除非有非常特殊的需求。
所有涉及存储的操作应通过业务层来完成，从而确保系统设计的层次清晰，职责明确。

### 控制层（Controller）

控制层负责接收并处理来自客户端的请求，
具体操作包括：**解析请求参数、进行参数校验、分发业务逻辑、整合处理结果并返回响应**。
它的主要职责是将请求路由到业务层进行处理，而不直接涉及业务逻辑的实现。

在控制层中，我们通过 `services.Servicer` 接口将请求分发给业务层（Service）。
业务逻辑处理完成后，控制层将结果整合并返回给客户端，从而实现业务路由的功能。

{{&lt; admonition warning &#34;FIXME&#34; ture &gt;}}
框图需要修改，应该严格按照具体操作来画
{{&lt; /admonition &gt;}}

{{&lt; figure src=&#34;/posts/check-rule/Controller.svg&#34; title=&#34;控制层结构示意图&#34; &gt;}}

### 业务层（Biz/Service）

业务层是整个系统的核心，负责处理所有的业务逻辑。
当控制层接收到请求并将其转发至业务层时，业务层将根据具体的业务规则，调用存储层（Store）进行数据的 CURD。

在此层级中，所有的业务逻辑代码应集中于此，确保业务逻辑与其他逻辑（如存储和控制）解耦。
业务层的设计目标是让代码更具扩展性和可维护性。

{{&lt; figure src=&#34;/posts/check-rule/service-store.svg&#34; title=&#34;业务层与存储层的交互&#34; &gt;}}

### 存储层（Store）

存储层是数据交互的入口，它负责与数据库 / 第三方服务进行 CURD 操作。
由于本接口不会存储用户数据，因此存储层的职责由**负责数据的持久化存储** 转变为了**负责外部数据进行交互**，并为上层提供所需的数据。

该层不会涉及任何业务逻辑，而仅专注于数据的存储与转换。

同时，存储层也负责数据的格式转换，例如：

- 将数据库或第三方服务返回的数据格式转换为业务层和控制层能处理的数据格式；
- 将业务层和控制层的数据转换为存储系统或外部服务能够识别的格式。

### 层间交互

在整个系统中，各层之间通过接口进行交互，确保功能的独立性和可扩展性。
层与层之间的通信遵循依赖倒置原则，以便实现模块化和插件化的设计目标，同时大大提高了系统的测试性。

- Controller 依赖于 Service 层：Controller 通过调用 Service 层接口处理业务逻辑，可利用 `golang/mock` 模拟 Service 层进行单元测试。
- Service 依赖于 Store 层：Service 层通过 Store 层接口与存储系统/第三方服务交互，可通过 `golang/mock` 模拟操作。
- Store 依赖于数据库和外部服务：Store 层与数据库或微服务进行直接交互，可以使用 [sqlmock](https://github.com/DATA-DOG/go-sqlmock) 模拟数据库操作，使用 [httpmock](https://github.com/jarcoal/httpmock) 模拟外部 HTTP 请求。

### 资格服务代码设计

在了解了三层架构的基础后，资格服务的代码设计也基于此结构实现。
在具体实现中，我们遵循面向接口编程的原则，以提高代码的扩展性和可测试性。

#### Controller

在 Controller 层中，我们定义了如下的代码结构：

{{&lt; figure src=&#34;/posts/check-rule/class-Controllers.v2.svg&#34; title=&#34;&#34; &gt;}}

它持有 `services.Servicer` 接口，并且实现了 `POST`/`HEAD`/`GET` 等 HTTP 方法，
用于处理 HTTP 请求的响应、请求参数的解析与合法性校验、Service 层业务逻辑的调用执行等操作。

#### Service

在 Service 层中，我们定义了如下的代码结构：

{{&lt; figure src=&#34;/posts/check-rule/class-Services.v2.svg&#34; title=&#34;&#34; &gt;}}

`Server` 接口定义了该服务所支持的功能，实现了接口就是规范的功能。
与此同时，在 `Server` 的实现类（例如，`DjcRuleService`） 中持有 `rule.Ruler` 接口的引用，用于执行资格校验规则。

#### Store

在存储层中，我们采用了工厂方法设计模式，以实现不同规则的动态校验。
`rule.Ruler` 接口定义了核心的 Check 方法，具体的校验规则类（如 `DJCFFriendsRule`, `DJCfmVipRule` 等）通过工厂方法创建，并在 `service.createRuler` 中创建具体规则校验实现类，用以实现具体校验逻辑。

代码结构定义如下：

{{&lt; figure src=&#34;/posts/check-rule/class-Store.v2.svg&#34; title=&#34;&#34; &gt;}}

## 微服务部署

### 服务发现

在上游客户端向某个服务发送请求时，它首先会根据**所请求的服务名称**（例如`check.rule`）在**配置管理中心**（例如`etcd`）查找该服务对应的配置文件：

- 服务配置文件: `/cfg/daoju/.../info/check/rule/check.rule.cfg`
- 环境配置文件: `/cfg/daoju/.../info/check/rule/check.rule.1_1.cfg`
- 部署配置文件: `/cfg/daoju/.../deployment/djc_rule_test.cfg`

需要注意的是，**部署配置文件的文件名是由环境部署文件中的 `deployment[_number].name` 配置项所确定**.

#### 服务配置

服务配置文件中通常包含以下关键消息：

```shell
# /cfg/daoju/.../info/check/rule/check.rule.cfg

[api]
name=check.rule
api=/cgi-bin/daoju/.../rule_check.cgi
timeout=5000
proto=http
method=post

[verify]
key=xxxxxxxxxxxxxxxxxxxx
...
```

其中:

- `api.name` 描述了该服务的名称
- `api.api` 描述了该服务的 URL 路径
- `api.timeout` 描述了该服务的请求超时时间
- `api.proto` 描述了该服务采用的协议格式
- `api.method` 描述了该服务具体采用哪种请求方法
- `verify.key` 描述了请求数据进行校验的密钥

#### 环境配置

环境配置文件中通常包含以下关键消息：

```shell
# /cfg/daoju/.../info/check/rule/check.rule.1_1.cfg

[weight]
total=100
depcnt=2
weight_0=40
weight_1=60

[maintenance]
status=0

[limit]
qps=1000

[deployment]
name=djc_rule_test

[deployment_1]
name=djc_check_rule_go_test
...
```

其中:

- `weight.total` 描述了所有部署环境的权重总和，通常是 100
- `weight.depcnt` 描述了该服务部署在本环境下（测试环境或生产环境）的服务数量
- `weight.weight_&lt;number&gt;` 描述了第 `&lt;number&gt;` 服务的获取请求数据的占比（`weight_&lt;number&gt; / weight.total`）
- `maintenance.status` 描述了当前环境是否已发布 / 正常
- `limit.qps` 描述了该服务所支持的最大 QPS
- `deployment[_number].name` 描述了该服务的环境部署名称

#### 部署配置

部署配置文件中通常包含以下关键消息：

```shell
# /cfg/daoju/.../deployment/djc_rule_test.cfg

[djc_check_rule_go_test]
modid=xxxxxxx:yyyyyyy
mod=xxxx
cmd=xxxxx
domain=
ip_num=1
defaultip_0=&lt;ip&gt;
defaultport_0=&lt;port&gt;

[polaris]
namespace=Development
service=gdp.aaa.bbb.ccc
...
```

其中:

- `djc_check_rule_go_test` 主要用于指向该部署
- `polaris.namespace` 描述了该部署处于何种环境下，例如 `Development` / `Production` / `Test`
- `polaris.service` 描述了该部署所指向的北极星服务地址

一旦通过上面的流程确定了 `polaris.service` 就可以确定 GDP 中的具体代码，大致流程如下：

{{&lt; figure src=&#34;/posts/check-rule/服务注册与发现.svg&#34; title=&#34;&#34; &gt;}}

## 【接口测试】流量回放与镜像

### 流量回放 (Traffic Replay)

流量回放，顾名思义，指的是通过复制线上真实流量（录制），然后在测试环境（或生产环境）进行模拟请求（回放）验证代码逻辑正确性的一种技术方法。
它通过采集线上流量在测试环境（或生产环境）回放逐一对比每个子调用差异和入口调用结果来发现接口代码是否存在问题。

通俗理解，流量回放和使用其他工具（比如 JMeter / postman）构造请求，然后根据返回的响应数据判断测试是否通过的本质相同。
两者的区别在于：流量回放是线上真实流量，而在传统的利用测试工具来发送请求的手段中，人工介入较多。

### 流量镜像 (Traffic Mirror)

流量镜像（Traffic Mirror） ，也称流量影子（Traffic Shadow）, 是一种强大的、无风险的测试应用版本的方法，它将实时流量的副本发送给被镜像的服务。
采用这种方法，可以轻松地测试新版本，而无需在生产环境中部署新版本。

### 验证资格服务的稳定性与准确性

#### 稳定性

由于本次资格服务接口项目属于重构项目，因此在代码重构结束之后，借助了流量回放与镜像技术，验证重构后的代码逻辑是否正确，保证重构后的代码逻辑与重构前的代码逻辑一致。

另外，为了充分验证代码的正确性，我们采用了两种逐层递进的验证方式，
即首先在 CLS 日志系统中抓取线上流量，然后通过流量回放技术手段，在测试环境（或生产环境）进行模拟请求，确保代码逻辑的稳定性，如下图所示；

{{&lt; figure src=&#34;/posts/check-rule/Traffic-replay.svg&#34; title=&#34;&#34; &gt;}}

在进行代码稳定性验证阶段，主要关注以下几个指标：

- 响应时间
- 吞吐量
- 并发用户数
- CPU 使用率
- 内存占用
- 磁盘 I/O
- 网络带宽使用情况
- 出现异常或错误的请求比例

#### 准确性

在充分验证重构代码的**稳定性**之后，接下来就需要验证代码的**准确性**。
这里采用了流量镜像的强大验证技术手段，也就是说，在测试环境（或生产环境）中，直接对线上流量进行镜像，并持续监控镜像流量，确保代码逻辑的准确性。

{{&lt; figure src=&#34;/posts/check-rule/Traffic-Mirror.svg&#34; title=&#34;&#34; &gt;}}

在进行代码准确性性验证阶段，主要关注以下几个指标：

- 功能正确性：
  - 输入、输出是否符合预期
  - 边界条件
  - 异常情况处理
- Online 与 Local 数据是否一致

## Reference

- [beego](https://github.com/beego/beego)
- [sqlmock](https://github.com/DATA-DOG/go-sqlmock)
- [httpmock](https://github.com/jarcoal/httpmock)
- [【道聚城】微服务建设实践总结](https://km.woa.com/articles/show/408014)


---

> : [kyden](https:github.com/kydance)  
> URL: http://kyden.us.kg/posts/check-rule/  

