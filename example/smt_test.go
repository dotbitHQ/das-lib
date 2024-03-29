package example

import (
	"context"
	"fmt"
	"github.com/dotbitHQ/das-lib/common"
	"github.com/dotbitHQ/das-lib/smt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"testing"
	"time"
)

func TestSparseMerkleTree(t *testing.T) {
	tree := smt.NewSparseMerkleTree(nil)
	key := common.Hex2Bytes("0000000000000000000000000000000000000000000000000000000000000000")
	value := common.Hex2Bytes("00ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff")
	_ = tree.Update(key, value)
	fmt.Println(tree.Root())

	key = common.Hex2Bytes("0100000000000000000000000000000000000000000000000000000000000000")
	value = common.Hex2Bytes("11ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff")
	_ = tree.Update(key, value)
	fmt.Println(tree.Root())

	key = common.Hex2Bytes("0200000000000000000000000000000000000000000000000000000000000000")
	value = common.Hex2Bytes("22ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff")
	_ = tree.Update(key, value)
	fmt.Println(tree.Root())

	key = common.Hex2Bytes("0300000000000000000000000000000000000000000000000000000000000000")
	value = common.Hex2Bytes("33ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff")
	_ = tree.Update(key, value)
	fmt.Println(tree.Root())
}

func TestMerkleProof(t *testing.T) {
	tree := smt.NewSparseMerkleTree(nil)
	key := common.Hex2Bytes("0000000000000000000000000000000000000000000000000000000000000000")
	value := common.Hex2Bytes("00ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff")
	_ = tree.Update(key, value)

	key = common.Hex2Bytes("0100000000000000000000000000000000000000000000000000000000000000")
	value = common.Hex2Bytes("11ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff")
	_ = tree.Update(key, value)

	key = common.Hex2Bytes("0300000000000000000000000000000000000000000000000000000000000000")
	value = common.Hex2Bytes("33ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff")
	_ = tree.Update(key, value)
	fmt.Println(tree.Root())

	var keys, values []smt.H256
	keys = append(keys, common.Hex2Bytes("0400000000000000000000000000000000000000000000000000000000000000"))
	values = append(values, smt.H256Zero())
	proof, err := tree.MerkleProof(keys, values)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("proof:", proof.String())
	root, err := tree.Root()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(smt.Verify(root, proof, keys, values))
}

func TestMerge(t *testing.T) {
	nodeKey := common.Hex2Bytes("0x00000000000000000000000000000000000000000000000000000000000000d0")
	lhs := smt.MergeValue{
		Value:     nil,
		BaseNode:  common.Hex2Bytes("0x9180ed6242e737f554d3c4f7b8f8f810581d810bcf3c1075070b45a6104d5ff8"),
		ZeroBits:  common.Hex2Bytes("0x26938181394b731558f2bcc40926fc1e38c3036f249319cf7ba845a4bbd76903"),
		ZeroCount: 251,
	}
	rhs := smt.MergeValueFromZero()
	//rhs := smt.MergeValueZero{
	//	BaseNode:  common.Hex2Bytes("0x9180ed6242e737f554d3c4f7b8f8f810581d810bcf3c1075070b45a6104d5ff8"),
	//	ZeroBits:  common.Hex2Bytes("0x26938181394b731558f2bcc40926fc1e38c3036f249319cf7ba845a4bbd76953"),
	//	ZeroCount: 255,
	//}
	res := smt.Merge(255, nodeKey, lhs, rhs)
	fmt.Println(res.String())
}

func TestMerkleProof2(t *testing.T) {
	tree := smt.NewSparseMerkleTree(nil)
	key1 := common.Hex2Bytes("0x88e6966ee9d691a6befe0664bb54c7b45bbda274a7cf5fa8cd07f56d94741223")
	value1 := common.Hex2Bytes("0xe19e9083ca4dbbee50e56c9825eed7fd750c1982f86412275c9efedf3440f83b")
	_ = tree.Update(key1, value1)
	fmt.Println(tree.Root())

	key2 := common.Hex2Bytes("0x26938181394b731558f2bcc40926fc1e38c3036f249319cf7ba845a4bbd769d3")
	value2 := common.Hex2Bytes("0x358305052f73809142a8b4c11f7becbef9d15ac718f86e0e7d51a7f1f2383718")
	_ = tree.Update(key2, value2)
	fmt.Println(tree.Root())

	key := common.Hex2Bytes("0x3b66d16df3f793044f09494c7d3fd540be1a94a3ed4c8a686c595cf144703e64")
	value := common.Hex2Bytes("0xb0aa5768b4893807d97ca1785b352fb8ec8f9b5521f549b0a285578b8a57ea97")
	_ = tree.Update(key, value)
	fmt.Println(tree.Root())

	key = common.Hex2Bytes("0x95377d6ba3f39fbfdd93f2fe7bb29ff1a52aa4baf0d1ae86d73f7ac9f5de31df")
	value = common.Hex2Bytes("0x61e0ef0afc4eaabbe68ee97bd54e629aa2914085ab3676f0f6df63eb233cc07b")
	_ = tree.Update(key, value)
	fmt.Println(tree.Root())

	var keys, values []smt.H256
	k := smt.H256Zero()
	v := smt.H256Zero()
	k[0] = '1'
	keys = append(keys, k)
	values = append(values, v)

	proof, err := tree.MerkleProof(keys, values)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(proof)
	root, err := tree.Root()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(smt.Verify(root, proof, keys, values))
}

func TestSmt(t *testing.T) {
	// 10000 4s
	// 100000 2min
	fmt.Println(time.Now().String())
	tree := smt.NewSparseMerkleTree(nil)
	count := 1000
	for i := 0; i < count; i++ {
		key := fmt.Sprintf("key-%d", i)
		value := fmt.Sprintf("value-%d", i)
		k := smt.Sha256(key)
		v := smt.Sha256(value)
		_ = tree.Update(k, v)
	}
	fmt.Println(time.Now().String())
	for i := 0; i < count; i++ {
		key := fmt.Sprintf("key-%d", i)
		value := fmt.Sprintf("value-%d", i)
		var keys, values []smt.H256
		k1 := smt.Sha256(key)
		keys = append(keys, k1)
		v1 := smt.Sha256(value)
		values = append(values, v1)
		proof, err := tree.MerkleProof(keys, values)
		if err != nil {
			t.Fatal(err)
		}
		root, err := tree.Root()
		if err != nil {
			t.Fatal(err)
		}
		fmt.Println(smt.Verify(root, proof, keys, values))
	}
	fmt.Println(time.Now().String())

}

func TestMongodbStoreDB(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		t.Fatal(err)
	}
	s := smt.NewMongoStore(ctx, client, "smt", "test")
	if err := s.Collection().Drop(context.Background()); err != nil {
		t.Fatal(err)
	}
	//key := BranchKey{
	//	Height:    0,
	//	NodeKey:   H256Zero(),
	//}
	//node := BranchNode{
	//	Left:  MergeValueFromZero(),
	//	Right: MergeValueFromZero(),
	//}
	//err = s.InsertBranch(key, node)
	//if err != nil {
	//	t.Fatal(err)
	//}
	//res, err := s.GetBranch(key)
	//if err != nil {
	//	t.Fatal(err)
	//}
	//fmt.Println(res)
	//fmt.Println(s.UpdateRoot(H256Zero()))
	//fmt.Println(s.Root())
}

func TestMongodbStore(t *testing.T) {
	// 1000 2min 30M
	// 10000 16min 330M
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		t.Fatal(err)
	}
	s := smt.NewMongoStore(ctx, client, "smt", "test")
	fmt.Println(time.Now().String())
	tree := smt.NewSparseMerkleTree(s)
	count := 100
	//for i := 0; i < count; i++ {
	//	key := fmt.Sprintf("key-%d", i)
	//	value := fmt.Sprintf("value-%d", i)
	//	k := smt.Sha256(key)
	//	v := smt.Sha256(value)
	//	if err := tree.Update(k, v); err != nil {
	//		t.Fatal(err)
	//	}
	//}
	fmt.Println(time.Now().String())

	for i := 0; i < count; i++ {
		key := fmt.Sprintf("key-%d", i)
		value := fmt.Sprintf("value-%d", i)
		var keys, values []smt.H256
		k1 := smt.Sha256(key)
		keys = append(keys, k1)
		v1 := smt.Sha256(value)
		//v1=smt.H256Zero()
		values = append(values, v1)
		proof, err := tree.MerkleProof(keys, values)
		if err != nil {
			t.Fatal(err)
		}
		//fmt.Println("proof:",proof.String())
		root, err := tree.Root()
		if err != nil {
			t.Fatal(err)
		}
		//fmt.Println("root:",root.String())
		fmt.Println(smt.Verify(root, proof, keys, values))
	}
}

func TestDelete(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		t.Fatal(err)
	}
	s := smt.NewMongoStore(ctx, client, "smt", "test")
	tree := smt.NewSparseMerkleTree(s)

	//tree := NewSparseMerkleTree("", nil)
	key := smt.H256Zero()
	value := smt.H256Zero()
	value[0] = '1'
	err = tree.Update(key, value)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(tree.Root())
	value = smt.H256Zero()
	err = tree.Update(key, value)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(tree.Root())
}

func Test2(t *testing.T) {
	k := []byte{91, 188, 110, 15, 128, 104, 148, 171, 83, 108, 156, 199, 122, 145, 190, 146, 84, 150, 206, 158, 161, 39, 206, 84, 182, 81, 148, 9, 34, 103, 13, 123}
	v := []byte{10, 234, 159, 195, 108, 210, 94, 183, 178, 143, 10, 96, 67, 42, 241, 169, 69, 205, 188, 246, 9, 20, 78, 125, 231, 127, 180, 223, 79, 171, 218, 34}
	tree := smt.NewSparseMerkleTree(nil)
	fmt.Println(tree.Update(k, v))
	root, err := tree.Root()
	fmt.Println("root:", root, err)

	k = []byte{35, 39, 151, 104, 126, 195, 63, 86, 168, 88, 152, 230, 236, 75, 146, 78, 65, 193, 18, 38, 122, 141, 162, 51, 42, 58, 32, 61, 19, 233, 11, 88}
	v = []byte{35, 170, 60, 60, 144, 195, 7, 48, 102, 130, 3, 182, 179, 244, 135, 57, 12, 225, 6, 169, 246, 125, 32, 198, 108, 107, 202, 217, 145, 67, 170, 177}
	fmt.Println(tree.Update(k, v))
	root, err = tree.Root()
	fmt.Println("root:", root, err)

	var ks, vs []smt.H256
	ks = append(ks, k)
	vs = append(vs, v)
	proof, err := tree.MerkleProof(ks, vs)
	fmt.Println("proof:", proof, err)
	fmt.Println(common.Hex2Bytes(proof.String()))

	fmt.Println(smt.Verify(root, proof, ks, vs))

}

var SmtRpcUrl = "http://127.0.0.1:10000"

func TestSmtDel(t *testing.T) {
	reverse1w := smt.NewSmtSrv(SmtRpcUrl, "reverse_1w")
	root, err := reverse1w.GetSmtRoot()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(root.String())
	ok, err := reverse1w.DeleteSmt()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("del:", ok)
	root, err = reverse1w.GetSmtRoot()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(root.String())

	// ps -ef |grep rpc
	// top -p pid
	//1w 0.92s 0.72s
	//3w 2.83s 2.22s
	//5w 5.23s 3.77s
	//10w 11.45s 8.56s
	//30w 38.97s 32.37s
	//50w 72.08s 59.43s
	//100w 152.68s 134.93s
}

func TestAddSmt(t *testing.T) {
	reverse1w := smt.NewSmtSrv(SmtRpcUrl, "reverse_1w")
	root, err := reverse1w.GetSmtRoot()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(root.String())

	//
	var list []smt.SmtKv
	for i := 0; i < 1000000; i++ {
		key := smt.H256(smt.Sha256(fmt.Sprintf("key-%d", i)))
		val := smt.H256(smt.Sha256(fmt.Sprintf("val-%d", i)))
		list = append(list, smt.SmtKv{
			Key:   key,
			Value: val,
		})
	}

	base := 5000
	for i, j := 0, base; i < len(list); i, j = j, j+base {
		fmt.Println(i, j, len(list[i:j]))
		_, err = reverse1w.UpdateSmt(list[i:j], smt.SmtOpt{
			GetProof: false,
			GetRoot:  false,
		})
		if err != nil {
			t.Fatal(err)
		}
	}

	//
	root, err = reverse1w.GetSmtRoot()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(root.String())

}
