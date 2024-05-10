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
			CodeHash: types.HexToHash("0xf1ef61b6977508d9ec56fe43399a01e576086a76cf0f7c687d1418335e8c401f"),
			HashType: types.HashTypeType,
			Args:     common.Hex2Bytes("0x01"),
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
		Account:     "20240509.bit",
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

func TestGetDidCellOccupiedCapacity2(t *testing.T) {
	dc, _ := getNewDasCoreTestnet2()

	anyLock := types.Script{
		CodeHash: types.HexToHash("0xf1ef61b6977508d9ec56fe43399a01e576086a76cf0f7c687d1418335e8c401f"),
		HashType: types.HashTypeType,
		Args:     common.Hex2Bytes("0x01"),
	}
	fmt.Println(dc.GetDidCellOccupiedCapacity(&anyLock, "20240509.bit"))
}
