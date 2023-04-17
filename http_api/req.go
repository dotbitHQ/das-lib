package http_api

import (
	"fmt"
	"github.com/parnurzeal/gorequest"
	"net/http"
	"time"
)

func SendReq(url string, req, data interface{}) error {
	var resp ApiResp
	resp.Data = &data

	res, _, errs := gorequest.New().Post(url).Retry(3, time.Second*5).
		Timeout(time.Second * 10).SendStruct(&req).EndStruct(&resp)

	if len(errs) > 0 {
		return fmt.Errorf("SendReq errs: %v", errs)
	} else if res.StatusCode != http.StatusOK {
		return fmt.Errorf("SendReq StatusCode: %d", res.StatusCode)
	}

	if resp.ErrNo != ApiCodeSuccess {
		return fmt.Errorf("%d - %s", resp.ErrNo, resp.ErrMsg)
	}
	return nil
}
