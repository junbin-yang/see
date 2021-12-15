package see

import (
	"net/http"
	"path"
	"reflect"
	"runtime"
	"strings"
)

var anyMethods []string = []string{
	http.MethodConnect, http.MethodDelete, http.MethodGet,
	http.MethodHead, http.MethodOptions, http.MethodPatch,
	http.MethodPost, http.MethodPut, http.MethodTrace,
}

type routerGroup struct {
	prefix      string        // 路由分组Url
	middlewares []HandlerFunc // 中间件
	engine      *Engine
	private     bool // 单个路由私有的分组，支持单路由中间件
}

func (this *routerGroup) Group(prefix string) *routerGroup {
	engine := this.engine
	newGroup := &routerGroup{
		prefix: this.prefix + prefix, // 上一个路由分组前缀加下一个
		engine: engine,
	}
	engine.groups = append(engine.groups, newGroup)
	return newGroup
}

// 中间件实现
func (this *routerGroup) Use(middlewares ...HandlerFunc) {
	this.middlewares = append(this.middlewares, middlewares...)
}

// 注册路由
func (this *routerGroup) addRoute(method string, pattern string, handler []HandlerFunc) {
	l := len(handler)
	lastHandler := handler[l-1]
	// 处理单路由中间件
	if l > 1 {
		engine := this.engine
		engine.groups = append(
			engine.groups,
			&routerGroup{
				prefix:      this.prefix + pattern, // 上一个路由分组前缀加下一个
				engine:      engine,
				private:     true,
				middlewares: handler[:l-1],
			},
		)
	}
	printRoute(method, this.prefix+pattern, lastHandler)
	if method == "Any" {
		for _, method = range anyMethods {
			this.engine.router.addRoute(method, this.prefix+pattern, lastHandler)
		}
		return
	}
	this.engine.router.addRoute(method, this.prefix+pattern, lastHandler)
}

func printRoute(httpMethod, absolutePath string, handler HandlerFunc) {
	if DebugPrintRouteFunc == nil {
		if access != nil {
			access.PrintRoute(httpMethod, absolutePath)
		}
	} else {
		handlerName := runtime.FuncForPC(reflect.ValueOf(handler).Pointer()).Name()
		DebugPrintRouteFunc(httpMethod, absolutePath, handlerName)
	}
}

func (this *routerGroup) Handle(method, pattern string, handler ...HandlerFunc) {
	this.addRoute(method, pattern, handler)
}

func (this *routerGroup) Any(pattern string, handler ...HandlerFunc) {
	this.addRoute("Any", pattern, handler)
}

func (this *routerGroup) GET(pattern string, handler ...HandlerFunc) {
	this.addRoute("GET", pattern, handler)
}

func (this *routerGroup) POST(pattern string, handler ...HandlerFunc) {
	this.addRoute("POST", pattern, handler)
}

func (this *routerGroup) PUT(pattern string, handler ...HandlerFunc) {
	this.addRoute("PUT", pattern, handler)
}

func (this *routerGroup) DELETE(pattern string, handler ...HandlerFunc) {
	this.addRoute("DELETE", pattern, handler)
}

func (this *routerGroup) PATCH(pattern string, handler ...HandlerFunc) {
	this.addRoute("PATCH", pattern, handler)
}

func (this *routerGroup) HEAD(pattern string, handler ...HandlerFunc) {
	this.addRoute("HEAD", pattern, handler)
}

func (this *routerGroup) OPTIONS(pattern string, handler ...HandlerFunc) {
	this.addRoute("OPTIONS", pattern, handler)
}

// 静态文件实现
func (this *routerGroup) StaticFile(relativePath, filepath string) {
	if strings.Contains(relativePath, ":") || strings.Contains(relativePath, "*") {
		panic("Dynamic parameters cannot be used when serving static files")
	}

	handler := func(c *Context) {
		c.File(filepath)
	}
	this.GET(relativePath, handler)
	this.HEAD(relativePath, handler)
}

func (this *routerGroup) Static(relativePath string, root string) {
	// 动态参数不能在静态文件系统里使用
	if strings.Contains(relativePath, ":") || strings.Contains(relativePath, "*") {
		panic("Dynamic parameters cannot be used when serving static files")
	}

	handler := this.createStaticHandler(relativePath, http.Dir(root))
	urlPattern := path.Join(relativePath, "/*filepath")
	// 注册方法
	this.GET(urlPattern, handler)
	this.HEAD(urlPattern, handler)
}

func (this *routerGroup) createStaticHandler(relativePath string, fs http.FileSystem) HandlerFunc {
	absolutePath := path.Join(this.prefix, relativePath)
	fileServer := http.StripPrefix(absolutePath, http.FileServer(fs))
	return func(c *Context) {
		file := c.Param("filepath")
		// 检查文件是否有权限打开
		if _, err := fs.Open(file); err != nil {
			c.Status(http.StatusNotFound)
			return
		}
		c.StatusCode = http.StatusOK
		fileServer.ServeHTTP(c.Writer, c.Req)
	}
}

type RESTful interface {
	Create(*Context)
	Query(*Context)
	Update(*Context)
	Delete(*Context)
}

func (this *routerGroup) REST(pattern string, rt RESTful) {
	this.GET(pattern, rt.Query)
	this.POST(pattern, rt.Create)
	this.PUT(pattern, rt.Update)
	this.DELETE(pattern, rt.Delete)
}
