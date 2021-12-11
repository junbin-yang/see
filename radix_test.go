package see

import (
	"fmt"
	"testing"
)

func TestTrie(t *testing.T) {
	tree := NewTree()

	routes := []string{
		"/",
		"/hi",
		"/contact",
		"/co",
		"/doc/go_faq.html",
		"/doc/go1.html",
		"/user/:name/*id",
	}
	for _, route := range routes {
		tree.Insert(route, f1)
	}

	url := []string{
		"/hi",
		"/user/zhangsan/10086",
	}
	ps := make(Params, 0, 20)
	var h HandlerFunc
	tree.Search(url[1], &ps, &h)
	h(&Context{})
	t.Log(ps)
}

func f1(c *Context) {
	fmt.Println("ok")
}
