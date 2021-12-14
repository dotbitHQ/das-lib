package sign

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/DeAccountSystems/das-lib/common"
	"github.com/pranksteess/go-ethereum/crypto"
	"github.com/pranksteess/go-ethereum/signer/core"
)
func EthSignature(data []byte, hexPrivateKey string) ([]byte, error) {
	if len(data) == 0 {
		return nil, errors.New("invalid raw data")
	}
	key, err := crypto.HexToECDSA(hexPrivateKey)
	if err != nil {
		return nil, err
	}

	return crypto.Sign(data, key)
}

func PersonalSignature(data string, hexPrivateKey string) ([]byte, error) {
	if len(data) == 0 {
		return nil, errors.New("invalid raw data")
	}

	data = common.EthMessageHeader + data

	key, err := crypto.HexToECDSA(hexPrivateKey)
	if err != nil {
		return nil, err
	}
	tmpHash := crypto.Keccak256([]byte(data))

	return crypto.Sign(tmpHash, key)
}
func VerifyPersonalSignature(sign []byte, rawData string, address string) bool {
	if len(sign) != 65 { // sign check
		return false
	}

	if sign[64] >= 27 {
		sign[64] = sign[64] - 27
	}

	rawData = common.EthMessageHeader + rawData

	pubKey, err := GetSignedPubKey(rawData, sign)
	if err != nil {
		return false
	}

	recoveredAddr := crypto.PubkeyToAddress(*pubKey)

	return recoveredAddr.Hex() == address
}
func VerifyEthSignature(sign []byte, rawData string, address string) bool {
	if len(sign) != 65 { // sign check
		return false
	}

	if sign[64] >= 27 {
		sign[64] = sign[64] - 27
	}

	pubKey, err := GetSignedPubKey(rawData, sign)
	if err != nil {
		return false
	}

	recoveredAddr := crypto.PubkeyToAddress(*pubKey)

	return recoveredAddr.Hex() == address
}

func EIP712Signature(typedData core.TypedData, hexPrivateKey string) ([]byte, []byte, error) {

	domainSeparator, err := typedData.HashStruct("EIP712Domain", typedData.Domain.Map())
	if err != nil {
		return nil, nil, err
	}
	typedDataHash, err := typedData.HashStruct(typedData.PrimaryType, typedData.Message)
	if err != nil {
		return nil, nil, err
	}
	fmt.Println("domainSeparator: ", common.Bytes2Hex(domainSeparator), "typedDataHash: ", common.Bytes2Hex(typedDataHash))
	rawData := []byte(fmt.Sprintf("\x19\x01%s%s", string(domainSeparator), string(typedDataHash)))
	dataHash := crypto.Keccak256(rawData)
	key, err := crypto.HexToECDSA(hexPrivateKey)
	if err != nil {
		return nil, nil, err
	}

	signature, err := crypto.Sign(dataHash, key)
	if err != nil {
		return nil, nil, err
	}

	if signature[64] < 27 {
		signature[64] += 27
	}

	return dataHash, signature, nil
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
