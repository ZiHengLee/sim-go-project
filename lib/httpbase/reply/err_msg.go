package reply

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

var msgFunc func(c *gin.Context, code int, args ...interface{}) (msg string)

func SetMsgFunc(f func(c *gin.Context, code int, args ...interface{}) (msg string)) {
	msgFunc = f
}

func Msg(c *gin.Context, code int, args ...interface{}) (msg string) {
	if msgFunc != nil {
		return msgFunc(c, code, args...)
	}
	if len(args) == 1 {
		msg = fmt.Sprintf("%v", args[0])
	} else if len(args) > 1 {
		msg = fmt.Sprintf("%v", args)
	}
	return
}
