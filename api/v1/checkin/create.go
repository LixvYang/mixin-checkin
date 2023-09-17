package checkin

import (
	"github.com/gin-gonic/gin"
	v1 "github.com/lixvyang/mixin-checkin/api/v1"
	"github.com/lixvyang/mixin-checkin/internal/logic"
	"github.com/lixvyang/mixin-checkin/internal/model"
	"github.com/lixvyang/mixin-checkin/internal/utils/errmsg"
	"github.com/rs/zerolog"
)

// 插入签到表 目前只有我一个用户
func CreateCheckIn(c *gin.Context) {
	xl := c.MustGet("logger").(*zerolog.Logger)
	checkInReq := new(model.CheckinReq)
	if err := c.ShouldBindJSON(checkInReq); err != nil {
		xl.Error().Err(err).Msg("c.ShouldBindJSON(checkInReq) error")
		v1.SendResponse(c, errmsg.ERROR_BIND, nil)
		return
	}

	// 业务处理
	if err := logic.CreateCheckInHandler(c, xl, checkInReq); err != nil {
		xl.Error().Err(err).Msg("logic.CreateCheckInHandler error")
		v1.SendResponse(c, errmsg.ERROR, nil)
		return
	}

	// 响应
	v1.SendResponse(c, errmsg.SUCCSE, nil)
}
