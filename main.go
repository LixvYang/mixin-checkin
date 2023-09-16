package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/lixvyang/mixin-checkin/internal/dao/mongo"
	"github.com/lixvyang/mixin-checkin/internal/dao/redis"
	"github.com/lixvyang/mixin-checkin/internal/router"
	"github.com/lixvyang/mixin-checkin/internal/utils/setting"
	"github.com/lixvyang/mixin-checkin/pkg/logger"

	"github.com/rs/zerolog/log"
)

func init() {
	if err := setting.Init("./configs/configs_example.yaml"); err != nil {
		log.Fatal().Err(err)
	}
	log.Info().Any("config", setting.Conf).Send()
	logger.Get(setting.Conf)
	if err := mongo.Init(setting.Conf.MongoConfig); err != nil {
		logger.Lg.Panic().Err(err).Msg("init mongo error")
	}
	if err := redis.Init(setting.Conf.RedisConfig); err != nil {
		logger.Lg.Panic().Err(err).Msg("init redis err")
	}
}

func main() {
	defer mongo.Close()
	defer redis.Close()
	// 5. 注册路由
	r := router.Init()
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", setting.Conf.Port),
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Lg.Fatal().Msgf("listen: %v\n", err)
		}
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	if err := srv.Shutdown(ctx); err != nil {
		logger.Lg.Fatal().Err(err).Msg("Server ShutDown.")
	}
	logger.Lg.Info().Msg("Server exiting.")
}
