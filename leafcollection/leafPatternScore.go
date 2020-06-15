package leafcollection

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/phantomv1989/hclustering/tree"
)

// LeafDataPatternScore uses float score for leaf key
type LeafDataPatternScore struct {
	Data             []float64
	MatchedPositions map[string]float64
}

// FindAllLeafPositions returns all matched leaves in a flat format, Key format as follows Tree.Path:LeafPos
func FindAllLeafPositionsPatternScore(prefix string, tree *tree.Tree, leafArray *[]LeafDataPatternScore, scoreLimit float64, withLeafPos bool) map[string]float64 {
	result := map[string]float64{}
	findAllLeafPositionsPatternScore(prefix, tree, leafArray, &result, scoreLimit, withLeafPos)
	return result
}

func findLeafPatternScore(prefix string, tree *tree.Tree, leafArray *[]LeafDataPatternScore, scoreLimit float64) (int, float64) {
	for i := range *leafArray {
		s := getScore(tree.Leaf, (*leafArray)[i].Data)
		if s < scoreLimit {
			if v, e := (*leafArray)[i].MatchedPositions[prefix]; e {
				return i, v
			}
			return i, 0
		}
	}
	return -1, 0
}

// FindInsertLeafPatternScore ..
func FindInsertLeafPatternScore(prefix string, tree *tree.Tree, leafArray *[]LeafDataPatternScore, scoreLimit float64, patternScore float64) int {
	leafPos, _ := findLeafPatternScore(prefix, tree, leafArray, scoreLimit)
	if leafPos == -1 {
		ld := LeafDataPatternScore{
			Data:             tree.Leaf,
			MatchedPositions: map[string]float64{},
		}
		*leafArray = append(*leafArray, ld)
		leafPos = len(*leafArray) - 1
	}
	(*leafArray)[leafPos].MatchedPositions[prefix] += patternScore
	return leafPos
}

// InsertLeavesRecursive inserts all leaves of a decomposed tree to an array
func InsertLeavesRecursivePatternScore(prefix string, tree *tree.Tree, leafArray *[]LeafDataPatternScore, scoreLimit float64, patternScore float64) {
	FindInsertLeafPatternScore(prefix, tree, leafArray, scoreLimit, patternScore)
	for i := range tree.Children {
		InsertLeavesRecursivePatternScore(prefix+"."+strconv.Itoa(i), tree.Children[i], leafArray, scoreLimit, patternScore)
	}
}

// SaveLeafCollectionPatternScore "./myfile"
func SaveLeafCollectionPatternScore(fileName string, ld *[]LeafDataPatternScore) {
	ldata, _ := json.Marshal(ld)
	err := ioutil.WriteFile(fileName, ldata, 0644)
	if err != nil {
		panic(err)
	}
}

// LoadLeafCollectionPatterScore returns empty result if path not exist
func LoadLeafCollectionPatterScore(folderPath string) []LeafDataPatternScore {
	ld := []LeafDataPatternScore{}
	if _, e := os.Stat(folderPath); e == nil {
		someBytes, err := ioutil.ReadFile(folderPath)
		if err != nil {
			panic(err)
		}
		err = json.Unmarshal(someBytes, &ld)
		if err != nil {
			panic(err)
		}
	}
	return ld
}

func findAllLeafPositionsPatternScore(prefix string, tree *tree.Tree, leafArray *[]LeafDataPatternScore, result *map[string]float64, scoreLimit float64, withLeafPos bool) {
	leafPos, patternScore := findLeafPatternScore(prefix, tree, leafArray, scoreLimit)
	if withLeafPos {
		(*result)[prefix+":"+strconv.Itoa(leafPos)] = patternScore
	} else {
		(*result)[prefix] = patternScore
	}

	for i := range tree.Children {
		findAllLeafPositionsPatternScore(prefix+"."+strconv.Itoa(i), tree.Children[i], leafArray, result, scoreLimit, withLeafPos)
	}
}
