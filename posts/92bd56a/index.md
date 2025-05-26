# Prometheus


{{< admonition type=abstract title="导语" open=true >}}
**这是导语部分**
{{< /admonition >}}

<!--more-->

## 简介

Prometheus 是目前主流的一个开源监控系统和告警工具包，它可以与 Kubernetes 等现代基础设施平台配合，轻松集成到云原生环境中，提供对容器华应用、微服务架构等的全面监控。

Prometheus 受启发于 Google 的 Brogmon 监控系统，从2012年开始由前Google工程师在 Soundcloud 以开源软件的形式进行研发，并且于2015年早期对外发布早期版本。

2016年5月继 Kubernetes 之后成为第二个正式加入CNCF基金会的项目，同年6月正式发布1.0版本。

2017年底发布了基于全新存储层的2.0版本，能更好地与容器平台、云平台配合。

Prometheus 于2016年加入云计算基金会，成为继 Kubernetes 之后的第二个托管项目。

Prometheus 收集并存储其指标作为时间序列数据，即指标信息与其记录的时间戳一起存储，同时存储的还有可选的称为标签的键值对。

### 特性

Prometheus 的主要特性包括:

- 多位数据模型：包含 Metric 名称和键值对标识的时间序列数据
- PromQL：一种可以灵活利用上述维度数据的查询语言
- 不依赖于分布式存储：单个服务器节点是自治的
- Pull 模式：通过基于 HTTP 的拉模式（Pull）进行时间序列数据收集
- 可以通过一个中间网关（Push Gateway）以推模式上报时间序列数据
- 通过服务发现或静态配置发现监控目标
- 支持多种模式的图表和仪表盘

### 何为 Metric ?

Metric 就是用数字来测量/度量。时间序列，指的是记录一段时间内的变化。

例如，
对于 web 服务器，可以测量请求耗时；
对于数据库，可以测量活动连接数、活动查询数等。

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

### 架构

Prometheus 可直接或间接通过推送网关（Push Gateway）抓取监控指标（**适用于短时任务**），它在本地存储所有抓取到的样本数据，并在这些数据上执行一系列规则，以从现有数据中汇总并记录新的时间序列或生成告警；

可以使用 Grafana 或其他 API 消费对收集到的数据进行可视化展示。

## 数据模型

Prometheus 将所有数据都存储为[时间序列](https://en.wikipedia.org/wiki/Time_series)：
属于同一指标（metric）和同一组标注维度（label）的带时间戳的值流。

除了存储的时间序列外，prometheus 还可以根据查询结果生成临时派生时间序列。

```
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

### Metric name

每个时间序列都由其指标名称和称为标签的可选键值对唯一标识。

- 指定要测量的系统的一般功能，例如 `http_request_total`
- 指标名称可以包含 ASCII 字符、数字、下划线和冒号，必须匹配正则表达式 `[a-zA-Z_:][a-zA-Z0-9_:]*`

> 冒号是为用户定义的录制规则保留的，exporter或直接仪器不应使用它们

### Matric labels

- 使 Prometheus 的维度数据模型能够识别同一指标名称的任何标签组合：
它标识了该度量的特定维度示例化，(例如，所有发送 POST 到 `/api/tracks` 的 HTTP 请求)。Prometheus 查询语言允许基于这些维度进行筛选和聚合；
- 任何标签值的更改，包括添加、删除标签，都将创建一个新的时间序列；
- label 可以包含 ASCII 字符、数字以及下划线，必须匹配 `[a-zA-Z_][a-zA-Z0-9_]*`
- 以 `__` (两个下划线) 开头的标签是 Prometheus 内部使用的
- 标签值可以包含任何 Unicode 字符。
- 标签值为空的标签被视为等同于不存在的标签。

更多内容可参考 [Metric name 和 Label 的最佳实践](https://prometheus.io/docs/practices/naming/)。

### Sample

Sample，样本构成实际的时间序列数据，其中每个样本包括：

- 一个 float64 值
- 毫秒精度的时间戳

## Install

```bash
go get github.com/prometheus/prometheus
go get github.com/prometheus/client_golang/prometheus/promhttp

```

## Reference

- [Prometheus](https://prometheus.io/)
- [Prometheus+Grafana+Go服务自建监控系统入门](https://www.xhyonline.com/?p=1492)
- [Prometheus Server](https://github.com/prometheus/prometheus)
- [clientlibs](https://prometheus.io/docs/instrumenting/clientlibs/)
- [Push Gateway](https://github.com/prometheus/pushgateway)
- [exporters](https://prometheus.io/docs/instrumenting/exporters/)
- [alertmanager](https://github.com/prometheus/alertmanager)
- [时间序列](https://en.wikipedia.org/wiki/Time_series)


---

> Author: [kyden](https://github.com/kydance)  
> URL: http://localhost:1313/posts/92bd56a/  

