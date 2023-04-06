package remote_sign

import (
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
)

func (r *RemoteSignClient) SignEvmTx(method, address string, tx *types.Transaction) (*types.Transaction, error) {
	reply := reqParam{}
	txRlpBys, err := rlp.EncodeToBytes(tx)
	if err != nil {
		return nil, fmt.Errorf("rlp.EncodeToBytes err: %s", err.Error())
	}
	param := struct {
		Address string `json:"address"`
		Tx      string `json:"tx"`
	}{
		Address: address,
		Tx:      hex.EncodeToString(txRlpBys),
	}
	if err := r.client.CallContext(r.ctx, &reply, method, param); err != nil {
		return nil, fmt.Errorf("client.CallContext err: %s", err.Error())
	}
	if reply.Errno == 0 {
		signTxStr := reply.Data.(string)
		signTxBys, err := hex.DecodeString(signTxStr)
		if err != nil {
			return nil, fmt.Errorf("hex.DecodeString signed tx err: %s", err.Error())
		}
		signTx := types.Transaction{}
		if err = rlp.DecodeBytes(signTxBys, &signTx); err != nil {
			return nil, fmt.Errorf("rlp.DecodeBytes signed tx err: %s", err.Error())
		}
		return &signTx, nil
	} else {
		return nil, fmt.Errorf("client.CallContext err: %s", reply.Errmsg)
	}
}
