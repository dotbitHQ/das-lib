package remote_sign

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/btcsuite/btcd/wire"
	"github.com/dotbitHQ/das-lib/common"
	"github.com/dotbitHQ/das-lib/http_api"
	"github.com/dotbitHQ/das-lib/sign"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
)

type SignType int

const (
	SignTypeTx     SignType = 0
	SignTypeMsg    SignType = 1
	SignTypeETH712 SignType = 2
)

type ReqRemoteSign struct {
	SignType   SignType          `json:"sign_type"`
	Address    string            `json:"address"`
	EvmChainID int64             `json:"evm_chain_id"`
	Data       string            `json:"data"`
	MMJson     *common.MMJsonObj `json:"mm_json"`
}

type RespRemoteSign struct {
	Data string `json:"data"`
}

func RemoteSign(url string, req ReqRemoteSign) (*http_api.ApiResp, *RespRemoteSign, error) {
	var data RespRemoteSign
	resp, err := http_api.SendReqV2(url, &req, &data)
	if err != nil {
		return nil, nil, fmt.Errorf("http_api.SendReqV2 err: %s", err.Error())
	}
	return resp, &data, nil
}

func SignTxForCKBHandle(url, addr string) sign.HandleSignCkbMessage {
	return func(data string) ([]byte, error) {
		return SignTxForCKB(url, addr, data)
	}
}

func SignTxForCKB(url, addr, data string) ([]byte, error) {
	resp, res, err := RemoteSign(url, ReqRemoteSign{
		SignType:   SignTypeTx,
		Address:    addr,
		EvmChainID: 0,
		Data:       data,
	})
	if err != nil {
		return nil, fmt.Errorf("RemoteSign err: %s", err.Error())
	}
	if resp.ErrNo != http_api.ApiCodeSuccess {
		return nil, fmt.Errorf("RemoteSign fail code: %d, msg: %s", resp.ErrNo, resp.ErrMsg)
	}
	bys, err := hex.DecodeString(res.Data)
	if err != nil {
		return nil, fmt.Errorf("hex.DecodeString err: %s", err.Error())
	}
	return bys, nil
}

func SignTxForEVM(url, addr string, evmChainId int64, tx *types.Transaction) (*types.Transaction, error) {
	dataBys, err := rlp.EncodeToBytes(tx)
	if err != nil {
		return nil, fmt.Errorf("rlp.EncodeToBytes err: %s", err.Error())
	}

	resp, res, err := RemoteSign(url, ReqRemoteSign{
		SignType:   SignTypeTx,
		Address:    addr,
		EvmChainID: evmChainId,
		Data:       hex.EncodeToString(dataBys),
	})
	if err != nil {
		return nil, fmt.Errorf("RemoteSign err: %s", err.Error())
	}
	if resp.ErrNo != http_api.ApiCodeSuccess {
		return nil, fmt.Errorf("RemoteSign fail code: %d, msg: %s", resp.ErrNo, resp.ErrMsg)
	}
	bys, err := hex.DecodeString(res.Data)
	if err != nil {
		return nil, fmt.Errorf("hex.DecodeString err: %s", err.Error())
	}
	sigTx := &types.Transaction{}
	if err = rlp.DecodeBytes(bys, sigTx); err != nil {
		return nil, fmt.Errorf("rlp.DecodeBytes err: %s", err.Error())
	}
	return sigTx, nil
}

func SignTxForTRON(url, addr string, data []byte) ([]byte, error) {
	resp, res, err := RemoteSign(url, ReqRemoteSign{
		SignType:   SignTypeTx,
		Address:    addr,
		EvmChainID: 0,
		Data:       hex.EncodeToString(data),
	})
	if err != nil {
		return nil, fmt.Errorf("RemoteSign err: %s", err.Error())
	}
	if resp.ErrNo != http_api.ApiCodeSuccess {
		return nil, fmt.Errorf("RemoteSign fail code: %d, msg: %s", resp.ErrNo, resp.ErrMsg)
	}
	bys, err := hex.DecodeString(res.Data)
	if err != nil {
		return nil, fmt.Errorf("hex.DecodeString err: %s", err.Error())
	}
	return bys, nil
}

func SignTxForDOGE(url, addr string, tx *wire.MsgTx) (*wire.MsgTx, error) {
	data := bytes.NewBuffer(make([]byte, 0, tx.SerializeSizeStripped()))
	if err := tx.SerializeNoWitness(data); err != nil {
		return nil, fmt.Errorf("SerializeNoWitness err: %s", err.Error())
	}
	resp, res, err := RemoteSign(url, ReqRemoteSign{
		SignType:   SignTypeTx,
		Address:    addr,
		EvmChainID: 0,
		Data:       hex.EncodeToString(data.Bytes()),
	})
	if err != nil {
		return nil, fmt.Errorf("RemoteSign err: %s", err.Error())
	}
	if resp.ErrNo != http_api.ApiCodeSuccess {
		return nil, fmt.Errorf("RemoteSign fail code: %d, msg: %s", resp.ErrNo, resp.ErrMsg)
	}
	bys, err := hex.DecodeString(res.Data)
	if err != nil {
		return nil, fmt.Errorf("hex.DecodeString err: %s", err.Error())
	}
	var sigTx wire.MsgTx
	if err = sigTx.DeserializeNoWitness(bytes.NewReader(bys)); err != nil {
		return nil, fmt.Errorf("sigTx.DeserializeNoWitness err: %s", err.Error())
	}
	return &sigTx, nil
}

func SignTxFor712(url, addr, data string, evmChainId int64, mmJson *common.MMJsonObj) (string, error) {
	resp, res, err := RemoteSign(url, ReqRemoteSign{
		SignType:   SignTypeETH712,
		Address:    addr,
		EvmChainID: evmChainId,
		Data:       data,
		MMJson:     mmJson,
	})
	if err != nil {
		return "", fmt.Errorf("RemoteSign err: %s", err.Error())
	}
	if resp.ErrNo != http_api.ApiCodeSuccess {
		return "", fmt.Errorf("RemoteSign fail code: %d, msg: %s", resp.ErrNo, resp.ErrMsg)
	}
	return res.Data, nil
}
