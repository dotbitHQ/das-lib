package example

import (
	"fmt"
	"github.com/dotbitHQ/das-lib/smt"
	"log"
	"testing"
	"time"
)

func TestRsTree(t *testing.T) {
	fmt.Println(time.Now().String())
	tree := smt.NewSmtSrv("http://localhost:10000", "")
	count := 2
	opt := smt.SmtOpt{
		GetRoot:  true,
		GetProof: true,
	}
	var kvTemp []smt.SmtKv
	for i := 0; i < count; i++ {
		key := fmt.Sprintf("key-%d", i)
		value := fmt.Sprintf("value-%d", i)
		k := smt.Sha256(key)
		v := smt.Sha256(value)
		kvTemp = append(kvTemp, smt.SmtKv{
			k,
			v,
		})
	}
	res, err := tree.UpdateSmt(kvTemp, opt)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("root", res.Root)
	for i, _ := range res.Proofs {
		fmt.Println(i, "proof: ", res.Proofs[i])
	}
	fmt.Println(time.Now().String())
}
