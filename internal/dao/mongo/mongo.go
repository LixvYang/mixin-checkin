package mongo

import (
	"context"
	"fmt"
	"sync"

	"github.com/lixvyang/mixin-checkin/internal/utils/setting"
	"github.com/lixvyang/mixin-checkin/pkg/logger"

	"github.com/qiniu/qmgo"
)

var mgo = new(Mgo)
var coll = new(Coll)

type Mgo struct {
	once     sync.Once
	mongoCli *qmgo.Client
	DB       *qmgo.Database
}

type Coll struct {
	UserColl      *qmgo.Collection
	CheckIn       *qmgo.Collection
	CheckInRecord *qmgo.Collection
}

// Init 初始化Mongo连接
func Init(cfg *setting.MongoConfig) (err error) {
	mgo.once.Do(func() {
		ctx := context.Background()
		mgo.mongoCli, err = qmgo.NewClient(ctx, &qmgo.Config{Uri: fmt.Sprintf("mongodb://%s:%d", cfg.Host, cfg.Port)})
		if err != nil {
			logger.Lg.Panic().Err(err).Msg("Init Mongo error.")
		}
		mgo.DB = mgo.mongoCli.Database(cfg.DB)
		logger.Lg.Info().Msg("init mongo success.")
		// coll := DB.Collection("user")
	})

	coll.UserColl = mgo.DB.Collection("user")
	coll.CheckIn = mgo.DB.Collection("checkin")
	coll.CheckInRecord = mgo.DB.Collection("checkin_record")
	return
}

// Close 关闭Mongo连接
func Close() {
	logger.Lg.Info().Msg("mongo success close.")
	_ = mgo.mongoCli.Close(context.Background())
}
