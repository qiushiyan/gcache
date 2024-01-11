package store

import (
	"reflect"
	"testing"
)

func TestGetter(t *testing.T) {
	var f GetterFunc = func(key Key) (Value, error) {
		return ByteView{[]byte(key)}, nil
	}

	expect := ByteView{[]byte("key")}
	if v, _ := f.Get("key"); !reflect.DeepEqual(v, expect) {
		t.Errorf("getterfunc failed expected %s got %s", expect, v)
	}
}
