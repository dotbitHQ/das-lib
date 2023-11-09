package sign

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dotbitHQ/das-lib/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
	"strings"
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
func VerifyEthSignature(sign []byte, rawByte []byte, address string) (bool, error) {
	if len(sign) != 65 { // sign check
		return false, fmt.Errorf("invalid param")
	}
	if sign[64] >= 27 {
		sign[64] -= 27
	}

	pub, err := crypto.Ecrecover(rawByte[:], sign)
	if err != nil {
		return false, err
	}
	pubKey, err := crypto.UnmarshalPubkey(pub)
	if err != nil {
		return false, err
	}

	recoveredAddr := crypto.PubkeyToAddress(*pubKey)
	//fmt.Println("recovered:", recoveredAddr.Hex(), "addr:", address)
	return strings.EqualFold(recoveredAddr.Hex(), address), nil
}

func PersonalSignature(data []byte, hexPrivateKey string) ([]byte, error) {
	l := len(data)
	if l == 0 {
		return nil, errors.New("invalid raw data")
	}

	data = append([]byte(fmt.Sprintf(common.EthMessageHeader, l)), data...)
	key, err := crypto.HexToECDSA(hexPrivateKey)
	if err != nil {
		return nil, err
	}
	tmpHash := crypto.Keccak256(data)

	return crypto.Sign(tmpHash, key)
}
func VerifyPersonalSignature(sign []byte, rawByte []byte, address string) (bool, error) {
	l := len(rawByte)
	if len(sign) != 65 || l == 0 { // sign check
		return false, fmt.Errorf("invalid param")
	}

	if sign[64] >= 27 {
		sign[64] -= 27
	}
	rawByte = append([]byte(fmt.Sprintf(common.EthMessageHeader, l)), rawByte...)
	hash := crypto.Keccak256(rawByte)

	pub, err := crypto.Ecrecover(hash[:], sign)
	if err != nil {
		return false, err
	}
	pubKey, err := crypto.UnmarshalPubkey(pub)
	if err != nil {
		return false, err
	}
	recoveredAddr := crypto.PubkeyToAddress(*pubKey)
	//fmt.Println("recovered:", recoveredAddr.Hex(), "addr:", address)
	return strings.EqualFold(recoveredAddr.Hex(), address), nil
}

func EIP712Signature(typedData apitypes.TypedData, hexPrivateKey string) ([]byte, []byte, error) {

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
	fmt.Println("sign dataHash:", common.Bytes2Hex(dataHash))

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

func VerifyEIP712Signature(typedData apitypes.TypedData, sign []byte, address string) (bool, error) {
	if len(sign) != 65 { // sign check
		return false, fmt.Errorf("invalid param")
	}
	if sign[64] >= 27 {
		sign[64] -= 27
	}

	domainSeparator, err := typedData.HashStruct("EIP712Domain", typedData.Domain.Map())
	if err != nil {
		return false, err
	}
	typedDataHash, err := typedData.HashStruct(typedData.PrimaryType, typedData.Message)
	if err != nil {
		return false, err
	}
	rawData := []byte(fmt.Sprintf("\x19\x01%s%s", string(domainSeparator), string(typedDataHash)))
	dataHash := crypto.Keccak256(rawData)
	//fmt.Println("verify dataHash:", common.Bytes2Hex(dataHash))
	pubKeyRaw, err := crypto.Ecrecover(dataHash, sign)
	if err != nil {
		return false, err
	}
	pubKey, err := crypto.UnmarshalPubkey(pubKeyRaw)
	if err != nil {
		return false, err
	}
	recoveredAddr := crypto.PubkeyToAddress(*pubKey)
	//fmt.Println("recovered:", recoveredAddr.Hex(), "addr:", address)
	return strings.EqualFold(recoveredAddr.Hex(), address), nil
}

func DoEIP712Sign(chainId int64, signMsg, private string, mmJsonObj *common.MMJsonObj) (string, error) {
	log.Info("DoSign:", chainId, signMsg)
	var signData []byte

	var obj3 apitypes.TypedData
	mmJson := mmJsonObj.String()
	oldChainId := fmt.Sprintf("chainId\":%d", chainId)
	newChainId := fmt.Sprintf("chainId\":\"%d\"", chainId)
	mmJson = strings.ReplaceAll(mmJson, oldChainId, newChainId)
	oldDigest := "\"digest\":\"\""
	newDigest := fmt.Sprintf("\"digest\":\"%s\"", signMsg)
	mmJson = strings.ReplaceAll(mmJson, oldDigest, newDigest)

	_ = json.Unmarshal([]byte(mmJson), &obj3)
	var mmHash, signature []byte
	mmHash, signature, err := EIP712Signature(obj3, private)
	if err != nil {
		return "", fmt.Errorf("EIP712Signature err: %s", err.Error())
	}

	signData = append(signature, mmHash...)

	hexChainId := fmt.Sprintf("%x", chainId)
	chainIdData := common.Hex2Bytes(fmt.Sprintf("%016s", hexChainId))
	signData = append(signData, chainIdData...)

	return common.Bytes2Hex(signData), nil
}
