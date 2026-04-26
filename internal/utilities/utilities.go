package utilities

import (
	"math/big"
)

func RatToFloat(r *big.Rat) float64 {
	f, _ := r.Float64()
	return f
}
