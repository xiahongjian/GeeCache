package lru

import "container/list"

type Cache struct {
	// 缓存能使用的最大内存
	maxBytes int64

	// 已使用的内存
	nbytes int64

	// 使用双向链表保存元素最近使用情况
	ll *list.List

	// 使用map保存元素
	cache map[string]*list.Element

	// 缓存失效的钩子函数
	OnEvicted EvictedCallback
}

type EvictedCallback func(string, Value)

// entry 双向链表中的节点数据结构，包含key和value
type entry struct {
	// 键
	key string
	// 值
	value Value
}

// Value 使用Len函数返回占用的字节数
type Value interface {
	Len() int
}

// New 创建一个缓存
// @param maxBytes int64 缓存的最大容量，如果传入的值小于或等于0则不限制大小
// @param onEvicted EvictedCallback 缓存失效时调用的函数
// @return *Cache 缓存对象指针
func New(maxBytes int64, onEvicted EvictedCallback) *Cache {
	return &Cache{
		maxBytes:  maxBytes,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
		OnEvicted: onEvicted,
	}
}

// Get 根据key从缓存中获取数据，如果找到了返回数据，并且将此数据移动到双向链表的头部
// @receiver c *Cache
// @param key string
// @return value Value
// @return ok bool
func (c *Cache) Get(key string) (value Value, ok bool) {
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		return kv.value, true
	}
	return
}

// RemoveOldest 删除缓存中最近最少使用的数据，因此删除的是双向链表的尾部
// @receiver c *Cache
func (c *Cache) RemoveOldest() {
	ele := c.ll.Back()
	if ele != nil {
		// 从双向链表中删除
		c.ll.Remove(ele)
		kv := ele.Value.(*entry)
		// 从map中删除
		delete(c.cache, kv.key)
		// 重新计算已使用的内存
		c.nbytes -= int64(len(kv.key)) + int64(kv.value.Len())
		if c.OnEvicted != nil {
			c.OnEvicted(kv.key, kv.value)
		}
	}
}

// Add 添加/更新数，如果添加数据后使用内存大于最大可以使用内存就移除最近最少使用的数据
// @receiver c *Cache
// @param key string
// @param value Value
func (c *Cache) Add(key string, value Value) {
	if ele, ok := c.cache[key]; ok {
		// 将数据移动到链表头部
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		// 重新计算已使用内存
		c.nbytes -= int64(kv.value.Len()) - int64(value.Len())
		// 更新数据
		kv.value = value
	} else {
		// 将数据移动到链表头部
		ele := c.ll.PushFront(&entry{key, value})
		c.cache[key] = ele
		// 重新计算已使用内存
		c.nbytes += int64(len(key)) + int64(value.Len())
	}
	// 淘汰旧数据
	for c.maxBytes > 0 && c.nbytes > c.maxBytes {
		c.RemoveOldest()
	}
}

// Len 返回缓存中数据个数
func (c *Cache) Len() int {
	return c.ll.Len()
}
