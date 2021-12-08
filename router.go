package see

import (
	"github.com/junbin-yang/golib/radix"
	"net/http"
	"sync"
	"unsafe"
)

type root struct {
	method string
}

type route struct {
	// 存储每种请求方式的树根节点
	roots   []*radix.Tree
	noRoute HandlerFunc
	pool    sync.Pool
}

// 初始化路由
func newRoute() *route {
	r := &route{
		roots: make([]*radix.Tree, len(anyMethods)),
	}
	r.pool.New = func() interface{} {
		params := make(radix.Params, 20)
		return &params
	}
	return r
}

func (r *route) getRootIndex(method string) int {
	switch method {
	case http.MethodConnect:
		return 0
	case http.MethodDelete:
		return 1
	case http.MethodGet:
		return 2
	case http.MethodHead:
		return 3
	case http.MethodOptions:
		return 4
	case http.MethodPatch:
		return 5
	case http.MethodPost:
		return 6
	case http.MethodPut:
		return 7
	case http.MethodTrace:
		return 8
	}
	return -1
}

// 注册路由
func (r *route) addRoute(method string, pattern string, handler HandlerFunc) {
	index := r.getRootIndex(method)
	if index == -1 {
		panic("The http method error")
	}
	if r.roots[index] == nil {
		r.roots[index] = radix.PrefixRoot()
	}

	r.roots[index].Insert(pattern, handler)
}

// 获取路由，并且返回所有动态参数。
func (r *route) getRoute(method string, path string) (HandlerFunc, radix.Params, int) {
	index := r.getRootIndex(method)
	root := r.roots[index]
	if root == nil {
		return nil, nil, 0
	}

	// 在该方法的路由树上查找该路径
	params := r.pool.Get().(*radix.Params)
	paramIndex := 0
	rt := root.Search(path, params, &paramIndex)
	if rt == nil {
		return nil, nil, 0
	}

	return rt.Value.(HandlerFunc), *params, paramIndex
}

// 找到并执行处理请求函数
func (r *route) handle(c *Context) {
	handler, params, index := r.getRoute(c.Method, c.Path)
	if handler != nil {
		// 将解析出来的路由参数赋值给了c.Params。这样就能够通过c.Param()访问到了
		c.Params = (*(*[]Param)(unsafe.Pointer(&params)))[:index]
		for i := 0; i < index-1; i++ {
			params[i] = radix.Param{}
		}
		r.pool.Put(&params)

		c.handlers[len(c.handlers)-1] = handler
	} else {
		// 没有匹配到路由
		if r.noRoute == nil {
			c.handlers = append(c.handlers, func(c *Context) {
				c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
			})
		} else {
			c.handlers = append(c.handlers, r.noRoute)
		}
	}
	c.Next()
}
