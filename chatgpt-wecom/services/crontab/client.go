package chatgpt_service

import (
	"chatgpt-wecom/pkg/config"
	"chatgpt-wecom/services"
	"chatgpt-wecom/services/client"
	"sync"
)

var pool services.ClientPool
var once sync.Once

type chatGPTDataClient struct {
	client.DefaultClient
}

func GetCrontabClientPool() services.ClientPool {
	once.Do(func() {
		cnf := config.GetConf()
		c := &chatGPTDataClient{}
		pool = c.GetPool(cnf.DependOnServices.ChatGPTService.Address)
	})
	return pool
}
