package handler

import (
	"github.com/gin-gonic/gin"
)

func (h *Handler) HelloWorld(c *gin.Context) {
	h.ReplyData(c, "service hello world")
	return
}
