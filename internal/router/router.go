package router

import (
	"github.com/gin-gonic/gin"
	"github.com/lixvyang/ginzero"
	"github.com/lixvyang/mixin-checkin/api/v1/checkin"
	"github.com/lixvyang/mixin-checkin/internal/middleware"
	"github.com/lixvyang/mixin-checkin/internal/utils/setting"
	"github.com/lixvyang/mixin-checkin/pkg/logger"
)

func Init(conf *setting.AppConfig) *gin.Engine {
	r := gin.New()
	gin.SetMode(conf.Mode)
	r.Use(ginzero.Ginzero(&logger.Lg), ginzero.RecoveryWithZero(&logger.Lg, true), middleware.GinXid(&logger.Lg))

	a := r.Group("/api/v1")
	{
		a.POST("/checkin", checkin.PostCheckIn)
		a.POST("/createcheckin", checkin.CreateCheckIn)
	}

	return r
}
