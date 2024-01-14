package consistenthash

import (
	"fmt"
	"strconv"
	"testing"
)

func TestHash(t *testing.T) {
	ring := New(3, func(key []byte) uint32 {
		i, _ := strconv.Atoi(string(key))
		return uint32(i)
	})

	// Given the above hash function, this will generate the map
	// {2: [20, 21, 22], 4: {40, 41, 42}, 6: {60, 61, 62}}
	ring.Add("6", "4", "2")

	fmt.Println(ring.nodeMap)
	testCases := map[string]string{
		"2":  "2",
		"23": "4",
		"59": "6",
		"63": "2",
	}

	for k, v := range testCases {
		if ring.Get(k) != v {
			t.Errorf("get failed, expected %s get %s", k, v)
		}
	}

}
