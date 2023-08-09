package base

import (
	"github.com/capell/capell_scan/lib/app/iapp"
	"github.com/capell/capell_scan/lib/httpbase/handler"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	handler.Handler
	App iapp.IApp
	opt *Option
}

func NewHandler(opt *Option) (bdlr *Handler, err error) {
	h := &Handler{
		opt: opt,
	}
	bdlr = h
	return
}

func (h *Handler) Route(r *gin.RouterGroup) {
	r.POST("/hello_world", h.HelloWorld)
}
