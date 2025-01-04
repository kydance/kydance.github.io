# 浅析 MySQL


{{&lt; admonition type=abstract title=&#34;导语&#34; open=true &gt;}}
导语内容
{{&lt; /admonition &gt;}}

&lt;!--more--&gt;

## 插入数据 `INSERT`

```sql
INSERT INTO TableName
    [columnName_1, columnName_2, ...]
VALUES(
    &lt;columnValue_1&gt;, &lt;columnValue_2, ...&gt;, ...
    );
```

## 删除数据 `DELETE`

```sql
DELETE FROM TableName
WHERE Condition;
```

{{&lt; admonition type=Tips title=&#34;省略 WHERE 子句&#34; open=true &gt;}}
当使用 `DELETE` 且省略 `WHERE` 子句时，它将删除所有行
{{&lt; /admonition &gt;}}

## 更新数据 `UPDATE`

```sql
UPDATE TableName
SET columnName_1 = columnValue_1,
    columnName_2 = columnValue_2, ...
WHERE Condition;
```

{{&lt; admonition type=Tips title=&#34;UPDATE&#34; open=true &gt;}}

- 更新多个 column 时，只需要使用一个 `SET` 命令，且每一个`columnName = columnValue`之间使用 `,` 分割（最后一列不需要）
- 可以嵌套**子查询**
{{&lt; /admonition &gt;}}

## 检索数据 `SELECT`

```sql
SELECT [DISTINCT]
    &lt;column1, column2, ...&gt;
FROM
    &lt;Table1 [[AS]&lt;aliasName&gt;], table2, ...&gt;
WHERE
    &lt;Condition1 [AND / OR / NOT] Condition2 ...&gt;
GROUP BY
    &lt;Column1, Column2, ...&gt;
HAVING
    &lt;Column1, Column2, ...&gt;
ORDER BY
    &lt;Column1 [ASC / DESC], Column2 [ASC / DESC], ...&gt;
LIMIT
    &lt;A&gt; OFFSET &lt;B&gt;;
```

{{&lt; admonition type=Tips title=&#34;SELECT&#34; open=true &gt;}}

- `DISTINCT` 指示数据库只返回不同的值，且**必须放在所有列名之前，作用于所有的列**
- `LIMIT` 指定返回的行数, `OFFSET` 指定从哪里开始计算，`LIMIT 4 OFFSET 3 &lt;=&gt; LIMIT 3, 4`
- `ORDER BY` 根据一个或多个列的名字进行排序，默认升序 `ASC`，降序 `DESC`
- `WHERE` 指定筛选条件：
  - `=` / `&lt;&gt;` / `!=` / `&lt;` / `&lt;=` / `!&lt;` / `&gt;` /
`&gt;=` / `!&gt;` / `BETWEEN ... AND ...` / `IS NULL`
  - `AND` / `OR` / `IN` / `NOT`
  - `LIKE`: `%` 表示匹配**任何字符出现任意次数**，`_` 表示匹配**任意单个字符**
- `GROUP BY` 吧数据进行逻辑分组，以便能对每一个组进行聚合计算
- `HAVING` 过滤分组，类似于 `WHERE`，且支持所有 `WHERE` 操作:
  - `WHERE` 过滤行，`HAVING` 过滤分组
  - `WHERE` 在分组前过滤，`HAVING` 在分组后过滤（`WHERE`排除掉的行不包括在分组内）
  - `HAVING` 应与 `GROUP BY` 结合使用

{{&lt; /admonition &gt;}}

---

## 表操作

### 查看表结构

```sql
DESC &lt;tbName&gt;;

SHOW CREATE TABLE &lt;tbName&gt;
```

### 建表

```SQL
CREATE TABLE &lt;tbName&gt; (
  ... [COMMENTS],
  ... [COMMENTS],
  ...
) [COMMENTS];
```

### 表结构更新

```SQL
-- 添加字段
ALTER TABLE &lt;tbName&gt; ADD &lt;columnName&gt; &lt;dataType&gt;(LENGTH) [&lt;COMMENT&gt; &#39;comments&#39;] [CONSTRAINT];

-- 修改字段数据类型
ALTER TABLE &lt;tbName&gt; MODIFY &lt;columnName&gt; &lt;newDataType&gt;(LENGTH);

-- 修改字段数据类型和字段名
ALTER TABLE &lt;tbName&gt; CHANGE &lt;oldName&gt; &lt;newName&gt; &lt;dataType&gt;(LENGTH) [COMMENT &#39;comments&#39;][CONSTRAINT];

-- 删除字段
ALTER TABLE &lt;tbName&gt; DROP &lt;columnName&gt;;
```

### 表删除

永久删除表结构

```SQL
DROP TABLE [IF EXISTS] &lt;TableName&gt;;
```

删除表格内容，保留表结构

```SQL
TRUNCATE [TABLE] &lt;TableName&gt;;
```

{{&lt; admonition type=Tips title=&#34;&#34; open=true &gt;}}

- `DELETE`: 逐行删除表中记录数据
- `TRUNCATE`: 直接删除原来的表，然后重新建立一个一模一样的表，执行速度比 `DELETE` 快

{{&lt; /admonition &gt;}}

### 表重命名

```SQL
ALTER TABLE oldTable RENAME [TO | AS] newTable;
```

## MySQL 中的数据类型

| Type | Size | Signed Range |
| --- | --- | --- |
| TINYINT | 1 Byte | (-128, 127) |
| SMALLINT | 2 Byte | (-32768, 32767) |
| MEDIUMINT | 3 Byte | (-8388608, 8388607) |
| INT / INTEGER | 4 Byte | $(−2^32, 2^32−1)$ |
| BIGINT | 8 Byte | (−2^63 , 2^63−1) |
| FLOAT | 4 Byte | (-3.402823466E&#43;38, 3.402823466351E&#43;38) |
| DOUBLE | 8 Byte | (-1.7976931348623157E&#43;308, 1.7976931348623157E&#43;308) |
| DECIMAL | | |

| Type | Size(Byte) | DESC |
| --- | --- | --- |
| CHAR | 0 - 255 | **定长字符串** (需指定长度），性能较VARCHAR更高些 |
| VARCHAR | 0 - 65535 | **变长字符串** (需指定长度) |
| TINYBLOB | 0 - 255 | 不超过 255 个字符的二进制数据 |
| TINYTEXT | 0 - 255 | 短文本字符串 |
| BLOB | 0 - 65535 | 二进制形式的长文本数据 |
| TEXT | 0 - 65535 | 长文本数据 |
| MEDIUMBLOB | 0 - 16777215 | 二进制形式的中等长度文本数据 |
| MEDIUMTEXT | 0 - 16777215 | 中等长度文本数据 |
| LONGBLOB | 0 - 4294967295 | 二进制形式的极大长度长文本数据 |
| LONGTEXT | 0 - 4294967295 | 极大长度文本数据 |

| Type | Size(Byte) | Range | Format | DESC |
| --- | --- | --- | --- | --- |
| DATE | 3 | 1000-01-01 至 9999-12-31 | YYYY-MM-DD | 日期值 |
| TIME | 3 | -838:59:59 至 838:59:59 | HH:MM:SS | 时间值或持续时间 |
| YEAR | 1 | 1901 至 2155 | YYYY | 年份值 |
| DATETIME | 8 | 1000-01-01 00:00:00 至 9999-12-31 23:59:59 | YYYY-MM-DD HH:MM:SS	| 混合日期和时间值 |
| TIMESTAMP | 4 | 1970-01-01 00:00:01 至 2038-01-19 03:14:07 | YYYY-MM-DD HH:MM:SS | 混合日期和时间值，时间戳 |


---

> : [kyden](https:github.com/kydance)  
> URL: http://kyden.us.kg/posts/mysql/  

