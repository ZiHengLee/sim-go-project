package base

import (
	"context"
	"fmt"
	"github.com/capell/capell_scan/lib/app/iapp"
	"github.com/capell/capell_scan/lib/httpbase/handler"
	swapProto "github.com/capell/capell_scan/proto/swap"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	handler.Handler
	App iapp.IApp
	opt *Option
	swapProto.UnimplementedMsgServer
}

func (h *Handler) BaseLp(ctx context.Context, req *swapProto.MsgBaseLp) (resp *swapProto.MsgBaseLpResponse, err error) {
	fmt.Println("hello")
	return &swapProto.MsgBaseLpResponse{}, nil
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
