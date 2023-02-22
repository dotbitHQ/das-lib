package bitcoin

import (
	"fmt"
	"github.com/parnurzeal/gorequest"
	"net/http"
	"time"
)

type resultBalance struct {
	Error       string `json:"error"`
	Success     int    `json:"success"`
	Balance     string `json:"balance"`
	Confirmed   string `json:"confirmed"`
	Unconfirmed string `json:"unconfirmed"`
}

func GetBalanceDoge(addr string) (result resultBalance, err error) {
	engine := gorequest.New().Timeout(time.Second * 30)
	url := fmt.Sprintf("https://dogechain.info/api/v1/address/balance/%s", addr)
	res, body, errs := engine.Get(url).EndStruct(&result)
	if len(errs) > 0 {
		err = fmt.Errorf("req errs: %v", errs)
		return
	}
	if res.StatusCode != http.StatusOK {
		err = fmt.Errorf("http code: %d, [%s]", res.StatusCode, body)
		return
	}
	if result.Success != 1 {
		err = fmt.Errorf("error: %s", result.Error)
		return
	}
	return
}
