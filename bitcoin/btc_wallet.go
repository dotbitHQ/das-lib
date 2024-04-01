package bitcoin

import (
	"encoding/hex"
	"fmt"
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"strings"
)

type BtcAddressType string

const (
	BtcAddressTypeP2PKH      BtcAddressType = "P2PKH"
	BtcAddressTypeP2SH       BtcAddressType = "P2SH"
	BtcAddressTypeP2WPKH     BtcAddressType = "P2WPKH"
	BtcAddressTypeP2SHP2WPKH BtcAddressType = "P2SH-P2WPKH"
	BtcAddressTypeP2TR       BtcAddressType = "P2TR"
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
	case BtcAddressTypeP2SH:
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

		pkScript, err := txscript.PayToAddrScript(addressPubKey.AddressPubKeyHash())
		if err != nil {
			return fmt.Errorf("txscript.PayToAddrScript err: %s", err.Error())
		}

		scriptAddr, err := btcutil.NewAddressScriptHash(pkScript, &mainNetParams)
		if err != nil {
			return fmt.Errorf("btcutil.NewAddressScriptHash err: %s", err.Error())
		}

		fmt.Println("WIF:", wif.String())
		//fmt.Println("PubKey:", addressPubKey.EncodeAddress())
		fmt.Println("ScriptAddr:", scriptAddr.EncodeAddress())
		fmt.Println("PubHash", hex.EncodeToString(addressPubKey.AddressPubKeyHash().Hash160()[:]))
		fmt.Println("pkScript:", hex.EncodeToString(pkScript))
		fmt.Println("pkScriptHash:", hex.EncodeToString(btcutil.Hash160(pkScript)))
		fmt.Println("PriKey:", hex.EncodeToString(key.Serialize()))
	case BtcAddressTypeP2WPKH:
		if compress == false {
			return fmt.Errorf("compress must be true")
		}
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
		fmt.Println("PubKey:", addressWPH.EncodeAddress(), len(addressWPH.EncodeAddress()))
		fmt.Println("PubHash", hex.EncodeToString(addressPubKey.AddressPubKeyHash().Hash160()[:]))
		fmt.Println("PriKey:", hex.EncodeToString(key.Serialize()))
	case BtcAddressTypeP2SHP2WPKH:
		if compress == false {
			return fmt.Errorf("compress must be true")
		}
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

		//
		pkScript, err := txscript.PayToAddrScript(addressWPH)
		if err != nil {
			return fmt.Errorf("txscript.PayToAddrScript err: %s", err.Error())
		}
		scriptAddr, err := btcutil.NewAddressScriptHash(pkScript, &mainNetParams)
		if err != nil {
			return fmt.Errorf("btcutil.NewAddressScriptHash err: %s", err.Error())
		}

		fmt.Println("WIF:", wif.String())
		//fmt.Println("PubKey:", addressPubKey.EncodeAddress())
		fmt.Println("ScriptAddr:", scriptAddr.EncodeAddress())
		fmt.Println("PubHash", hex.EncodeToString(addressPubKey.AddressPubKeyHash().Hash160()[:]))
		fmt.Println("pkScript:", hex.EncodeToString(pkScript))
		fmt.Println("pkScriptHash:", hex.EncodeToString(btcutil.Hash160(pkScript)), btcutil.Hash160(pkScript))
		fmt.Println("PriKey:", hex.EncodeToString(key.Serialize()))
	case BtcAddressTypeP2TR:
		return fmt.Errorf("unsupport P2TR")
	}
	return nil
}

func FormatBTCAddr(addr string) (BtcAddressType, string, error) {
	addrType := BtcAddressType("")
	if strings.HasPrefix(addr, "bc1q") {
		if len(addr) != 42 {
			return "", "", fmt.Errorf("unspport address [%s]", addr)
		}
		addrType = BtcAddressTypeP2WPKH
	} else if strings.HasPrefix(addr, "1") {
		addrType = BtcAddressTypeP2PKH
	} else {
		return "", "", fmt.Errorf("unspport address [%s]", addr)
	}

	netParams := GetBTCMainNetParams()
	addrDecode, err := btcutil.DecodeAddress(addr, &netParams)
	if err != nil {
		return "", "", fmt.Errorf("btcutil.DecodeAddress [%s] err: %s", addr, err.Error())
	}
	return addrType, hex.EncodeToString(addrDecode.ScriptAddress()), nil
}
