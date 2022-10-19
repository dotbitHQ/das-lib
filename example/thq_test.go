package example

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/dotbitHQ/das-lib/common"
	"github.com/dotbitHQ/das-lib/core"
	"github.com/minio/blake2b-simd"
	"github.com/nervosnetwork/ckb-sdk-go/address"
	"github.com/nervosnetwork/ckb-sdk-go/indexer"
	"strings"
	"testing"
)

func TestTHQ(t *testing.T) {
	dc, err := getNewDasCoreTestnet2()
	if err != nil {
		t.Fatal(err)
	}
	list, err := dc.GetTimeCellList()
	if err != nil {
		t.Fatal(err)
	}
	for i, v := range list {
		fmt.Println(i, v.Timestamp())
	}
}

func TestLuckNumber(t *testing.T) {
	str := `rocket.bit
cheat.bit
success.bit
perfect.bit
school.bit
monolith.bit
flashback.bit
century.bit
detective.bit
professor.bit
doctor.bit
worker.bit
payment.bit
follower.bit
capital.bit
nominate.bit
monetary.bit
millenium.bit
economics.bit
tycoon.bit
warlord.bit
product.bit
shine.bit
glitter.bit
person.bit
trainer.bit
raise.bit
jackpot.bit
bless.bit
cancel.bit
limited.bit
original.bit
emperor.bit
buyback.bit
monarch.bit
warlock.bit
invisible.bit
neymar.bit`
	list := strings.Split(str, "\n")
	dc, err := getNewDasCoreMainNet()
	if err != nil {
		t.Fatal(err)
	}

	configRelease, err := dc.ConfigCellDataBuilderByTypeArgs(common.ConfigCellTypeArgsRelease)
	if err != nil {
		t.Fatal(err)
	}
	luckyNumber, _ := configRelease.LuckyNumber()
	fmt.Println("config release lucky number: ", luckyNumber, len(list))
	for _, v := range list {
		resNum, _ := Blake256AndFourBytesBigEndian([]byte(v))
		fmt.Println(resNum > luckyNumber, v)
	}

	return

}

func Blake256AndFourBytesBigEndian(data []byte) (uint32, error) {
	bys, err := Blake256(data)
	if err != nil {
		return 0, err
	}
	bytesBuffer := bytes.NewBuffer(bys[0:4])
	var res uint32
	if err = binary.Read(bytesBuffer, binary.BigEndian, &res); err != nil {
		return 0, err
	}
	return res, nil
}

func Blake256(data []byte) ([]byte, error) {
	tmpConfig := &blake2b.Config{
		Size:   32,
		Person: []byte("2021-07-22 12:00"),
	}
	hash, err := blake2b.New(tmpConfig)
	if err != nil {
		return nil, err
	}
	hash.Write(data)
	return hash.Sum(nil), nil
}

func TestNormalCell(t *testing.T) {
	dc, err := getNewDasCoreMainNet()
	if err != nil {
		t.Fatal(err)
	}

	parseAddr, err := address.Parse("ckb1qyqxmyfg2a5w0jt0rn4qzu7gzead5t87405qs8cqan")
	if err != nil {
		t.Fatal(err)
	}
	liveCells, total, err := dc.GetBalanceCells(&core.ParamGetBalanceCells{
		DasCache:          nil,
		LockScript:        parseAddr.Script,
		CapacityNeed:      0,
		CapacityForChange: 0,
		SearchOrder:       indexer.SearchOrderDesc,
	})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("total:", total, len(liveCells))
	for _, v := range liveCells {
		fmt.Println(v.OutPoint.TxHash.String(), v.OutPoint.Index, v.Output.Capacity/common.OneCkb)
	}
}
