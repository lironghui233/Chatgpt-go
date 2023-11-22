package chat_context

import (
	"chatgpt-service/pkg/config"
	"chatgpt-service/pkg/db/redis"
	"chatgpt-service/pkg/log"
	"chatgpt-service/proto"
	"context"
	"encoding/json"
	"time"

	redis2 "github.com/redis/go-redis/v9"
)

// qq、微信、企微
type hasID struct {
}

func (c *hasID) Get(id, group string, endpoint proto.ChatEndpoint) (value interface{}, err error) {
	key := getKey(id, group, endpoint)
	redisPool := redis.GetPool()
	redisClient := redisPool.Get()
	defer redisPool.Put(redisClient)

	list := make([]*ChatMessage, 0)
	value = list
	str, err := redisClient.Get(context.Background(), key).Result()
	if err == redis2.Nil {
		err = nil
		return
	}
	if err != nil {
		log.Error(err)
		return
	}
	err = json.Unmarshal([]byte(str), &list)
	if err != nil {
		log.Error(err)
		return
	}
	value = list
	return
}
func (c *hasID) Set(id, group string, endpoint proto.ChatEndpoint, value interface{}, ttl int) error {
	key := getKey(id, group, endpoint)
	redisPool := redis.GetPool()
	redisClient := redisPool.Get()
	defer redisPool.Put(redisClient)

	cnf := config.GetConf()
	list := value.([]*ChatMessage)
	if len(list) == 0 {
		return nil
	}
	if len(list) > cnf.Chat.ContextLen {
		list = list[:cnf.Chat.ContextLen]
	}
	bytes, err := json.Marshal(list)
	if err != nil {
		log.Error(err)
		return err
	}
	err = redisClient.SetEx(context.Background(), key, string(bytes), time.Duration(ttl)*time.Second).Err()
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}
