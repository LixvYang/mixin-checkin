package mongo

import (
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lixvyang/mixin-checkin/internal/model"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
)

var (
	ErrorCheckInLate    = errors.New("签到太迟了")
	ErrorCheckInNoExist = errors.New("用户签到表不存在")
)

// 检查签到是否已经迟到
func CheckCheckInLate(ctx *gin.Context, xl *zerolog.Logger, req *model.CheckinReq) (err error) {
	checkin := new(model.Checkin)
	err = coll.CheckIn.Find(ctx, bson.M{"uid": req.Uid}).One(checkin)
	if err != nil {
		xl.Error().Err(err).Msg("coll.CheckinRecord.CountDocuments err")
		return err
	}

	reqTime, err := time.Parse(TimeFormat_HHmmss, strings.Split(req.Time, " ")[1])
	if err != nil {
		xl.Error().Err(err).Msg("time.Parse(TimeFormat, req.Time) err")
		return err
	}
	checkinTime, err := time.Parse(TimeFormat_HHmmss, checkin.Time)
	if err != nil {
		xl.Error().Err(err).Msg("time.Parse(TimeFormat, req.Time) err")
		return err
	}
	if reqTime.Before(checkinTime) {
		// 符合签到条件
		return nil
	}
	xl.Error().Err(ErrorCheckInLate).Send()
	return ErrorCheckInLate
}

// 用户签到表是否存在
func CheckCheckInExist(ctx *gin.Context, xl *zerolog.Logger, req *model.CheckinReq) (err error) {
	count, err := coll.CheckIn.Find(ctx, bson.M{
		"uid": req.Uid,
	}).Count()

	if err != nil {
		xl.Error().Err(err).Msg("coll.CheckIn.Find err")
		return err
	}

	if count == 0 {
		return ErrorCheckInNoExist
	}

	return
}

// 插入用户签到表
func InsertCheckIn(ctx *gin.Context, xl *zerolog.Logger, req *model.CheckinReq) (err error) {
	_, err = coll.CheckIn.InsertOne(ctx, *req)
	if err != nil {
		xl.Error().Err(err).Msg("coll.CheckIn.InsertOne err")
		return err
	}

	return
}
