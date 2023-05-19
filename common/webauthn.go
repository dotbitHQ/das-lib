package common

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/hex"
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
	hash := sha256.Sum256([]byte(cid))
	for i := 0; i < 4; i++ {
		hash = sha256.Sum256(hash[:])
	}
	return hash[22:]
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
	return hash[22:]
}
