package main

import (
	"errors"
	"math"
	"math/bits"
	"hash/fnv"
)

const (
	HLL_MIN_PRECISION = 4
	HLL_MAX_PRECISION = 16
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

	m := uint64(math.Pow(2, float64(p)))
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

func (hll *HLL) fingerprint(b []byte) uint32 {
	hash := fnv.New32a()
	hash.Write(b)
	return hash.Sum32()
}

func (hll *HLL) Add(data []byte) {
	//shift on right for precision amount of bits
	x := hll.fingerprint(data)
	k := uint32(32-hll.p)
	r := uint8(1 + bits.LeadingZeros32(x << hll.p))
	i    := x >> uint8(k)
	if r > hll.reg[i] {
		hll.reg[i] = r
	}
}

func (hll *HLL) Estimate() float64 {
	sum := 0.0
	for _, val := range hll.reg {
		sum += math.Pow(math.Pow(2.0, float64(val)),-1)
	}

	alpha := 0.7213 / (1.0 + 1.079/float64(hll.m))
	estimation := alpha * math.Pow(float64(hll.m), 2.0) / sum
	emptyRegs := hll.emptyCount()
	if estimation <= 2.5*float64(hll.m) { // do small range correction
		if emptyRegs > 0 {
			estimation = float64(hll.m) * math.Log(float64(hll.m)/float64(emptyRegs))
		}
	} else if estimation > 1/30.0*math.Pow(2.0, 32.0) { // do large range correction
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

func (hll *HLL) emptyCount() int {
	sum := 0
	for _, val := range hll.reg {
		if val == 0 {
			sum++
		}
	}
	return sum
}
