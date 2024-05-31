package common

type DasNetType = int

const (
	DasNetTypeMainNet  DasNetType = 1
	DasNetTypeTestnet2 DasNetType = 2
	DasNetTypeTestnet3 DasNetType = 3
)

type DasAlgorithmId int

const (
	DasAlgorithmIdCkb       DasAlgorithmId = 0
	DasAlgorithmIdCkbMulti  DasAlgorithmId = 1
	DasAlgorithmIdCkbSingle DasAlgorithmId = 2
	DasAlgorithmIdEth       DasAlgorithmId = 3
	DasAlgorithmIdTron      DasAlgorithmId = 4
	DasAlgorithmIdEth712    DasAlgorithmId = 5
	DasAlgorithmIdEd25519   DasAlgorithmId = 6
	DasAlgorithmIdDogeChain DasAlgorithmId = 7
	DasAlgorithmIdWebauthn  DasAlgorithmId = 8
	DasAlgorithmIdBitcoin   DasAlgorithmId = 9
	DasAlgorithmIdAnyLock   DasAlgorithmId = 99
)

type DasSubAlgorithmId int

const (
	DasWebauthnSubAlgorithmIdES256     DasSubAlgorithmId = 7
	DasSubAlgorithmIdBitcoinP2PKH      DasSubAlgorithmId = 1
	DasSubAlgorithmIdBitcoinP2WPKH     DasSubAlgorithmId = 2
	DasSubAlgorithmIdBitcoinP2SHP2WPKH DasSubAlgorithmId = 3
	DasSubAlgorithmIdBitcoinP2TR       DasSubAlgorithmId = 4
)

func (d DasAlgorithmId) ToCoinType() CoinType {
	switch d {
	case DasAlgorithmIdEth, DasAlgorithmIdEth712:
		return CoinTypeEth
	case DasAlgorithmIdTron:
		return CoinTypeTrx
	case DasAlgorithmIdDogeChain:
		return CoinTypeDogeCoin
	case DasAlgorithmIdWebauthn:
		return CoinTypeCKB
	case DasAlgorithmIdBitcoin:
		return CoinTypeBTC
	case DasAlgorithmIdAnyLock:
		return CoinTypeCKB
	default:
		return ""
	}
}

func (d DasAlgorithmId) Bytes() []byte {
	return []byte{uint8(d)}
}

func (d DasAlgorithmId) ToSoScriptType() SoScriptType {
	switch d {
	case DasAlgorithmIdCkbSingle:
		return SoScriptTypeCkbSingle
	case DasAlgorithmIdCkbMulti:
		return SoScriptTypeCkbMulti
	case DasAlgorithmIdEth, DasAlgorithmIdEth712:
		return SoScriptTypeEth
	case DasAlgorithmIdTron:
		return SoScriptTypeTron
	case DasAlgorithmIdEd25519:
		return SoScriptTypeEd25519
	case DasAlgorithmIdDogeChain:
		return SoScriptTypeDogeCoin
	case DasAlgorithmIdBitcoin:
		return SoScriptBitcoin
	case DasAlgorithmIdWebauthn:
		return SoScriptWebauthn
	default:
		return SoScriptTypeCkbSingle
	}
}

func (d DasAlgorithmId) ToChainType() ChainType {
	switch d {
	case DasAlgorithmIdCkbSingle:
		return ChainTypeCkbSingle
	case DasAlgorithmIdCkbMulti:
		return ChainTypeCkbMulti
	case DasAlgorithmIdEth, DasAlgorithmIdEth712:
		return ChainTypeEth
	case DasAlgorithmIdTron:
		return ChainTypeTron
	case DasAlgorithmIdEd25519:
		return ChainTypeMixin
	case DasAlgorithmIdDogeChain:
		return ChainTypeDogeCoin
	case DasAlgorithmIdBitcoin:
		return ChainTypeBitcoin
	case DasAlgorithmIdWebauthn:
		return ChainTypeWebauthn
	case DasAlgorithmIdAnyLock:
		return ChainTypeAnyLock
	default:
		return ChainTypeCkb
	}
}
