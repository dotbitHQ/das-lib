package sign

import (
	"crypto/ecdsa"
	"errors"
	"github.com/DeAccountSystems/das-lib/common"
	"github.com/pranksteess/go-ethereum/crypto"
)

func TronSignature(signType bool, data string, hexPrivateKey string) ([]byte, error) {
	if len(data) == 0 {
		return nil, errors.New("invalid raw data")
	}

	if signType {
		data = common.TronMessageHeader + data
	}

	tmpHash := crypto.Keccak256([]byte(data))

	privateKey, err := crypto.HexToECDSA(hexPrivateKey)
	if err != nil {
		return nil, err
	}

	signData, err := crypto.Sign(tmpHash[:], privateKey)
	if err == nil && len(signData) == 65 && signData[64] < 27 {
		signData[64] = signData[64] + 27
	}
	return signData, err
}

func TronVerifySignature(signType bool, sign []byte, rawData string, base58Addr string) bool {
	if len(sign) != 65 { // sign check
		return false
	}

	if sign[64] >= 27 {
		sign[64] = sign[64] - 27
	}

	if signType {
		rawData = common.TronMessageHeader + rawData
	}

	pubKey, err := GetSignedPubKey(rawData, sign)
	if err != nil {
		return false
	}

	recoveredAddr := crypto.PubkeyToAddress(*pubKey)
	tronAddr := common.TronPreFix + recoveredAddr.String()[2:]

	base58address, err := common.TronHexToBase58(tronAddr)
	if err != nil {
		return false
	}
	return base58address == base58Addr
}

func GetSignedPubKey(rawData string, sign []byte) (*ecdsa.PublicKey, error) {
	if len(sign) != 65 { // sign check
		return nil, errors.New("invalid transaction signature, should be 65 length bytes")
	}
	rawByte := []byte(rawData)
	hash := crypto.Keccak256(rawByte)

	pub, err := crypto.Ecrecover(hash[:], sign)
	if err != nil {
		return nil, err
	}
	return crypto.UnmarshalPubkey(pub)
}
