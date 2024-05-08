package example

import (
	"fmt"
	"github.com/dotbitHQ/das-lib/common"
	"github.com/nervosnetwork/ckb-sdk-go/types"
	"testing"
)

func TestCapacity2Str(t *testing.T) {
	fmt.Println(common.Capacity2Str(210000000000))

	fmt.Println(common.GetAccountLength("unit.bit"))
}

func TestConvertAddressKey(t *testing.T) {
	fmt.Println(common.ConvertRecordsAddressCoinType("address.0"))
	fmt.Println(common.ConvertRecordsAddressKey("address.btc"))
}

func TestGetMaxHashLenScriptForNormalCell(t *testing.T) {
	str := common.GetMaxHashLenScriptForNormalCell(&types.Script{
		CodeHash: types.HexToHash("0x0b1f412fbae26853ff7d082d422c2bdd9e2ff94ee8aaec11240a5b34cc6e890f"),
		HashType: "type",
		Args:     common.Hex2Bytes("0x"),
	})
	fmt.Println(str)
}
