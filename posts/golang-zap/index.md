# Golang Zap


{{&lt; admonition type=abstract title=&#34;导语&#34; open=true &gt;}}
导语内容
{{&lt; /admonition &gt;}}

&lt;!--more--&gt;

如果你的应用采用容器化部署，其实更建议将日志输出到标准输出。
容器平台一般都具有采集容器日志的能力。
采集日志时，可以选择从标准输出采集或者容器中的日志文件采集，如果是从日志文件进行采集，通常需要配置日志采集路径，但如果是从标准输出采集，则不用。
所以，如果将日志直接输出到标准输出，则可以不加配置直接复用容器平台已有的能力，做到记录日志和采集日志完全解耦。

定制开发步骤分为以下几步：

创建一个封装了 zap.Logger 的自定义 Logger；

编写创建函数，创建 zapLogger 对象；

创建 *zap.Logger 对象；

实现日志接口。


---

> Author: [kyden](https:github.com/kydance)  
> URL: http://kyden.us.kg/posts/golang-zap/  

