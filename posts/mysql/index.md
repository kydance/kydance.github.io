# MySQL 核心操作精讲：从基础语法到实战应用


{{< admonition type=abstract title="导语" open=true >}}
数据库是现代应用程序的核心基础设施，而 MySQL 作为最流行的关系型数据库之一，其重要性不言而喻。本文将带你系统掌握 MySQL 的核心操作，从基础的增删改查到表结构管理，从数据类型选择到性能优化，为你提供一份全面且实用的 MySQL 开发指南。无论你是数据库新手，还是想要提升数据库开发技能的开发者，都能从本文中获得实用的知识和技巧。
{{< /admonition >}}

<!--more-->

## 插入数据 `INSERT`

```sql
INSERT INTO TableName
    [columnName_1, columnName_2, ...]
VALUES(
    <columnValue_1>, <columnValue_2, ...>, ...
    );
```

## 删除数据 `DELETE`

```sql
DELETE FROM TableName
WHERE Condition;
```

{{< admonition type=Tips title="省略 WHERE 子句" open=true >}}
当使用 `DELETE` 且省略 `WHERE` 子句时，它将删除所有行
{{< /admonition >}}

## 更新数据 `UPDATE`

```sql
UPDATE TableName
SET columnName_1 = columnValue_1,
    columnName_2 = columnValue_2, ...
WHERE Condition;
```

{{< admonition type=Tips title="UPDATE" open=true >}}

- 更新多个 column 时，只需要使用一个 `SET` 命令，且每一个`columnName = columnValue`之间使用 `,` 分割（最后一列不需要）
- 可以嵌套**子查询**
{{< /admonition >}}

## 检索数据 `SELECT`

```sql
SELECT [DISTINCT]
    <column1, column2, ...>
FROM
    <Table1 [[AS]<aliasName>], table2, ...>
WHERE
    <Condition1 [AND / OR / NOT] Condition2 ...>
GROUP BY
    <Column1, Column2, ...>
HAVING
    <Column1, Column2, ...>
ORDER BY
    <Column1 [ASC / DESC], Column2 [ASC / DESC], ...>
LIMIT
    <A> OFFSET <B>;
```

{{< admonition type=Tips title="SELECT" open=true >}}

- `DISTINCT` 指示数据库只返回不同的值，且**必须放在所有列名之前，作用于所有的列**
- `LIMIT` 指定返回的行数, `OFFSET` 指定从哪里开始计算，`LIMIT 4 OFFSET 3 <=> LIMIT 3, 4`
- `ORDER BY` 根据一个或多个列的名字进行排序，默认升序 `ASC`，降序 `DESC`
- `WHERE` 指定筛选条件：
  - `=` / `<>` / `!=` / `<` / `<=` / `!<` / `>` /
`>=` / `!>` / `BETWEEN ... AND ...` / `IS NULL`
  - `AND` / `OR` / `IN` / `NOT`
  - `LIKE`: `%` 表示匹配**任何字符出现任意次数**，`_` 表示匹配**任意单个字符**
- `GROUP BY` 吧数据进行逻辑分组，以便能对每一个组进行聚合计算
- `HAVING` 过滤分组，类似于 `WHERE`，且支持所有 `WHERE` 操作:
  - `WHERE` 过滤行，`HAVING` 过滤分组
  - `WHERE` 在分组前过滤，`HAVING` 在分组后过滤（`WHERE`排除掉的行不包括在分组内）
  - `HAVING` 应与 `GROUP BY` 结合使用

{{< /admonition >}}

---

## 表操作

### 查看表结构

```sql
DESC <tbName>;

SHOW CREATE TABLE <tbName>
```

### 建表

```SQL
CREATE TABLE <tbName> (
  ... [COMMENTS],
  ... [COMMENTS],
  ...
) [COMMENTS];
```

### 表结构更新

```SQL
-- 添加字段
ALTER TABLE <tbName> ADD <columnName> <dataType>(LENGTH) [<COMMENT> 'comments'] [CONSTRAINT];

-- 修改字段数据类型
ALTER TABLE <tbName> MODIFY <columnName> <newDataType>(LENGTH);

-- 修改字段数据类型和字段名
ALTER TABLE <tbName> CHANGE <oldName> <newName> <dataType>(LENGTH) [COMMENT 'comments'][CONSTRAINT];

-- 删除字段
ALTER TABLE <tbName> DROP <columnName>;
```

### 表删除

永久删除表结构

```SQL
DROP TABLE [IF EXISTS] <TableName>;
```

删除表格内容，保留表结构

```SQL
TRUNCATE [TABLE] <TableName>;
```

{{< admonition type=Tips title="" open=true >}}

- `DELETE`: 逐行删除表中记录数据
- `TRUNCATE`: 直接删除原来的表，然后重新建立一个一模一样的表，执行速度比 `DELETE` 快

{{< /admonition >}}

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
| FLOAT | 4 Byte | (-3.402823466E+38, 3.402823466351E+38) |
| DOUBLE | 8 Byte | (-1.7976931348623157E+308, 1.7976931348623157E+308) |
| DECIMAL | | |

| Type | Size(Byte) | DESC |
| --- | --- | --- |
| CHAR | 0 - 255 | **定长字符串** (需指定长度)，性能较VARCHAR更高些 |
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
| DATETIME | 8 | 1000-01-01 00:00:00 至 9999-12-31 23:59:59 | YYYY-MM-DD HH:MM:SS | 混合日期和时间值 |
| TIMESTAMP | 4 | 1970-01-01 00:00:01 至 2038-01-19 03:14:07 | YYYY-MM-DD HH:MM:SS | 混合日期和时间值，时间戳 |

## 关系型数据库三大范式

### 第一范式

1NF，**数据库表中的字段都是单一属性的，不可再分；每一个属性都是原子项，不可分割**。
如果数据库设计不满足第一范式，就不称为关系型数据库。

### 第二范式

2NF，**数据库表中，取保非主键列完全依赖于主键列**，即每个非主键列都完全依赖于主键列，而不是部分依赖。

如果一个表中某一个字段 A 的值是由另一个字段或一组字段 B 的值来决定，那么 A 就依赖于 B。

### 第三范式

3NF，**取保非主键之间不存在传递依赖**，即一个非主键列都只依赖于主键列，而不依赖于其他非主键列。

## MySQL 数据库面试题

### 1. MySQL 数据库的索引类型？它们有何作用？

索引，是对数据库表中一个或多个列的值进行排序的结构

优点

- 大大加快数据的检索速度
- 创建唯一性索引，保证数据库表中每一行数据的唯一性
- 加速表与表之间的连接
- 在使用分组、排序子句进行数据检索时，可以显著减少查询时分组和排序的时间

缺点

- 索引需要占用数据表以外的物理存储空间
- 创建索引和维护索引需要花费一定的时间
- 当对表进行更新操作时，索引需要被重建，降低了数据的维护速度

索引类型

- 唯一索引 `UNIQUE`: 此索引的每一个索引值只对应唯一的数据记录
  - 对于单列唯一性索引，它将保证数据表中该列的每一个值都是唯一的
  - 对于多列唯一性索引，它将保证数据表中该组列的每一组值都是唯一的
- 主键索引 `PRIMARY KEY`: 数据库表中常常有一列或列组合，其值为宜标识数据表中的唯一记录
  - 主键索引时唯一索引的特定类型，要求主键中的每一个值都是唯一的
  - 主键索引还可以是组合索引，即主键中的多个字段值可以组成唯一的索引值

索引实现方式

- B+Tree 索引
- Hash Table 索引
- 位图索引

---

### 2. 聚簇索引和非聚簇索引的区别是什么？

聚簇索引

- **表示表中存储的数据按照索引的顺序存储，检索效率比非聚簇索引高，但对数据更新影响较大**
- **一个表只能有一个**
- **在物理存储方面是连续存储的**

非聚簇索引

- **表示数据存储在一个地方，索引存储在另一个地方，索引带有指向数据的存储位置，索引检索效率比聚簇索引低，但对数据更新影响小**
- **一个表可以有多个**
- **在逻辑上连续的，物理存储可以不连续**

---

### 3. MySQL 中常见索引类型，以及它们的优势和劣势

Hash Table 索引模型

- 键 - 值（key - value）存储数据
- 优点：查询快速、记录增加速度快
- 劣势：key 的存储不是采用有序存储，区间查询速度慢（需要多次调用 Hash 函数）
- 使用场景：只需要等值查询（没有区间查询），例如 `Memcached`, `NoSQL`

有序数组索引模型

- 等值查询、范围查询场景下的性能都很优秀
- 只适用于**静态存储**引擎（更新数据的成本太高）

N 叉树搜索树索引模型

- N 值的大小取决于数据块的大小，InnoDB 的一个整数字段索引，N 差不多为 1200（当树高为 4 时，可以存储 $N^3$ 个数值，将近 17 亿）；
考虑到树根的数据块总是在内存中，一个 10 亿行的表上的一个整数字段的索引，查询一个值最多只需要 3 次访盘；
- 主键长度越小，普通索引的叶子节点就越小，普通索引占用的空间也就越小

### 4. 覆盖索引


---

> Author: [kyden](https://github.com/kydance)  
> URL: http://kydance.github.io/posts/mysql/  

