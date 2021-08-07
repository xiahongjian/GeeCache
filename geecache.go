package geecache

type Getter interface {
	Get(key string) (ByteView, error)
}

type GetterFunc func(key string) (ByteView, error)

func (f GetterFunc) Get(key string) (ByteView, error) {
	return f(key)
}
