package send

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/redis/go-redis/v9"
)

type Sender struct{}

type DirectMessageBody struct {
	AccountId string
	Content   string
}

func (s *Sender) DirectMessage(accountId, content string) error {
	rc := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "",
		DB:       0,
	})
	ctx := context.Background()
	wsIp, err := rc.Get(ctx, accountId).Result()
	if err != nil {
		log.Printf("[gochatter-api] error getting result %v", rc)
		return err
	}
	dm := DirectMessageBody{accountId, content}
	body, err := json.Marshal(dm)
	if err != nil {
		log.Printf("[gochatter-api] error marshalling body %v", err)
		return err
	}
	log.Printf("[gochatter-api] found websocket server ip %s", wsIp)
	url := "https://" + wsIp + ":8444" + "/direct_message"
	r, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	c := &http.Client{}
	_, err = c.Do(r)
	if err != nil {
		return err
	}

	return nil
}
