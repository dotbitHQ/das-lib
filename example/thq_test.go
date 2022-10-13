package example

import (
	"fmt"
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
