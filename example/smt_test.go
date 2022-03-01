package example

import (
	"fmt"
	"github.com/DeAccountSystems/das-lib/common"
	"github.com/DeAccountSystems/das-lib/smt"
	"github.com/nervosnetwork/ckb-sdk-go/crypto/blake2b"
	"testing"
	"time"
)

func TestSparseMerkleTree(t *testing.T) {
	tree := smt.NewSparseMerkleTree(nil)
	key := common.Hex2Bytes("0000000000000000000000000000000000000000000000000000000000000000")
	value := common.Hex2Bytes("00ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff")
	_ = tree.Update(key, value)
	fmt.Println("root:", tree.Root.String())

	key = common.Hex2Bytes("0100000000000000000000000000000000000000000000000000000000000000")
	value = common.Hex2Bytes("11ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff")
	_ = tree.Update(key, value)
	fmt.Println("root:", tree.Root.String())

	key = common.Hex2Bytes("0200000000000000000000000000000000000000000000000000000000000000")
	value = common.Hex2Bytes("22ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff")
	_ = tree.Update(key, value)
	fmt.Println("root:", tree.Root.String())

	key = common.Hex2Bytes("0300000000000000000000000000000000000000000000000000000000000000")
	value = common.Hex2Bytes("33ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff")
	_ = tree.Update(key, value)
	fmt.Println("root:", tree.Root.String())
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
	fmt.Println("root:", tree.Root.String())

	var keys, values []smt.H256
	keys = append(keys, common.Hex2Bytes("0400000000000000000000000000000000000000000000000000000000000000"))
	values = append(values, smt.H256Zero())
	proof, err := tree.MerkleProof(keys, values)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("proof:", proof.String())
	fmt.Println(smt.Verify(tree.Root, proof, keys, values))
}

func TestMerge(t *testing.T) {
	nodeKey := common.Hex2Bytes("0x00000000000000000000000000000000000000000000000000000000000000d0")
	lhs := smt.MergeValueZero{
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
	res := smt.Merge(255, nodeKey, &lhs, rhs)
	fmt.Println(res.String())
}

func TestMerkleProof2(t *testing.T) {
	tree := smt.NewSparseMerkleTree(nil)
	key1 := common.Hex2Bytes("0x88e6966ee9d691a6befe0664bb54c7b45bbda274a7cf5fa8cd07f56d94741223")
	value1 := common.Hex2Bytes("0xe19e9083ca4dbbee50e56c9825eed7fd750c1982f86412275c9efedf3440f83b")
	_ = tree.Update(key1, value1)
	fmt.Println("root:", tree.Root, tree.Root.String())

	key2 := common.Hex2Bytes("0x26938181394b731558f2bcc40926fc1e38c3036f249319cf7ba845a4bbd769d3")
	value2 := common.Hex2Bytes("0x358305052f73809142a8b4c11f7becbef9d15ac718f86e0e7d51a7f1f2383718")
	_ = tree.Update(key2, value2)
	fmt.Println("root:", tree.Root, tree.Root.String())

	key := common.Hex2Bytes("0x3b66d16df3f793044f09494c7d3fd540be1a94a3ed4c8a686c595cf144703e64")
	value := common.Hex2Bytes("0xb0aa5768b4893807d97ca1785b352fb8ec8f9b5521f549b0a285578b8a57ea97")
	_ = tree.Update(key, value)
	fmt.Println("root:", tree.Root, tree.Root.String())

	key = common.Hex2Bytes("0x95377d6ba3f39fbfdd93f2fe7bb29ff1a52aa4baf0d1ae86d73f7ac9f5de31df")
	value = common.Hex2Bytes("0x61e0ef0afc4eaabbe68ee97bd54e629aa2914085ab3676f0f6df63eb233cc07b")
	_ = tree.Update(key, value)
	fmt.Println("root:", tree.Root, tree.Root.String())

	var keys, values []smt.H256
	k := smt.H256Zero()
	v := smt.H256Zero()
	k[0] = '1'
	keys = append(keys, k)
	values = append(values, v)
	//keys = appen(keys, key2)
	//values = append(values, value2)

	proof, err := tree.MerkleProof(keys, values)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(proof)
	fmt.Println(smt.Verify(tree.Root, proof, keys, values))
}

func TestSmt(t *testing.T) {
	fmt.Println(time.Now().String())
	tree := smt.NewSparseMerkleTree(nil)
	count := 4
	for i := 0; i < count; i++ {
		key := fmt.Sprintf("key-%d", i)
		value := fmt.Sprintf("value-%d", i)
		k, _ := blake2b.Blake256([]byte(key))
		v, _ := blake2b.Blake256([]byte(value))
		//fmt.Println("k:",common.Bytes2Hex(k))
		//fmt.Println("v:",common.Bytes2Hex(v))
		_ = tree.Update(k, v)
	}
	fmt.Println(time.Now().String())
	//for i:=0;i<count;i++{
	//	key := fmt.Sprintf("key-%d", i)
	//	value := fmt.Sprintf("value-%d", i)
	//	var keys, values []smt.H256
	//	k1, _ := blake2b.Blake256([]byte(key))
	//	keys = append(keys, k1)
	//	v1, _ := blake2b.Blake256([]byte(value))
	//	values = append(values, v1)
	//	proof, err := tree.MerkleProof(keys, values)
	//	if err != nil {
	//		t.Fatal(err)
	//	}
	//	fmt.Println(smt.Verify(tree.Root, proof, keys, values))
	//}
	//fmt.Println(time.Now().String())

}

func TestByte(t *testing.T) {
	a := smt.MergeValueZero{
		BaseNode:  []byte("111"),
		ZeroBits:  nil,
		ZeroCount: 0,
	}
	fun := func(value smt.MergeValue) {
		if z, ok := (value).(*smt.MergeValueZero); ok {
			b := z.BaseNode
			b = []byte("222")
			fmt.Println(b, z.BaseNode)
		}
	}
	fun(&a)
}

func TestMergeWithZero(t *testing.T) {
	height := byte(255)
	nodeKey := smt.H256Zero()
	value := smt.MergeValueZero{
		BaseNode:  smt.H256Zero(),
		ZeroBits:  smt.H256Zero(),
		ZeroCount: 1,
	}

	smt.MergeWithZero(height, nodeKey, &value, true)

	a := []byte{'0'} //smt.H256Zero()
	b := []byte{'0'} //smt.H256Zero()
	fmt.Printf("%p-%p\n", a, b)
	b = a
	fmt.Printf("%p-%p\n", a, b)
}
