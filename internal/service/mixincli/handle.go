package mixincli

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/fox-one/mixin-sdk-go"
	"github.com/gofrs/uuid"
	"github.com/lixvyang/mixin-checkin/internal/dao/mongo"
	"github.com/lixvyang/mixin-checkin/pkg/logger"
	"github.com/rs/zerolog"
)

func AutoReplay() {
	h := func(ctx context.Context, msg *mixin.MessageView, userID string) error {
		if userID, _ := uuid.FromString(msg.UserID); userID == uuid.Nil {
			return nil
		}

		return handleMsg(ctx, msg)
	}

	ctx := context.Background()

	for {
		if err := MixinCli.LoopBlaze(ctx, mixin.BlazeListenFunc(h)); err != nil {
			logger.Lg.Info().Caller().Msgf("LoopBlaze: %v", err)
		}

		time.Sleep(time.Second)
	}
}

func handleMsg(ctx context.Context, msg *mixin.MessageView) (err error) {
	// 处理文本消息
	// 1. 解码消息内容
	// 查看用户是否存在 如果不存在 则返回无此用户的消息
	// 2. 根据uid去查找用户信息，然后发起请求
	// 发起请求后，需要处理请求 如果发现发送的是checkin，则进行下一步
	// 下一步是先要求深呼吸1分钟后，再进行鼓励和和记录打卡信息
	// Decode the message content
	xl := logger.Lg
	msgContent, err := base64.StdEncoding.DecodeString(msg.Data)
	if err != nil {
		return err
	}

	err = mongo.CheckUserExist(ctx, &xl, msg.UserID)
	if err != nil {
		// no exist user
		xl.Error().Msgf("UsedID <%s> no exists", msg.UserID)
		return err
	}

	session := getSession(msg.UserID)
	if session == nil {
		setSession(msg.UserID, &UserSession{
			State: UserSessionStateInit,
		})
	}

	session = getSession(msg.UserID)

	if session.State == UserSessionStateWaiting {
		data := []byte("您在等待中...")
		err = respondTextMsg(ctx, msg, mixin.MessageCategoryPlainText, data, rand.Intn(100))
		return
	}

	if session.State == UserSessionStateInit {
		if strings.ToUpper(string(msgContent)) == "CHECKIN" {
			if err = mongo.CheckCheckInRecord(ctx, &xl, msg.UserID); err != nil {
				data := []byte("今天已经签过到了😯~")
				err = respondTextMsg(ctx, msg, mixin.MessageCategoryPlainText, data, 1)
				return
			}
			// 签到
			// 1. 发送请跟着以下图进行深呼吸文字
			data := []byte("请跟着下方图片呼吸30秒")
			err = respondTextMsg(ctx, msg, mixin.MessageCategoryPlainText, data, 1)
			if err != nil {
				xl.Error().Err(err).Msg("发送文本出错")
			}
			time.Sleep(1 * time.Second)
			// 2. 发送图文
			err = respondImageMsg(ctx, &xl, msg, mixin.MessageCategoryPlainImage, 2)
			if err != nil {
				xl.Error().Err(err).Msg("发送图片出错")
			}

			setSession(msg.UserID, &UserSession{
				State: UserSessionStateWaiting,
			})
			// 3. 等待30秒
			time.Sleep(30 * time.Second)
			// 4. 发送 checkin 标签
			err = respondButtons(ctx, &xl, msg, mixin.MessageCategoryAppButtonGroup, 3)
			if err != nil {
				xl.Error().Err(err).Msg("发送按钮组出错")
			}
			// 5. 表示发送成功 转账+激励
			setSession(msg.UserID, &UserSession{
				State: UserSessionStateSuccess,
			})
		}
	} else {
		// 这里要获取用户session状态
		if string(msgContent) == "RECHECK" {
			err = SendSuccessMessage(ctx, &xl, msg.UserID)
			if err != nil {
				xl.Error().Err(err).Msg("发送签到成功消息出错")
			}
			time.Sleep(1 * time.Second)
			err = Transfer(ctx, &xl, msg.UserID)
			if err != nil {
				xl.Error().Err(err).Msg("转账出错")
			}
			// 清空状态
			setSession(msg.UserID, nil)
			return
		}
	}

	return nil
}

func respondTextMsg(ctx context.Context, msg *mixin.MessageView, category string, data []byte, step int) error {
	payload := base64.StdEncoding.EncodeToString(data)
	id, _ := uuid.FromString(msg.MessageID)
	// Create a request
	reply := &mixin.MessageRequest{
		ConversationID: msg.ConversationID,
		RecipientID:    msg.UserID,
		MessageID:      uuid.NewV5(id, fmt.Sprintf("reply %d", step)).String(),
		Category:       category,
		Data:           payload,
	}
	// Send the response
	return MixinCli.SendMessage(ctx, reply)
}

func respondImageMsg(ctx context.Context, xl *zerolog.Logger, msg *mixin.MessageView, category string, step int) error {
	id, _ := uuid.FromString(msg.MessageID)
	img := mixin.ImageMessage{
		AttachmentID: "a96dad4b-6b1e-49ec-8a6e-fef3055c694f",
		MimeType:     "image/gif",
		Size:         212654,
		Width:        500,
		Height:       500,
		Thumbnail:    "L3T9L#xufQxu~qfQfQa|?bazfQay",
	}
	data, err := json.Marshal(img)
	if err != nil {
		xl.Error().Err(err).Send()
		return err
	}
	payload := base64.StdEncoding.EncodeToString(data)

	// Create a request
	reply := &mixin.MessageRequest{
		ConversationID: msg.ConversationID,
		RecipientID:    msg.UserID,
		MessageID:      uuid.NewV5(id, fmt.Sprintf("reply %d", step)).String(),
		Category:       mixin.MessageCategoryPlainImage,
		Data:           payload,
	}
	// Send the response
	return MixinCli.SendMessage(ctx, reply)
}

func respondButtons(ctx context.Context, xl *zerolog.Logger, msg *mixin.MessageView, category string, step int) error {
	id, _ := uuid.FromString(msg.MessageID)
	button := []mixin.AppButtonMessage{
		{
			Label:  "Check In",
			Color:  "#7983C2",
			Action: "input:RECHECK",
		},
	}
	data, err := json.Marshal(button)
	if err != nil {
		xl.Error().Err(err).Send()
		return err
	}
	payload := base64.StdEncoding.EncodeToString(data)

	// Create a request
	reply := &mixin.MessageRequest{
		ConversationID: msg.ConversationID,
		RecipientID:    msg.UserID,
		MessageID:      uuid.NewV5(id, fmt.Sprintf("reply %d", step)).String(),
		Category:       mixin.MessageCategoryAppButtonGroup,
		Data:           payload,
	}
	// Send the response
	return MixinCli.SendMessage(ctx, reply)
}
