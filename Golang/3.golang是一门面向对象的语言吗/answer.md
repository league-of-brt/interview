# golang是面向对象的语言吗？

* [golang是面向对象的语言吗？](#golang是面向对象的语言吗)
  * [1 Yes and No](#1-yes-and-no)
  * [2 官方说法](#2-官方说法)
  * [3 深入理解](#3-深入理解)

参考：https://www.zhihu.com/question/315995798/answer/630360043

Java有三大面向对象核心特征：
1. 封装。
2. 继承。
3. 多态。

golang虽然可以通过各种方式达到类似的功能，但golan 算是面向对象的语言吗？

## 1 Yes and No

> 这个咱就看go官网的说辞就好了，他们自己说是Yes and No。
>
> 明显go是允许OO的编程风格的，但又缺乏一些Java和C++中的常见类型继承结构。Go的interface也和Java中的用法大相径庭, 这也是我经常吹捧的隐式继承。
>
> Go自己觉得这一套挺好的，更加的容易使用且通用性更强。很多时候我们用OO的思想来组织我们的项目，但需要注意的是java的继承关系是一种非常强的耦合，有时候会给以后的升级带来麻烦（这个可以看为什么java8的升级带来了interface的default method）。我觉得隐式继承在超大型的monorepo项目中是非常有帮助的，当然小型的项目可能好处不是很明显。

## 2 官方说法

https://golang.org/doc/faq#Is_Go_an_object-oriented_language

> Yes and no. Although Go has types and methods and allows an object-oriented style of programming, there is no type hierarchy. The concept of “interface” in Go provides a different approach that we believe is easy to use and in some ways more general. There are also ways to embed types in other types to provide something analogous—but not identical—to subclassing. Moreover, methods in Go are more general than in C++ or Java: they can be defined for any sort of data, even built-in types such as plain, “unboxed” integers. They are not restricted to structs (classes).
>
> Also, the lack of a type hierarchy makes “objects” in Go feel much more lightweight than in languages such as C++ or Java.

## 3 深入理解

https://www.zhihu.com/question/315995798/answer/633618997

大家应该深入理解下面向对象的编程思想，思考下面向对象如何设计、如何解决问题。

如果不用面向对象，又该怎么设计你当前的功能，你能分得清楚两者的区别吗？