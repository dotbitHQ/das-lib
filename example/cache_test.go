package example

import (
	"context"
	"fmt"
	"github.com/dotbitHQ/das-lib/core"
	"github.com/dotbitHQ/das-lib/dascache"
	"sync"
	"testing"
	"time"
)

func TestDasCache(t *testing.T) {
	var wg sync.WaitGroup
	dc := dascache.NewDasCache(context.Background(), &wg)
	dc.RunClearExpiredOutPoint(time.Minute * 3)

	go func() {
		for i := 0; i < 10; i++ {
			dc.AddOutPoint([]string{fmt.Sprintf("0x-%d", i)})
			time.Sleep(time.Second)
		}
	}()
	wg.Wait()
}

func TestCache(t *testing.T) {
	dc, err := getNewDasCoreMainNet()
	if err != nil {
		t.Fatal(err)
	}
	//fmt.Println(common.CharSetTypeEmojiMap)
	dc.RunSetConfigCellByCache([]core.CacheConfigCellKey{
		core.CacheConfigCellKeyCharSet,
		core.CacheConfigCellKeyReservedAccounts,
		core.CacheConfigCellKeyBase,
	})

	//str, err := dc.GetConfigCellByCache(core.CacheConfigCellKeyReservedAccounts)
	//if err != nil {
	//	t.Fatal(err)
	//}
	//fmt.Println(str)

	str, err := dc.GetConfigCellByCache(core.CacheConfigCellKeyBase)
	if err != nil {
		t.Log(err)
	}
	fmt.Println(str)
	select {}
}
