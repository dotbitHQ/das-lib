package bitcoin

import (
	"fmt"
	"github.com/btcsuite/btcd/wire"
)

func (t *TxTool) NewBTCTxNewTx(uos []UnspentOutputs, addresses []string, values []int64, opReturn string) (*wire.MsgTx, error) {
	if len(uos) == 0 || (len(addresses) != len(values)) {
		return nil, fmt.Errorf("param is invalid:%v,%v,%v", uos, addresses, values)
	}

	return nil, nil
}

//func main1() {
//	// 创建比特币网络参数
//	params := &chaincfg.TestNet3Params
//
//	// 创建一个新的比特币交易
//	tx := wire.NewMsgTx(wire.TxVersion)
//
//	// 添加输入
//	prevTxHash, err := wire.NewShaHashFromStr("previous_transaction_hash")
//	if err != nil {
//		log.Fatal(err)
//	}
//	prevOut := wire.NewOutPoint(prevTxHash, 0) // 0 是输出索引
//	txIn := wire.NewTxIn(prevOut, nil, nil)
//	tx.AddTxIn(txIn)
//
//	// 添加输出
//	addressStr := "segwit_address"
//	segwitAddress, err := btcutil.DecodeAddress(addressStr, params)
//	if err != nil {
//		log.Fatal(err)
//	}
//	witnessProgram, err := txscript.PayToAddrScript(segwitAddress)
//	if err != nil {
//		log.Fatal(err)
//	}
//	txOut := wire.NewTxOut(1000000, witnessProgram) // 1000000 是发送的比特币数量（单位：聪）
//	tx.AddTxOut(txOut)
//
//	// 签名交易
//	privKeyWIF := "your_private_key_wif"
//	privKey, err := btcutil.DecodeWIF(privKeyWIF)
//	if err != nil {
//		log.Fatal(err)
//	}
//	sigScript, err := txscript.SignatureScript(tx, 0, witnessProgram, txscript.SigHashAll, privKey.PrivKey, true)
//	if err != nil {
//		log.Fatal(err)
//	}
//	txIn.SignatureScript = sigScript
//
//	// 打印交易信息
//	fmt.Printf("Transaction: %v\n", tx)
//}
