package example

import (
	"context"
	"fmt"
	"github.com/dotbitHQ/das-lib/common"
	"github.com/dotbitHQ/das-lib/witness"
	"github.com/nervosnetwork/ckb-sdk-go/rpc"
	"github.com/nervosnetwork/ckb-sdk-go/types"
	"testing"
	"time"
)

func TestSubAccountMintSign(t *testing.T) {
	sams := witness.SubAccountMintSign{
		Version:            witness.SubAccountMintSignVersion1,
		Signature:          []byte{},
		ExpiredAt:          uint64(time.Now().Unix()),
		AccountListSmtRoot: []byte{},
	}
	dataBys := sams.GenSubAccountMintSignBytes()

	var sanb witness.SubAccountNewBuilder
	res, _ := sanb.ConvertSubAccountMintSignFromBytes(dataBys)
	fmt.Println(res.Version, res.ExpiredAt, res.Signature, res.AccountListSmtRoot)
}

func TestSubAccountNew(t *testing.T) {
	san := witness.SubAccountNew{
		Version:   0,
		Signature: nil,
		SignRole:  nil,
		NewRoot:   nil,
		Proof:     nil,
		Action:    "",
		SubAccountData: &witness.SubAccountData{
			Lock: &types.Script{
				CodeHash: types.Hash{},
				HashType: "",
				Args:     nil,
			},
			AccountId:            common.Bytes2Hex(common.GetAccountIdByAccount("aaa.bit")),
			AccountCharSet:       nil,
			Suffix:               "",
			RegisteredAt:         0,
			ExpiredAt:            0,
			Status:               0,
			Records:              nil,
			Nonce:                0,
			EnableSubAccount:     0,
			RenewSubAccountPrice: 0,
		},
		EditKey:        "",
		EditValue:      nil,
		EditLockArgs:   nil,
		EditRecords:    nil,
		RenewExpiredAt: 0,
		PrevRoot:       nil,
		CurrentRoot:    nil,
	}
	dataBys, err := san.GenSubAccountNewBytes()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(dataBys)

	var sanb witness.SubAccountNewBuilder
	subAcc, err := sanb.ConvertSubAccountNewFromBytes(dataBys)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(subAcc.SubAccountData.AccountId, subAcc.Version)

	subAcc.Version = witness.SubAccountNewVersion2
	dataBys, err = subAcc.GenSubAccountNewBytes()
	fmt.Println(dataBys)
	subAcc, err = sanb.ConvertSubAccountNewFromBytes(dataBys)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(subAcc.SubAccountData.AccountId, subAcc.Version)
}

func TestSubAccountNewMapFromTx(t *testing.T) {
	dc, err := getNewDasCoreTestnet2()
	if err != nil {
		t.Fatal(err)
	}
	hash := "0xa2178d7bd194fcd9f9d7533081ee51a0ba76e4028448052a02473a59958a50c7"
	if res, err := dc.Client().GetTransaction(context.Background(), types.HexToHash(hash)); err != nil {
		t.Fatal(err)
	} else {
		var sanb witness.SubAccountNewBuilder
		resMap, err := sanb.SubAccountNewMapFromTx(res.Transaction)
		if err != nil {
			t.Fatal(err)
		}
		for k, v := range resMap {
			fmt.Println(k, v.SubAccountData.AccountId, v.EditKey, v.EditRecords, v.EditLockArgs, v.RenewExpiredAt)
		}
	}

}

func TestUpdateSubAccountTx(t *testing.T) {
	str := ` {"version":"0x0","cell_deps":[{"out_point":{"tx_hash":"0x12fceb12d070d519c035d4b5893675a2043e0f96606d9aa9e156d73f19ec9af1","index":"0x0"},"dep_type":"code"},{"out_point":{"tx_hash":"0xd1e83290fb5e50edba9063988f48f139ec671466867289374c9c5ce5a7e45893","index":"0x0"},"dep_type":"code"},{"out_point":{"tx_hash":"0x8c6a6a31afc1375f3c49d4ef345f302847f76a7dc6bfb5c4a33df294560a9856","index":"0x0"},"dep_type":"code"},{"out_point":{"tx_hash":"0xd478afbac379544ed61b8c469cc9b1de5ae43be54d922248ba44858651d653f3","index":"0x1"},"dep_type":"code"},{"out_point":{"tx_hash":"0xf1a0b0f007ca3d456142f7b03fd36cf70ad8a171cdcaadffc239d784c8753ac5","index":"0x1"},"dep_type":"code"},{"out_point":{"tx_hash":"0x10b000eb473abf0847655a02ad8384ea808e8f9b88a3240f695475984f7e674d","index":"0x0"},"dep_type":"code"},{"out_point":{"tx_hash":"0xb1e95aa34d8a3f27207868d88751757b7b7607c7430a083d982195179d8bd440","index":"0x0"},"dep_type":"code"},{"out_point":{"tx_hash":"0xa918c7abee1a612d3d8222642638ca51ed442c855ed9d11c6b37d9df71aeed06","index":"0x0"},"dep_type":"code"},{"out_point":{"tx_hash":"0x5f1e03c7dc777d3cbec80aec8d3cefbf220ada32de5568155ff2c81acf99cce3","index":"0x0"},"dep_type":"code"},{"out_point":{"tx_hash":"0x243d467d3a0c3355e64a03baab4f4850ebe2133b6ac34586e7f21b2248d997f0","index":"0x0"},"dep_type":"code"},{"out_point":{"tx_hash":"0x4c817bc81165aae004f0961d583492e95759212edaa210afc434766998ce2670","index":"0x0"},"dep_type":"code"},{"out_point":{"tx_hash":"0xf8de3bb47d055cdf460d93a2a6e1b05f7432f9777c8c474abf4eec1d4aee5d37","index":"0x0"},"dep_type":"dep_group"},{"out_point":{"tx_hash":"0xa7ff448225fc131d657af882a3f97a8219be230d7e25d070a9282de89302c640","index":"0x0"},"dep_type":"code"}],"header_deps":[],"inputs":[{"since":"0x0","previous_output":{"tx_hash":"0xe5412417a5cb49071cb6be684e7f3bbae06ee64f96dfb02ffacf407c1d80b67c","index":"0x1"}},{"since":"0x0","previous_output":{"tx_hash":"0xaedb65e1eaf1b3e1d197ddd2e3d896574d17ea9001eee3112eac4e6e56abb704","index":"0x1"}},{"since":"0x0","previous_output":{"tx_hash":"0xaedb65e1eaf1b3e1d197ddd2e3d896574d17ea9001eee3112eac4e6e56abb704","index":"0x2"}},{"since":"0x0","previous_output":{"tx_hash":"0xaedb65e1eaf1b3e1d197ddd2e3d896574d17ea9001eee3112eac4e6e56abb704","index":"0x3"}}],"outputs":[{"capacity":"0x53d1a72e0","lock":{"code_hash":"0xf1ef61b6977508d9ec56fe43399a01e576086a76cf0f7c687d1418335e8c401f","hash_type":"type","args":"0x"},"type":{"code_hash":"0x8bb0413701cdd2e3a661cc8914e6790e16d619ce674930671e695807274bd14c","hash_type":"type","args":"0x64331471e7f3a760b8caf97811c4cb80565e76ff"}},{"capacity":"0x9d0373ea2f","lock":{"code_hash":"0x9bd7e06f3ecf4be0f2fcd2188b23f1b9fcc88e5d4b65a8637b17723bbda3cce8","hash_type":"type","args":"0xd0e1f9a79ab9361821cbc3b31fccee094cafacd0"},"type":null}],"outputs_data":["0x0065cd1d000000000000000000000000","0x"],"witnesses":["0x0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000","0x5500000010000000550000005500000041000000635770ca8ac7aa8ef35f756da96b0a844791bc0ffbf2e3635bd9ced2559719d85a59a1160660b8ff1aab8faeeaffd2a23ee0dac76d52147addcb49f1ef9f3ddb01","0x","0x","0x64617300000000260000000c00000022000000120000007570646174655f7375625f6163636f756e7400000000","0x0400000001000000410000002f96b3a5cae015ce8d6b096a375b09d2f3040eb5551368f83c8e9ee4a786c1f03c867be35066611910088dc5600d05b3f083fb5a75372799fe757969d7f440d801010000000008000000b75c906300000000200000002f96b3a5cae015ce8d6b096a375b09d2f3040eb5551368f83c8e9ee4a786c1f03c867be35066611910088dc5600d05b3f083fb5a75372799fe757969d7f440d801","0x646173080000000400000002000000000000000000000000000000080000000000000000000000200000004dd5d8f750a9997ff93ea42d8c8494fe05aa4dfa354a55e0ea5bf901fd6e4b0ecd0000004c4f9b519bc5735e3b9a08c3ca63fd852a58ac3984125159bb23d9350a7a1a11531934ffd7e856edd2cc15a0869b6fe3d455949a754cac82060000000000000000000000004f02519ef28c6e5e93dfdfc31be43d42fb396fdcd6883230940af9643a11ace79681a7131258b3d10a673e855b0a23a0653d350863d7bf2d000000000000000000000000510122b4bcfc5aa67276bafeecebbd3f40a5cc430177b6521264e80fc84a8a7bc0e800000000000000000000000000000000000000000000000000000000000000004f60a6010000a6010000300000008f000000a30000006f01000080010000880100009001000091010000950100009d0100009e0100005f000000100000003000000031000000326df166e3f0a900a0aee043e31a4dea0f01ea3307e6e235f09d1b4220b75fbd012a00000005c9f53b1d85356b60453f867610888d89a0b667ad05c9f53b1d85356b60453f867610888d89a0b667ad9b145673b11048b062f989f7ed5c4c36bd6d5ee8cc00000024000000390000004e00000063000000780000008d000000a2000000b7000000150000000c00000010000000020000000100000074150000000c00000010000000020000000100000065150000000c00000010000000020000000100000073150000000c00000010000000020000000100000074150000000c00000010000000010000000100000030150000000c00000010000000010000000100000031150000000c0000001000000001000000010000002d150000000c000000100000000100000001000000300d0000002e32303232313133302e6269742823876300000000a856686500000000000400000000000000000000000000000000000000000000000000000000","0x646173640000008c0000003c00000040000000480000005000000054000000580000005c000000640000006c000000740000007c0000008000000084000000880000002a000000000edbcb0400000000e1f50500000000008d27002c0100008813000010270000000000001027000000000000102700000000000010270000000000002c0100002c0100002c0100002c010000","0x6461737100000080000000300000003800000040000000480000005000000058000000600000006800000070000000780000007c00000000c817a804000000009435770000000000e1f5050000000000e1f50500000000a086010000000000a086010000000000a086010000000000a086010000000000a0860100000000002c0100002c010000","0x646173a286010039000000006100620063006400650066006700680069006a006b006c006d006e006f0070007100720073007400750076007700780079007a00","0x646173a18601001b0000000130003100320033003400350036003700380039002d00","0x64617368000000b90200001000000011000000c501000001b401000034000000540000007400000094000000b4000000d4000000f400000014010000340100005401000074010000940100001106d9eaccde0995a7e07e80dd0ce7509f21752538dfdd1ee2526d24574846b10fbff871dd05aee1fda2be38786ad21d52a2765c6025d1ef6927d761d51a3cd14ff58f2c76b4ac26fdf675aa82541e02e4cf896279c6d6982d17b959788b2f0c08d1cdc6ab92d9cabe0096a2c7642f73d0ef1b24c94c43f21c6c3a32ffe0bb5e6c8441233f00741955f65e476721a1a5417997c1e4368801c99c7f617f8b754467d48c0911e406518de2116bd91c6af37c05f1db23334ca829d2af3042427e449438124abdf4cbbfd61065e8b64523172bef5eefe27cb769c40acaf036aa89c200000000000000000000000000000000000000000000000000000000000000001a3f02aa89651a18112f0c21d0ae370a86e13f6a060c378184cd859a7bb6520361711416468fa5211ead5f24c6f3efadfbbc332274c5d40e50c6feadcb5f96068bb0413701cdd2e3a661cc8914e6790e16d619ce674930671e695807274bd14c4fd085557b4ef857b0577723bbf0a2e94081bbe3114de847cd9db01abaeb4f4ef40000001c000000400000006400000088000000ac000000d0000000209b35208da7d20d882f0871f3979c68c53981bcc4caa71274c035449074d08200000000747411fb3914dd7ca5488a0762c6f4e76f56387e83bcbb24e3a01afef1d5a5b4000000000000000000000000000000000000000000000000000000000000000000000000000000008ffa409ba07d74f08f63c03f82b7428d36285fe75b2173fc2476c0f7b80c707a000000009e0823959e5b76bd010cc503964cced4f8ae84f3b03e94811b083f9765534ff100000000a706f46e58e355a6d29d7313f548add21b875639ea70605d18f682c1a08740d600000000"]}`
	tx, err := rpc.TransactionFromString(str)
	if err != nil {
		t.Fatal(err)
	}
	var sanb witness.SubAccountNewBuilder
	res, err := sanb.SubAccountNewMapFromTx(tx)
	if err != nil {
		t.Fatal(err)
	}
	for k, v := range res {
		fmt.Println(k, v.Version, v.Action, v.Account)
		fmt.Println(v.SubAccountData.AccountId, v.SubAccountData.RegisteredAt, common.Bytes2Hex(v.SubAccountData.Lock.Args))
	}

	mintSign, err := sanb.ConvertSubAccountMintSignFromBytes(common.Hex2Bytes("0x0400000001000000410000002f96b3a5cae015ce8d6b096a375b09d2f3040eb5551368f83c8e9ee4a786c1f03c867be35066611910088dc5600d05b3f083fb5a75372799fe757969d7f440d801010000000008000000b75c906300000000200000002f96b3a5cae015ce8d6b096a375b09d2f3040eb5551368f83c8e9ee4a786c1f03c867be35066611910088dc5600d05b3f083fb5a75372799fe757969d7f440d801"))
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(mintSign.Version, common.Bytes2Hex(mintSign.Signature), common.Bytes2Hex(mintSign.AccountListSmtRoot))
}
