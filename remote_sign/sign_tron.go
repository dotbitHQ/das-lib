package remote_sign

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/fbsobreira/gotron-sdk/pkg/proto/api"
	"github.com/fbsobreira/gotron-sdk/pkg/proto/core"
	"github.com/golang/protobuf/proto"
)

func (r *RemoteSignClient) SignTrxTx(address string, tx *api.TransactionExtention) (*api.TransactionExtention, error) {
	rawTx, err := TransactionToHexString(tx.Transaction)
	if err != nil {
		return nil, fmt.Errorf("TransactionToHexString err: %s", err.Error())
	}
	if err, signTx := r.signTrxTx(address, rawTx); err != nil {
		return nil, fmt.Errorf("signTrxTx err: %s", err.Error())
	} else {
		if coreTx, err := NewTransactionFromHexString(signTx); err != nil {
			return nil, fmt.Errorf("NewTransactionFromHexString err: %s", err.Error())
		} else {
			if raw, err := proto.Marshal(coreTx.GetRawData()); err != nil {
				return nil, fmt.Errorf("Marshal err: %s", err.Error())
			} else {
				txIdStr := fmt.Sprintf("%x", sha256.Sum256(raw))
				txId, _ := hex.DecodeString(txIdStr)
				retTx := api.TransactionExtention{
					Transaction:    coreTx,
					Txid:           txId,
					ConstantResult: tx.ConstantResult,
					Result:         tx.Result,
				}
				return &retTx, nil
			}
		}
	}
}

func (r *RemoteSignClient) signTrxTx(address, tronTxHexStr string) (error, string) {
	reply := reqParam{}
	type addressInfo struct {
		Address string `json:"address"`
		Index   int64  `json:"index"`
	}
	param := struct {
		Addrs []addressInfo `json:"addrs"`
		Tx    string        `json:"tx"`
	}{
		Addrs: []addressInfo{
			{
				Address: address,
			},
		},
		Tx: tronTxHexStr,
	}
	if err := r.client.CallContext(r.ctx, &reply, SignMethodTron, param); err != nil {
		return fmt.Errorf("client.CallContext err: %s", err.Error()), ""
	}
	if reply.Errno == 0 {
		return nil, reply.Data.(string)
	} else {
		return fmt.Errorf("client.CallContext err: %s", reply.Errmsg), ""
	}
}

func TransactionToHexString(tx *core.Transaction) (string, error) {
	data, err := proto.Marshal(tx)
	if err != nil {
		return "", fmt.Errorf("marshal tx:%v", err)
	}
	return hex.EncodeToString(data), nil
}

func NewTransactionFromHexString(raw string) (*core.Transaction, error) {
	data, err := hex.DecodeString(raw)
	if err != nil {
		return nil, fmt.Errorf("hex decode:%s %v", raw, err)
	}
	tx := core.Transaction{}
	if err := proto.Unmarshal(data, &tx); err != nil {
		return nil, fmt.Errorf("unmashal err:%v", err)
	}
	return &tx, nil
}
