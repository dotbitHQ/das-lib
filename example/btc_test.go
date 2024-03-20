package example

import (
	"fmt"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/dotbitHQ/das-lib/bitcoin"
	"github.com/dotbitHQ/das-lib/common"
	"testing"
)

func TestCreateBTCWallet(t *testing.T) {
	err := bitcoin.CreateBTCWallet(bitcoin.BtcAddressTypeP2WPKH, true)
	if err != nil {
		t.Fatal(err)
	}
	//WIF: L2vKWmpxVFsRCQPxnhvjsLiYB3hTSV85fAm1Jo6CcAJkvgKqjxoh
	//PubKey: 147VZrBkaWy5zJhpuGAa7EZ9B9YBLu8MuM
	//PubHash 2222b81757f47ebe58881573f64fb8c5f59ba533
	//PriKey: aa13ee7c615ef80c9063bf6875fb894b3936c9551d73bfe0361a4682ae7efe8f

	//WIF: KwVZNWG6fyqSh1uhVM25iNgNL89wxdbZcr3M5dnTtqdq4T4ZQfBt
	//PubKey: bc1q88cy67dd4q2aag30ezhlrt93wwvpapsruefmrf
	//PubHash 39f04d79ada815dea22fc8aff1acb173981e8603
	//PriKey: 082720675b373fbaa6c24fb099867dfbbdeba98ab3c7c83c9ecb2ea26b5fa97d
}

func TestDecodeAddr(t *testing.T) {
	addrStr := "bc1qar0srrr7xfkvy5l643lydnw9re59gtzzwf5mdq"
	p := bitcoin.GetBTCMainNetParams()
	addr, err := btcutil.DecodeAddress(addrStr, &p)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(addr.String(), addr.EncodeAddress(), common.Bytes2Hex(addr.ScriptAddress()))
}
