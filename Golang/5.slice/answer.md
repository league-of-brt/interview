# 谈谈你对slice的理解？

* [谈谈你对slice的理解？](#谈谈你对slice的理解)
  * [1 slice底层结构](#1-slice底层结构)
  * [2 new，make的不同](#2-newmake的不同)
    * [2\.1 new](#21-new)
    * [2\.2 make](#22-make)
  * [3 练习，画出一个string slice](#3-练习画出一个string-slice)
    * [3\.1 错误示例](#31-错误示例)
    * [3\.2 正确示例](#32-正确示例)
  * [4 共用内存区域](#4-共用内存区域)
  * [5 扩容机制](#5-扩容机制)
    * [5\.1 预估规则](#51-预估规则)
    * [5\.2 需要多大的内存](#52-需要多大的内存)
    * [5\.3 申请内存](#53-申请内存)
      * [5\.3\.1 内存管理模块](#531-内存管理模块)
      * [5\.3\.2 匹配到合适的内存规格](#532-匹配到合适的内存规格)
  * [6 扩容例子](#6-扩容例子)

## 1 slice底层结构

slice是什么结构，这个问题可以解构为：

1. 元素存哪里（data指针，指向连续内存的开头元素位置）
2. 存了多少个元素（len，指导连续内存的有效范围，即存放元素的内存的范围）
3. 可以存多少元素（cap，指导是否进行扩容）

也就是说一个slice由data、len、cap三个部分组成。

如果我们声明一个整型slice：

```golang
var s []int
```

或者：

```golang
var s = []int{}
```

|data|len|cap|ram|
|----|----|----|----|
|uintptr|0|0|nil|

1. data字段是一个指针，是一个uintptr，零值是0，指向连续内存的开头元素位置。 `var s []int` 和 `var s = []int{}` 只分配了s这个切片结构，没有分配实际存放元素的内存，所以data没有指向任何地址，为0值。
2. s中显然没有任何元素，所以连续内存的有效范围为0，len = 0。
3. s没有连续内存，不能放入任何元素，cap = 0。

有些同学可能会问，为啥cap为0呢，我不是可以用append怼一些元素到s里面吗？那就是能存元素咯，怎么也不会是0吧。

如果执行：

```golang
s := []int{}
s[0] = 1 // panic
```

len是控制是否可以用下标访问的，因为s的len = 0，所以直接进行赋值操作会panic。cap是控制是否需要扩容的，空数组没有必要指导扩容，所以这里的cap为0。

之所以能append，是因为append中有分配连续内存/扩容的操作，所以执行append肯定是可以怼元素的。但是在append之前，是不能操作的。

## 2 new，make的不同

### 2.1 new

```golang
var s = new([]int)
```

|data|len|cap|ram|
|----|----|----|----|
|uintptr|0|0|nil|

1. new不会分配内存。**new的返回值，就是slice结构的起始地址。所以这里的s并不是一个slice，而是一个地址。**
2. len = 0。
3. cap = 0。

如果此时使用append，就会给slice申请内存：

```golang
var s = new([]int)
*s = append(*s, 1)
```

### 2.2 make

```golang
var s = make([]int,2,5)
```

|data|len|cap|ram|
|----|----|----|----|
|uintptr|2|5|[0,0,0,0,0]|

1. make会分配能存放5个元素的连续内存（这里并不是刚好仅能存放），并且认为已经填入了2个元素，且前2个元素都是int类型初始值，即为0。
2. 已经存了2个元素，len = 2。
3. 能存5个元素，cap = 5。

对元素进行赋值：

```golang
var s = make([]int,2,5)
s = append(s, 1)
```

|data|len|cap|ram|
|----|----|----|----|
|uintptr|3|5|[0,0,1,0,0]|

1. data指向的连续内存，在两个元素之后，刚好能依次append进去3个元素。如果再多append进去一个元素，比如4，就会发生扩容。
2. len = 2 + 1 = 3。
3. 未发生扩容，cap = 5。

有些同学可能觉得，诶这个内存区域看起来像是放得下5个元素，我能不能这样操作：

```golang
var s = make([]int,2,5)
s = append(s, 1)
s[3] = 3
s[4] = 4
```

会panic。

这个内存区域看似能放得下5个元素，但是len控制slice的下标访问，len = 3，即有效位只有前3位。如果超过len，访问后2位会爆数组越界错误。

## 3 练习，画出一个string slice

画出以下代码的内存示意图：

```golang
var s = new([]string)
var str = "test"
*s = append(*s, str)
```

### 3.1 错误示例

1. s是一个指针，指向slice。
2. 定义了一个str对象。
2. str对象由指针和长度组成，指针指向内存中`test`的首个字节，长度为4。
3. slice由data、len、cap组成，data指向str对象。

![pic](https://brt-1303999354.cos.ap-shanghai.myqcloud.com/QQ%E6%88%AA%E5%9B%BE20210313120253.png)

为什么错了？

因为通过append，申请了新的连续内存，slice的data指向新的连续内存中的首个全新的string对象，不再是外层的str对象。

### 3.2 正确示例

1. s是一个指针，指向slice。
2. 定义了一个str对象。
2. str对象由指针和长度组成，指针指向内存中`test`的首个字节，长度为4。
3. slice由data、len、cap组成，append的时候申请了新的连续内存。首个元素为一个新的string，内容和str一样，但是地址不同。data指向这个新的string。
4. 由于字面量一样，新旧string的指针，指向的都是连续内存中的首个元素（在Golang的常量区）。

![pic](https://brt-1303999354.cos.ap-shanghai.myqcloud.com/QQ%E6%88%AA%E5%9B%BE20210313121055.png)

## 4 共用内存区域

```golang
// Test5 ...
func Test5(t *testing.T) {
	array := [10]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	s1 := array[1:3]
	s2 := array[5:9]
	fmt.Println(s1)      // [1 2]
	fmt.Println(len(s1)) // 2
	fmt.Println(cap(s1)) // 9
	fmt.Println(s2)      // [5 6 7 8]
	fmt.Println(len(s2)) // 4
	fmt.Println(cap(s2)) // 5
}
```

|var|data|len|cap|ram|
|----|----|----|----|----|
|s1|uintptr|2|9|[1,2,0,0,0,0,0,0,0]|
|s2|uintptr|4|5|[5,6,7,8,0]|

这里大家可能会奇怪，为啥s1中只有2个元素，len为2，但是cap是9？

画张图大家就明白了：

![pic](https://brt-1303999354.cos.ap-shanghai.myqcloud.com/QQ%E6%88%AA%E5%9B%BE20210312231334.png)

s1、s2已经指明共用array的内存区域了，如果array中的元素改变，s1、s2中的元素也会跟着变：

```golang
// Test5 ...
func Test5(t *testing.T) {
	array := [10]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	s1 := array[1:3]
	s2 := array[5:9]
	fmt.Println(s1)      // [1 2]
	fmt.Println(len(s1)) // 2
	fmt.Println(cap(s1)) // 9
	fmt.Println(s2)      // [5 6 7 8]
	fmt.Println(len(s2)) // 4
	fmt.Println(cap(s2)) // 5

	array[2] = 3
	fmt.Println(s1)      // [1 3] 我变了！
}
```

如果此时s2增加了两个元素，发生了扩容，那么就不一样了：

```golang
// Test5 ...
func Test5(t *testing.T) {
	array := [10]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	s1 := array[1:3]
	s2 := array[5:9]
	fmt.Println(s1)      // [1 2]
	fmt.Println(len(s1)) // 2
	fmt.Println(cap(s1)) // 9
	fmt.Println(s2)      // [5 6 7 8]
	fmt.Println(len(s2)) // 4
	fmt.Println(cap(s2)) // 5

	s2 = append(s2, 9, 10)
	fmt.Println(len(s2)) // 6
	fmt.Println(cap(s2)) // 10

	array[5] = 1    // 尝试修改array，看看s2的底层数组会不会改变
	fmt.Println(s2) // [5 6 7 8 9 10]
}
```

1. `s2 = append(s2, 9, 10)` 时，发现已经超出cap，发生扩容操作。
2. 申请内存，copy原来的元素5、6、7、8，插入9、10。

现在内存示意图就会变成这样：

![pic](https://brt-1303999354.cos.ap-shanghai.myqcloud.com/QQ%E6%88%AA%E5%9B%BE20210312232753.png)

s2的data指向的已经不是array的内存区域了，改变array的元素不会导致s2的元素改变。

## 5 扩容机制

需要注意这三点：

扩容步骤：
1. 预测扩容后容量（根据预估规则计算newCap）
2. 计算需要多大的内存（类型*容量）
3. 申请内存（匹配到合适的内存规格）

### 5.1 预估规则

1. if oldCap * 2 < cap, newCap = cap
2. if oldCap * 2 >= cap：
    1. oldLen < 1024, newCap >= oldCap * 2
    2. oldLen > 1024, newCap >= oldCap * 1.25

现在网上有很多资料认为 

1. `oldLen < 1024, newCap = oldCap * 2`
2. `oldLen > 1024, newCap = oldCap * 1.25`

但是这篇文章否认了这一观点：https://www.bookstack.cn/read/qcrao-Go-Questions/%E6%95%B0%E7%BB%84%E5%92%8C%E5%88%87%E7%89%87-%E5%88%87%E7%89%87%E7%9A%84%E5%AE%B9%E9%87%8F%E6%98%AF%E6%80%8E%E6%A0%B7%E5%A2%9E%E9%95%BF%E7%9A%84.md

结论是计算完 `newCap = oldCap * 1.25/2` 之后，根据内存分配策略，还会对newcap做一个内存对齐，进行内存对齐之后，newCap会**大于等于**oldCap的1.25倍或者两倍。

### 5.2 需要多大的内存

newCap个元素需要多大的内存，要怎么计算呢？

难道是简单的：newCap * 元素大小吗？比如我一个int占8个字节，扩容到5，就申请8*5 = 40个字节的内存吗？

原则上是这样的，我们只需要40个字节就能完成扩容操作，现在我们需要进行申请。

### 5.3 申请内存

在Golang中，想要申请内存，并不是直接和操作系统交涉的，并不是你向操作系统要40字节，它就会给你。为什么？程序向操作系统申请内存，是要经过用户态到内核态的转换的，大家想想，如果能随便向操作系统要内存，就会导致用户态和内核态的频繁转换，非常消耗性能。

#### 5.3.1 内存管理模块

解决思路：临时向操作系统申请内存开销大，我们能不能整个“内存池”，预先一次性向操作系统申请一批内存过来，放在“内存池”里面，用的时候就直接取，那不就不需要临时申请内存了？

那么，这个内存池又该由谁去维护呢？操作系统？虚拟机？我们自己写程序？

那当然是虚拟机帮我们做这个功能了，在Golang中，有个内存管理模块，他会提前向操作系统申请一批内存。我们的程序如果要内存，得向这个内存管理模块申请，不能向操作系统申请。

#### 5.3.2 匹配到合适的内存规格

内存管理模块会提前向操作系统申请一批内存，然后把这些内存分为常用的规格管理起来，比如8字节、16字节、32字节、48字节......

我们申请内存时，内存管理模块会帮我们匹配足够大并且最接近的规格。比如我们申请40字节，就会匹配到48字节的内存块，我们实际上拿到的是48字节，能装下6个int，是超过我们的扩容需要的。

## 6 扩容例子

```golang
a := []string{"My", "name", "is"}
a = append(a, "eggo")
```

1. 调用append方法，发现数据量大于cap。开始扩容。
2. 计算newCap:
    1. 现在oldCap = 3, len = 3
    2. 我们要插入一个元素，需要cap = 4
    3. 3 x 2 = 6 > 4 && 3 < 1024
    4. newCap = 3 * 2 = 6
3. 计算需要多大的内存:
    1. 6 * 16 Byte（string由指针和长度组成，指针8字节，长度int也为8字节） = 96 Byte
4. 匹配到合适的内存规格：
    1. 从小到大逐一匹配，匹配到96字节大小的内存块，申请完成。

接下来的步骤是：

5. 将老slice中的数据复制到新的内存区域。
6. 将append的元素添加到新的底层数组中。
7. a的data指向新的内存区域的首个元素。
8. append方法返回新的slice，新的slice的data变化了，长度没有变化，容量增大了。
