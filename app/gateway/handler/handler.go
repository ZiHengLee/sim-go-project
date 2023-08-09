package handler

import (
	"github.com/capell/capell_scan/app/gateway/handler/base"
	"github.com/capell/capell_scan/lib/logger"
	"github.com/capell/capell_scan/lib/time"
	"github.com/gin-gonic/gin"
)

type Collection struct {
	base *base.Handler
}

var Default *Collection

func NewCollection(opt *Option) (col *Collection, err error) {
	c := &Collection{}
	if opt.Base != nil {
		opt.Base.App = opt.App
		c.base, err = base.NewHandler(opt.Base)
		if err != nil {
			logger.Error("handler new base handler err:%v", err)
			return
		}
	}
	col = c
	return
}

func Init(opt *Option) (err error) {
	c, err := NewCollection(opt)
	if err != nil {
		return err
	}
	Default = c
	return
}

func (c *Collection) Route(e *gin.Engine) {
	e.OPTIONS("/*any", func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.JSON(200, map[string]interface{}{
			"code": 0,
		})
	})

	if c.base != nil {
		g := e.Group("/api/base")
		c.base.Route(g)
	}
}

func Route(e *gin.Engine) {
	Default.Route(e)
	e.GET("/healthcheck", func(c *gin.Context) {
		c.JSON(200, map[string]interface{}{
			"timestamp": time.Now().Unix(),
		})
	})
}
