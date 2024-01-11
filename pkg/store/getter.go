package store

// A Getter loads data for a key.
type Getter interface {
	Get(key Key) (Value, error)
}

// A GetterFunc implements Getter with a function.
type GetterFunc func(key Key) (Value, error)

// Get implements Getter interface function
func (f GetterFunc) Get(key Key) (Value, error) {
	if f == nil {
		return nil, nil
	}
	return f(key)
}
