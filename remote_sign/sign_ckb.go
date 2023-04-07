package remote_sign

import (
	"encoding/hex"
	"fmt"
	"github.com/dotbitHQ/das-lib/common"
)

func (r *RemoteSignClient) SignCkbMessage(ckbSignerAddress, message string) ([]byte, error) {
	if common.Has0xPrefix(message) {
		message = message[2:]
	}
	reply := reqParam{}
	param := struct {
		Address     string `json:"address"`
		CkbBuildRet string `json:"ckb_build_ret"`
		Tx          string `json:"tx"`
	}{
		Address:     ckbSignerAddress,
		CkbBuildRet: "",
		Tx:          message,
	}
	if err := r.client.CallContext(r.ctx, &reply, SignMethodCkb, param); err != nil {
		return nil, fmt.Errorf("remoteRpcClient.Call err: %s", err.Error())
	}
	if reply.Errno == 0 {
		signTxStr := reply.Data.(string)
		signTxBys, err := hex.DecodeString(signTxStr)
		if err != nil {
			return nil, fmt.Errorf("hex.DecodeString signed tx err: %s", err.Error())
		}
		return signTxBys, nil
	} else {
		return nil, fmt.Errorf("remoteRpcClient.Call err: %s", reply.Errmsg)
	}
}
