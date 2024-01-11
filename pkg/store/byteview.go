package store

// readonly view of byte slice
type ByteView struct {
	B []byte
}

func (v ByteView) String() string {
	return string(v.B)
}

func CloneBytes(b []byte) []byte {
	c := make([]byte, len(b))
	copy(c, b)
	return c
}

func (v ByteView) Len() int {
	return len(v.B)
}
