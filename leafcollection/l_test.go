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
	a, b := []float64{}, []float64{}
	for i := 0; i < 100; i++ {
		a = append(a, float64(i))
		b = append(b, float64(i))
	}
	for i := 0; i < 100; i++ {
		a = append(a, float64(i*2))
		b = append(b, float64(i*-2))
	}

	h := func(x []float64) tree.Tree {
		x1 := tree.CreateTree(x, 0, 0, leafSize, branchCount)
		x1.DecomposeMax()
		return x1
	}

	aTree := h(a)
	bTree := h(b)

	leafCollection := []LeafData{}

	InsertLeafCollectionRecursive(&aTree, &leafCollection, scoreLim, InsertMode.Insert)
	InsertLeafCollectionRecursive(&bTree, &leafCollection, scoreLim, InsertMode.Insert)

	if len(leafCollection[2].MatchedPositions) != 4 {
		panic("len(leafCollection[2].MatchedPositions)!=4")
	}
	if leafCollection[2].MatchedPositions["2.0"] != 4 {
		panic("")
	}
	result := map[string]int{}
	FindAllLeafPositions(&aTree, &leafCollection, &result, scoreLim)
	println("asd")

	InsertLeafCollectionRecursive(&aTree, &leafCollection, scoreLim, InsertMode.Delete)
	InsertLeafCollectionRecursive(&bTree, &leafCollection, scoreLim, InsertMode.Delete)

	if leafCollection[2].MatchedPositions["2.0"] != 0 {
		panic("")
	}

}
