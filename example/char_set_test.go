package example

import (
	"fmt"
	"github.com/dotbitHQ/das-lib/common"
	"testing"
)

func TestGetAccountCharType(t *testing.T) {
	var res = make(map[common.AccountCharType]struct{})
	var list []common.AccountCharSet
	list = append(list, common.AccountCharSet{
		CharSetName: common.AccountCharTypeEmoji,
		Char:        "",
	})
	common.GetAccountCharType(res, list)
	fmt.Println(res)

	list = append(list, common.AccountCharSet{
		CharSetName: common.AccountCharTypeEn,
		Char:        "",
	})
	common.GetAccountCharType(res, list)
	fmt.Println(res)
}

func TestGetDotBitAccountLength(t *testing.T) {
	acc := "ให้😊บริก1-าร.ให้😊บริก1-าร.bit"
	fmt.Println(common.GetDotBitAccountLength(acc))
}

func TestAccountLength(t *testing.T) {
	//acc := "ให้😊บริก1-าร.bit"
	////fmt.Println(common.GetDotBitAccountLength(acc))
	//
	//reg:=regexp.MustCompile("[\u0E00-\u0E7F][\u0E31\u0E33-\u0E3A\u0E47-\u0E4E]*")
	//res:=reg.FindAllStringIndex(acc,-1)
	//fmt.Println(len(acc))
	//for _,v:=range res{
	//	fmt.Println(acc[v[0]:v[1]])
	//}

	acc := "ให้😊บริก1-ารฬั่.bit"
	dc, err := getNewDasCoreTestnet2()
	if err != nil {
		t.Fatal(err)
	}

	list, err := dc.GetAccountCharSetList(acc)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(list)

	//acc := ""
	//accLen := 0
	//
	//tmpMap := common.CharSetTypeEmojiMap
	//accLen += len(tmpMap)
	//for k, _ := range tmpMap {
	//	//fmt.Println(k, len(k), []byte(k),utf8.RuneCountInString(k))
	//	acc += k
	//}
	//tmpMap = common.CharSetTypeThMap
	//accLen += len(tmpMap)
	//for k, _ := range tmpMap {
	//	//fmt.Println(k, len(k), []byte(k),utf8.RuneCountInString(k))
	//	acc += k
	//}
	//fmt.Println(accLen)
	//fmt.Println(common.GetDotBitAccountLength(acc + ".bit"))

}
