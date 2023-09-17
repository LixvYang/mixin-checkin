package mixincli

import (
	"context"

	"github.com/fox-one/mixin-sdk-go"
	"github.com/lixvyang/mixin-checkin/internal/utils/setting"
	"github.com/rs/zerolog"
	"github.com/shopspring/decimal"
)

func Transfer(ctx context.Context, xl *zerolog.Logger, uid string) (err error) {
	amount, err := decimal.NewFromString("1")
	if err != nil {
		xl.Err(err).Caller().Send()
	}
	_, err = MixinCli.Transfer(ctx, &mixin.TransferInput{
		AssetID:    "4d8c508b-91c5-375b-92b0-ee702ed2dac5",
		OpponentID: uid,
		Memo:       "测试一下",
		TraceID:    mixin.UniqueConversationID(setting.Conf.MixinConfig.ClientId, uid),
		Amount:     amount,
	}, setting.Conf.MixinConfig.Pin)
	if err != nil {
		xl.Err(err).Caller().Send()
	}
	return err
}
