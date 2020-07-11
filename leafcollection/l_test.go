package leafcollection

import (
	"math/rand"
	"testing"
	"time"

	"github.com/phantomv1989/hclustering/tree"
)

func createRandomFloat(size int) []float64 {
	r := []float64{}
	rand.Seed(time.Now().Unix())
	for i := 0; i < size; i++ {
		r = append(r, rand.Float64())
	}
	return r
}

func TestA(t *testing.T) {
	leafSize := 10
	branchCount := 4
	scoreLim := 0.001
	a, b, c := []float64{}, []float64{}, []float64{}
	for i := 0; i < 100; i++ {
		a = append(a, float64(i))
		b = append(b, float64(i))
		c = append(c, 1)
	}
	for i := 0; i < 100; i++ {
		a = append(a, float64(i*2))
		b = append(b, float64(i*-2))
	}

	h := func(x []float64) tree.Tree {
		x1 := tree.CreateTreeDecomposeMax(x, leafSize, branchCount)
		return x1
	}

	aTree := h(a)
	bTree := h(b)
	cTree := h(c)

	leafCollection := []LeafData{}

	InsertLeavesRecursive("0", &aTree, &leafCollection, scoreLim)
	InsertLeavesRecursive("0", &bTree, &leafCollection, scoreLim)
	InsertLeavesRecursive("0", &cTree, &leafCollection, scoreLim)

	FindAllLeafPositions("0", &aTree, &leafCollection, scoreLim, true)
	GetHierarchicalMatches("0", &bTree, &leafCollection, scoreLim, true)
	InsertLeavesRecursive("0", &aTree, &leafCollection, scoreLim)
	InsertLeavesRecursive("0", &bTree, &leafCollection, scoreLim)

	//SaveLeafCollection("./qweqwe", &leafCollection)
	// println(qwe)

	// asd := LoadLeafCollection("./qweqwe")
	// println(len(asd))
}
