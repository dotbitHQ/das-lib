package witness

import (
	"encoding/json"
	"github.com/dotbitHQ/das-lib/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/nervosnetwork/ckb-sdk-go/types"
	"testing"
)

type TestReverseSmtRecord struct {
	Version       ReverseSmtRecordVersion
	Action        ReverseSmtRecordAction
	Signature     []byte
	SignType      uint8
	Address       []byte
	Proof         []byte
	PrevNonce     uint32 `witness:",omitempty"`
	PrevAccount   string
	NextRoot      []byte
	NextAccount   string
	TestInterface DasWitness
}

type TestInterfaceType string

func (t TestInterfaceType) Gen() ([]byte, error) {
	return []byte(t), nil
}

func (t TestInterfaceType) Parse(data []byte) (DasWitness, error) {
	return TestInterfaceType(data), nil
}

func TestGenWitnessData(t *testing.T) {
	data, err := GenWitnessData(&TestReverseSmtRecord{
		Version:       ReverseSmtRecordVersion1,
		Action:        ReverseSmtRecordActionUpdate,
		Signature:     common.Hex2Bytes("0xd56f475c74374450d912eba19aae40b98669e1b0bf436caf9e045ca26a78dddc69d344424d7102c426fb79599c47745691e1903882f53be785eca7ab630297c91b"),
		SignType:      3,
		Address:       common.Hex2Bytes("0xdeefc10a42cd84c072f2b0e2fa99061a74a0698c"),
		Proof:         common.Hex2Bytes(""),
		PrevNonce:     7,
		PrevAccount:   "reverse-smt.bit",
		NextRoot:      common.Hex2Bytes(""),
		NextAccount:   "reverse-smt01.bit",
		TestInterface: TestInterfaceType("test_interface"),
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(common.Bytes2Hex(data))

	trsr := &TestReverseSmtRecord{
		TestInterface: TestInterfaceType(""),
	}
	if err := ParseFromBytes(data, trsr); err != nil {
		t.Fatal(err)
	}
	t.Logf("%#v", trsr)
}

func TestParseFromTx(t *testing.T) {
	txJson := `
{
    "version":"0x0",
    "cell_deps":[
        {
            "out_point":{
                "tx_hash":"0x6678cf2da360945b031170b2e23e776b684eac29b3f9c1f8a2cf62cb88f2a4f7",
                "index":"0x0"
            },
            "dep_type":"code"
        },
        {
            "out_point":{
                "tx_hash":"0xdf4337204e1fb77c9deac20afe0e4fcb83a8d986db0f7250874485dffd48e66e",
                "index":"0x0"
            },
            "dep_type":"code"
        },
        {
            "out_point":{
                "tx_hash":"0xd819d634d5593d1c7f22d8f954fd743c50eaea427fe669595b7f7a1109bbac6f",
                "index":"0x0"
            },
            "dep_type":"code"
        },
        {
            "out_point":{
                "tx_hash":"0x77cdb8d076e3780ef46c42e8f473e9ec2ea1d9521e1cf8ee0db9efb01671d341",
                "index":"0x0"
            },
            "dep_type":"code"
        },
        {
            "out_point":{
                "tx_hash":"0xf249a946f1302c34d63d437eaf345ce77b96c91f142cef3c356ec16f0ecc3f34",
                "index":"0x0"
            },
            "dep_type":"code"
        },
        {
            "out_point":{
                "tx_hash":"0x4811e5a7877b8b5cec6562e44a89006dbd180dc583dd114fd25f341e8d46db09",
                "index":"0x0"
            },
            "dep_type":"code"
        },
        {
            "out_point":{
                "tx_hash":"0x8ffa409ba07d74f08f63c03f82b7428d36285fe75b2173fc2476c0f7b80c707a",
                "index":"0x0"
            },
            "dep_type":"code"
        },
        {
            "out_point":{
                "tx_hash":"0x9e0823959e5b76bd010cc503964cced4f8ae84f3b03e94811b083f9765534ff1",
                "index":"0x0"
            },
            "dep_type":"code"
        },
        {
            "out_point":{
                "tx_hash":"0xa706f46e58e355a6d29d7313f548add21b875639ea70605d18f682c1a08740d6",
                "index":"0x0"
            },
            "dep_type":"code"
        },
        {
            "out_point":{
                "tx_hash":"0x747411fb3914dd7ca5488a0762c6f4e76f56387e83bcbb24e3a01afef1d5a5b4",
                "index":"0x0"
            },
            "dep_type":"code"
        },
        {
            "out_point":{
                "tx_hash":"0x209b35208da7d20d882f0871f3979c68c53981bcc4caa71274c035449074d082",
                "index":"0x0"
            },
            "dep_type":"code"
        },
        {
            "out_point":{
                "tx_hash":"0xf8de3bb47d055cdf460d93a2a6e1b05f7432f9777c8c474abf4eec1d4aee5d37",
                "index":"0x0"
            },
            "dep_type":"dep_group"
        },
        {
            "out_point":{
                "tx_hash":"0xa7ff448225fc131d657af882a3f97a8219be230d7e25d070a9282de89302c640",
                "index":"0x0"
            },
            "dep_type":"code"
        }
    ],
    "header_deps":[
        "0x6293045bb57025c28e724217420837d7a7d2a4f71f1364f58e6681353b2d0ddd"
    ],
    "inputs":[
        {
            "since":"0x0",
            "previous_output":{
                "tx_hash":"0x17789c41d0bd890fe305a3dc50ce92911111518ec3b29145bf6ecb868052ed36",
                "index":"0x0"
            }
        },
        {
            "since":"0x0",
            "previous_output":{
                "tx_hash":"0x5b5406b6bd97a0dd8a3acf123e9cc9844c99ea4b4bd83655e1d171dbc88e3efb",
                "index":"0x0"
            }
        }
    ],
    "outputs":[
        {
            "capacity":"0x4a817c800",
            "lock":{
                "code_hash":"0xf1ef61b6977508d9ec56fe43399a01e576086a76cf0f7c687d1418335e8c401f",
                "hash_type":"type",
                "args":"0x"
            },
            "type":{
                "code_hash":"0x8041560ab6bd812c4523c824f2dcf5843804a099cb2f69fcbd57c8afcef2ed5f",
                "hash_type":"type",
                "args":"0x"
            }
        },
        {
            "capacity":"0xe8d4a4e8f0",
            "lock":{
                "code_hash":"0x9bd7e06f3ecf4be0f2fcd2188b23f1b9fcc88e5d4b65a8637b17723bbda3cce8",
                "hash_type":"type",
                "args":"0xda44ed9db97056a06e471d3a1b6a1b82219e7232"
            },
            "type":null
        }
    ],
    "outputs_data":[
        "0xb4bdcdec0653e52b55db4567a303cf8df35392e9aa687667808ca3cac3cfa5e0",
        "0x"
    ],
    "witnesses":[
        "0x6461730a0000000400000001000000060000007570646174654100000006b3abdf1a885d2a4741d39250a1080d66e3ba47add98c091574b1feb886a68e20e587add97c600064f0cace958671ebabe83c351f5d6265f1808d16e2ec653601010000000314000000deefc10a42cd84c072f2b0e2fa99061a74a0698c030000004c4f00000000000000000020000000b4bdcdec0653e52b55db4567a303cf8df35392e9aa687667808ca3cac3cfa5e00f000000726576657273652d736d742e626974",
        "0x6461730a0000000400000001000000060000007570646174654100000006b3abdf1a885d2a4741d39250a1080d66e3ba47add98c091574b1feb886a68e20e587add97c600064f0cace958671ebabe83c351f5d6265f1808d16e2ec653601010000000314000000deefc10a42cd84c072f2b0e2fa99061a74a0698c030000004c4f00000000000000000020000000b4bdcdec0653e52b55db4567a303cf8df35392e9aa687667808ca3cac3cfa5e00f000000726576657273652d736d742e626974"
    ]
}
`
	var inTransaction = &struct {
		Version     hexutil.Uint    `json:"version"`
		HeaderDeps  []types.Hash    `json:"header_deps"`
		OutputsData []hexutil.Bytes `json:"outputs_data"`
		Witnesses   []hexutil.Bytes `json:"witnesses"`
	}{}

	if err := json.Unmarshal([]byte(txJson), inTransaction); err != nil {
		t.Fatal(err)
	}

	tx := &types.Transaction{}
	_ = gconv.Struct(inTransaction, tx)

	res := make([]*ReverseSmtRecord, 0)
	if err := ParseFromTx(tx, common.ActionDataTypeReverseSmt, &res); err != nil {
		t.Fatal(err)
	}
	t.Log(res)
}
