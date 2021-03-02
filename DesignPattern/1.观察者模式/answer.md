# 观察者模式

* [观察者模式](#观察者模式)
  * [1 事件驱动](#1-事件驱动)
  * [2 观察者模式](#2-观察者模式)
  * [3 坏实现](#3-坏实现)
  * [4 好实现](#4-好实现)
  * [5 设计问题](#5-设计问题)
  * [6 思考题](#6-思考题)
  * [7 坏实现](#7-坏实现)
  * [8 好实现](#8-好实现)

参考我以前写的文章：

1. https://xie4ever.com/2018/04/10/%e4%bd%bf%e7%94%a8%e4%ba%8b%e4%bb%b6%e9%a9%b1%e5%8a%a8%e7%9a%84%e7%ae%80%e5%8d%95%e4%be%8b%e5%ad%90%ef%bc%88%e8%a7%82%e5%af%9f%e8%80%85%e6%a8%a1%e5%bc%8f%ef%bc%89/
2. https://www.xie4ever.com/2018/04/12/%e8%a7%82%e5%af%9f%e8%80%85%e6%a8%a1%e5%bc%8f%e7%bb%83%e4%b9%a0%ef%bc%9a%e5%8d%87%e7%ba%a7%e4%ba%8b%e4%bb%b6/

原本我是用java写的，现在转用golang实现一遍，顺便复习下面向对象（go）。

> 假设你正在开发一个博客系统。现在有这样一个需求：“当作者创建（写完并保存）、修改、删除了一篇文章后，把文章推上最新发布榜，并且增加作者的积分”。你会怎么做呢？

实际上，这就是“达到什么什么条件，就做什么什么”的场景。完全可以用设计模式来解耦。

## 1 事件驱动

首先我们了解下事件驱动：

> 1. 实体。对实体的操作将会发出事件。在上面的场景中，实体就是“文章”。对实体进行操作的时候，我们会抛出一个相应的事件对象。（举个例子：创建新文章时，抛出一个文章事件。）
> 2. 事件。事件也是一个对象，事件对象应该“说明”实体发生了什么操作，并且封装了处理器类需要的所有信息，供处理器类使用。（举个例子：创建新文章时，抛出一个文章事件。事件类型是“新建文章”，封装了文章对象本身、创建日期等等信息。）
> 3. 事件监听器。事件监听器储存（管理）了所有的事件处理器，并且负责监听相应事件的发生。一旦监听到了某个事件，监听器就会告知所有的事件处理器，对事件进行处理。（举个例子：事件监听器A管理着B、C、D三个事件处理器。当事件A收到了事件E时，就会告知B、C、D对事件E进行相应的处理。）
> 4. 事件处理器。事件处理器负责处理事件。（举个例子：排名事件处理器收到了“新建文章”这个事件后，就会发进行“把文章推上最新发布榜”的处理。）

其实就是把实体直接干一个部分，拆分成实体、事件、事件监听器、事件处理器四个部分，实现解耦。

## 2 观察者模式

> 如果只从字面上去理解，“观察者”一直呆在一边持续观察，就像篮球场上的裁判一样，一直观察着球员的动作；一旦球员犯规了，裁判就会出示一张红牌。在这个过程中，球员只需要专注比赛，裁判会自发地完成所有的观察工作。

> 但是，如果要用语言实现这个场景（实现“持续”这个效果），就没这么优雅了。根据面向对象的设计原则，裁判的“观察”应该写成一个observe方法，“持续观察”就应该在observe方法里写一个while或者for循环。最后，裁判的表现就变成了：“我观察观察观察观察…”。暂且不论这个场景有多么滑稽，从性能上来说，这种做法就是非常低效的。

（顺带一提，这样是不是就变成了生产者消费者模式？）

> 如果理解不了上面的场景，不妨把上面的裁判想像成一个非常关心你的老妈。从现在开始，老妈每秒钟都问你一次“你吃了吗？”，如果发现你吃了，她才去收拾碗筷。这样一来，如果你一个小时内没吃饭，老妈就要不停询问3600次…作为一个孝顺儿子，我们当然不能让这种状况发生。

个人认为，观察者模式就是把“主动观察”变成“主动告知”，要求被观察者自觉地通知观察者。拿上面的篮球裁判做个例子，只要球员一犯规就主动通知裁判（发送给裁判犯规事件），裁判就不用一直盯着球员看了（不需要轮询），直接给出一张红牌即可（处理事件）。

再拿上面的老妈作个例子，现在老妈再也不问你“你吃了吗？”（不需要轮序），而是等你自己吃完，告诉她“我吃完了”（发送给老妈吃完事件），老妈才会去收拾碗筷（处理事件）。

## 3 坏实现

github：https://github.com/xie4ever/design-pattern/tree/master/observer/badexample/example1

```golang
package article

import (
	"log"
	"time"
)

// Article ...
type Article struct {
	ID      int64  `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

// NewArticle ...
func NewArticle(title, content string) Article {
	return Article{
		ID:      time.Now().Unix(),
		Title:   title,
		Content: content,
	}
}

// Add ...
func (a Article) Add() error {
	log.Print("do something")
	log.Print("rank change")
	log.Print("point change")
	return nil
}

// Delete ...
func (a Article) Delete() error {
	log.Print("do something")
	log.Print("rank change")
	log.Print("point change")
	return nil
}

// Modify ...
func (a Article) Modify() error {
	log.Print("do something")
	log.Print("rank change")
	log.Print("point change")
	return nil
}
```

实际上，只有创建、修改、删除这几个功能和方法是相关的，而“对作者进行一次排名”、“增加作者的积分”都是附加的功能，如果都写在同一个方法里，就会显得很臃肿。

那么有些同学就会说了，我在创建方法里面开个事务单独处理文章，然后别的附加功能全都走消息队列，就一行发送消息代码，怎么就臃肿了？

这是因为我只模拟了2个功能，如果现在有50个附加功能，那么随随便便就能多出上50行代码。每个方法都随便多出50行代码，整个文件就没什么可读性了。

## 4 好实现

github：https://github.com/xie4ever/design-pattern/blob/master/observer/goodexample/example1

关键是通过event、observer、processer进行解耦，这里只贴一块核心代码：

``` golang
package article

import "sync"

var (
	obs observer
)

type observer struct {
	ProcessorMap sync.Map `json:"processor_map"`
}

// GetObs ...
func GetObs() *observer {
	return &obs
}

// AddProcessor ...
func (o *observer) AddProcessor(p processor) {
	o.ProcessorMap.Store(p.GetID(), p)
}

// DeleteProcessor ...
func (o *observer) DeleteProcessor(id int64) {
	o.ProcessorMap.Delete(id)
}

// PostEvent ...
func (o *observer) PostEvent(e Event) error {
	if e.ID == 0 {
		return nil
	}
	if e.Type == 0 {
		return nil
	}

	o.ProcessorMap.Range(
		func(k, v interface{}) bool {
			p := v.(processor)
			if e.Type == TypeEventAdd {
				if err := p.EntryAdded(e); err != nil {
					return false
				}
			}
			if e.Type == TypeEventDelete {
				if err := p.EntryDeleted(e); err != nil {
					return false
				}
			}
			if e.Type == TypeEventModify {
				if err := p.EntryModified(e); err != nil {
					return false
				}
			}
			return true
		})

	return nil
}
```

## 5 设计问题

上面的代码确实实现了解耦，但是实现得不够优雅。

为什么？我个人觉得处理器接口设计得不好：

```golang
type processor interface {
	GetID() int64
	EntryAdded(e Event) error
	EntryDeleted(e Event) error
	EntryModified(e Event) error
}
```

这样是对文章实体的操作实现全面监控。这个例子中，只有3个操作文章的方法，所以可以这样写。但是，假如有100个操作方法，你要怎么办...写100个接口方法吗？

其实我个人觉得，这样设计就可以了：

```golang
type processor interface {
	GetType() int
	DoSomething(e Event) error
}
```

切记，不要在一个处理器里面塞太多方法，要像消息队列那样去设计！

## 6 思考题

这里给出第二题：

> 玩家打怪升级后，立刻修改玩家等级排名、给予玩家升级奖励、公告栏展示玩家信息。怎么实现？

切记，要像消息队列一样去设计！反正就是要把处理器给拆细！每种变化就一个处理器，然后对应好事件枚举就可以了，就像发布订阅一样。

## 7 坏实现

github：https://github.com/xie4ever/design-pattern/blob/master/observer/badexample/example2

反正就是堆到一起：

```golang
package player

import (
	"log"
	"time"
)

// Player ...
type Player struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// NewPlayer ...
func NewPlayer(name string) Player {
	return Player{
		ID:   time.Now().Unix(),
		Name: name,
	}
}

// Attack ...
func (p Player) Attack() error {
	log.Print("attack somebody")
	log.Print("level up")
	log.Print("level up reward change")
	log.Print("level rank change")
	log.Print("level up announce")
	return nil
}
```

## 8 好实现

github：https://github.com/xie4ever/design-pattern/blob/master/observer/goodexample/example2

这个解耦我觉得写得不错，应该是最佳实践了：

``` golang
package player

import (
	"sync"
)

var (
	obs observer
)

type observer struct {
	ProcessorMap sync.Map `json:"processor_map"`
}

// GetObs ...
func GetObs() *observer {
	return &obs
}

// AddProcessor ...
func (o *observer) AddProcessor(p processor) {
	v, ok := o.ProcessorMap.Load(p.GetType())
	if ok {
		pList := v.([]processor)
		pList = append(pList, p)
		o.ProcessorMap.Store(p.GetType(), pList)
	} else {
		o.ProcessorMap.Store(p.GetType(), []processor{p})
	}
}

// DeleteProcessorByType ...
func (o *observer) DeleteProcessorByType(t int) {
	o.ProcessorMap.Delete(t)
}

// PostEvent ...
func (o *observer) PostEvent(e Event) error {
	if e.ID == 0 {
		return nil
	}
	v, ok := o.ProcessorMap.Load(e.Type)
	if ok {
		pList := v.([]processor)
		for _, p := range pList {
			if err := p.DoSomething(e); err != nil {
				return err
			}
		}
	}
	return nil
}
```
