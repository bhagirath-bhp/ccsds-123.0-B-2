package ccsds

type LocalSumMode int

const (
	WideNeighbourOriented LocalSumMode = iota
	WideColumnOriented
	NarrowNeighbourOriented
	NarrowColumnOriented
)

func (ic *ImgCube) localSum(z, y, x int) int {
	switch ic.p.LocalSumMode {
	case WideNeighbourOriented:
		return ic.wideNeighbourOrientedLocalSum(z, y, x)
	case WideColumnOriented:
		return ic.wideColumnOrientedLocalSum(z, y, x)
	case NarrowNeighbourOriented:
		return ic.narrowNeighbourOrientedLocalSum(z, y, x)
	case NarrowColumnOriented:
		return ic.narrowColumnOrientedLocalSum(z, y, x)
	default:
		panic("Invalid scope!")
	}
}

func (ic *ImgCube) wideNeighbourOrientedLocalSum(z, y, x int) int {
	switch {
	case y > 0 && (x > 0 && x < ic.p.NX-1):
		sZYXm1 := ic.sampleRepresentativeAt(z, y, x-1)
		sZYm1Xm1 := ic.sampleRepresentativeAt(z, y-1, x-1)
		sZYm1X := ic.sampleRepresentativeAt(z, y-1, x)
		sZYm1Xp1 := ic.sampleRepresentativeAt(z, y-1, x+1)
		return sZYXm1 + sZYm1Xm1 + sZYm1X + sZYm1Xp1

	case y == 0 && x > 0:
		sZYXm1 := ic.sampleRepresentativeAt(z, y, x-1)
		return 4 * sZYXm1

	case y > 0 && x == 0:
		sZYm1X := ic.sampleRepresentativeAt(z, y-1, x)
		sZYm1Xp1 := ic.sampleRepresentativeAt(z, y-1, x+1)
		return 2 * (sZYm1X + sZYm1Xp1)

	case y > 0 && x == ic.p.NX-1:
		sZYXm1 := ic.sampleRepresentativeAt(z, y, x-1)
		sZYm1Xm1 := ic.sampleRepresentativeAt(z, y-1, x-1)
		sZYm1X := ic.sampleRepresentativeAt(z, y-1, x)
		return sZYXm1 + sZYm1Xm1 + (2 * sZYm1X)

	default:
		panic("Invalid coords!")
	}
}

func (ic *ImgCube) wideColumnOrientedLocalSum(z, y, x int) int {
	switch {
	case y > 0:
		sZYm1X := ic.sampleRepresentativeAt(z, y-1, x)
		return 4 * sZYm1X
	case y == 0 && x > 0:
		sZYXm1 := ic.sampleRepresentativeAt(z, y, x-1)
		return 4 * sZYXm1
	default:
		panic("Invalid coords!")
	}
}

func (ic *ImgCube) narrowNeighbourOrientedLocalSum(z, y, x int) int {
	switch {
	case y > 0 && (x > 0 && x < ic.p.NX-1):
		sZYm1Xm1 := ic.sampleRepresentativeAt(z, y-1, x-1)
		sZYm1X := ic.sampleRepresentativeAt(z, y-1, x)
		sZYm1Xp1 := ic.sampleRepresentativeAt(z, y-1, x+1)
		return sZYm1Xm1 + (2 * sZYm1X) + sZYm1Xp1
	case y == 0 && x > 0 && z > 0:
		sZm1YXm1 := ic.sampleRepresentativeAt(z-1, y, x-1)
		return 4 * sZm1YXm1
	case y > 0 && x == 0:
		sZYm1X := ic.sampleRepresentativeAt(z, y-1, x)
		sZYm1Xp1 := ic.sampleRepresentativeAt(z, y-1, x+1)
		return 2 * (sZYm1X + sZYm1Xp1)
	case y > 0 && x == ic.p.NX-1:
		sZYm1Xm1 := ic.sampleRepresentativeAt(z, y-1, x-1)
		sZYm1X := ic.sampleRepresentativeAt(z, y-1, x)
		return 2 * (sZYm1Xm1 + sZYm1X)
	case y == 0 && x > 0 && z == 0:
		return 4 * ic.sMid()
	default:
		panic("Invalid coords!")
	}
}

func (ic *ImgCube) narrowColumnOrientedLocalSum(z, y, x int) int {
	switch {
	case y > 0:
		sZYm1X := ic.sampleRepresentativeAt(z, y-1, x)
		return 4 * sZYm1X
	case y == 0 && x > 0 && z > 0:
		sZm1YXm1 := ic.sampleRepresentativeAt(z-1, y, x-1)
		return 4 * sZm1YXm1
	case y == 0 && x > 0 && z == 0:
		return 4 * ic.sMid()
	default:
		panic("Invalid coords!")
	}
}
