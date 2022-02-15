package example

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/DeAccountSystems/das-lib/common"
	"github.com/DeAccountSystems/das-lib/sign"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/signer/core"
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

	privateKey := ""
	address := "0xdD3b3D0F3FA9546a5616d0200b83f784a5220ae8"

	fmt.Println("----------- struct -------------")
	var typesStandard = core.Types{
		"EIP712Domain": {
			{
				"chainId",
				"uint256",
			},
			{
				"name",
				"string",
			},
			{
				"verifyingContract",
				"address",
			},
			{
				"version",
				"string",
			},
		},
		"Action": {
			{
				"action",
				"string",
			},
			{
				"params",
				"string",
			},
		},
		"Cell": {
			{
				"capacity",
				"string",
			},
			{
				"lock",
				"string",
			},
			{
				"type",
				"string",
			},
			{
				"data",
				"string",
			},
			{
				"extraData",
				"string",
			},
		},
		"Transaction": {
			{
				"DAS_MESSAGE",
				"string",
			},
			{
				"inputsCapacity",
				"string",
			},
			{
				"outputsCapacity",
				"string",
			},
			{
				"fee",
				"string",
			},
			{
				"action",
				"Action",
			},
			{
				"inputs",
				"Cell[]",
			},
			{
				"outputs",
				"Cell[]",
			},
			{
				"digest",
				"bytes32",
			},
		},
	}
	var domainStandard = core.TypedDataDomain{
		ChainId:           math.NewHexOrDecimal256(1),
		Name:              "da.systems",
		VerifyingContract: "0x0000000000000000000000000000000020210722",
		Version:           "1",
	}
	var messageStandard = map[string]interface{}{
		"DAS_MESSAGE":     "EDIT RECORDS OF ACCOUNT 5ph2lc3zs6x.bit",
		"inputsCapacity":  "221.9993 CKB",
		"outputsCapacity": "221.9992 CKB",
		"fee":             "0.0001 CKB",
		"action": map[string]interface{}{
			"action": "edit_records",
			"params": "0x01",
		},
		"inputs": []interface{}{
			map[string]interface{}{
				"capacity":  "221.9993 CKB",
				"lock":      "das-lock,0x01,0x05c9f53b1d85356b60453f867610888d89a0b667...",
				"type":      "account-cell-type,0x01,0x",
				"data":      "{ account: 5ph2lc3zs6x.bit, expired_at: 1658835295 }",
				"extraData": "{ status: 0, records_hash: 0x55478d76900611eb079b22088081124ed6c8bae21a05dd1a0d197efcc7c114ce }",
			},
		},
		"outputs": []interface{}{
			map[string]interface{}{
				"capacity":  "221.9992 CKB",
				"lock":      "das-lock,0x01,0x05c9f53b1d85356b60453f867610888d89a0b667...",
				"type":      "account-cell-type,0x01,0x",
				"data":      "{ account: 5ph2lc3zs6x.bit, expired_at: 1658835295 }",
				"extraData": "{ status: 0, records_hash: 0x17970d6aa6704f8d9084fbd5ae02c374eaa9152589062af29d3dc3d15b9e7802 }",
			},
		},
		"digest": "0x2277c4591b9fdf7403289bbaa9a8d43dc0e9cf9ecf46e416a057bc594a899dcb",
	}
	var typedData = core.TypedData{
		Types:       typesStandard,
		PrimaryType: "Transaction",
		Domain:      domainStandard,
		Message:     messageStandard,
	}
	_, signature, err := sign.EIP712Signature(typedData, privateKey)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("sig: ", common.Bytes2Hex(signature))
	ok, err := sign.VerifyEIP712Signature(typedData, signature, address)
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Fatal(fmt.Errorf("verify failed"))
	}

	fmt.Println("----------- DAS edit record -------------")
	data712 := `{
    "types": {
        "EIP712Domain": [
            {
                "name": "chainId",
                "type": "uint256"
            },
            {
                "name": "name",
                "type": "string"
            },
            {
                "name": "verifyingContract",
                "type": "address"
            },
            {
                "name": "version",
                "type": "string"
            }
        ],
        "Action": [
            {
                "name": "action",
                "type": "string"
            },
            {
                "name": "params",
                "type": "string"
            }
        ],
        "Cell": [
            {
                "name": "capacity",
                "type": "string"
            },
            {
                "name": "lock",
                "type": "string"
            },
            {
                "name": "type",
                "type": "string"
            },
            {
                "name": "data",
                "type": "string"
            },
            {
                "name": "extraData",
                "type": "string"
            }
        ],
        "Transaction": [
            {
                "name": "DAS_MESSAGE",
                "type": "string"
            },
            {
                "name": "inputsCapacity",
                "type": "string"
            },
            {
                "name": "outputsCapacity",
                "type": "string"
            },
            {
                "name": "fee",
                "type": "string"
            },
            {
                "name": "action",
                "type": "Action"
            },
            {
                "name": "inputs",
                "type": "Cell[]"
            },
            {
                "name": "outputs",
                "type": "Cell[]"
            },
            {
                "name": "digest",
                "type": "bytes32"
            }
        ]
    },
    "primaryType": "Transaction",
    "domain": {
        "chainId": "1",
        "name": "da.systems",
        "verifyingContract": "0x0000000000000000000000000000000020210722",
        "version": "1"
    },
    "message": {
        "DAS_MESSAGE": "EDIT RECORDS OF ACCOUNT 5ph2lc3zs6x.bit",
        "inputsCapacity": "221.9993 CKB",
        "outputsCapacity": "221.9992 CKB",
        "fee": "0.0001 CKB",
        "action": {
            "action": "edit_records",
            "params": "0x01"
        },
        "inputs": [
            {
                "capacity": "221.9993 CKB",
                "lock": "das-lock,0x01,0x05c9f53b1d85356b60453f867610888d89a0b667...",
                "type": "account-cell-type,0x01,0x",
                "data": "{ account: 5ph2lc3zs6x.bit, expired_at: 1658835295 }",
                "extraData": "{ status: 0, records_hash: 0x55478d76900611eb079b22088081124ed6c8bae21a05dd1a0d197efcc7c114ce }"
            }
        ],
        "outputs": [
            {
                "capacity": "221.9992 CKB",
                "lock": "das-lock,0x01,0x05c9f53b1d85356b60453f867610888d89a0b667...",
                "type": "account-cell-type,0x01,0x",
                "data": "{ account: 5ph2lc3zs6x.bit, expired_at: 1658835295 }",
                "extraData": "{ status: 0, records_hash: 0x17970d6aa6704f8d9084fbd5ae02c374eaa9152589062af29d3dc3d15b9e7802 }"
            }
        ],
        "digest": "0x2277c4591b9fdf7403289bbaa9a8d43dc0e9cf9ecf46e416a057bc594a899dcb"
    }
}`
	var obj core.TypedData
	_ = json.Unmarshal([]byte(data712), &obj)
	_, signature, err = sign.EIP712Signature(obj, privateKey)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("sig: ", common.Bytes2Hex(signature))
	ok, err = sign.VerifyEIP712Signature(typedData, signature, address)
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Fatal(fmt.Errorf("verify failed"))
	}

	fmt.Println("----------- sample -------------")
	data := `{
    "types": {
        "EIP712Domain": [
            {
                "name": "name",
                "type": "string"
            },
            {
                "name": "version",
                "type": "string"
            },
            {
                "name": "chainId",
                "type": "uint256"
            },
            {
                "name": "verifyingContract",
                "type": "address"
            }
        ],
        "Person": [
            {
                "name": "name",
                "type": "string"
            },
            {
                "name": "test",
                "type": "uint8"
            },
            {
                "name": "wallet",
                "type": "address"
            }
        ],
        "Mail": [
            {
                "name": "from",
                "type": "Person"
            },
            {
                "name": "to",
                "type": "Person"
            },
            {
                "name": "contents",
                "type": "string"
            }
        ]
    },
    "primaryType": "Mail",
    "domain": {
        "name": "Ether Mail",
        "version": "1",
        "chainId": "1",
        "verifyingContract": "0xCCCcccccCCCCcCCCCCCcCcCccCcCCCcCcccccccC"
    },
    "message": {
        "from": {
            "name": "Cow",
            "test": "3",
            "wallet": "0xcD2a3d9F938E13CD947Ec05AbC7FE734Df8DD826"
        },
        "to": {
            "name": "Bob",
            "wallet": "0xbBbBBBBbbBBBbbbBbbBbbbbBBbBbbbbBbBbbBBbB",
            "test": "2"
        },
        "contents": "Hello, Bob!"
    }
}`
	var obj2 core.TypedData
	err = json.Unmarshal([]byte(data), &obj2)
	if err != nil {
		t.Fatal(err)
	}
	_, signature, err = sign.EIP712Signature(obj2, privateKey)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("sig: ", common.Bytes2Hex(signature))
	ok, err = sign.VerifyEIP712Signature(obj2, signature, address)
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Fatal(fmt.Errorf("verify failed"))
	}

	fmt.Println("----------- DAS withdraw -------------")
	withdrawStr := `{
    "types": {
        "EIP712Domain": [
            {
                "name": "chainId",
                "type": "uint256"
            },
            {
                "name": "name",
                "type": "string"
            },
            {
                "name": "verifyingContract",
                "type": "address"
            },
            {
                "name": "version",
                "type": "string"
            }
        ],
        "Action": [
            {
                "name": "action",
                "type": "string"
            },
            {
                "name": "params",
                "type": "string"
            }
        ],
        "Cell": [
            {
                "name": "capacity",
                "type": "string"
            },
            {
                "name": "lock",
                "type": "string"
            },
            {
                "name": "type",
                "type": "string"
            },
            {
                "name": "data",
                "type": "string"
            },
            {
                "name": "extraData",
                "type": "string"
            }
        ],
        "Transaction": [
            {
                "name": "DAS_MESSAGE",
                "type": "string"
            },
            {
                "name": "inputsCapacity",
                "type": "string"
            },
            {
                "name": "outputsCapacity",
                "type": "string"
            },
            {
                "name": "fee",
                "type": "string"
            },
            {
                "name": "action",
                "type": "Action"
            },
            {
                "name": "inputs",
                "type": "Cell[]"
            },
            {
                "name": "outputs",
                "type": "Cell[]"
            },
            {
                "name": "digest",
                "type": "bytes32"
            }
        ]
    },
    "primaryType": "Transaction",
    "domain": {
        "chainId": "5",
        "name": "da.systems",
        "verifyingContract": "0x0000000000000000000000000000000020210722",
        "version": "1"
    },
    "message": {
        "DAS_MESSAGE": "TRANSFER FROM 0xc9f53b1d85356b60453f867610888d89a0b667ad(1001.9998 CKB) TO 0xc9f53b1d85356b60453f867610888d89a0b667ad(1001.9997 CKB)",
        "inputsCapacity": "1001.9998 CKB",
        "outputsCapacity": "1001.9997 CKB",
        "fee": "0.0001 CKB",
        "digest": "0x2277c4591b9fdf7403289bbaa9a8d43dc0e9cf9ecf46e416a057bc594a899dcb",
        "action": {
            "action": "withdraw_from_wallet",
            "params": "0x00"
        },
        "inputs": [],
        "outputs": []
    }
}`
	var obj3 core.TypedData
	_ = json.Unmarshal([]byte(withdrawStr), &obj3)
	_, signature, err = sign.EIP712Signature(obj3, privateKey)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("sig: ", common.Bytes2Hex(signature))
	ok, err = sign.VerifyEIP712Signature(obj3, signature, address)
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Fatal(fmt.Errorf("verify failed"))
	}
}
