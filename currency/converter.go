package currency

import (
	"math"
)

// ConvertAmountMicros converts a monetary amount in micros from a base currency
// to a target currency using the given exchange rate and platform margin percentage.
//
// Formula:
//
//	Target Amount (Micros) = round(Base Amount (Micros) * Rate * (1 + MarginPercent / 100))
//
// 1 USD = 1,000,000 micros.
func ConvertAmountMicros(baseMicros int64, rate float64, marginPercent float64) int64 {
	if baseMicros == 0 || rate == 0 {
		return 0
	}

	multiplier := 1.0 + (marginPercent / 100.0)
	effectiveRate := rate * multiplier
	converted := float64(baseMicros) * effectiveRate

	return int64(math.Round(converted))
}

// ConvertTargetToBaseMicros inverts the calculation to find the approximate base USD micros
// from a target currency amount.
func ConvertTargetToBaseMicros(targetMicros int64, rate float64, marginPercent float64) int64 {
	if targetMicros == 0 || rate == 0 {
		return 0
	}

	multiplier := 1.0 + (marginPercent / 100.0)
	effectiveRate := rate * multiplier
	base := float64(targetMicros) / effectiveRate

	return int64(math.Round(base))
}
