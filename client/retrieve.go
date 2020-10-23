package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type retrieveMessage struct {
	Id []byte
	AesKey []byte
}
func (c *httpClient) Retrieve(id, aeskey []byte) (payload []byte, err error) {
	m := retrieveMessage{Id: id, AesKey: aeskey}
	body, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/retrieve", c.endpoint), bytes.NewBuffer(body))
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	_, respBody, err := c.do(ctx, req)
	return respBody, nil
}

