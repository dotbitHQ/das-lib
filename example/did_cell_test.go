package example

import (
	"fmt"
	"github.com/dotbitHQ/das-lib/common"
	"github.com/dotbitHQ/das-lib/molecule"
	"github.com/dotbitHQ/das-lib/witness"
	"github.com/nervosnetwork/ckb-sdk-go/types"
	"testing"
)

func TestGetDidCellOccupiedCapacity(t *testing.T) {
	didCell := types.CellOutput{
		Capacity: 0,
		Lock: &types.Script{
			CodeHash: types.HexToHash("0xf329effd1c475a2978453c8600e1eaf0bc2087ee093c3ee64cc96ec6847752cb"),
			HashType: types.HashTypeType,
			Args:     common.Hex2Bytes("0x0115a33588908cf8edb27d1abe3852bf287abd389100"),
		},
		Type: &types.Script{
			CodeHash: types.HexToHash("0x0b1f412fbae26853ff7d082d422c2bdd9e2ff94ee8aaec11240a5b34cc6e890f"),
			HashType: types.HashTypeType,
			Args:     nil,
		},
	}

	defaultWitnessHash := molecule.Byte20Default()
	didCellData := witness.DidCellData{
		ItemId:      witness.ItemIdDidCellDataV0,
		Account:     "20240507.bit",
		ExpireAt:    0,
		WitnessHash: common.Bytes2Hex(defaultWitnessHash.RawData()),
	}
	didCellDataBys, err := didCellData.ObjToBys()
	if err != nil {
		t.Fatal(err)
	}

	didCellCapacity := didCell.OccupiedCapacity(didCellDataBys)
	fmt.Println(didCellCapacity)
}
