package http

import (
	"fmt"
	"net/http/httputil"
	"os"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
)

func loggerHandler(c *gin.Context) {
	// Start timer
	start := time.Now()
	path := c.Request.URL.Path
	raw := c.Request.URL.RawQuery
	method := c.Request.Method

	// Process request
	c.Next()

	// Stop timer
	end := time.Now()
	latency := end.Sub(start)
	statusCode := c.Writer.Status()
	errCode := c.GetInt(ErrCodeKey)
	clientIP := c.ClientIP()
	if raw != "" {
		path = path + "?" + raw
	}
	glog.Infof("METHOD:%s | PATH:%s | CODE:%d | IP:%s | TIME:%d | ECODE:%d", method, path, statusCode, clientIP, latency/time.Millisecond, errCode)
}

func recoverHandler(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			const size = 64 << 10
			buf := make([]byte, size)
			buf = buf[:runtime.Stack(buf, false)]
			httprequest, _ := httputil.DumpRequest(c.Request, false)
			pnc := fmt.Sprintf("[Recovery] %s panic recovered:\n%s\n%s\n%s", time.Now().Format("2006-01-02 15:04:05"), string(httprequest), err, buf)
			fmt.Fprintf(os.Stderr, pnc)
			glog.Error(pnc)
			c.AbortWithStatus(500)
		}
	}()
	c.Next()
}
