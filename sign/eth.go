package sign

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/DeAccountSystems/das-lib/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/signer/core"
	"math/big"
)

func EthSignature(signType bool, data string, hexPrivateKey string) ([]byte, error) {
	if len(data) == 0 {
		return nil, errors.New("invalid raw data")
	}

	if signType {
		data = common.EthMessageHeader + data
	}

	key, err := crypto.HexToECDSA(hexPrivateKey)
	if err != nil {
		return nil, err
	}
	tmpHash := crypto.Keccak256([]byte(data))

	return crypto.Sign(tmpHash, key)
}

func EthVerifySignature(signType bool, sign []byte, rawData string, address string) bool {
	if len(sign) != 65 { // sign check
		return false
	}

	if sign[64] >= 27 {
		sign[64] = sign[64] - 27
	}

	if signType {
		rawData = common.EthMessageHeader + rawData
	}

	pubKey, err := GetSignedPubKey(rawData, sign)
	if err != nil {
		return false
	}

	recoveredAddr := crypto.PubkeyToAddress(*pubKey)

	return recoveredAddr.Hex() == address
}

func EthSignature712(str string, obj *common.MMJsonObj, digest, hexPrivateKey string) ([]byte, error) {
	obj.Message.Digest = digest
	objData, _ := json.Marshal(obj)
	var typedData core.TypedData
	_ = json.Unmarshal(objData, &typedData)
	err := json.Unmarshal([]byte(str), &typedData)
	if err != nil {
		return nil, err
	}

	domainMap := typedData.Domain.Map()
	domainMap["chainId"] = big.NewInt(5)
	log.Info("typedData.Message:", typedData.Message)
	log.Info("typedData.Domain.Map():", domainMap)

	typedDataHash, _ := typedData.HashStruct(typedData.PrimaryType, typedData.Message)
	domainSeparator, _ := typedData.HashStruct("EIP712Domain", domainMap)
	rawData := []byte(fmt.Sprintf("\x19\x01%s%s", string(domainSeparator), string(typedDataHash)))
	dataHash := crypto.Keccak256(rawData)

	key, err := crypto.HexToECDSA(hexPrivateKey)
	if err != nil {
		return nil, err
	}
	signature, err := crypto.Sign(dataHash, key)
	if err != nil {
		return nil, err
	}

	if signature[64] < 27 {
		signature[64] += 27
	}

	return signature, nil
}

func EthVerifySignature712(obj *common.MMJsonObj, sign []byte, digest, address string) bool {
	if len(sign) != 65 { // sign check
		return false
	}

	if sign[64] >= 27 {
		sign[64] = sign[64] - 27
	}

	obj.Message.Digest = digest
	objData, _ := json.Marshal(obj)
	var typedData core.TypedData
	_ = json.Unmarshal(objData, &typedData)

	typedDataHash, _ := typedData.HashStruct(typedData.PrimaryType, typedData.Message)
	domainSeparator, _ := typedData.HashStruct("EIP712Domain", typedData.Domain.Map())
	rawData := []byte(fmt.Sprintf("\x19\x01%s%s", string(domainSeparator), string(typedDataHash)))
	dataHash := crypto.Keccak256(rawData)

	pubKeyRaw, err := crypto.Ecrecover(dataHash, sign)
	if err != nil {
		return false
	}
	pubKey, err := crypto.UnmarshalPubkey(pubKeyRaw)
	if err != nil {
		return false
	}
	recoveredAddr := crypto.PubkeyToAddress(*pubKey)

	return recoveredAddr.Hex() == address
}
