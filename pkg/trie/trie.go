package trie;

import (
	"encoding/gob"
	"errors"
	"fmt"
	"os"

	"github/pm/pkg/common"
)

type TrieNode struct {
	Character rune
	IsEnd     bool
	Children  map[rune]*TrieNode
	Parents   []*TrieNode
	Value     string
}

type Trie struct {
	Id   string
	Root *TrieNode
}

func NewReconcilableTrie(storageKey string) common.Reconcilable {
	trieAlphaList := common.NewAlphaList()
	trieStorage := NewTrie(storageKey)

	return common.Reconcilable{
		AlphaList:     trieAlphaList,
		DataStructure: trieStorage,
	}
}

func NewTrie(id string) *Trie {
	return &Trie{
		Id: id,
		Root: &TrieNode{
			Character: 0,
			IsEnd:     false,
			Children:  make(map[rune]*TrieNode),
			Parents:   []*TrieNode{},
			Value:     "",
		},
	}
}

func (t *Trie) Update(alpha common.Alpha) (error) {
	alphaType := alpha.GetType()
	switch alphaType {
	case common.AddTrieNodeAlpha:
		addTrieNodeAlpha := alpha.(AddTrieNodeAlpha)
		t.AddFile(addTrieNodeAlpha.fileName, addTrieNodeAlpha.fileLocation)
	case common.RemoveTrieNodeAlpha:
		removeTrieNodeAlpha := alpha.(RemoveTrieNodeAlpha)
		t.RemoveWord(removeTrieNodeAlpha.fileName)
	}

	return nil
}

func (t *Trie) Rewind(alpha common.Alpha) (error) {
	alphaType := alpha.GetType()
	switch alphaType {
	case common.AddTrieNodeAlpha:
		addTrieNodeAlpha := alpha.(AddTrieNodeAlpha)
		t.RemoveWord(addTrieNodeAlpha.fileName)
	case common.RemoveTrieNodeAlpha:
		removeTrieNodeAlpha:= alpha.(RemoveTrieNodeAlpha)
		t.AddFile(removeTrieNodeAlpha.fileName, removeTrieNodeAlpha.fileLocation)
	}

	return nil
}

func (t *Trie) Validate(alpha common.Alpha) (bool) {
	alphaType := alpha.GetType()

	switch alphaType {
	case common.AddTrieNodeAlpha:
		addTrieNodeAlpha := alpha.(AddTrieNodeAlpha)
		node, error := t.RetrieveValue(addTrieNodeAlpha.fileName)
		if error != nil && node == "" {
			return false
		}

		return true
	case common.RemoveTrieNodeAlpha:
		addTrieNodeAlpha := alpha.(AddTrieNodeAlpha)
		node, error := t.RetrieveValue(addTrieNodeAlpha.fileName)
		if error != nil && node == "" {
			return true 
		}

		return false 
	}

	return false 
}

type AddTrieNodeAlpha struct {
fileName string
fileLocation string
}

type RemoveTrieNodeAlpha struct {
fileName string
fileLocation string
}

func (atna AddTrieNodeAlpha) GetType() byte {
	return common.AddTrieNodeAlpha
}

func (atna AddTrieNodeAlpha) GetId() string {
	return atna.fileName + atna.fileLocation + string(common.AddTrieNodeAlpha)
}

func (rtna RemoveTrieNodeAlpha) GetType() byte {
	return common.RemoveTrieNodeAlpha
}

func (rtna RemoveTrieNodeAlpha) GetId() string {
	return rtna.fileName + string(common.AddTrieNodeAlpha)
}

func (t *Trie) AddFile(fileName string, fileLocation string) error {
	currentNode := t.Root

	if len(fileName) == 0 {
		return errors.New("Empty fileName not allowed")
	}

	for _, charac := range fileName {
		if _, exist := currentNode.Children[charac]; !exist {
			currentNode.Children[charac] = &TrieNode{
				Character: charac,
				IsEnd:     false,
				Children:  make(map[rune]*TrieNode),
			}

			childNode := currentNode.Children[charac]
			_ = append(childNode.Parents, currentNode)
		}

		currentNode = currentNode.Children[charac]
	}

	currentNode.IsEnd = true
	currentNode.Value = fileLocation

	return nil
}

func (t *Trie) RemoveWord(word string) {
	baseTrie := t.WalkWord(word)
	if baseTrie == nil {
		fmt.Println("Cannot remove non-existant word")
		return
	}

	if baseTrie.IsEnd {
		baseTrie.IsEnd = false
	} else {
		fmt.Println("Word does not exist")
		return
	}

	if len(baseTrie.Children) != 0 {
		return
	}

	for _, parent := range baseTrie.Parents {
		t.RemoveParents(parent, baseTrie)
	}
}

func (t *Trie) RemoveParents(parent *TrieNode, child *TrieNode) {
	_, childExist := parent.Children[child.Character]
	if childExist {
		delete(parent.Children, child.Character)
	}

	if parent.IsEnd {
		return
	}

	if len(parent.Children) > 0 {
		return
	}

	for _, grandParent := range parent.Parents {
		t.RemoveParents(grandParent, parent)
	}
}

func (t *Trie) IsBarren(node *TrieNode) bool {
	return (len(node.Children) == 0) && !node.IsEnd
}

// TODO: Remove the duplicated code
func (t *Trie) LoadWordsFromPrefix(prefix string) []string {
	var words []string
	baseTrie := t.WalkWord(prefix)

	if baseTrie == nil {
		fmt.Println("No words from prefix")
		return words
	}

	if baseTrie.IsEnd {
		words = append(words, prefix)
	}

	if len(baseTrie.Children) == 0 {
		return words
	}

	for _, child := range baseTrie.Children {
		BuildWordsFromChildren(prefix, child, &words)
	}

	return words
}

func BuildWordsFromChildren(base string, node *TrieNode, words *[]string) {
	newBase := base + string(node.Character)
	if node.IsEnd {
		*words = append(*words, newBase)
	}

	if len(node.Children) == 0 {
		return
	}

	for _, child := range node.Children {
		BuildWordsFromChildren(newBase, child, words)
	}
}

func (t *Trie) WalkWord(word string) *TrieNode {
	currentNode := t.Root
	if currentNode == nil {
		return nil
	}

	for _, charac := range word {
		if _, characExist := currentNode.Children[charac]; !characExist {
			return nil
		}

		currentNode = currentNode.Children[charac]
	}

	return currentNode
}

func (t *Trie) RetrieveValue(word string) (string, error) {
	wordNode := t.WalkWord(word)
	if wordNode != nil {
		return wordNode.Value, nil
	}

	err := errors.New("Cannot retrieve value")
	return "", err
}

func (t *Trie) LoadAllWords() ([]string, error) {
	baseTrie := t.Root
	var words []string
	prefix := ""

	for _, child := range baseTrie.Children {
		BuildWordsFromChildren(prefix, child, &words)
	}

	return words, nil
}

func (t *Trie) UpdateValue(word string, fileLocation string) error {
	wordNode := t.WalkWord(word)
	if wordNode != nil {
		wordNode.Value = fileLocation
		return nil
	}

	err := errors.New("Cannot update fileLocation")
	return err
}

func (t *Trie) Save() {
	file, err := os.Create("./.pm/trie/" + t.Id)
	if err != nil {
		fmt.Printf(err.Error())
		fmt.Println("Error creating file")
		return
	}
	defer file.Close()

	encoder := gob.NewEncoder(file)
	encodingErr := encoder.Encode(t)
	if encodingErr != nil {
		fmt.Println("Error encoding trie")
		return
	}
}

func Load(fileName string) *common.Reconcilable {
	path := "./.pm/trie/" + fileName
	file, fileErr := os.Open(path)
	if fileErr != nil {
		return nil
	}
	defer file.Close()

	info, err := os.Stat(path)
	if err != nil {
		return nil
	}

	if info.Size() == 0 {
		newTrie := NewReconcilableTrie(fileName)
		newTrie.SaveReconcilable(path)
		return &newTrie
	}

	decoder := gob.NewDecoder(file)

	var loadedTrie common.Reconcilable
	decodingErr := decoder.Decode(&loadedTrie)
	if decodingErr != nil {
		fmt.Println("Error decoding trie")
		return nil
	}

	return &loadedTrie
}
