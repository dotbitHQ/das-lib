package bitcoin

import (
	"encoding/hex"
	"fmt"
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
)

// doge net params
func GetDogeMainNetParams() chaincfg.Params {
	// https://github.com/dogecoin/dogecoin/blob/master/src/chainparams.cpp#L167
	mainNetParams := chaincfg.MainNetParams
	mainNetParams.PubKeyHashAddrID = 0x1E //30
	mainNetParams.ScriptHashAddrID = 0x16 //33
	mainNetParams.PrivateKeyID = 0x9E     //158
	return mainNetParams
}

func CreateDogeWallet() error {
	mainNetParams := GetDogeMainNetParams()
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
