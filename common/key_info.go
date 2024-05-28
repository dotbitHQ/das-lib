package common

import (
	"encoding/hex"
	"fmt"
	"github.com/btcsuite/btcd/btcutil/base58"
	"regexp"
	"strings"
)

type CoinType string // EIP-155

// https://github.com/satoshilabs/slips/blob/master/slip-0044.md
const (
	CoinTypeBTC      CoinType = "0"
	CoinTypeEth      CoinType = "60"
	CoinTypeTrx      CoinType = "195"
	CoinTypeBNB      CoinType = "714"
	CoinTypeBSC      CoinType = "9006"
	CoinTypeMatic    CoinType = "966"
	CoinTypeDogeCoin CoinType = "3"
	CoinTypeCKB      CoinType = "309"
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
	case CoinTypeCKB:
		return ChainTypeCkb
	case CoinTypeBTC:
		return ChainTypeBitcoin
	}
	return -1
}
func FormatDasChainTypeToCoinType(chainType ChainType) CoinType {
	switch chainType {
	case ChainTypeEth:
		return CoinTypeEth
	case ChainTypeTron:
		return CoinTypeTrx
	case ChainTypeDogeCoin:
		return CoinTypeDogeCoin
	case ChainTypeCkb, ChainTypeWebauthn, ChainTypeAnyLock:
		return CoinTypeCKB
	}
	return "-1"
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

var TokenId2RecordKeyMap = map[string][]string{
	"eth_eth":        {"eth", "60"},
	"bsc_bnb":        {"bnb", "9006"},
	"tron_trx":       {"trx", "195"},
	"eth_erc20_usdt": {"eth", "60"},
	"bsc_bep20_usdt": {"bnb", "9006"},
	"doge_doge":      {"doge", "3"},
	"did_point":      {"ckb", "309"},
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
		if strings.HasPrefix(address, DogeCoinBase58PreFix) {
			if _, err := Base58CheckDecode(address, DogeCoinBase58Version); err != nil {
				return "", fmt.Errorf("Base58CheckDecode err: %s", err.Error())
			} else {
				return address, nil
			}
		} else if addr, err := Base58CheckEncode(address, DogeCoinBase58Version); err != nil {
			return "", fmt.Errorf("Base58CheckEncode err: %s", err.Error())
		} else {
			return addr, nil
		}
	case CoinTypeCKB:
		return address, nil
	}
	return "", fmt.Errorf("unknow coin-type [%s]", coinType)
}

func Base58CheckDecode(addr string, version byte) (string, error) {
	payload, v, err := base58.CheckDecode(addr)
	if err != nil {
		return "", fmt.Errorf("base58.CheckDecode err: %s[%s]", err.Error(), addr)
	} else if v != version {
		return "", fmt.Errorf("base58.CheckDecode version diff: %d[%s]", v, addr)
	}
	return hex.EncodeToString(payload), nil
}

func Base58CheckEncode(payload string, version byte) (string, error) {
	bys, err := hex.DecodeString(payload)
	if err != nil {
		return "", fmt.Errorf("payload DecodeString err: %s", err.Error())
	}
	return base58.CheckEncode(bys, version), nil
}
