package bitcoin

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/btcutil/base58"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/dotbitHQ/das-lib/common"
)

func CreateDogecoinWallet() error {
	// doge net params
	// https://github.com/dogecoin/dogecoin/blob/master/src/chainparams.cpp#L167
	mainNetParams := chaincfg.MainNetParams
	mainNetParams.PubKeyHashAddrID = 0x1E //30
	mainNetParams.ScriptHashAddrID = 0x16 //33
	mainNetParams.PrivateKeyID = 0x9E     //158

	key, err := btcec.NewPrivateKey()
	if err != nil {
		return fmt.Errorf("NewPrivateKey err: %s", err.Error())
	}
	wif, err := btcutil.NewWIF(key, &mainNetParams, true)
	if err != nil {
		return fmt.Errorf("btcutil.NewWIF err: %s", err.Error())
	}
	addressPubKey, err := btcutil.NewAddressPubKey(wif.PrivKey.PubKey().SerializeUncompressed(), &mainNetParams)
	if err != nil {
		return fmt.Errorf("btcutil.NewAddressPubKey err: %s", err.Error())
	}
	fmt.Println("PubKey:", hex.EncodeToString(addressPubKey.AddressPubKeyHash().Hash160()[:]))
	fmt.Println("PubKey:", addressPubKey.EncodeAddress())
	fmt.Println("WIF:", wif.String())
	fmt.Println("PriKey:", hex.EncodeToString(key.Serialize()))
	return nil
}

func FormatAddressToPayload(addr string) (payload string, err error) {
	decode := base58.Decode(addr)
	payloadBys := decode[1 : len(decode)-4]
	//fmt.Println(hex.EncodeToString(decode))
	payload = hex.EncodeToString(payloadBys)
	//fmt.Println(payload)

	h := sha256.Sum256(decode[:len(decode)-4])
	h2 := sha256.Sum256(h[:])
	if bytes.Compare(h2[:4], decode[len(decode)-4:]) != 0 {
		err = fmt.Errorf("failed to checksum")
		return
	}
	return
}

func FormatPayloadToAddress(id common.DasAlgorithmId, payload string) (addr string, err error) {
	switch id {
	case common.DasAlgorithmIdDogecoin:
		payload = "1e" + payload
	default:
		err = fmt.Errorf("unknow DasAlgorithmId[%d]", id)
		return
	}
	bys := common.Hex2Bytes(payload)
	h := sha256.Sum256(bys)
	h2 := sha256.Sum256(h[:])
	bys = append(bys, h2[:4]...)
	addr = base58.Encode(bys)
	return
}
