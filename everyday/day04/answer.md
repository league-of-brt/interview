
CPU是按块从内存读取数据，不同的操作系统，数据块的大小不一样。  
当一份数据刚好处在数据块边界，会给cpu增加额外的操作。  
````
|----xxxx|xxxx----| // 需要两次读取内存才能拿到完整的数据
````
golang自动帮我们进行了内存对齐，64位golang将每8字节看作一个内存块，根据结构体的字段长度给结构体分配内存空间。  
但是内存对齐是一种使用空间换时间的做法，字段的类型和顺序会影响到内存的大小。

### 例如：
```
	var x struct {
		a bool
		b int32
		c int64
	}
	fmt.Println(unsafe.Sizeof(x.a))
	fmt.Println(unsafe.Sizeof(x.b))
	fmt.Println(unsafe.Sizeof(x.c))
	fmt.Println(unsafe.Sizeof(x))

	fmt.Println()

	var y struct {
		a bool
		c int64
		b int32
	}
	fmt.Println(unsafe.Sizeof(y.a))
	fmt.Println(unsafe.Sizeof(y.b))
	fmt.Println(unsafe.Sizeof(y.c))
	fmt.Println(unsafe.Sizeof(y))
```
### 输出：
```
1
4
8
16

1
4
8
24
```
unsafe.Sizeof可以字段查看占用内存大小，可以看出y比x多了8位。

内存情况：  
```
x: |x---xxxx|xxxxxxxx|
y: |x-------|xxxxxxxx|xxxx----|
```

### 优化：
* 使用适合数据类型。
* 修改字段的顺序。

### 优点：
* 方便跨平台
* 减少占用内存，在操作大量数据时比较明显
* 减少cpu操作，提高性能

### 延申：
* mysql、pg等数据库建表时，长度使用2的n次方，例如varchar(30) 改为 varchar(32)
