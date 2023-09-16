package router

import (
	"github.com/gin-gonic/gin"
	"github.com/lixvyang/ginzero"
	"github.com/lixvyang/mixin-checkin/internal/middleware"
	"github.com/lixvyang/mixin-checkin/internal/utils/setting"
	"github.com/lixvyang/mixin-checkin/pkg/logger"
)

func Init(conf *setting.AppConfig) *gin.Engine {
	r := gin.New()
	gin.SetMode(conf.Mode)
	r.Use(ginzero.Ginzero(&logger.Lg), ginzero.RecoveryWithZero(&logger.Lg, true), middleware.GinXid(&logger.Lg))

	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "hello")
	})

	r.GET("/panic", func(c *gin.Context) {
		panic("/panic")
	})

	return r
}
