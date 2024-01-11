package store

// abstract representation of inner cache store
type Store interface {
	Create(int64, EvictedCallback) Store
	Get(key Key) (Value, bool)
	Set(key Key, value Value)
	Delete(key Key) bool
}
