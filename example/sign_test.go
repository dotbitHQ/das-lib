package example

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/dotbitHQ/das-lib/common"
	"github.com/dotbitHQ/das-lib/sign"
	"github.com/dotbitHQ/das-lib/txbuilder"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
	"github.com/nervosnetwork/ckb-sdk-go/transaction"
	"github.com/nervosnetwork/ckb-sdk-go/types"
	"strings"
	"testing"
)

func TestTronSignature(t *testing.T) {
	//res:="0xd5556e62653347b6b95d3d5c5c00439d7bae8f22708483a1d970d22be1ca40b43414733532aab98ee25bf68cbf215143778e835f0a4bd70942899d7fe564107f1c"
	signType := true
	data := common.Hex2Bytes("0xd3b4b3ed69dcfbbdc593bedb06b083e417792ddb0aef6fae293071f42b2d824804")
	privateKey := ""
	address := "TQoLh9evwUmZKxpD1uhFttsZk3EBs8BksV"
	signature, err := sign.TronSignature(signType, data, privateKey)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(common.Bytes2Hex(signature))

	fmt.Println(sign.TronVerifySignature(signType, signature, data, address))
}

func TestEthSignature(t *testing.T) {
	data := common.Hex2Bytes("0x15f92d66997823cbc225c806e2160cada949765eee0a50c467e439d53e225254")
	privateKey := ""
	address := "0xdD3b3D0F3FA9546a5616d0200b83f784a5220ae8"
	signature, err := sign.EthSignature(data, privateKey)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(common.Bytes2Hex(signature))
	ok, err := sign.VerifyEthSignature(signature, data, address)
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Fatal(fmt.Errorf("verify failed"))
	}
}
func TestPersonalSignature(t *testing.T) {
	data0 := hex.EncodeToString([]byte("0xADD EMAIL - 1639644121"))
	fmt.Println(data0)
	data1 := common.Hex2Bytes(data0)
	privateKey := ""
	address := ""
	signature, err := sign.PersonalSignature(data1, privateKey)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(common.Bytes2Hex(signature))
	ok, err := sign.VerifyPersonalSignature(signature, data1, address)
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Fatal(fmt.Errorf("verify failed"))
	}
}

func TestEthSignature712(t *testing.T) {
	//0x5f410f286333decd069582eb991d965e123f9c3bef2079198bf7145cc1ead0ac005489e0e317317fc4fc13809e838c3bf99492f3f078f088a5864f886bb3a7ef0183df5ddcefaecac331dde3d78ca45b0a664bd65a292d87a8b7530e79c341f8690000000000000005
	digest := "0x6f416b786af70e4ce946c51ff2caffb26f24116bb433fb6ec7fe00a992fabc6d"
	mmJson := `{"types":{"EIP712Domain":[{"name":"chainId","type":"uint256"},{"name":"name","type":"string"},{"name":"verifyingContract","type":"address"},{"name":"version","type":"string"}],"Action":[{"name":"action","type":"string"},{"name":"params","type":"string"}],"Cell":[{"name":"capacity","type":"string"},{"name":"lock","type":"string"},{"name":"type","type":"string"},{"name":"data","type":"string"},{"name":"extraData","type":"string"}],"Transaction":[{"name":"DAS_MESSAGE","type":"string"},{"name":"inputsCapacity","type":"string"},{"name":"outputsCapacity","type":"string"},{"name":"fee","type":"string"},{"name":"action","type":"Action"},{"name":"inputs","type":"Cell[]"},{"name":"outputs","type":"Cell[]"},{"name":"digest","type":"bytes32"}]},"primaryType":"Transaction","domain":{"chainId":5,"name":"da.systems","verifyingContract":"0x0000000000000000000000000000000020210722","version":"1"},"message":{"DAS_MESSAGE":"EDIT RECORDS OF ACCOUNT 0001.bit","inputsCapacity":"214.9989 CKB","outputsCapacity":"214.9988 CKB","fee":"0.0001 CKB","digest":"","action":{"action":"edit_records","params":"0x01"},"inputs":[{"capacity":"214.9989 CKB","lock":"das-lock,0x01,0x05c9f53b1d85356b60453f867610888d89a0b667...","type":"account-cell-type,0x01,0x","data":"{ account: 0001.bit, expired_at: 1916807174 }","extraData":"{ status: 0, records_hash: 0xa7a206f7b378a0181909a98bf9fe5a167f72cdd0edcb35749a79f49c0ecf3c61 }"}],"outputs":[{"capacity":"214.9988 CKB","lock":"das-lock,0x01,0x05c9f53b1d85356b60453f867610888d89a0b667...","type":"account-cell-type,0x01,0x","data":"{ account: 0001.bit, expired_at: 1916807174 }","extraData":"{ status: 0, records_hash: 0x55478d76900611eb079b22088081124ed6c8bae21a05dd1a0d197efcc7c114ce }"}]}}`
	privateKey := ""
	address := "0x15a33588908cF8Edb27D1AbE3852Bf287Abd3891"

	chainId := 5
	var obj3 apitypes.TypedData
	oldChainId := fmt.Sprintf("chainId\":%d", chainId)
	newChainId := fmt.Sprintf("chainId\":\"%d\"", chainId)
	mmJson = strings.ReplaceAll(mmJson, oldChainId, newChainId)
	oldDigest := "\"digest\":\"\""
	newDigest := fmt.Sprintf("\"digest\":\"%s\"", digest)
	mmJson = strings.ReplaceAll(mmJson, oldDigest, newDigest)

	_ = json.Unmarshal([]byte(mmJson), &obj3)
	var mmHash, signature []byte
	mmHash, signature, err := sign.EIP712Signature(obj3, privateKey)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println("EIP712Signature mmHash:", common.Bytes2Hex(mmHash))
	fmt.Println("EIP712Signature signature:", common.Bytes2Hex(signature))

	signData := append(signature, mmHash...)
	hexChainId := fmt.Sprintf("%x", chainId)
	chainIdData := common.Hex2Bytes(fmt.Sprintf("%016s", hexChainId))
	signData = append(signData, chainIdData...)
	fmt.Println("signData:", common.Bytes2Hex(signData))

	fmt.Println(sign.VerifyEIP712Signature(obj3, signature, address))
	//0x5f410f286333decd069582eb991d965e123f9c3bef2079198bf7145cc1ead0ac005489e0e317317fc4fc13809e838c3bf99492f3f078f088a5864f886bb3a7ef0183df5ddcefaecac331dde3d78ca45b0a664bd65a292d87a8b7530e79c341f8690000000000000005
	//0x5f410f286333decd069582eb991d965e123f9c3bef2079198bf7145cc1ead0ac005489e0e317317fc4fc13809e838c3bf99492f3f078f088a5864f886bb3a7ef0183df5ddcefaecac331dde3d78ca45b0a664bd65a292d87a8b7530e79c341f8690000000000000005
}

func TestSig(t *testing.T) {
	byteBuf := bytes.NewBuffer([]byte{})
	_ = binary.Write(byteBuf, binary.LittleEndian, uint8(3))
	fmt.Println(byteBuf.Bytes())

	data, _ := transaction.EmptyWitnessArg.Serialize()
	fmt.Println(data, common.Bytes2Hex(data))
	fmt.Println(len(data))

	//signatureNum := make([]byte, 1)
	//binary.LittleEndian.PutUint64(signatureNum, 3)
	//fmt.Println(signatureNum)
	//
	//signatureNum = make([]byte, 1)
	//binary.LittleEndian.PutUint32(signatureNum, 3)
	//fmt.Println(signatureNum)
	//
	//molecule.Bytes2GoU8()
}

func TestGenerateMultiSignDigest(t *testing.T) {
	hash := "0x1b45d2d9524665f84f73780252cf6074a0d42c00db42d2cd670f36459ee7d507"
	client, err := getClientMainNet()
	if err != nil {
		t.Fatal(err)
	}
	res, err := client.GetTransaction(context.Background(), types.HexToHash(hash))
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res.Transaction.Hash)
	var txBuilder txbuilder.DasTxBuilder
	var dasTxBuilderTransaction txbuilder.DasTxBuilderTransaction
	dasTxBuilderTransaction.Transaction = res.Transaction
	txBuilder.DasTxBuilderTransaction = &dasTxBuilderTransaction

	sortArgsList := [][]byte{
		common.Hex2Bytes("0x567419c40d0f2c3566e7630ee32697560fa97a7b"),
		common.Hex2Bytes("0x543d8ec90d784f60cf920e76a359ae83839a5e7a"),
		common.Hex2Bytes("0x14dd22136ce74aee2a007c71e5440143dab7b326"),
		common.Hex2Bytes("0x619b019a75910e04d5f215ace571e5600d48b676"),
		common.Hex2Bytes("0x6d6a5e1df00e2cf82dd4dcfbba444a94119ae2de"),
	}
	//
	//wit := "0x3f010000100000003f0100003f0100002b01000000000305567419c40d0f2c3566e7630ee32697560fa97a7b543d8ec90d784f60cf920e76a359ae83839a5e7a14dd22136ce74aee2a007c71e5440143dab7b326619b019a75910e04d5f215ace571e5600d48b6766d6a5e1df00e2cf82dd4dcfbba444a94119ae2de534976631ae05be9873967c5c50ad69ebb88c4b3c61748b6aba52e5f9304eb6632aa7127d1303386f34b0a5fffbe70ad607d84d047674eecd66afa8b54a366530181b90f5ed1581ff01bef521d60aea6378a4d4243b98db679d25a561968e0d2e0192887024ad575ad2430c763b05bd85f74152ce232c5ab66f8e6d3311c23951f000dd35c5f2abbae1f55f6d8d33ea67783800070a19a64f8d3e3e52827d38e087f4e4331c28bd1c65ab29dd04c6157a77e7cb78a77309c10e63c1a09b31524ba7000"
	//bys := common.Hex2Bytes(wit)
	//fmt.Println("bys",common.Bytes2Hex(bys[:len(bys)-195]))

	digest, err := txBuilder.GenerateMultiSignDigest([]int{0}, 0, [][]byte{{}, {}, {}}, sortArgsList)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(common.Bytes2Hex(digest))
}
