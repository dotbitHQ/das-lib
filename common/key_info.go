package common

type CoinType string // EIP-155

const (
	CoinTypeEth   CoinType = "60"
	CoinTypeTrx   CoinType = "195"
	CoinTypeBNB   CoinType = "714"
	CoinTypeMatic CoinType = "966"
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
	case CoinTypeEth, CoinTypeBNB, CoinTypeMatic:
		return ChainTypeEth
	case CoinTypeTrx:
		return ChainTypeTron
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

var AddressKeyMap = map[string]string{
	"address.1815":  "address.ada",
	"address.118":   "address.atom",
	"address.9000":  "address.avalanche",
	"address.145":   "address.bch",
	"address.519":   "address.bsc",
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

func ConvertAddressKey(addressKey string) string {
	if item, ok := AddressKeyMap[addressKey]; ok {
		return item
	}
	return addressKey
}
