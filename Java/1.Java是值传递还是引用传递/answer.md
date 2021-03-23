# Java是值传递还是引用传递

* [Java是值传递还是引用传递](#java是值传递还是引用传递)
  * [1 常见错误](#1-常见错误)
    * [1\.1 传递基本类型](#11-传递基本类型)
    * [1\.2 传递“没有提供改变自身方法的引用类型”](#12-传递没有提供改变自身方法的引用类型)
    * [1\.3 提供了改变自身方法的引用类型](#13-提供了改变自身方法的引用类型)
    * [1\.4 提供了改变自身方法的引用类型，但是不使用，而是使用赋值运算符](#14-提供了改变自身方法的引用类型但是不使用而是使用赋值运算符)
  * [2 实例分析](#2-实例分析)
    * [2\.1 局部变量/方法参数](#21-局部变量方法参数)
    * [2\.2 数组类型引用和对象](#22-数组类型引用和对象)
    * [2\.3 多维数组](#23-多维数组)
    * [2\.4 String](#24-string)
  * [3 一些例子](#3-一些例子)

首先说一下前提，纠结Java是值传递还是引用传递是没有意义的。重要的是我们要知道这些个例子，内存的情况是怎么样的，为啥有的会改变，有的又不会改变。

来源于：https://www.zhihu.com/question/20628016/answer/28970414 ，目前最准确的解释是：**参数藉由值传递方式，传递的值是个引用。（句中两个“值”不是一个意思，第一个值是evaluation result，第二个值是value content）**

这个标准答案对我来说很难理解，因为是从各种概念的角度提出的。根据Java中各种现象，我的理解是：Java是值传递，并且有以下特点：

> 1. 如果参数是基本类型，传递的是基本类型的字面量值的拷贝。
>
> 2. 如果参数是引用类型，传递的是该参量所引用的对象在堆中地址值的拷贝。

## 1 常见错误

一定要知道怎么分析内存。

### 1.1 传递基本类型

这个是最经典的例子了，大家都知道答案是10，但是执行过程中发生了什么呢？

```java
int num = 10;
void foo(int value) {
    value = 100;
}
foo(num);
System.out.println(num); // 100 × 10 √
```

![pic](https://brt-1303999354.cos.ap-shanghai.myqcloud.com/QQ%E6%88%AA%E5%9B%BE20210323124854.png)

从这个例子就能看出java的基本类型，肯定不是引用传递了。

### 1.2 传递“没有提供改变自身方法的引用类型”

也是老经典了，关键在于声明s的方式，在虚拟机中的实现不同。别的都一样。

```java
String s = new String("happy");
void foo(String value) {
    value = "sad";
}
foo(s);
System.out.println(s) // sad × happy √
```

![pic](https://brt-1303999354.cos.ap-shanghai.myqcloud.com/QQ%E6%88%AA%E5%9B%BE20210323125510.png)

```java
String s = "happy";
void foo(String value) {
    value = "sad";
}
foo(s);
System.out.println(s); // sad × happy √
```

![pic](https://brt-1303999354.cos.ap-shanghai.myqcloud.com/QQ%E6%88%AA%E5%9B%BE20210323124520.png)

从这个例子可以看出，java的对象肯定也不是引用传递。

### 1.3 提供了改变自身方法的引用类型

因为以值传递的形式传递地址，所以变量s和形参builder都指向堆中同一个StringBuilder对象，所以修改的是同一个对象。

```java
StringBuilder s = new StringBuilder("My");
void foo(StringBuilder builder) {
    builder.append("SQL");
}
foo(s);
System.out.println(s.toString()); // My × MySQL √
```

![pic](https://brt-1303999354.cos.ap-shanghai.myqcloud.com/QQ%E6%88%AA%E5%9B%BE20210323130021.png)

有些同学就是被这一步搞晕的。因为从效果上看，确实builder对象发生了改变，像是引用传递。于是他们得出一个错误结论：“Java的基本类型是值传递，对象是引用传递”。

但是实际上传递的是不是builder对象，而是指向builder对象的地址。修改地址指向的对象，当然会导致形参实参一起变化。所以这里依然是值传递。

### 1.4 提供了改变自身方法的引用类型，但是不使用，而是使用赋值运算符

```java
StringBuilder s = new StringBuilder("My");
void foo(StringBuilder builder) {
    builder = new StringBuilder("SQL");
}
foo(s);
System.out.println(s.toString()); // MySQL × My √
```

![pic](https://brt-1303999354.cos.ap-shanghai.myqcloud.com/QQ%E6%88%AA%E5%9B%BE20210323131443.png)

如果Java是引用传递，那么执行 `builder = new StringBuilder("SQL")` 之后，实参s肯定会被重新赋值。实际上是形参builder指向了new出来的新的StringBuilder对象，而实参s依旧指向之前的StringBuilder对象。

## 2 实例分析

### 2.1 局部变量/方法参数

局部变量和方法参数在JVM中的储存方法是相同的，都是在栈上开辟空间来储存的，随着进入方法开辟，退出方法回收。

以32位JVM为例，boolean/byte/short/char/int/float以及引用都是分配 4 字节空间，long/double分配 8 字节空间。对于每个方法来说，最多占用多少空间是一定的，这在编译时就可以计算好。

我们都知道JVM内存模型中有 Stack 和 Heap 的存在，但是更准确的说，是每个线程都分配一个独享的 Stack，所有线程共享一个 Heap。对于每个方法的局部变量来说，是绝对无法被其他方法，甚至其他线程的同一方法所访问到的，更遑论修改。

当我们在方法中声明一个 `int i = 0`，或者 `Object obj = null` 时，仅仅涉及 Stack，不影响到 Heap。当我们 `new Object()` 时，会在 Heap 中开辟一段内存并初始化Object对象。当我们将这个对象赋予obj变量时，仅仅是stack中代表obj的那4个字节变更为这个对象的地址（其实就是引用的指向改变了）。

### 2.2 数组类型引用和对象

数组涉及一堆指针的东西，这里要特别注意内存结构。

当我们声明一个数组时，如 `int[] arr = new int[]{0, 0}`，因为数组也是对象，arr实际上是引用，Stack 上仅仅占用4字节空间，new int[2]会在 Heap 中开辟一个数组对象，然后arr指向它。

于是有些同学会想象出这样一张图：

![pic](https://brt-1303999354.cos.ap-shanghai.myqcloud.com/QQ%E6%88%AA%E5%9B%BE20210323134007.png)

但是，如果内存结构真的是这样的，那么就会有这样一个例子：

```java
int[] arr = new int[]{0, 0};
int a = arr[1];
arr[1] = 1;
System.out.println(a); // 1 × 0 √
```

然后大家就晕了，我这里明明通过 `arr[1] = 1` 改变了数组中的值，而 `a = arr[1]` 明显又是指向这个值的，为什么不是输出1呢？

其实是内存结构理解错了：

![pic](https://brt-1303999354.cos.ap-shanghai.myqcloud.com/QQ%E6%88%AA%E5%9B%BE20210323135139.png)

数组中存的不是值，而是对内容的引用，也就是说 `arr[1]` 是一个引用。

1. 执行 `int a = arr[1]` 时，变量a指向 `arr[1]` 指向的堆中变量0。
2. 执行 `arr[1] = 1`，实际上是 `arr[1]` 指向了堆中新的变量1。
3. 变量a的指向永远是这个0，所以输出只会为0。

### 2.3 多维数组

当我们声明一个二维数组时，如 `int[][] arr2 = new int[2][2]`，arr2同样仅在stack中占用4个字节，会在内存中开辟一个长度为2的，类型为int[]的数组，然后arr2指向这个数组。这个数组内部有两个引用（大小为4字节），分别指向两个长度为2的类型为int的数组。

![pic](https://brt-1303999354.cos.ap-shanghai.myqcloud.com/QQ%E6%88%AA%E5%9B%BE20210323140529.png)

再举个例子：

```java
int[][] arr2 = new int[2][2];

int[] tmp0 = arr2[0];
int[] tmp1 = arr2[1];

arr2[0] = new int[]{1};
arr2[1] = new int[]{2, 3, 4};
```

![pic](https://brt-1303999354.cos.ap-shanghai.myqcloud.com/QQ%E6%88%AA%E5%9B%BE20210323144503.png)

再举个例子：

```java
int[][] arr2 = new int[2][2];

int[] tmp0 = new int[]{1};
int[] tmp1 = new int[]{2, 3, 4};

arr2[0] = tmp0;
arr2[1] = tmp1;

int[][] arr3 = new int[3][];
arr3[0] = arr2[0];
arr3[1] = arr2[1];
```

![pic](https://brt-1303999354.cos.ap-shanghai.myqcloud.com/QQ%E6%88%AA%E5%9B%BE20210323145544.png)

### 2.4 String

```java
String s = new String("happy");
```

有很多同学把String想得太过于简单了，其实String的结构还算有点复杂：

![pic](https://brt-1303999354.cos.ap-shanghai.myqcloud.com/QQ%E6%88%AA%E5%9B%BE20210323150257.png)

JVM的版本、String的定义方式，都会导致String的内存结构不同。

## 3 一些例子

```java
public void tricky(Point arg1, Point arg2) {
    arg1.x = 100;
    arg1.y = 100;
    Point temp = arg1;
    arg1 = arg2;
    arg2 = temp;

    System.out.println("Inside func arg1 x: " + arg1.x + ", y: " + arg1.y); // 2 3
    System.out.println("Inside func arg2 x: " + arg2.x + ", y: " + arg2.y); // 100 100
}

public static void main(String[] args) {
    Point p1 = new Point(2, 3);
    Point p2 = new Point(2, 3);
    System.out.println("p1 x: " + p1.x + ", y: " + p1.y); // 2 3
    System.out.println("p2 x: " + p2.x + ", y: " + p2.y); // 2 3

    tricky(p1, p2); 

    System.out.println("p1 x: " + p1.x + ", y: " + p1.y); // 100 100
    System.out.println("p2 x: " + p2.x + ", y: " + p2.y); // 2 3
}
```

建议做一下这道题，如果能做出来，就算不懂什么值传递引用传递，也不会轻易用错。

最重要理解这三步：

```java
Point p1 = new Point(2, 3);
Point p2 = new Point(2, 3);
public void tricky(Point arg1, Point arg2) {
    arg1.x = 100;
    arg1.y = 100;
}
```
![pic](https://brt-1303999354.cos.ap-shanghai.myqcloud.com/QQ%E6%88%AA%E5%9B%BE20210323153136.png)

```java
Point temp = arg1;
arg1 = arg2;
arg2 = temp;
```

![pic](https://brt-1303999354.cos.ap-shanghai.myqcloud.com/QQ%E6%88%AA%E5%9B%BE20210323153038.png)

跳出tricky方法，形参和临时变量全部被回收。到了最外层main方法后，只剩下：

![pic](https://brt-1303999354.cos.ap-shanghai.myqcloud.com/QQ%E6%88%AA%E5%9B%BE20210323152949.png)

一定要动手画画图。
