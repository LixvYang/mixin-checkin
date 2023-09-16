package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lixvyang/mixin-checkin/internal/utils/errmsg"
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func SendResponse(c *gin.Context, code int, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    code,
		Message: errmsg.GetErrMsg(code),
		Data:    data,
	})
}
