package sign

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/dotbitHQ/das-lib/common"
	"math/big"
)

type ClientDataJson struct {
	Type        string `json:"type"`
	Challenge   string `json:"challenge"`
	Origin      string `json:"origin"`
	CrossOrigin string `json:"croosOrigin"`
}

func VerifyWebauthnSignature(challenge, dataBys []byte, signAddressPk1 string) (res bool, err error) {
	if len(signAddressPk1) != 20 {
		return false, fmt.Errorf("signAddressPk1 length error : ", signAddressPk1)
	}
	index, indexLen1, indexLen2, dataLen := uint16(0), uint16(1), uint16(2), uint16(0)
	//pkIndex
	dataLen = uint16(dataBys[index])
	pkIndex := dataBys[index+indexLen1 : index+indexLen1+dataLen]
	log.Info("pkIndex: ", pkIndex[0])
	index = index + indexLen1 + dataLen

	dataLen = uint16(dataBys[index])
	signature := dataBys[index+indexLen1 : index+indexLen1+dataLen]
	log.Info("signature: ", common.Bytes2Hex(signature))
	index = index + indexLen1 + dataLen

	dataLen = uint16(dataBys[index])
	pubKeyBytes := dataBys[index+indexLen1 : index+indexLen1+dataLen]
	log.Info("pubKeyBytes: ", common.Bytes2Hex(pubKeyBytes))
	index = index + indexLen1 + dataLen
	//verify pubKey
	var pubKey ecdsa.PublicKey
	pubKey.Curve = elliptic.P256()
	pubKey.X = new(big.Int).SetBytes(pubKeyBytes[:32])
	pubKey.Y = new(big.Int).SetBytes(pubKeyBytes[32:])
	pk1 := common.CalculatePk1(&pubKey)
	if signAddressPk1 != hex.EncodeToString(pk1) {
		log.Info("signAddressPk1: ", signAddressPk1, " pk1: ", hex.EncodeToString(pk1))
		return false, nil
	}

	dataLen = uint16(dataBys[index])
	authnticatorData := dataBys[index+indexLen1 : index+indexLen1+dataLen]
	log.Info("authnticatorData: ", common.Bytes2Hex(authnticatorData))
	index = index + indexLen1 + dataLen

	dataLen = binary.LittleEndian.Uint16(dataBys[index : index+indexLen2])
	clientDataJsonBys := dataBys[index+indexLen2 : index+indexLen2+dataLen]
	log.Info("clientDataJsonBys: ", common.Bytes2Hex(clientDataJsonBys))
	index = index + indexLen1 + dataLen

	log.Info("json clientDataJsonData ", string(clientDataJsonBys))
	var clientDataJsonData map[string]interface{}
	err = json.Unmarshal(clientDataJsonBys, &clientDataJsonData)
	if err != nil {
		log.Info("unmarshal err :", err)
		return false, fmt.Errorf("json.Unmarshal(clientDataJsonData) err: %s", err.Error())
	}
	if _, ok := clientDataJsonData["challenge"]; !ok {
		return false, fmt.Errorf("There is no challenge in clientDataJson")
	}

	// verify challenge
	challengeBase64url := base64.URLEncoding.EncodeToString(challenge)
	log.Info("challengeBase64url: ", challengeBase64url)
	log.Info("clientDataJsonData.challenge", clientDataJsonData["challenge"])
	if challengeBase64url != clientDataJsonData["challenge"] {
		return false, fmt.Errorf("clientDataJsonData.challenge  error")
	}

	clientDataJsonHash := sha256.Sum256(clientDataJsonBys)
	signMsg := append(authnticatorData, clientDataJsonHash[:]...)
	hash := sha256.Sum256(signMsg)
	R := new(big.Int).SetBytes(signature[:32])
	S := new(big.Int).SetBytes(signature[32:])
	pubkey := new(ecdsa.PublicKey)
	pubkey.X = new(big.Int).SetBytes(pubKeyBytes[:32])
	pubkey.Y = new(big.Int).SetBytes(pubKeyBytes[32:])
	res, err = VerifyEcdsaP256Signature(hash[:], R, S, pubkey)
	return
}

func VerifyEcdsaP256Signature(hash []byte, R, S *big.Int, pubkey *ecdsa.PublicKey) (res bool, err error) {
	if len(hash) != 32 {
		return false, fmt.Errorf("hash length error: ", hash)
	}
	//P' = (z*G*S^-1 + Qa*R*S^-1) mod p
	curve := elliptic.P256()
	N := curve.Params().N
	z := new(big.Int).SetBytes(hash)
	//s^-1
	sInv := new(big.Int).ModInverse(S, N)
	u1 := new(ecdsa.PublicKey)
	//G*z
	u1.X, u1.Y = curve.ScalarBaseMult(z.Bytes())
	//G*z*s^-1
	u1.X, u1.Y = curve.ScalarMult(u1.X, u1.Y, sInv.Bytes())

	rs := new(big.Int).Mul(R, sInv)
	u2 := new(ecdsa.PublicKey)

	u2.X, u2.Y = curve.ScalarMult(pubkey.X, pubkey.Y, rs.Bytes())

	X, _ := curve.Add(u1.X, u1.Y, u2.X, u2.Y)
	if X.Cmp(R) == 0 {
		return true, nil
	} else {
		return false, nil
	}
}
