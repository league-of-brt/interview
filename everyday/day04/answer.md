
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
	fmt.Println("SIZE")
	fmt.Println(unsafe.Sizeof(x))
	fmt.Println(unsafe.Sizeof(x.a))
	fmt.Println(unsafe.Sizeof(x.b))
	fmt.Println(unsafe.Sizeof(x.c))
	fmt.Println("Alignof")
	fmt.Println(unsafe.Alignof(x.a))
	fmt.Println(unsafe.Alignof(x.b))
	fmt.Println(unsafe.Alignof(x.c))
	fmt.Println("Offsetof")
	fmt.Println(unsafe.Offsetof(x.a))
	fmt.Println(unsafe.Offsetof(x.b))
	fmt.Println(unsafe.Offsetof(x.c))

	fmt.Println()

	var y struct {
		a bool
		c int64
		b int32
	}
	fmt.Println("SIZE")
	fmt.Println(unsafe.Sizeof(y))
	fmt.Println(unsafe.Sizeof(y.a))
	fmt.Println(unsafe.Sizeof(y.b))
	fmt.Println(unsafe.Sizeof(y.c))
	fmt.Println("Alignof")
	fmt.Println(unsafe.Alignof(y.a))
	fmt.Println(unsafe.Alignof(y.b))
	fmt.Println(unsafe.Alignof(y.c))
	fmt.Println("Offsetof")
	fmt.Println(unsafe.Offsetof(y.a))
	fmt.Println(unsafe.Offsetof(y.b))
	fmt.Println(unsafe.Offsetof(y.c))
```
### 输出：
```
SIZE
16
1
4
8
Alignof
1
4
8
Offsetof
0
4
8

SIZE
24
1
4
8
Alignof
1
4
8
Offsetof
0
16
8
```
unsafe.Sizeof函数返回操作数在内存中的字节大小。  
unsafe.Alignof 函数返回对应参数的类型需要对齐的倍数。  
unsafe.Offsetof函数的参数必须是一个字段 x.f, 然后返回 f 字段相对于 x 起始地址的偏移量, 包括可能的空洞。  

内存情况：  
```
x: |x---xxxx|xxxxxxxx|
y: |x-------|xxxxxxxx|xxxx----|
```
可以看出y比x多了8位。

### 优化：
* 使用适合数据类型。
* 修改字段的顺序。

### 优点：
* 方便跨平台
* 减少占用内存，在操作大量数据时比较明显
* 减少cpu操作，提高性能

### 延申：
* mysql、pg等数据库建表时，长度使用2的n次方，例如varchar(30) 改为 varchar(32)
* mysql、pg等数据库建表时，长度使用2的n次方，例如varchar(30) 改为 varchar(32)内存
* mysql、pg等数据库建表时，长度使用2的n次方，例如varchar(30) 改为 varchar(32)
* go的内存逃逸和内存回收
