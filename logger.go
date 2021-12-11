package see

import (
	"bytes"
	"io"
	"time"
)

func Logger() HandlerFunc {
	if access == nil {
		access = (&acc{}).New()
	}

	return func(c *Context) {
		// 开始时间
		startTime := time.Now()
		// 处理请求
		var body []byte
		if Mode() != ReleaseMode {
			body, _ = c.GetRawData()
			c.Req.Body = io.NopCloser(bytes.NewBuffer(body))
		}

		c.Next()

		// 结束时间
		endTime := time.Now()
		// 执行时间
		latencyTime := endTime.Sub(startTime)
		// 请求方式
		reqMethod := c.Method
		// 请求路由
		reqUri := c.RequestURI
		// 请求参数
		userAgent := c.Req.Header["User-Agent"]
		// 状态码
		statusCode := c.StatusCode
		// 请求IP
		clientIP := c.Req.RemoteAddr
		// 写入
		access.Writer(map[string]interface{}{"userAgent": userAgent, "clientIP": clientIP, "reqMethod": reqMethod, "reqUri": reqUri, "statusCode": statusCode, "latencyTime": latencyTime, "Body": body})
	}
}
