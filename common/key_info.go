package common

import (
	"fmt"
	"github.com/DeAccountSystems/das-lib/core"
)

type CoinType string // EIP-155

const (
	CoinTypeEth   CoinType = "60"
	CoinTypeTrx   CoinType = "195"
	CoinTypeBNB   CoinType = "714"
	CoinTypeMatic CoinType = "966"
)

func FormatCoinTypeToDasChainType(coinType CoinType) ChainType {
	switch coinType {
	case CoinTypeEth, CoinTypeBNB, CoinTypeMatic:
		return ChainTypeEth
	case CoinTypeTrx:
		return ChainTypeTron
	}
	return -1
}

type ChainId string //BIP-44

const (
	ChainIdEthMainNet     ChainId = "1"
	ChainIdBscMainNet     ChainId = "56"
	ChainIdPolygonMainNet ChainId = "137"

	ChainIdEthTestNet     ChainId = "5" // Goerli
	ChainIdBscTestNet     ChainId = "97"
	ChainIdPolygonTestNet ChainId = "80001"
)

func FormatChainIdToDasChainType(netType DasNetType, chainId ChainId) ChainType {
	if netType == DasNetTypeMainNet {
		switch chainId {
		case ChainIdEthMainNet, ChainIdBscMainNet, ChainIdPolygonMainNet:
			return ChainTypeEth
		}
	} else {
		switch chainId {
		case ChainIdEthTestNet, ChainIdBscTestNet, ChainIdPolygonTestNet:
			return ChainTypeEth
		}
	}
	return -1
}

type ChainTypeAddress struct {
	Type    string  `json:"type"` // blockchain
	KeyInfo KeyInfo `json:"key_info"`
}

type KeyInfo struct {
	CoinType CoinType `json:"coin_type"`
	ChainId  ChainId  `json:"chain_id"`
	Key      string   `json:"key"`
}

func (c *ChainTypeAddress) FormatChainTypeAddress(net DasNetType, is712 bool) (*core.DasAddressHex, error) {
	if c.Type != "blockchain" {
		return nil, fmt.Errorf("not support type[%s]", c.Type)
	}
	dasChainType := FormatCoinTypeToDasChainType(c.KeyInfo.CoinType)
	if dasChainType == -1 {
		dasChainType = FormatChainIdToDasChainType(net, c.KeyInfo.ChainId)
	}
	if dasChainType == -1 {
		return nil, fmt.Errorf("not support coin type[%s]-chain id[%s]", c.KeyInfo.CoinType, c.KeyInfo.ChainId)
	}

	daf := core.DasAddressFormat{DasNetType: net}
	addrHex, err := daf.NormalToHex(core.DasAddressNormal{
		ChainType:     dasChainType,
		AddressNormal: c.KeyInfo.Key,
		Is712:         is712,
	})
	if err != nil {
		return nil, fmt.Errorf("address NormalToHex err: %s", err.Error())
	}

	return &addrHex, nil
}
