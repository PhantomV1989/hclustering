package tree

import (
	"gonum.org/v1/gonum/floats"
)

type Leaf struct {
	Orig []float64
	Body []float64
}

type Tree struct {
	LeafSize    int
	BranchCount int
	Orig        []float64
	Leaf        []float64
	Children    []*Tree
}

func CreateTreeDecomposeMax(f []float64, LeafSize, BranchCount int, scaleOne bool) Tree {
	leaf := Downsample(f, LeafSize)
	normalizeMaxMin(leaf, scaleOne)

	tree := Tree{
		Orig:        f,
		Leaf:        leaf,
		LeafSize:    LeafSize,
		BranchCount: BranchCount,
	}
	tree.DecomposeMax(scaleOne)
	return tree
}

func CreateTree(f []float64, LeafSize, BranchCount int, scaleOne bool) Tree {
	leaf := Downsample(f, LeafSize)
	normalizeMaxMin(leaf, scaleOne)

	return Tree{
		Orig:        f,
		Leaf:        leaf,
		LeafSize:    LeafSize,
		BranchCount: BranchCount,
	}
}

// DecomposeMax decomposes tree until the smallest original leaf array is >= defined unit leaf size
func (t *Tree) DecomposeMax(scaleOne bool) {
	if len(t.Orig) >= t.BranchCount*t.LeafSize {
		branchArrs := partitionFloatArr(t.Orig, t.BranchCount)
		for b := range branchArrs {
			tt := CreateTree(branchArrs[b], t.LeafSize, t.BranchCount, scaleOne)
			tt.DecomposeMax(scaleOne)
			t.Children = append(t.Children, &tt)
		}
	}
}

// Decompose breaks tree down by n levels
func (t *Tree) Decompose(level int, scaleOne bool) {
	if level <= 0 {
		return
	}
	branchArrs := partitionFloatArr(t.Orig, t.BranchCount)
	for b := range branchArrs {
		tt := CreateTree(branchArrs[b], t.LeafSize, t.BranchCount, scaleOne)
		tt.Decompose(level-1, scaleOne)
		t.Children = append(t.Children, &tt)
	}
}

func normalizeMaxMin(f []float64, scaleOne bool) {
	miv := floats.Min(f)
	floats.AddConst(-miv, f)

	if scaleOne {
		mxv := floats.Max(f)
		miv = floats.Min(f)
		if mxv != miv {
			floats.Scale(1/mxv, f)
		}
	}
}

// Downsample ...
func Downsample(f []float64, toSize int) []float64 {
	if len(f) < toSize {
		return []float64{}
	}
	rv := []float64{}
	pf := partitionFloatArr(f, toSize)
	for i := range pf {
		v := floats.Sum(pf[i]) / float64(len(pf[i]))
		rv = append(rv, v)
	}
	return rv
}

func partitionFloatArr(f []float64, cnt int) [][]float64 {
	rv := [][]float64{}
	partitionSize := len(f) / cnt
	for i := 0; i < cnt; i++ {
		subf := f[i*partitionSize : (i+1)*partitionSize]
		rv = append(rv, subf)
	}
	return rv
}
