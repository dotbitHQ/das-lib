package witness

import (
	common2 "github.com/dotbitHQ/das-lib/common"
	"testing"
)

func TestSubAccountNewBuilder_ConvertSubAccountCellOutputData(t *testing.T) {
	data := common2.Hex2Bytes("0xc545e0beadaf1d9c988790585f0b2fe2896727e22cbecfca6febe28247d7f7e400ab9041000000000000000000000000ff0100000000000000000000")
	dataDetail := ConvertSubAccountCellOutputData(data)
	t.Log(dataDetail)
}
