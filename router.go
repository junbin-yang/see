package see

import (
	"github.com/junbin-yang/see/trie"
	"net/http"
	"unsafe"
	//"strings"
)

type route struct {
	// 存储每种请求方式的树根节点
	roots map[string]*trie.Node
	// [请求方法-请求路径]处理请求的函数([]HandlerFunc最后一个才是实际处理函数，前面部分均为中间件)
	handlers map[string][]HandlerFunc

	noRoute HandlerFunc
}

// 初始化路由
func newRoute() *route {
	return &route{
		roots:    make(map[string]*trie.Node),
		handlers: make(map[string][]HandlerFunc),
	}
}

// 注册路由
func (r *route) addRoute(method string, pattern string, handler []HandlerFunc) {
	if _, ok := r.roots[method]; !ok {
		r.roots[method] = trie.NewTree()
	}

	r.roots[method].Insert(pattern)

	key := method + "-" + pattern
	r.handlers[key] = handler
}

// 获取路由，并且返回所有动态参数。例如
// 例如/p/go/doc匹配到/p/:lang/doc，解析结果为：{lang: "go"}，/static/css/h.css匹配到/static/*filepath，解析结果为{filepath: "css/h.css"}。
func (r *route) getRoute(method string, path string) (string, []trie.Param) {
	root := r.roots[method]
	if root == nil {
		return "", nil
	}

	// 在该方法的路由树上查找该路径
	params := make([]trie.Param, 0)
	for i = 0; i < trie.CountParams(path); i++ {
		params = append(params, trie.Param{})
	}
	fullpath := root.Search(path, &params)
	return fullpath, params
}

// 找到并执行处理请求函数
func (r *route) handle(c *Context) {
	fullPath, params := r.getRoute(c.Method, c.Path)
	if fullPath != "" {
		// 将解析出来的路由参数赋值给了c.Params。这样就能够通过c.Param()访问到了
		c.Params = *(*[]Param)(unsafe.Pointer(&params))
		key := c.Method + "-" + fullPath
		c.handlers = append(c.handlers, r.handlers[key]...)
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

/*
// 注册路由
func (r *route) addRoute(method string, pattern string, handler []HandlerFunc) {
	parts := parsePattern(pattern)
	if _, ok := r.roots[method]; !ok {
		r.roots[method] = &node{children: make(map[string]*node)}
	}
	r.roots[method].insert(pattern, parts)

	key := method + "-" + pattern
	r.handlers[key] = handler
}

// 获取路由，并且返回所有动态参数。例如
// 例如/p/go/doc匹配到/p/:lang/doc，解析结果为：{lang: "go"}，/static/css/h.css匹配到/static/*filepath，解析结果为{filepath: "css/h.css"}。
func (r *route) getRoute(method string, path string) (*node, []Param) {
	root := r.roots[method]
	if root == nil {
		return nil, nil
	}

	// 在该方法的路由树上查找该路径
	searchParts := parsePattern(path)
	return root.search(searchParts)
}

// 找到并执行处理请求函数
func (r *route) handle(c *Context) {
	n, params := r.getRoute(c.Method, c.Path)
	if n != nil {
		// 将解析出来的路由参数赋值给了c.Params。这样就能够通过c.Param()访问到了
		c.Params = params
		key := c.Method + "-" + n.pattern
		c.handlers = append(c.handlers, r.handlers[key]...)
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

func parsePattern(pattern string) []string {
	vs := strings.Split(pattern, "/")
	parts := make([]string, 0)
	for _, item := range vs {
		if item != "" {
			parts = append(parts, item)
			if item[0] == '*' {
				break
			}
		}
	}
	return parts
}
*/
