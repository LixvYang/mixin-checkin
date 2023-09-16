package mixincli

import (
	"github.com/lixvyang/mixin-checkin/internal/utils/setting"
	"github.com/lixvyang/mixin-checkin/pkg/logger"

	"github.com/fox-one/mixin-sdk-go"
)

var (
	MixinCli *mixin.Client
	err      error
)

func InitMixinClient(conf setting.MixinConfig) {
	MixinCli, err = mixin.NewFromKeystore(&mixin.Keystore{
		ClientID:   conf.ClientId,
		SessionID:  conf.SessionId,
		PrivateKey: conf.PrivateKey,
		PinToken:   conf.PinToken,
	})
	if err != nil {
		logger.Lg.Fatal()
	}
}
