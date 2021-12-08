package see

import (
	"github.com/junbin-yang/golib/radix"
	"net/http"
	"unsafe"
)

type route struct {
	// 存储每种请求方式的树根节点
	roots   map[string]*radix.Tree
	noRoute HandlerFunc
}

// 初始化路由
func newRoute() *route {
	r := &route{
		roots: make(map[string]*radix.Tree),
	}
	return r
}

// 注册路由
func (r *route) addRoute(method string, pattern string, handler HandlerFunc) {
	if _, ok := r.roots[method]; !ok {
		r.roots[method] = radix.PrefixRoot()
	}
	r.roots[method].Insert(pattern, handler)
}

// 获取路由，并且返回所有动态参数。
func (r *route) getRoute(method string, path string) (HandlerFunc, radix.Params) {
	root := r.roots[method]
	if root == nil {
		return nil, nil
	}

	// 在该方法的路由树上查找该路径
	params := make(radix.Params, 20)
	rt := root.Search(path, &params)
	if rt == nil {
		return nil, nil
	}
	return rt.Value.(HandlerFunc), params
}

// 找到并执行处理请求函数
func (r *route) handle(c *Context) {
	handler, params := r.getRoute(c.Method, c.Path)
	if handler != nil {
		// 将解析出来的路由参数赋值给了c.Params。这样就能够通过c.Param()访问到了
		c.Params = *(*[]Param)(unsafe.Pointer(&params))
		c.handlers = append(c.handlers, handler)
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
