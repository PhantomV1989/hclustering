package tree

import (
	"math/rand"
	"testing"
	"time"
)

func createRandomFloat(size int) []float64 {
	r := []float64{}
	rand.Seed(time.Now().Unix())
	for i := 0; i < size; i++ {
		r = append(r, rand.Float64())
	}
	return r
}

// TestX ...
func TestDownsample(t *testing.T) {
	x := createRandomFloat(100)
	sz := 20
	f := downsample(x, sz)
	if len(f) != sz {
		panic("")
	}

	sz = 80
	f = downsample(x, sz)
	if len(f) != sz || f[0] != x[0] {
		panic("")
	}

	sz = 33
	f = downsample(x, sz)
	if len(f) != sz {
		panic("")
	}
}

func TestCreateDecomposeTree(t *testing.T) {
	qwe := createRandomFloat(120)
	t1 := CreateTree(qwe, 0, 0, 20, 2)
	t1.Decompose(1)
	if len(t1.Children) != 2 {
		panic("len(t1.Children)!=2")
	}
	if t1.Children[0].Level != 1 {
		panic("t1.Children[0].Level!=1")
	}
	if t1.Children[0].Orig[0] != qwe[0] {
		panic("t1.Children[0].Orig[0] != qwe[0]")
	}
	if t1.Children[1].Orig[0] != qwe[60] {
		panic("t1.Children[1].Orig[0] != qwe[60]")
	}
	if t1.Children[0].Pos != 0 {
		panic("t1.Children[0].Pos != 0")
	}
	if t1.Children[1].Pos != 1 {
		panic("t1.Children[1].Pos != 1")
	}

}

func TestDecomposeMax(t *testing.T) {
	qwe := createRandomFloat(120)
	t1 := CreateTree(qwe, 0, 0, 20, 2)
	t1.DecomposeMax()
	if len(t1.Children) != 2 {
		panic("len(t1.Children)!=2")
	}
	if t1.Children[0].Level != 1 {
		panic("t1.Children[0].Level!=1")
	}
	if t1.Children[0].Orig[0] != qwe[0] {
		panic("t1.Children[0].Orig[0] != qwe[0]")
	}
	if t1.Children[1].Orig[0] != qwe[60] {
		panic("t1.Children[1].Orig[0] != qwe[60]")
	}
	if t1.Children[0].Pos != 0 {
		panic("t1.Children[0].Pos != 0")
	}
	if t1.Children[1].Pos != 1 {
		panic("t1.Children[1].Pos != 1")
	}

}
