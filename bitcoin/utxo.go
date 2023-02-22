package bitcoin

import (
	"fmt"
	"github.com/parnurzeal/gorequest"
	"net/http"
	"time"
)

func (t *TxTool) GetUnspentOutputsDoge(addr string, value int64) ([]UnspentOutputs, error) {
	var uos []UnspentOutputs
	total := int64(0)
	value += t.DustLimit

	for i := 1; total < value; i++ {
		result, err := t.getUnspentOutputsDoge(addr, i)
		if err != nil {
			return nil, fmt.Errorf("getUnspentOutputsDoge err: %s", err.Error())
		}
		if len(result.UnspentOutputs) == 0 {
			break
		}
		for _, v := range result.UnspentOutputs {
			tmp := UnspentOutputs{
				Private: t.PrivateKey,
				Address: addr,
				Hash:    v.TxHash,
				Index:   v.TxOutputN,
				Value:   v.Value,
			}
			uos = append(uos, tmp)
			total += v.Value
		}
	}
	if total < value {
		return nil, InsufficientBalanceError
	}
	return uos, nil
}

type resultUnspentOutputs struct {
	Error          string               `json:"error"`
	Success        int                  `json:"success"`
	UnspentOutputs []dataUnspentOutputs `json:"unspent_outputs"`
}

type dataUnspentOutputs struct {
	TxHash        string `json:"tx_hash"`
	TxOutputN     uint32 `json:"tx_output_n"`
	Script        string `json:"script"`
	Address       string `json:"address"`
	Value         int64  `json:"value"`
	Confirmations uint64 `json:"confirmations"`
}

func (t *TxTool) getUnspentOutputsDoge(addr string, page int) (result resultUnspentOutputs, err error) {
	engine := gorequest.New().Timeout(time.Second * 30)

	url := fmt.Sprintf("https://dogechain.info/api/v1/address/unspent/%s/%d", addr, page)
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
