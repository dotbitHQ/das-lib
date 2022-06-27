package example

import (
	"context"
	"fmt"
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
