# 谈谈你对CDN的理解

* [谈谈你对CDN的理解](#谈谈你对cdn的理解)
  * [1 什么是内容分发网络？](#1-什么是内容分发网络)
  * [2 解决了什么问题？](#2-解决了什么问题)
  * [3 CDN加速原理](#3-cdn加速原理)
  * [4 一些坑](#4-一些坑)

一句话概括，提供更好的网络质量，做资源的缓存。

参考：
1. https://cloud.tencent.com/document/product/228/2939
2. https://www.zhihu.com/question/36514327/answer/1604554133

## 1 什么是内容分发网络？

内容分发网络（Content Delivery Network，CDN），是在现有 Internet 中增加的一层新的网络架构，由遍布全球的高性能加速节点构成。

这些高性能的服务节点都会按照一定的缓存策略存储业务内容，当用户向某一业务内容发起请求时，请求会被调度至最接近用户的服务节点，直接由服务节点快速响应，有效降低用户访问延迟，提升可用性。

## 2 解决了什么问题？

1. 用户与业务服务器地域间物理距离较远，需要进行多次网络转发，传输延时较高且不稳定。
2. 用户使用运营商与业务服务器所在运营商不同，请求需要运营商之间进行互联转发。
3. 业务服务器网络带宽、处理能力有限，当接收到海量用户请求时，会导致响应速度降低、可用性降低。

## 3 CDN加速原理

假设您业务源站域名为：www.test.com

域名接入 CDN 开始使用加速服务后，当您的发起 HTTP 请求时，实际的处理流程如下图所示：

![pic](https://brt-1303999354.cos.ap-shanghai.myqcloud.com/cdn.png)

1. 用户向 www.test.com 下的某图片资源（如：1.jpg）发起请求，会先向 Local DNS 发起域名解析请求。
2. 当 Local DNS 解析 www.test.com 时，会发现已经配置了 CNAME:
www.test.com.cdn.dnsv1.com ，然后就会递归查询到运营商提供的DNS服务器，运营商的DNS服务器会为请求分配最佳节点IP。
3. Local DNS 获取运营商的DNS服务器返回的解析 IP。
4. 用户获取解析 IP。
5. 用户向获取的 IP 发起对资源 1.jpg 的访问请求。
6. 若该 IP 对应的节点缓存有 1.jpg，则会将数据直接返回给用户（10），此时请求结束。
7. 若该节点未缓存 1.jpg，则节点会向业务源站发起对 1.jpg 的请求（6、7、8），获取资源后，结合用户自定义配置的缓存策略（缓存过期配置），将资源缓存至节点（9），并返回给用户（10），此时请求结束。

## 4 一些坑

以淘宝 CDN 为例：

> 淘宝的图片业务除了访问量大，还会面临更新频繁的问题。图片的频繁更新，一方面会由于商品上的图片url变化，导致商品缓存失效，另一方面会大幅降低CDN的图片访问缓存命中率。
>
> 针对图片url变化导致商品缓存失效的问题，我们通过刷新cdn缓存，用户访问时重新回源的方式，实现了改图保持图片url不变，这个过程中，我们解决了一些列的问题，包括：OSS三地同步更新、图片尺寸收敛、图片域名收敛、客户端及浏览器本地缓存。
>
> 针对改图降低CDN图片缓存命中率的问题，我们根据业务的特点，提前合成不同波段的图片，并预热到CDN，保障了源站的安全。

主要问题为：

1. 如果在用户访问高峰期，图片内容大批量发生变化，大量用户的访问就会穿透cdn，对源站造成巨大的压力。
2. 如果网站新增内容，CDN没有收录，对新内容的访问会直接打到源站。
3. 难以区分盗刷和正常访问，如果CDN按流量计费，可能损失很大。
4. 一些资源实时性要求很高，比如淘宝的商品图片改价格，CDN也需要实时同步。
5. 存在客户端及浏览器缓存。
6. CDN同步需要时间，不会立刻生效。

难点是实时性和大批量修改，这就对架构和技术有很高的要求。

一些技术手段：

1. 改图保持图片URL不变。
2. OSS三地同步。
3. 图片尺寸收敛。
4. 多副本清除CDN缓存。
5. 图片域名收敛。
6. 客户端及浏览器缓存。
7. 提前预热CDN图片。
