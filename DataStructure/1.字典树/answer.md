# 字典树 trie-tree

`4.3我认为是本节内容最难的，可能需要讨论`

* [字典树 trie\-tree](#字典树-trie-tree)
  * [1 字典树的概念](#1-字典树的概念)
    * [1\.1 怎么做？](#11-怎么做)
    * [1\.2 用什么数据结构存单词？](#12-用什么数据结构存单词)
    * [1\.3 这个数据结构有什么优点？](#13-这个数据结构有什么优点)
    * [1\.4 为什么？](#14-为什么)
    * [1\.5 可能会遇到什么问题？](#15-可能会遇到什么问题)
    * [1\.6 怎么优化？](#16-怎么优化)
  * [2 白板作图，画出字典树的简单结构](#2-白板作图画出字典树的简单结构)
  * [3 字典树的简单操作](#3-字典树的简单操作)
    * [3\.1 简单的初始化、增、删](#31-简单的初始化增删)
    * [3\.2 改操作？](#32-改操作)
  * [4 功能设计问题](#4-功能设计问题)
    * [4\.1 怎么持久化？](#41-怎么持久化)
    * [4\.2 怎么做缓存？](#42-怎么做缓存)
    * [4\.3 同时加入多个单词，怎么保证字段树的正确性？](#43-同时加入多个单词怎么保证字段树的正确性)
    * [4\.4 同时修改呢？](#44-同时修改呢)
    * [4\.5 同时删除呢？](#45-同时删除呢)

## 1 字典树的概念

现在我有一大堆单词，要做一个前缀搜索提示功能。比如我有trigger，try，time三个单词，当我依次输入前缀t、tr、tri，要找到相应的全部单词。

### 1.1 怎么做？

前缀是个很明显的提示，要想到使用字典树。

### 1.2 用什么数据结构存单词？

字典树还有更多变式，包括：
  - 压缩trie树，优化了一定的空间，但是增加维护成本。
  - 后缀树。
  - 三分字典树。

并不是一成不变的，当然这里还是使用最基本的字典树。

### 1.3 这个数据结构有什么优点？

利用字符串的公共前缀来减少查询时间，最大限度的减少无谓的字符串比较，查询效率比哈希树高。经典的利用空间换取时间的一种策略。

### 1.4 为什么？

简单的数据结构问题。
- 假设所有字符串长度之和为n，构建字典树的时间复杂度为O(n)。
- 假设要查找的字符串长度为k，查找的时间复杂度为O(k)。

### 1.5 可能会遇到什么问题？

字典树每个节点都需要用一个数组来存储子节点的指针，即便实际只有两三个子节点，但依然需要一个完整大小的数组。所以，字典树比较耗内存，空间复杂度较高。

### 1.6 怎么优化？

- 可以牺牲一点查询的效率。将每个节点的子节点数组用其他数据结构代替。例如有序数组，红黑树，散列表等。例如，当子节点数组采用有序数组时，可以使用二分查找来查找下一个字符。

- 缩点优化。将末尾一些只有一个子节点的节点，可以进行合并，但是增加了编码的难度。

## 2 白板作图，画出字典树的简单结构

![trie-tree](https://brt-1303999354.cos.ap-shanghai.myqcloud.com/QQ%E6%88%AA%E5%9B%BE20201220152453.png)

## 3 字典树的简单操作

### 3.1 简单的初始化、增、删

定义数据结构TrieNode：

```java
package com.xie.leetcode.struct;

import lombok.Data;

import java.util.HashMap;

/**
 * @author xie4ever
 * @date 2020/12/20 11:10
 */
@Data
public class TrieNode {
    // 当前节点储存的字符
    public char val;
    // 当前节点为止的字符串
    public String string;
    // 当前节点是否一个字符串的结尾
    public boolean isEnd;
    // 字符串的下一个节点
    public HashMap<Character, TrieNode> next = new HashMap<>();

    public TrieNode(char val, String string) {
        this.val = val;
        this.string = string;
    }
}
```

实际操作：
```java
package com.xie.leetcode.dataStructure;

import com.xie.leetcode.struct.TrieNode;

import java.util.ArrayList;
import java.util.List;
import java.util.Stack;

/**
 * @author xie4ever
 * @date 2020/12/20 11:18
 */
public class TrieTree {
    // 根节点
    TrieNode root;

    // 初始化根节点
    public TrieTree() {
        root = new TrieNode(' ', "");
    }

    // 全文搜索是否存在
    public boolean isFullWordExist(String word) {
        TrieNode current = root;
        for (int i = 0; i < word.length(); i++) {
            char c = word.charAt(i);
            if (!current.next.containsKey(c)) {
                return false;
            } else {
                current = current.next.get(c);
            }
        }
        // 最后一个节点是否标志为单词
        return current.isEnd;
    }

    // 前缀搜索是否存在
    public boolean isPrefixWordSelect(String word) {
        TrieNode current = root;
        for (int i = 0; i < word.length(); i++) {
            char c = word.charAt(i);
            if (!current.next.containsKey(c)) {
                return false;
            } else {
                current = current.next.get(c);
            }
        }
        // 只要有节点，下面一定有单词
        return true;
    }

    // 返回最深的前缀节点，没有就返回null
    public TrieNode getDeepestPrefixNode(String word) {
        TrieNode current = root;
        for (int i = 0; i < word.length(); i++) {
            char c = word.charAt(i);
            if (!current.next.containsKey(c)) {
                return null;
            } else {
                current = current.next.get(c);
            }
        }
        return current;
    }

    // 返回前缀搜索的所有单词，没有就返回空列表
    public List<String> listPrefixWord(String word) {
        List<String> result = new ArrayList<>();
        if (word.isEmpty()) {
            return result;
        }
        // 先找到前缀节点
        TrieNode prefixNode = getDeepestPrefixNode(word);
        if (prefixNode == null) {
            return result;
        }
        // 遍历前缀节点下面的所有节点
        return dfs4ListPrefixWord(result, prefixNode);
    }

    private List<String> dfs4ListPrefixWord(List<String> result, TrieNode node) {
        if (node == null) {
            return result;
        }
        if (node.isEnd) {
            result.add(node.string);
        }
        for (Character c : node.next.keySet()) {
            result = dfs4ListPrefixWord(result, node.next.get(c));
        }
        return result;
    }

    // 插入
    public void insert(String word) {
        TrieNode current = root;
        for (int i = 0; i < word.length(); i++) {
            char c = word.charAt(i);
            if (!current.next.containsKey(c)) {
                current.next.put(c, new TrieNode(c, word.substring(0, i + 1)));
            }
            current = current.next.get(c);
        }
        // 最后一个节点标志为单词
        current.isEnd = true;
    }

    // 删除
    public void delete(String word) {
        if (word.isEmpty()) {
            return;
        }
        Stack<TrieNode> stack = new Stack<>();
        TrieNode current = root;
        for (int i = 0; i < word.length(); i++) {
            char c = word.charAt(i);
            if (!current.next.containsKey(c)) {
                return;
            } else {
                stack.push(current);
                current = current.next.get(c);
            }
        }

        if (!current.next.isEmpty()) {
            current.isEnd = false;
            return;
        } else {
            // 处理父节点
            stack.peek().next.put(current.val, null);
        }

        while (!stack.empty()) {
            if (stack.peek() == root) {
                return;
            }
            TrieNode node = stack.pop();
            boolean flag = true;
            for (Character c : node.next.keySet()) {
                if (node.next.get(c) != null) {
                    flag = false;
                }
            }
            if (flag) {
                stack.peek().next.put(node.val, null);
            }
        }
    }

    public static void main(String[] args) {
        TrieTree root = new TrieTree();
        root.insert("ap");
        root.insert("app");
        root.insert("apk");

        List<String> list = root.listPrefixWord("ap");
        for (String s : list) {
            System.out.println(s);
        }

        root.delete("bb");

        list = root.listPrefixWord("ap");
        for (String s : list) {
            System.out.println(s);
        }

        root.delete("apk");

        list = root.listPrefixWord("ap");
        for (String s : list) {
            System.out.println(s);
        }
    }
}
```

### 3.2 改操作？

个人认为，不存在什么改操作，直接删掉新增就可以。

## 4 功能设计问题

### 4.1 怎么持久化？

- 我个人觉得这个是看规模的问题。设计一个功能之前永远先估计数据量。
- 像我现在这样的个人博客，前缀关键字可能在1w+，那么我觉得存在库里问题不大，普通存储都能搞定，随时还原成trie树就行。
- 如果是工业级别的前缀搜索，肯定有推荐系统之类的东西做支持。我个人可能会先把关键字分布式存起来，然后还原成很多trie树。如果有关键字进来，就做个hash，分配给相关的trie树做处理。
- 这里还涉及一个问题，如果让你去设计一个储存系统，你怎么设计。

### 4.2 怎么做缓存？

- 像我现在这样的静态博客，甚至能整个trie索引文件，完全交给前端去做搜索。缺点就是索引文件比较大容易耗带宽，修改不方便。
- 工业级别肯定使用Redis集群来顶，这里就牵涉到Redis的问题。

### 4.3 同时加入多个单词，怎么保证字段树的正确性？

- 既然是并发问题，就绕不开消息队列和锁。
- 增加单词，就扔到消息队列去处理。
- 必须加锁。如果不加锁，我先插入一个apple，再插入一个apples，可能apples的isEnd=false会把apple的e节点的isEnd=true给顶掉，就找不到apple这个词了。
- 所以插入apple时候，必须先加分布式锁，就锁住a、ap、app、appl、apple。如果apples进来，必须去拿a、ap、app、appl、apple、apples锁（锁住一整条路径）。这样会不会导致竞争太激烈？我也没想到更好的解决方法。
- 注意这里的死锁问题。比如app和apple同时进来，app拿到了a锁，想要请求ap锁。apple拿到了ap锁，想要请求a锁，这样肯定就死锁了。所以这里的拿锁肯定是要有顺序的，apple必须先拿到a锁再去拿ap锁。
- 我个人认为这里可以合并任务。比如设计一个任务等待机制，一分钟内的请求一起更新一次trie树。这样apple和apples的任务就能合并到一起，有效减少锁竞争。
- 这里的锁肯定要有优先级，有超时机制，又可以拓展。

### 4.4 同时修改呢？

不存在修改，所以不存在同时修改问题。

### 4.5 同时删除呢？

* 首先这里要想好，为什么要删除节点？如果不想让一个推荐词出现，直接加个黑名单就好了。
* 如果要清理trie树，肯定要先计划好清理哪些东西，找个相对低峰期，然后复制一份线上的，清理，同步，增量处理。
