package keywords

import (
	"chatgpt-service/pkg/config"
	"chatgpt-service/services"
	"chatgpt-service/services/client"
	"sync"
)

var pool services.ClientPool
var once sync.Once

type keywordsClient struct {
	client.DefaultClient
}

func GetKeywordsClientPool() services.ClientPool {
	once.Do(func() {
		cnf := config.GetConf()
		c := &keywordsClient{}
		pool = c.GetPool(cnf.DependOnServices.Keywords.Address)
	})
	return pool
}
