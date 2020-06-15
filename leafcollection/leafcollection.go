package leafcollection

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/gonum/floats"
	"github.com/phantomv1989/hclustering/tree"
)

var InsertMode = struct {
	Insert int
	Delete int
}{Insert: 0, Delete: 1}

// LeafData uses cumulative count for leaf key
type LeafData struct {
	Data             []float64
	MatchedPositions map[string]int
}

// SaveLeafCollection "./myfile"
func SaveLeafCollection(fileName string, ld *[]LeafData) {
	ldata, _ := json.Marshal(ld)
	err := ioutil.WriteFile(fileName, ldata, 0644)
	if err != nil {
		panic(err)
	}
}

// LoadLeafCollection returns empty result if path not exist
func LoadLeafCollection(folderPath string) []LeafData {
	ld := []LeafData{}
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

// GetHierarchicalMatches returns all matched leaves in a hierarchical structure
func GetHierarchicalMatches(prefix string, tree *tree.Tree, leafArray *[]LeafData, scoreLimit float64, withLeafPos bool) map[int]map[string]int {
	allpos := FindAllLeafPositions(prefix, tree, leafArray, scoreLimit, withLeafPos)
	result := map[int]map[string]int{}
	for i := range allpos {
		k := strings.Count(i, ".")
		if _, e := result[k]; !e {
			result[k] = map[string]int{}
		}
		result[k][i] = allpos[i]
	}
	return result
}

// FindAllLeafPositions returns all matched leaves in a flat format, Key format as follows Tree.Path:LeafPos
func FindAllLeafPositions(prefix string, tree *tree.Tree, leafArray *[]LeafData, scoreLimit float64, withLeafPos bool) map[string]int {
	result := map[string]int{}
	findAllLeafPos(prefix, tree, leafArray, &result, scoreLimit, withLeafPos)
	return result
}

// InsertLeavesRecursive inserts all leaves of a decomposed tree to an array
func InsertLeavesRecursive(prefix string, tree *tree.Tree, leafArray *[]LeafData, scoreLimit float64, insertMode int) {
	FindInsertLeaf(prefix, tree, leafArray, scoreLimit, insertMode)
	for i := range tree.Children {
		InsertLeavesRecursive(prefix+"."+strconv.Itoa(i), tree.Children[i], leafArray, scoreLimit, insertMode)
	}
}

func GetFlatLeafCollection(ld *[]LeafData) map[string]int {
	r := map[string]int{}
	for l := range *ld {
		for ll := range (*ld)[l].MatchedPositions {
			r[ll+":"+strconv.Itoa(l)] = (*ld)[l].MatchedPositions[ll]
		}
	}
	return r
}

func PruneLeafCollectionByMinBracket(leafArray []LeafData, minLeafKeyBracketCount int) []LeafData {
	la2 := []LeafData{}
	for l := range leafArray {
		mp := map[string]int{}
		for ll := range leafArray[l].MatchedPositions {
			if !(leafArray[l].MatchedPositions[ll] < minLeafKeyBracketCount && leafArray[l].MatchedPositions[ll] > -minLeafKeyBracketCount) {
				mp[ll] = leafArray[l].MatchedPositions[ll]
			}
		}
		if len(mp) > 0 {
			ldt := LeafData{Data: leafArray[l].Data, MatchedPositions: mp}
			la2 = append(la2, ldt)
		}
	}
	return la2
}

func FindInsertLeaf(prefix string, tree *tree.Tree, leafArray *[]LeafData, scoreLimit float64, insertMode int) int {
	leafPos, treePosCnt := findLeaf(prefix, tree, leafArray, scoreLimit)
	if leafPos == -1 {
		ld := LeafData{
			Data:             tree.Leaf,
			MatchedPositions: map[string]int{},
		}
		*leafArray = append(*leafArray, ld)
		leafPos = len(*leafArray) - 1
	}
	if treePosCnt == 0 {
		(*leafArray)[leafPos].MatchedPositions[prefix] = 1
		if insertMode == 1 {
			(*leafArray)[leafPos].MatchedPositions[prefix] = -1
		}
	} else {
		if insertMode == 0 {
			(*leafArray)[leafPos].MatchedPositions[prefix]++
		} else {
			(*leafArray)[leafPos].MatchedPositions[prefix]--
		}

	}
	return leafPos
}

func findAllLeafPos(prefix string, tree *tree.Tree, leafArray *[]LeafData, result *map[string]int, scoreLimit float64, withLeafPos bool) {
	leafPos, treeCnt := findLeaf(prefix, tree, leafArray, scoreLimit)
	if withLeafPos {
		(*result)[prefix+":"+strconv.Itoa(leafPos)] = treeCnt
	} else {
		(*result)[prefix] = treeCnt
	}

	for i := range tree.Children {
		findAllLeafPos(prefix+"."+strconv.Itoa(i), tree.Children[i], leafArray, result, scoreLimit, withLeafPos)
	}
}

func findLeaf(prefix string, tree *tree.Tree, leafArray *[]LeafData, scoreLimit float64) (int, int) {
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

func getScore(a, b []float64) float64 {
	return floats.Distance(a, b, 2)
}
