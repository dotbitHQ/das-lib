package example

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/DeAccountSystems/das-lib/common"
	"github.com/DeAccountSystems/das-lib/molecule"
	"github.com/DeAccountSystems/das-lib/witness"
	"github.com/nervosnetwork/ckb-sdk-go/types"
	"strings"
	"testing"
)

func TestActionDataBuilderFromTx(t *testing.T) {
	dc, err := getNewDasCoreTestnet2()
	if err != nil {
		t.Fatal(err)
	}
	hash := "0x114dcdb52147d5886b4fa62757dff30aa3144800d6b2583018b5c7a793ce61ff"
	if res, err := dc.Client().GetTransaction(context.Background(), types.HexToHash(hash)); err != nil {
		t.Fatal(err)
	} else {
		builder, err := witness.ActionDataBuilderFromTx(res.Transaction)
		if err != nil {
			t.Fatal(err)
		}
		fmt.Println("action name:", builder.Action)
		if builder.Action == common.DasActionBuyAccount {
			inviterScript, err := molecule.ScriptFromSlice(builder.Params[0], false)
			if err != nil {
				t.Fatal(err)
			}
			fmt.Println(common.Bytes2Hex(inviterScript.Args().RawData()))
			channelScript, err := molecule.ScriptFromSlice(builder.Params[1], false)
			if err != nil {
				t.Fatal(err)
			}
			fmt.Println(common.Bytes2Hex(channelScript.Args().RawData()))
		}
	}
}

func TestGenActionDataWitness(t *testing.T) {
	dc, err := getNewDasCoreTestnet2()
	if err != nil {
		t.Fatal(err)
	}
	inviter := "0xc9f53b1d85356b60453f867610888d89a0b667ad"
	inviterScript, _, err := dc.FormatAddressToDasLockScript(common.ChainTypeEth, inviter, true)
	if err != nil {
		t.Fatal(err)
	}
	channel := "0x15a33588908cf8edb27d1abe3852bf287abd3891"
	channelScript, _, err := dc.FormatAddressToDasLockScript(common.ChainTypeEth, channel, true)
	if err != nil {
		t.Fatal(err)
	}
	params := witness.GenBuyAccountParams(inviterScript, channelScript)
	witBys, err := witness.GenActionDataWitness(common.DasActionBuyAccount, params)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(common.Bytes2Hex(witBys))
}

func TestActionDataFromTx(t *testing.T) {
	txStr := `{"version":"0x0","cell_deps":[{"out_point":{"tx_hash":"0x9ec8d03db8bb89c5a87280e2e60ef7b130025b3726c67ee726db9fc8a6e6d9a7","index":"0x4"},"dep_type":"code"},{"out_point":{"tx_hash":"0x4cc3d61c0239a10afa0ec8b096ed2b5ce982469672946ade6cb78305868f8a88","index":"0x0"},"dep_type":"code"},{"out_point":{"tx_hash":"0x77cdb8d076e3780ef46c42e8f473e9ec2ea1d9521e1cf8ee0db9efb01671d341","index":"0x0"},"dep_type":"code"},{"out_point":{"tx_hash":"0x823a0a983c36ce967b80abd91fc4daa19ad67253ad599c8926d00107fccd6fdb","index":"0x0"},"dep_type":"code"},{"out_point":{"tx_hash":"0x88e11044559418c6f4960a164cd8883ed37778597538ea5c8b1227f4e70f21c8","index":"0x0"},"dep_type":"code"},{"out_point":{"tx_hash":"0xecc95d73969c83c2182975d0e30aec03d1cb652a150c912a7295fc27ec46e4e1","index":"0x0"},"dep_type":"code"},{"out_point":{"tx_hash":"0xf8de3bb47d055cdf460d93a2a6e1b05f7432f9777c8c474abf4eec1d4aee5d37","index":"0x0"},"dep_type":"dep_group"},{"out_point":{"tx_hash":"0x09a9eafc93d3a452190caa7a389c53261be2207c1561f27f2ddee2ce658118e7","index":"0x0"},"dep_type":"code"},{"out_point":{"tx_hash":"0x7dc4ae8fe597045fbd7fe78f2bd26435644a69b755de3824a856f681bacb732b","index":"0x0"},"dep_type":"code"}],"header_deps":[],"inputs":[{"since":"0x0","previous_output":{"tx_hash":"0xf49884367b9bacbdccb112ad31c877dd92470c161ecc866337e193a416d0c193","index":"0x0"}}],"outputs":[{"capacity":"0x4ae0d81f0","lock":{"code_hash":"0x326df166e3f0a900a0aee043e31a4dea0f01ea3307e6e235f09d1b4220b75fbd","hash_type":"type","args":"0x04a2ac25bf43680c05abe82c7b1bcc1a779cff8d5d04a2ac25bf43680c05abe82c7b1bcc1a779cff8d5d"},"type":{"code_hash":"0x61711416468fa5211ead5f24c6f3efadfbbc332274c5d40e50c6feadcb5f9606","hash_type":"type","args":"0x"}}],"outputs_data":["0x747a6830312e626974"],"witnesses":["0x550000001000000055000000550000004100000088e870cc4a8ad32db2f9a5b4e75a26a19e7abe1dd951d5f7b544a8f3fa18362032d4a2a6a34db2adf651b89cbf44847ae36f909c1bca354b7b8b6c2501ea750a01","0x646173000000002d0000000c000000280000001800000072656465636c6172655f726576657273655f7265636f72640100000000","0x646173010000000a010000100000000a0100000a010000fa0000001000000014000000180000000000000002000000de000000de0000002400000038000000b9000000c1000000c9000000d1000000d9000000da000000a35ea5d5ef43a74e95351254802c334237bdde2b81000000180000002d00000042000000570000006c000000150000000c00000010000000020000000100000074150000000c0000001000000002000000010000007a150000000c00000010000000020000000100000068150000000c00000010000000010000000100000030150000000c0000001000000001000000010000003180cbde60000000000000000000000000000000000000000000000000000000000104000000","0x646173700000002800000010000000180000002000000000c817a80400000000e1f505000000001027000000000000","0x646173680000004902000010000000110000007d010000016c0100002c0000004c0000006c0000008c000000ac000000cc000000ec0000000c0100002c0100004c0100001106d9eaccde0995a7e07e80dd0ce7509f21752538dfdd1ee2526d24574846b10fbff871dd05aee1fda2be38786ad21d52a2765c6025d1ef6927d761d51a3cd14ff58f2c76b4ac26fdf675aa82541e02e4cf896279c6d6982d17b959788b2f0c08d1cdc6ab92d9cabe0096a2c7642f73d0ef1b24c94c43f21c6c3a32ffe0bb5e6c8441233f00741955f65e476721a1a5417997c1e4368801c99c7f617f8b754467d48c0911e406518de2116bd91c6af37c05f1db23334ca829d2af3042427e446ce893d9e64720388ee2faa570dd4b81986f7c4743fcbfda177b91bad6de681f0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000061711416468fa5211ead5f24c6f3efadfbbc332274c5d40e50c6feadcb5f9606cc000000180000003c0000006000000084000000a8000000209b35208da7d20d882f0871f3979c68c53981bcc4caa71274c035449074d08200000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000f273e6c581ad6bda82315bc06f1b9df4efae20fc5394995231ea96ab2b0ee3dd000000007dc4ae8fe597045fbd7fe78f2bd26435644a69b755de3824a856f681bacb732b00000000"]}`
	var tx types.Transaction
	_ = json.Unmarshal([]byte(txStr), &tx)
	builder, err := witness.ActionDataBuilderFromTx(&tx)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(builder.Action)
}

func TestActionDataBuilderFromWitness(t *testing.T) {
	dc, err := getNewDasCoreTestnet2()
	if err != nil {
		t.Fatal(err)
	}
	str := `0x01d68d10ef23478f8c8e7561cbee73a3902c4600d8392c9c63cedd59abca0c79-1
0x001ecc01993d7f79a5c7c048e5e602229b65ee2bc9cd8c161b3ae24f92179a85-0
0x00bb9156878a2d82d3274d9ad4b41ecf1f66c5b674a531da64f1fa4912fccf3c-0
0x007c69fc5476179fac1de11bd0644d44fb28d5b793f2dff96489f62410353658-1
0x132d8de77dc067a777d5364be0ba25084ba81edd555e0b2dffbef341e6cff61a-0
0x0084c5ac42f97f89145b2e7aa201b84fc5b7487a4c399dd6248c0e848e285da8-1
0x02cdbbcdfec4961b3310132c90af4bdb9e1a723776f56791abda4ebbbbee1a06-0
0x1728f9931d6f4d77227c6e6646419b230cfbb598269ed2cfc446329ea8a67df9-0
0x02285de49f5cac9770af3e1556508310c7ddf9a16bfe9dacad9f2977dca3e68e-0
0x023fed77f3a72c34dc9642aa7d542128a2668020193923852009be364981171d-0
0x00db8f2e203cd07d988aa60faa98f9f2b6aa3bf60f40d5754b243b4bbe8e3593-0
0x009e80169f96c4dc853f7758940e06c338a35db7672a76a0773a21c28f1e1f4f-0
0x01d68d10ef23478f8c8e7561cbee73a3902c4600d8392c9c63cedd59abca0c79-0
0x05180eee3fdbf89c68739969e8fe475316cf1afb35d507f0f8c4653ae45ea81d-0
0x519889f006591352cb689df374ccf738ecf2a28b812da78d3d7072e19ce43e81-0
0x00c4fe4f67f8d2611254783280698be9ac429168cebff76ec79b9ac968eacb16-0
0x0043936cb7bf7ac69dcc990aad84324c9cdc3f09742af715d29d5d0f212a9eb0-1
0x0343c250842fc57daef9fc30e5b9e1270c753a43215db46b19563c588417fcae-0
0x0074c7f0a69ad5f4f5705865389aebf27497c50de67e1c273be3c13525907690-0
0x1987deda6808681a8548f692528e84268b634f00fc5f73d40c14a5552444b8b3-0
0x00bb9156878a2d82d3274d9ad4b41ecf1f66c5b674a531da64f1fa4912fccf3c-1
0x0086e14df3c02fa3fe60969cd9305df685615c1968bcce2d6d9263196a8e171e-1
0x0550b7212b07002849ec615a1f83aac538a6dc561c83a76b30273439c867c026-0
0x0226a7bb24e1aace9c20db372841d89c51af83c5595879036a8bb2133c4ba4db-0
0x040094a604ad88cbb70190a4a602963851f8e106b80a92d1943a519b3c34b42d-0
0x0213c46b98b9d18cb30cb23db2ad47636192a25c67e5399a34f339ff12e7d084-0`
	list := strings.Split(str, "\n")
	for _, v := range list {
		outpoint := common.String2OutPointStruct(v)
		res, err := dc.Client().GetTransaction(context.Background(), outpoint.TxHash)
		if err != nil {
			t.Fatal(err)
		}
		builder, err := witness.ActionDataBuilderFromWitness(res.Transaction.Witnesses[len(res.Transaction.Inputs)])
		if err != nil {
			t.Fatal(err)
		}
		fmt.Println(builder.Action, outpoint)
	}

}
