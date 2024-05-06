package core

import (
	"fmt"
	"github.com/dotbitHQ/das-lib/common"
	"github.com/nervosnetwork/ckb-sdk-go/address"
	"github.com/nervosnetwork/ckb-sdk-go/types"
)

type ChainTypeAddress struct {
	Type    string  `json:"type"` // blockchain
	KeyInfo KeyInfo `json:"key_info"`
}

type KeyInfo struct {
	CoinType common.CoinType `json:"coin_type"`
	ChainId  common.ChainId  `json:"chain_id"`
	Key      string          `json:"key"`
}

func (c *ChainTypeAddress) GetChainId(net common.DasNetType) (chainId int64) {
	switch net {
	case common.DasNetTypeMainNet:
		switch c.KeyInfo.CoinType {
		case common.CoinTypeEth:
			chainId = 1
		case common.CoinTypeBSC, common.CoinTypeBNB:
			chainId = 56
		case common.CoinTypeMatic:
			chainId = 137
		}
	default:
		switch c.KeyInfo.CoinType {
		case common.CoinTypeEth:
			chainId = 17000
		case common.CoinTypeBSC, common.CoinTypeBNB:
			chainId = 97
		case common.CoinTypeMatic:
			chainId = 80001
		}
	}
	return
}

func (c *ChainTypeAddress) FormatChainTypeAddress(net common.DasNetType, is712 bool) (*DasAddressHex, error) {
	if c.Type != "blockchain" {
		return nil, fmt.Errorf("not support type[%s]", c.Type)
	}
	dasChainType := common.FormatCoinTypeToDasChainType(c.KeyInfo.CoinType)
	if dasChainType == -1 {
		dasChainType = common.FormatChainIdToDasChainType(net, c.KeyInfo.ChainId)
	}
	if dasChainType == -1 {
		return nil, fmt.Errorf("not support coin type[%s]-chain id[%s]", c.KeyInfo.CoinType, c.KeyInfo.ChainId)
	}

	daf := DasAddressFormat{DasNetType: net}
	addrHex, err := daf.NormalToHex(DasAddressNormal{
		ChainType:     dasChainType,
		AddressNormal: c.KeyInfo.Key,
		Is712:         is712,
	})
	if err != nil {
		return nil, fmt.Errorf("address NormalToHex err: %s", err.Error())
	}

	return &addrHex, nil
}

func (c *ChainTypeAddress) FormatAnyLock() (bool, *address.ParsedAddress, error) {
	if c.KeyInfo.CoinType != common.CoinTypeCKB {
		return false, nil, nil
	}
	addrParse, err := address.Parse(c.KeyInfo.Key)
	if err != nil {
		return false, nil, fmt.Errorf("address.Parse err: %s", err.Error())
	}
	contractDispatch, err := GetDasContractInfo(common.DasContractNameDispatchCellType)
	if err != nil {
		return false, nil, fmt.Errorf("GetDasContractInfo err: %s", err.Error())
	}
	if !contractDispatch.IsSameTypeId(addrParse.Script.CodeHash) {
		return true, addrParse, nil
	}
	return false, nil, nil
}

func FormatChainTypeAddress(net common.DasNetType, chainType common.ChainType, key string) ChainTypeAddress {
	var coinType common.CoinType
	switch chainType {
	case common.ChainTypeEth:
		coinType = common.CoinTypeEth
	case common.ChainTypeTron:
		coinType = common.CoinTypeTrx
	case common.ChainTypeWebauthn:
		coinType = common.CoinTypeCKB
	case common.ChainTypeDogeCoin:
		coinType = common.CoinTypeDogeCoin
	}

	var chainId common.ChainId
	if net == common.DasNetTypeMainNet {
		switch chainType {
		case common.ChainTypeEth:
			chainId = common.ChainIdEthMainNet
		}
	} else {
		switch chainType {
		case common.ChainTypeEth:
			chainId = common.ChainIdEthTestNet
		}
	}

	return ChainTypeAddress{
		Type: "blockchain",
		KeyInfo: KeyInfo{
			CoinType: coinType,
			ChainId:  chainId,
			Key:      key,
		},
	}
}

func (c *ChainTypeAddress) FormatChainTypeAddressToScript(net common.DasNetType, is712 bool) (*types.Script, *types.Script, error) {
	addHex, err := c.FormatChainTypeAddress(net, is712)
	if err != nil {
		return nil, nil, err
	}
	daf := DasAddressFormat{DasNetType: net}
	return daf.HexToScript(*addHex)
}
