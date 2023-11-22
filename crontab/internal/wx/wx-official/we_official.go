package wx_official

import (
	"crontab/internal/wx"
	"fmt"
)

type wxOfficial struct {
	*wx.DefaultToken
}

func NewWxOfficial(id, secret string) wx.Token {
	url := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s", id, secret)
	return &wxOfficial{
		DefaultToken: &wx.DefaultToken{
			Id:     id,
			Secret: secret,
			Url:    url,
			App:    "",
		},
	}
}
