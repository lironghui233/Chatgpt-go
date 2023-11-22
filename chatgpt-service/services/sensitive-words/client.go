package sensitive_words

import (
	"chatgpt-service/pkg/config"
	"chatgpt-service/services"
	"chatgpt-service/services/client"
	"sync"
)

var pool services.ClientPool
var once sync.Once

type sensitiveWordsClient struct {
	client.DefaultClient
}

func GetSensitiveWordsClientPool() services.ClientPool {
	once.Do(func() {
		cnf := config.GetConf()
		c := &sensitiveWordsClient{}
		pool = c.GetPool(cnf.DependOnServices.SensitiveWords.Address)
	})
	return pool
}
