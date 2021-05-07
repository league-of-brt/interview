# 谈谈你对数据库隔离级别和mvcc的理解

## 1 什么叫一段锁，两段锁？

func test() {
    mutex.lock()
    dosomething()
    mutex.unlock()
}

## 2 事务的隔离级别

|隔离级别|脏写|脏读|不可重复读|幻读|
|----|----|----|----|----|
|未提交读RU|√|√|√|√|
|已提交读RC|×|×|√|√|
|可重复读RR|×|×|×|√|
|串行化S|×|×|×|×|

## 3 Mysql的innodb事务的隔离级别

|隔离级别|脏写|脏读|不可重复读|幻读|
|----|----|----|----|----|
|未提交读RU|√|√|√|√|
|已提交读RC|×|×|√|√|
|可重复读RR|×|×|×|×|
|串行化S|×|×|×|×|

## 4 Mysql中锁的种类

1. 表锁
2. 行锁
    1. 读意向锁
    2. 写意向锁
    3. 读锁
    4. 写锁
    5. 自增锁

## 5 Mysql中锁的粒度

1. 表锁
2. 行锁
    1. 记录锁
    2. 间隙锁
    3. Next-Key锁
    4. 插入意向锁

## 7 什么是脏写

|id|name|
|----|----|
|1|name1|

====================================================

|A|B|
|----|----|
|begin|begin|
|UPDATE user SET name = 'name2' where id = 1||
||UPDATE user SET name = 'name3' where id = 1|
||commit|
|rollback||

## 8 什么是脏读

|id|name|age|
|----|----|----|
|1|name1|2|

====================================================

|A|B|
|----|----|
|begin|begin|
|UPDATE user SET name = 'name2' where id = 1||
||SELECT FROM user where id = 1|
||do something|
|rollback||
||SELECT FROM user where id = 1|

## 9 什么是不可重复读

|id|name|
|----|----|
|1|name2|

====================================================

|A|B|
|----|----|
|begin|begin|
|SELECT FROM user where id = 1||
||UPDATE user SET name = 'name2' where id = 1|
||commit|
|SELECT FROM user where id = 1||

## 10 什么是幻读

|id|name|
|----|----|
|1|name1|
|2|name2|

====================================================

|A|B|
|----|----|
|begin|begin|
|SELECT FROM user where id > 1||
||INSERT INTO user VALUES(2,'name2')|
||commit|
|SELECT FROM user where id > 1||

## 11 不可重复读和幻读的区别

## 12 什么是mvcc，什么是READ VIEW，什么是版本链

mysql row 有预留字段 4 

1. 预留主键
2. 删除标识
3. A开了一个事务，ID = 10，我在A事务中修改或者创建一行数据，就会修改row_trx_id
4. ptr，指向了undolog

A 10 修改name = name1
B 20 修改name = name2

原本行 trx_id = 20 name = name2 ptr   
            |
版本2 trx_id = 10 data = name1 ptr
            |    
版本1 trx_id = 0 data = name0 null

READ VIEW

开启事务A，执行一条SQL语句，就会开启READ VIEW

A 20 [20,30] 20 31
B 30 [20,30] 20 31

READ VIEW

1. 正在并发的事务id数组 
2. 当前所有事务最小的id
3. 记录所有事务最大的id+1



## 13 为什么RC能解决脏读问题

A
select
A 20 [20,30] 20 31

select
A [20] 20 31 20

B update

最新记录trx_id = 30

A 20 [20,30] 20 31
B 30 [20,30] 20 31

## 14 为什么RC有不可重复读问题

## 15 为什么RC有幻读问题

## 16 为什么RR能解决不可重复读

## 17 为什么innodb的RR能解决幻读问题

select * from user where id > 1

A 20 [20,30] 20 31
B 30 [20,30] 20 31

B插入两条新纪录 30 30

## 18 快照读，当前读

select * from user where nid > 1 for update

id > 1 加next-key锁

id > 1全部锁住

0 name0
1 name1
2 name2

select mvcc for update

update
delete
insert

|id|name|age|
|----|----|----|
|1|name1|1|

|A|B|
|----|----|
|begin|begin|
|UPDATE user SET age=age+1 where id = 1||
||UPDATE user SET age=age+1 where id = 1|
||commit|
|SELECT FROM user where id = 1||
|commit||

A [20,30] 20 31 20
B [20,30] 20 31 30
C [20,30,40] 20 41 40

A B C

1. 全表
2. A B C

A 50% m 50%f

100 99m 1f
50 50

1. 重新采样
2. sql强制走索引 force index
