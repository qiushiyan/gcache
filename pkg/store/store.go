package store

// abstract representation of inner cache store
type Store interface {
	Get(key Key) (Value, bool)
	Set(key Key, value Value)
	Delete(key Key) bool
}
