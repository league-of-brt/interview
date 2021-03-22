# 手撕LFU

* [手撕LFU](#手撕lfu)
  * [1 概述](#1-概述)
  * [2 图解](#2-图解)
    * [2\.1 Put方法具体实现](#21-put方法具体实现)
    * [2\.2 Get方法具体实现](#22-get方法具体实现)
  * [3 简单实现](#3-简单实现)
  * [4 一些缺点](#4-一些缺点)

有兴趣的同学可以先看看这题：
https://leetcode-cn.com/problems/lfu-cache/

## 1 概述

一般来说，无论是面试还是实际生产，好像LRU的应用都比LFU要多。但是实际上我个人认为LFU还比较难实现一点。

## 2 图解

可以用很多数据结构实现LFU（比如两个HashMap、两条链表啥的），而且这些设计和内存、空间使用率有很大关系。

但是我觉得好理解才是最重要的，先实现，然后再优化。所以我推荐使用两个HashMap和多个双向链表的组合。

![pic](https://brt-1303999354.cos.ap-shanghai.myqcloud.com/QQ%E6%88%AA%E5%9B%BE20210322164629.png)

只提供两个方法：

1. Get(key)
2. Put(key,value)

**对于LFU和LRU的实现，我们最需要关注的是节点的移动，也就是节点指针的指向。**

### 2.1 Put方法具体实现

首先注意我们使用了两个HashMap：

1. freq map，key=访问次数，value=双链表。双链表的头是更近的元素，优先淘汰链表尾部的元素。
2. value map，key=元素的key，value=元素的value。这个HashMap用于确认元素是否存在。

所以，每次操作都有必要操作这两个HashMap。

假设我们创建一个能够容纳两个元素的cache。

我们先加入一个元素key=1，value=1。将节点移动到次数为1的双链表的头部：

![pic](https://brt-1303999354.cos.ap-shanghai.myqcloud.com/QQ%E6%88%AA%E5%9B%BE20210322165607.png)

在加入一个元素key=2，value=2。由于新的节点最近被操作，所以要放在双链表的头部，在旧节点之前：

![pic](https://brt-1303999354.cos.ap-shanghai.myqcloud.com/QQ%E6%88%AA%E5%9B%BE20210322165514.png)

现在修改元素key=2，将其value设置为3。那么现在节点key=2被访问了两次，将其“升级”，放到下一个访问次数的双链表的头部。

而节点key=1也顶替了原来节点key=2的位置，相当于做了一个移动操作：

![pic](https://brt-1303999354.cos.ap-shanghai.myqcloud.com/QQ%E6%88%AA%E5%9B%BE20210322165919.png)

现在加入一个元素key=3，value=3，发现超出容量。这里有两个选择（注意同时干掉value map中的key=1）：

1. 先把节点key=3插到访问次数为1的双链表的头结点，再删除双链表的尾节点key=1。
2. 先删除访问次数为1的双链表的尾结点key=1，再插入新节点key=3。

不过效果都是一样的：

![pic](https://brt-1303999354.cos.ap-shanghai.myqcloud.com/QQ%E6%88%AA%E5%9B%BE20210322170239.png)

至此就是Put方法的全部功能。

### 2.2 Get方法具体实现

现在我们查询key=114514，在HashMap中找不到对应指针，那么返回一个-1。

如果查询key=2，那么需要再增加key=2的节点的访问次数，将其移动到下一访问次数双链表的表头：

![pic](https://brt-1303999354.cos.ap-shanghai.myqcloud.com/QQ%E6%88%AA%E5%9B%BE20210322175259.png)

## 3 简单实现

顺着上面的思路，简单翻译出来：

```golang
package cache

// LFUCache ...
type LFUCache struct {
	size     int
	capacity int
	maxFreq  int
	freqMap  map[int]*dLinkedNode
	valMap   map[int]*dLinkedNode
}

type dLinkedNode struct {
	key   int
	value int
	freq  int
	prev  *dLinkedNode
	next  *dLinkedNode
}

// Constructor ...
func Constructor(capacity int) LFUCache {
	return LFUCache{
		size:     0,
		capacity: capacity,
		freqMap:  make(map[int]*dLinkedNode),
		valMap:   make(map[int]*dLinkedNode),
	}
}

// Get ...
func (c *LFUCache) Get(key int) int {
	if c.capacity <= 0 {
		return -1
	}
	if node := c.valMap[key]; node == nil {
		return -1
	}

	node := c.valMap[key]

	if node.prev == nil {
		// 已经是头节点
		if node.next != nil {
			// 处理下一个节点
			node.next.prev = nil
			c.freqMap[node.freq] = node.next
		} else {
			// 没有下一个节点直接置为空
			c.freqMap[node.freq] = nil
		}

		node.freq++
		if node.freq > c.maxFreq {
			c.maxFreq = node.freq
		}

		if oldHead := c.freqMap[node.freq]; oldHead == nil {
			// 下一频率没有头节点
			node.next = nil
			node.prev = nil
			c.freqMap[node.freq] = node
		} else {
			node.next = oldHead
			node.prev = nil
			oldHead.prev = node
			c.freqMap[node.freq] = node
		}

		return node.value
	}

	// 不是头结点
	node.prev.next = node.next
	if node.next != nil {
		node.next.prev = node.prev
	}

	node.freq++
	if node.freq > c.maxFreq {
		c.maxFreq = node.freq
	}

	if oldHead := c.freqMap[node.freq]; oldHead == nil {
		// 下一频率没有头节点
		node.next = nil
		node.prev = nil
		c.freqMap[node.freq] = node
	} else {
		node.next = oldHead
		node.prev = nil
		oldHead.prev = node
		c.freqMap[node.freq] = node
	}

	return node.value
}

// Put ...
func (c *LFUCache) Put(key int, value int) {
	if c.capacity <= 0 {
		return
	}

	if node := c.valMap[key]; node == nil {
		c.size++
		if c.size > c.capacity {
			// 先干掉访问频率最少的那个节点
			c.size--
			for i := 1; i <= c.maxFreq; i++ {
				if node := c.freqMap[i]; node == nil {
					continue
				} else {
					// 找到最低的频率，找到最后的节点
					for node.next != nil {
						node = node.next
					}

					if node.prev == nil {
						// 如果是头节点，直接置为nil
						c.freqMap[node.freq] = nil
						node.prev = nil
					} else {
						// 直接删除节点
						node.prev.next = nil
						node.prev = nil
					}

					c.valMap[node.key] = nil
					break
				}
			}
		}

		// 加入新节点
		freq := 1
		node := &dLinkedNode{
			key:   key,
			value: value,
			freq:  freq,
			prev:  nil,
			next:  nil,
		}

		if oldHead := c.freqMap[freq]; oldHead != nil {
			oldHead.prev = node
			node.next = oldHead
		}

		c.valMap[key] = node
		c.freqMap[freq] = node

		if freq > c.maxFreq {
			c.maxFreq = freq
		}

		return
	}

	node := c.valMap[key]
	node.value = value

	if node.prev == nil {
		// 已经是头节点
		if node.next != nil {
			// 处理下一个节点
			node.next.prev = nil
			c.freqMap[node.freq] = node.next
		} else {
			// 没有下一个节点直接置为空
			c.freqMap[node.freq] = nil
		}
	} else {
		// 不是头结点
		node.prev.next = node.next
		if node.next != nil {
			node.next.prev = node.prev
		}
	}

	node.freq++
	if node.freq > c.maxFreq {
		c.maxFreq = node.freq
	}

	if oldHead := c.freqMap[node.freq]; oldHead == nil {
		// 下一频率没有头节点
		node.next = nil
		node.prev = nil
		c.freqMap[node.freq] = node
	} else {
		node.next = oldHead
		node.prev = nil
		oldHead.prev = node
		c.freqMap[node.freq] = node
	}
}
```

## 4 一些缺点

这种简单的LFU，不经过优化是不好用的，简单举两个例子：

> 1. 元素A是个热点信息，被访问了超过100w次。访问次数实在是太多，如果现在我们想把这个元素赶出去，没有办法通过正常的访问次数机制将其淘汰。可能借助一些手段或者一些什么新的策略。
>
> 2. 现在新加入一个元素B，LFU容量已经满了。再加入一个元素C，对于之前被访问很多次的老油条，刚加入的元素B显然有更高几率会被赶出去。这就容易导致新的热点数据还没来得及堆积起访问次数就频繁被淘汰。

总之我认为，LFU不适合用于做热度的场景，比较适合做历史累计的场景。
