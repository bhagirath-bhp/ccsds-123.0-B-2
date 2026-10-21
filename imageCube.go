package ccsds

import "fmt"

type PredictionMode int

const (
	FullPrediction PredictionMode = iota
	ReducedPrediction
)

type ImgCubeParams struct {
	DynamicRange    int
	NX, NY, NZ      int
	Representatives []int
	LocalSumMode    LocalSumMode
	PredictionMode  PredictionMode
	NPreviousBands  int
	Omega           int
	// Weight initialization (4.6.3). Lambda and Q are only used when
	// WeightInit == CustomWeightInit; Lambda[z] must have C_z entries.
	WeightInit WeightInitMode
	Lambda     [][]int
	Q          int
	// Weight update scaling exponent parameters (4.10.2): v_min, v_max, t_inc.
	VMin, VMax int
	TInc       int
}

type ImgCube struct {
	p ImgCubeParams
}

func NewImageCube(p ImgCubeParams) (ImgCube, error) {
	if err := p.validate(); err != nil {
		return ImgCube{}, err
	}
	return ImgCube{p: p}, nil
}

func (p ImgCubeParams) validate() error {
	switch {
	case p.NX < 2 || p.NY < 1 || p.NZ < 1:
		return fmt.Errorf("invalid dimensions %dx%dx%d", p.NX, p.NY, p.NZ)
	case len(p.Representatives) != p.NX*p.NY*p.NZ:
		return fmt.Errorf("len(representatives) = %d, want %d", len(p.Representatives), p.NX*p.NY*p.NZ)
	case p.DynamicRange < 2 || p.DynamicRange > 32:
		return fmt.Errorf("dynamic range %d out of range [2, 32]", p.DynamicRange)
	case p.Omega < 4 || p.Omega > 19:
		return fmt.Errorf("omega %d out of range [4, 19]", p.Omega)
	case p.NPreviousBands < 0 || p.NPreviousBands > 15:
		return fmt.Errorf("nPreviousBands %d out of range [0, 15]", p.NPreviousBands)
	case p.VMin < -6 || p.VMin > p.VMax || p.VMax > 9:
		return fmt.Errorf("vMin %d, vMax %d must satisfy -6 <= vMin <= vMax <= 9", p.VMin, p.VMax)
	case p.TInc < 16 || p.TInc > 2048 || p.TInc&(p.TInc-1) != 0:
		return fmt.Errorf("tInc %d must be a power of 2 in [2^4, 2^11]", p.TInc)
	}
	return p.validateWeightInit()
}

func (p ImgCubeParams) validateWeightInit() error {
	switch p.WeightInit {
	case DefaultWeightInit:
		return nil
	case CustomWeightInit:
		if p.Q < 3 || p.Q > p.Omega+3 {
			return fmt.Errorf("weight initialization resolution Q %d is out of range [3, %d]", p.Q, p.Omega+3)
		}
		if len(p.Lambda) != p.NZ {
			return fmt.Errorf("len(Lambda) = %d, want NZ = %d", len(p.Lambda), p.NZ)
		}
		lMin, lMax := -(1 << (p.Q - 1)), (1<<(p.Q-1))-1 // signed Q-bit range
		for z, lambda := range p.Lambda {
			if want := p.weightVectorLen(z); len(lambda) != want {
				return fmt.Errorf("len(Lambda[%d]) = %d, want C_z = %d", z, len(lambda), want)
			}
			for i, l := range lambda {
				if l < lMin || l > lMax {
					return fmt.Errorf("Lambda[%d][%d] = %d out of signed %d-bit range [%d, %d]", z, i, l, p.Q, lMin, lMax)
				}
			}
		}
		return nil
	default:
		return fmt.Errorf("invalid weight initialization mode %d", p.WeightInit)
	}
}
