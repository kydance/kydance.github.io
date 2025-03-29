# JSON 完全指南：从基础概念到编解码实战


{{&lt; admonition type=abstract title=&#34;导语&#34; open=true &gt;}}
在现代 Web 开发中，JSON 已成为最受欢迎的数据交换格式。它不仅轻量级、易读易写，更因其语言无关性而被广泛应用于 API 设计和前后端通信。本文将带你全面了解 JSON，从其设计理念到实际应用，从基础语法到编码最佳实践。无论你是前端开发者还是后端工程师，掌握 JSON 都是提升开发效率的必备技能。
{{&lt; /admonition &gt;}}

&lt;!--more--&gt;

## I. 何为 JSON ?

{{&lt; figure src=&#34;/posts/json-encoder-decoder/JSON_vector_logo.png&#34; title=&#34;&#34; width=100 height=10 &gt;}}

&gt; JSON 源于对实时服务器到浏览器会话通信协议的需求，无需使用 Flash 或 Java 小程序等浏览器插件，这是 2000 年代初期使用的主要方法。

JSON(JavsScript Object Notation，JavsScript 对象表示法)，由美国程序员道格拉斯·克罗克福特构想和设计的一种轻量级资料交换格式。
其内容由属性和值所组成，具有易于阅读和处理的优势。

JSON是独立于编程语言的资料格式，其不仅是JavaScript的子集，也采用了C语言家族的习惯用法，目前也有许多编程语言都能够将其解析和字符串化，其广泛使用的程度也使其成为通用的资料格式。

{{&lt; admonition tip &#34;JSON&#34; ture &gt;}}
扩展名：`.json`

互联网媒体类型: `application/json`

类型代码: `TEXT`

统一类型标识: `public.json`

格式类型: 数据交换

扩展自: [JavaScript](https://zh.wikipedia.org/wiki/JavaScript)

标准: [RFC 7159](https://tools.ietf.org/html/rfc7159), [ECMA-404](http://www.ecma-international.org/publications/files/ECMA-ST/ECMA-404.pdf)

网站: [json.org](http://json.org/)
{{&lt; /admonition &gt;}}

---

## II. JSON 基本数据类型

### 数值

十进制数值，不能有前导0，可以为负数，可以有小数部分，不区分整数与浮点数。
也可以用 `e`/`E` 表示指数部分，不能包含非数（如 `NaN`）。

{{&lt; admonition tip &#34;JSON&#34; ture &gt;}}
JavaScript 用双精度浮点数表示所有数值（后来也支持 BigInt）
{{&lt; /admonition &gt;}}

---

### 字符串

以双引号 `&#34;&#34;` 括起来的零个或多个 [Unicode](https://zh.wikipedia.org/wiki/Unicode) [码位](https://zh.wikipedia.org/wiki/%E7%A0%81%E4%BD%8D)，支持反斜杠开始的转义字符序列。

---

### 布尔值

表示为 `true` 或 `false`

---

### 数组

有序的零个或多个值

每个值可以为任意类型，并使用方括号`[]`包裹，元素之间使用逗号`,`分隔，例如`[val, val]`.

---

### 对象

若干无序的&#34;key-value&#34;对（key-value pairs），其中 `key` 只能是字符串，并以花括号`{}`包裹，多个&#34;key-value&#34;之间使用逗号`,`分隔，`key` 与 `value` 之间使用冒号`:`分隔

建议但不强制要求对象中的键是独一无二的。

---

### 空值

表示为 `null`

---

{{&lt; admonition note &#34;JSON&#34; ture &gt;}}
JSON 交换时必须编码为[UTF-8](https://zh.wikipedia.org/wiki/UTF-8)。

转义序列可以为: `\\`, `\&#34;`, `\/`, `\b`, `\f`, `\n`, `\r`, `\t` 或 Unicode16 进制转义字符序列（`\u`后面跟随 4 位 16 进制数字）。
{{&lt; /admonition &gt;}}

## III. Example

```JSON
{
    &#34;firstName&#34;: &#34;John&#34;,
    &#34;lastName&#34;: &#34;Smith&#34;,
    &#34;sex&#34;: &#34;male&#34;,
    &#34;age&#34;: 25,
    &#34;address&#34;: 
    {
        &#34;streetAddress&#34;: &#34;21 2nd Street&#34;,
        &#34;city&#34;: &#34;New York&#34;,
        &#34;state&#34;: &#34;NY&#34;,
        &#34;postalCode&#34;: &#34;10021&#34;
    },
    &#34;phoneNumber&#34;: 
    [
        {
        &#34;type&#34;: &#34;home&#34;,
        &#34;number&#34;: &#34;212 555-1234&#34;
        },
        {
        &#34;type&#34;: &#34;fax&#34;,
        &#34;number&#34;: &#34;646 555-4567&#34;
        }
    ]
 }
```

---

## IV. 与其他格式的比较

### XML

**JSON与XML最大的不同，在于XML是一个完整的标记语言，而JSON不是.**

这使得XML在程序判读上需要比较多的功夫。主要的原因在于XML的设计理念与JSON不同。XML利用标记语言的特性提供了绝佳的延展性（如XPath），在数据存储，扩展及高级检索方面具备对JSON的优势，而JSON则由于比XML更加小巧，以及浏览器的内置快速解析支持，使得其更适用于网络数据传输领域。

---

### YAML

**在功能和语法上，JSON 都是 YAML 语言的一个子集**

特别是，YAML 1.2规范指定&#34;任何JSON格式的文件都是YAML格式的有效文件&#34;。最常见的 YAML 解析器也能够处理 JSON。

版本 1.2 之前的 YAML 规范没有完全涵盖 JSON，主要是由于 YAML 中缺乏本机 UTF-32 支持，以及对逗号分隔空格的要求；此外，JSON 规范还包括 `/* */` 样式的注释。

YAML 最重要的区别是语法扩展集，它们在 JSON 中没有类似物：

- 关系数据支持：在 YAML 文档中，可以引用以前在文件/流中找到的锚点；通过这种方式，您可以表达递归结构。
- 支持除基元之外的可扩展数据类型，如字符串、数字、布尔值等。
- 支持带缩进的块语法；
- 它允许您在不使用不必要的符号的情况下描述结构化数据：各种括号、引号等。

## V. Reference

- [JSON](https://zh.wikipedia.org/wiki/JSON)
- [JavaScript](https://zh.wikipedia.org/wiki/JavaScript)
- [RFC 7159](https://tools.ietf.org/html/rfc7159)
- [ECMA-404](http://www.ecma-international.org/publications/files/ECMA-ST/ECMA-404.pdf)
- [json.org](http://json.org/)
- [Unicode](https://zh.wikipedia.org/wiki/Unicode)
- [UTF-8](https://zh.wikipedia.org/wiki/UTF-8)


---

> Author: [kyden](https://github.com/kydance)  
> URL: http://kydance.github.io/posts/json-encoder-decoder/  

