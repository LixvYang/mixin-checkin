package mongo

import (
	"context"
	"fmt"
	"sync"

	"github.com/lixvyang/mixin-checkin/internal/utils/setting"
	"github.com/lixvyang/mixin-checkin/pkg/logger"

	"github.com/qiniu/qmgo"
)

var (
	once        sync.Once
	mongoClient *qmgo.Client
	DB          *qmgo.Database
)

// Init 初始化Mongo连接
func Init(cfg *setting.MongoConfig) (err error) {
	once.Do(func() {
		ctx := context.Background()
		mongoClient, err := qmgo.NewClient(ctx, &qmgo.Config{Uri: fmt.Sprintf("mongodb://%s:%d", cfg.Host, cfg.Port)})
		if err != nil {
			logger.Lg.Panic().Err(err).Msg("Init Mongo error.")
		}
		DB = mongoClient.Database(cfg.DB)
		logger.Lg.Info().Msg("init mongo success.")
		// coll := db.Collection("user")
	})
	return
}

// Close 关闭Mongo连接
func Close() {
	mongoClient.Close(context.Background())
}
