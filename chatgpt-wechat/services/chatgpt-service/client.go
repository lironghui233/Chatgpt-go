package chatgpt_service

import (
	"chatgpt-wechat/pkg/config"
	"chatgpt-wechat/services"
	"chatgpt-wechat/services/client"
	"sync"
)

var pool services.ClientPool
var once sync.Once

type chatGPTDataClient struct {
	client.DefaultClient
}

func GetChatGPTServiceClientPool() services.ClientPool {
	once.Do(func() {
		cnf := config.GetConf()
		c := &chatGPTDataClient{}
		pool = c.GetPool(cnf.DependOnServices.ChatGPTService.Address)
	})
	return pool
}
