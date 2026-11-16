package ccsds_spike

import (
	"fmt"
)

func (ccsds *CCSDS) dotProduct(vec1, vec2 []int) int {
	if len(vec1) != len(vec2) {
		panic(fmt.Errorf("slice lengths must be same to perform dot product!"))
	}
	sum := 0
	for i, v1 := range vec1 {
		sum += v1 * vec2[i]
	}
	return sum
}

func (ccsds *CCSDS) clip(value, min, max int) int {
	switch {
	case value < min:
		return min
	case value > max:
		return max
	default:
		return value
	}
}

func (ccsds *CCSDS) signedMod(x, r int) int {
	exp := 1 << r
	offset := exp / 2
	return ((x + offset) & (exp - 1)) - offset
}

func (ccsds *CCSDS) floorDiv(a, b int) int {
	res := a / b
	rem := a % b
	// If the signs are different and there is a remainder, floor down by 1
	if rem != 0 && ((a < 0) != (b < 0)) {
		res--
	}
	return res
}

func (ccsds *CCSDS) min(a, b int) int {
	return min(a, b)
}

func (ccsds *CCSDS) getSign(v int) int {
	switch {
	case v > 0:
		return 1
	case v < 0:
		return -1
	default:
		return 0
	}
}

func (ccsds *CCSDS) getSign_p(v int) int {
	if v > 0 {
		return 1
	} else {
		return -1
	}
}
