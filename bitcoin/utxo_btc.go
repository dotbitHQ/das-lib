package bitcoin

import (
	"fmt"
	"github.com/parnurzeal/gorequest"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

type btcUTXO struct {
	Txid  string `json:"txid"`
	Vout  uint32 `json:"vout"`
	Value string `json:"value"`
}

func GetUnspentOutputsBtc(addr, privateKey, url, apiKey string, value int64) (int64, []UnspentOutputs, error) {
	var uos []UnspentOutputs
	total := int64(0)
	value += DustLimitBtc

	engine := gorequest.New().Timeout(time.Second * 30)
	// https://btc.nownodes.io/api/v2/utxo/%s?padge=1&limit=20
	rand.Seed(time.Now().UnixNano())
	randNum := rand.Intn(1000)
	url = fmt.Sprintf("%s/%s?p=%d", url, addr, randNum)
	var data []btcUTXO
	engine = engine.Get(url).AppendHeader("api-key", apiKey)

	res, body, errs := engine.EndStruct(&data)
	if len(errs) > 0 {
		return 0, nil, fmt.Errorf("req errs: %v[%s][%d]", errs, body, randNum)
	}
	if res.StatusCode != http.StatusOK {
		return 0, nil, fmt.Errorf("http code: %d, [%s]", res.StatusCode, body)
	}
	for _, v := range data {
		outValue, err := strconv.ParseInt(v.Value, 10, 64)
		if err != nil {
			return 0, nil, fmt.Errorf("strconv.ParseInt err: %s", err.Error())
		}
		tmp := UnspentOutputs{
			Private: privateKey,
			Address: addr,
			Hash:    v.Txid,
			Index:   v.Vout,
			Value:   outValue,
		}
		uos = append(uos, tmp)
		total += outValue
		if total > value {
			break
		}
	}

	return total, uos, nil
}
