package currency

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConvertAmountMicros(t *testing.T) {
	tests := []struct {
		name          string
		baseMicros    int64
		rate          float64
		marginPercent float64
		wantMicros    int64
	}{
		{
			name:          "zero amount returns zero",
			baseMicros:    0,
			rate:          1.5,
			marginPercent: 2.0,
			wantMicros:    0,
		},
		{
			name:          "zero rate returns zero",
			baseMicros:    1000000,
			rate:          0.0,
			marginPercent: 2.0,
			wantMicros:    0,
		},
		{
			name:          "1:1 rate with zero margin",
			baseMicros:    1000000, // 1.00 USD
			rate:          1.0,
			marginPercent: 0.0,
			wantMicros:    1000000,
		},
		{
			name:          "USD to EUR (rate 0.92) zero margin",
			baseMicros:    10000000, // 10.00 USD
			rate:          0.92,
			marginPercent: 0.0,
			wantMicros:    9200000, // 9.20 EUR
		},
		{
			name:          "USD to NGN (rate 1550.0) with 1.5% margin",
			baseMicros:    1000000, // 1.00 USD
			rate:          1550.0,
			marginPercent: 1.5,
			// 1,000,000 * 1550.0 * 1.015 = 1,573,250,000 micros
			wantMicros: 1573250000,
		},
		{
			name:          "USD to GBP (rate 0.7854) with 2% margin rounding half up",
			baseMicros:    500000, // 0.50 USD
			rate:          0.7854,
			marginPercent: 2.0,
			// 500,000 * 0.7854 * 1.02 = 400554
			wantMicros: 400554,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ConvertAmountMicros(tt.baseMicros, tt.rate, tt.marginPercent)
			assert.Equal(t, tt.wantMicros, got)
		})
	}
}

func TestConvertTargetToBaseMicros(t *testing.T) {
	tests := []struct {
		name          string
		targetMicros  int64
		rate          float64
		marginPercent float64
		wantMicros    int64
	}{
		{
			name:          "zero target returns zero",
			targetMicros:  0,
			rate:          1.5,
			marginPercent: 2.0,
			wantMicros:    0,
		},
		{
			name:          "zero rate returns zero",
			targetMicros:  1000000,
			rate:          0.0,
			marginPercent: 2.0,
			wantMicros:    0,
		},
		{
			name:          "1:1 rate with zero margin",
			targetMicros:  1000000,
			rate:          1.0,
			marginPercent: 0.0,
			wantMicros:    1000000,
		},
		{
			name:          "invert NGN conversion with margin",
			targetMicros:  1573250000,
			rate:          1550.0,
			marginPercent: 1.5,
			// 1573250000 / (1550.0 * 1.015) = 1000000
			wantMicros: 1000000,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ConvertTargetToBaseMicros(tt.targetMicros, tt.rate, tt.marginPercent)
			assert.Equal(t, tt.wantMicros, got)
		})
	}
}
