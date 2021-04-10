# 谈谈你对IO多路复用的理解

## 1 网络IO硬件的理解

1. 发送者 接受者 连接 文件描述符 IP 端口 传输协议 状态
2. TCP是操作系统内核实现的 内核缓冲区
3. 用户缓冲区

1. 文件描述符
2. 网卡
3. 内核缓冲区
4. 用户缓冲区

## 2 阻塞式的read

开始read 

网卡

内核缓冲区

用户缓冲区

read到

开多线程、多进程

1 read 2 read 3 read

1M 1000 1G 16G 16w

## 3 非阻塞式的read

开始read 

网卡 read = -1

内核缓冲区 read = -1

内核缓冲区接收了所有数据，改写文件描述符

开始read 

内核缓冲区 -> 用户缓冲区

read到

文件描述符可以放到一个线程去管理，维护一个数组

用户态做遍历
    -1       -1       阻塞       -1
[socket1, socket2 , socket3 , socket4]

文件描述符和文件的区别

1. 用户态去查看内核态状态，用户态去查看文件，是一个系统调用

## 4 select

select

用户维护数组
[socket1, socket2 , socket3，socket5, socket6] 维护副本 tmp

tmp交给内核遍历
[socket1, socket2 , socket3，socket5, socket6] 2

维护数组的操作是系统调用

socket3文件描述符发生改变，会去返回数组变化

select = 2

用户再次遍历数组，找出变化的文件描述符

## 5 poll

没有根本性的变化，加大了可监听的文件描述符的数量，突破了1024的限制

数据结构从数组变成了链表

## 6 epoll

1 epoll_create
2 epoll_ctl
3 epoll_wait

1 epoll_create同样是维护了一个文件描述符数组，给到内核
2 假如数组中的文件描述符发生了变化，epoll_ctl去处理
3 事件触发机制，内核缓冲区接收了所有数据，改写文件描述符的时候，知道谁改变了，epoll_wait拿到变化的文件描述符的结果
