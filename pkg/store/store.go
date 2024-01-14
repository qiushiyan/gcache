package store

// abstract representation of inner cache store
type Store interface {
	Get(key Key) (v Value, ok bool)
	Set(key Key, value Value)
	Delete(key Key) bool
}
