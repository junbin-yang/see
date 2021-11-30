package see

import (
	"net/http"
	"testing"
)

func BenchmarkOneRoute(B *testing.B) {
	router := New()
	router.GET("/ping", func(c *Context) {})
	runRequest(B, router, "GET", "/ping")
}

func BenchmarkRecoveryMiddleware(B *testing.B) {
	router := New()
	router.Use(Recovery())
	router.GET("/", func(c *Context) {})
	runRequest(B, router, "GET", "/")
}

func BenchmarkLoggerMiddleware(B *testing.B) {
	router := New()
	SetLoggerConfig("access", "/var/log/see", true, 7)
	router.Use(Logger())
	router.GET("/log", func(c *Context) {})
	runRequest(B, router, "GET", "/log")
}

func BenchmarkManyHandlers(B *testing.B) {
	router := New()
	SetLoggerConfig("access", "/var/log/see", true, 7)
	router.Use(Recovery(), Logger())
	router.Use(func(c *Context) {})
	router.Use(func(c *Context) {})
	router.GET("/ping", func(c *Context) {})
	runRequest(B, router, "GET", "/ping")
}

func Benchmark5Params(B *testing.B) {
	router := New()
	router.Use(func(c *Context) {})
	router.GET("/param/:param1/:params2/:param3/:param4/:param5", func(c *Context) {})
	runRequest(B, router, "GET", "/param/path/to/parameter/john/12345")
}

func BenchmarkOneRouteJSON(B *testing.B) {
	router := New()
	data := struct {
		Status string `json:"status"`
	}{"ok"}
	router.GET("/json", func(c *Context) {
		c.JSON(http.StatusOK, data)
	})
	runRequest(B, router, "GET", "/json")
}

func BenchmarkOneRouteSet(B *testing.B) {
	router := New()
	router.GET("/ping", func(c *Context) {
		c.Set("key", "value")
	})
	runRequest(B, router, "GET", "/ping")
}

func BenchmarkOneRouteString(B *testing.B) {
	router := New()
	router.GET("/text", func(c *Context) {
		c.String(http.StatusOK, "this is a plain text")
	})
	runRequest(B, router, "GET", "/text")
}

func BenchmarkManyRoutesFist(B *testing.B) {
	router := New()
	router.Any("/ping", func(c *Context) {})
	runRequest(B, router, "GET", "/ping")
}

func BenchmarkManyRoutesLast(B *testing.B) {
	router := New()
	router.Any("/ping", func(c *Context) {})
	runRequest(B, router, "OPTIONS", "/ping")
}

func Benchmark404(B *testing.B) {
	router := New()
	router.Any("/something", func(c *Context) {})
	router.NoRoute(func(c *Context) {})
	runRequest(B, router, "GET", "/ping")
}

func Benchmark404Many(B *testing.B) {
	router := New()
	router.GET("/", func(c *Context) {})
	router.GET("/path/to/something", func(c *Context) {})
	router.GET("/post/:id", func(c *Context) {})
	router.GET("/view/:id", func(c *Context) {})
	router.GET("/favicon.ico", func(c *Context) {})
	router.GET("/robots.txt", func(c *Context) {})
	router.GET("/delete/:id", func(c *Context) {})
	router.GET("/user/:id/:mode", func(c *Context) {})

	router.NoRoute(func(c *Context) {})
	runRequest(B, router, "GET", "/viewfake")
}

type mockWriter struct {
	headers http.Header
}

func newMockWriter() *mockWriter {
	return &mockWriter{
		http.Header{},
	}
}

func (m *mockWriter) Header() (h http.Header) {
	return m.headers
}

func (m *mockWriter) Write(p []byte) (n int, err error) {
	return len(p), nil
}

func (m *mockWriter) WriteString(s string) (n int, err error) {
	return len(s), nil
}

func (m *mockWriter) WriteHeader(int) {}

func runRequest(B *testing.B, r *Engine, method, path string) {
	// create fake request
	req, err := http.NewRequest(method, path, nil)
	if err != nil {
		panic(err)
	}
	w := newMockWriter()
	B.ReportAllocs()
	B.ResetTimer()
	for i := 0; i < B.N; i++ {
		r.ServeHTTP(w, req)
	}
}
