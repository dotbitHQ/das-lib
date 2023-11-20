package common

import (
	"fmt"
	"math"
	"math/big"
	"strconv"
)

const (
	StartPremium int64 = 100000000
)

func Premium(startTime, nowTime int64) float64 {
	if startTime > nowTime {
		return float64(StartPremium)
	}
	numerator := new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(8), nil))
	//2^(n/86400)
	//n/86400
	exponent := new(big.Float).Quo(big.NewFloat(float64(nowTime-startTime)), big.NewFloat(86400))
	exponentFloat, _ := exponent.Float64()
	//2^(n/86400)
	denominator := math.Pow(2, exponentFloat)
	// 10^8 / 2^(n/86400)
	result := new(big.Float).Quo(numerator, new(big.Float).SetFloat64(denominator))

	res, _ := result.Float64()
	num, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", res), 64)
	return num
}
