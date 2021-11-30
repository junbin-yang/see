package see

import (
	"fmt"
	"log"
	"net/http"
	"runtime"
	"strings"
)

func Recovery() HandlerFunc {
	return func(c *Context) {
		defer func() {
			if err := recover(); err != nil {
				message := fmt.Sprintf("%s", err)
				log.Printf("%s\n\n", trace(message))
				c.fail(http.StatusInternalServerError, "Internal Server Error")
			}
		}()

		c.Next()
	}
}

// 打印堆栈信息
func trace(message string) string {
	var pcs [32]uintptr
	n := runtime.Callers(3, pcs[:]) // 跳过上3个调用对象，函数本身、trace和defer

	var str strings.Builder
	str.WriteString(fmt.Sprintf("[Recovery] panic recovered:\n%s", message))

	frames := runtime.CallersFrames(pcs[:n])
	for {
		frame, more := frames.Next()
		str.WriteString(fmt.Sprintf("\n%s:%d\n\t%s", frame.File, frame.Line, frame.Function))
		if !more {
			break
		}
	}

	return str.String()
}
