package store

type Key string
type EvictedCallback func(key Key, value Value)

func (k Key) Empty() bool {
	return k == ""
}

type Value interface {
	Len() int
}
