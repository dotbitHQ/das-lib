package http_api

import (
	"fmt"
	"github.com/dotbitHQ/das-lib/common"
	"github.com/dotbitHQ/das-lib/sign"
	"github.com/dotbitHQ/das-lib/txbuilder"
	"github.com/scorpiotzh/toolib"
)

type SignInfo struct {
	SignKey     string               `json:"sign_key"`               // sign tx key
	SignAddress string               `json:"sign_address,omitempty"` // sign address
	SignList    []txbuilder.SignData `json:"sign_list"`              // sign list
	MMJson      *common.MMJsonObj    `json:"mm_json"`                // 712 mmjson
}

func (s *SignInfo) SignListString() string {
	return toolib.JsonString(s.SignList)
}

func fixSignature(signMsg string) string {
	if len(signMsg) >= 132 && signMsg[130:132] == "1b" {
		signMsg = signMsg[0:130] + "00" + signMsg[132:]
	}
	if len(signMsg) >= 132 && signMsg[130:132] == "1c" {
		signMsg = signMsg[0:130] + "01" + signMsg[132:]
	}
	return signMsg
}

func VerifySignature(signType common.DasAlgorithmId, signMsg, signature, address string) (bool, string, error) {
	signOk := false
	var err error
	switch signType {
	case common.DasAlgorithmIdEth:
		signature = fixSignature(signature)
		signOk, err = sign.VerifyPersonalSignature(common.Hex2Bytes(signature), []byte(signMsg), address)
		if err != nil {
			return false, signature, fmt.Errorf("VerifyPersonalSignature err: %s", err.Error())
		}
	case common.DasAlgorithmIdTron:
		signature = fixSignature(signature)
		if address, err = common.TronHexToBase58(address); err != nil {
			return false, signature, fmt.Errorf("TronHexToBase58 err: %s [%s]", err.Error(), address)
		}
		signOk = sign.TronVerifySignature(true, common.Hex2Bytes(signature), []byte(signMsg), address)
	case common.DasAlgorithmIdEd25519:
		signOk = sign.VerifyEd25519Signature(common.Hex2Bytes(address), common.Hex2Bytes(signMsg), common.Hex2Bytes(signMsg))
	case common.DasAlgorithmIdDogeChain:
		signOk, err = sign.VerifyDogeSignature(common.Hex2Bytes(signature), []byte(signMsg), address)
		if err != nil {
			return false, signature, fmt.Errorf("VerifyDogeSignature err: %s [%s]", err.Error(), address)
		}
	case common.DasAlgorithmIdWebauthn:
		signOk, err = sign.VerifyWebauthnSignature([]byte(signMsg), common.Hex2Bytes(signature), address[20:])
		if err != nil {
			return false, signature, fmt.Errorf("VerifyWebauthnSignature err: %s [%s]", err.Error(), address)
		}
	default:
		return false, signature, fmt.Errorf("not exist sign type[%d]", signType)
	}

	if !signOk {
		return false, signature, nil
	}
	return true, signature, nil
}
