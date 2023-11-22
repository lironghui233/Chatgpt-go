package chatgpt_service

import (
	"chatgpt-qq/config"
	"chatgpt-qq/services"
	"chatgpt-qq/services/client"
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
		pool = c.GetPool(cnf.ChatGPTService.Address)
	})
	return pool
}
