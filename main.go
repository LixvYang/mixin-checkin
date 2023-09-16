package main

import (
	"github.com/lixvyang/mixin-checkin/internal/dao/mongo"
	"github.com/lixvyang/mixin-checkin/internal/utils/setting"
	"github.com/lixvyang/mixin-checkin/pkg/logger"

	"github.com/rs/zerolog/log"
)

func init() {
	if err := setting.Init("./configs/configs.yaml"); err != nil {
		log.Fatal().Err(err)
	}
	log.Info().Any("config", setting.Conf).Send()
	logger.Init(setting.Conf)
}

func main() {
	if err := mongo.Init(setting.Conf.MongoConfig); err != nil {
		logger.Lg.Panic().Err(err).Msg("err")
	}
}
