package we_com

import (
	"crontab/internal/wx"
	"fmt"
)

type weCom struct {
	*wx.DefaultToken
}

func NewWecom(id, secret, app string) wx.Token {
	url := fmt.Sprintf("https://qyapi.weixin.qq.com/cgi-bin/gettoken?corpid=%s&corpsecret=%s", id, secret)
	return &weCom{
		DefaultToken: &wx.DefaultToken{
			Id:     id,
			Secret: secret,
			Url:    url,
			App:    app,
		},
	}
}
