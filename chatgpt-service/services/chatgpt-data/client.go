package chatgpt_data

import (
	"chatgpt-service/pkg/config"
	"chatgpt-service/services"
	"chatgpt-service/services/client"
	"sync"
)

var pool services.ClientPool
var once sync.Once

type chatGPTDataClient struct {
	client.DefaultClient
}

func GetChatGPTDataClientPool() services.ClientPool {
	once.Do(func() {
		cnf := config.GetConf()
		c := &chatGPTDataClient{}
		pool = c.GetPool(cnf.DependOnServices.ChatGPTData.Address)
	})
	return pool
}
