package bitcoin

import (
	"encoding/hex"
	"fmt"
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/schnorr"
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
	mainNetParams.Bech32HRPSegwit = "bc"
	return mainNetParams
}

func GetBTCTestNetParams() chaincfg.Params {
	//https: //github.com/bitcoin/bitcoin/blob/3d216baf91ca754e46e89788205513a956ec6d0a/src/kernel/chainparams.cpp#L145
	mainNetParams := chaincfg.MainNetParams
	mainNetParams.PubKeyHashAddrID = 0x6f
	mainNetParams.ScriptHashAddrID = 0xc4
	mainNetParams.PrivateKeyID = 0xef
	mainNetParams.Bech32HRPSegwit = "tb"
	return mainNetParams
}

func CreateBTCWallet(netParams chaincfg.Params, addrType BtcAddressType, compress bool) error {
	switch addrType {
	case BtcAddressTypeP2PKH:
		key, err := btcec.NewPrivateKey()
		if err != nil {
			return fmt.Errorf("NewPrivateKey err: %s", err.Error())
		}
		wif, err := btcutil.NewWIF(key, &netParams, compress)
		if err != nil {
			return fmt.Errorf("btcutil.NewWIF err: %s", err.Error())
		}
		addressPubKey, err := btcutil.NewAddressPubKey(wif.SerializePubKey(), &netParams)
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
		wif, err := btcutil.NewWIF(key, &netParams, compress)
		if err != nil {
			return fmt.Errorf("btcutil.NewWIF err: %s", err.Error())
		}
		addressPubKey, err := btcutil.NewAddressPubKey(wif.SerializePubKey(), &netParams)
		if err != nil {
			return fmt.Errorf("btcutil.NewAddressPubKey err: %s", err.Error())
		}

		pkScript, err := txscript.PayToAddrScript(addressPubKey.AddressPubKeyHash())
		if err != nil {
			return fmt.Errorf("txscript.PayToAddrScript err: %s", err.Error())
		}

		scriptAddr, err := btcutil.NewAddressScriptHash(pkScript, &netParams)
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
		wif, err := btcutil.NewWIF(key, &netParams, compress)
		if err != nil {
			return fmt.Errorf("btcutil.NewWIF err: %s", err.Error())
		}
		addressPubKey, err := btcutil.NewAddressPubKey(wif.SerializePubKey(), &netParams)
		if err != nil {
			return fmt.Errorf("btcutil.NewAddressPubKey err: %s", err.Error())
		}
		pkHash := addressPubKey.AddressPubKeyHash().Hash160()[:]

		addressWPH, err := btcutil.NewAddressWitnessPubKeyHash(pkHash, &netParams)
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
		wif, err := btcutil.NewWIF(key, &netParams, compress)
		if err != nil {
			return fmt.Errorf("btcutil.NewWIF err: %s", err.Error())
		}
		addressPubKey, err := btcutil.NewAddressPubKey(wif.SerializePubKey(), &netParams)
		if err != nil {
			return fmt.Errorf("btcutil.NewAddressPubKey err: %s", err.Error())
		}
		pkHash := addressPubKey.AddressPubKeyHash().Hash160()[:]
		addressWPH, err := btcutil.NewAddressWitnessPubKeyHash(pkHash, &netParams)
		if err != nil {
			return fmt.Errorf("NewAddressWitnessPubKeyHash err: %s", err.Error())
		}

		//

		pkScript, err := txscript.PayToAddrScript(addressWPH)
		if err != nil {
			return fmt.Errorf("txscript.PayToAddrScript err: %s", err.Error())
		}
		scriptAddr, err := btcutil.NewAddressScriptHash(pkScript, &netParams)
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
		if compress == false {
			return fmt.Errorf("compress must be true")
		}
		key, err := btcec.NewPrivateKey()
		if err != nil {
			return fmt.Errorf("NewPrivateKey err: %s", err.Error())
		}
		wif, err := btcutil.NewWIF(key, &netParams, compress)
		if err != nil {
			return fmt.Errorf("btcutil.NewWIF err: %s", err.Error())
		}
		addressPubKey, err := btcutil.NewAddressPubKey(wif.SerializePubKey(), &netParams)
		if err != nil {
			return fmt.Errorf("btcutil.NewAddressPubKey err: %s", err.Error())
		}
		//pkHash := addressPubKey.AddressPubKeyHash().Hash160()[:]

		//
		tapKey := txscript.ComputeTaprootKeyNoScript(addressPubKey.PubKey())
		addrTR, err := btcutil.NewAddressTaproot(
			schnorr.SerializePubKey(tapKey),
			&netParams,
		)
		if err != nil {
			return fmt.Errorf("btcutil.NewAddressTaproot err: %s", err.Error())
		}
		fmt.Println("WIF:", wif.String())
		//fmt.Println("PubKey:", addressPubKey.EncodeAddress())
		fmt.Println("TRAddr:", addrTR.EncodeAddress())
		fmt.Println("PubHash", hex.EncodeToString(addressPubKey.AddressPubKeyHash().Hash160()[:]))
		fmt.Println("PriKey:", hex.EncodeToString(key.Serialize()))
	}
	return nil
}

func FormatBTCAddr(addr string) (chaincfg.Params, BtcAddressType, string, error) {
	addrType := BtcAddressType("")
	netParams := GetBTCMainNetParams()
	if strings.HasPrefix(addr, "tb1q") || strings.HasPrefix(addr, "m") {
		netParams = GetBTCTestNetParams()
	}
	if strings.HasPrefix(addr, "bc1q") || strings.HasPrefix(addr, "tb1q") {
		if len(addr) != 42 {
			return netParams, "", "", fmt.Errorf("unspport address [%s]", addr)
		}
		addrType = BtcAddressTypeP2WPKH
	} else if strings.HasPrefix(addr, "1") || strings.HasPrefix(addr, "m") {
		addrType = BtcAddressTypeP2PKH
	} else {
		return netParams, "", "", fmt.Errorf("unspport address [%s]", addr)
	}

	addrDecode, err := btcutil.DecodeAddress(addr, &netParams)
	if err != nil {
		return netParams, "", "", fmt.Errorf("btcutil.DecodeAddress [%s] err: %s", addr, err.Error())
	}
	return netParams, addrType, hex.EncodeToString(addrDecode.ScriptAddress()), nil
}
