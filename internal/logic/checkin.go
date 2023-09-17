package logic

import (
	"github.com/gin-gonic/gin"
	"github.com/lixvyang/mixin-checkin/internal/dao/mongo"
	"github.com/lixvyang/mixin-checkin/internal/model"
	"github.com/rs/zerolog"
)

func CheckInHandler(c *gin.Context, xl *zerolog.Logger, req *model.CheckinReq) (err error) {
	// 1. 检查用户是否存在
	// 1.1 不存在则退出
	if err = mongo.CheckUserExist(c, xl, req.Uid); err != nil {
		xl.Info().Err(err).Msg("用户不存在~")
		return err
	}

	// 检查用户的签到表是否存在
	if err = mongo.CheckCheckInExist(c, xl, req); err != nil {
		xl.Info().Err(err).Msg("用户签到表不存在~")
		return err
	}

	// 2.0 检查是否当天已经签过到了 如果签过了 则退出 返回已经签过到了
	if err = mongo.CheckCheckInRecord(c, xl, req); err != nil {
		xl.Info().Err(err).Msg("今天已经签过到了~")
		return err
	}

	// 2.1 存在则继续 检查时间 如果超过则退出
	if err = mongo.CheckCheckInLate(c, xl, req); err != nil {
		xl.Info().Err(err).Msg("今天签到迟到了~")
		return err
	}

	// 3. 插入文档
	if err = mongo.InsertCheckInRecord(c, xl, req); err != nil {
		xl.Info().Err(err).Msg("插入文档失败!")
		return err
	}
	return
}
