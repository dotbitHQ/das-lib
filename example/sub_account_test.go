package example

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/DeAccountSystems/das-lib/common"
	"github.com/DeAccountSystems/das-lib/core"
	"github.com/DeAccountSystems/das-lib/dascache"
	"github.com/DeAccountSystems/das-lib/smt"
	"github.com/DeAccountSystems/das-lib/witness"
	"github.com/nervosnetwork/ckb-sdk-go/indexer"
	"github.com/nervosnetwork/ckb-sdk-go/types"
	"strings"
	"sync"
	"testing"
)

func TestGetSMTRoot(t *testing.T) {
	dc, err := getNewDasCoreTestnet2()
	if err != nil {
		t.Fatal(err)
	}
	hash := "0xd2ed490f6cec9543291b3b730d0f38a2e46258c8848c6ec7ac12a6f9fa0ffd7f"
	if res, err := dc.Client().GetTransaction(context.Background(), types.HexToHash(hash)); err != nil {
		t.Fatal(err)
	} else {
		fmt.Println("SMTRoot")
		contract, err := core.GetDasContractInfo(common.DASContractNameSubAccountCellType)
		if err != nil {
			t.Fatal(err)
		}
		for k, v := range res.Transaction.Outputs {
			if v.Type != nil && contract.IsSameTypeId(v.Type.CodeHash) {
				fmt.Println(common.Bytes2Hex(res.Transaction.OutputsData[k]))
			}
		}
	}
}

func TestSubAccountBuilderMapFromTx(t *testing.T) {
	dc, err := getNewDasCoreTestnet2()
	if err != nil {
		t.Fatal(err)
	}
	hash := "0x14bbfd4b1e576264bc844e74e7c7e5f39891877e3865621ef40f83d9c3d6cf20"
	if res, err := dc.Client().GetTransaction(context.Background(), types.HexToHash(hash)); err != nil {
		t.Fatal(err)
	} else {
		builderMaps, err := witness.SubAccountBuilderMapFromTx(res.Transaction)
		if err != nil {
			t.Fatal(err)
		}

		for _, builder := range builderMaps {
			fmt.Println("--------------------------------------------")
			fmt.Println(common.Bytes2Hex(builder.PrevRoot))
			fmt.Println(common.Bytes2Hex(builder.CurrentRoot))
			fmt.Println(builder.Version)
			fmt.Println(builder.Account)
			fmt.Println(builder.SubAccount)
			fmt.Println(builder.SubAccount.ExpiredAt, builder.CurrentSubAccount.ExpiredAt)
			fmt.Println(len(builder.SubAccount.Records), len(builder.CurrentSubAccount.Records))
			fmt.Println(common.Bytes2Hex(builder.SubAccount.Lock.Args), common.Bytes2Hex(builder.CurrentSubAccount.Lock.Args))
			fmt.Println(common.Bytes2Hex(builder.SubAccount.ToH256()))
		}
	}
}

func TestGenActionDataWitnessV2(t *testing.T) {
	fmt.Println(witness.GenActionDataWitness(common.DasActionCreateSubAccount, nil))
	fmt.Println(witness.GenActionDataWitnessV2(common.DasActionCreateSubAccount, nil, common.ParamManager))
}

func TestGetRefund(t *testing.T) {
	client, _ := getClientMainNet()
	res, _ := client.GetTransaction(context.Background(), types.HexToHash("0xfd1153209e99d26cfffc1b3583b223548cb8259af09ea9e4dbb2c8391d9dde46"))
	for i, v := range res.Transaction.Outputs {
		fmt.Println(common.Bytes2Hex(v.Lock.Args), i, v.Capacity)
	}
}

func TestPre(t *testing.T) {
	var wg sync.WaitGroup
	ca := dascache.NewDasCache(context.Background(), &wg)
	searchKey := &indexer.SearchKey{
		Script: &types.Script{
			CodeHash: types.HexToHash("0x18ab87147e8e81000ab1b9f319a5784d4c7b6c98a9cec97d738a5c11f69e7254"),
			HashType: types.HashTypeType,
			Args:     nil,
		},
		ScriptType: indexer.ScriptTypeType,
		Filter: &indexer.CellsFilter{
			Script: &types.Script{
				CodeHash: types.HexToHash("0x303ead37be5eebfcf3504847155538cb623a26f237609df24bd296750c123078"),
				HashType: types.HashTypeType,
				Args:     nil,
			},
			OutputDataLenRange: &[2]uint64{52, 53}, // hash + account id
		},
	}
	ca.AddOutPoint([]string{"0x204f2b221eaef83af915139e6ff5cc4ca1a08e245ec7a513b9c8049f0920b6a9-0"})
	client, _ := getClientMainNet()
	cells, _ := core.GetSatisfiedLimitLiveCell(client, ca, searchKey, 400, indexer.SearchOrderAsc)
	for _, v := range cells {
		fmt.Println(common.OutPointStruct2String(v.OutPoint))
	}
}

func TestSubAccountBuilderFromBytes(t *testing.T) {
	bys := common.Hex2Bytes("0x6461730800000041000000b0ad48aad0062a684a46aa98e2cffd55c20982ee0dbcbf29dc747c0b06c43a9618b66f607d80d0029f1ae9bdfd347552db2ac4fe51459f5af8880f1f28452c7301010000000020000000d8f23be95bdff4e509bbfd711de7ed3b88eb9c0af44cc15e2c67608c86a73a49200000009b9f8b3b4f7e1121a6a48d4b359285b55c6160bb889436d29d5ef4bb58821b9ecd0000004c4f9c519c2791e985d68babbaf7ed736189b2b82a42e55c54aeb81511c1e6bc07e92a6e4e80aff85fccab124dad8e50716abbfa5ed2a174040000000000000000000000004f01519e2a28071e0204a84ab07697c06f40c5c48637703e6320942f119f0f258414b821851e47a17042cc4febb4be64a66c27eb615ecd210000000000000000000000005101e50bb68f17990e17086a911458dae14dee75cfe41de23aee43841b9c03dbe08000000000000000000000000000000000000000400000000000000000000000004f6004000000010000005701000057010000300000008f000000a30000002401000031010000390100004101000042010000460100004e0100004f0100005f000000100000003000000031000000326df166e3f0a900a0aee043e31a4dea0f01ea3307e6e235f09d1b4220b75fbd012a00000005c9f53b1d85356b60453f867610888d89a0b667ad05c9f53b1d85356b60453f867610888d89a0b667addbcaa515cbd79477e17502a6e51dcdccadad869081000000180000002d00000042000000570000006c000000150000000c00000010000000010000000100000030150000000c00000010000000010000000100000030150000000c00000010000000010000000100000030150000000c00000010000000010000000100000030150000000c00000010000000010000000100000031090000002e303030312e62697482a03d620000000002d41e640000000000040000000000000000000000000000000000000000070000006d616e616765725f0000005f000000100000003000000031000000326df166e3f0a900a0aee043e31a4dea0f01ea3307e6e235f09d1b4220b75fbd012a00000005c9f53b1d85356b60453f867610888d89a0b667ad0515a33588908cf8edb27d1abe3852bf287abd3891")
	res, err := witness.SubAccountBuilderFromBytes(bys[common.WitnessDasTableTypeEndIndex:])
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(common.Bytes2Hex(res.SubAccount.Lock.Args), res.Account)
}

func TestGenSubAccountBytes(t *testing.T) {
	account := "aaa.bit"

	accountCharSet, err := common.AccountToAccountChars(account)
	if err != nil {
		t.Fatal(err)
	}
	subAccount := witness.SubAccount{
		Lock: &types.Script{
			CodeHash: types.HexToHash("0x8bb0413701cdd2e3a661cc8914e6790e16d619ce674930671e695807274bd14c"),
			HashType: types.HashTypeType,
			Args:     common.Hex2Bytes("0x05c9f53b1d85356b60453f867610888d89a0b667ad0515a33588908cf8edb27d1abe3852bf287abd3891"),
		},
		AccountId:            "0x338e9410a195ddf7fedccd99834ea6c5b6e5449c",
		AccountCharSet:       accountCharSet,
		Suffix:               ".aaa.bit",
		RegisteredAt:         1,
		ExpiredAt:            2,
		Status:               0,
		Records:              nil,
		Nonce:                0,
		EnableSubAccount:     0,
		RenewSubAccountPrice: 0,
	}
	param := witness.SubAccountParam{
		Signature:      nil,
		SignRole:       nil,
		PrevRoot:       []byte{2},
		CurrentRoot:    []byte{3},
		Proof:          []byte{4},
		SubAccount:     &subAccount,
		EditKey:        "",
		EditLockArgs:   nil,
		EditRecords:    nil,
		RenewExpiredAt: 0,
	}
	bys, err := param.NewSubAccountWitness()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(common.Bytes2Hex(bys))
}

func TestTestGenSubAccountBytes2(t *testing.T) {
	var param witness.SubAccountParam
	str := `{"Signature":null,"SignRole":null,"PrevRoot":"AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=","CurrentRoot":"2rrwuKes5atj5TKCbLBcyCroH8/yw/7b2fvBtOH9Pb8=","Proof":"TE8A","SubAccount":{"lock":{"code_hash":"0x326df166e3f0a900a0aee043e31a4dea0f01ea3307e6e235f09d1b4220b75fbd","hash_type":"type","args":"Bcn1Ox2FNWtgRT+GdhCIjYmgtmetBcn1Ox2FNWtgRT+GdhCIjYmgtmet"},"account_id":"0xdbcaa515cbd79477e17502a6e51dcdccadad8690","account_char_set":[{"char_set_name":1,"char":"0"},{"char_set_name":1,"char":"0"},{"char_set_name":1,"char":"0"},{"char_set_name":1,"char":"0"},{"char_set_name":1,"char":"1"}],"suffix":".0001.bit","registered_at":1648195840,"expired_at":1679731840,"status":0,"records":null,"nonce":0,"enable_sub_account":0,"renew_sub_account_price":0},"EditKey":"","EditLockScript":null,"EditRecords":null,"RenewExpiredAt":0}`
	_ = json.Unmarshal([]byte(str), &param)
	bys, err := param.NewSubAccountWitness()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(common.Bytes2Hex(bys))

}

func TestSMTRootVerify(t *testing.T) {
	dc, err := getNewDasCoreTestnet2()
	if err != nil {
		t.Fatal(err)
	}
	hashList := []string{
		"0x23d91172a788a407ecc01c65bc02e6129f109d3a62b21ed7b3cc16daa93c465c",
		"0xff5f6472261175c588607f1a4a3d70829dcd48c6382a8f858ccd2a533c3e344b",
		"0xe7135538f29787308a06a9c5f789b6ecea647ba242631c38d8ee29281a556652",
		"0x14bbfd4b1e576264bc844e74e7c7e5f39891877e3865621ef40f83d9c3d6cf20",
	}
	tree := smt.NewSparseMerkleTree(nil)

	for _, hash := range hashList {
		if res, err := dc.Client().GetTransaction(context.Background(), types.HexToHash(hash)); err != nil {
			t.Fatal(err)
		} else {
			builderMaps, err := witness.SubAccountBuilderMapFromTx(res.Transaction)
			if err != nil {
				t.Fatal(err)
			}

			for _, builder := range builderMaps {
				fmt.Println(strings.Repeat("-", 100))
				fmt.Println(fmt.Sprintf("%-20s %s", "current root", common.Bytes2Hex(builder.CurrentRoot)))

				key := smt.AccountIdToSmtH256(builder.SubAccount.AccountId)
				value := builder.SubAccount.ToH256()
				subAccount, _ := builder.ConvertToEditValue()
				switch string(builder.EditKey) {
				case common.EditKeyOwner, common.EditKeyManager:
					builder.SubAccount.Nonce++
					builder.SubAccount.Lock.Args = common.Hex2Bytes(subAccount.LockArgs)
					value = builder.SubAccount.ToH256()
				case common.EditKeyRecords:
					builder.SubAccount.Nonce++
					builder.SubAccount.Records = subAccount.Records
					value = builder.SubAccount.ToH256()
				case common.EditKeyExpiredAt:
					builder.SubAccount.Nonce++
					builder.SubAccount.ExpiredAt = subAccount.ExpiredAt
					value = builder.SubAccount.ToH256()
				}
				_ = tree.Update(key, value)
				root, _ := tree.Root()
				fmt.Println(fmt.Sprintf("%-20s %s", "tree value", common.Bytes2Hex(value)))
				fmt.Println(fmt.Sprintf("%-20s %s", "tree root", common.Bytes2Hex(root)))
			}
		}
	}
}
