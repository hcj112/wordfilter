package http

import (
	"github.com/gin-gonic/gin"
)

const (
	// OK ok
	OK = 0
	// RequestErr request error
	RequestErr = -400
	// ServerErr server error
	ServerErr = -500

	ErrCodeKey = "err/code"
)

type response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func errors(c *gin.Context, code int, msg string) {
	c.Set(ErrCodeKey, code)
	c.JSON(200, response{
		Code:    code,
		Message: msg,
	})
}

func result(c *gin.Context, data interface{}, code int) {
	c.Set(ErrCodeKey, code)
	c.JSON(200, response{
		Code: code,
		Data: data,
	})
}
