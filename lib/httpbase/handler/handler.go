package handler

import (
	"github.com/gin-gonic/gin"

	"github.com/capell/capell_scan/lib/httpbase/reply"
)

type Handler struct {
}

func (h Handler) BindJson(c *gin.Context, param IParam) (resp *reply.BaseResp) {
	err := c.ShouldBindJSON(param)
	if err != nil {
		resp = &reply.BaseResp{
			Code: reply.ErrcodeParamError,
			Err:  err,
		}
		h.Reply(c, *resp, nil)
		return
	}
	resp = param.Validate(c)
	if resp != nil && (resp.Err != nil || resp.Code != 0) {
		h.Reply(c, *resp, nil, resp.Err)
		return
	}
	return
}

func (h *Handler) Reply(c *gin.Context, resp reply.BaseResp, data interface{}, args ...interface{}) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Credentials", "true")
	if resp.Err != nil || resp.Code != 0 {
		code := resp.Code
		c.JSON(200, map[string]interface{}{
			"code": code,
			"msg":  code.ServiceErrCodeName(args...),
		})
	} else {
		c.JSON(200, map[string]interface{}{
			"code": reply.ErrcodeOk,
			"data": data,
		})
	}
}

func (h *Handler) ReplyData(c *gin.Context, data interface{}) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Credentials", "true")
	c.JSON(200, map[string]interface{}{
		"code": reply.ErrcodeOk,
		"data": data,
	})
}

func (h *Handler) ReplyErr(c *gin.Context, code reply.ErrCode, args ...interface{}) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Credentials", "true")
	c.JSON(200, map[string]interface{}{
		"code": code,
		"msg":  code.ServiceErrCodeName(args...),
	})
}
