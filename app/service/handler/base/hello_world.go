package base

import (
	"github.com/capell/capell_scan/lib/httpbase/handler"
	"github.com/gin-gonic/gin"
)

type NftAddTaskRq struct {
	handler.Param
}

func (h *Handler) HelloWorld(c *gin.Context) {
	h.ReplyData(c, "service hello world")
	return
}
