# GMP问题集

1. 为什么需要 P 这个组件，直接把 runqueues 放到 M 不行吗？
2. gpm到底是什么
3. scheduler是如何调度的
4. 什么时候会触发调度
5. 当在M上运行的goroutine发生阻塞时，会怎么工作
6. 为什么每个m都会对应一个g0 (g0是用于调度每个线程中的goroutine，包括gc等等，拥有比较大的栈内存)
7. 什么时候会抢占P
8. 调度的本质
9. 多个线程与多个M如何一一对应？
10. 为什么要把工作线程与m对应
11. 为什么在创建goroutine的newproc函数要传入参数大小
12. 什么时候调用的main函数？
13. g0到main goroutine的转换过程
14. 非main goroutine是如何返回到goexit函数的；
15. mcall函数如何从用户goroutine切换到g0继续执行
16. 调度循环
17. 使用什么策略来挑选下一个进入运行的goroutine
18. 如何把挑选出来的goroutine放到CPU上运行
19. schdule的三种调度方式