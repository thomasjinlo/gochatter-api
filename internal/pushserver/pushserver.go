package pushserver

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type PushServer struct {
	ip string
}

type MessageBody struct {
	Author string
	Content string
}

func (p *PushServer) Broadcast(author, content string) (*http.Response, error) {
	msg := MessageBody {
		Author: author,
		Content: content,
	}
	body, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}
	url := "http://" + p.ip + ":8444" + "/send_message"
	r, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	c := &http.Client{}
	res, err := c.Do(r)
	if err != nil {
		return nil, err
	}
	return res, err
}

func NewPushServer(ip string) *PushServer {
	return &PushServer{ip: ip}
}
