package sign

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/sha256"
	"encoding/base64"
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

func VerifyWebauthnSignature(challenge, webauthnSignMsgBytes []byte, signAddress string) (res bool, err error) {

	//webauthnSignMsgBytes := common.Hex2Bytes(webauthnSignMsg)
	fmt.Println(webauthnSignMsgBytes)
	//[1 0 64 88 9 250 76 123 95 7 1 140 63 169 80 141 189 172 220 124 5 157 223 16 14 90 85 16 91 12 251 193 57 16 9 153 253 44 94 70 149 163 104 101 10 64 70 24 184 226 143 120 40 127 162 56 155 181 145 24 38 21 159 114 129 234 216 64 62 220 79 109 27 163 28 174 47 142 122 240 182 216 45 121 198 149 87 108 131 125 240 16 91 244 209 215 133 131 28 85 109 209 26 156 203 19 221 175 168 201 201 120 201 169 139 78 116 87 153 86 255 211 107 42 0 240 159 133 142 138 34 36
	//authnticatorData
	//37 73 150 13 229 136 14 140 104 116 52 23 15 100 118 96 91 143 228 174 185 162 134 50 199 153 92 243 186 131 29 151 99 5 0 0 0 0
	//
	//95 0 123 34 116 121 112 101 34 58 34 119 101 98 97 117 116 104 110 46 103 101 116 34 44 34 99 104 97 108 108 101 110 103 101 34 58 34 89 87 70 104 34 44 34 111 114 105 103 105 110 34 58 34 104 116 116 112 58 47 47 108 111 99 97 108 104 111 115 116 58 56 48 48 49 34 44 34 99 114 111 115 115 79 114 105 103 105 110 34 58 102 97 108 115 101 125]
	signature := webauthnSignMsgBytes[3:67]
	fmt.Println("signature: ", common.Bytes2Hex(signature))
	//return
	pubKeyBytes := webauthnSignMsgBytes[68:132]
	fmt.Println("pubkey ", pubKeyBytes)
	//验证公钥
	
	//return
	authnticatorLenth := int(webauthnSignMsgBytes[132])
	fmt.Println("authnticatorLenth ", authnticatorLenth)
	//return
	authnticatorData := webauthnSignMsgBytes[133 : 133+authnticatorLenth]
	fmt.Println("authnticatorData ", authnticatorData)

	clientDataJsonData := webauthnSignMsgBytes[133+authnticatorLenth+2:]
	fmt.Println("clientDataJsonData ", clientDataJsonData)

	fmt.Println("json clientDataJsonData ", string(clientDataJsonData))
	// 解析JSON字符串到一个空接口类型
	var data map[string]interface{}

	err = json.Unmarshal(clientDataJsonData, &data)
	if err != nil {
		fmt.Println("unmarshal err :", err)
		return
	}

	// 验证challenge
	challengeBase64url := base64.URLEncoding.EncodeToString(challenge)
	fmt.Println("challengeBase64url", challengeBase64url)
	fmt.Println(data["challenge"])
	if challengeBase64url != data["challenge"] {
		return
	}
	clientDataJsonHash := sha256.Sum256(clientDataJsonData)
	signMsg := append(authnticatorData, clientDataJsonHash[:]...)
	hash := sha256.Sum256(signMsg)
	R := new(big.Int).SetBytes(signature[:32])
	S := new(big.Int).SetBytes(signature[32:])
	//fmt.Println("R ", signature[:32], len(signature[:32]))
	//fmt.Println("S ", signature[32:], len(signature[32:]))
	//return
	pubkey := new(ecdsa.PublicKey)
	pubkey.X = new(big.Int).SetBytes(pubKeyBytes[:32])
	pubkey.Y = new(big.Int).SetBytes(pubKeyBytes[32:])

	res, err = VerifyEcdsaP256Signature(hash[:], R, S, pubkey)
	return
}

func VerifyEcdsaP256Signature(hash []byte, R, S *big.Int, pubkey *ecdsa.PublicKey) (res bool, err error) {
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
