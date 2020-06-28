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

// LeafData uses cumulative count for leaf key
type LeafData struct {
	Data []float64
}

// SaveLeafCollection "./myfile"
func SaveLeafCollection(fileName string, ld []LeafData) {
	ldata, err := json.Marshal(ld)
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile(fileName, ldata, 0644)
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
	findAllLeafPos(prefix, tree, leafArray, &result, scoreLimit)
	return result
}

// InsertLeavesRecursive inserts all leaves of a decomposed tree to an array
func InsertLeavesRecursive(prefix string, tree *tree.Tree, leafArray *[]LeafData, scoreLimit float64) {
	FindInsertLeaf(prefix, tree, leafArray, scoreLimit)
	for i := range tree.Children {
		InsertLeavesRecursive(prefix+"."+strconv.Itoa(i), tree.Children[i], leafArray, scoreLimit)
	}
}

func FindInsertLeaf(prefix string, tree *tree.Tree, leafArray *[]LeafData, scoreLimit float64) int {
	leafPos := findLeaf(prefix, tree, leafArray, scoreLimit)
	if leafPos == -1 {
		ld := LeafData{
			Data: tree.Leaf,
		}
		*leafArray = append(*leafArray, ld)
		leafPos = len(*leafArray) - 1
	}
	return leafPos
}

func findAllLeafPos(prefix string, tree *tree.Tree, leafArray *[]LeafData, result *map[string]int, scoreLimit float64) {
	leafPos := findLeaf(prefix, tree, leafArray, scoreLimit)
	(*result)[prefix] = leafPos
	for i := range tree.Children {
		findAllLeafPos(prefix+"."+strconv.Itoa(i), tree.Children[i], leafArray, result, scoreLimit)
	}
}

func findLeaf(prefix string, tree *tree.Tree, leafArray *[]LeafData, scoreLimit float64) int {
	for i := range *leafArray {
		s := getScore(tree.Leaf, (*leafArray)[i].Data)
		if s < scoreLimit {
			return i
		}
	}
	return -1
}

func getScore(a, b []float64) float64 {
	return floats.Distance(a, b, 2)
}
