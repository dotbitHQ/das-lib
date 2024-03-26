package example

import (
	"fmt"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/dotbitHQ/das-lib/bitcoin"
	"github.com/dotbitHQ/das-lib/common"
	"testing"
)

func TestCreateBTCWallet(t *testing.T) {
	err := bitcoin.CreateBTCWallet(bitcoin.BtcAddressTypeP2SHP2WPKH, true)
	if err != nil {
		t.Fatal(err)
	}
	//WIF: L2vKWmpxVFsRCQPxnhvjsLiYB3hTSV85fAm1Jo6CcAJkvgKqjxoh
	//PubKey: 147VZrBkaWy5zJhpuGAa7EZ9B9YBLu8MuM
	//PubHash 2222b81757f47ebe58881573f64fb8c5f59ba533
	//PriKey: aa13ee7c615ef80c9063bf6875fb894b3936c9551d73bfe0361a4682ae7efe8f

	//WIF: L3t7wxUjYs5A11kajfdQy2w1CnTKCbSxYFMMgstuYX7QraQt7nwb
	//ScriptAddr: 35Y6PCZk4zuP1GJkjrqqR7PpvgWbiMVuvx
	//PubHash d6c09590c8515eaaae150871b19a11cb44c54771
	//pkScript: 76a914d6c09590c8515eaaae150871b19a11cb44c5477188ac
	//pkScriptHash: 2a307b6ee071be7d8f484f1f0c06369742e46919
	//PriKey: c6c8a6bf98b562089e93e5f5270ea4468f3a442a88cccfcc74692bad458c32d3

	//WIF: KwVZNWG6fyqSh1uhVM25iNgNL89wxdbZcr3M5dnTtqdq4T4ZQfBt
	//PubKey: bc1q88cy67dd4q2aag30ezhlrt93wwvpapsruefmrf
	//PubHash 39f04d79ada815dea22fc8aff1acb173981e8603
	//PriKey: 082720675b373fbaa6c24fb099867dfbbdeba98ab3c7c83c9ecb2ea26b5fa97d

	//WIF: KyMDvdf11J1CydwBNuMQ6uYVJXbV93j2FCi5ts2XZbVRPm7PeVvZ
	//ScriptAddr: 3A3basSqtJZPdA9mKCC1KtQkgXjKSSJnWc
	//PubHash 3a6274d504078fd35d21aff131eb22c7b1af13ef
	//pkScript: 00143a6274d504078fd35d21aff131eb22c7b1af13ef
	//pkScriptHash: 5ba56c93f710da685871a01afd2e47da5ca069b2
	//PriKey: 3f8a2671be95d5301e0bd7239a87ed9bb357e71545b3e8efbe89dfb1e932fdce
}

func TestDecodeAddr(t *testing.T) {
	addrStr := "35qVLYDmdnh8hC8VEJMqPPmqQ3S5K9ya5U"
	p := bitcoin.GetBTCMainNetParams()
	addr, err := btcutil.DecodeAddress(addrStr, &p)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(addr.String(), addr.EncodeAddress(), common.Bytes2Hex(addr.ScriptAddress()))
}
