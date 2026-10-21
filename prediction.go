package ccsds

import "fmt"

func (ic *ImgCube) predictCentralLocalDiff(W []int, U []int) ([]int, error) {
	if len(W) != len(U) {
		return nil, fmt.Errorf("Invalid weight/local diff vector!")
	}
	cld_predicted := make([]int, len(U))
	for i := 0; i < len(W); i++ {
		cld_predicted[i] = W[i] * U[i]
	}
	return cld_predicted, nil
}

func (ic *ImgCube) predict(z, y, x int) {
	U, err := ic.localDifferenceVector(z, y, x)
	if err != nil {
		panic()
	}
	W, 
}
