package store

// readonly view of byte slice
type ByteView struct {
	b []byte
}

func cloneBytes(b []byte) []byte {
	c := make([]byte, len(b))
	copy(c, b)
	return c
}

func (v ByteView) Len() int {
	return len(v.b)
}

func (v ByteView) B() []byte {
	return cloneBytes(v.b)
}
