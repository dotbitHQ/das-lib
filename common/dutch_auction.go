package common

import (
	"fmt"
	"github.com/dotbitHQ/das-lib/core"
	"math"
	"math/big"
	"strconv"
)

const (
	StartPremium int64 = 100000000
	GRACE_PERIOD int64 = 90 * 24 * 3600
)

func Premium(expires, nowTime int64) float64 {
	expires = expires + GRACE_PERIOD
	if expires > nowTime {
		return float64(StartPremium)
	}
	numerator := new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(8), nil))
	//2^(n/86400)
	//n/86400
	exponent := new(big.Float).Quo(big.NewFloat(float64(nowTime-expires)), big.NewFloat(86400))
	exponentFloat, _ := exponent.Float64()
	//2^(n/86400)
	denominator := math.Pow(2, exponentFloat)
	// 10^8 / 2^(n/86400)
	result := new(big.Float).Quo(numerator, new(big.Float).SetFloat64(denominator))

	res, _ := result.Float64()
	num, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", res), 64)
	return num
}

type AuctionConfig struct {
	GracePeriodTime, AuctionPeriodTime, DeliverPeriodTime uint32
}

func GetAuctionConfig(dasCore *core.DasCore) (res *AuctionConfig, err error) {
	builderConfigCell, err := dasCore.ConfigCellDataBuilderByTypeArgs(ConfigCellTypeArgsAccount)
	if err != nil {
		err = fmt.Errorf("ConfigCellDataBuilderByTypeArgs err: %s", err.Error())
		return
	}
	gracePeriodTime, err := builderConfigCell.ExpirationGracePeriod()
	if err != nil {
		err = fmt.Errorf("ExpirationGracePeriod err: %s", err.Error())
		return
	}
	auctionPeriodTime, err := builderConfigCell.ExpirationAuctionPeriod()
	if err != nil {
		err = fmt.Errorf("ExpirationAuctionPeriod err: %s", err.Error())
		return
	}
	deliverPeriodTime, err := builderConfigCell.ExpirationDeliverPeriod()
	if err != nil {
		err = fmt.Errorf("ExpirationDeliverPeriod err: %s", err.Error())
		return
	}
	res = &AuctionConfig{
		GracePeriodTime:   gracePeriodTime,
		AuctionPeriodTime: auctionPeriodTime,
		DeliverPeriodTime: deliverPeriodTime,
	}
	return
}
