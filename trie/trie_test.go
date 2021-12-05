package trie

import (
	"testing"
)

func TestTrie(t *testing.T) {
	tree := NewTree()

	routes := []string{
		"/",
		"/hi",
		"/contact",
		"/co",
		"/c",
		"/a",
		"/ab",
		"/doc/",
		"/doc/go_faq.html",
		"/doc/go1.html",
		"/α",
		"/β",
		"/user/:name/*id",
	}
	for _, route := range routes {
		tree.Insert(route)
	}

	url := []string{
		"/hi",
		"/doc/go1.html",
		"/user/zhangsan/100",
		"/",
	}
	for _, i := range url {
		t.Log(tree.Search(i))
	}
}
