package common

import (
	"encoding/gob"
	"errors"
	"log"
	"os"
)

type AlphaList struct {
	Alphas   []Alpha
	Branches map[string]int
}

func NewAlphaList() AlphaList {
	return AlphaList{
		Alphas:   []Alpha{},
		Branches: make(map[string]int),
	}
}

func (al AlphaList) MergeIn() {
}

func (al AlphaList) Diff() {
}

type Reconcilable struct {
	AlphaList     AlphaList
	DataStructure DataStructure
	FilePath      string
}

// Rewind dataStructure to state at Alpha
// Remove all alphas after Alpha input
func (r Reconcilable) Reset(input Alpha) error {
	incumbentLength := len(r.AlphaList.Alphas)
	alphasToRewind := []Alpha{}

	alphaFound := false
	inputAlphaIndex := 0
	for index := incumbentLength; index > 0; index-- {
		if input.GetId() == r.AlphaList.Alphas[index].GetId() {
			alphaFound = true
			inputAlphaIndex = index
			break
		}

		_ = append(alphasToRewind, r.AlphaList.Alphas[index])
	}

	if !alphaFound {
		return errors.New("Alpha cannot be found, Reset failed")
	}

	r.AlphaList.Alphas = r.AlphaList.Alphas[:inputAlphaIndex+1]

	for index := 0; index < len(alphasToRewind); index++ {
		return r.DataStructure.Rewind(alphasToRewind[index])
	}

	return nil
}

// Update DataStructure using Alpha
// Append Alpha to AlphaList
func (r Reconcilable) Commit(input Alpha) error {
	return r.DataStructure.Update(input)
}

// Appends a list of Alphas to AlphaList
// Updates the DataStructure with the list of Alphas in order
func (r Reconcilable) FastForward(input []Alpha) error {
	for index := len(input); index > 0; index-- {
		error := r.DataStructure.Update(input[index])

		if error != nil {
			return error
		}
	}

	return nil
}

func (r Reconcilable) SaveReconcilable() {
	file, err := os.Create(r.FilePath)
	if err != nil {
		log.Printf(err.Error())
		log.Println("Error creating file")
		return
	}
	defer file.Close()

	encoder := gob.NewEncoder(file)
	gob.Register(r.DataStructure)
	encodingErr := encoder.Encode(r)
	if encodingErr != nil {
		log.Println("Error encoding dag", encodingErr.Error())
		return
	}
}

func LoadReconcilable(filePath string) *Reconcilable {
	file, fileErr := os.Open(filePath)

	if fileErr != nil {
		log.Println("Error opening binary file")
		return nil
	}
	defer file.Close()

	decoder := gob.NewDecoder(file)
	var loadedReconcilable *Reconcilable
	decodingErr := decoder.Decode(&loadedReconcilable)
	if decodingErr != nil {
		log.Println("Error decoding", decodingErr.Error())
		return nil
	}

	return loadedReconcilable
}

// type DepNode struct {
// 	Parents []Alpha
// 	Children []Alpha
// }
//
// type DepTree struct {
// 	Tree map[string]DepNode
// }
//
// func NewDepTree(alphas []Alpha) *DepTree {
// 	tree := make(map[string]DepNode)
//
// 	for i := 0; i < len(alphas); i++ {
// 		dependencies := alphas[i].GetAlphaDependencies()
// 		var emptyChildren []Alpha
//
// 		id := alphas[i].GetId()
//
// 		node := DepNode{
// 			Parents: dependencies,
// 			Children: emptyChildren,
// 		}
//
// 		tree[id] = node
//
// 		// Add children for existing parents
// 		for dependency := 0; i < len(dependencies); i++ {
// 			parentId := dependencies[dependency].GetId()
// 			parent, ok := tree[parentId];
// 			if !ok {
// 				continue
// 			}
//
// 			_ = append(parent.Children, alphas[i])
// 		}
// 	}
//
// 	return &DepTree{
// 		Tree: make(map[string]DepNode),
// 	}
// }
