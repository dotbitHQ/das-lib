package example

import (
	"fmt"
	"github.com/dotbitHQ/das-lib/common"
	"github.com/dotbitHQ/das-lib/smt"
	"log"
	"testing"
)

func getTree(smtName string) *smt.SmtServer {
	tree := smt.NewSmtSrv("http://127.0.0.1:10000", smtName)
	return tree
}

func TestRsTreeGetRoot(t *testing.T) {
	tree := getTree("tree1")
	r, err := tree.GetSmtRoot()
	if err != nil {
		t.Error(err)
	}
	t.Log(r)
}

func TestRsTreeUpdate(t *testing.T) {
	tree := getTree("tree1")
	var kvTemp []smt.SmtKv
	kvTemp = append(kvTemp, smt.SmtKv{
		Key:   common.Hex2Bytes("0200000000000000000000000000000000000000000000000000000000000000"),
		Value: common.Hex2Bytes("22ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"),
	})
	kvTemp = append(kvTemp, smt.SmtKv{
		Key:   common.Hex2Bytes("0300000000000000000000000000000000000000000000000000000000000000"),
		Value: common.Hex2Bytes("33ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"),
	})
	opt := smt.SmtOpt{GetRoot: true, GetProof: true}
	r, err := tree.UpdateSmt(kvTemp, opt)
	if err != nil {
		t.Error(err)
	}
	t.Log(r)
}

func TestRsTreeUpdateMiddle(t *testing.T) {
	tree := getTree("tree66")
	var kvTemp []smt.SmtKv
	kvTemp = append(kvTemp, smt.SmtKv{
		Key:   common.Hex2Bytes("0000000000000000000000000000000000000000000000000000000000000000"),
		Value: common.Hex2Bytes("00ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"),
	}, smt.SmtKv{
		Key:   common.Hex2Bytes("0100000000000000000000000000000000000000000000000000000000000000"),
		Value: common.Hex2Bytes("11ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"),
	}, smt.SmtKv{
		Key:   common.Hex2Bytes("0200000000000000000000000000000000000000000000000000000000000000"),
		Value: common.Hex2Bytes("22ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"),
	}, smt.SmtKv{
		Key:   common.Hex2Bytes("0300000000000000000000000000000000000000000000000000000000000000"),
		Value: common.Hex2Bytes("33ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"),
	})
	opt := smt.SmtOpt{GetProof: true, GetRoot: true}
	r, err := tree.UpdateMiddleSmt(kvTemp, opt)
	if err != nil {
		t.Error(err)
	}
	t.Log(r)
}

func TestRsTreeDelete(t *testing.T) {
	tree := getTree("tree66")
	r, err := tree.DeleteSmt()
	if err != nil {
		t.Error(err)
	}
	t.Log(r)
}

func TestRsTree(t *testing.T) {
	tree := getTree("test")
	//5000 6s
	//8000 9s
	count := 10000
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
}
