package sign

import (
	"bytes"
	"crypto/elliptic"
	"crypto/sha256"
	"encoding/hex"
	"errors"
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

func DogeSignature(data []byte, hexPrivateKey string, compress bool, segwitTypes SegwitType) ([]byte, error) {
	bys, err := magicHash(data)
	if err != nil {
		return nil, fmt.Errorf("magicHash err: %s", err.Error())
	}
	fmt.Println("magicHash:", common.Bytes2Hex(bys))
	key, err := crypto.HexToECDSA(hexPrivateKey)
	if err != nil {
		return nil, fmt.Errorf("crypto.HexToECDSA err: %s", err.Error())
	}

	sig, err := crypto.Sign(bys, key)
	if err != nil {
		return nil, fmt.Errorf("crypto.Sign err: %s", err.Error())
	}
	fmt.Println("DogeSignature:", common.Bytes2Hex(sig))
	if compress {
		sig = append(sig, []byte{1}...)
	} else {
		sig = append(sig, []byte{0}...)
	}

	switch segwitTypes {
	case P2WPKH:
		sig = append(sig, []byte{0}...)
	case P2SH_P2WPKH:
		sig = append(sig, []byte{1}...)
	default:
		sig = append(sig, []byte{0}...)
	}

	return sig, nil
}

func VerifyDogeSignature(sig []byte, data []byte, payload string) (bool, error) {
	bys, err := magicHash(data)
	if err != nil {
		return false, fmt.Errorf("magicHash err: %s", err.Error())
	}
	fmt.Println("magicHash:", common.Bytes2Hex(bys))

	if len(sig) != 67 { // sign check
		return false, fmt.Errorf("invalid param")
	}
	if sig[64] >= 27 {
		sig[64] -= 27
	}

	if bytes.Compare(sig[65:66], []byte{1}) == 0 { // compressed
		//pub, err := crypto.Ecrecover(bys, sig[:65])
		//if err != nil {
		//	return false, fmt.Errorf("crypto.Ecrecover err: %s", err.Error())
		//}
		//fmt.Println("publicKey:", hex.EncodeToString(pub), len(pub))

		sigToPub, err := crypto.SigToPub(bys, sig[:65])
		if err != nil {
			return false, fmt.Errorf("crypto.SigToPub err: %s", err.Error())
		}
		compressPublicKey := elliptic.MarshalCompressed(sigToPub.Curve, sigToPub.X, sigToPub.Y)

		//pub, err := crypto.Ecrecover(bys, sig[:65])
		//if err != nil {
		//	return false, fmt.Errorf("crypto.Ecrecover err: %s", err.Error())
		//}
		//fmt.Println("publicKey:", hex.EncodeToString(pub))
		//compressPublicKey := append([]byte{2}, pub[1:33]...)

		fmt.Println("compressPublicKey:", hex.EncodeToString(compressPublicKey))
		resPayload := hex.EncodeToString(btcutil.Hash160(compressPublicKey))
		fmt.Println("VerifyDogeSignature:", resPayload)
		return resPayload == payload, nil
	} else {
		pub, err := crypto.Ecrecover(bys, sig[:65])
		if err != nil {
			return false, fmt.Errorf("crypto.Ecrecover err: %s", err.Error())
		}
		fmt.Println("publicKey:", hex.EncodeToString(pub), len(pub))

		resPayload := hex.EncodeToString(btcutil.Hash160(pub))
		fmt.Println("VerifyDogeSignature:", resPayload)

		return resPayload == payload, nil
	}
}

type SegwitType int

const (
	P2WPKH SegwitType = iota
	P2SH_P2WPKH
)

type SignatureInfo struct {
	Compressed bool
	SegwitType *SegwitType
	Recovery   int
	Signature  []byte
}

func DecodeSignature(buffer []byte) (*SignatureInfo, error) {
	if len(buffer) != 65 {
		return nil, errors.New("invalid signature length")
	}

	flagByte := int(buffer[0]) - 27
	if flagByte > 15 || flagByte < 0 {
		return nil, errors.New("invalid signature parameter")
	}

	var segwitType *SegwitType = nil
	if (flagByte & 8) != 0 {
		if (flagByte & 4) != 0 {
			segwitType = new(SegwitType)
			*segwitType = P2WPKH
		} else {
			segwitType = new(SegwitType)
			*segwitType = P2SH_P2WPKH
		}
	}

	return &SignatureInfo{
		Compressed: (flagByte & 12) != 0,
		SegwitType: segwitType,
		Recovery:   flagByte & 3,
		Signature:  buffer[1:],
	}, nil
}

func (s *SignatureInfo) ToSig() []byte {
	res := append(s.Signature, byte(s.Recovery))
	if s.Compressed {
		res = append(res, byte(1))
	} else {
		res = append(res, byte(0))
	}
	if s.SegwitType == nil {
		res = append(res, byte(0))
	} else {
		switch *s.SegwitType {
		case P2WPKH:
			res = append(res, byte(0))
		case P2SH_P2WPKH:
			res = append(res, byte(1))
		default:
			res = append(res, byte(0))
		}
	}
	return res
}
