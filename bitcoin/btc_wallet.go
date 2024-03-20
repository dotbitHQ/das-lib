package bitcoin

import (
	"encoding/hex"
	"fmt"
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
)

type BtcAddressType string

const (
	BtcAddressTypeP2PKH      BtcAddressType = "P2PKH"
	BtcAddressTypeP2WPKH     BtcAddressType = "P2WPKH"
	BtcAddressTypeP2SHP2WPKH BtcAddressType = "P2SH-P2WPKH"
)

// btc net params
func GetBTCMainNetParams() chaincfg.Params {
	//https: //github.com/bitcoin/bitcoin/blob/3d216baf91ca754e46e89788205513a956ec6d0a/src/kernel/chainparams.cpp#L145
	mainNetParams := chaincfg.MainNetParams
	mainNetParams.PubKeyHashAddrID = 0x00
	mainNetParams.ScriptHashAddrID = 0x05
	mainNetParams.PrivateKeyID = 0x80
	return mainNetParams
}

func CreateBTCWallet(addrType BtcAddressType, compress bool) error {
	mainNetParams := GetBTCMainNetParams()
	switch addrType {
	case BtcAddressTypeP2PKH:
		key, err := btcec.NewPrivateKey()
		if err != nil {
			return fmt.Errorf("NewPrivateKey err: %s", err.Error())
		}
		wif, err := btcutil.NewWIF(key, &mainNetParams, compress)
		if err != nil {
			return fmt.Errorf("btcutil.NewWIF err: %s", err.Error())
		}
		addressPubKey, err := btcutil.NewAddressPubKey(wif.SerializePubKey(), &mainNetParams)
		if err != nil {
			return fmt.Errorf("btcutil.NewAddressPubKey err: %s", err.Error())
		}

		fmt.Println("WIF:", wif.String())
		fmt.Println("PubKey:", addressPubKey.EncodeAddress())
		fmt.Println("PubHash", hex.EncodeToString(addressPubKey.AddressPubKeyHash().Hash160()[:]))
		fmt.Println("PriKey:", hex.EncodeToString(key.Serialize()))
	case BtcAddressTypeP2WPKH:
		key, err := btcec.NewPrivateKey()
		if err != nil {
			return fmt.Errorf("NewPrivateKey err: %s", err.Error())
		}
		wif, err := btcutil.NewWIF(key, &mainNetParams, compress)
		if err != nil {
			return fmt.Errorf("btcutil.NewWIF err: %s", err.Error())
		}
		addressPubKey, err := btcutil.NewAddressPubKey(wif.SerializePubKey(), &mainNetParams)
		if err != nil {
			return fmt.Errorf("btcutil.NewAddressPubKey err: %s", err.Error())
		}
		pkHash := addressPubKey.AddressPubKeyHash().Hash160()[:]

		addressWPH, err := btcutil.NewAddressWitnessPubKeyHash(pkHash, &mainNetParams)
		if err != nil {
			return fmt.Errorf("NewAddressWitnessPubKeyHash err: %s", err.Error())
		}

		fmt.Println("WIF:", wif.String())
		fmt.Println("PubKey:", addressWPH.EncodeAddress())
		fmt.Println("PubHash", hex.EncodeToString(addressPubKey.AddressPubKeyHash().Hash160()[:]))
		fmt.Println("PriKey:", hex.EncodeToString(key.Serialize()))
	case BtcAddressTypeP2SHP2WPKH:
	}
	return nil
}
