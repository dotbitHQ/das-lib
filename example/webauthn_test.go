package example

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"das_database/dao"
	"encoding/hex"
	"fmt"
	"github.com/dotbitHQ/das-lib/common"
	"github.com/dotbitHQ/das-lib/core"
	"github.com/dotbitHQ/das-lib/witness"
	"github.com/nervosnetwork/ckb-sdk-go/types"
	"math/big"
	"testing"
)

func TestCaculateCid1(t *testing.T) {
	daf := core.DasAddressFormat{DasNetType: common.DasNetTypeTestnet2}
	dc, err := getNewDasCoreTestnet2()
	if err != nil {
		t.Fatal(err)
	}
	hash := "0x175b2ef3fd663fc3cfcc64da60814acafe02f2de47b538574a9391ce895b050d"
	res, err := dc.Client().GetTransaction(context.Background(), types.HexToHash(hash))
	if err != nil {
		fmt.Println(err)
		return
	}

	builder, err := witness.WebAuthnKeyListDataBuilderFromTx(res.Transaction, common.DataTypeNew)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("args: ", common.Bytes2Hex(res.Transaction.Outputs[0].Lock.Args))
	ownerHex, _, err := daf.ArgsToHex(res.Transaction.Outputs[0].Lock.Args)
	if err != nil {
		t.Error(err)
	}
	var masterCidPk1 dao.TableCidPk

	masterCidPk1.Cid = common.Bytes2Hex(ownerHex.AddressPayload[:10])
	masterCidPk1.Pk = common.Bytes2Hex(ownerHex.AddressPayload[10:])
	masterCidPk1.Outpoint = common.OutPoint2String(res.Transaction.Hash.Hex(), 0)
	webauthnSignLv, err := witness.GetWebAuthnSignLvByWitness0(res.Transaction.Witnesses[0])
	if err != nil {
		fmt.Println(err)
		return
	}
	if webauthnSignLv.PkIndex == 255 {
		masterCidPk1.OriginPk = webauthnSignLv.PubKey
	}
	fmt.Println("master: ", masterCidPk1)
	var pubKey ecdsa.PublicKey
	pubKey.Curve = elliptic.P256()
	pubKey.X = new(big.Int).SetBytes(common.Hex2Bytes(webauthnSignLv.PubKey)[:32])
	pubKey.Y = new(big.Int).SetBytes(common.Hex2Bytes(webauthnSignLv.PubKey)[32:])
	signAddrPk1 := common.CaculatePk1(&pubKey)
	keyList := witness.ConvertToWebauthnKeyList(builder.DeviceKeyListCellData.Keys())
	//var authorize []dao.TableAuthorize
	//更新master的cid1, pk1, originPk
	//非master的cidpk
	var slaveCidPks []dao.TableCidPk
	var authorize []dao.TableAuthorize

	for i := 0; i < len(keyList); i++ {
		var slaveCidPk dao.TableCidPk
		cid1 := keyList[i].Cid
		pk1 := keyList[i].PubKey
		//非master
		fmt.Println("key: ", keyList[i])
		if cid1 != masterCidPk1.Cid {
			slaveCidPk.Cid = keyList[i].Cid
			slaveCidPk.Pk = keyList[i].PubKey
			//slave签的名
			if webauthnSignLv.PkIndex != 255 && common.Bytes2Hex(signAddrPk1) == pk1 {
				slaveCidPk.OriginPk = webauthnSignLv.PubKey
			}
			slaveCidPks = append(slaveCidPks, slaveCidPk)
		}
		authorize = append(authorize, dao.TableAuthorize{
			MasterAlgId:    common.DasAlgorithmIdWebauthn,
			MasterSubAlgId: common.DasAlgorithmId(7),
			MasterCid:      masterCidPk1.Cid,
			MasterPk:       masterCidPk1.Pk,
			SlaveAlgId:     common.DasAlgorithmId(keyList[i].MinAlgId),
			SlaveSubAlgId:  common.DasAlgorithmId(keyList[i].SubAlgId),
			SlaveCid:       keyList[i].Cid,
			SlavePk:        keyList[i].PubKey,
			Outpoint:       common.OutPoint2String(res.Transaction.Hash.Hex(), 0),
		})
	}
	fmt.Println("master: ", masterCidPk1)
	fmt.Println("slave: ", slaveCidPks)
	fmt.Println("authorize: ", authorize)
}

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
	address := "ckt1qqexmutxu0c2jq9q4msy8cc6fh4q7q02xvr7dc347zw3ks3qka0m6qggqamft065ws39jxjms4tndppk0kszg6u86gyqwa54ha28ggjerfdc24ekssm8mgpydwrayhfhhx2"
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
		return
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

func TestAddPkIndexForSignMsg(t *testing.T) {
	dc, err := getNewDasCoreTestnet2()
	if err != nil {
		t.Fatal(err)
	}
	signMsg := "4002709dae863bdf506131edfc8a38242b5cb126aed6477695e0b7ce7b2408b756501a116029f8f8aa923877a775bc1d0a449ce02fabf8420dcb7c1b0d59b1426f40a6ae565f4a6137a8ed08e33988cbbe24698ea906ec84215ce042e4812c19502f33b03f6bcc027b41f503f2d25de9e346591cbd03aef5ce5826b3151fdc2aec212549960de5880e8c687434170f6476605b8fe4aeb9a28632c7995cf3ba831d97630500000000bf007b2274797065223a22776562617574686e2e676574222c226368616c6c656e6765223a22526e4a7662534175596d6c304f69426d4d6a63794e7a67314e6a51344e7a51354d6a51334d7a45355a6a4e68596a4e6b595755345a4755304f54646d4f44426c4d7a417a4d5459344e7a45794d575a695a54646b4f574a6b4e7a59304e474d344e7a6332222c226f726967696e223a22687474703a2f2f6c6f63616c686f73743a38303031222c2263726f73734f726967696e223a66616c73657d"
	dc.AddPkIndexForSignMsg(&signMsg, 1)
	fmt.Println(signMsg)
}
