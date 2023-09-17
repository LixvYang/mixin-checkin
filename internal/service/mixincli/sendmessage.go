package mixincli

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/fox-one/mixin-sdk-go"
	v1 "github.com/lixvyang/mixin-checkin/api/v1"
	"github.com/lixvyang/mixin-checkin/internal/dao/mongo"
	"github.com/lixvyang/mixin-checkin/internal/model"
	"github.com/lixvyang/mixin-checkin/internal/utils/errmsg"
	"github.com/lixvyang/mixin-checkin/internal/utils/setting"
	"github.com/lixvyang/mixin-checkin/pkg/logger"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

var randomQuotes = []string{
	"不畏将来 不念过往",
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

	res := fmt.Sprintf("今天是 %d 年 %d 月 %d 日\n 你在 %d:%d:%s 完成了今日份打卡\n \n%s\n \n这是你早起的奖励⬇️⬇️⬇️", year, month, day, hour, minute, second, generateQuotes())
	return res
}

func SendSuccessMessage(ctx context.Context, xl *zerolog.Logger, uid string) (err error) {
	// 向本地发送 /checkin http请求
	err = sendCheckIn(&model.CheckinReq{
		Uid:  uid,
		Time: time.Now().Format(mongo.TimeFormat_YYMMDD_HHmmss),
	})

	if err != nil {
		logger.Lg.Error().Err(err).Send()
		return err
	}

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

func sendCheckIn(checkReq *model.CheckinReq) (err error) {
	url := fmt.Sprintf("http://127.0.0.1:%d/api/v1/checkin", setting.Conf.Port)
	var data []byte
	data, err = json.Marshal(checkReq)
	if err != nil {
		return err
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", url, strings.NewReader(string(data)))
	if err != nil {
		logger.Lg.Error().Err(err).Send()
		return err
	}

	req.Header.Add("Content-Type", "application/json")
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	var resp v1.Response
	if err = json.Unmarshal(body, &resp); err != nil {
		return err
	}
	if resp.Code == errmsg.ERROR_RECHECKIN {
		return errors.New("重复签到")
	}

	return
}
