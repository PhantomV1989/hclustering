package leafcollection

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strconv"

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

// // GetHierarchicalMatches returns all matched leaves in a hierarchical structure
// func GetHierarchicalMatches(prefix string, tree *tree.Tree, leafArray *[]LeafData, scoreLimit float64, withLeafPos bool) map[int]map[string]int {
// 	allpos, allscore := FindAllLeafPositions(prefix, tree, leafArray, withLeafPos)
// 	result := map[int]map[string]int{}
// 	for i := range allpos {
// 		k := strings.Count(i, ".")
// 		if _, e := result[k]; !e {
// 			result[k] = map[string]int{}
// 		}
// 		result[k][i] = allpos[i]
// 	}
// 	return result
// }

// FindAllLeafPositions returns all matched leaves in a flat format, Key format as follows Tree.Path:LeafPos
func FindAllLeafPositions(prefix string, tree *tree.Tree, leafArray *[]LeafData, withLeafPos bool) (map[string]int, map[string]float64) {
	result := map[string]int{}
	rscore := map[string]float64{}
	findAllLeafPos(prefix, tree, leafArray, &result, &rscore)
	return result, rscore
}

// InsertLeavesRecursive inserts all leaves of a decomposed tree to an array
func InsertLeavesRecursive(prefix string, tree *tree.Tree, leafArray *[]LeafData, scoreLimit float64) {
	FindInsertLeaf(tree, leafArray, scoreLimit)
	for i := range tree.Children {
		InsertLeavesRecursive(prefix+"."+strconv.Itoa(i), tree.Children[i], leafArray, scoreLimit)
	}
}

// FindInsertLeaf find, insert new leaf if not found
func FindInsertLeaf(tree *tree.Tree, leafArray *[]LeafData, scoreLimit float64) int {
	leafPos, score := findLeaf(tree, leafArray)
	if score > scoreLimit {
		ld := LeafData{
			Data: tree.Leaf,
		}
		*leafArray = append(*leafArray, ld)
		leafPos = len(*leafArray) - 1
	}
	return leafPos
}

func findAllLeafPos(prefix string, tree *tree.Tree, leafArray *[]LeafData, result *map[string]int, rscore *map[string]float64) {
	leafPos, score := findLeaf(tree, leafArray)
	(*result)[prefix] = leafPos
	(*rscore)[prefix] = score
	for i := range tree.Children {
		findAllLeafPos(prefix+"."+strconv.Itoa(i), tree.Children[i], leafArray, result, rscore)
	}
}

func findLeaf(tree *tree.Tree, leafArray *[]LeafData) (int, float64) {
	bestMatchPos := 0
	bestMatchScore := getScore(tree.Leaf, (*leafArray)[0].Data)

	for i := range *leafArray {
		s := getScore(tree.Leaf, (*leafArray)[i].Data)
		if s < bestMatchScore {
			bestMatchPos = i
			bestMatchScore = s
		}
	}
	return bestMatchPos, bestMatchScore
}

func getScore(a, b []float64) float64 {
	return floats.Distance(a, b, 2)
}
