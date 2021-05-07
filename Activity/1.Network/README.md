# 计网第一期分享

# 1 TCP

1. 三次握手
    1. 为什么需要三次？怎么握手？
    2. 两次不行吗，会遇到什么问题？
    3. 四次不行吗，会遇到什么问题？
    4. 我们平时写代码怎么没有自己实现握手，在哪里进行了握手步骤？

2. 四次挥手
    1. 为什么需要四次挥手？怎么挥手？
    2. 怎么理解time-wait和close-wait状态？
    3. 在四次挥手中，为什么客户端进入TIME_WAIT状态后需要等待 2MSL 之后才CLOSE？ 如果不等待 2MSL 直接CLOSE，会导致什么结果？
    4. 在四次挥手的各阶段，假如因为网络原因丢包，会导致什么结果？

3. 谈谈你对粘包问题的理解？

4. TCP怎么保证传输的可靠性？
    1. 超时重传
    2. 快速重传
    3. SACK
    4. D-SACK

5. 谈谈你对TCP流量控制的理解？
    1. 滑动窗口
    2. 操作系统缓冲区和滑动窗口的关系
    3. 窗口关闭
    4. 糊涂窗口综合征（少问）

6. 谈谈你对TCP拥塞控制的理解？
    1. 慢启动
    2. 拥塞避免
    3. 拥塞发生怎么处理
    4. 快速恢复

7. 你认为TCP有哪些优点？有哪些缺点？什么场景适合使用TCP？

8. 如果让你优化TCP，你会做些什么？

9. 问题处理
    1. 什么是SYN攻击？
    2. 什么是全连接攻击？
    3. 什么是DDOS攻击？
    4. 一个服务器最多支持多少TCP连接？
    5. 怎么打满一个服务器的连接？
    6. 怎么处理大量time-wait和close-wait？
    7. 什么是中间人攻击？

# 2 UDP

1. 谈谈TCP和UDP的联系和区别？

2. 你认为UDP有哪些优点？有哪些缺点？什么场景适合使用UDP？

3. 如果让你优化UDP，你会做些什么？

# 3 QUIC

谈谈你对QUIC的理解？

# 4 HTTP

1. 谈谈HTTP的发展历史

2. 怎么区分header和body

3. 谈谈content-length

4. 谈谈分块传输

5. HTTP 和 TCP/UDP 的关系

6. 长连接短链接

7. 队头阻塞

8. TCP 会导致 HTTP请求乱序吗

9. 谈谈 HTTP 1.1 的缺陷

10. 谈谈 HTTP 2 的优点

11. 实际问题
    1. 我现在给服务器传递一堆信息，希望服务器将其压缩，怎么处理流程
    2. 谈谈请求方式的选择
    3. 谈谈状态码（301、304...）

# 5 HTTPS

1. 对称加密和非对称加密

2. 什么是中间人攻击？

3. 为什么HTTP不够安全？

4. HTTPS是对称加密还是非对称加密？如果是对称加密会有什么问题？

5. 什么是证书？

6. 如何建立HTTPS连接？能不能画出来？

7. 为什么HTTPS安全？

8. HTTPS是否绝对安全？

9. 问题处理
    1. 什么是重放攻击？
    2. 什么是慢速攻击？
        1. Slow headers
        2. Slow body
        3. Slow read

# 6 在浏览器地址栏输入一个URL后回车，背后会进行哪些技术步骤？

1. 浏览器对URL的预处理

2. DNS流程做了什么？

3. IP做了什么？

4. ARP做了什么？

5. MAC做了什么？

6. 什么是根服务器？

7. 网卡、交换机、路由器

7. CDN

8. 后端技术

9. 实际问题
    1. 一个页面能开多少TCP连接？
    2. 怎么加快访问速度？