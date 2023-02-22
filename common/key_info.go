package common

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/btcsuite/btcd/btcutil/base58"
	"regexp"
	"strings"
)

type CoinType string // EIP-155

// https://github.com/satoshilabs/slips/blob/master/slip-0044.md
const (
	CoinTypeEth      CoinType = "60"
	CoinTypeTrx      CoinType = "195"
	CoinTypeBNB      CoinType = "714"
	CoinTypeBSC      CoinType = "9006"
	CoinTypeMatic    CoinType = "966"
	CoinTypeDogeCoin CoinType = "3"
)

type ChainId string //BIP-44

const (
	ChainIdEthMainNet     ChainId = "1"
	ChainIdBscMainNet     ChainId = "56"
	ChainIdPolygonMainNet ChainId = "137"

	ChainIdEthTestNet     ChainId = "5" // Goerli
	ChainIdBscTestNet     ChainId = "97"
	ChainIdPolygonTestNet ChainId = "80001"
)

func FormatCoinTypeToDasChainType(coinType CoinType) ChainType {
	switch coinType {
	case CoinTypeEth, CoinTypeBNB, CoinTypeBSC, CoinTypeMatic:
		return ChainTypeEth
	case CoinTypeTrx:
		return ChainTypeTron
	case CoinTypeDogeCoin:
		return ChainTypeDogeCoin
	}
	return -1
}

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

// coin-type to address-key

var RecordsAddressCoinTypeMap = map[string]string{
	"address.1815":  "address.ada",
	"address.118":   "address.atom",
	"address.9000":  "address.avalanche",
	"address.145":   "address.bch",
	"address.9006":  "address.bsc",
	"address.236":   "address.bsv",
	"address.0":     "address.btc",
	"address.52752": "address.celo",
	"address.309":   "address.ckb",
	"address.5":     "address.dash",
	"address.223":   "address.dfinity",
	"address.3":     "address.doge",
	"address.354":   "address.dot",
	"address.194":   "address.eos",
	"address.61":    "address.etc",
	"address.60":    "address.eth",
	"address.461":   "address.fil",
	"address.539":   "address.flow",
	"address.1010":  "address.heco",
	"address.291":   "address.iost",
	"address.4218":  "address.iota",
	"address.434":   "address.ksm",
	"address.2":     "address.ltc",
	"address.397":   "address.near",
	"address.966":   "address.polygon",
	"address.1991":  "address.sc",
	"address.501":   "address.sol",
	"address.5757":  "address.stacks",
	"address.330":   "address.terra",
	"address.195":   "address.trx",
	"address.818":   "address.vet",
	"address.43":    "address.xem",
	"address.148":   "address.xlm",
	"address.128":   "address.xmr",
	"address.144":   "address.xrp",
	"address.1729":  "address.xtz",
	"address.133":   "address.zec",
	"address.313":   "address.zil",
}
var RecordsAddressKeyMap = map[string]string{
	"address.xmr":       "address.128",
	"address.xrp":       "address.144",
	"address.iota":      "address.4218",
	"address.ltc":       "address.2",
	"address.xem":       "address.43",
	"address.terra":     "address.330",
	"address.celo":      "address.52752",
	"address.eth":       "address.60",
	"address.sc":        "address.1991",
	"address.trx":       "address.195",
	"address.zil":       "address.313",
	"address.ada":       "address.1815",
	"address.atom":      "address.118",
	"address.bsv":       "address.236",
	"address.fil":       "address.461",
	"address.near":      "address.397",
	"address.sol":       "address.501",
	"address.eos":       "address.194",
	"address.etc":       "address.61",
	"address.zec":       "address.133",
	"address.avalanche": "address.9000",
	"address.dash":      "address.5",
	"address.dfinity":   "address.223",
	"address.ckb":       "address.309",
	"address.stacks":    "address.5757",
	"address.bsc":       "address.9006",
	"address.dot":       "address.354",
	"address.polygon":   "address.966",
	"address.flow":      "address.539",
	"address.heco":      "address.1010",
	"address.iost":      "address.291",
	"address.ksm":       "address.434",
	"address.vet":       "address.818",
	"address.bch":       "address.145",
	"address.btc":       "address.0",
	"address.doge":      "address.3",
	"address.xlm":       "address.148",
	"address.xtz":       "address.1729",
	//

}

func ConvertRecordsAddressKey(addressKey string) string {
	if item, ok := RecordsAddressKeyMap[addressKey]; ok {
		return item
	}
	return addressKey
}

func ConvertRecordsAddressCoinType(addressCoinType string) string {
	if item, ok := RecordsAddressCoinTypeMap[addressCoinType]; ok {
		return item
	}
	return addressCoinType
}

func FormatAddressByCoinType(coinType string, address string) (string, error) {
	switch CoinType(coinType) {
	case CoinTypeEth, CoinTypeBNB, CoinTypeBSC, CoinTypeMatic:
		if ok, err := regexp.MatchString("^0x[0-9a-fA-F]{40}$", address); err != nil {
			return "", fmt.Errorf("regexp.MatchString err: %s", err.Error())
		} else if ok {
			return address, nil
		} else {
			return "", fmt.Errorf("regexp.MatchString false")
		}
	case CoinTypeTrx:
		if strings.HasPrefix(address, TronBase58PreFix) {
			if _, err := TronBase58ToHex(address); err != nil {
				return "", fmt.Errorf("TronBase58ToHex err: %s", err.Error())
			} else {
				return address, nil
			}
		} else if strings.HasPrefix(address, TronPreFix) {
			if addr, err := TronHexToBase58(address); err != nil {
				return "", fmt.Errorf("TronHexToBase58 err: %s", err.Error())
			} else {
				return addr, nil
			}
		}
	case CoinTypeDogeCoin:
		if _, _, err := base58.CheckDecode(address); err != nil {
			return "", fmt.Errorf("base58.CheckDecode err: %s", err.Error())
		}
	}
	return "", fmt.Errorf("unknow coin-type [%s]", coinType)
}

func FormatDogeCoinAddressToPayload(addr string) (payload string, err error) {
	decode := base58.Decode(addr)
	payloadBys := decode[1 : len(decode)-4]
	//fmt.Println(hex.EncodeToString(decode))
	payload = hex.EncodeToString(payloadBys)
	//fmt.Println(payload)

	h := sha256.Sum256(decode[:len(decode)-4])
	h2 := sha256.Sum256(h[:])
	if bytes.Compare(h2[:4], decode[len(decode)-4:]) != 0 {
		err = fmt.Errorf("failed to checksum")
		return
	}
	return
}
