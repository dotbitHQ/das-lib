package example

import (
	"crypto/ed25519"
	"crypto/rand"
	"fmt"
	"github.com/DeAccountSystems/das-lib/common"
	"github.com/DeAccountSystems/das-lib/core"
	"testing"
)

func TestMixin(t *testing.T) {
	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		println(err)
		return
	}
	fmt.Println(common.Bytes2Hex(privateKey))
	fmt.Println(common.Bytes2Hex(publicKey), len(common.Bytes2Hex(publicKey)), len(publicKey))
	// 0x70c756ecfa897dc71cbfce48931fbb261f2e593fe234902a57f36aa9c27dec17 66 32
	args1 := "0x0515a33588908cF8Edb27D1AbE3852Bf287Abd38910515a33588908cF8Edb27D1AbE3852Bf287Abd3891"
	args2 := "0x0670c756ecfa897dc71cbfce48931fbb261f2e593fe234902a57f36aa9c27dec170670c756ecfa897dc71cbfce48931fbb261f2e593fe234902a57f36aa9c27dec17"
	fmt.Println(len(args1), len(args2))
	fmt.Println(len(common.Hex2Bytes(args1)), len(common.Hex2Bytes(args2)))
}

func TestFormatDasLockToOwnerAndManager(t *testing.T) {
	//args := "0x0515a33588908cF8Edb27D1AbE3852Bf287Abd38910515a33588908cF8Edb27D1AbE3852Bf287Abd3891"
	//args := "0x0670c756ecfa897dc71cbfce48931fbb261f2e593fe234902a57f36aa9c27dec170670c756ecfa897dc71cbfce48931fbb261f2e593fe234902a57f36aa9c27dec17"
	//args:="0x0670c756ecfa897dc71cbfce48931fbb261f2e593fe234902a57f36aa9c27dec170515a33588908cF8Edb27D1AbE3852Bf287Abd3891"
	args := "0x0515a33588908cf8edb27d1abe3852bf287abd38910670c756ecfa897dc71cbfce48931fbb261f2e593fe234902a57f36aa9c27dec17"
	owner, manager := core.FormatDasLockToOwnerAndManager(common.Hex2Bytes(args))
	fmt.Println(common.Bytes2Hex(owner), common.Bytes2Hex(manager))
}

func TestFormatAddressToHex(t *testing.T) {
	fmt.Println(core.FormatAddressToHex(common.ChainTypeCkb, "ckb1qjfhdsa4syv599s2s3nfrctwga70g0tu07n9gpnun9ydlngf5vsnwqeerx5wkcvuet3jlw5g6vec9xffmvh5xfqr8yv636mpnn9wxta63rfn8q5e98dj7sey04pzkp"))
}

func TestFormatOwnerManagerAddressToArgs2(t *testing.T) {
	oCT, mCT, oA, mA := common.ChainTypeEth, common.ChainTypeTron, "0xc9f53b1d85356B60453F867610888D89a0B667Ad", "TEooRfPxhqJ7AJfmsRg5hZWEX95VeNxvtX"
	args := core.FormatOwnerManagerAddressToArgs(oCT, mCT, oA, mA)
	fmt.Println(common.Bytes2Hex(args))
	oCT, oA = common.ChainTypeMixin, "0x70c756ecfa897dc71cbfce48931fbb261f2e593fe234902a57f36aa9c27dec17"
	args = core.FormatOwnerManagerAddressToArgs(oCT, mCT, oA, mA)
	fmt.Println(common.Bytes2Hex(args))
	mCT, mA = common.ChainTypeMixin, "0x70c756ecfa897dc71cbfce48931fbb261f2e593fe234902a57f36aa9c27dec17"
	args = core.FormatOwnerManagerAddressToArgs(oCT, mCT, oA, mA)
	fmt.Println(common.Bytes2Hex(args))
}

func TestFormatDasLockToHexAddress(t *testing.T) {
	args := "0x053919a8eb619ccae32fba88d333829929db2f432405c9f53b1d85356b60453f867610888d89a0b667ad"
	fmt.Println(core.FormatDasLockToHexAddress(common.Hex2Bytes(args)))
	args = "0x0670c756ecfa897dc71cbfce48931fbb261f2e593fe234902a57f36aa9c27dec170515a33588908cF8Edb27D1AbE3852Bf287Abd3891"
	fmt.Println(core.FormatDasLockToHexAddress(common.Hex2Bytes(args)))
}
