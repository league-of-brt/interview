# 手撕LRU

有兴趣的同学可以先看看这题：
https://leetcode-cn.com/problems/lru-cache/

## 1 概述

LRU，Least Recent Used，最近最少使用，只要是为了节省储存空间、淘汰掉最近最少使用的数据的算法，可以统称为LRU。

LRU也有很多种实现方法，最主要的区别的是淘汰数据的策略。

我们写的实验性的LRU，和生产环境真实使用的LRU区别非常大，主要体现于对热点数据的保护，和对性能的追求。

再牛逼的LRU，也是从基本形态一点点改进过来的。我们可以先把简单的实现了，然后理解下Mysql和Redis对LRU的改进，最后再实现一个改进版的LRU。

## 2 图解

总之先记住一点，一个简单的LRUCache，是由一个HashMap和一个双向链表组成的。

![pic](https://brt-1303999354.cos.ap-shanghai.myqcloud.com/v2-09f037608b1b2de70b52d1312ef3b307_720w.png)

只提供两个方法：

1. Get(key)
2. Put(key,value)

### 2.1 Put方法具体实现：

当我们插入一个键值对时，实际上我们是创建了一个node。在HashMap的key保存的是node的指针，在双向链表中存的是node本体。

所谓的Put操作，其实就是新建一个node，然后调整它在HashMap和双链表中的位置，仅此而已。

我们整一个容量为2的LRUCache，整个结构的初始状态是这样：

![pic](https://brt-1303999354.cos.ap-shanghai.myqcloud.com/QQ%E6%88%AA%E5%9B%BE20210305155810.png)

> 这里要非常注意一点，双向链表使用了两个虚拟节点head和tail，这两个节点是固定的，也就是说双向链表的头尾一直都是这两个节点，这是为了编程方便。

当我们插入一个键值对，比如`1:1`，首先在HashMap中查询。如果没有那就直接插入，就会变成这样：

![pic](https://brt-1303999354.cos.ap-shanghai.myqcloud.com/QQ%E6%88%AA%E5%9B%BE20210305160320.png)

当我们再插入一个键值对，比如`2:2`，还是一样在HashMap中查询。没有就是插入，那就把最新的数据插入到head节点之后（最新的数据最靠近head的原则），就会变成这样：

![pic](https://brt-1303999354.cos.ap-shanghai.myqcloud.com/QQ%E6%88%AA%E5%9B%BE20210305160803.png)

同理，现在插入一个`3:3`，变成了这样：

![pic](https://brt-1303999354.cos.ap-shanghai.myqcloud.com/QQ%E6%88%AA%E5%9B%BE20210305161355.png)

我们发现不对劲了，容量限制为2，但是现在已经有3个节点了，那么需要干掉最久未使用的数据，即最靠近tail的节点。于是把tail节点之前的`1:1`干掉：

![pic](https://brt-1303999354.cos.ap-shanghai.myqcloud.com/QQ%E6%88%AA%E5%9B%BE20210305161555.png)

最后我们尝试更新一个已有的键值对，尝试插入`2:4`，先在HashMap中查询，能通过key=2找到指针，于是修改节点`2:2`的值为`2:4`。根据最新的数据最靠近head的原则，把`2:4`节点移到head的后面：

![pic](https://brt-1303999354.cos.ap-shanghai.myqcloud.com/QQ%E6%88%AA%E5%9B%BE20210305162025.png)

至此就是Put方法的全部功能。

### 2.2 Get方法具体实现：

现在我们查询key=114514，在HashMap中找不到对应指针，那么返回一个nil。

如果查询key=3，能找到`3:3`，那么就把节点的值value=3返回。

除此之外，根据最新的数据最靠近head的原则，需要把`3:3`移动到head的后面，变成这样：

![pic](https://brt-1303999354.cos.ap-shanghai.myqcloud.com/QQ%E6%88%AA%E5%9B%BE20210305162936.png)

至此就是Get方法的全部功能。

## 3 简单实现

基本实现：https://github.com/xie4ever/practice/tree/master/golang/leetcode/lru/example1

可以根据2的图解先自己实现一个，我实现的初版如下：

```golang
package cache

// LRUCache ...
type LRUCache struct {
	size     int
	capacity int
	cacheMap map[int]*dLinkedNode
	head     *dLinkedNode
	tail     *dLinkedNode
}

type dLinkedNode struct {
	key   int
	value int
	prev  *dLinkedNode
	next  *dLinkedNode
}

// Constructor ...
func Constructor(capacity int) LRUCache {
	head := &dLinkedNode{
		key:   0,
		value: 0,
	}
	tail := &dLinkedNode{
		key:   0,
		value: 0,
	}
	head.next = tail
	tail.prev = head

	return LRUCache{
		size:     0,
		capacity: capacity,
		cacheMap: map[int]*dLinkedNode{},
		head:     head,
		tail:     tail,
	}
}

// Get ...
func (c *LRUCache) Get(key int) int {
	if node, ok := c.cacheMap[key]; !ok {
		return -1
	} else {
		node.prev.next = node.next
		node.next.prev = node.prev

		node.prev = c.head
		node.next = c.head.next

		c.head.next.prev = node
		c.head.next = node

		return node.value
	}
}

// Put ...
func (c *LRUCache) Put(key int, value int) {
	if node, ok := c.cacheMap[key]; !ok {
		deleteNode := c.tail.prev
		node := &dLinkedNode{
			key:   key,
			value: value,
			prev:  c.head,
			next:  c.head.next,
		}
		c.cacheMap[key] = node

		c.head.next.prev = node
		c.head.next = node

		c.size++
		if c.size > c.capacity {
			c.size--
			delete(c.cacheMap, deleteNode.key)
			c.tail.prev.prev.next = c.tail
			c.tail.prev = c.tail.prev.prev
		}
		return
	} else {
		node.value = value

		node.prev.next = node.next
		node.next.prev = node.prev

		node.prev = c.head
		node.next = c.head.next

		c.head.next.prev = node
		c.head.next = node
	}
}
```

写一个测试用例：

```golang
package cache

import (
	"fmt"
	"testing"
)

// Test ...
func Test(t *testing.T) {
	c := Constructor(2)
	c.Put(1, 1)
	c.Put(2, 2)
	fmt.Println(c.Get(1))
	c.Put(3, 3)
	fmt.Println(c.Get(2))
	c.Put(4, 4)
	fmt.Println(c.Get(1))
	fmt.Println(c.Get(3))
	fmt.Println(c.Get(4))
}

// Test2 ...
func Test2(t *testing.T) {
	c := Constructor(2)
	c.Put(1, 0)
	c.Put(2, 2)
	fmt.Println(c.Get(1))
	c.Put(3, 3)
	fmt.Println(c.Get(2))
	c.Put(4, 4)
	fmt.Println(c.Get(1))
	fmt.Println(c.Get(3))
	fmt.Println(c.Get(4))
}
```

## 4 改进

### 4.1 抽取通用方法，简化节点操作

简化节点操作：

https://github.com/xie4ever/practice/tree/master/golang/leetcode/lru/example2

### 4.2 支持interface

初版使用int，真实环境的map都是支持interface，优化参数：

https://github.com/xie4ever/practice/tree/master/golang/leetcode/lru/example3

### 4.3 加锁支持并发安全

因为涉及双链表操作有大量nil，目前开个并发就会panic，需要加锁：

https://github.com/xie4ever/practice/tree/master/golang/leetcode/lru/example4

### 4.4 更换线程安全的HashMap

如果使用golang，原则上只要操作加了锁，就算使用普通map也是线程安全的。但是在java中不是这样，出于安全考虑换用Sync.Map。

https://github.com/xie4ever/practice/tree/master/golang/leetcode/lru/example5

### 4.5 增加过期时间，读取时懒删除

实现过期机制和懒删除，缺点是如果直接不调用Get方法，直接访问内部结构还是会拿到过期数据。

https://github.com/xie4ever/practice/tree/master/golang/leetcode/lru/example6

### 4.6 增加清理协程，主动删除过期节点

开协程清理过期节点：

https://github.com/xie4ever/practice/tree/master/golang/leetcode/lru/example7

### 4.7 清理协程添加到构建方法中

清理协程应该和LRUCache的构筑方法绑定：

https://github.com/xie4ever/practice/tree/master/golang/leetcode/lru/example8

## 5 生产环境使用的LRU

最经典的就是MySQL InnoDB的LRU：

具体可以参考这一篇文章：https://zhuanlan.zhihu.com/p/142087506

![pic](https://brt-1303999354.cos.ap-shanghai.myqcloud.com/QQ%E6%88%AA%E5%9B%BE20210305165706.png)

如果使用我们简单的LRU，可能会出现什么情况？

> 在MySQL中经常会出现全表扫描，一种是开发人员对索引的使用不当导致的，一种是业务如此，无法避免。
>
> 当出现全表扫描时，InnoDB会将该表的数据页全部从磁盘文件加载进缓存页中，这些缓存页会被加入到LRU链表中。如果进行全表扫描的对象是一张非常大的表，可能是几十 GB 的数据，而且这张表记录的是类似于账户流水、操作日志等使用不频繁的数据，这个时候如果 LRU 链表已经满了，现在我们就要淘汰一部分缓存页，腾出空间来存放全表扫描出来的数据。这样就会因为全表扫描的数据量大，需要淘汰的缓存页多，导致在淘汰的过程中，极有可能将需要频繁使用到的缓存页给淘汰了，而放进来的新数据却是使用频率很低的数据，甚至是这一次使用之后，后面几乎再也不用，如操作日志等。
>
> 最终导致的现象就是，当我们在对这些使用不频繁的大表进行全表扫描之后，在一段时间内，Buffer Pool 缓存的命中率明显下降，SQL 的性能也明显下降，因为常用的缓存页被淘汰了，再进行查询时，需要从重新磁盘读取，发生磁盘 IO，性能下降。所以，如果 MySQL 只是简单的使用 LRU 算法，那么碰到全表扫描时，就会存在性能下降的问题，甚至在高并发场景下，成为性能瓶颈。
>
> 除此之外，MySQL的预读机制也会替换LRU中的热点数据，导致缓存失效。

对此，MySQL优化了LRU的结构，使用冷热分离，尽量保存住热点数据，减少大量数据操作对热点数据的影响。

## 6 Redis使用的LRU

如果面试问到这个就坑了，因为Redis使用的LRU很有特点，是基于时钟的LRU实现的，这个不专门去研究一下还真搞不懂。

简单举个例子，我刚才实现的例子，过期策略是基于本机时间。但是Redis是有集群的，可能每台物理机的时间都不一致，怎么实现这么多机器的数据同时过期？

这里我猜测是使用了时间中心，具体怎么实现的还要看。

先去看下官方文档：
https://redis.io/topics/lru-cache

然后可以看一下这个：
http://haoran.tech/2018/07/19/4-Redis%E8%BF%87%E6%9C%9F%E7%AD%96%E7%95%A5-%E6%89%8B%E5%86%99LRU/

（我觉得这个问题需要大家讨论一下）

