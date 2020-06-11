package main

import (
	"sync"
)

var wg sync.WaitGroup

func main() {
	q := []float64{}
	for i := 0; i < 100; i++ {
		q = append(q, float64(i))
	}
}
