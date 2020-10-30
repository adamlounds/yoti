package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type storeMessage struct {
	Id      []byte
	Payload []byte
}

func (c *httpClient) Store(id, payload []byte) (aesKey []byte, err error) {
	m := storeMessage{Id: id, Payload: payload}
	body, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/store", c.endpoint), bytes.NewBuffer(body))
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	_, respBody, err := c.do(ctx, req)
	if err != nil {
		return nil, err
	}
	return respBody, nil
}
