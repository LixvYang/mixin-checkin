package oauth

import (
	"log"
	"net/http"

	"github.com/fox-one/mixin-sdk-go"
	"github.com/gin-gonic/gin"
	"github.com/lixvyang/mixin-checkin/internal/dao/mongo"
	"github.com/lixvyang/mixin-checkin/internal/model"
	"github.com/lixvyang/mixin-checkin/internal/service/auth"
	"github.com/lixvyang/mixin-checkin/internal/utils/setting"
	"github.com/rs/zerolog"
)

// func MixinOauth(c *gin.Context) {
// 	var code = c.Query("code")
// 	access_token, _, err := mixin.AuthorizeToken(c, setting.Conf.ClientId, setting.Conf.AppSecret, code, "")
// 	if err != nil {
// 		log.Printf("AuthorizeToken: %v", err)
// 		return
// 	}

// 	userinfo, err := service.GetUserInfo(access_token)
// 	if err != nil {
// 		log.Println("Get userInfo fail!!!")
// 		if setting.Conf.Mode == "release" {
// 			c.Redirect(http.StatusFound, "https://betxin.one")
// 		} else if setting.Conf.Mode == "debug" {
// 			c.Redirect(http.StatusFound, "http://localhost:8080")
// 		}
// 	}

// 	user := model.User{
// 		AvatarUrl:      userinfo.AvatarURL,
// 		FullName:       userinfo.FullName,
// 		MixinId:        userinfo.IdentityNumber,
// 		IdentityNumber: userinfo.IdentityNumber,
// 		MixinUuid:      userinfo.UserID,
// 		SessionId:      userinfo.SessionID,
// 	}

// 	session := sessions.Default(c)

// 	// 如果用户不存在
// 	if checked := model.CheckUser(userinfo.UserID); checked != errmsg.SUCCSE {
// 		if coded := model.CreateUser(&user); coded != errmsg.SUCCSE {
// 			log.Println("Get userInfo fail!!!")
// 		}

// 		sessionToken := uuid.NewV4().String()
// 		session.Set("userId", user.MixinUuid)
// 		session.Set("token", sessionToken)
// 		_ = session.Save()
// 	} else {
// 		//用户存在 就更新数据
// 		if coded := model.UpdateUser(userinfo.UserID, &user); coded != errmsg.SUCCSE {
// 			log.Println("Update userInfo fail!!!")
// 		}
// 		session.Clear()
// 		sessionToken := uuid.NewV4().String()
// 		session.Set("userId", user.MixinUuid)
// 		session.Set("token", sessionToken)
// 		session.Save()
// 	}
// 	if setting.Conf.Mode == "release" {
// 		c.Redirect(http.StatusPermanentRedirect, "https://betxin.one")
// 	} else if setting.Conf.Mode == "debug" {
// 		c.Redirect(http.StatusPermanentRedirect, "http://localhost:8080")
// 	}
// }

func MixinOauth(c *gin.Context) {
	xl := c.MustGet("logger").(*zerolog.Logger)
	var code = c.Query("code")
	access_token, _, err := mixin.AuthorizeToken(c, setting.Conf.ClientId, setting.Conf.AppSecret, code, "")
	if err != nil {
		xl.Info().Msgf("AuthorizeToken: %v", err)
		return
	}

	userinfo, err := auth.GetUserInfo(access_token)
	if err != nil {
		log.Println("Get userInfo fail!!!")
		if setting.Conf.Mode == "release" {
			c.Redirect(http.StatusFound, "https://betxin.one")
		} else if setting.Conf.Mode == "debug" {
			c.Redirect(http.StatusFound, "http://localhost:8080")
		}
	}

	user := &model.User{
		FullName: userinfo.FullName,
		Uid:      userinfo.UserID,
	}

	err = mongo.CheckUserExist(c, xl, userinfo.UserID)
	if err != nil {
		// 如果用户不存在
		mongo.CreateUser(c, xl, user)
	} else {
		// 如果用户存在
	}

	// // 如果用户不存在
	// if checked := mongo.CheckUserExist(c, userinfo.UserID); checked != errmsg.SUCCSE {
	// 	if coded := model.CreateUser(&user); coded != errmsg.SUCCSE {
	// 		log.Println("Get userInfo fail!!!")
	// 	}

	// 	sessionToken := uuid.NewV4().String()
	// 	session.Set("userId", user.MixinUuid)
	// 	session.Set("token", sessionToken)
	// 	_ = session.Save()
	// } else {
	// 	//用户存在 就更新数据
	// 	if coded := model.UpdateUser(userinfo.UserID, &user); coded != errmsg.SUCCSE {
	// 		log.Println("Update userInfo fail!!!")
	// 	}
	// 	session.Clear()
	// 	sessionToken := uuid.NewV4().String()
	// 	session.Set("userId", user.MixinUuid)
	// 	session.Set("token", sessionToken)
	// 	session.Save()
	// }
	if setting.Conf.Mode == "release" {
		c.Redirect(http.StatusPermanentRedirect, "https://betxin.one")
	} else if setting.Conf.Mode == "debug" {
		c.Redirect(http.StatusPermanentRedirect, "http://localhost:8080")
	}
}
