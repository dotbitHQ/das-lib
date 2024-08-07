package example

import (
	"bytes"
	"context"
	"crypto/ecdsa"
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
	"math/big"
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
	digest := "0xf56c76e7c6c5b55a240ced59829b5d82e52f9c1e2cd2daa84425f24c3bcd4706"
	mmJson := `{"types":{"EIP712Domain":[{"name":"name","type":"string"},{"name":"version","type":"string"},{"name":"chainId","type":"uint256"},{"name":"verifyingContract","type":"address"}],"Action":[{"name":"action","type":"string"},{"name":"params","type":"string"}],"Cell":[{"name":"capacity","type":"string"},{"name":"lock","type":"string"},{"name":"type","type":"string"},{"name":"data","type":"string"},{"name":"extraData","type":"string"}],"Transaction":[{"name":"DAS_MESSAGE","type":"string"},{"name":"inputsCapacity","type":"string"},{"name":"outputsCapacity","type":"string"},{"name":"fee","type":"string"},{"name":"action","type":"Action"},{"name":"inputs","type":"Cell[]"},{"name":"outputs","type":"Cell[]"},{"name":"digest","type":"bytes32"}]},"primaryType":"Transaction","domain":{"name":"d.id","version":"1","chainId":1,"verifyingContract":"0x0000000000000000000000000000000020210722"},"message":{"DAS_MESSAGE":"TRANSFER FROM 0x444f01464dbbf88e13292ae42772d9143a107e8e(1442.15530562 CKB) TO ckb1qyqxmyfg2a5w0jt0rn4qzu7gzead5t87405qs8cqan(480.4228 CKB), 0x444f01464dbbf88e13292ae42772d9143a107e8e(961.73247312 CKB)","inputsCapacity":"1442.15530562 CKB","outputsCapacity":"1442.15527312 CKB","fee":"0.0000325 CKB","digest":"","action":{"action":"transfer","params":"0x00"},"inputs":[],"outputs":[]}}`
	address := "0x444f01464dbbf88e13292ae42772d9143a107e8e"

	chainId := 1
	var obj3 apitypes.TypedData
	oldChainId := fmt.Sprintf("chainId\":%d", chainId)
	newChainId := fmt.Sprintf("chainId\":\"%d\"", chainId)
	mmJson = strings.ReplaceAll(mmJson, oldChainId, newChainId)
	oldDigest := "\"digest\":\"\""
	newDigest := fmt.Sprintf("\"digest\":\"%s\"", digest)
	mmJson = strings.ReplaceAll(mmJson, oldDigest, newDigest)

	_ = json.Unmarshal([]byte(mmJson), &obj3)
	//var mmHash, signature []byte
	//mmHash, signature, err := sign.EIP712Signature(obj3, privateKey)
	//if err != nil {
	//	t.Fatal(err)
	//}

	//fmt.Println("EIP712Signature mmHash:", common.Bytes2Hex(mmHash))
	//fmt.Println("EIP712Signature signature:", common.Bytes2Hex(signature))

	//signData := append(signature, mmHash...)
	//hexChainId := fmt.Sprintf("%x", chainId)
	//chainIdData := common.Hex2Bytes(fmt.Sprintf("%016s", hexChainId))
	//signData = append(signData, chainIdData...)
	//fmt.Println("signData:", common.Bytes2Hex(signData))
	signature := "0x3570f9d1ce34c08b9e0c7369f522077e339c52695c7c902f261de5f5e55850476d10e22e14cfff87de3ff511194f453804ee8e1d6772bd33a40b2a1378b0237d1c"

	fmt.Println(len(signature))
	fmt.Println(sign.VerifyEIP712Signature(obj3, common.Hex2Bytes(signature), address))
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

func TestEcdsaP256Signature(t *testing.T) {

	//pubKey
	var pubkey *ecdsa.PublicKey
	pubkey.X = new(big.Int).SetBytes([]byte{216, 152, 197, 85, 225, 214, 251, 8, 101, 200, 149, 118, 238, 212, 67, 118, 33, 24, 122, 12, 126, 168, 234, 163, 164, 63, 57, 160, 100, 129, 107, 120})
	pubkey.Y = new(big.Int).SetBytes([]byte{185, 230, 151, 223, 144, 136, 150, 177, 230, 140, 106, 80, 73, 45, 143, 8, 237, 244, 33, 112, 238, 245, 116, 157, 155, 253, 75, 69, 60, 165, 32, 45})
	//signData
	signData := []byte{207, 196, 244, 62, 96, 221, 141, 119, 28, 192, 239, 225, 161, 130, 124, 133, 105, 188, 245, 104, 249, 88, 19, 245, 63, 142, 56, 142, 231, 252, 149, 236}
	//signature
	R := new(big.Int).SetBytes([]byte{58, 72, 38, 75, 25, 243, 23, 37, 15, 120, 166, 186, 35, 146, 95, 244, 128, 34, 235, 216, 7, 234, 102, 213, 162, 14, 56, 139, 232, 5, 211, 107})
	S := new(big.Int).SetBytes([]byte{171, 248, 163, 95, 53, 35, 169, 242, 19, 38, 48, 97, 199, 242, 102, 161, 60, 225, 214, 218, 45, 14, 175, 66, 34, 147, 76, 242, 223, 238, 131, 162})
	res, _ := sign.VerifyEcdsaP256Signature(signData, R, S, pubkey)
	fmt.Println(res)
}

func TestVerifyWebauthnSignature(t *testing.T) {
	//signData  LV
	//#1+1
	//length(pubkey_index) + pubkey_index +
	//#1+64
	//length(signature) + signature +
	//#1+64
	//length(pubkey)+ pubkey +
	//#1+ *
	//length(authnticator_data) + authnticator_data +
	//#2 + *
	//length(clientDataJson) + clientDataJson
	signData := "010040bca95af564cd3756efadc658e3e920ca09f4a2ab6e1283cb903b9a8a935fa39a4db720ffdec9317336607b55791c4b04540a97e5bd036f3164e1cee2464b4c434096e07df8713895932052ce68061c208aab9210fe30adb501b32729c24a250470ddce694298aae92e415031caa81dec6767c53fea9300db49ce10ea68bc8a08052549960de5880e8c687434170f6476605b8fe4aeb9a28632c7995cf3ba831d976305000000005f007b2274797065223a22776562617574686e2e676574222c226368616c6c656e6765223a2259574668222c226f726967696e223a22687474703a2f2f6c6f63616c686f73743a38303031222c2263726f73734f726967696e223a66616c73657d"
	data := "aaa"
	res, err := sign.VerifyWebauthnSignature([]byte(data), common.Hex2Bytes(signData), "0f76d29f1f522b440e99")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("verify result: ", res)
}
