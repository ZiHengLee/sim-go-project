package base

import (
	"fmt"
	"github.com/capell/capell_scan/lib/httpbase/handler"
	swapProto "github.com/capell/capell_scan/proto/base"
	swapRpc "github.com/capell/capell_scan/rpc/base"
	"github.com/gin-gonic/gin"
)

type NftAddTaskRq struct {
	handler.Param
}

func (h *Handler) HelloWorld(c *gin.Context) {
	var req swapProto.MsgBaseLp
	r, err := swapRpc.BaseLp(c, &req)
	fmt.Println(r, err)
	h.ReplyData(c, nil)
	return
}
