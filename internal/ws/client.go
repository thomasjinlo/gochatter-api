package ws

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/redis/go-redis/v9"
)

type Client struct {
	rc *redis.Client
	hc *http.Client
}

func NewClient(rc *redis.Client, hc *http.Client) *Client {
	return &Client{rc: rc, hc: hc}
}

type directMessage struct {
	SourceAccountId string
	TargetAccountId string
	Content         string
}

func (c *Client) SendDirectMessage(targetAccountId, sourceAccountId, content string) error {
	ctx := context.Background()
	wsIps, err := c.rc.SMembers(ctx, targetAccountId).Result()
	if err != nil {
		log.Printf("[gochatter-api] error retrieving ws ip: %v", err)
		return err
	}
	dm := directMessage{
		TargetAccountId: targetAccountId,
		SourceAccountId: sourceAccountId,
		Content:         content,
	}
	body, err := json.Marshal(dm)
	if err != nil {
		return err
	}
	for _, wsIp := range wsIps {
		url := "https://" + wsIp + ":8444" + "/direct_message"
		r, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
		r.Header.Add("Content-Type", "application/json")
		res, err := c.hc.Do(r)
		if err != nil {
			return err
		}
		if res.StatusCode != http.StatusOK {
			return errors.New(fmt.Sprintf("dm request failed with status: %v", res.StatusCode))
		}
	}
	return nil
}
