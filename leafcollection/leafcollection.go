package leafcollection

import (
	"strconv"

	"github.com/gonum/floats"
	"github.com/phantomv1989/hclustering/tree"
)

var InsertMode = struct {
	Insert int
	Delete int
}{Insert: 0, Delete: 1}

type LeafData struct {
	Data             []float64
	MatchedPositions map[string]int
}

func FindAllLeafPositions(tree *tree.Tree, leafArray *[]LeafData, result *map[string]int, scoreLimit float64) {
	// leafPos.TreeLevel.TreePos
	leafPos, treeCnt := FindLeaf(tree, leafArray, scoreLimit)
	if leafPos > -1 {
		(*result)[strconv.Itoa(leafPos)+"."+tree.GetPositionString()] = treeCnt
	}
	for i := range tree.Children {
		FindAllLeafPositions(tree.Children[i], leafArray, result, scoreLimit)
	}
}

func InsertLeafCollectionRecursive(tree *tree.Tree, leafArray *[]LeafData, scoreLimit float64, insertMode int) {
	FindInsertLeaf(tree, leafArray, scoreLimit, insertMode)
	for i := range tree.Children {
		InsertLeafCollectionRecursive(tree.Children[i], leafArray, scoreLimit, insertMode)
	}
}

func FindInsertLeaf(tree *tree.Tree, leafArray *[]LeafData, scoreLimit float64, insertMode int) int {
	leafPos, treePosCnt := FindLeaf(tree, leafArray, scoreLimit)
	if leafPos == -1 {
		ld := LeafData{
			Data:             tree.Leaf,
			MatchedPositions: map[string]int{},
		}
		*leafArray = append(*leafArray, ld)
		leafPos = len(*leafArray) - 1
	}
	if treePosCnt == 0 {
		(*leafArray)[leafPos].MatchedPositions[tree.GetPositionString()] = 1
		if insertMode == 1 {
			(*leafArray)[leafPos].MatchedPositions[tree.GetPositionString()] = -1
		}
	} else {
		if insertMode == 0 {
			(*leafArray)[leafPos].MatchedPositions[tree.GetPositionString()]++
		} else {
			(*leafArray)[leafPos].MatchedPositions[tree.GetPositionString()]--
		}

	}
	return leafPos
}

func FindLeaf(tree *tree.Tree, leafArray *[]LeafData, scoreLimit float64) (int, int) {
	for i := range *leafArray {
		s := GetScore(tree.Leaf, (*leafArray)[i].Data)
		if s < scoreLimit {
			if v, e := (*leafArray)[i].MatchedPositions[tree.GetPositionString()]; e {
				return i, v
			}
			return i, 0
		}
	}
	return -1, 0
}

func GetScore(a, b []float64) float64 {
	return floats.Distance(a, b, 1)
}
