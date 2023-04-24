package chain_evm

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/scorpiotzh/mylog"
	"github.com/shopspring/decimal"
	"strings"
	"sync"
)

var (
	log = mylog.NewLogger("chain_evm", mylog.LevelDebug)
)

type ChainEvm struct {
	Client         *ethclient.Client
	Ctx            context.Context
	Node           string
	RefundAddFee   float64
	NotEnabledLock bool

	lock sync.Mutex
}

func NewChainEvm(ctx context.Context, node string, refundAddFee float64) (*ChainEvm, error) {
	ethClient, err := ethclient.Dial(node)
	if err != nil {
		return nil, fmt.Errorf("ethclient.Dial err: %s", err.Error())
	}
	return &ChainEvm{
		Client:       ethClient,
		Ctx:          ctx,
		Node:         node,
		RefundAddFee: refundAddFee,
		lock:         sync.Mutex{},
	}, nil
}

func (c *ChainEvm) EstimateGas(from, to string, value decimal.Decimal, input []byte, addFee float64) (gasPrice, gasLimit decimal.Decimal, err error) {
	fromAddr := common.HexToAddress(from)
	toAddr := common.HexToAddress(to)
	call := ethereum.CallMsg{From: fromAddr, To: &toAddr, Value: value.BigInt(), Data: input}
	limit, err := c.Client.EstimateGas(c.Ctx, call)
	if err != nil {
		return decimal.Zero, decimal.Zero, fmt.Errorf("EstimateGas err: %s", err.Error())
	}
	gasLimit, _ = decimal.NewFromString(fmt.Sprintf("%d", limit))
	fee, err := c.Client.SuggestGasPrice(c.Ctx)
	if err != nil {
		return decimal.Zero, decimal.Zero, fmt.Errorf("SuggestGasPrice err: %s", err.Error())
	}
	gasPrice, _ = decimal.NewFromString(fmt.Sprintf("%d", fee))

	log.Info("EstimateGas:", from, to, value, gasPrice, gasLimit, addFee)
	if addFee > 1 && addFee < 5 {
		gasPrice = gasPrice.Mul(decimal.NewFromFloat(addFee))
	}
	return
}

func (c *ChainEvm) NewTransaction(from, to string, value decimal.Decimal, data []byte, nonce uint64, gasPrice, gasLimit decimal.Decimal) (*types.Transaction, error) {
	toAddr := common.HexToAddress(to)
	log.Info("NewTransaction:", from, to, value, nonce, gasPrice, gasLimit)
	tx := types.NewTx(&types.LegacyTx{
		Nonce:    nonce,
		To:       &toAddr,
		Value:    value.BigInt(),
		Gas:      gasLimit.BigInt().Uint64(),
		GasPrice: gasPrice.BigInt(),
		Data:     data,
	})
	return tx, nil
}

func (c *ChainEvm) NonceAt(address string) (uint64, error) {
	return c.Client.NonceAt(c.Ctx, common.HexToAddress(address), nil)
}

func (c *ChainEvm) SignWithPrivateKey(private string, tx *types.Transaction) (*types.Transaction, error) {
	privateKey, err := crypto.HexToECDSA(HexFormat(private))
	if err != nil {
		return nil, fmt.Errorf("crypto.HexToECDSA err: %s", err.Error())
	}

	chainID, err := c.Client.NetworkID(c.Ctx)
	if err != nil {
		return nil, fmt.Errorf("NetworkID err: %s", err.Error())
	}
	sigTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		return nil, fmt.Errorf("SignTx err: %s", err.Error())
	}
	return sigTx, nil
}

func (c *ChainEvm) SendTransaction(tx *types.Transaction) error {
	return c.Client.SendTransaction(c.Ctx, tx)
}

// PackMessage
// go build -ldflags -s -v -o main cmd/abigen/*.go
// ./main --abi erc20.json --pkg chain_evm --type Erc20 --out erc20.go --alias _totalSupply=TotalSupply1
func PackMessage(name string, args ...interface{}) ([]byte, error) {
	cAbi, err := abi.JSON(strings.NewReader(Erc20MetaData.ABI))
	if err != nil {
		return nil, fmt.Errorf("new abi instance err:%v", err)
	}
	data, err := cAbi.Pack(name, args...)
	if err != nil {
		return nil, fmt.Errorf("abi package err:%s-%v", name, err)
	}
	return data, nil
}
