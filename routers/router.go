package routers

import (
	"bytes"
	"io"
	"io/ioutil"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type Controller interface {
	Register(router *gin.RouterGroup)
}

func InitRouter() *gin.Engine {
	r := gin.Default()
	r.Use(RequestLoggerMiddleware())
	r.GET("/ping", Ping)
	return r
}

func Ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func RequestLoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var buf bytes.Buffer
		tee := io.TeeReader(c.Request.Body, &buf)
		body, _ := ioutil.ReadAll(tee)
		c.Request.Body = ioutil.NopCloser(&buf)
		log.Debug(string(body))
		log.Debug(c.Request.Header)
		c.Next()
	}
}
