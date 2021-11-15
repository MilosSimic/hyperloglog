package main

import (
	"fmt"
)

func main() {
	hll, err := NewFromPrecision(16)
	if err != nil {
		fmt.Println(err)
	}

	hll.Add([]byte("Test1"))
	hll.Add([]byte("Test1"))
	hll.Add([]byte("Test1"))
	hll.Add([]byte("Test1"))

	hll.Add([]byte("124"))
	hll.Add([]byte("pdm390"))
	hll.Add([]byte("abcd"))
	hll.Add([]byte("yk3801083841jhfnvofn es"))
	hll.Add([]byte("nrnrvnevn22480jn"))
	hll.Add([]byte("Milos"))

	fmt.Println(hll.Estimate())
}
