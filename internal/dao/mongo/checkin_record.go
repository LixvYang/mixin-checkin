package mongo

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lixvyang/mixin-checkin/internal/model"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
)

const TimeFormat_YYMMDD_HHmmss = "2006-01-02 15:04:05"
const TimeFormat_HHmmss = "15:04:05"

var (
	CheckInRecordExist   = errors.New("签到记录已经存在")
	CheckInRecordNoExist = errors.New("签到记录不存在")
)

// 签到当天记录是否已经存在
func CheckCheckInRecord(ctx context.Context, xl *zerolog.Logger, uid string) (err error) {
	// todayTime := strings.Split(time.Now().Format(TimeFormat_YYMMDD_HHmmss), " ")[0]
	// todayTime := strings.Split(req.Time, " ")[0]
	todayTime := strings.Split(time.Now().Format(TimeFormat_YYMMDD_HHmmss), " ")[0]
	count, err := coll.CheckInRecord.Find(ctx, bson.M{
		"time": bson.M{
			"$regex": fmt.Sprintf("^%s", todayTime),
		},
		"uid": uid,
	}).Count()

	if err != nil {
		xl.Error().Err(err).Msg("coll.CheckinRecord.CountDocuments err")
		return err
	}

	if count > 0 {
		return CheckInRecordExist
	}

	return
}

func InsertCheckInRecord(ctx *gin.Context, xl *zerolog.Logger, req *model.CheckinReq) (err error) {
	_, err = coll.CheckInRecord.InsertOne(ctx, req)
	if err != nil {
		xl.Error().Err(err).Msg("coll.CheckInRecord.InsertOne err")
		return err
	}

	return
}

// 查找当天的签到记录
func FindCheckInRecordToday(ctx context.Context, xl *zerolog.Logger, uid string) (err error) {
	todayNow := string(strings.Split(time.Now().Format("2006-01-02 03:04:05"), " ")[0])
	count, err := coll.CheckInRecord.Find(ctx, bson.M{
		"time": bson.M{
			"$regex": fmt.Sprintf("^%s", todayNow),
		},
		"uid": uid,
	}).Count()

	if err != nil {
		xl.Error().Err(err).Msg("coll.CheckInRecord.Find")
		return err
	}

	if count > 0 {
		return nil
	}

	return CheckInRecordNoExist
}
