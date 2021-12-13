package see

import (
	"context"
	"github.com/cloudwego/netpoll"
	"github.com/cloudwego/netpoll-http2"
	"net/http"
	"strings"
	"sync"
	"time"
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
		debugPrint("Using port :8080 by default")
		return http.ListenAndServe(":8080", this)
	case 1:
		debugPrint("Using port", addr[0])
		return http.ListenAndServe(addr[0], this)
	default:
		panic("too many parameters")
	}
}

func debugPrint(o ...interface{}) {
	if access != nil {
		access.Println(o...)
	}
}

func (this *Engine) RunH2s(addr ...string) {
	var Addr string
	switch len(addr) {
	case 0:
		debugPrint("Http2 using port :8080 by default")
		Addr = ":8080"
	case 1:
		debugPrint("Http2 using port", addr[0])
		Addr = addr[0]
	default:
		panic("too many parameters")
	}

	// 创建 listener
	listener, err := netpoll.CreateListener("tcp", Addr)
	if err != nil {
		panic("create netpoll listener fail")
	}

	server := http2.Server{Handler: this}
	// OnRequest: 是指连接上发生读事件时触发的回调
	var onRequest netpoll.OnRequest = server.ServeConn

	// options: EventLoop 初始化自定义配置项
	opts := []netpoll.Option{
		/*
		 *  空闲超时,利用 TCP KeepAlive 机制来踢出死连接并减少维护开销。使用 Netpoll 时，一般不需要频繁创建和关闭连接，所以通常来说，空闲连接影响不大。
		 *  当连接长时间处于非活动状态时，为了防止出现假死、对端挂起、异常断开等造成的死连接，在空闲超时后，主动关闭连接。
		 */
		netpoll.WithIdleTimeout(10 * time.Minute),
		//初始化新链接,当接收新连接时，会自动执行注册的 OnPrepare 方法来完成准备工作
		netpoll.WithOnPrepare(func(conn netpoll.Connection) context.Context {
			conn.SetReadTimeout(3 * time.Minute)
			//CloseCallback 是指连接关闭时触发的回调
			var cb netpoll.CloseCallback = func(connection netpoll.Connection) error {
				debugPrint("http2 client off")
				return nil
			}
			conn.AddCloseCallback(cb)

			return context.Background()
		}),
	}

	// 创建 EventLoop
	eventLoop, err := netpoll.NewEventLoop(onRequest, opts...)
	if err != nil {
		panic("create netpoll event-loop fail")
	}

	// 运行 Server
	err = eventLoop.Serve(listener)
	if err != nil {
		panic("netpoll server exit")
	}
}
