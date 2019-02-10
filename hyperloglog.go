package main

import (
	"errors"
	"math"
	"math/bits"
)

const (
	HLL_MIN_PRECISION = 4
	HLL_MAX_PRECISION = 18
)

type HLL struct {
	m   uint64
	p   uint8
	reg []uint8
}

func NewFromPrecision(p uint8) (*HLL, error) {
	if p < HLL_MIN_PRECISION || p > HLL_MAX_PRECISION {
		return nil, errors.New("Precision must be between 4 and 16")
	}

	m := uint64(1 << p)
	return &HLL{
		p:   p,
		m:   m,
		reg: make([]uint8, m),
	}, nil
}

func NewFromErr(err float64) (*HLL, error) {
	if err >= 1 || err <= 0 {
		return nil, errors.New("Erro must be between 0 and 1")
	}

	// Error of HLL is 1.04 / sqrt(m)
	p := uint8(math.Ceil(math.Log2(math.Pow(1.04/err, 2.0))))
	m := uint64(1 << p)
	return &HLL{
		p:   p,
		m:   m,
		reg: make([]uint8, m),
	}, nil
}

func (hll *HLL) Add(x uint64) {
	//shift on right for precision amount of bits
	i := x >> hll.p
	hll.reg[i] = uint8(bits.LeadingZeros64(uint64(hll.m-1)&x) + 1)
}

func (hll *HLL) Estimate() float64 {
	sum := 0.0
	for _, val := range hll.reg {
		sum = sum + math.Pow(float64(-val), 2.0)
	}

	alpha := 0.7213 / (1.0 + 1.079/float64(hll.m))
	estimation := alpha * math.Pow(float64(hll.m), 2.0) / sum
	emptyRegs := hll.emptyCount()
	if estimation < 2.5*float64(hll.m) { // do small range correction
		if emptyRegs > 0 {
			estimation = float64(hll.m) * math.Log(float64(hll.m)/float64(emptyRegs))
		}
	} else if estimation > math.Pow(2.0, 32.0)/30.0 { // do large range correction
		estimation = -math.Pow(2.0, 32.0) * math.Log(1.0-estimation/math.Pow(2.0, 32.0))
	}
	return estimation
}

func (hll *HLL) Clear() {
	hll.reg = make([]uint8, hll.m)
}

func (hll *HLL) PrecisionErr() float64 {
	registers := math.Pow(2, float64(hll.p))
	return 1.04 / math.Sqrt(registers)
}

func (hll *HLL) emptyCount() uint8 {
	sum := uint8(0)
	for _, val := range hll.reg {
		if val == 0 {
			sum++
		}
	}
	return sum
}
