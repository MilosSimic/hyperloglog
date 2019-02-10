package main

import (
	"fmt"
	"hash/fnv"
)

func fingerprint(b []byte) uint64 {
	hash := fnv.New64a()
	hash.Write(b)
	return hash.Sum64()
}

func main() {
	hll, err := NewFromPrecision(4)
	if err != nil {
		fmt.Println(err)
	}

	n := fingerprint([]byte("Test"))
	n1 := fingerprint([]byte("Test"))
	n2 := fingerprint([]byte("Test2"))
	n3 := fingerprint([]byte("Test3"))

	hll.Add(n)
	hll.Add(n1)
	hll.Add(n2)
	hll.Add(n3)

	fmt.Println(hll.Estimate())
}
