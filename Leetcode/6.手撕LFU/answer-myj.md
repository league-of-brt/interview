```
type LFUCacheNode struct {
	k, v int
	req  int // req times
	t    int // t get key tick, 0 mean deleted
}

type LFUCache struct {
	cap  int
	len  int
	data []*LFUCacheNode // data
	dict map[int]*LFUCacheNode
	tick int // tick
}

func Constructor(capacity int) LFUCache {
	c := LFUCache{
		cap:  capacity,
		data: make([]*LFUCacheNode, 0, capacity),
		dict: make(map[int]*LFUCacheNode, capacity),
		tick: 1,
	}
	return c
}

func (this *LFUCache) Get(key int) int {
	if this.cap == 0 {
		return -1
	}
	this.tick++
	if len(this.data) == 0 {
		return -1
	}
	if v, ok := this.dict[key]; ok && v.k == key && v.t > 0 {
		v.req++
		v.t = this.tick
		return v.v
	}
	return -1
}

func (this *LFUCache) Put(key int, value int) {
	if this.cap == 0 {
		return
	}
	this.tick++
	if v, ok := this.dict[key]; ok && v.k == key && v.t > 0 {
		v.req++
		v.t = this.tick
		v.v = value
		return
	}
	var slot *LFUCacheNode
	if this.len+1 > this.cap {
		slot = this.delCache()
	}
	if slot == nil {
		slot = new(LFUCacheNode)
		this.data = append(this.data, slot)
	}
	this.len++
	slot.k = key
	slot.v = value
	slot.req = 1
	slot.t = this.tick
	this.dict[key] = slot
}

func (this *LFUCache) delCache() *LFUCacheNode {
	var slot *LFUCacheNode
	for _, node := range this.data {
		if node.t == 0 {
			continue
		}
		if slot == nil {
			slot = node
			continue
		}
		if (node.req == slot.req && node.t < slot.t) || (node.req < slot.req) {
			slot = node
		}
	}
	if slot == nil {
		panic("no slot")
	}

	slot.t = 0
	this.len--
	return slot
}
```

- map容易膨胀，有点麻烦
