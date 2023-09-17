package mongo

import (
	"context"

	"github.com/lixvyang/mixin-checkin/internal/model"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
)

var (
	ErrorUserExist       = errors.New("用户已存在")
	ErrorUserNoExist     = errors.New("用户不存在")
	ErrorInvaildPassword = errors.New("用户名或密码不正确")
)

// 查找用户是否存在
func CheckUserExist(ctx context.Context, xl *zerolog.Logger, uid string) (err error) {
	count, err := coll.UserColl.Find(ctx, bson.M{"uid": uid}).Count()
	if err != nil {
		xl.Error().Err(err).Msg("coll.UserColl.Find err")
		return err
	}
	if count > 0 {
		return ErrorUserExist
	}
	return
}

func CreateUser(ctx context.Context, xl *zerolog.Logger, userInfo *model.User) (err error) {
	_, err = coll.UserColl.InsertOne(ctx, userInfo)
	if err != nil {
		xl.Error().Err(err).Msg("coll.UserColl.InsertOne err")
		return err
	}
	return
}

// 查找所有用户
func FindAllUser(ctx context.Context, xl *zerolog.Logger) (checkins []model.Checkin, err error) {
	err = coll.CheckIn.Find(ctx, bson.M{}).All(&checkins)
	if err != nil {
		xl.Error().Err(err).Msg("coll.CheckIn.Find err")
		return checkins, err
	}
	return checkins, nil
}
