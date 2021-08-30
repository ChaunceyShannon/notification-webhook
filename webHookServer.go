package main

import (
	"bytes"
	"io/ioutil"

	"github.com/gin-gonic/gin"
)

type messageStruct struct {
	messageHandler []string
	message        string
}

type webhookServer struct {
	g       *gin.Engine
	msgchan chan messageStruct
	bind    string
}

func (c *webhookServer) init(msgchan chan messageStruct, bind string) {
	c.msgchan = msgchan
	c.bind = bind

	// Disable debug message
	if !itemInArray(lg.levelString, []string{"trace", "debug"}) {
		gin.SetMode(gin.ReleaseMode)
		c.g = gin.New()
	} else { // Enable debug message
		c.g = gin.Default()
		c.g.Use(func(c *gin.Context) {
			buf, err := ioutil.ReadAll(c.Request.Body)
			panicerr(err)

			request := c.Request.Method + " " + c.Request.RequestURI + " " + c.Request.Proto + "\n"
			for k, v := range c.Request.Header {
				request += k + ": " + strJoin(", ", v) + "\n"
			}
			request += "\n" + str(buf)
			lg.trace("Get HTTP request: \n" + request)

			rdr := ioutil.NopCloser(bytes.NewBuffer(buf))
			c.Request.Body = rdr

			c.Next()
		})
	}

	c.g.NoRoute(func(c *gin.Context) {
		c.Header("Content-Type", "text/html; charset=utf-8")
		c.String(404, "")
	})
}

func (c *webhookServer) register(path string, msgHandler string, handler func(c *gin.Context) string) {
	msgHandlerList := strSplit(msgHandler, ",")
	c.g.POST(path, func(gc *gin.Context) {
		msg := handler(gc)
		c.msgchan <- messageStruct{
			messageHandler: msgHandlerList,
			message:        msg,
		}
	})
}

func (c *webhookServer) run() {
	lg.info("Web server started on " + c.bind)
	c.g.Run(c.bind)
}
