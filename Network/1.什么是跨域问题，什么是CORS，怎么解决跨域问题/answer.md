# 什么是跨域问题？什么是CORS？怎么解决跨域问题？

* [什么是跨域问题？什么是CORS？怎么解决跨域问题？](#什么是跨域问题什么是cors怎么解决跨域问题)
  * [1 跨域问题](#1-跨域问题)
    * [1\.1 什么是跨域问题？](#11-什么是跨域问题)
    * [1\.2 请你演示一下跨域问题？](#12-请你演示一下跨域问题)
    * [1\.3 现实开发中遇过跨域问题吗？](#13-现实开发中遇过跨域问题吗)
  * [2 同源策略](#2-同源策略)
    * [2\.1 什么是同源策略？](#21-什么是同源策略)
    * [2\.2 为什么需要同源策略？](#22-为什么需要同源策略)
    * [2\.3 同源策略限制哪些行为](#23-同源策略限制哪些行为)
    * [2\.4 先天允许跨域请求的资源](#24-先天允许跨域请求的资源)
  * [3 如何解决跨域问题？](#3-如何解决跨域问题)
    * [3\.1 解决思路](#31-解决思路)
    * [3\.2 CORS](#32-cors)
  * [4 CORS实现机制](#4-cors实现机制)
    * [4\.1 两种请求](#41-两种请求)
    * [4\.2 简单请求](#42-简单请求)
    * [4\.3 简单请求的回应](#43-简单请求的回应)
    * [4\.4 Access\-Control\-Allow\-Origin = \*](#44-access-control-allow-origin--)
    * [4\.5 非简单请求](#45-非简单请求)
    * [4\.6 非简单请求的回应](#46-非简单请求的回应)
    * [4\.7 浏览器的正常请求和回应](#47-浏览器的正常请求和回应)
    * [4\.8 CORS和JSONP的对比](#48-cors和jsonp的对比)

参考：
1. https://www.bilibili.com/video/BV1Kt411E76z、
2. http://www.ruanyifeng.com/blog/2016/04/cors.html
3. https://zhuanlan.zhihu.com/p/145837536

## 1 跨域问题

### 1.1 什么是跨域问题？

前端调用的后端接口不属于同一个域（域名或端口不同），就会产生跨域问题，也就是说你的应用访问了该应用域名或端口之外的域名或端口。

> 由于浏览器同源策略的限制，非同源下的请求，都会产生跨域问题。
>
> 同源策略即：同一协议，同一域名，同一端口号。当其中一个不满足时，我们的请求即会发生跨域问题。

举个简单的例子：

1. http://www.abc.com:3000 到 https://www.abc.com:3000 的请求会出现跨域（域名、端口相同但协议不同）。
2. http://www.abc.com:3000 到 http://www.abc.com:3001 的请求会出现跨域（域名、协议相同但端口不同）。
3. http://www.abc.com:3000 到 http://www.def.com:3000 的请求会出现跨域（域名不同）。

### 1.2 请你演示一下跨域问题？

给我一个chrome就可以演示了。首先访问 www.baidu.com，然后在控制台调用fetch

``` javascript
fetch("https://www.baidu.com") // ok
fetch("https://www.taobao.com") // not ok
```

结果为：

```http
Access to fetch at 'https://www.taobao.com/' from origin 'https://www.baidu.com' has been blocked by CORS policy: No 'Access-Control-Allow-Origin' header is present on the requested resource. If an opaque response serves your needs, set the request's mode to 'no-cors' to fetch the resource with CORS disabled.
```

### 1.3 现实开发中遇过跨域问题吗？

本地开发一个前端项目，这个前端项目是通过node运行的，端口是10001（前端页面通过10001获取）。而服务端是通过spring boot提供的，端口号是10002。

当你通过前端去调用一个服务端接口时，很可能会得到一个异常，返回的response是undefined，并且message消息中只有一个Network Error。

于是你调试一下，发现后端确实收到了请求，也确实处理了请求，但是浏览器确实收到了一个错误。

需要明确的一点，前端去请求后端，这个请求已经发出去了，服务端也接收到并处理了。但是返回的响应结果不是浏览器想要的结果，所以浏览器将这个响应的结果给拦截了，这就是一个经典的跨域错误。

为什么浏览器要拦截响应的结果呢？

这是因为有同源策略，是出于安全方面的考虑。

## 2 同源策略

### 2.1 什么是同源策略？

早在1995年，Netscape 公司就在浏览器中引入了同源策略。

最初的同源策略，主要是限制Cookie的访问，A网页设置的 Cookie，B网页无法访问，除非B网页和A网页是同源的。

### 2.2 为什么需要同源策略？

因为如果没有同源策略的保护，浏览器将没有任何安全可言。

> 老李是一个钓鱼爱好者，经常在 http://51mai.com 网站上买各种钓鱼的工具，并且通过银行 http://yinhang.com 以账号密码的方式直接支付。
>
> 这天老李又在 http://51mai.com 上买了一根鱼竿，输入银行账号密码支付成功后，在支付成功页看到一个叫钓鱼 http://diaoyu.com 的网站投放的一个"免费领取鱼饵"的广告。
>
> 老李什么都没想就点击了这个广告，跳转到了钓鱼的网站，殊不知这真是一个 “钓鱼” 网站，老李银行账户里面钱全部被转走了。

为什么钱会被转走呢？

1. 老李购买鱼竿，并登录了银行的网站输入账号密码进行了支付，浏览器在本地缓存了银行的Cookie。
2. 老李点击钓鱼网站，钓鱼网站使用老李登录银行之后的Cookie，伪造成自己是老李进行了转账操作。

这个过程就是著名的CSRF，跨站请求伪造，正是由于可能存在的伪造请求，导致了浏览器的不安全。

原则上所有浏览器都要遵守同源策略，禁止获得不同源网站的Cookie，防止伪造身份。

### 2.3 同源策略限制哪些行为

据2.2，同源策略是一个安全机制，他本质是限制了从一个源加载的文档或脚本如何与来自另一个源的资源进行交互，用于隔离潜在恶意文件。

随着互联网的发展，同源策略越来越严格，不仅限于Cookie的读取。目前，如果非同源，共有三种行为受到限制：

1. Cookie、LocalStorage 和 IndexDB 无法读取（防止窃取信息）。
2. DOM 无法获得（防止篡改页面）。
3. 请求的响应被拦截。

看到这里你应该明白，为什么1.3的请求会被拦截了，原因就是请求的源和服务端的源不是同源，而服务端又没有设置允许的跨域资源共享，所以请求的响应被浏览器给拦截掉了。

那么问题来了，同源策略确实很安全，但是有时候我们确实需要跨域资源共享，能不能通过某种规则，绕过同源策略呢？

这就需要CORS了。当然CORS也只是解决跨域问题的其中一个策略而已，还有很多别的方法也能解决跨域问题。

### 2.4 先天允许跨域请求的资源

其实有很多资源是允许跨域的，比如图片、css、script。浏览器不会对这些资源进行拦截。

所以有时候大家都说盗链盗链，就是在 www.a.com 未经授权去展示 www.b.com 中的图片，这种行为就等于白嫖。而且盗链是可以层层叠加的，www.a.com 盗链 www.b.com，www.c.com 又去盗链 www.a.com ...

这样层层加码，如果盗链太多甚至会耗尽 www.b.com 的带宽。

这都是题外话了，这种先天允许跨域请求的资源，需要防盗链。

## 3 如何解决跨域问题？

### 3.1 解决思路

我们有三种解决思路：

1. 谁不行我们就干掉谁！既然是浏览器拦截响应的结果，那我们就搞定浏览器，解除跨域限制（不现实）。
2. 发送JSONP请求替代XHR请求（JSONP只支持GET，并不能适用所有的请求方式，不推荐）。
3. 跨域问题本质是前后端交互的问题。对此，我们可以制定一套前后端交互的方案，只要前端后端一起实现，就能支持一定程度的跨域（推荐）。

关于XHR，实际上是XMLHttpRequest，后端同学可能不了解。建议阅读：http://www.ruanyifeng.com/blog/2012/09/xmlhttprequest_level_2.html

### 3.2 CORS

在思路3中，要求制定一套前后端交互的方案。

对此，业界制定了CORS方案
（很多同学可能会以为CORS就是一个报错，其实CORS是为了解决跨域问题的一种方案）。

> CORS是一个W3C标准，全称是"跨域资源共享"（Cross-origin resource sharing）。
>
> 它允许浏览器向跨源服务器，发出XMLHttpRequest请求，从而克服了AJAX只能同源使用的限制
> CORS需要浏览器和服务器同时支持。目前，所有浏览器都支持该功能，IE浏览器不能低于IE10。
> 
> 整个CORS通信过程，都是浏览器自动完成，不需要用户参与。对于开发者来说，CORS通信与同源的AJAX通信没有差别，代码完全一样。浏览器一旦发现AJAX请求跨源，就会自动添加一些附加的头信息，有时还会多出一次附加的请求，但用户不会有感觉。
>
> 因此，实现CORS通信的关键是服务器。只要服务器实现了CORS接口，就可以跨源通信。

## 4 CORS实现机制

对于这个问题，可以参考阮一峰大佬的文章：http://www.ruanyifeng.com/blog/2016/04/cors.html ，我这里只是简单描述一下。

### 4.1 两种请求

浏览器将CORS请求分成两类：简单请求（simple request）和非简单请求（not-so-simple request）。

只要同时满足以下两大条件，就属于简单请求。

>（1) 请求方法是以下三种方法之一：
>   - HEAD
>   - GET
>   - POST
>
>（2）HTTP的头信息不超出以下几种字段：
>   - Accept
>   - Accept-Language
>   - Content-Language
>   - Last-Event-ID
>   - Content-Type：只限于三个值application/x-www-form-urlencoded、multipart/form-data、text/plain

这是为了兼容表单（form），因为历史上表单一直可以发出跨域请求。AJAX的跨域设计就是，只要表单可以发，AJAX就可以直接发。（这里我暂时没有理解?）

凡是不同时满足上面两个条件，就属于非简单请求。

浏览器对这两种请求的处理，是不一样的。

### 4.2 简单请求

对于简单请求，浏览器直接发出CORS请求。具体来说，就是在头信息之中，增加一个Origin字段。

下面是一个例子，浏览器发现这次跨源AJAX请求是简单请求，就自动在头信息之中，添加一个Origin字段。

```http
GET /cors HTTP/1.1
Origin: http://api.bob.com
Host: api.alice.com
Accept-Language: en-US
Connection: keep-alive
User-Agent: Mozilla/5.0...
```

上面的头信息中，Origin字段用来说明，本次请求来自哪个源（协议 + 域名 + 端口）。服务器根据这个值，决定是否同意这次请求。

### 4.3 简单请求的回应

如果Origin指定的源，不在许可范围内，服务器会返回一个正常的HTTP回应。

浏览器发现，这个回应的头信息没有包含Access-Control-Allow-Origin字段（详见下文），就知道出错了，从而抛出一个错误，被XMLHttpRequest的onerror回调函数捕获。

注意，这种错误无法通过状态码识别，因为HTTP回应的状态码有可能是200。

如果Origin指定的域名在许可范围内，服务器返回的响应，会多出几个Access-Control开头的头信息字段。

```http
Access-Control-Allow-Origin: http://api.bob.com
Access-Control-Allow-Credentials: true
Access-Control-Expose-Headers: FooBar
Content-Type: text/html; charset=utf-8
```

其中只有Access-Control-Allow-Origin是必须有的，别的都是可选的。

如果能收到，就说明服务器同一跨域请求了，浏览器不再拦截响应的结果。

### 4.4 Access-Control-Allow-Origin = *

现在你终于知道，为啥网上很多解决跨域问题的回答，要你在服务端请求header加上
`Access-Control-Allow-Origin = *`了吧，其实这就是CORS的一种处理标准，浏览器不再拦截响应。

知其然要知其所以然哦。

### 4.5 非简单请求

非简单请求是那种对服务器有特殊要求的请求，比如请求方法是PUT或DELETE，或者Content-Type字段的类型是application/json。

非简单请求的CORS请求，会在正式通信之前，增加一次HTTP查询请求，称为"预检"请求（preflight）。

浏览器先询问服务器，当前网页所在的域名是否在服务器的许可名单之中，以及可以使用哪些HTTP动词和头信息字段。只有得到肯定答复，浏览器才会发出正式的XMLHttpRequest请求，否则就报错。

下面是一段浏览器的JavaScript脚本：

```javascript
var url = 'http://api.alice.com/cors';
var xhr = new XMLHttpRequest();
xhr.open('PUT', url, true);
xhr.setRequestHeader('X-Custom-Header', 'value');
xhr.send();
```

上面代码中，HTTP请求的方法是PUT，并且发送一个自定义头信息X-Custom-Header。

浏览器发现，这是一个非简单请求，就自动发出一个"预检"请求，要求服务器确认可以这样请求。

下面是这个"预检"请求的HTTP头信息:

```http
OPTIONS /cors HTTP/1.1
Origin: http://api.bob.com
Access-Control-Request-Method: PUT
Access-Control-Request-Headers: X-Custom-Header
Host: api.alice.com
Accept-Language: en-US
Connection: keep-alive
User-Agent: Mozilla/5.0...
```

"预检"请求用的请求方法是OPTIONS，表示这个请求是用来询问的。头信息里面，关键字段是Origin，表示请求来自哪个源。

除了Origin字段，"预检"请求的头信息包括两个特殊字段。

1. Access-Control-Request-Method。该字段是必须的，用来列出浏览器的CORS请求会用到哪些HTTP方法，上例是PUT。
2. Access-Control-Request-Headers。该字段是一个逗号分隔的字符串，指定浏览器CORS请求会额外发送的头信息字段，上例是X-Custom-Header。

### 4.6 非简单请求的回应

服务器收到"预检"请求以后，检查了Origin、Access-Control-Request-Method和Access-Control-Request-Headers字段以后，确认允许跨源请求，就可以做出回应:

```
HTTP/1.1 200 OK
Date: Mon, 01 Dec 2008 01:15:39 GMT
Server: Apache/2.0.61 (Unix)
Access-Control-Allow-Origin: http://api.bob.com
Access-Control-Allow-Methods: GET, POST, PUT
Access-Control-Allow-Headers: X-Custom-Header
Content-Type: text/html; charset=utf-8
Content-Encoding: gzip
Content-Length: 0
Keep-Alive: timeout=2, max=100
Connection: Keep-Alive
Content-Type: text/plain
```

上面的HTTP回应中，关键的是Access-Control-Allow-Origin字段，表示 http://api.bob.com 可以请求数据。该字段也可以设为星号，表示同意任意跨源请求。

如果服务器否定了"预检"请求，会返回一个正常的HTTP回应，但是没有任何CORS相关的头信息字段。

这时，浏览器就会认定，服务器不同意预检请求，因此触发一个错误，被XMLHttpRequest对象的onerror回调函数捕获。控制台会打印出如下的报错信息。

```http
XMLHttpRequest cannot load http://api.alice.com.
Origin http://api.bob.com is not allowed by Access-Control-Allow-Origin.
```

### 4.7 浏览器的正常请求和回应

一旦服务器通过了"预检"请求，以后每次浏览器正常的CORS请求，就都跟简单请求一样，会有一个Origin头信息字段。服务器的回应，也都会有一个Access-Control-Allow-Origin头信息字段。

下面是"预检"请求之后，浏览器的正常CORS请求:

```http
PUT /cors HTTP/1.1
Origin: http://api.bob.com
Host: api.alice.com
X-Custom-Header: value
Accept-Language: en-US
Connection: keep-alive
User-Agent: Mozilla/5.0...
```

上面头信息的Origin字段是浏览器自动添加的。

下面是服务器正常的回应:

```http
Access-Control-Allow-Origin: http://api.bob.com
Content-Type: text/html; charset=utf-8
```

上面头信息中，Access-Control-Allow-Origin字段是每次回应都必定包含的。

### 4.8 CORS和JSONP的对比

在2.3中提到过，CORS只是解决跨域问题的一种方案，在适合的场景下，JSONP也能解决跨域问题。

两种方案只能说互有优势：
1. JSONP只支持GET请求，CORS支持所有类型的HTTP请求。
2. JSONP的优势在于支持老式浏览器，以及可以向不支持CORS的网站请求数据。

其实只要懂得基本原理即可，实现方式可以现查现用。
