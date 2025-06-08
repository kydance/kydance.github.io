# Prometheus


{{< admonition type=abstract title="导语" open=true >}}
**Prometheus已成为云原生环境中事实上的标准监控解决方案。本文将带您全面了解Prometheus的核心概念、架构设计、数据模型及实践应用，帮助您构建高效可靠的监控系统。**
{{< /admonition >}}

<!--more-->

## 简介

Prometheus 是目前主流的一个开源监控系统和告警工具包，它可以与 Kubernetes 等现代基础设施平台配合，轻松集成到云原生环境中，提供对容器华应用、微服务架构等的全面监控。

Prometheus 受启发于 Google 的 Brogmon 监控系统，从2012年开始由前Google工程师在 Soundcloud 以开源软件的形式进行研发，并且于2015年早期对外发布早期版本。

2016年5月继 Kubernetes 之后成为第二个正式加入CNCF基金会的项目，同年6月正式发布1.0版本。

2017年底发布了基于全新存储层的2.0版本，能更好地与容器平台、云平台配合。

Prometheus 于2016年加入云计算基金会，成为继 Kubernetes 之后的第二个托管项目。

Prometheus 收集并存储其指标作为时间序列数据，即指标信息与其记录的时间戳一起存储，同时存储的还有可选的称为标签的键值对。

---

### 特性

Prometheus 的主要特性包括:

- 多位数据模型：包含 Metric 名称和键值对标识的时间序列数据
- PromQL：一种可以灵活利用上述维度数据的查询语言
- 不依赖于分布式存储：单个服务器节点是自治的
- Pull 模式：通过基于 HTTP 的拉模式（Pull）进行时间序列数据收集
- 可以通过一个中间网关（Push Gateway）以推模式上报时间序列数据
- 通过服务发现或静态配置发现监控目标
- 支持多种模式的图表和仪表盘

---

### 何为 Metric ?

Metric 就是用数字来测量/度量。时间序列，指的是记录一段时间内的变化。

例如，
对于 web 服务器，可以测量请求耗时；
对于数据库，可以测量活动连接数、活动查询数等。

---

### 组件

Prometheus 的生态系统由多个组成部分构成，其中许多是可选的：

- [Prometheus Server](https://github.com/prometheus/prometheus): 用于抓取和存储时间序列数据
- [clientlibs](https://prometheus.io/docs/instrumenting/clientlibs/): 用于检测应用程序代码的客户端库
- [Push Gateway](https://github.com/prometheus/pushgateway): 用于短时任务的推送接收器
- [exporters](https://prometheus.io/docs/instrumenting/exporters/):
用于 HAProxy、StatsD、Graphite 等服务的专用输出程序
- [alertmanager](https://github.com/prometheus/alertmanager): 处理告警的警报管理器
- 各种支持工具

**大多数的 Prometheus 组件都是使用 Go 语言开发的 => 很容易构建、部署为静态二进制文件。**

---

### 架构

Prometheus 可直接或间接通过推送网关（Push Gateway）抓取监控指标（**适用于短时任务**），它在本地存储所有抓取到的样本数据，并在这些数据上执行一系列规则，以从现有数据中汇总并记录新的时间序列或生成告警；

可以使用 Grafana 或其他 API 消费对收集到的数据进行可视化展示。

---

## 数据模型

Prometheus 将所有数据都存储为[时间序列](https://en.wikipedia.org/wiki/Time_series)：
属于同一指标（metric）和同一组标注维度（label）的带时间戳的值流。

除了存储的时间序列外，prometheus 还可以根据查询结果生成临时派生时间序列。

```text
^
│   . . . . . . . . . . . . . . . . .   . .   go_gc_duration_seconds_count 12
│     . . . . . . . . . . . . . . . . . . .   go_goroutines 32
│     . . . . . . . . . .   . . . . . . . .   go_info{version="go1.22.3"} 1
│     . . . . . . . . . . . . . . . .   . .  
v
  <------------------ 时间 ---------------->
```

在时间序列中的每一个点称为一个样本（sample），样本由以下三部分组成：

- 指标 Metric：Metric name 和描述当前样本特征的 label sets
- 时间戳 Timestamp：一个精确到毫秒的时间戳
- 样本值 Value：一个 float64 的浮点型数据表示当前的样本值

---

### Metric name

每个时间序列都由其指标名称和称为标签的可选键值对唯一标识。

- 指定要测量的系统的一般功能，例如 `http_request_total`
- 指标名称可以包含 ASCII 字符、数字、下划线和冒号，必须匹配正则表达式 `[a-zA-Z_:][a-zA-Z0-9_:]*`

> 冒号是为用户定义的录制规则保留的，exporter或直接仪器不应使用它们

---

### Matric labels

- 使 Prometheus 的维度数据模型能够识别同一指标名称的任何标签组合：
它标识了该度量的特定维度示例化，(例如，所有发送 POST 到 `/api/tracks` 的 HTTP 请求)。Prometheus 查询语言允许基于这些维度进行筛选和聚合；
- 任何标签值的更改，包括添加、删除标签，都将创建一个新的时间序列；
- label 可以包含 ASCII 字符、数字以及下划线，必须匹配 `[a-zA-Z_][a-zA-Z0-9_]*`
- 以 `__` (两个下划线) 开头的标签是 Prometheus 内部使用的
- 标签值可以包含任何 Unicode 字符。
- 标签值为空的标签被视为等同于不存在的标签。

更多内容可参考 [Metric name 和 Label 的最佳实践](https://prometheus.io/docs/practices/naming/)。

---

### Sample

Sample，样本构成实际的时间序列数据，其中每个样本包括：

- 一个 float64 值
- 毫秒精度的时间戳

{{< admonition type=note title="原生直方图支持" open=true >}}
从 Prometheus v2.40 开始，实验性地支持原生直方图（histograms）。采样值不再是简单的 float64，而是一个完整的直方图。
{{< /admonition >}}

---

### Notation (表达式)

给定一个指标名称和一组标签，时间序列通常使用以下符合进行表示：`<metric name>{<label name>=<label value>, ...}`

例如，指标名称 `api_http_request_total`，带有 `method="POST"` 和 `handler="/messages"` label 的时间序列可以写为：
`api_http_request_total{method="POST", handler="/messages"}`

与 [OpenTSDB](http://opentsdb.net/) 使用的表示方式相同。

---

## Metric 类型

Prometheus 客户端提供四种核心指标类型。
这些类型目前仅在客户端库（以便根据特定类型的使用情况定制应用程序接口）和传输协议中有所区别。

目前，Prometheus 服务端还没有使用类型信息，而是将所有数据平铺为无类型的时间序列。  *未来版本可能会有所改变*

---

### Counter (计数器)

**`Counter` 是一种累积度量，表示单个递增的计数器（只增不减），其值只能在重新启动时增加或重置为零**。
例如，可以使用 `Counter` 统计已服务的请求数、已完成的任务数、错误数等。

**不能使用 `Counter` 记录可能减小的值**。

---

### Gauge (仪表盘)

**`Gauge` 是一种度量标准，代表一个可以任意升降的单一数值**。

Gauge 通常用于测量温度、当前内存使用量、当前连接数等，也可用于上下变化的 "计数"，如并发请求的数量.

---

### Histogram (直方图)

**`Histogram` 对观测结果进行采样（通常时请求耗时、响应体大小），并按可配置的桶进行计数，还提供了所有观察值的总和**。

基本度量名称为 `<basename>` 的 `Histogram` 会在抓取过程中暴露多个时间序列：

- 观察桶的累积计数器，对外展示为：`<basename>_bucket{le="<upper inclusive bound>"}`
- 所有观测值的总和，对外展示为`<basename>_sum`
- 已观察到的事件数，对外展示为`<basename>_count（与上面的<basename>_bucket{le="+Inf"}相同）`

使用 `histogram_quantile()` 函数可以根据 histogram 甚至 histogram 的聚合计算分位数。
Histogram 也适用于计算 Apdex 得分。在对bucket进行操作时，请记住 histogram 是累积的。

{{< admonition type=note title="直方图" open=true >}}
从普罗米修斯v2.40开始，就有对原生直方图的实验支持。原生直方图只需要一个时间序列，除了观测值的总和和计数外，还包括动态数量的桶。原生直方图允许以很小的成本获得更高的分辨率。一旦本机直方图接近成为一个稳定的功能，详细的文档将随之而来。
{{< /admonition >}}

---

### Summary (摘要)

与 histogram 类似，summary对观察结果（通常是请求耗时和响应体大小）进行采样。
虽然它还提供了观测的总计数和所有观测值的总和，但它在滑动时间窗口内计算可配置的分位数。

基本度量名称为 `<basename>` 的 summary 会在抓取过程中暴露多个时间序列：

- 观测事件的流式 `φ-quantiles(0 ≤ φ ≤ 1)` 分位数，对外展示为`<basename>{quantile="<φ>"}`
- 所有观测值的总和，对外展示为`<basename>_sum`
- 已观察到的事件数，对外展示为`<basename>_count`

关于 histogram 和 summary 的区别，可以简单概括为 histogram分桶记录数据，后续可在服务端使用表达式函数进行各种计算；而summary在客户端上报时就按配置上报计算好的φ-分位数。

  1. 如果需要多个实例的数据进行汇总，请选择 histogram。
  2. 除此以外，如果对将要观察的值的范围和分布有所了解，请选择 histogram。无论值的范围和分布如何，如果需要准确的分位数，请选择 summary。

更多内容请查看HISTOGRAMS AND SUMMARIES。

---

## job 和 instance

使用 Prometheus 的术语来说，一个可以抓取的端点被称为一个 `instance`，通常对应于一个进程。

具有相同功能的实例集合（例如，为提高可扩展性/可靠性而创建的副本进程）称为 job。

例如，具有四个副本 instance 的api-server job：

- job: api-server
  - instance 1: 1.2.3.4:5670
  - instance 2: 1.2.3.4:5671
  - instance 3: 5.6.7.8:5670
  - instance 4: 5.6.7.8:5671

当 Prometheus 抓取目标时，它会自动在抓取的时间序列上附加下面的标签，用于区分不同的目标：

- job：目标所属的已配置作业名
- instance：被抓取的目标URL的`<host>:<port>` 部分。

对于每一次抓取，prometheus 都会按照以下时间序列存储一个样本：

- `up{job="<job-name>", instance="<instance-id>"}`: 如果实例是健康的，即可访问的，就是1或者如果抓取失败，则为0。
- `scrape_duration_seconds{job="<job-name>", instance="<instance-id>"}`
- `scrape_samples_post_metric_relabeling{job="<job-name>", instance="<instance-id>"}`
- `scrape_samples_scraped{job="<job-name>", instance="<instance-id>"}`
- `scrape_series_added{job="<job-name>", instance="<instance-id>"}`

---

## Usage

Prometheus 支持预编译二进制文件安装、源码安装、docker等方式。

```bash
go get github.com/prometheus/prometheus
go get github.com/prometheus/client_golang/prometheus/promhttp
```

### Docker 使用示例

1. 使用 Docker 拉取、开启一个实例: 启动 Prometheus Server

    ```bash
    $docker run -d --name prometheus \
        -p 9090:9090 \
        prom/prometheus
    ```

2. 通过浏览器访问 `http://localhost:9090/` 查看 Prometheus UI

    {{< figure src="/posts/Prometheus/images/promethus.png" title="MCP 资源工作流程" >}}

### 配置文件

`prometheus.yaml`

```yaml
# my global config
global:
  scrape_interval: 15s # Set the scrape interval to every 15 seconds. Default is every 1 minute.
  evaluation_interval: 15s # Evaluate rules every 15 seconds. The default is every 1 minute.
  # scrape_timeout is set to the global default (10s).

# Alertmanager configuration
alerting:
  alertmanagers:
    - static_configs:
        - targets:
          # - alertmanager:9093

# Load rules once and periodically evaluate them according to the global 'evaluation_interval'.
rule_files:
  # - "first_rules.yml"
  # - "second_rules.yml"

# A scrape configuration containing exactly one endpoint to scrape:
# Here it's Prometheus itself.
scrape_configs:
  # The job name is added as a label `job=<job_name>` to any timeseries scraped from this config.
  - job_name: "prometheus"

    # metrics_path defaults to '/metrics'
    # scheme defaults to 'http'.

    static_configs:
      - targets: ["localhost:9090"]
```

#### 配置文件说明

- `global`: 全局配置
- `alerting`: Alertmanager 相关配置
- `rule_files`: 规则相关配置，可以指定记录规则和告警规则
- `scrape_configs`: 采集配置

{{< admonition type="note" open="true" >}}
完整配置信息请查看官方文档：[prometheus configuration](https://prometheus.io/docs/prometheus/latest/configuration/configuration/)
{{< /admonition >}}

```bash
$ docker run -d \
    --name=prometheus \
    -p 9090:9090 \
    -v /path/to/prometheus.yaml:/etc/prometheus/prometheus.yaml \
    prom/prometheus 
```

Prometheus 数据存储在容器内的 `/prometheus` 目录中，因此每次重启容器时都会清除数据。
如果想要保存数据，需要为容器设置持久存储（或绑定挂载）：

```bash
# Create persistent volume 
$ docker volume create prometheus-data

# Start Prometheus container
$ docker run -d \
    --name = prometheus \
    -p 9090:9090 \
    -v /path/to/prometheus.yaml:/etc/prometheus/prometheus.yaml \
    -v prometheus-data:/prometheus \
    prom/prometheus 
```

---

### 采集指标

Prometheus 通过在目标节点的 HTTP 端口上采集 metric 数据来监控目标节点。
可以通过 [http://127.0.0.1:9090/metrics](http://127.0.0.1:9090/metrics) 查看指标数据。

- `#` 开头的时 metric 相关注释
- `prometheus_target_xxxxx` 是相应的 metric 数据

---

### Grafana 可视化

Grafana 是一款开源的数据可视化工具，支持多种数据源（如 Graphite、InfluxDB、OpenTSDB、Prometheus、Elasticsearch等）并且具有快速灵活的客户端图表，
提供了丰富的仪表盘插件和面板插件，支持多种展示方式，如折线图、柱状图、饼图、点状图等，满足用户不同的可视化需求。

#### 安装 grafana

使用以下命令快速开启一个轻量级 grafana 容器环境。

```bash
$ docker run -d --name=grafana \
    -p 3000:3000 grafana/grafana-oss
```

启动成功后，使用浏览器打开 [http://localhost:3000](http://localhost:3000)。默认的登录账号是"admin" / "admin"。

{{< admonition type="note" title="Grafana 配置 Prometheus 地址" open="true" >}}
由于 prometheus 和 grafana 都采用了 docker 环境，因此在 Grafana 中配置 Connection 的 `Prometheus server URL` 使用的是 `http://host.docker.internal:9090` 而不是 `http://localhost:9090`
{{< /admonition >}}

---

### Exporter 采集数据

**在 Prometheus 的架构设计中，Prometheus Server 并不直接负责监控特定的目标，其主要任务负责数据的收集、存储以及对外提供数据查询支持。**

因此为了能够能够监控到某些指标，如主机的CPU使用率，我们需要使用到 Exporter。

**Prometheus 周期性的从 Exporter暴露的HTTP服务地址（通常是/metrics）拉取监控样本数据。**

广义上讲所有可以向 Prometheus 提供监控样本数据的程序都可以被称为一个 Exporter。
而一个 Exporter 实例被称为 target，Prometheus 通过轮询的方式定期从这些 target 中获取样本数据，Exporter 有两种运行方式：

- 独立运行（需使用独立运行的 Exporter 上报运行状态）
    1. 不能直接提供 HTTP 接口，如监控 Linux 系统状态指标。
    2. 项目发布时间较早，不支持 Prometheus 监控接口，如 MySQL、Redis；
- 集成到应用中（主动暴露运行状态给 Prometheus）
    1. 适用于需要较多自定义监控指标的项目。目前一些开源项目就增加了对 Prometheus 监控的原生支持，如 Kubernetes，ETCD 等。
    2. 可以在业务代码中增加自定义指标数据上报至 Prometheus 。

#### 社区提供的 exporter

Prometheus 社区提供了丰富的 exporter 实现，涵盖了从基础设施、数据库、中间件等各个方面的监控功能。这些 exporter 可以实现大部分通用的监控需求。

完整的第三方库支持[Exporter 文档](https://prometheus.io/docs/instrumenting/exporters/)

|领域|Exporter|
|---|---|
| 数据库 |  MySQL server exporter (official)、MSSQL server exporter、Elasticsearch exporter、MongoDB exporter、Redis exporter 等 |
| 消息队列 | Kafka exporter, RabbitMQ exporter, RocketMQ exporter, NSQ exporter等 |
| API服务 |AWS ECS exporter，Azure Health exporter, Cloudflare exporter等|

还有一些第三方软件默认提供 Prometheus 格式的指标数据，因此不需要单独的 Exporter：

- [Envoy](https://www.envoyproxy.io/docs/envoy/latest/operations/admin.html#get--stats?format=prometheus)
- [Etcd (direct)](https://github.com/coreos/etcd)
- [Flink](https://github.com/apache/flink)
- [Grafana](https://grafana.com/docs/grafana/latest/administration/view-server/internal-metrics/)
- [Kong](https://github.com/Kong/kong-plugin-prometheus)
- [Kubernetes (direct)](https://github.com/kubernetes/kubernetes)
- [RabbitMQ](https://rabbitmq.com/prometheus.html)

---

### 自己应用代码添加 Prometheus 监控

**当想为自己的应用程序添加 Prometheus 监控时，就需要使用 Prometheus Client 库编写代码，并在应用程序实例上的 HTTP 端点定义和公开内部 metric**

#### Prometheus Go Client 示例

Install requirements

```bash
go get github.com/gin-gonic/gin
go get github.com/prometheus/client_golang/prometheus
go get github.com/prometheus/client_golang/prometheus/promhttp
```

```go
package main

import (
 "math/rand/v2"
 "regexp"
 "strconv"

 "github.com/gin-gonic/gin"
 "github.com/prometheus/client_golang/prometheus"
 "github.com/prometheus/client_golang/prometheus/collectors"
 "github.com/prometheus/client_golang/prometheus/promhttp"
)

// 自定义业务指标：业务状态码 Counter
var statusCounter = prometheus.NewCounterVec(
 prometheus.CounterOpts{
  Name: "api_response_status_count",
 },
 []string{"method", "path", "status"},
)

func initRegistry() *prometheus.Registry {
 // New Registry
 reg := prometheus.NewRegistry()

 // Add Go 编译信息
 reg.MustRegister(collectors.NewBuildInfoCollector())

 // Go runtime metrics
 reg.MustRegister(collectors.NewGoCollector(
    collectors.WithGoCollectorRuntimeMetrics(
            collectors.GoRuntimeMetricsRule{Matcher: regexp.MustCompile("/.*")})))

 // 注册自定义的业务指标
 reg.MustRegister(statusCounter)

 return reg
}

func main() {
 r := gin.Default()
 r.GET("/ping", func(ctx *gin.Context) {
  // Mock 业务逻辑
  status := 0
  if rand.IntN(10)%3 == 0 {
   status = 1
  }

  // 记录业务指标
  statusCounter.WithLabelValues(
   ctx.Request.Method,
   ctx.Request.URL.Path,
   strconv.Itoa(status),
  ).Inc()

  ctx.JSON(200, gin.H{
   "status": status,
   "msg":    "pong",
  })
 })

 // 对外提供 /metrics 接口，支持 prometheus 采集
 reg := initRegistry()
 r.GET("/metrics", gin.WrapH(
  promhttp.HandlerFor(
   reg,
   promhttp.HandlerOpts{
    Registry: reg,
   },
  )))

 _ = r.Run("127.0.0.1:5568")
}
```

---

## Reference

- [Prometheus](https://prometheus.io/)
- [Prometheus+Grafana+Go服务自建监控系统入门](https://www.xhyonline.com/?p=1492)
- [Prometheus Server](https://github.com/prometheus/prometheus)
- [clientlibs](https://prometheus.io/docs/instrumenting/clientlibs/)
- [Push Gateway](https://github.com/prometheus/pushgateway)
- [exporters](https://prometheus.io/docs/instrumenting/exporters/)
- [alertmanager](https://github.com/prometheus/alertmanager)
- [时间序列](https://en.wikipedia.org/wiki/Time_series)
- [Prometheus 官方配置文档](https://prometheus.io/docs/prometheus/latest/configuration/configuration/)


---

> Author: [kyden](https://github.com/kydance)  
> URL: http://kydance.github.io/posts/92bd56a/  

