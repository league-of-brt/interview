# Redis有几大数据类型？说说操作方法及使用场景。

* [Redis有几大数据类型？说说操作方法及使用场景。](#redis有几大数据类型说说操作方法及使用场景)
  * [1 string](#1-string)
    * [1\.1 操作](#11-操作)
    * [1\.2 使用场景](#12-使用场景)
      * [1\.2\.1 高频信息](#121-高频信息)
      * [1\.2\.2 session共享](#122-session共享)
      * [1\.2\.3 阅读数统计](#123-阅读数统计)
      * [1\.2\.4 分布式锁](#124-分布式锁)
      * [1\.2\.5 分布式锁的坑](#125-分布式锁的坑)
  * [2 hash](#2-hash)
    * [2\.1 操作](#21-操作)
    * [2\.2 使用场景](#22-使用场景)
  * [3 list](#3-list)
    * [3\.1 特性](#31-特性)
    * [3\.2 操作](#32-操作)
    * [3\.3 使用场景](#33-使用场景)
      * [3\.3\.1 朋友圈点赞功能（注意用法）](#331-朋友圈点赞功能注意用法)
      * [3\.3\.2 点菜系统](#332-点菜系统)
  * [4 set](#4-set)
    * [4\.1 特性](#41-特性)
    * [4\.2 操作](#42-操作)
    * [4\.3 使用场景](#43-使用场景)
      * [4\.3\.1 抽奖](#431-抽奖)
      * [4\.3\.2 共同好友统计](#432-共同好友统计)
  * [5 zset](#5-zset)
    * [5\.1 特性](#51-特性)
    * [5\.2 操作](#52-操作)
    * [5\.3 使用场景](#53-使用场景)
      * [5\.3\.1 考试成绩排名](#531-考试成绩排名)
      * [5\.3\.2 新闻热度榜](#532-新闻热度榜)
  * [6 总结](#6-总结)

string、hash、list、set、sorted_set这五个基础数据类型，大家都快背烂了，基本操作多多少少也有了解。

但是据我所知，很多人开发中基本只用到了一个string，别的数据类型没有太多的使用场景。所以这次我准备举几个例子，看看怎么用redis解决。

## 1 string

### 1.1 操作

> 1. set key value: 新增或更新字符串键值对
> 2. mset key value [key1 value1 ...]：批量新增或更新键值对
> 3. setnx key value：如果key不存在就添加，否则就失败
> 4. setex key seconds value：设置简直对的时同时设置过期时间
> 5. get key ：获取指定key的值
> 6. mget key [key1 key2 ...]：获取多个key的值
> 7. del key [key1 key2 ...]：删除指定key
> 8. expire key seconds：设置指定key过期时间，以秒为单位
> 9. ttl key：查看指定key还剩余多长时间
> 10. incr key：将指定key存储的数值加1
> 11. decr key：将指定key存储的数值减1
> 12. incrby key step：将指定key存储的数值加上step
> 13. decrby key step ：将指定key存储的数值减去step

### 1.2 使用场景

#### 1.2.1 高频信息

微博大V主页高频的访问，对于粉丝数、关注数、微博数都需要时时更新。

```bash
127.0.0.1:6379> set user:id:12:fans 100
OK
127.0.0.1:6379> set user:id:12:blogs 200
OK
127.0.0.1:6379> set user:id:12:following 300
OK
```

如果增加了粉丝数，可以直接：

```bash
127.0.0.1:6379> incr user:id:12:fans
(integer) 101
```

还有一种是json封装一下，直接跟上所有信息：

```bash
127.0.0.1:6379> set user:id:12:info {fans:100,blogs:200,following:300}
OK
127.0.0.1:6379> get user:id:12:info
"{fans:100,blogs:200,following:300}"
```

具体使用哪种方案，取决于业务。第二种方式显然可以维护更少的key，但是如果字段更新频繁，就会经常读写，效率比较差。

#### 1.2.2 session共享

对于分布式系统或者集群系统，有时候需要共享用户session，这时候可以转换成json进行缓存。

```bash
127.0.0.1:6379> set user:id:12:session "{id:12,name:xie}"
OK
127.0.0.1:6379> get user:id:12:session
"{id:12,name:xie}"
```

#### 1.2.3 阅读数统计

博客的阅读数统计。

```bash
127.0.0.1:6379> set post:id:1:reader:num 100
OK
127.0.0.1:6379> incr post:id:1:reader:num
(integer) 101
```

#### 1.2.4 分布式锁

分布式锁借助setnx完成，set if not exists，没有值就能新增成功，否则就失败，就和锁的抢占是一样的，同时需要考虑值的过期和删除。

先提一嘴，Redis有三种部署方式：

1. 单机模式
2. Master-Slave + Sentinel 选举模式
3. Redis Cluster 模式

示例代码在单机Redis上可以正常运行，另外两种涉及分布式，需要考虑更多问题，不是简单的执行Redis命令就完事了。

```bash
127.0.0.1:6379> setnx lock:product:id:1 true # 抢占不过期的锁
(integer) 1
127.0.0.1:6379> setnx lock:product:id:1 false # 再想抢占会失败，说明被锁住了
(integer) 0
127.0.0.1:6379> del lock:product:id:1
(integer) 1
127.0.0.1:6379> setnx lock:product:id:1 true # 删除之后可以再抢占
(integer) 1
127.0.0.1:6379> expire lock:product:id:1 60 # 可以设置超时时间，但是这个操作是错误的，不应该使用
(integer) 1
127.0.0.1:6379> ttl lock:product:id:1 # 倒计时
(integer) 49
127.0.0.1:6379> ttl lock:product:id:1
(integer) 17
127.0.0.1:6379> ttl lock:product:id:1
(integer) 1
127.0.0.1:6379> ttl lock:product:id:1 # 锁过期会得到赋值
(integer) -2
127.0.0.1:6379> get lock:product:id:1
(nil)
127.0.0.1:6379> set lock:product:id:1 true ex 10 nx # 正确的设置过期时间，保证操作的原子性
OK
```

#### 1.2.5 分布式锁的坑

这里只是稍微提醒一下想要写分布式锁的同学，如果你实现的是单机Redis的分布式锁，是玩票性质的，那没问题，执行Redis命令就搞定。

但是要考虑清楚，如果是单机Redis，你的Master挂掉了，整个系统的分布式锁就都没了，要事先想好能否承受这种后果。

于是你醒悟了，为了高可用，采用Master-Slave模式，通过Sentinel做了高可用。但是，如果加锁的时候只对Master加锁，此时需要对Slave进行同步，这时一旦master挂掉了，选举上来的slave是没有这个锁的，就会出现锁丢失的问题。

基于以上的考虑，Redis的作者也考虑到这个问题，他提出了一个RedLock的算法。

这个算法的意思大概是这样的：假设Redis的部署模式是Redis Cluster，总共有5个Master节点。

通过以下步骤获取一把锁：

1. 获取当前时间戳，单位是毫秒。
2. 轮流尝试在每个 Master 节点上创建锁，过期时间设置较短，一般就几十毫秒。
3. 尝试在大多数节点上建立一个锁，比如 5 个节点就要求是 3 个节点（n / 2 +1）。
4. 客户端计算建立好锁的时间，如果建立锁的时间小于超时时间，就算建立成功了。
5. 要是锁建立失败了，那么就依次删除这个锁。
6. 只要别人建立了一把分布式锁，你就得不断轮询去尝试获取锁。

但是这样的这种算法还是颇具争议的，可能还会存在不少的问题，无法保证加锁的过程一定正确。

个人认为，如果你要用Redis实现真正意义上的分布式锁，几乎是不太可能的，一般建议使用开源框架，比如Redission。

甚至，就算是Redis的主创，也在自己的博客上贴出了这样一个问题：

![pic](https://brt-1303999354.cos.ap-shanghai.myqcloud.com/unsafe-lock.png)

总之，如果非要用Redis做分布式锁，建议使用Redission，更好的解决方案是Zookeeper，如果你坚持要自己实现，那只能祝你好运了。

## 2 hash

### 2.1 操作

> 1. hset key field value: 新增或更新key对应字段的值
> 2. hsetnx key field value：新增一个不存在Key的字段值
> 3. hmset key field value [field value ...]：在指定Key上存储多个字段和值
> 4. hget key field：获取指定key中指定字段的值
> 5. hdel key field [field1...]：删除指定Key值的指定字段
> 6. hlen key：获取指定key中的字段的数量
> 7. hgetall key：获取指定key中所有的字段值
> 8. hincrby key field step：指定key中字段值增加step

### 2.2 使用场景

某对象存在多重属性，即存在大量key-value映射，同时又想进行一定的操作，可以使用hash。

个人认为，难点在于怎么设计一个合理的结构，比较考验抽象的能力。

购物车的基本操作：

![pic](https://brt-1303999354.cos.ap-shanghai.myqcloud.com/1591066533641-f8c09431-59d4-4e28-b6e3-a5168dd9ddea.png)

```bash
127.0.0.1:6379> hmset shoppingcar:id:1 1 2 2 4
OK
127.0.0.1:6379> hget shoppingcar:id:1 1
"2"
127.0.0.1:6379> hget shoppingcar:id:1 2
"4"
127.0.0.1:6379> hgetall shoppingcar:id:1 # 拿到所有物品
1) "1"
2) "2"
3) "2"
4) "4"
127.0.0.1:6379> hlen shoppingcar:id:1 # 物品种类
(integer) 2
127.0.0.1:6379> hincrby shoppingcar:id:1 1 10 # 增加物品
(integer) 12
127.0.0.1:6379> hdel shoppingcar:id:1 2 # 删除物品
(integer) 1
127.0.0.1:6379> hgetall shoppingcar:id:1
1) "1"
2) "12"
```

有些同学可能就要问了，这种场景，我使用string不一样能实现吗？

确实可以，但要看你怎么实现了。如果一个购物车里面存了一大堆东西，如果你对整体做序列化，成本还是比较高的。如果你拆分成一个个key，你就得维护这些分别key，在我看来也比较麻烦。

## 3 list

### 3.1 特性

比较需要注意，列表可以进行左右操作，所以用list还得有一定的方向感。

### 3.2 操作

> 1. lpush key value [value1 ...]：在指定key的列表左边插入一个或多个值
> 2. rpush key value [value1 ...]：在指定key的列表右边插入一个或多个值
> 3. lpop key：从指定key的列表左边取出第一个值
> 4. rpop key：从指定key的列表右边取出第一个值
> 5. lrange key start end：从指定key列表中获取指定区间内的数据
> 6. blpop key [key1 ...] timeout：从指定key列表中左边取出第一个值，若列表中没有元素就等待timeout时间，如果timeout为0就一直等待
> 7. brpop key [key1 ...] timeout：从指定key列表中右边取出第一个值，若列表中没有元素就等待timeout时间，如果timeout为0就一直等待
> 8. lset key index value：将指定下标的值更新为value

### 3.3 使用场景

因为list可以实现队列的功能，只要自己定义生产者消费者，就能整出一个低配消息队列，所以很多人似乎拿着list当消息队列去使用。

我个人是不建议的，因为缺少很多必要的功能，比如Redis崩溃消息可能会丢失，消费者遇到错误也要维护一套自己的重放机制。

个人认为，list就老老实实扮演一个存东西的容器就好了。

#### 3.3.1 朋友圈点赞功能（注意用法）

```bash
127.0.0.1:6379> rpush like:id:1 1 2 3 4 # 1、2、3、4用户点赞
(integer) 4
127.0.0.1:6379> lrange like:id:1 0 -1 # 展示
1) "1"
2) "2"
3) "3"
4) "4"
127.0.0.1:6379> lrem like:id:1 0 3 # 3用户取消点赞
(integer) 1
127.0.0.1:6379> rpush like:id:1 5 # 5用户点赞
(integer) 4
127.0.0.1:6379> lrange like:id:1 0 -1 # 排序依然按照时间顺序
1) "1"
2) "2"
3) "4"
4) "5"
```

一般来说，朋友圈不支持重复点赞，但是List支持放入重复的元素。

所以这个需求，在放入元素之前，得先找找是不是已经存在重复元素。

#### 3.3.2 点菜系统

再想一个可重复的场景，比如点菜系统，用户可以点很多重复的菜式，厨师根据点的菜顺序去做菜：

```bash
127.0.0.1:6379> rpush cook:id:1 {table_id:1,meal_id:2,meal_num:3}
(integer) 1
127.0.0.1:6379> rpush cook:id:1 {table_id:2,meal_id:3,meal_num:4}
(integer) 2
127.0.0.1:6379> rpush cook:id:1 {table_id:1,meal_id:2,meal_num:1}
(integer) 3
127.0.0.1:6379> lrange cook:id:1 0 -1
1) "{table_id:1,meal_id:2,meal_num:3}"
2) "{table_id:2,meal_id:3,meal_num:4}"
3) "{table_id:1,meal_id:2,meal_num:1}"
127.0.0.1:6379> rpop cook:id:1 # 这里需要注意，rpop弹出元素，并不会删除
"{table_id:1,meal_id:2,meal_num:1}"
127.0.0.1:6379> lrem cook:id:1 0 {table_id:1,meal_id:2,meal_num:1} # 取消订单
(integer) 0
127.0.0.1:6379> lrange cook:id:1 0 -1
1) "{table_id:1,meal_id:2,meal_num:3}"
2) "{table_id:2,meal_id:3,meal_num:4}"
127.0.0.1:6379> rpush cook:id:1 {table_id:3,meal_id:4,meal_num:5} # 追加订单
(integer) 3
127.0.0.1:6379> lrange cook:id:1 0 -1
1) "{table_id:1,meal_id:2,meal_num:3}"
2) "{table_id:2,meal_id:3,meal_num:4}"
3) "{table_id:3,meal_id:4,meal_num:5}"
```

## 4 set

### 4.1 特性

元素无序且不重复。

### 4.2 操作

> 1. sadd key member [member ...]：在集合中增加一个或多个元素
> 2. srem key member [member ...]：从集合中删除一个或多个元素
> 3. smembers key：获取集合中的所欲元素
> 4. scard key：获取集合中的元素个数
> 5. sismember key member：判断指定member是否在集合中
> 6. srandmember key [count]：从集合中获取count个元素，不从集合中删除
> 7. spop key [count]：从集合中获取count个元素，从集合中删除
> 8. sinter key [key1 ...]：指定多个集合进行交集运算
> 9. sinterstore dest key [key1 ...]：指定多个集合进行交集运算，存入dest集合
> 10. sunion key [key1 ...]：指定多个集合进行并集运算
> 11. sunionstore dest key [key1 ...]：指定多个集合进行并集运算，存入dest集合
> 12. sdiff key [key1 ...]：指定多个集合进行差集运算
> 13. sdiffstore dest key [key1 ...]：指定多个集合进行差集运算，并存入dest集合

### 4.3 使用场景

#### 4.3.1 抽奖

抽奖的人不能重复，理想的使用Set的场景。

```bash
127.0.0.1:6379> sadd lottery:id:1 1 2 3 4 5 6
(integer) 6
127.0.0.1:6379> sadd lottery:id:1 1 # 不可重复加入
(integer) 0
127.0.0.1:6379> sadd lottery:id:1 7 8
(integer) 2
127.0.0.1:6379> smembers lottery:id:1
1) "1"
2) "2"
3) "3"
4) "4"
5) "5"
6) "6"
7) "7"
8) "8"
127.0.0.1:6379> srandmember lottery:id:1 # 随机弹出1个元素，继续在Set中保留
"6"
127.0.0.1:6379> smembers lottery:id:1
1) "1"
2) "2"
3) "3"
4) "4"
5) "5"
6) "6"
7) "7"
8) "8"
127.0.0.1:6379> spop lottery:id:1 5 # 随机弹出5个元素，不会在Set中保留
1) "7"
2) "6"
3) "2"
4) "1"
5) "8"
127.0.0.1:6379> smembers lottery:id:1
1) "3"
2) "4"
3) "5"
```

#### 4.3.2 共同好友统计

```bash
127.0.0.1:6379> sadd a 1 2 3
(integer) 3
127.0.0.1:6379> sadd b 1 3 5 7
(integer) 4
127.0.0.1:6379> sinter a b # a、b的共同好友
1) "1"
2) "3"
127.0.0.1:6379> sdiff a b # a有b没有的好友
1) "2"
127.0.0.1:6379> sdiff b a # b有a没有的好友
1) "5"
2) "7"
127.0.0.1:6379> sunion a b # a和b的所有好友
1) "1"
2) "2"
3) "3"
4) "5"
5) "7"
```

## 5 zset

### 5.1 特性

元素是有序不可重复的。

和Set用法基本一样，只是每个元素中多了一个分值，用于元素排序。

### 5.2 操作

> 1. zadd key score member [(score member)...]：往有序集合中添加带分值的元素
> 2. zrem key member [member...]：从有序集合中删除成员
> 3. zscore key member：返回集合中指定成员的分值
> 4. zcard key：统计集合中元素个数
> 5. zrange key start stop [withscores]：返回指定范围的元素，withscores代表返回的元素包含对应的分值
> 6. zreverange key start stop [withscores]：返回指定范围的倒序元素，withscores代表返回的元素包含对应的分值
> 7. 同set一样也可以进行交集、并集、差集的集合运算

### 5.3 使用场景

#### 5.3.1 考试成绩排名

一般来说，学生成绩考出来就不怎么变，以学生成绩为分值，学生姓名为元素。

```bash
127.0.0.1:6379> zadd exam:id:1 90 a 100 b 59 c 60 d 61 e
(integer) 5
127.0.0.1:6379> zrange exam:id:1 0 -1 # 从小到大输出
1) "c"
2) "d"
3) "e"
4) "a"
5) "b"
127.0.0.1:6379> zrevrange exam:id:1 0 -1 # 从大到小输出
1) "b"
2) "a"
3) "e"
4) "d"
5) "c"
127.0.0.1:6379> zrange exam:id:1 0 -1 withscores # 同时输出分数
 1) "c"
 2) "59"
 3) "d"
 4) "60"
 5) "e"
 6) "61"
 7) "a"
 8) "90"
 9) "b"
10) "100"
127.0.0.1:6379> zrevrange exam:id:1 0 -1 withscores
 1) "b"
 2) "100"
 3) "a"
 4) "90"
 5) "e"
 6) "61"
 7) "d"
 8) "60"
 9) "c"
10) "59"
```

#### 5.3.2 新闻热度榜

新闻热度可能随时变化，以热度为分值，以新闻ID为元素。

```bash
127.0.0.1:6379> zadd rank:id:1 500 1
(integer) 1
127.0.0.1:6379> zadd rank:id:1 114 2
(integer) 1
127.0.0.1:6379> zadd rank:id:1 514 3
(integer) 1
127.0.0.1:6379> zadd rank:id:1 114514 4
(integer) 1
127.0.0.1:6379> zrange rank:id:1 0 10 withscores
1) "2"
2) "114"
3) "1"
4) "500"
5) "3"
6) "514"
7) "4"
8) "114514"
127.0.0.1:6379> zincrby rank:id:1 100 1 # 热度增加，排名也可能有变化
"600"
127.0.0.1:6379> zrange rank:id:1 0 10 withscores
1) "2"
2) "114"
3) "3"
4) "514"
5) "1"
6) "600"
7) "4"
8) "114514"
```

## 6 总结

基本数据类型是一定要背的，特性也是一定要背的，方法很难记忆，最好记几个经典使用场景，顺带着就把方法记下来了。
