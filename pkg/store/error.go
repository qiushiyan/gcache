package store

import "fmt"

type errKind int

const (
	_ errKind = iota
	KeyEmpty
	GetterError
)

var (
	ErrKeyEmpty = Error{kind: KeyEmpty}
	ErrGetter   = Error{kind: GetterError}
)

type Error struct {
	kind  errKind
	value string
	err   error
}

func (e Error) Error() string {
	switch e.kind {
	case KeyEmpty:
		return "key must be non-empty"
	case GetterError:
		return fmt.Sprintf("getter error: %v", e.err)
	default:
		return e.err.Error()
	}
}

func (e Error) With(err error) Error {
	e.err = err
	return e
}
