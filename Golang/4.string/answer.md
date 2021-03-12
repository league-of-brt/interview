# golang中的字符咋存？utf8咋编码？string啥结构？

* [golang中的字符咋存？utf8咋编码？string啥结构？](#golang中的字符咋存utf8咋编码string啥结构)
  * [1 计算机怎么表达字符？](#1-计算机怎么表达字符)
  * [2 字符集是什么？](#2-字符集是什么)
  * [3 怎么存储字符串？如何划分字符串中的字符边界？](#3-怎么存储字符串如何划分字符串中的字符边界)
    * [3\.1 直接编码](#31-直接编码)
    * [3\.2 定长编码](#32-定长编码)
    * [3\.3 变长编码](#33-变长编码)
  * [4 字符串变量是什么结构？怎么知道内存中字符串的起止点？](#4-字符串变量是什么结构怎么知道内存中字符串的起止点)
  * [4\.1 终结符](#41-终结符)
  * [4\.2 字节长度处理](#42-字节长度处理)
  * [5 golang中的字符串能修改吗？](#5-golang中的字符串能修改吗)

参考：https://www.bilibili.com/video/BV1ff4y1m72A

## 1 计算机怎么表达字符？

1. 比特 = bit = b
2. 字节 = byte = 8b

> 1. 1个b可以是0，可以是1。这样1个b就能表达2个数字。
>
> 2. 2个b可以是00、01、10、11，能表达0、1、2、3四个数字。也就是说位数决定了能表达的数字的范围。
>
> 3. 8个b组成了1个byte，能表达 2 ^ 8 = 256 个数字，也就是从0到255。
>
> 4. 依次类推，2个byte = 16个b，可以表达2 ^ 16 = 65536，从0到65535。 

数字可以用b来表示（例如7用八位二进制表示就是00000111），是可以直接通过计算得出来的，但是现在有个字符A，那怎么通过计算得到呢？

很简单，可以通过编号去中转一下：

|字符|编号|二进制|
|----|----|----|
|A|65|01000001|

如果要存储字符A，那么就存着这个01000001数值，读取字符时，就可以根据一个映射表，把数值还原成编号，找到映射的字符就行了。

## 2 字符集是什么？

现在我们知道A对应65、B对应66......

|字符|编号|二进制|
|----|----|----|
|A|65|01000001|
|B|66|01000010|
|C|67|01000011|
.....

像这样把你想要表达的字符都用编号表示，然后收集起来，最终就会得到一张映射表，也就是一个字符集，比如我们最早接触的字符集就是ASCII。

最早的ASCII，只能表示英文字母、阿拉伯数字、西文符、控制字符，那没有中文咋办？

很简单，只要拓展中文字符集就好。理论上，需要表现的字符是有限的，但是数值是无限的。所以无论你要什么符号，你都可以用一个数值去映射。

于是在ASCII的基础上，1980年的GB2312就支持了简体中文。

有简体，没有繁体怎么行？

那就在拓展字符集咯。1984年的BIG5就支持了繁体字。

那么问题来了，拓展字符的需求是无限的，人们刚刚确定了一个字符集，就发现不支持某些字符，又要去确定下一个字符集.....这样就周而复始没个完了，非常不方便开发使用。

所以，与其不停推出自定义字符集，不如制定全球通用的标准字符集，于是1994年公布全球第一个通用字符集Unicode。

总之，字符集促成了字符和二进制的映射关系，字符集就是一个映射表。我们一般使用的都是标准化的字符集，当然你也可以定义自己的字符集。

## 3 怎么存储字符串？如何划分字符串中的字符边界？

现在有这么个需求，一串字符串`eggo世界`，想要存起来，要怎么存？

### 3.1 直接编码

肯定是先按照字符集，去找到每个字符对应的映射然后存二进制咯？

|字符|编号|二进制|
|----|----|----|
|e|101|01100101|
|g|103|01100111|
|g|103|01100111|
|o|111|01101111|
|世|19990|01001110 00010110|
|界|30028|01110101 01001100|

现在我们得到一大串二进制，直接拼起来吗？

如果照搬编号，那不就变成了`011001010110011101100111011011110100111000010110 0111010101001100`，那么问题来了，怎么知道这串数值要怎么划分呢？

1. 如果按照`01100101 01100111 01100111 01101111 0100111000010110 0111010101001100`划分，那么就会得到6个字符`eggo世界`。
2. 但是如果按照`0110010101100111 0110011101101111 0100111000010110 0111010101001100`划分，就会的带4个字符`敧杧世界`。
3. 所以说，就算是同一组数值，按照不同的规则去划分，会影响字符集映射的结果。

### 3.2 定长编码

那我们统一补齐最大位数不就好了？我们先设定用32位数（够大了吧？）表达一个字符，然后高位补零就行啦！

|字符|编号|二进制|
|----|----|----|
|e|101|00000000 00000000 00000000 01100101|
|g|103|00000000 00000000 00000000 01100111|
|g|103|00000000 00000000 00000000 01100111|
|o|111|00000000 00000000 00000000 01101111|
|世|19990|00000000 00000000 01001110 00010110|
|界|30028|00000000 00000000 01110101 01001100|

这样确实可以，就是需要更多的空间去存无用的字符，太浪费了。

前面提到，字符集收录的字符越多，就要更多的数值去映射。数值越大，就要用到更多的位数，定长编码就会造成更多的浪费。

### 3.3 变长编码

我们可以根据编号去决定区间，设置一定的规则。

|编号|字节|位数|模板|
|----|----|----|----|
|0-127|1|8|0XXXXXXX|
|128-2047|2|16|110XXXXX 10XXXXXX|
|2048-65535|3|24|1110XXXX 10XXXXXX 10XXXXXX|
.....

练习一下，现在有一串数据str，是`01100101 11100100 10111000 10010110 .....`，我们把这串字符串还原出来。

1. 开头是`01100101`，第一个标识位就是0。对照着模板，我们得知str这串数据，开头前8位1个字节代表了1个字符。那么取出前8位，拿到`01100101`，减去开头的标示位0，就能得到`1100101`，直接转换成十进制，就变成了101，就是字符`e`。所以str字符串第一个字符就是`e`。
2. 接下来，发现第9-12位有个1110，根据模板，后面24位一共3个字节代表了第2个字符。
于是拿到`11100100 10111000 10010110`，截取标志位，得到`0100 111000 010110`，组合起来，得到`01001110 00010110`二进制编号，转化成十进制得到19990，对应汉字`世`。

再练习一下，尝试获得`界`这个字符的数值。

1. `界`这个字符，查表得到30028，在2048-65535之间，于是使用`1110XXXX 10XXXXXX 10XXXXXX`模板。30028的二进制是`01110101 01001100`,从头至尾填到模板中，最终得到`11100111 10010101 10001100`。
2. 还原一下`11100111 10010101 10001100`，从开头`111`，就能知道套用`1110XXXX 10XXXXXX 10XXXXXX`模板，截取标志位，得到`0111 010101 001100`，组合得到`01110101 01001100`，转换成十进制就变成了30028，查表得到`界`。

**其实这个就是golang默认的utf-8的编码规则了。**

总之，现在拿到一串数值，然后和你约定字符集，就能根据字符集的编码规则把字符还原出来。

## 4 字符串变量是什么结构？怎么知道内存中字符串的起止点？

现在我们定义一个string变量str，内容是`eggo`。

```golang
str := "eggo"
```

string变量，结构中有一个data指针，指向`eggo`开头的`e`：

|01|02|03|04|05|06|07|08|09|
|----|----|----|----|----|----|----|----|----|
|e|g|g|o|..|..|..|..|..|

那么问题又来了，我怎么知道这个字符串到内存的哪里结束呢？

万一内存是这样的：

|01|02|03|04|05|06|07|08|09|
|----|----|----|----|----|----|----|----|----|
|e|g|g|o|a|b|..|..|..|

我通过data指针拿str的数值的时候，连续的读过去，不就读成`eggoab`了吗。

## 4.1 终结符

为了搞清楚字符串到哪里终结，C语言的处理方式是在字符串结尾放一个终结符`\0`。

|01|02|03|04|05|06|07|08|09|
|----|----|----|----|----|----|----|----|----|
|e|g|g|o|\0|a|b|..|..|

这样读到\0，就知道这个字符串结束了。

那么问题又来了，假如我的字符串本身就带有一个终结符，那么这个字符串读的时候岂不是会提早结束？

比如我的字符串是`eggo\0`:

|01|02|03|04|05|06|07|08|09|
|----|----|----|----|----|----|----|----|----|
|e|g|g|o|\0|\0|a|b|..|

读到`\0`就结束了，那我的字符串就永远表示不全了！

## 4.2 字节长度处理

golang选了更加聪明的方式。

string变量中除了data外，还存有一个len。data指针指向内存中字符串的开头，len决定了读取的字节数。

比如上面的：

```golang
str := "eggo"
```

len就是`eggo`的字节长度，就是4。

如果：

```golang
str := "eggo世界"
```

世界每个字符占了三个字节，那么`eggo世界`就是十个字节长度，len就是10。

## 5 golang中的字符串能修改吗？

不能。可以读，没办法改。

```golang
str := "eggo"
fmt.Printf("%c\n",str[1]) // ok
str[2]='o' // not ok
```

这是因为：

1. golang默认字符串内容是不会修改的，编译器会把字符串内容分配到只读内存段，就只能读不能改了。
2. 字符串变量可以共用底层字符串内容，如果允许修改字符串，会造成不可预测的错误。

比如：

``` golang
str1 := "eggo"
str2 := "go"
```

str2中的data指针可能指向str1内存中的`go`中的`g`，如果我们能修改s1的内容,那么s2的内容也会被修改，这显然是不能预测的。

如果我执行这样的代码：
```golang
str := "eggo"
str = "hello"
```
内存中的变量改了吗？

不会的，只是另外在内存中开辟了一块存`hello`的区域，str的data指向了`hello`的`h`而已，之前的`eggo`区域没有任何变动。

这样呢？
```golang
str := "eggo"
bs := ([]byte)str
fmt.Printf("%c\n",str[2])
```
这只是在内存中申请了一块slice内存，并且把str的内容复制了一份过来而已，原来的str还是不会有任何变动。