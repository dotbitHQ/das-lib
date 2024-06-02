package witness

import (
	"crypto/elliptic"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/btcsuite/btcd/btcec/v2/schnorr"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/txscript"
	"github.com/dotbitHQ/das-lib/bitcoin"
	"github.com/dotbitHQ/das-lib/common"
	"github.com/dotbitHQ/das-lib/molecule"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/nervosnetwork/ckb-sdk-go/crypto/blake2b"
)

type ReverseSmtRecordAction string
type ReverseSmtRecordVersion uint32

const (
	ReverseSmtRecordVersion1 ReverseSmtRecordVersion = 1

	ReverseSmtRecordActionUpdate ReverseSmtRecordAction = "update"
	ReverseSmtRecordActionRemove ReverseSmtRecordAction = "remove"
)

type ReverseSmtRecord struct {
	Version     ReverseSmtRecordVersion
	Action      ReverseSmtRecordAction
	Signature   []byte
	SignType    uint8
	Address     []byte
	Proof       []byte
	PrevNonce   uint32 `witness:",omitempty"`
	PrevAccount string
	NextRoot    []byte
	NextAccount string
}

func (r *ReverseSmtRecord) GetP2SHP2WPKH(netType common.DasNetType) (string, error) {
	if common.DasAlgorithmId(r.SignType) == common.DasAlgorithmIdBitcoin {
		pkHash := r.Address
		net := bitcoin.GetBTCMainNetParams()
		if netType != common.DasNetTypeMainNet {
			net = bitcoin.GetBTCTestNetParams()
		}

		addressWPH, err := btcutil.NewAddressWitnessPubKeyHash(pkHash, &net)
		if err != nil {
			return "", fmt.Errorf("btcutil.NewAddressWitnessPubKeyHash err: %s", err.Error())
		}
		pkScript, err := txscript.PayToAddrScript(addressWPH)
		if err != nil {
			return "", fmt.Errorf("txscript.PayToAddrScript err: %s", err.Error())
		}
		scriptAddr, err := btcutil.NewAddressScriptHash(pkScript, &net)
		if err != nil {
			return "", fmt.Errorf("btcutil.NewAddressScriptHash err: %s", err.Error())
		}
		addr := scriptAddr.EncodeAddress()
		log.Info("GetP2SHP2WPKH:", common.Bytes2Hex(r.Address), addr)
		return addr, nil
	}
	return "", nil
}

func (r *ReverseSmtRecord) GetP2TR(netType common.DasNetType) (string, error) {
	if common.DasAlgorithmId(r.SignType) == common.DasAlgorithmIdBitcoin {
		net := bitcoin.GetBTCMainNetParams()
		if netType != common.DasNetTypeMainNet {
			net = bitcoin.GetBTCTestNetParams()
		}

		data := make([]byte, 0)
		data = append(data, molecule.GoU32ToBytes(r.PrevNonce+1)...)
		data = append(data, []byte(r.NextAccount)...)
		dataBys, _ := blake2b.Blake256(data)
		signMsg := common.DotBitPrefix + hex.EncodeToString(dataBys)
		log.Info("GetP2TR signMsg:", signMsg)

		bys, err := magicHashBTC([]byte(signMsg))
		if err != nil {
			return "", fmt.Errorf("magicHashBTC err: %s", err.Error())
		}
		log.Info("magicHashBTC:", common.Bytes2Hex(bys))
		log.Info("r.Signature:", common.Bytes2Hex(r.Signature))

		sig := r.Signature
		if sig[64] >= 27 {
			sig[64] -= 27
		}
		fmt.Println(common.Bytes2Hex(bys))
		fmt.Println(common.Bytes2Hex(sig[:65]))
		sigToPub, err := crypto.SigToPub(bys, sig[:65])
		if err != nil {
			return "", fmt.Errorf("crypto.SigToPub err: %s", err.Error())
		}
		compressPublicKey := elliptic.MarshalCompressed(sigToPub.Curve, sigToPub.X, sigToPub.Y)
		log.Info("compressPublicKey:", common.Bytes2Hex(compressPublicKey))
		resPayload := hex.EncodeToString(btcutil.Hash160(compressPublicKey))
		log.Info("VerifyBitcoinSignature:", resPayload)

		//
		addressPubKey, err := btcutil.NewAddressPubKey(compressPublicKey, &net)
		if err != nil {
			return "", fmt.Errorf("btcutil.NewAddressPubKey err: %s", err.Error())
		}
		tapKey := txscript.ComputeTaprootKeyNoScript(addressPubKey.PubKey())
		addrTR, err := btcutil.NewAddressTaproot(
			schnorr.SerializePubKey(tapKey),
			&net,
		)
		if err != nil {
			return "", fmt.Errorf("btcutil.NewAddressTaproot err: %s", err.Error())
		}
		addr := addrTR.EncodeAddress()
		log.Info("GetP2TR:", common.Bytes2Hex(r.Address), addr)
		return addr, nil
	}
	return "", nil
}

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
