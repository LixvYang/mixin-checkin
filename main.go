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
	"github.com/lixvyang/mixin-checkin/internal/service/mixincli"
	"github.com/lixvyang/mixin-checkin/internal/utils/cron"
	"github.com/lixvyang/mixin-checkin/internal/utils/setting"
	"github.com/lixvyang/mixin-checkin/pkg/logger"

	"github.com/rs/zerolog/log"
)

func init() {
	if err := setting.Init("./configs/configs.yaml"); err != nil {
		log.Fatal().Err(err)
	}
	logger.Get(setting.Conf)
	if err := mongo.Init(setting.Conf.MongoConfig); err != nil {
		logger.Lg.Panic().Err(err).Msg("init mongo error")
	}
	if err := redis.Init(setting.Conf.RedisConfig); err != nil {
		logger.Lg.Panic().Err(err).Msg("init redis err")
	}
	if err := mixincli.Init(setting.Conf.MixinConfig); err != nil {
		logger.Lg.Panic().Err(err).Msg("init mixincli err")
	}
	go cron.Sched.Init()
}

func main() {
	// if err := mixincli.SendMessage(context.Background(), "6a87e67f-02fb-47cf-b31f-32a13dd5b3d9"); err != nil {
	// 	logger.Lg.Error().Err(err).Send()
	// }
	// loc, _ := time.LoadLocation("Asia/Shanghai")
	// t := time.Now().UTC()
	// // 也可以直接使用24小时制格式
	// fmt.Println(t.In(loc).Format("15:04:05"))
	defer mongo.Close()
	defer redis.Close()
	r := router.Init(setting.Conf)
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
