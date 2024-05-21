package lru

import (
	"fmt"
	"testing"
)

func TestLru(t *testing.T) {
	lru := NewLRU(4)
	lru.Insert("a", 1)
	lru.Insert("b", "b")
	lru.Insert("c", 123434)
	lru.dl.Show()
	fmt.Printf("get a: %v\n", lru.Get("a"))
	lru.dl.Show()
	lru.Insert("c", 12343433666)
	lru.dl.Show()
}
