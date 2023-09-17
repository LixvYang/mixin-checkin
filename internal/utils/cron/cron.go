package cron

import (
	"context"
	"sync"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/lixvyang/mixin-checkin/internal/dao/mongo"
	"github.com/lixvyang/mixin-checkin/pkg/logger"
)

var Sched = new(CheckInSched)

func init() {
	timezone, _ := time.LoadLocation("Asia/Shanghai")
	Sched.sched = gocron.NewScheduler(timezone)
}

type CheckInSched struct {
	sched *gocron.Scheduler
}

func (c *CheckInSched) Init() {
	c.EveryDayJob()
}

// 每天早上7点检查是否早起了的逻辑
func (c *CheckInSched) EveryDayJob() {
	job := func() {
		logger.Lg.Info().Msg("开始执行每日任务")
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		checkins, err := mongo.FindAllUser(ctx, &logger.Lg)
		if err != nil {
			logger.Lg.Err(err).Str("查找用户错误", "err").Send()
			return
		}
		var wg sync.WaitGroup
		checkinsLen := len(checkins)
		wg.Add(checkinsLen)
		for i := 0; i < checkinsLen; i++ {
			go func(i int) {
				defer wg.Done()
				err := mongo.FindCheckInRecordToday(ctx, &logger.Lg, checkins[i].Uid)
				if err != nil {
					// 今天没签到
					// mixincli.MixinCli.SendMessage(ctx, &mixin.MessageRequest{})
					logger.Lg.Err(err).Msg("今天没签到")
					return
				}
				// 今天签到了
				logger.Lg.Err(err).Msg("今天签到了")
			}(i)
		}
		wg.Wait()
	}

	c.sched.Every(1).Day().At("10:44:00").Do(job)
	// c.sched.Every(1).Second().Do(f)
	c.sched.StartAsync()
}
