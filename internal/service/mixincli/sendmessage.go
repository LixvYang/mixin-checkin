package mixincli

import (
	"context"
	"encoding/base64"
	"fmt"
	"math/rand"
	"time"

	"github.com/fox-one/mixin-sdk-go"
	"github.com/lixvyang/mixin-checkin/internal/utils/setting"
)

var randomQuotes = []string{
	"不畏将来，不念过往",
	"生活有时候就像马尔科夫过程，凡是过去的都成了序章，真正决定未来的是每一个当下",
	"待人友善是修养 独来独往是性格",
	"阅己，悦己，越己",
	"生命太短暂，不要去做根本没人想要的东西",
	"两列波在相遇后相互穿过，仍然保持各自的运动状态继续传播，彼此之间好像未曾相遇",
	"或许有一天，你会拍打着教科书说：编者还欠费功夫",
	"太空浩瀚，岁月悠长，我始终乐于和她分享同一颗行星和同一个时代",
	"质子在许许多多个夏天后死去",
	"改变或者离开，需要的可不仅仅是那一点儿聪明",
	"我们曾经仰望着浩瀚的星空 思考自身的存在",
	"人生不可逆，我亦是过程；独立且自由，孤独且执着",
	"自我教育是唯一的教育形式",
	"这世上只有一种成功 就是能够用自己喜欢的方式度过自己的一生",
	"作为大佬,请学会像资产阶级一样思考问题",
}

func generateQuotes() string {
	rand.Seed(time.Now().UnixNano())
	index := rand.Intn(len(randomQuotes))
	return randomQuotes[index]
}

func generateCheckinSuccessMsg() string {
	now := time.Now()
	year := now.Year()
	month := now.Month()
	day := now.Day()

	hour := now.Hour()
	minute := now.Minute()
	intsecond := now.Second()
	second := fmt.Sprintf("%d", intsecond)
	if int(intsecond) < 10 {
		second = fmt.Sprintf("0%d", intsecond)
	}

	res := fmt.Sprintf("今天是 %d 年 %d 月 %d 日\n 你在 %d:%d:%s 完成了今日份打卡\n \n%s", year, month, day, hour, minute, second, generateQuotes())
	return res
}

func SendMessage(ctx context.Context, uid string) (err error) {
	msg := generateCheckinSuccessMsg()
	return MixinCli.SendMessage(ctx, &mixin.MessageRequest{
		MessageID:      mixin.RandomTraceID(),
		ConversationID: mixin.UniqueConversationID(setting.Conf.ClientId, uid),
		RecipientID:    uid,
		Category:       mixin.MessageCategoryPlainText,
		Data:           msg,
		DataBase64:     base64.RawURLEncoding.EncodeToString([]byte(msg)),
	})
}
