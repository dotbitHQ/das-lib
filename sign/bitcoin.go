package sign

import (
	"crypto/elliptic"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/dotbitHQ/das-lib/common"
	"github.com/ethereum/go-ethereum/crypto"
)

func magicHashBTC(data []byte) ([]byte, error) {
	l := len(data)
	if l == 0 {
		return nil, fmt.Errorf("invalid raw data")
	}

	data = append(append([]byte(common.BitcoinMessageHeader), byte(l)), data...)
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

func BitcoinSignature(data []byte, hexPrivateKey string, compress bool, segwitType SegwitType) ([]byte, error) {
	bys, err := magicHashBTC(data)
	if err != nil {
		return nil, fmt.Errorf("magicHashBTC err: %s", err.Error())
	}
	//fmt.Println("magicHash:", common.Bytes2Hex(bys))
	//log.Info("magicHashBTC:", common.Bytes2Hex(bys))
	key, err := crypto.HexToECDSA(hexPrivateKey)
	if err != nil {
		return nil, fmt.Errorf("crypto.HexToECDSA err: %s", err.Error())
	}

	sig, err := crypto.Sign(bys, key)
	if err != nil {
		return nil, fmt.Errorf("crypto.Sign err: %s", err.Error())
	}

	if compress {
		sig = append(sig, []byte{1}...)
	} else {
		sig = append(sig, []byte{0}...)
	}
	sig = append(sig, byte(segwitType))
	//log.Info("BitcoinSignature:", common.Bytes2Hex(sig))
	return sig, nil
}

func VerifyBitcoinSignature(sig []byte, data []byte, payload string) (bool, error) {
	bys, err := magicHashBTC(data)
	if err != nil {
		return false, fmt.Errorf("magicHashBTC err: %s", err.Error())
	}
	log.Info("magicHashBTC:", common.Bytes2Hex(bys))

	if len(sig) != 67 { // sign check
		return false, fmt.Errorf("invalid param")
	}
	if sig[64] >= 27 {
		sig[64] -= 27
	}

	compress := sig[65]
	segwitType := sig[66]
	switch segwitType {
	case byte(P2WPKH): // P2WPKH
		switch compress {
		case byte(0):
			return false, fmt.Errorf("invalid segwitType [%d] and compress [%d]", segwitType, compress)
			//pub, err := crypto.Ecrecover(bys, sig[:65])
			//if err != nil {
			//	return false, fmt.Errorf("crypto.Ecrecover err: %s", err.Error())
			//}
			//log.Info("publicKey:", common.Bytes2Hex(pub), len(pub))
			//
			//resPayload := hex.EncodeToString(btcutil.Hash160(pub))
			//log.Info("VerifyDogeSignature:", resPayload)
			//return resPayload == payload, nil
		case byte(1): // compressed
			sigToPub, err := crypto.SigToPub(bys, sig[:65])
			if err != nil {
				return false, fmt.Errorf("crypto.SigToPub err: %s", err.Error())
			}
			compressPublicKey := elliptic.MarshalCompressed(sigToPub.Curve, sigToPub.X, sigToPub.Y)

			//log.Info("compressPublicKey:", common.Bytes2Hex(compressPublicKey))
			resPayload := hex.EncodeToString(btcutil.Hash160(compressPublicKey))
			log.Info("VerifyDogeSignature:", resPayload)
			return resPayload == payload, nil
		default:
			return false, fmt.Errorf("unsupport compress[%d]", compress)
		}
	case byte(P2SH_P2WPKH): // P2SH_P2WPKH
		return false, fmt.Errorf("unsupport segwitType[%d]", segwitType)
	case byte(P2PKH): // P2PKH
		switch compress {
		case byte(0):
			pub, err := crypto.Ecrecover(bys, sig[:65])
			if err != nil {
				return false, fmt.Errorf("crypto.Ecrecover err: %s", err.Error())
			}
			log.Info("publicKey:", common.Bytes2Hex(pub), len(pub))

			resPayload := hex.EncodeToString(btcutil.Hash160(pub))
			//log.Info("VerifyDogeSignature:", resPayload)
			return resPayload == payload, nil
		case byte(1): // compressed
			sigToPub, err := crypto.SigToPub(bys, sig[:65])
			if err != nil {
				return false, fmt.Errorf("crypto.SigToPub err: %s", err.Error())
			}
			compressPublicKey := elliptic.MarshalCompressed(sigToPub.Curve, sigToPub.X, sigToPub.Y)

			//log.Info("compressPublicKey:", common.Bytes2Hex(compressPublicKey))
			resPayload := hex.EncodeToString(btcutil.Hash160(compressPublicKey))
			//log.Info("VerifyDogeSignature:", resPayload)
			return resPayload == payload, nil
		default:
			return false, fmt.Errorf("unsupport compress[%d]", compress)
		}
	default:
		return false, fmt.Errorf("unsupport segwitType[%d]", segwitType)
	}
}
