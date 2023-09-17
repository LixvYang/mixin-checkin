package logic

import (
	"github.com/gin-gonic/gin"
	"github.com/lixvyang/mixin-checkin/internal/dao/mongo"
	"github.com/lixvyang/mixin-checkin/internal/model"
	"github.com/rs/zerolog"
)

func CreateCheckInHandler(c *gin.Context, xl *zerolog.Logger, req *model.CheckinReq) (err error) {
	return mongo.InsertCheckIn(c, xl, req)
}
