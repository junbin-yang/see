package see

import (
	"net/http"
	"strings"
)

const defaultMultipartMemory = 32 << 20 // 32 MB
const ReleaseMode = "release"
const DebugMode = "debug"

var mode = DebugMode
var access *acc
var DebugPrintRouteFunc func(httpMethod, absolutePath, handlerName string)

func Mode() string {
	return mode
}

func SetMode(value string) {
	if mode != value {
		switch value {
		case DebugMode, "":
			mode = DebugMode
		case ReleaseMode:
			mode = ReleaseMode
		default:
			panic("mode unknown: " + value)
		}
	}
}

type HandlerFunc func(c *Context)

type Engine struct {
	// Engine继承routerGroup所有属性和方法
	*routerGroup
	router *route
	groups []*routerGroup

	// Value of 'maxMemory' param that is given to http.Request's ParseMultipartForm method call.
	MaxMultipartMemory int64
}

func New() *Engine {
	engine := &Engine{router: newRoute()}
	engine.routerGroup = &routerGroup{engine: engine}
	engine.groups = []*routerGroup{engine.routerGroup}
	engine.MaxMultipartMemory = defaultMultipartMemory
	return engine
}

// filename = path + "/" + prefixname +"-2020-01-01.log"
func Enable(prefixname, path string, rotate bool, keepdays int64) *Engine {
	SetLoggerConfig(prefixname, path, rotate, keepdays)
	engine := New()
	engine.Use(Logger(), Recovery())
	return engine
}

func SetLoggerConfig(prefixname, path string, rotate bool, keepdays int64) {
	access = (&acc{FileName: prefixname, Path: path, Rotate: rotate, KeepDays: keepdays}).New()
}

func Default() *Engine {
	access = (&acc{}).New()
	engine := New()
	engine.Use(Logger(), Recovery())
	return engine
}

func (this *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 找出所有路由的中间件函数
	var middlewares []HandlerFunc
	for _, group := range this.groups {
		if strings.HasPrefix(r.URL.Path, group.prefix) {
			middlewares = append(middlewares, group.middlewares...)
		}
	}

	c := NewContext(w, r)
	c.handlers = middlewares
	c.engine = this
	this.router.handle(c)
}

// 找不到路由时的回调
func (this *Engine) NoRoute(handler HandlerFunc) {
	this.router.noRoute = handler
}

func (this *Engine) Run(addr ...string) (err error) {
	switch len(addr) {
	case 0:
		access.Println("Using port :8080 by default")
		return http.ListenAndServe(":8080", this)
	case 1:
		access.Println("Using port", addr[0])
		return http.ListenAndServe(addr[0], this)
	default:
		panic("too many parameters")
	}
}
