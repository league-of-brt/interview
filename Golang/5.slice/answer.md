# 谈谈你对slice的理解？

## 1 slice底层结构

slice是什么结构，这个问题可以解构为：

1. 元素存哪里（array）
2. 存了多少个元素（len）
3. 可以存多少元素（cap）

也就是说一个slice由array、len、cap三个部分组成。

如果我们声明一个整型slice：

### 1.1 var s []int

```golang
var s []int
```

|array|len|cap|
|----|----|----|
|nil|0|0|

1. array为nil，因为 `var ints []int` 只分配了s这个切片结构，但是并没有分配底层数组，所以底层数组为nil，即array为nil。
2. s中显然没有任何元素，len = 0。
3. s中一个元素也没法存，cap = 0。

### 1.2 s := []int{}

```golang
var s = []int{}
```

|array|len|cap|
|----|----|----|
|[]|0|0|

1. array为空数组，`var s = []int{}` 会构建底层数组，但是因为数组中没有元素，所以构建一个空数组。
2. 存了0个元素，len为0。
3. 能存0个元素，cap为0。

有些同学可能会问，为啥cap为0呢，我不是可以用append怼一些元素到s里面吗？那就是能存元素咯，怎么也不会是0吧。

如果执行：

```golang
s := []int{}
s[0] = 1 // panic
```

因为底层数组为空数组，操作空数组当然会panic，所以cap也是0。

能append的原因，是因为append中有对底层数组初始化/扩容的操作，所以执行append肯定是可以怼元素的。但是在append之前，是不能操作的。

## 2 new，make的不同

### 2.1 new

```golang
var s = new([]int)
```

|array|len|cap|
|----|----|----|
|nil|0|0|

1. new不会分配底层数组，array为nil。**new的返回值，就是slice结构的起始地址。所以这里的s并不是一个slice，而是一个地址。**
2. len = 0。
3. cap = 0。

如果此时使用append，就会给slice开辟底层数组：

```golang
var s = new([]int)
*s = append(*s, 1)
```

### 2.2 make

```golang
var s = make([]int,2,5)
```

|array|len|cap|
|----|----|----|
|[0,0,0,0,0]|2|5|

1. make会分配底层数组，array为长度为5的数组，并且认为已经填入了2个元素，且前2个元素都是int类型初始值，即为0。
2. 已经存了2个元素，占用了底层数组的0、1下标，len = 2。
3. 一共能存5个元素，cap = 5。

如果此时执行：

```golang
var s = make([]int,2,5)
s = append(s, 1)
```

|array|len|cap|
|----|----|----|
|[0,0,1,0,0]|3|5|

1. 在array的两个元素之后，刚好能依次append进去3个元素。如果再多append进去一个元素，比如4，就会发生扩容，根据一定规则去拓展array和cap，再插入元素。
2. len = 2 + 1 = 3。
3. 未发生扩容，cap = 5。

有些同学可能觉得，诶我这个array看起来像是有5个元素，我能不能这样操作：

```golang
var s = make([]int,2,5)
s = append(s, 1)
s[3] = 3
s[4] = 4
```

这样是不行的，因为底层数组看似有5位，但是有效位只有前3位，如果访问后两位会爆数组越界错误。

## 3 练习，画出一个string slice

画出以下代码的内存示意图：

```golang
var s = new([]string)
var str = "test"
*s = append(*s, str)
```

### 3.1 错误示例

1. s是一个指针，指向slice。
2. 定义了一个str变量。
2. str变量由指针和长度组成，指针指向内存中存的字节`test`，长度为4。
3. slice由array、len、cap组成，array存的是str。

![pic](https://brt-1303999354.cos.ap-shanghai.myqcloud.com/QQ%E6%88%AA%E5%9B%BE20210312161849.png)

为什么错了？因为通过append方法，slice的array存的是一个全新的string对象，不是外层的str变量。

在append进slice之后，slice中的对象，和外层的对象已经没有关系了。

### 3.2 正确示例

1. s是一个指针，指向slice。
2. 定义了一个str变量。
2. str变量由指针和长度组成，指针指向内存中存的字节`test`，长度为4。
3. slice由array、len、cap组成，append的时候，创建了一个新的string，内容和str一样，但是地址完全不同。array存的是这个新的string。

![pic](https://brt-1303999354.cos.ap-shanghai.myqcloud.com/QQ%E6%88%AA%E5%9B%BE20210312162828.png)

## 4 扩容机制

## 5 计算机如何分配内存

## 6 扩容例子