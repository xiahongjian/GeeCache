package geecache

// ByteView 为缓存数据的包装结构
type ByteView struct {
	b []byte
}

// Len 返回内部字节数组的长度
// @receiver v ByteView
// @return int 字节数组的长度
func (v ByteView) Len() int {
	return len(v.b)
}

// ByteSlice 返回内部数据的一个副本
// @receiver v ByteView
// @return []byte 返回内部字节数组的slice
func (v ByteView) ByteSlice() []byte {
	return cloneBytes(v.b)
}

// String 返回内部数据的字符串形式
// @receiver v ByteView
// @return string
func (v ByteView) String() string {
	return string(v.b)
}

// cloneBytes 克隆所给字节数组的数据，并返回
// @param b []byte
// @return []byte
func cloneBytes(b []byte) []byte {
	dump := make([]byte, len(b))
	copy(dump, b)
	return dump
}
