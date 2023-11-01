package common

import (
	"fmt"
	"math/big"
	"time"
)

const (
	Bit1  int64 = 999989423469314432
	Bit2  int64 = 999978847050491904
	Bit3  int64 = 999957694548431104
	Bit4  int64 = 999915390886613504
	Bit5  int64 = 999830788931929088
	Bit6  int64 = 999661606496243712
	Bit7  int64 = 999323327502650752
	Bit8  int64 = 998647112890970240
	Bit9  int64 = 997296056085470080
	Bit10 int64 = 994599423483633152
	Bit11 int64 = 989228013193975424
	Bit12 int64 = 978572062087700096
	Bit13 int64 = 957603280698573696
	Bit14 int64 = 917004043204671232
	Bit15 int64 = 840896415253714560
	Bit16 int64 = 707106781186547584

	PRECISION    int64 = 1e18
	StartPremium int64 = 100000000
	GRACE_PERIOD int64 = 90 * 24 * 3600
	TOTALDAYS    int64 = 27
)

func Premium(expires int64) int64 {
	endValue := StartPremium >> TOTALDAYS
	expires = expires + GRACE_PERIOD
	if expires > time.Now().Unix() {
		return StartPremium
	}
	//default now time 1698140024
	//nowTime := int64(1698140024)
	nowTime := time.Now().Unix()
	elapsed := nowTime - expires
	fmt.Println("elapsed:", elapsed)
	premium := decayedPremium(elapsed)
	fmt.Println("premium:", premium)
	if premium >= endValue {
		return (premium - endValue)
	}
	return 0
}

func decayedPremium(elapsed int64) int64 {
	elapsedBI := new(big.Int).SetInt64(elapsed)
	PRECISIONBI := new(big.Int).SetInt64(PRECISION)
	//daysPast := (elapsed * PRECISION) / 24 * 3600
	daysPastBI := new(big.Int).Div(new(big.Int).Mul(elapsedBI, PRECISIONBI), new(big.Int).SetInt64(24*3600))
	//intDays := daysPast / PRECISION
	initDaysBI := new(big.Int).Div(daysPastBI, PRECISIONBI)
	//premium := _startPremium >> intDays
	premiumBI := new(big.Int).Rsh(new(big.Int).SetInt64(StartPremium), uint(initDaysBI.Int64()))
	//partDay := (daysPast - intDays*PRECISION)
	partDayBI := new(big.Int).Sub(daysPastBI, new(big.Int).Mul(initDaysBI, PRECISIONBI))
	//fraction := (partDay * (2 ^ 16)) / PRECISION
	temp := new(big.Int).Exp(new(big.Int).SetInt64(2), new(big.Int).SetInt64(16), nil)
	fractionBI := new(big.Int).Div(new(big.Int).Mul(partDayBI, temp), PRECISIONBI)
	totalPremium := addFractionalPremium(fractionBI, premiumBI)
	return totalPremium.Int64()
}

func addFractionalPremium(fraction, premium *big.Int) *big.Int {
	PRECISIONBI := new(big.Int).SetInt64(PRECISION)

	//if fraction&(1<<0) != 0 {
	//	premium = (premium * Bit1) / PRECISION
	//}
	var temp *big.Int
	var res *big.Int
	temp = new(big.Int).Lsh(new(big.Int).SetInt64(1), uint(0))
	res = new(big.Int).And(fraction, temp)
	if isZero := res.Cmp(new(big.Int).SetInt64(0)); isZero != 0 {
		premium = new(big.Int).Div(new(big.Int).Mul(premium, new(big.Int).SetInt64(Bit1)), PRECISIONBI)
	}

	//if fraction&(1<<1) != 0 {
	//	premium = (premium * Bit2) / PRECISION
	//}
	temp = new(big.Int).Lsh(new(big.Int).SetInt64(1), uint(1))
	res = new(big.Int).And(fraction, temp)
	if isZero := res.Cmp(new(big.Int).SetInt64(0)); isZero != 0 {
		premium = new(big.Int).Div(new(big.Int).Mul(premium, new(big.Int).SetInt64(Bit2)), PRECISIONBI)
	}

	//if fraction&(1<<2) != 0 {
	//	premium = (premium * Bit3) / PRECISION
	//}
	temp = new(big.Int).Lsh(new(big.Int).SetInt64(1), uint(2))
	res = new(big.Int).And(fraction, temp)
	if isZero := res.Cmp(new(big.Int).SetInt64(0)); isZero != 0 {
		premium = new(big.Int).Div(new(big.Int).Mul(premium, new(big.Int).SetInt64(Bit3)), PRECISIONBI)
	}
	//if fraction&(1<<3) != 0 {
	//	premium = (premium * Bit4) / PRECISION
	//}
	temp = new(big.Int).Lsh(new(big.Int).SetInt64(1), uint(3))
	res = new(big.Int).And(fraction, temp)
	if isZero := res.Cmp(new(big.Int).SetInt64(0)); isZero != 0 {
		premium = new(big.Int).Div(new(big.Int).Mul(premium, new(big.Int).SetInt64(Bit4)), PRECISIONBI)
	}
	//if fraction&(1<<4) != 0 {
	//	premium = (premium * Bit5) / PRECISION
	//}
	temp = new(big.Int).Lsh(new(big.Int).SetInt64(1), uint(4))
	res = new(big.Int).And(fraction, temp)
	if isZero := res.Cmp(new(big.Int).SetInt64(0)); isZero != 0 {
		premium = new(big.Int).Div(new(big.Int).Mul(premium, new(big.Int).SetInt64(Bit5)), PRECISIONBI)
	}
	//if fraction&(1<<5) != 0 {
	//	premium = (premium * Bit6) / PRECISION
	//}
	temp = new(big.Int).Lsh(new(big.Int).SetInt64(1), uint(5))
	res = new(big.Int).And(fraction, temp)
	if isZero := res.Cmp(new(big.Int).SetInt64(0)); isZero != 0 {
		premium = new(big.Int).Div(new(big.Int).Mul(premium, new(big.Int).SetInt64(Bit6)), PRECISIONBI)
	}
	//if fraction&(1<<6) != 0 {
	//	premium = (premium * Bit7) / PRECISION
	//}
	temp = new(big.Int).Lsh(new(big.Int).SetInt64(1), uint(6))
	res = new(big.Int).And(fraction, temp)
	if isZero := res.Cmp(new(big.Int).SetInt64(0)); isZero != 0 {
		premium = new(big.Int).Div(new(big.Int).Mul(premium, new(big.Int).SetInt64(Bit7)), PRECISIONBI)
	}
	//if fraction&(1<<7) != 0 {
	//	premium = (premium * Bit8) / PRECISION
	//}
	temp = new(big.Int).Lsh(new(big.Int).SetInt64(1), uint(7))
	res = new(big.Int).And(fraction, temp)
	if isZero := res.Cmp(new(big.Int).SetInt64(0)); isZero != 0 {
		premium = new(big.Int).Div(new(big.Int).Mul(premium, new(big.Int).SetInt64(Bit8)), PRECISIONBI)
	}
	//if fraction&(1<<8) != 0 {
	//	premium = (premium * Bit9) / PRECISION
	//}
	temp = new(big.Int).Lsh(new(big.Int).SetInt64(1), uint(8))
	res = new(big.Int).And(fraction, temp)
	if isZero := res.Cmp(new(big.Int).SetInt64(0)); isZero != 0 {
		premium = new(big.Int).Div(new(big.Int).Mul(premium, new(big.Int).SetInt64(Bit9)), PRECISIONBI)
	}

	//if fraction&(1<<9) != 0 {
	//	premium = (premium * Bit10) / PRECISION
	//}
	temp = new(big.Int).Lsh(new(big.Int).SetInt64(1), uint(9))
	res = new(big.Int).And(fraction, temp)
	if isZero := res.Cmp(new(big.Int).SetInt64(0)); isZero != 0 {
		premium = new(big.Int).Div(new(big.Int).Mul(premium, new(big.Int).SetInt64(Bit10)), PRECISIONBI)
	}

	//if fraction&(1<<10) != 0 {
	//	premium = (premium * Bit11) / PRECISION
	//}
	temp = new(big.Int).Lsh(new(big.Int).SetInt64(1), uint(10))
	res = new(big.Int).And(fraction, temp)
	if isZero := res.Cmp(new(big.Int).SetInt64(0)); isZero != 0 {
		premium = new(big.Int).Div(new(big.Int).Mul(premium, new(big.Int).SetInt64(Bit11)), PRECISIONBI)
	}

	//if fraction&(1<<11) != 0 {
	//	premium = (premium * Bit12) / PRECISION
	//}
	temp = new(big.Int).Lsh(new(big.Int).SetInt64(1), uint(11))
	res = new(big.Int).And(fraction, temp)
	if isZero := res.Cmp(new(big.Int).SetInt64(0)); isZero != 0 {
		premium = new(big.Int).Div(new(big.Int).Mul(premium, new(big.Int).SetInt64(Bit12)), PRECISIONBI)
	}
	//if fraction&(1<<12) != 0 {
	//	premium = (premium * Bit13) / PRECISION
	//}
	temp = new(big.Int).Lsh(new(big.Int).SetInt64(1), uint(12))
	res = new(big.Int).And(fraction, temp)
	if isZero := res.Cmp(new(big.Int).SetInt64(0)); isZero != 0 {
		premium = new(big.Int).Div(new(big.Int).Mul(premium, new(big.Int).SetInt64(Bit13)), PRECISIONBI)
	}
	//if fraction&(1<<13) != 0 {
	//	premium = (premium * Bit14) / PRECISION
	//}
	temp = new(big.Int).Lsh(new(big.Int).SetInt64(1), uint(13))
	res = new(big.Int).And(fraction, temp)
	if isZero := res.Cmp(new(big.Int).SetInt64(0)); isZero != 0 {
		premium = new(big.Int).Div(new(big.Int).Mul(premium, new(big.Int).SetInt64(Bit14)), PRECISIONBI)
	}

	//if fraction&(1<<14) != 0 {
	//	premium = (premium * Bit15) / PRECISION
	//}
	temp = new(big.Int).Lsh(new(big.Int).SetInt64(1), uint(14))
	res = new(big.Int).And(fraction, temp)
	if isZero := res.Cmp(new(big.Int).SetInt64(0)); isZero != 0 {
		premium = new(big.Int).Div(new(big.Int).Mul(premium, new(big.Int).SetInt64(Bit15)), PRECISIONBI)
	}

	//if fraction&(1<<15) != 0 {
	//	premium = (premium * Bit16) / PRECISION
	//}
	temp = new(big.Int).Lsh(new(big.Int).SetInt64(1), uint(15))
	res = new(big.Int).And(fraction, temp)
	if isZero := res.Cmp(new(big.Int).SetInt64(0)); isZero != 0 {
		premium = new(big.Int).Div(new(big.Int).Mul(premium, new(big.Int).SetInt64(Bit16)), PRECISIONBI)
	}
	return premium
}
