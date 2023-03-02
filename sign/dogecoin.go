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
	data = append([]byte(fmt.Sprintf(common.DogeMessageHeader, l)), data...)

	m := sha256.New()
	m.Write(data)
	data = m.Sum(nil)

	n := sha256.New()
	n.Write(data)
	data = n.Sum(nil)
	return data, nil
}

func DogeSignature(data []byte, hexPrivateKey string) ([]byte, error) {
	bys, err := magicHash(data)
	if err != nil {
		return nil, fmt.Errorf("magicHash err: %s", err.Error())
	}
	fmt.Println(hex.EncodeToString(bys))

	tmpHash := crypto.Keccak256(data)
	fmt.Println(hex.EncodeToString(tmpHash))

	//decodePrvKey, err := hex.DecodeString(hexPrivateKey)
	//if err != nil {
	//	return nil, fmt.Errorf("hex.DecodeString err: %s", err.Error())
	//}
	//prvKey, _ := btcec.PrivKeyFromBytes(decodePrvKey)
	//params := bitcoin.GetDogeMainNetParams()
	//wif, err := btcutil.NewWIF(prvKey, &params, true)
	//if err != nil {
	//	return nil, fmt.Errorf("btcutil.NewWIF err: %s", err.Error())
	//}

	//key, err := crypto.HexToECDSA(hexPrivateKey)
	//if err != nil {
	//	return nil, err
	//}
	//
	//return crypto.Sign(bys, key)
	return nil, nil
}

func VerifyDogeSignature(sign []byte, data []byte, payload string) (bool, error) {
	bys, err := magicHash(data)
	if err != nil {
		return false, fmt.Errorf("magicHash err: %s", err.Error())
	}

	if len(sign) != 65 { // sign check
		return false, fmt.Errorf("invalid param")
	}
	if sign[64] >= 27 {
		sign[64] -= 27
	}

	pub, err := crypto.Ecrecover(bys[:], sign)
	if err != nil {
		return false, err
	}

	return hex.EncodeToString(btcutil.Hash160(pub)) == payload, nil
}
