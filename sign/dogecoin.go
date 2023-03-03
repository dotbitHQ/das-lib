package sign

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/dotbitHQ/das-lib/common"
	"github.com/ethereum/go-ethereum/crypto"
)

func magicHash(data []byte) ([]byte, error) {
	l := len(data)
	if l == 0 {
		return nil, fmt.Errorf("invalid raw data")
	}

	data = append(append([]byte(common.DogeMessageHeader), byte(l)), data...)
	//fmt.Println(hex.EncodeToString(data))

	h1 := sha256.New()
	h1.Write(data)
	data = h1.Sum(nil)
	//fmt.Println(hex.EncodeToString(data))

	h2 := sha256.New()
	h2.Write(data)
	data = h2.Sum(nil)
	//fmt.Println(hex.EncodeToString(data))
	return data, nil
}

func DogeSignature(data []byte, hexPrivateKey string, prex []byte) ([]byte, error) {
	bys, err := magicHash(data)
	if err != nil {
		return nil, fmt.Errorf("magicHash err: %s", err.Error())
	}
	fmt.Println(common.Bytes2Hex(bys))
	key, err := crypto.HexToECDSA(hexPrivateKey)
	if err != nil {
		return nil, fmt.Errorf("crypto.HexToECDSA err: %s", err.Error())
	}

	sig, err := crypto.Sign(bys, key)
	if err != nil {
		return nil, fmt.Errorf("crypto.Sign err: %s", err.Error())
	}

	if len(prex) > 0 {
		sig = append(prex, sig[:len(sig)-1]...)
	}
	return sig, nil
}

func VerifyDogeSignature(sig []byte, data []byte, payload string) (bool, error) {
	bys, err := magicHash(data)
	if err != nil {
		return false, fmt.Errorf("magicHash err: %s", err.Error())
	}

	if len(sig) != 65 { // sign check
		return false, fmt.Errorf("invalid param")
	}
	sigFormat := append(sig[1:], sig[:1]...)
	if sigFormat[64] >= 27 {
		sigFormat[64] -= 27
	}

	pub, err := crypto.Ecrecover(bys[:], sigFormat)
	if err != nil {
		return false, fmt.Errorf("crypto.Ecrecover err: %s", err.Error())
	}
	fmt.Println("pub:", hex.EncodeToString(pub))

	resPayload := hex.EncodeToString(btcutil.Hash160(pub))
	fmt.Println("VerifyDogeSignature:", resPayload)

	return resPayload == payload, nil
}
