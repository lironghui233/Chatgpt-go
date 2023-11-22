package cron

import (
	we_com "crontab/internal/wx/we-com"
	wx_official "crontab/internal/wx/wx-official"
	"crontab/pkg/config"
	"crontab/pkg/log"

	"github.com/robfig/cron/v3"
)

func Run() {
	// cron.WithSeconds() 启用每秒定时任务
	// cron.WithLocation(time.Local) 设置时区为本地时区
	cnf := config.GetConf()
	c := cron.New()
	c.AddFunc("*/5 * * * *", func() {
		for _, item := range cnf.WeComs {
			weCom := we_com.NewWecom(item.CorpId, item.CorpSecret, item.App)
			err := weCom.RefreshToken()
			if err != nil {
				log.Error(err)
				continue
			}
		}
		for _, item := range cnf.WxOfficials {
			official := wx_official.NewWxOfficial(item.AppId, item.Secret)
			err := official.RefreshToken()
			if err != nil {
				log.Error(err)
				continue
			}

		}
	})
	c.Run()
}
