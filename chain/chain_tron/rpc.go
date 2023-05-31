package chain_tron

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/dotbitHQ/das-lib/chain/chain_evm"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/fbsobreira/gotron-sdk/pkg/proto/api"
	"github.com/fbsobreira/gotron-sdk/pkg/proto/core"
	"github.com/golang/protobuf/proto"
)

func (c *ChainTron) GetBlockNumber() (int64, error) {
	block, err := c.Client.GetNowBlock2(c.Ctx, new(api.EmptyMessage))
	if err != nil {
		return 0, err
	}
	return block.BlockHeader.RawData.Number, nil
}

func (c *ChainTron) GetBlockByNumber(blockNumber uint64) (*api.BlockExtention, error) {
	num := int64(blockNumber)
	block, err := c.Client.GetBlockByNum2(c.Ctx, &api.NumberMessage{Num: num})
	if err != nil {
		return nil, err
	}
	return block, nil
}

func (c *ChainTron) CreateTransaction(fromHex, toHex, memo string, amount int64) (*api.TransactionExtention, error) {
	fromAddr, err := hex.DecodeString(fromHex)
	if err != nil {
		return nil, fmt.Errorf("decode from hex:%s %v", fromHex, err)
	}
	toAddr, err := hex.DecodeString(toHex)
	if err != nil {
		return nil, fmt.Errorf("decode to hex:%s %v", toHex, err)
	}
	in := &core.TransferContract{
		OwnerAddress: fromAddr,
		ToAddress:    toAddr,
		Amount:       amount,
	}
	tx, err := c.Client.CreateTransaction2(c.Ctx, in)
	if err != nil {
		return nil, fmt.Errorf("create tx err:%v", err)
	}
	if tx.Result.Code != api.Return_SUCCESS {
		return nil, fmt.Errorf("create tx failed:%s", tx.Result.Message)
	}
	if memo != "" {
		tx.Transaction.RawData.Data = []byte(memo)
	}
	return tx, nil
}

// AddSign Deprecated
func (c *ChainTron) AddSign(tx *core.Transaction, private string) (*api.TransactionExtention, error) {
	pri, err := hex.DecodeString(private)
	if err != nil {
		return nil, fmt.Errorf("decode private:%v", err)
	}

	ts, err := c.Client.AddSign(c.Ctx, &core.TransactionSign{Transaction: tx, PrivateKey: pri})
	if err != nil {
		return nil, fmt.Errorf("sign err:%v", err)
	}
	if ts.Result.Code != api.Return_SUCCESS {
		return nil, fmt.Errorf("sign failed:%s", ts.Result.Message)
	}
	return ts, nil
}

func (c *ChainTron) LocalSign(tx *api.TransactionExtention, privateKey string) error {
	if tx == nil || tx.Transaction == nil {
		return fmt.Errorf("tx is nil")
	}
	rawData, err := proto.Marshal(tx.Transaction.GetRawData())
	if err != nil {
		return fmt.Errorf("proto.Marshal err: %s", err.Error())
	}
	h256h := sha256.New()
	h256h.Write(rawData)
	hash := h256h.Sum(nil)

	private, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		return fmt.Errorf("crypto.HexToECDSA err: %s", err.Error())
	}
	signData, err := crypto.Sign(hash, private)
	if err != nil {
		return fmt.Errorf("crypto.Sign err: %s", err.Error())
	}
	tx.Transaction.Signature = append(tx.Transaction.Signature, signData)
	tx.Txid = hash
	return nil
}

func (c *ChainTron) SendTransaction(in *core.Transaction) error {
	ret, err := c.Client.BroadcastTransaction(c.Ctx, in)
	if err != nil {
		return fmt.Errorf("broadcast tx err:%v", err)
	}
	if ret.Code != api.Return_SUCCESS {
		return fmt.Errorf("broadcast tx failed:%s", ret.Message)
	}
	return nil
}

func (c *ChainTron) TransferTrc20(contractHex, fromHex, toHex string, amount int64, feeLimit int64) (*api.TransactionExtention, error) {
	conAddr, err := hex.DecodeString(contractHex)
	if err != nil {
		return nil, fmt.Errorf("hex decode:%v", err)
	}
	fromAddr, err := hex.DecodeString(fromHex)
	if err != nil {
		return nil, fmt.Errorf("hex decode:%v", err)
	}

	data, err := chain_evm.PackMessage("transfer", common.HexToAddress(toHex), amount)
	if err != nil {
		return nil, fmt.Errorf("decode str:%v", err)
	}

	in := core.TriggerSmartContract{
		OwnerAddress:    fromAddr,
		ContractAddress: conAddr,
		Data:            data,
	}
	tx, err := c.Client.TriggerContract(c.Ctx, &in)
	if err != nil {
		return nil, fmt.Errorf("TriggerContract:%v", err)
	}
	if tx.Result.Code != api.Return_SUCCESS {
		return nil, fmt.Errorf("TriggerContract failed:%s %s", tx.Result.Code.String(), tx.Result.Message)
	}
	tx.Transaction.RawData.FeeLimit = feeLimit
	return tx, nil
}
