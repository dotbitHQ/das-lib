package example

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/DeAccountSystems/das-lib/common"
	"github.com/DeAccountSystems/das-lib/sign"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
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
