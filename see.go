package see

import (
	"net/http"
	"strings"
	"sync"
)

const defaultMultipartMemory = 32 << 20 // 32 MB
const ReleaseMode = "release"
const DebugMode = "debug"

var mode = DebugMode
var access *acc
var DebugPrintRouteFunc func(httpMethod, absolutePath, handlerName string)
var maxParams int

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
	// context的临时对象池
	pool sync.Pool

	maxParams      uint8
	maxMiddlewares uint8
}

func New(opt ...uint8) *Engine {
	engine := &Engine{router: newRoute()}
	engine.routerGroup = &routerGroup{engine: engine}
	engine.groups = []*routerGroup{engine.routerGroup}
	engine.MaxMultipartMemory = defaultMultipartMemory
	switch len(opt) {
	case 1:
		engine.maxParams = opt[0]
	case 2:
		engine.maxParams = opt[0]
		engine.maxMiddlewares = opt[1]
	default:
		engine.maxParams = 20
		engine.maxMiddlewares = 100
	}
	engine.pool.New = func() interface{} {
		return &Context{index: -1, Params: make(Params, 0, engine.maxParams), handlers: make([]HandlerFunc, 0, engine.maxMiddlewares)}
	}
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
	// 取一个临时Context对象
	c := this.pool.Get().(*Context)
	c.SetContext(w, r)

	// 找出所有路由的中间件函数
	for _, group := range this.groups {
		if strings.HasPrefix(r.URL.Path, group.prefix) {
			if group.private == false || (group.private && (r.URL.Path == group.prefix)) {
				//c.handlers = append(c.handlers, group.middlewares...)
				for _, middle := range group.middlewares {
					i := len(c.handlers)
					c.handlers = (c.handlers)[:i+1]
					(c.handlers)[i] = middle
				}
			}
		}
	}

	c.engine = this
	this.router.handle(c)

	// 重置标记后放回对象池
	c.Reset()
	this.pool.Put(c)
}

// 找不到路由时的回调
func (this *Engine) NoRoute(handler HandlerFunc) {
	this.router.noRoute = handler
}

func (this *Engine) Run(addr ...string) (err error) {
	switch len(addr) {
	case 0:
		//access.Println("Using port :8080 by default")
		return http.ListenAndServe(":8080", this)
	case 1:
		//access.Println("Using port", addr[0])
		return http.ListenAndServe(addr[0], this)
	default:
		panic("too many parameters")
	}
}
