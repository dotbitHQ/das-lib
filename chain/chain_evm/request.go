package chain_evm

import (
	"fmt"
	"github.com/parnurzeal/gorequest"
	"net/http"
	"time"
)

func (c *ChainEvm) Request(url, method string, result interface{}) (resp *BaseResp, err error) {
	c.lock.Lock()
	defer c.lock.Unlock()

	baseResp := BaseResp{Result: result}
	ret, body, errs := gorequest.New().
		Timeout(time.Second*30).Post(url).
		Set("Content-Type", "application/json").
		Send(method).EndStruct(&baseResp)

	if errs != nil {
		log.Error("request error:", method, errs, "body:", string(body))
		return &baseResp, fmt.Errorf("request err: %v", errs)
	} else if ret != nil && ret.StatusCode != http.StatusOK {
		log.Error("request error:", method, errs, "body:", string(body))
		return &baseResp, fmt.Errorf("request err: %d", ret.StatusCode)
	}
	return &baseResp, nil
}

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type BaseResp struct {
	JsonRpc string      `json:"jsonrpc"`
	Id      int32       `json:"id"`
	Result  interface{} `json:"result"`
	Error   Error       `json:"error"`
}
