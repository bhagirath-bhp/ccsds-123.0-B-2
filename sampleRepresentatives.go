package ccsds

func (ic *ImgCube) sampleRepresentativeAt(z, y, x int) int {
	idx := (z*ic.p.NY+y)*ic.p.NX + x
	return ic.p.Representatives[idx]
}

func (ic *ImgCube) sMin() int { return 0 }
func (ic *ImgCube) sMid() int { return 1 << (ic.p.DynamicRange - 1) }
func (ic *ImgCube) sMax() int { return (1 << ic.p.DynamicRange) - 1 }
