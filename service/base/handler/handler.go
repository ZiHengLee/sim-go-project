package handler

import (
	"context"
	"fmt"
	"github.com/capell/capell_scan/lib/app/iapp"
	"github.com/capell/capell_scan/lib/httpbase/handler"
	baseProto "github.com/capell/capell_scan/proto/base"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	handler.Handler
	App iapp.IApp
	opt *Option
	baseProto.UnimplementedMsgServer
}

func NewHandler(opt *Option) (bdlr *Handler, err error) {
	h := &Handler{
		opt: opt,
	}
	if opt.NeedGrpc {
		fmt.Println("register grpc")
		baseProto.RegisterMsgServer(opt.App.GetGrpcServer(), h)
	}
	bdlr = h
	return
}

func (h *Handler) BaseLp(context.Context, *baseProto.MsgBaseLp) (*baseProto.MsgBaseLpResponse, error) {
	fmt.Println("base lp")
	return &baseProto.MsgBaseLpResponse{}, nil
}

func (h *Handler) Route(r *gin.RouterGroup) {
	r.POST("/hello_world", h.HelloWorld)
}
