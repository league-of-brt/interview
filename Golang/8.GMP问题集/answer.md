# GMP问题集

## 1 为什么需要 P 这个组件，直接把 runqueues 放到 M 不行吗？

1. 在没有P的前提下，如果M直接去绑定G，一旦M发生阻塞，为了跑剩下的G，就要在等待队列找现成的M，如果没有现成的M就要开M。这样，一旦G很多并且阻塞很多，就会开一大堆M，加重线程切换的负担，并不能提升并发效率。
2. 一旦等待队列和M捆绑，这样就导致每个M都带有等待队列，阻塞的M的等待队列为空，浪费资源。
3. 系统调用阻塞的时候只阻塞一个g，p可以带着剩下的g到别的M执行，提高效率。

## 2 gpm到底是什么

go的协程调度器，使用gmp模型。该模型的意义是减少阻塞对并发的影响，让cpu每时每刻都有事干，榨干cpu的性能。

## 3 scheduler是如何调度的

有以下几个特点：
1. work stealing
2. hand off
3. 利用多核
4. 抢占
5. 全局G队列

## 4 什么时候会触发调度

1. 执行go关键字，就是新起goroutine。
2. GC，虚拟机判断该执行GC时，会去找空闲的M。
3. 发生系统调用，当前G会和当前M绑定，P会找到另外的M，继续执行剩下的G。
4. 发生阻塞，同上。

## 5 当在M上运行的goroutine发生系统调用时，会怎么工作

1. 同步，当前G会和当前M绑定，P会找到另外的M，继续执行剩下的G。
2. 异步，goroutine会被network poller接手，M会继续执行LRQ其他的runnable的goroutines，goroutine在执行完系统调用后，依然会回到P的LRQ（？）

## 6 为什么每个m都会对应一个g0 (g0是用于调度每个线程中的goroutine，包括gc等等，拥有比较大的栈内存)

每个P都会自带一个原生的g0，g0是用来协调调度的协程。

比如一个M刚执行完一个G，此时会调度P上的g0，g0负责找到下一个执行的G，交给M执行。

## 7 什么时候会抢占P

个人理解：抢占的意思就是协程之间能不能抢占时间片，如果必须执行完释放，那就是非抢占式。如果运行抢占，就是抢占式，这个和golang版本有关系，1.14之前是非抢占式的，1.14之后是抢占式的。

runtime.main会创建一个额外的M运行sysmon函数, 抢占就是在sysmon中实现的.
sysmon会进入一个无限循环, 第一轮回休眠20us, 之后每次休眠时间倍增, 最终每一轮都会休眠10ms。

sysmon中有netpool(获取fd事件), retake(抢占), forcegc(按时间强制执行gc), scavenge heap(释放自由列表中多余的项减少内存占用)等处理.

retake函数负责处理抢占, 流程是:

1. 枚举所有的P
2. 如果P在系统调用中(_Psyscall), 且经过了一次sysmon循环(20us~10ms), 则抢占这个P
3. 调用handoffp解除M和P之间的关联。如果P在运行中(_Prunning), 且经过了一次sysmon循环并且G运行时间超过forcePreemptNS(10ms), 则抢占这个P。
4. 调用preemptone函数
    1. 设置g.preempt = true
    2. 设置g.stackguard0 = stackPreempt

具体抢占有如下几步：

如果不是GC引起的则调用gopreempt_m函数完成抢占，gopreempt_m函数会调用goschedImpl函数, goschedImpl函数的流程是:

1. 把G的状态由运行中(_Grunnable)改为待运行(_Grunnable)
2. 调用dropg函数解除M和G之间的关联
3. 调用globrunqput把G放到全局运行队列
4. 调用schedule函数继续调度

schedule会继续找下一个可运行的G。

因为全局运行队列的优先度比较低, 各个M会经过一段时间再去重新获取这个G执行,
抢占机制保证了不会有一个G长时间的运行导致其他G无法运行的情况发生。

## 8 调度的本质

调度的本质其实都是，修改寄存器得值来做到CPU切换系统进程/goroutine的调度。

go 程序执行是由program和runtime组成，用户进行的系统调用，都是由runtime来拦截，以此帮助他进行垃圾回收以及调度等工作，runtime维护了所有goroutine，并通过go scheduler调度，goroutine和thread是相互独立的，但是goroutine依赖于thread才能执行。

runtime起始时会启动一些g，比如垃圾回收的g，运行代码的g，执行调度的g，并且会创建一个M用来执行这些g。


![pic](https://brt-1303999354.cos.ap-shanghai.myqcloud.com/image2.png)

