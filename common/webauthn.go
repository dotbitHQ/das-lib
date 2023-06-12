package common

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
)

func GetWebauthnPayload(cid string, pk *ecdsa.PublicKey) (payload string) {
	cid1 := CaculateCid1(cid)
	pk1 := CaculatePk1(pk)
	payload = CaculateWebauthnPayload(cid1, pk1)
	return
}

func CaculateWebauthnPayload(cid1, pk1 []byte) (payload string) {
	payload = hex.EncodeToString(append(cid1, pk1...))
	return
}

//cid' = hash(cid)*5 [:10]
func CaculateCid1(cid string) (cid1 []byte) {
	hash := sha256.Sum256(common.Hex2Bytes(cid))
	for i := 0; i < 4; i++ {
		hash = sha256.Sum256(hash[:])
	}
	return hash[:10]
}

//pk' = hash(X+Y)*5 [:10]
func CaculatePk1(pk *ecdsa.PublicKey) (cid1 []byte) {
	if pk == nil {
		return
	}
	xy := append(pk.X.Bytes(), pk.Y.Bytes()...)
	hash := sha256.Sum256(xy)
	for i := 0; i < 4; i++ {
		hash = sha256.Sum256(hash[:])
	}
	return hash[:10]
}

//Recover two possible ecdsa publicKeys
func GetEcdsaPossiblePubkey(curve elliptic.Curve, hash []byte, R, S *big.Int) ([2]*ecdsa.PublicKey, error) {
	var possiblePubkey [2]*ecdsa.PublicKey
	N := curve.Params().N
	z := new(big.Int).SetBytes(hash)
	//ModInverse : P/s = P * s^-1，s^-1 =  new(bit.Int).ModInverse(s,N)
	sInv := new(big.Int).ModInverse(S, N)
	rInv := new(big.Int).ModInverse(R, N)
	x := R
	//caculate y by x
	ySquared := new(big.Int).Exp(x, new(big.Int).SetInt64(3), curve.Params().P)
	ySquared.Sub(ySquared, new(big.Int).Mul(x, big.NewInt(int64(3))))
	ySquared.Add(ySquared, curve.Params().B)
	y := new(big.Int).ModSqrt(ySquared, curve.Params().P)
	if y == nil {
		return possiblePubkey, fmt.Errorf("ModSqrt err")
	}

	for j := 0; j < 2; j++ {
		if j == 1 {
			y = new(big.Int).Neg(y)
		}
		p := new(ecdsa.PublicKey)
		p.Curve = curve
		p.X = x
		p.Y = y
		//u1 := new(big.Int).Mul(z, rInv)
		//u1.Mod(u1, N)
		u1 := new(ecdsa.PublicKey)
		u1.X, u1.Y = curve.ScalarBaseMult(z.Bytes())
		u1.X, u1.Y = curve.ScalarMult(u1.X, u1.Y, sInv.Bytes())

		//p-u1
		u2 := new(ecdsa.PublicKey)
		u1.Y = new(big.Int).Neg(u1.Y)
		u2.X, u2.Y = curve.Add(p.X, p.Y, u1.X, u1.Y)

		Qa := new(ecdsa.PublicKey)
		Qa.Curve = curve
		//Qa = u2 * SR^-1
		tempX, tempY := curve.ScalarMult(u2.X, u2.Y, S.Bytes())
		Qa.X, Qa.Y = curve.ScalarMult(tempX, tempY, rInv.Bytes())
		recoverPubKey := new(ecdsa.PublicKey)
		recoverPubKey.Curve = curve
		recoverPubKey.X = Qa.X
		recoverPubKey.Y = Qa.Y
		//isValid := ecdsa.Verify(recoverPubKey, hash[:], R, S)
		//fmt.Println(isValid)
		possiblePubkey[j] = recoverPubKey
	}
	return possiblePubkey, nil
}
