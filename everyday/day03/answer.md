# MySQL锁机制探讨
前提： MySQL8.0 InnoDB
参考书籍：MySQL技术内幕InnoDB存储引擎

## 1、锁机制、锁粒度

### 锁概念
InnoDB中存在两种锁概念：lock、latch
#### lock
通常情况下，我们讲的数据库锁就是指lock
- lock的对象
  - 事务，用来锁定的是数据库中的对象，如：表、页、行
  
- lock持续时间
  - 整个事务过程，一般lock的对象仅在事务commit或rollback后释放
  
- lock的模式
  - 行锁（InnoDB默认）
  - 表锁（MyISAM默认，InnoDB不探讨）
  - 意向锁（InnoDB表级锁）
    
##### 行锁
- 锁类型
  - 共享锁（S Lock），允许事务读一行数据
  - 排他锁（X Lock），允许事务删除或更新一行数据
- 锁算法（排查死锁时会看到算法名相关的字眼，如果要搞清楚为什么死锁就得去了解锁算法）
  - Record Lock
  - Gap Lock
  - Next-Key Lock
  
##### 意向锁
- 锁类型
  - 意向共享锁（IS Lock），事务想要获取一张表某几行的共享锁
  - 意向排他锁（IX Lock），事务想要获取一张表某几行的排他锁
    
#### latch
latch称为闩锁，要求锁定时间非常短，其目的是用来保证并发线程操作临界资源（暂不深入）
- latch的对象
  - 线程
- latch持续时间
  - 临界资源


## 2、讲清楚一条insert、select语句的锁使用情况
## 3、常见死锁场景或sql，如何避免
