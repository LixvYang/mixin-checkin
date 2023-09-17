package mixincli

import (
	"github.com/lixvyang/mixin-checkin/internal/utils/setting"
	"github.com/lixvyang/mixin-checkin/pkg/logger"

	"github.com/fox-one/mixin-sdk-go"
)

var (
	MixinCli = new(mixin.Client)
)

func Init(conf *setting.MixinConfig) (err error) {
	MixinCli, err = mixin.NewFromKeystore(&mixin.Keystore{
		ClientID:   conf.ClientId,
		SessionID:  conf.SessionId,
		PrivateKey: conf.PrivateKey,
		PinToken:   conf.PinToken,
	})
	if err != nil {
		logger.Lg.Fatal().Err(err).Msg("init mixin client error.")
		return err
	}
	go AutoReplay()
	return
}