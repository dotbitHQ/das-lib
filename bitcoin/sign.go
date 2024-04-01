package bitcoin

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
)

type reqParam struct {
	Errno  int         `json:"errno"`
	Errmsg interface{} `json:"errmsg"`
	Data   interface{} `json:"data"`
}

type addressInfo struct {
	Address string `json:"address"`
	Value   int64  `json:"value"`
}

type RemoteSignMethod = string

const (
	RemoteSignMethodDogeTx RemoteSignMethod = "wallet_dogSignMsg"
	RemoteSignMethodBtcTx  RemoteSignMethod = "wallet_btcSignMsg"
)

func (t *TxTool) RemoteSignTx(method RemoteSignMethod, tx *wire.MsgTx, uos []UnspentOutputs) (*wire.MsgTx, error) {
	if tx == nil || len(uos) == 0 {
		return nil, fmt.Errorf("tx is nil")
	}
	if len(tx.TxIn) != len(uos) {
		return nil, fmt.Errorf("len of txin != len of uts")
	}
	reply := reqParam{}
	var data string
	reply.Data = &data

	rawTx, err := txToString(tx)
	if err != nil {
		return nil, fmt.Errorf("txToString err: %s", err.Error())
	}

	param := struct {
		Addresses []addressInfo `json:"addrs"`
		Tx        string        `json:"tx"`
	}{Tx: rawTx}
	for _, unspent := range uos {
		param.Addresses = append(param.Addresses, addressInfo{Address: unspent.Address, Value: unspent.Value})
	}

	if err := t.RemoteSignClient.CallContext(t.Ctx, &reply, method, param); err != nil {
		return nil, fmt.Errorf("client.CallContext err: %s", err.Error())
	}
	if reply.Errno == 0 {
		signTx, err := stringToTx(data)
		if err != nil {
			return nil, fmt.Errorf("stringToTx err: %s", err.Error())
		}
		return signTx, nil
	} else {
		return nil, fmt.Errorf("client.CallContext err: %s [%d]", reply.Errmsg, reply.Errno)
	}
}

func txToString(tx *wire.MsgTx) (string, error) {
	buf := bytes.NewBuffer(make([]byte, 0, tx.SerializeSize()))
	if err := tx.SerializeNoWitness(buf); err != nil {
		return "", fmt.Errorf("SerializeNoWitness err: %s", err)
	}
	return hex.EncodeToString(buf.Bytes()), nil
}

func stringToTx(hexTx string) (*wire.MsgTx, error) {
	transaction, err := hex.DecodeString(hexTx)
	if err != nil {
		return nil, fmt.Errorf("DecodeString err: %s", err.Error())
	}
	var msgTx wire.MsgTx
	if err := msgTx.DeserializeNoWitness(bytes.NewReader(transaction)); err != nil {
		return nil, fmt.Errorf("DeserializeNoWitness err: %s", err.Error())
	}
	return &msgTx, nil
}

func (t *TxTool) LocalSignTx(tx *wire.MsgTx, uos []UnspentOutputs) (string, error) {
	if tx == nil || len(uos) == 0 {
		return "", fmt.Errorf("tx is nil")
	}
	if len(tx.TxIn) != len(uos) {
		return "", fmt.Errorf("len of txin != len of uts")
	}

	for i := 0; i < len(uos); i++ {
		item := uos[i]
		if item.Private == "" {
			return "", fmt.Errorf("PrivateKey is nil")
		}
		pkScript, privateKey, compress, err := HexPrivateKeyToScript(item.Address, t.Params, item.Private)
		if err != nil {
			return "", fmt.Errorf("HexPrivateKeyToScript err: %s", err.Error())
		}
		sig, err := txscript.SignatureScript(tx, i, pkScript, txscript.SigHashAll, privateKey, compress)
		if err != nil {
			return "", fmt.Errorf("SignatureScript err: %s", err.Error())
		}
		tx.TxIn[i].SignatureScript = sig
	}
	buf := bytes.NewBuffer(make([]byte, 0, tx.SerializeSizeStripped()))
	_ = tx.SerializeNoWitness(buf)
	return hex.EncodeToString(buf.Bytes()), nil
}

func HexPrivateKeyToScript(addr string, params chaincfg.Params, privateKeyHex string) (pkScript []byte, privateKey *btcec.PrivateKey, compress bool, e error) {
	// pkScriptBytes
	scriptAddr, err := btcutil.DecodeAddress(addr, &params)
	if err != nil {
		e = fmt.Errorf("btcutil.DecodeAddress err: %s", err.Error())
		return
	}
	pkScript, err = txscript.PayToAddrScript(scriptAddr)
	if err != nil {
		e = fmt.Errorf("txscript.PayToAddrScript err: %s", err.Error())
		return
	}

	// privateKey
	privateKeyBys, err := hex.DecodeString(privateKeyHex)
	if err != nil {
		e = fmt.Errorf("hex.DecodeString err: %s", err.Error())
		return
	}
	privateKeyTmp, publicKey := btcec.PrivKeyFromBytes(privateKeyBys)
	privateKey = privateKeyTmp
	compressPubKeyHash := hex.EncodeToString(btcutil.Hash160(publicKey.SerializeCompressed()))
	pubKeyHash := hex.EncodeToString(btcutil.Hash160(publicKey.SerializeUncompressed()))

	encodeAddress := hex.EncodeToString(scriptAddr.ScriptAddress())
	//log.Info("HexPrivateKeyToScript:", encodeAddress, compressPubKeyHash, pubKeyHash)
	if encodeAddress == compressPubKeyHash {
		compress = true
	} else if encodeAddress == pubKeyHash {
		compress = false
	} else {
		e = fmt.Errorf("compress check err: ")
		return
	}

	return
}
