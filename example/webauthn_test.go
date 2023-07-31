package example

import (
	"context"
	"crypto/elliptic"
	"encoding/hex"
	"fmt"
	"github.com/dotbitHQ/das-lib/common"
	"github.com/dotbitHQ/das-lib/core"
	"github.com/dotbitHQ/das-lib/witness"
	"math/big"
	"testing"
)

func TestEcdsaRecover(t *testing.T) {
	curve := elliptic.P256()
	hash := common.Hex2Bytes("0xc1e6af5868ebf57c58db788df1d8014a3a3ff1990dcb526984acbae05861fd7d")
	R, _ := new(big.Int).SetString("104087844134925704986103407704658370369646035975913171373716846992819522079834", 0)
	S, _ := new(big.Int).SetString("10390831784974169865608269584110793970838780340007159803163726570481115790157", 0)
	possiblePubkey, err := common.GetEcdsaPossiblePubkey(curve, hash, R, S)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("possible 0: ", possiblePubkey[0].X, possiblePubkey[0].Y)
	fmt.Println("possible 1: ", possiblePubkey[1].X, possiblePubkey[1].Y)

	R1, _ := new(big.Int).SetString("42981296685980515483583785598431923091896009544991939339794301360948677609585", 0)
	S1, _ := new(big.Int).SetString("45429925069316295008842370314391818069763181289459716330659064524204951258103", 0)
	possiblePubkey1, err := common.GetEcdsaPossiblePubkey(curve, hash, R1, S1)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("-----")
	fmt.Println("possible 0: ", possiblePubkey1[0].X, possiblePubkey1[0].Y)
	fmt.Println("possible 1: ", possiblePubkey1[1].X, possiblePubkey1[1].Y)

	//realPubkey:  12098966267413439708728706199315115894307800943856814227612321598814731375752 3188317386184053029652564183251176637199913181249076473808524973789124060714
}

func TestGetkeylistCell(t *testing.T) {
	daf := core.DasAddressFormat{DasNetType: common.DasNetTypeTestnet2}
	dc, err := getNewDasCoreTestnet2()
	if err != nil {
		t.Fatal(err)
	}
	address := "ckt1qqexmutxu0c2jq9q4msy8cc6fh4q7q02xvr7dc347zw3ks3qka0m6qggq7w79h22yxg9h5r3vdw79yhka5vqn48t9yyq080zm49zryzm6pckxh0zjtmw6xqf6n4jj9r9323"
	addressHex, _ := daf.NormalToHex(core.DasAddressNormal{
		ChainType:     8,
		AddressNormal: address,
		Is712:         true,
	})
	fmt.Println("addressHex ", addressHex.AddressHex)
	dasLockKey := core.DasAddressHex{
		DasAlgorithmId:    common.DasAlgorithmIdWebauthn,
		DasSubAlgorithmId: common.DasWebauthnSubAlgorithmIdES256,
		AddressHex:        addressHex.AddressHex,
		AddressPayload:    common.Hex2Bytes(addressHex.AddressHex),
		ChainType:         common.ChainTypeWebauthn,
	}
	fmt.Println("dasLockKey ", dasLockKey)
	lockArgs, err := daf.HexToArgs(dasLockKey, dasLockKey)
	fmt.Println("lockArgs ", common.Bytes2Hex(lockArgs))
	cell, err := dc.GetKeyListCell(lockArgs)
	if err != nil {
		fmt.Println("GetKeyListCell(webauthn keyListCell) : ", err.Error())
	}
	if cell != nil {
		fmt.Println(common.OutPoint2String(cell.OutPoint.TxHash.Hex(), 0))
	} else {
		fmt.Println("not found cell")
	}

	keyListConfigTx, err := dc.Client().GetTransaction(context.Background(), cell.OutPoint.TxHash)
	if err != nil {
		fmt.Println(err)
	}
	webAuthnKeyListConfigBuilder, err := witness.WebAuthnKeyListDataBuilderFromTx(keyListConfigTx.Transaction, common.DataTypeNew)
	if err != nil {
		fmt.Println(err)
	}
	dataBuilder := webAuthnKeyListConfigBuilder.DeviceKeyListCellData.AsBuilder()
	deviceKeyListCellDataBuilder := dataBuilder.Build()
	keyList := deviceKeyListCellDataBuilder.Keys()
	for i := 0; i < int(keyList.Len()); i++ {
		mainAlgId := common.DasAlgorithmId(keyList.Get(uint(i)).MainAlgId().RawData()[0])
		subAlgId := common.DasSubAlgorithmId(keyList.Get(uint(i)).SubAlgId().RawData()[0])
		cid1 := keyList.Get(uint(i)).Cid().RawData()
		pk1 := keyList.Get(uint(i)).Pubkey().RawData()
		addressHex := hex.EncodeToString(append(cid1, pk1...))
		fmt.Println(mainAlgId, subAlgId, addressHex)
	}
}
