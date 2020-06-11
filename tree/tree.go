package tree

import (
	"strconv"

	"gonum.org/v1/gonum/floats"
)

type Leaf struct {
	Orig []float64
	Body []float64
}

type Tree struct {
	Level       int
	Pos         int
	LeafSize    int
	BranchCount int
	Orig        []float64
	Leaf        []float64
	Children    []*Tree
}

func CreateTree(f []float64, Level, Pos, LeafSize, BranchCount int) Tree {
	leaf := downsample(f, LeafSize)
	miv := floats.Min(leaf)
	floats.AddConst(-miv, leaf)

	return Tree{
		Level:       Level,
		Pos:         Pos,
		Orig:        f,
		Leaf:        leaf,
		LeafSize:    LeafSize,
		BranchCount: BranchCount,
	}
}

func (t *Tree) DecomposeMax() {
	if len(t.Orig) > 2*t.LeafSize {
		branchArrs := partitionFloatArr(t.Orig, t.BranchCount)
		for b := range branchArrs {
			tt := CreateTree(branchArrs[b], t.Level+1, b, t.LeafSize, t.BranchCount)
			tt.DecomposeMax()
			t.Children = append(t.Children, &tt)
		}
	}

}

func (t *Tree) Decompose(level int) {
	if level <= 0 {
		return
	}
	branchArrs := partitionFloatArr(t.Orig, t.BranchCount)
	for b := range branchArrs {
		tt := CreateTree(branchArrs[b], t.Level+1, b, t.LeafSize, t.BranchCount)
		tt.Decompose(level - 1)
		t.Children = append(t.Children, &tt)
	}
}

func (t *Tree) GetPositionString() string {
	return strconv.Itoa(t.Level) + "." + strconv.Itoa(t.Pos)
}

func downsample(f []float64, toSize int) []float64 {
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
