package see

import (
	"net/http"
	"strings"
)

type route struct {
	// 存储每种请求方式的树根节点
	roots []*node
	// [请求方法-请求路径]处理请求的函数([]HandlerFunc最后一个才是实际处理函数，前面部分均为中间件)
	handlers map[string][]HandlerFunc

	noRoute HandlerFunc
}

// 初始化路由
func newRoute() *route {
	return &route{
		roots:    make([]*node, len(anyMethods)),
		handlers: make(map[string][]HandlerFunc),
	}
}

// 注册路由
func (r *route) addRoute(method string, pattern string, handler []HandlerFunc) {
	parts := parsePattern(pattern)
	key := method + "-" + pattern
	arrid, status := anyMethods.GetKey(method)
	if status == false {
		panic("Error registering route method")
	}
	if r.roots[arrid] == nil {
		r.roots[arrid] = &node{}
	}
	r.roots[arrid].insert(pattern, parts, 0)

	r.handlers[key] = handler
}

// 获取路由，并且返回所有动态参数。例如
// 例如/p/go/doc匹配到/p/:lang/doc，解析结果为：{lang: "go"}，/static/css/h.css匹配到/static/*filepath，解析结果为{filepath: "css/h.css"}。
func (r *route) getRoute(mid int32, path string) (*node, []Param) {
	searchParts := parsePattern(path)
	var params []Param
	root := r.roots[mid]
	if root == nil {
		return nil, nil
	}

	n := root.search(searchParts, 0)
	if n != nil {
		parts := parsePattern(n.pattern)
		for index, part := range parts {
			if part[0] == ':' {
				params = append(params, Param{part[1:], searchParts[index]})
			}
			if part[0] == '*' && len(part) > 1 {
				params = append(params, Param{part[1:], strings.Join(searchParts[index:], "/")})
				break
			}
		}
		return n, params
	}
	return nil, nil
}

// 找到并执行处理请求函数
func (r *route) handle(c *Context) {
	n, params := r.getRoute(c.methodId, c.Path)
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

// 只能匹配一个*
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
