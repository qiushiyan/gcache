package store

type Key string
type EvictedCallback func(key Key, value Value)

type Value interface {
	Len() int
}
