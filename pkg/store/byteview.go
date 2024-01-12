package store

// readonly view of byte slice
type ByteView struct {
	b []byte
}

func NewByteView(b []byte) ByteView {
	return ByteView{b: CloneBytes(b)}
}

func NewByteViewFromStr(s string) ByteView {
	return ByteView{b: []byte(s)}
}

func (v ByteView) ByteSlice() []byte {
	return CloneBytes(v.b)
}

func (v ByteView) String() string {
	return string(v.b)
}

func CloneBytes(b []byte) []byte {
	c := make([]byte, len(b))
	copy(c, b)
	return c
}

func (v ByteView) Len() int {
	return len(v.b)
}
