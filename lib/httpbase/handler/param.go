package handler

import (
	"github.com/gin-gonic/gin"

	"github.com/capell/capell_scan/lib/httpbase/reply"
)

type IParam interface {
	Validate(c *gin.Context) *reply.BaseResp
}

type Param struct {
}

func (Param) Validate(c *gin.Context) (resp *reply.BaseResp) {
	return
}
