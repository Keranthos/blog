### 基础语法

#### 标识符

- 第一个字符必须是字母（a-z, A-Z）或下划线 **_**
- 标识符的其他的部分由字母、数字和下划线组成
- 对大小写敏感，count 和 Count 是不同的标识符
- 对长度无硬性限制，但禁止使用保留关键字

> Python 3 允许使用 Unicode 字符作为标识符，可以用中文作为变量名，非 ASCII 标识符也是允许的了
>
> ```python
> 姓名 = "张三"  # 合法
> π = 3.14159   # 合法
> ```

#### 保留关键字（部分）

| **类别**     | **关键字** | **说明**                         |
| :----------- | :--------- | :------------------------------- |
| **逻辑运算** | `and`      | 逻辑与运算                       |
|              | `or`       | 逻辑或运算                       |
|              | `not`      | 逻辑非运算                       |
| **条件控制** | `elif`     | 否则如果（else if 的缩写）       |
| **异常处理** | `try`      | 尝试执行代码块                   |
|              | `except`   | 捕获异常                         |
|              | `finally`  | 无论是否发生异常都会执行的代码块 |
|              | `raise`    | 抛出异常                         |
| **函数定义** | `def`      | 定义函数                         |
|              | `lambda`   | 创建匿名函数                     |
| **类与对象** | `class`    | 定义类                           |
|              | `del`      | 删除对象引用                     |
| **模块导入** | `from`     | 从模块导入特定部分               |
|              | `as`       | 为导入的模块或对象创建别名       |
| **作用域**   | `global`   | 声明全局变量                     |
|              | `nonlocal` | 声明非局部变量（用于嵌套函数）   |
| **异步编程** | `async`    | 声明异步函数                     |
|              | `await`    | 等待异步操作完成                 |
| **其他**     | `assert`   | 断言，用于测试条件是否为真       |
|              | `in`       | 检查成员关系                     |
|              | `is`       | 检查对象身份（是否是同一个对象） |
|              | `pass`     | 空语句，用于占位                 |
|              | `with`     | 上下文管理器，用于资源管理       |
|              | `yield`    | 从生成器函数返回值               |

#### 注释

- 单行注释以`#`开头
- 多行注释使用`'''`或者`"""`包裹
- 无法嵌套注释，如果必须建议在三引号中使用单行注释

#### 行与缩进

- 使用缩进表示代码块而非大括号`{}`
- 缩进的空格数可变，但同一代码块的每条语句缩进空格数必须相同

#### 多行语句

- 如果语句很长，使用反斜杠`\`实现多行语句
- 在`[]` `{}`或`()`中的多行语句无需使用`\`可以直接换行

#### 运算符（部分）

| 运算符类型   | 运算符                   | 描述                                                         |
| ------------ | ------------------------ | ------------------------------------------------------------ |
| 算术         | + - * / % ** //          |                                                              |
| 比较（关系） | == ！= > < >= <=         |                                                              |
| 赋值         | = += -= *= /= %= **= //= |                                                              |
|              | :=                       | 在表达式中同时进行赋值和返回赋值的值<br />如：if (n := len(a)) > 10: |
| 逻辑         | and / or / not           | 即为&& \|\| ！<br />如：if(a and b)                          |
| 位           | & ^ ~ << >>              |                                                              |
| 成员         | in / not in              | 在指定序列中寻找值，返回True / False<br />如：if(a in list)  |
| 身份         | is / is not              | 用于判断两个标识符是否引用自一个对象，类似于id(x) == id(y) id()函数用于获取对象内存地址<br />如：if(a is b) |

#### 数字类型

PY中数字只有四种类型：

- `int`：整数，为长整型
- `bool`
- `float`：浮点数，比如1.23、3E-2
- `complex`：复数，由实部和虚部组成，形式为`a+bj`，比如1+2j

#### 字符串

- 单引号`''`与双引号`""`完全相同
- 三引号`'''`或`"""`可以用于包裹一个多行字符串
- `\`用于转义，在字符串前面加上`r`（表示raw string）则会让字符串中的`\`显示而非转义
- 按字面意义级联字符串会自动转换连接，比如”A“"B""C"会被自动转换为ABC
- 字符串可以用`+`连接，用`*`重复
- 字符串索引从左开始为0，从右开始为-1
- 没有单独的字符类型，字符为长度为1的字符串；字符v和不能改变
- 字符串切片`str[start:end]`中包含`start`、不包含`end`

#### 模块导入

在 python 用 import 或者 from...import 来导入相应的模块。

将整个模块(`somemodule`)导入，格式为：` import somemodule`

从某个模块中导入某个函数,格式为：` from somemodule import somefunction`

从某个模块中导入多个函数,格式为： `from somemodule import firstfunc, secondfunc, thirdfunc`

将某个模块中的全部函数导入，格式为： `from somemodule import \`

#### 其他

- 一行可以有多条语句，中间以`;`分割

- 缩进相同的一组语句构成一个代码块，称之为代码组。像if、while、def和class这样的复合语句，首行以关键字开始，以冒号( : )结束，该行之后的一行或多行代码构成代码组

- `print`默认输出是换行的，不需要换行在变量末尾加上`end=""`

  ```
  print( x, end=" " )
  ```



### 基本数据类型

#### 基本认知

- 变量不需要声明，但在赋值后才会被创建，可以用`del`彻底移除该变量

  ```python
  A = 1
  B = 2
  del A, B
  ```

- 变量没有类型，我们所说的”类型“是变量所指的内存中对象的类型
- 允许同时为多个变量赋值如`a = b = c = 1`或者`a, b, c = 1, 2, "r"`
- 可以通过`type(变量名)`或者 `isinstance` 来判断类型,使用方法：`isinstance(变量名, 类型)`返回True / False

> `isinstance`与`type`的区别在于`type()`不会认为子类是一种父类，而`isinstance()`会

- 常见的数据类型包括以下几类：
  - 不可变数据（如果需要改变值则需要重新分配内存空间）：Number数字、String字符串、Tuple元组
  - 可变数据：List列表、Dictionary字典、Set集合

#### Number数字

- PY中，`bool`是`int`的子类，True(1) 和 False(0) 可以数字相加

- 特殊计算

  ```python
  A / B # 取A除以B的整数部分
  A % B # 取余
  A ** B # 得A的B次方
  ```

- 复数用`a+bj`或者`complex(a,b)`表示，ab都是浮点型

- 数学函数

  | 函数            | 返回值 ( 描述 )                                              |
  | :-------------- | :----------------------------------------------------------- |
  | abs(x)          | 返回数字的绝对值，如abs(-10) 返回 10                         |
  | ceil(x)         | 返回数字的上入整数，如math.ceil(4.1) 返回 5                  |
  | exp(x)          | 返回e的x次幂(ex),如math.exp(1) 返回2.718281828459045         |
  | fabs(x)         | 以浮点数形式返回数字的绝对值，如math.fabs(-10) 返回10.0      |
  | floor(x)        | 返回数字的下舍整数，如math.floor(4.9)返回 4                  |
  | log(x)          | 如math.log(math.e)返回1.0,math.log(100,10)返回2.0            |
  | log10(x)        | 返回以10为基数的x的对数，如math.log10(100)返回 2.0           |
  | max(x1, x2,...) | 返回给定参数的最大值，参数可以为序列。                       |
  | min(x1, x2,...) | 返回给定参数的最小值，参数可以为序列。                       |
  | modf(x)         | 返回x的整数部分与小数部分，两部分的数值符号与x相同，整数部分以浮点型表示。 |
  | pow(x, y)       | x**y 运算后的值。                                            |
  | round(x [,n\])  | 返回浮点数 x 的四舍五入值，如给出 n 值，则代表舍入到小数点后的位数。**其实准确的说是保留值将保留到离上一位更近的一端。** |
  | sqrt(x)         | 返回数字x的平方根。                                          |

- 随机数函数

  | 函数                               | 描述                                                         |
  | :--------------------------------- | :----------------------------------------------------------- |
  | choice(seq)                        | 从序列的元素中随机挑选一个元素，比如random.choice(range(10))，从0到9中随机挑选一个整数。 |
  | randrange ([start,\] stop [,step]) | 从指定范围内，按指定基数递增的集合中获取一个随机数，基数默认值为 1 |
  | random()                           | 随机生成下一个实数，它在[0,1)范围内。                        |
  | seed([x\])                         | 改变随机数生成器的种子seed。如果你不了解其原理，你不必特别去设定seed，Python会帮你选择seed。 |
  | shuffle(lst)                       | 将序列的所有元素随机排序                                     |
  | uniform(x, y)                      | 随机生成下一个实数，它在[x,y]范围内。                        |

- 三角函数

  | 函数        | 描述                                              |
  | :---------- | :------------------------------------------------ |
  | acos(x)     | 返回x的反余弦弧度值。                             |
  | asin(x)     | 返回x的反正弦弧度值。                             |
  | atan(x)     | 返回x的反正切弧度值。                             |
  | atan2(y, x) | 返回给定的 X 及 Y 坐标值的反正切值。              |
  | cos(x)      | 返回x的弧度的余弦值。                             |
  | hypot(x, y) | 返回欧几里德范数 sqrt(x*x + y*y)。                |
  | sin(x)      | 返回的x弧度的正弦值。                             |
  | tan(x)      | 返回x弧度的正切值。                               |
  | degrees(x)  | 将弧度转换为角度,如degrees(math.pi/2) ， 返回90.0 |
  | radians(x)  | 将角度转换为弧度                                  |

#### 字符串

- 字符串格式化

  ```python
  print ("我叫 %s 今年 %d 岁!" % ('小明', 10))
  ```

- f-string

  - 以 `f` 开头，后面跟着字符串，字符串中的表达式用大括号 {} 包起来，它会将变量或表达式计算后的值替换进去

  - 实例如下：

    ```python
    >>> name = 'Runoob'
    >>> f'Hello {name}'  # 替换变量
    'Hello Runoob'
    >>> f'{1+2}'         # 使用表达式
    '3'
    ```

- 字符串函数

  | 函数                                         | 描述                                                         |
  | -------------------------------------------- | ------------------------------------------------------------ |
  | count(str, beg = 0, end = len(string))       | 返回 str 在 string 里面出现的次数，如果 beg 或者 end 指定则返回指定范围内 str 出现的次数 |
  | endswith(suffix, beg = 0, end = len(string)) | 检查字符串是否以 suffix 结束，如果 beg 或者 end 指定则检查指定的范围内是否以 suffix 结束，如果是，返回 True,否则返回 False |
  | find(str, beg=0, end=len(string))            | 检测 str 是否包含在字符串中，如果指定范围 beg 和 end ，则检查是否包含在指定范围内，如果包含返回开始的索引值，否则返回-1 |
  | isalnum()                                    | 检查字符串是否由字母和数字组成，即字符串中的所有字符都是字母或数字。如果字符串至少有一个字符，并且所有字符都是字母或数字，则返回 True；否则返回 False |
  | islower()                                    | 如果字符串中包含至少一个区分大小写的字符，并且所有这些(区分大小写的)字符都是小写，则返回 True，否则返回 False |
  | join(seq)                                    | 以指定字符串作为分隔符，将 seq 中所有的元素(的字符串表示)合并为一个新的字符串 |
  | len(string)                                  | 返回字符串长度                                               |
  | max(str)                                     | 返回字符串 str 中最大的字母                                  |
  | replace(old, new [, max\])                   | 把 将字符串中的 old 替换成 new,如果 max 指定，则替换不超过 max 次。 |
  | upper()                                      | 转换字符串中的小写字母为大写                                 |

#### bool布尔类型

- 布尔类型只有两个值：True 和 False
- 布尔类型可以和其他数据类型进行比较，比如数字、字符串等。在比较时，Python 会将 True 视为 1，False 视为 0；同时也可以转换成其他类型
- 布尔类型可以和逻辑运算符一起使用，包括 and、or 和 not。这些运算符可以用来组合多个布尔表达式，生成一个新的布尔值
- 可以使用 `bool()` 函数将其他类型的值转换为布尔值。以下值在转换为布尔值时为 `False`：`None`、`False`、零 (`0`、`0.0`、`0j`)、空序列（如 `''`、`()`、`[]`）和空映射（如 `{}`）。其他所有值转换为布尔值时均为 `True`

#### List列表

- 列表写在方括号之间，元素用逗号隔开

- 列表的截取（`变量名[头下标：尾下标：步长]`，下标不填写即表示到达尽头）、计算（`+`与`*`的用法）与字符串类似，但是列表可更改

- 相关操作

  ```python
  list = ['Google', 'Runoob', 1997, 2000] # 创建列表
  
  list[2] = 2001 # 修改元素
  del list[2] # 移除元素（后面元素递补）
  ```

- 嵌套列表：类似于多重数组

- 列表函数与方法

  | 函数                                   | 描述                                           |
  | -------------------------------------- | ---------------------------------------------- |
  | len(list)                              | 列表元素个数                                   |
  | max(list) / min(list)                  | 列表元素最大/小值                              |
  | list(seq)                              | 元组转换为列表                                 |
  | list.append(obj)                       | 末尾添加新对象                                 |
  | list.count(obj)                        | 统计元素在列表中出现的次数                     |
  | list.extend(seq)                       | 末尾追加另一个序列中的多个值                   |
  | list.index(obj)                        | 列表中找出某个值第一个匹配项的索引位置         |
  | list.insert(index, obj)                | 对象插入列表                                   |
  | list.pop([index = -1])                 | 移除列表中一个元素（默认为尾部元素）并返回其值 |
  | list.remove(obj)                       | 移除列表中某个值的第一个匹配项                 |
  | list.reverse()                         | 反向列表元素                                   |
  | list.sort(key = None, reverse = False) | 对列表排序                                     |
  | list.clear()                           | 清空列表                                       |
  | list.copy()                            | 复制列表                                       |
  | operator.eq(lista, listb)              | 比较两个列表返回True / False                   |

#### Tuple元组

- 元组元素不能修改，写在小括号中，元素间用逗号隔开

- 截取、计算与字符串 / 列表类似

- 创建一个元素的元组时要在该元素后加`,`，防止被误解为普通的值

- 与字符串相同，可以使用`+` \ `+=` \ `*`号进行运算

- 元组方法与描述

  | 方法            | 描述                   |
  | --------------- | ---------------------- |
  | len(tuple)      | 计算元组元素个数       |
  | max(tuple)      | 返回元组中元素最大值   |
  | min(tuple)      | 返回元组中元素最小值   |
  | tuple(iterable) | 将可迭代系列转换为元组 |

  

#### Set集合

- 无序可变的数据类型，用于存储唯一的元素（重复的元素会在创建后自动去除）

- 可以用`{}`表示，元素间用`,`分隔；或者使用`set()`创建集合

  > 创建空集合必须使用`set()`而不是`{}`，`{}`被用于创建空字典

- 可以进行正常的集合计算：`a - b`差集、`a | b`并集、`a & b`交集、`a ^ b`ab中不同时存在的元素

- 集合方法与描述

  | 方法                          | 描述                                                         |
  | :---------------------------- | :----------------------------------------------------------- |
  | set.add()                     | 为集合添加元素                                               |
  | clear()                       | 移除集合中的所有元素                                         |
  | copy()                        | 拷贝一个集合                                                 |
  | difference()                  | 返回多个集合的差集                                           |
  | difference_update()           | 移除集合中的元素，该元素在指定的集合也存在。                 |
  | discard()                     | 删除集合中指定的元素                                         |
  | intersection()                | 返回集合的交集                                               |
  | intersection_update()         | 返回集合的交集。                                             |
  | isdisjoint()                  | 判断两个集合是否包含相同的元素，如果没有返回 True，否则返回 False。 |
  | issubset()                    | 判断指定集合是否为该方法参数集合的子集。                     |
  | issuperset()                  | 判断该方法的参数集合是否为指定集合的子集                     |
  | pop()                         | 随机移除元素                                                 |
  | remove()                      | 移除指定元素                                                 |
  | symmetric_difference()        | 返回两个集合中不重复的元素集合。                             |
  | symmetric_difference_update() | 移除当前集合中在另外一个指定集合相同的元素，并将另外一个指定集合中不同的元素插入到当前集合中。 |
  | union()                       | 返回两个集合的并集                                           |
  | update()                      | 给集合添加元素                                               |
  | len()                         | 计算集合元素个数                                             |

#### Dictionary字典

- 用`{}`标识，无序的键值对集合（键必须为不可变类型且不能重复）

- 创建方式

  ```python
  dict = {}
  dict['one'] = "1 - 菜鸟教程"
  dict[2]     = "2 - 菜鸟工具"
  
  dict([('Runoob', 1), ('Google', 2), ('Taobao', 3)])
  {x: x**2 for x in (2, 4, 6)} # 结果为{2: 4, 4: 16, 6: 36}
  dict(Runoob=1, Google=2, Taobao=3)
  ```

- 字典方法与描述

  | 方法                               | 描述                                                         |
  | ---------------------------------- | ------------------------------------------------------------ |
  | dict.clear()                       | 删除字典内所有元素                                           |
  | dict.copy()                        | 返回一个字典的浅复制                                         |
  | dict.fromkeys()                    | 创建一个新字典，以序列seq中元素做字典的键，val为字典所有键对应的初始值 |
  | dict.get(key, default=None)        | 返回指定键的值，如果键不在字典中返回 default 设置的默认值    |
  | key in dict                        | 如果键在字典dict里返回true，否则返回false                    |
  | dict.items()                       | 以列表返回一个视图对象                                       |
  | dict.keys()                        | 返回一个视图对象                                             |
  | dict.setdefault(key, default=None) | 和get()类似, 但如果键不存在于字典中，将会添加键并将值设为default |
  | dict.update(dict2)                 | 把字典dict2的键/值对更新到dict里                             |
  | dict.values()                      | 返回一个视图对象                                             |
  | dict.pop(key[,default])            | 删除字典 key（键）所对应的值，返回被删除的值。               |
  | dict.popitem()                     | 返回并删除字典中的最后一对键和值。                           |

  

#### bytes类型

- 表示不可变的二进制序列

- 类型中的元素是整数值（0~255）而非`Unicode`字符

- 常用于处理二进制数据如图象、音频、视频等等

- 创建与转化

  ```py
  x = b"hello" # 创建：使用b前缀
  
  x = bytes("hello", encoding="utf-8") # 转化：第二个参数为编码方式，默认为UTF-8编码
  ```

- 该类型支持许多操作方法比如切片、拼接、查找等等，同时由于其不可变因此在修改时需要创建新的对象

  ```python
  x = b"hello"
  y = x[1:3]  # 切片操作，得到 b"el"
  z = x + b"world"  # 拼接操作，得到 b"helloworld"
  ```

- 由于该类型中的元素是整数值，在比较时需要使用相应的整数值

  ```python
  x = b"hello"
  if x[0] == ord("h"):
      print("The first element is 'h'")
  ```

#### 数据类型转换

##### 隐式

整数（较低数据类型）与浮点数运算时自动转换为浮点数（较高数据类型）

##### 显式

显示转换函数

| 函数                  | 描述                                                |
| :-------------------- | :-------------------------------------------------- |
| int(x [,base])        | 将x转换为一个整数（base为进制）                     |
| float(x)              | 将x转换到一个浮点数                                 |
| complex(real [,imag]) | 创建一个复数                                        |
| str(x)                | 将对象 x 转换为字符串                               |
| repr(x)               | 将对象 x 转换为表达式字符串                         |
| eval(str)             | 用来计算在字符串中的有效Python表达式,并返回一个对象 |
| tuple(s)              | 将序列 s 转换为一个元组                             |
| list(s)               | 将序列 s 转换为一个列表                             |
| set(s)                | 转换为可变集合                                      |
| dict(d)               | 创建一个字典。d 必须是一个 (key, value)元组序列。   |
| frozenset(s)          | 转换为不可变集合                                    |
| chr(x)                | 将一个整数转换为一个字符                            |
| ord(x)                | 将一个字符转换为它的整数值                          |
| hex(x)                | 将一个整数转换为一个十六进制字符串                  |
| oct(x)                | 将一个整数转换为一个八进制字符串                    |

> 注：[, 参数]表示这个为可选参数





### 条件与循环

#### 条件控制

一般形式：

```python
if condition_1:
    statement_block_1
elif condition_2:
    statement_block_2
else:
    statement_block_3
```

**注意：**

1、每个条件后面要使用冒号 **:**，表示接下来是满足条件后要执行的语句块。

2、使用缩进来划分语句块，相同缩进数的语句在一起组成一个语句块。

3、在 Python 中没有 `switch...case` 语句，但在 Python3.10 版本添加了 `match...case`，功能也类似

```python
match subject:
    case <pattern_1>:
        <action_1>
    case <pattern_2>:
        <action_2>
    case <pattern_3>:
        <action_3>
    case _: # case _: 当其他 case 都无法匹配时，匹配这条，保证永远会匹配成功
        <action_wildcard> 
```

#### 循环语句

##### while 语句

一般形式

```python
while 判断条件(condition)：
    执行语句(statements)……
```

**注意：**

- python中没有`do...while`循环
- 如果while循环体只有一条语句，可以将其与 while 写在同一行中

使用 else 语句

```python
while <expr>:
    <statement(s)>
else:
    <additional_statement(s)>
```

##### for 语句

一般形式

```python
for <variable> in <sequence>:
    <statements>
else:
    <statements>
```

##### range()函数

如果需要遍历数字序列，可以使用内置 range() 函数。它会生成数列

```python
for i in range(5): # 遍历从0到5
for i in range(5,9): # 指定区间从5到8
for i in range(0, 10, 3): # 设置增量（步长），此处为遍历0/3/6/9
    
# 结合 range() 和 len() 函数以遍历一个序列的索引
a = ['Google', 'Baidu', 'Runoob', 'Taobao', 'QQ']
for i in range(len(a)):
```

##### 其他语句

- `break`：跳出循环体，但是此时不会执行后面的`else`语句（正常终止/条件判别错误会执行）
- `continue`：跳过当前循环块中的剩余语句，然后继续进行下一轮循环
- `pass`：为了保持程序结构的完整性，不做任何事情，一般用做占位语句



### 推导式

一种独特的数据处理方式，可以从一个数据序列构建另一个新的数据序列的结构体，适用于生成列表、字典、集合和生成器

#### 列表推导式

```python
[表达式 for 变量 in 列表] 
[out_exp_res for out_exp in input_list]
或者 
[表达式 for 变量 in 列表 if 条件]
[out_exp_res for out_exp in input_list if condition]
```

- 表达式：列表生成元素表达式，可以是有返回值的函数
- `for 变量 in 列表`：迭代`input_list`将`out_exp`传入`out_exp_res`表达式中
- `if 条件`：用于过滤列表中不符合条件的值

示例

```python
names = ['Bob','Tom','alice','Jerry','Wendy','Smith']
new_names = [name.upper()for name in names if len(name)>3]
# 结果为['ALICE', 'JERRY', 'WENDY', 'SMITH']
```

#### 字典推导式

```python
{ key_expr: value_expr for value in collection }
或
{ key_expr: value_expr for value in collection if condition }
```

示例

```python
listdemo = ['Google','Runoob', 'Taobao']
# 将列表中各字符串值为键，各字符串的长度为值，组成键值对
newdict = {key:len(key) for key in listdemo}
# 结果为{'Google': 6, 'Runoob': 6, 'Taobao': 6}
```

#### 集合推导式

```python
{ expression for item in Sequence }
或
{ expression for item in Sequence if conditional }
```

示例

```python
a = {x for x in 'abracadabra' if x not in 'abc'}
# 结果为{'d', 'r'}
```

#### 元组推导式（生成器表达式）

元组推导式可以利用 range 区间、元组、列表、字典和集合等数据类型，快速生成一个满足指定需求的元组

```python
(expression for item in Sequence )
或
(expression for item in Sequence if conditional )
```

元组推导式和列表推导式的用法也完全相同，只是元组推导式是用 **()** 圆括号将各部分括起来，而列表推导式用的是中括号 **[]**，另外元组推导式返回的结果是一个生成器对象

```python
a = (x for x in range(1,10)) # 返回的是一个生成器对象
tuple(a)       # 使用 tuple() 函数，可以直接将生成器对象转换成元组(1, 2, 3, 4, 5, 6, 7, 8, 9)
```



### 迭代器与生成器

#### 迭代器

##### 相关介绍

- 迭代器是访问集合元素的一种方式
- 迭代器是一个可以记住遍历的位置的对象，迭代器对象从集合的第一个元素开始访问，直到所有的元素被访问完结束。迭代器只能往前不会后退
- 迭代器有两个基本的方法：`iter()` 和 `next()`
- 字符串，列表或元组对象都可用于创建迭代器

##### 基础用法

```python
list=[1,2,3,4] # 字符串，列表或元组对象都可用于创建迭代器
it = iter(list) # 创建迭代器对象
print (next(it)) # 输出迭代器的下一个元素，结果为1

for x in it:
    print (x, end=" ") # 迭代器对象可以使用常规for语句进行遍历

while True: # 也可以使用next()函数遍历
    try:
        print (next(it))
    except StopIteration:
        sys.exit()
```

##### 进阶用法

**创建迭代器**

把一个类作为一个迭代器使用需要在类中实现两个方法 `__iter__()` 与 `__next__()` 

`__iter__()` 方法返回一个特殊的迭代器对象， 这个迭代器对象实现了 __next__() 方法并通过 `StopIteration` 异常标识迭代的完成。

`__next__()` 方法（Python 2 里是 next()）会返回下一个迭代器对象

（示例）创建一个返回数字的迭代器，初始值为 1，逐步递增 1

```python
class MyNumbers:
  def __iter__(self):
    self.a = 1
    return self
 
  def __next__(self):
    x = self.a
    self.a += 1
    return x
 
myclass = MyNumbers()
myiter = iter(myclass)
 
print(next(myiter))
print(next(myiter))
```

**StopIteration**

`StopIteration` 异常用于标识迭代的完成，防止出现无限循环的情况，在 `__next__()` 方法中我们可以设置在完成指定循环次数后触发 `StopIteration` 异常来结束迭代

（示例）在20次迭代后停止进行

```python
class MyNumbers:
  def __iter__(self):
    self.a = 1
    return self
 
  def __next__(self):
    if self.a <= 20:
      x = self.a
      self.a += 1
      return x
    else:
      raise StopIteration
 
myclass = MyNumbers()
myiter = iter(myclass)
 
for x in myiter:
  print(x)
```

#### 生成器

使用了 **yield** 的函数被称为生成器（generator）。

**yield** 是一个关键字，用于定义生成器函数，生成器函数是一种特殊的函数，可以在迭代过程中逐步产生值，而不是一次性返回所有结果。

跟普通函数不同的是，生成器是一个返回迭代器的函数，只能用于迭代操作，更简单点理解生成器就是一个迭代器。

当在生成器函数中使用 **yield** 语句时，函数的执行将会暂停，并将 **yield** 后面的表达式作为当前迭代的值返回。

然后，每次调用生成器的 **next()** 方法或使用 **for** 循环进行迭代时，函数会从上次暂停的地方继续执行，直到再次遇到 **yield** 语句。这样，生成器函数可以逐步产生值，而不需要一次性计算并返回所有结果。

调用一个生成器函数，返回的是一个迭代器对象

（示例）

```python
def countdown(n):
    while n > 0:
        yield n
        n -= 1
 
# 创建生成器对象
generator = countdown(5)
 
# 通过迭代生成器获取值
print(next(generator))  # 输出: 5
print(next(generator))  # 输出: 4
print(next(generator))  # 输出: 3
 
# 使用 for 循环迭代生成器
for value in generator:
    print(value)  # 输出: 2 1
```

以上实例中，**countdown** 函数是一个生成器函数。它使用 yield 语句逐步产生从 n 到 1 的倒数数字。在每次调用 yield 语句时，函数会返回当前的倒数值，并在下一次调用时从上次暂停的地方继续执行。

通过创建生成器对象并使用 next() 函数或 for 循环迭代生成器，我们可以逐步获取生成器函数产生的值。在这个例子中，我们首先使用 next() 函数获取前三个倒数值，然后通过 for 循环获取剩下的两个倒数值。

生成器函数的优势是它们可以按需生成值，避免一次性生成大量数据并占用大量内存。此外，生成器还可以与其他迭代工具（如for循环）无缝配合使用，提供简洁和高效的迭代方式



### with关键词

用于简化资源管理代码

#### 基础用法

```python
with expression [as variable]:
    # 代码块
```

- `expression` 返回一个支持上下文管理协议的对象
- `as variable` 是可选的，用于将表达式结果赋值给变量
- 代码块执行完毕后，自动调用清理方法

示例（文件操作）

```python
with open('example.txt', 'r') as file:
    content = file.read()
    print(content)
# 文件已自动关闭
```

未简化版本

```python
file = open('example.txt', 'r')
try:
    content = file.read()
    # 处理文件内容
finally:
    file.close()
```

#### 工作原理

##### 上下文管理协议

`with` 语句背后是 Python 的上下文管理协议，该协议要求对象实现两个方法：

1. `__enter__()`：进入上下文时调用，返回值赋给 `as` 后的变量
2. `__exit__()`：退出上下文时调用，处理清理工作

##### 异常处理机制

`__exit__()` 方法接收三个参数：

- `exc_type`：异常类型
- `exc_val`：异常值
- `exc_tb`：异常追踪信息

如果 `__exit__()` 返回 `True`，则表示异常已被处理，不会继续传播；返回 `False` 或 `None`，异常会继续向外传播

#### 实际应用场景

```python
# 文件操作
with open('input.txt', 'r') as infile, open('output.txt', 'w') as outfile: # 同时打开多个文件
    content = infile.read()
    outfile.write(content.upper())
# 数据库连接
import sqlite3

with sqlite3.connect('database.db') as conn:
    cursor = conn.cursor()
    cursor.execute('SELECT * FROM users')
    results = cursor.fetchall() # 连接自动关闭
# 线程锁
import threading

lock = threading.Lock()

with lock:
    # 临界区代码
    print("这段代码是线程安全的")
```



### 函数

一般格式

```python
def 函数名（参数列表）:
    "函数说明/" # 函数的第一行语句可以选择性地使用文档字符串—用于存放函数说明
    函数体
```

#### 参数传递

**注意：**在 python 中，类型属于对象，对象有不同类型的区分，变量是没有类型的

```python
a=[1,2,3]
a="Runoob"
```

以上代码中，`[1,2,3]` 是 List 类型，`"Runoob"` 是 String 类型，而变量 a 是没有类型，它仅仅是一个对象的引用（一个指针），可以是指向 List 类型对象，也可以是指向 String 类型对象。

##### 可更改与不可更改对象

- 不可变类型：包括strings, tuples, 和 numbers 。变量赋值 `a=5` 后再赋值 `a=10`，这里实际是新生成一个 int 值对象 10，再让 a 指向它，而 5 被丢弃，不是改变 a 的值，相当于新生成了 a。
- 可变类型：包括list,dict 等等。变量赋值 `la=[1,2,3,4]` 后再赋值 `la[2]=5` 则是将 list la 的第三个元素值更改，本身la没有动，只是其内部的一部分值被修改了

##### 参数传递

- 不可变类型：类似 C++ 的值传递，如整数、字符串、元组。如 fun(a)，传递的只是 a 的值，没有影响 a 对象本身。如果在 fun(a) 内部修改 a 的值，则是新生成一个 a 的对象。
- 可变类型：类似 C++ 的引用传递，如 列表，字典。如 fun(la)，则是将 la 真正的传过去，修改后 fun 外部的 la 也会受影响

#### 参数

调用函数时可使用的正式参数：

- 必需参数：须以正确的顺序传入函数。调用时的数量必须和声明时的一样

- 关键字参数：使用关键字参数来确定传入的参数值，允许函数调用时参数的顺序与声明时不一致

  ```python
  def printinfo( name, age ): # 函数说明
  
  printinfo( age=50, name="runoob" ) # 关键词参数
  ```

- 默认参数：调用函数时，如果没有传递参数，则会使用默认参数

  ```python
  def printinfo( name, age = 35 ): # 默认参数age
  ```

- 不定长参数：需要一个函数能处理比当初声明时更多的参数，声明时不会命名

  ```python
  # 基本语法
  def functionname([formal_args,] *var_args_tuple ):
     "函数_文档字符串"
     function_suite
     return [expression]
  
  # 示例
  def printinfo( arg1, *vartuple ):
     "打印任何传入的参数"
     print ("输出: ")
     print (arg1)
     print (vartuple)
   
  printinfo( 70, 60, 50 ) # 输出结果为70 (60,50)
  
  # 基本语法2：参数带两个星号**，会以字典形式导入
  def functionname([formal_args,] **var_args_dict ):
     "函数_文档字符串"
     function_suite
     return [expression]
  
  # 示例
  def printinfo( arg1, **vardict ):
     "打印任何传入的参数"
     print ("输出: ")
     print (arg1)
     print (vardict)
  
  printinfo(1, a=2,b=3) # 输出1 {'a': 2, 'b': 3}
  
  # 星号单独出现，则*后面的参数必须使用关键词传入
  def f(a,b,*,c):
      return a+b+c
  ```

#### lambda匿名函数

特点：

- lambda 函数是匿名的，它们没有函数名称不需要使用`def`，只能通过赋值给变量或作为参数传递给其他函数来使用
- 通常只包含一行代码（可以具有任意数量的参数但只能有一个表达式），这使得它们适用于编写简单的函数
- 通常用于编写简单的、单行的函数，通常在需要函数作为参数传递的情况下使用，例如在 map()、filter()、reduce() 等函数中
  - `map()`语法`map(function, iterable, ...)`，将函数应用于可迭代对象的每个元素，并返回一个迭代器
  - `fliter()`语法`filter(function, iterable)`，过滤可迭代对象，保留使函数返回 True 的元素
  - `reduce()`语法`reduce(function, iterable[, initializer])`，对可迭代对象中的元素进行累积操作（初始值为可选参数）

语法格式：

```python
lambda arguments: expression
```

- `lambda`是 Python 的关键字，用于定义 lambda 函数。
- `arguments` 是参数列表，可以包含零个或多个参数，但必须在冒号(`:`)前指定。
- `expression` 是一个表达式，用于计算并返回函数的结果

示例

```python
x = lambda a : a + 10
print(x(5)) # 结果为15

# 使用 lambda 函数与 filter() 一起，筛选偶数
numbers = [1, 2, 3, 4, 5, 6, 7, 8]
even_numbers = list(filter(lambda x: x % 2 == 0, numbers))
print(even_numbers)  # 输出：[2, 4, 6, 8]
```



### 装饰器

#### 特点

- 允许你动态地修改函数或类的行为
- 装饰器是一种函数，它接受一个函数作为参数，并返回一个新的函数或修改原来的函数
- 装饰器的语法使用 `@decorator_name` 来应用在函数或方法上
- Python 还提供了一些内置的装饰器，比如 `@staticmethod` 和 `@classmethod`，用于定义静态方法和类方法
- 应用场景
  - 日志记录: 装饰器可用于记录函数的调用信息、参数和返回值
  - 性能分析: 可以使用装饰器来测量函数的执行时间
  - 权限控制: 装饰器可用于限制对某些函数的访问权限
  - 缓存: 装饰器可用于实现函数结果的缓存，以提高性能

#### 基本语法

装饰器：一个接收函数作为输入并返回一个新的包装过后的函数的对象

```python
def decorator_function(original_function):
    def wrapper(*args, **kwargs):
        # 这里是在调用原始函数前添加的新功能
        before_call_code()
        
        result = original_function(*args, **kwargs)
        
        # 这里是在调用原始函数后添加的新功能
        after_call_code()
        
        return result
    return wrapper

# 使用装饰器
@decorator_function
def target_function(arg1, arg2):
    pass  # 原始函数的实现
```

**解析：**decorator 是一个装饰器函数，它接受一个函数` func` 作为参数，并返回一个内部函数 wrapper，在 wrapper 函数内部，你可以执行一些额外的操作，然后调用原始函数 `func`，并返回其结果。

- `decorator_function` 是装饰器，它接收一个函数 `original_function` 作为参数
- `wrapper` 是内部函数，它是实际会被调用的新函数，它包裹了原始函数的调用，并在其前后增加了额外的行为
- 当我们使用 `@decorator_function` 前缀在 `target_function` 定义前，Python会自动将 `target_function` 作为参数传递给 `decorator_function`，然后将返回的 `wrapper` 函数替换掉原来的 `target_function`

示例

```python
def my_decorator(func):
    def wrapper(*args, **kwargs):
        print("在原函数之前执行")
        func(*args, **kwargs)
        print("在原函数之后执行")
    return wrapper

@my_decorator
def greet(name):
    print(f"Hello, {name}!")

greet("Alice")

# 输出结果为：
在原函数之前执行
Hello, Alice!
在原函数之后执行
```

#### 带参数的装饰器

```python
def repeat(num_times):
    def decorator(func):
        def wrapper(*args, **kwargs):
            for _ in range(num_times):
                func(*args, **kwargs)
        return wrapper
    return decorator

@repeat(3)
def say_hello():
    print("Hello!")

say_hello()

# 输出结果为：
Hello!
Hello!
Hello!
```

repeat 函数是一个装饰器工厂，它接受一个参数 num_times，返回一个装饰器 decorator。decorator 接受一个函数 func，并返回一个 wrapper 函数。wrapper 函数会调用 `func` 函数 num_times 次。使用 @repeat(3) 装饰s ay_hell 函数后，调用 say_hello 会打印 "Hello!" 三次

#### 类装饰器

种用于动态修改类行为的装饰器，它接收一个类作为参数，并返回一个新的类或修改后的类。

类装饰器可以用于：

- 添加/修改类的方法或属性
- 拦截实例化过程
- 实现单例模式、日志记录、权限检查等功能

#### 内置装饰器

Python 提供了一些内置的装饰器，例如：

1. `@staticmethod`: 将方法定义为静态方法，不需要实例化类即可调用。
2. `@classmethod`: 将方法定义为类方法，第一个参数是类本身（通常命名为 `cls`）。
3. `@property`: 将方法转换为属性，使其可以像属性一样访问。

示例

```python
class MyClass:
    @staticmethod
    def static_method():
        print("This is a static method.")

    @classmethod
    def class_method(cls):
        print(f"This is a class method of {cls.__name__}.")

    @property
    def name(self):
        return self._name

    @name.setter
    def name(self, value):
        self._name = value

# 使用
MyClass.static_method()
MyClass.class_method()

obj = MyClass()
obj.name = "Alice"
print(obj.name)
```

#### 多个装饰器堆叠

可以将多个装饰器堆叠在一起，它们会按照从下到上的顺序依次应用



### 模块

#### import 语句

`import module1[, module2[,... moduleN]`

一个模块只会被导入一次

**搜索路径**（在PY编译安装时确定，记录在`sys`模块的`path`变量中）：

1. 当前目录
2. 环境变量 `PYTHONPATH` 指定的目录
3. Python 标准库目录
4. `.pth` 文件中指定的目录

#### from … import 语句

从模块中导入一个指定的部分到当前命名空间中

```python
from modname import name1[, name2[, ... nameN]]
from modname import * # 把一个模块的所有内容全都导入到当前的命名空间
```

**给模块起别名**

使用`as`关键字

```python
import numpy as np  # 将 numpy 模块别名设置为 np
from math import sqrt as square_root  # 将 sqrt 函数别名设置为 square_root
```

#### `__name__ `属性

一个模块被另一个程序第一次引入时，其主程序将运行

如果我们想在模块被引入时，模块中的某一程序块不执行，我们可以用 `__name__` 属性来使该程序块仅在该模块自身运行时执行

```python
if __name__ == '__main__':
   print('程序自身在运行') # 被导入时不会运行
else:
   print('我来自另一模块') # 被导入时运行
```

说明：每个模块都有一个 `__name__` 属性

- 如果模块是被直接运行，`__name__` 的值为 `__main__`
- 如果模块是被导入的，`__name__` 的值为模块名

#### `dir()`函数

内置的函数 `dir()` 可以找到模块内定义的所有名称，以一个字符串列表的形式返回

```python
>>> dir(sys)  
['__displayhook__', '__doc__', '__excepthook__', '__loader__', '__name__',
 '__package__', '__stderr__'...]
```

如果没有给定参数，那么 dir() 函数会罗列出当前模块中定义的所有名称



### 文件

#### `open()`

`open()` 将会返回一个 file 对象

`open(filename, mode)`

- filename：包含了你要访问的文件名称的字符串值。
- mode：决定了打开文件的模式：只读，写入，追加等。所有可取值见如下的完全列表。这个参数是非强制的，默认文件访问模式为只读(r

不同模式打开文件的完全列表：

| 模式 | 描述                                                         |
| :--- | :----------------------------------------------------------- |
| r    | 以只读方式打开文件。文件的指针将会放在文件的开头。这是默认模式 |
| rb   | 以二进制格式打开一个文件用于只读。文件指针将会放在文件的开头 |
| r+   | 打开一个文件用于读写。文件指针将会放在文件的开头             |
| rb+  | 以二进制格式打开一个文件用于读写。文件指针将会放在文件的开头 |
| w    | 打开一个文件只用于写入。如果该文件已存在则打开文件，并从开头开始编辑，即原有内容会被删除。如果该文件不存在，创建新文件 |
| wb   | 以二进制格式打开一个文件只用于写入。如果该文件已存在则打开文件，并从开头开始编辑，即原有内容会被删除。如果该文件不存在，创建新文件 |
| w+   | 打开一个文件用于读写。如果该文件已存在则打开文件，并从开头开始编辑，即原有内容会被删除。如果该文件不存在，创建新文件 |
| wb+  | 以二进制格式打开一个文件用于读写。如果该文件已存在则打开文件，并从开头开始编辑，即原有内容会被删除。如果该文件不存在，创建新文件 |
| a    | 打开一个文件用于追加。如果该文件已存在，文件指针将会放在文件的结尾。也就是说，新的内容将会被写入到已有内容之后。如果该文件不存在，创建新文件进行写入 |
| ab   | 以二进制格式打开一个文件用于追加。如果该文件已存在，文件指针将会放在文件的结尾。也就是说，新的内容将会被写入到已有内容之后。如果该文件不存在，创建新文件进行写入 |
| a+   | 打开一个文件用于读写。如果该文件已存在，文件指针将会放在文件的结尾。文件打开时会是追加模式。如果该文件不存在，创建新文件用于读写 |
| ab+  | 以二进制格式打开一个文件用于追加。如果该文件已存在，文件指针将会放在文件的结尾。如果该文件不存在，创建新文件用于读写 |

#### 文件对象的方法

##### `f.read()`

为了读取一个文件的内容，调用 `f.read(size)`, 这将读取一定数目的数据, 然后作为字符串或字节对象返回

size 是一个可选的数字类型的参数。 当 size 被忽略了或者为负, 那么该文件的所有内容都将被读取并且返回

##### `f.readline()` / `f.readlines()`

`f.readline()` 会从文件中读取单独的一行。换行符为 '\n'。`f.readline()` 如果返回一个空字符串, 说明已经已经读取到最后一行

`f.readlines()` 将返回该文件中包含的所有行。

如果设置可选参数 `sizehint`, 则读取指定长度的字节, 并且将这些字节按行分割

##### `f.tell()`

`f.tell()` 用于返回文件当前的读/写位置（即文件指针的位置）。文件指针表示从文件开头开始的字节数偏移量。

`f.tell()` 返回一个整数，表示文件指针的当前位置

##### `f.close()`

在文本文件中 (那些打开文件的模式下没有 b 的), 只会相对于文件起始位置进行定位。

当你处理完一个文件后, 调用 `f.close()` 来关闭文件并释放系统的资源，如果尝试再调用该文件，则会抛出异常（建议使用`with`关键字）



### 错误和异常

#### 异常处理

##### `try/except(...else)`

```python
while True:
    try:
        x = int(input("请输入一个数字: "))
        break
    except ValueError:
        print("您输入的不是数字，请再次尝试输入！")
```

try 语句按照如下方式工作；

- 首先，执行 try 子句（在关键字 try 和关键字 except 之间的语句）
- 如果没有异常发生，忽略 except 子句，try 子句执行后结束
- 如果在执行 try 子句的过程中发生了异常，那么 try 子句余下的部分将被忽略。如果异常的类型和 except 之后的名称相符，那么对应的 except 子句将被执行
- 如果一个异常没有与任何的 except 匹配，那么这个异常将会传递给上层的 try 中

一个 try 语句可能包含多个except子句，分别来处理不同的特定的异常。最多只有一个分支会被执行。

一个except子句可以同时处理多个异常，这些异常将被放在一个括号里成为一个元组如`**except** (RuntimeError, TypeError, NameError):`

最后一个except子句可以忽略异常的名称，它将被当作通配符使用:

```python
except:
    print("Unexpected error:", sys.exc_info()[0])
    raise # 再次抛出该异常
```

**`else`**

`..else`子句可选，但必须放在所有的 except 子句之后。其将在 try 子句没有发生任何异常的时候执行

##### `try-finally`

`finally`在`else`的后面，无论是否有异常都会执行

如果一个异常在 try 子句里（或者在 except 和 else 子句里）被抛出，而又没有任何的 except 把它截住，那么这个异常会在 finally 子句执行后被抛出

#### 抛出异常

使用 raise 语句抛出一个指定的异常

```python
raise [Exception [, args [, traceback]]]
```

`raise`只有唯一一个参数（其他参数为`Exception`的），其必须时一个异常的实例或者类（即`Exception`的子类）



### 面向对象

#### 关键词介绍

- **类(Class):** 用来描述具有相同的属性和方法的对象的集合。它定义了该集合中每个对象所共有的属性和方法。对象是类的实例
- **方法：**类中定义的函数
- **类变量：**类变量在整个实例化的对象中是公用的。类变量定义在类中且在函数体之外。类变量通常不作为实例变量使用
- **数据成员：**类变量或者实例变量用于处理类及其实例对象的相关的数据
- **方法重写：**如果从父类继承的方法不能满足子类的需求，可以对其进行改写，这个过程叫方法的覆盖（override），也称为方法的重写
- **局部变量：**定义在方法中的变量，只作用于当前实例的类
- **实例变量：**在类的声明中，属性是用变量来表示的，这种变量就称为实例变量，实例变量就是一个用 self 修饰的变量
- **继承：**即一个派生类（derived class）继承基类（base class）的字段和方法。继承也允许把一个派生类的对象作为一个基类对象对待。例如，有这样一个设计：一个Dog类型的对象派生自Animal类，这是模拟"是一个（is-a）"关系（例图，Dog是一个Animal）
- **实例化：**创建一个类的实例，类的具体对象
- **对象：**通过类定义的数据结构实例。对象包括两个数据成员（类变量和实例变量）和方法

#### 类对象

支持两种操作：属性引用和实例化

```python
class MyClass:
    """一个简单的类实例"""
    i = 12345
    def f(self):
        return 'hello world'
 
x = MyClass() # 实例化类
 
# 访问类的属性和方法
print("MyClass 类的属性 i 为：", x.i)
print("MyClass 类的方法 f 输出为：", x.f())
```

- **构造函数**：`__init__()`
- **self** 代表类的实例，而非类
  - 实例方法中必须以`self`作为自己的第一个参数
  - 第一个参数可以使用其他名称（不建议），但含义都是指向自己这个实例

#### 继承

```python
class DerivedClassName(modname.BaseClassName):
# 如果基类与派生类在同一个模块中可以不写模块名

# 示例
class people:
    #定义基本属性
    name = ''
    age = 0
    #定义私有属性,私有属性在类外部无法直接进行访问
    __weight = 0
    #定义构造方法
    def __init__(self,n,a,w):
        self.name = n
        self.age = a
        self.__weight = w
    def speak(self):
        print("%s 说: 我 %d 岁。" %(self.name,self.age))
 
# 单继承示例
class student(people):
    grade = ''
    def __init__(self,n,a,w,g):
        #调用父类的构函
        people.__init__(self,n,a,w)
        self.grade = g
    #覆写父类的方法
    def speak(self):
        print("%s 说: 我 %d 岁了，我在读 %d 年级"%(self.name,self.age,self.grade))
        
#另一个类，多继承之前的准备
class speaker():
    topic = ''
    name = ''
    def __init__(self,n,t):
        self.name = n
        self.topic = t
    def speak(self):
        print("我叫 %s，我是一个演说家，我演讲的主题是 %s"%(self.name,self.topic))
 
#多继承
class sample(speaker,student):
    a =''
    def __init__(self,n,a,w,g,t):
        student.__init__(self,n,a,w,g)
        speaker.__init__(self,n,t)
```

#### 方法重写

如果父类方法的功能不能满足需求，可以在子类重写你父类的方法

```python
class Parent:        # 定义父类
   def myMethod(self):
      print ('调用父类方法')
 
class Child(Parent): # 定义子类
   def myMethod(self):
      print ('调用子类方法')
 
c = Child()          # 子类实例
c.myMethod()         # 子类调用重写方法
super(Child,c).myMethod() #用子类对象调用父类已被覆盖的方法
```

super() 函数是用于调用父类(超类)的一个方法。

#### 类的属性与方法

类的私有属性`__private_attrs`：两个下划线开头，声明该属性为私有，不能在类的外部被使用或直接访问。在类内部的方法中使用时 `self.__private_attrs`。

类的私有方法`__private_method`：两个下划线开头，声明该方法为私有方法，只能在类的内部调用 ，不能在类的外部调用。在类内部的方法中使用时`self.__private_methods`

示例

```python
class JustCounter:
    __secretCount = 0  # 私有变量
    publicCount = 0    # 公开变量
 
    def count(self):
        self.__secretCount += 1
        self.publicCount += 1
        print (self.__secretCount)
 
counter = JustCounter()
counter.count()
counter.count()
print (counter.publicCount)
print (counter.__secretCount)  # 报错，实例不能访问私有变量
```

**类的专有方法**

- **__init__ :** 构造函数，在生成对象时调用
- **__del__ :** 析构函数，释放对象时使用
- **__repr__ :** 打印，转换
- **__setitem__ :** 按照索引赋值
- **__getitem__:** 按照索引获取值
- **__len__:** 获得长度
- **__cmp__:** 比较运算
- **__call__:** 函数调用
- **__add__:** 加运算
- **__sub__:** 减运算
- **__mul__:** 乘运算
- **__truediv__:** 除运算
- **__mod__:** 求余运算
- **__pow__:** 乘方

##### 运算符重载

```python
 class Vector:
   def __init__(self, a, b):
      self.a = a
      self.b = b
 
   def __str__(self):
      return 'Vector (%d, %d)' % (self.a, self.b)
   
   def __add__(self,other):
      return Vector(self.a + other.a, self.b + other.b)
 
v1 = Vector(2,10)
v2 = Vector(5,-2)
print (v1 + v2)
```





















