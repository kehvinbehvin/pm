package trie

import (
	"crypto/sha1"
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
	filePath := "./.pm/trie/" + storageKey

	return common.Reconcilable{
		AlphaList:     trieAlphaList,
		DataStructure: trieStorage,
		FilePath:      filePath,
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

func (t *Trie) Update(alpha common.Alpha) error {
	alphaType := alpha.GetType()
	switch alphaType {
	case common.AddTrieNodeAlpha:
		addTrieNodeAlpha := alpha.(AddTrieNodeAlpha)
		t.AddFile(addTrieNodeAlpha.FileName, addTrieNodeAlpha.FileLocation)
	case common.RemoveTrieNodeAlpha:
		removeTrieNodeAlpha := alpha.(RemoveTrieNodeAlpha)
		t.RemoveWord(removeTrieNodeAlpha.FileName)
	}

	return nil
}

func (t *Trie) Rewind(alpha common.Alpha) error {
	alphaType := alpha.GetType()
	switch alphaType {
	case common.AddTrieNodeAlpha:
		addTrieNodeAlpha := alpha.(AddTrieNodeAlpha)
		t.RemoveWord(addTrieNodeAlpha.FileName)
	case common.RemoveTrieNodeAlpha:
		removeTrieNodeAlpha := alpha.(RemoveTrieNodeAlpha)
		t.AddFile(removeTrieNodeAlpha.FileName, removeTrieNodeAlpha.FileLocation)
	}

	return nil
}

func (t *Trie) Validate(alpha common.Alpha) bool {
	alphaType := alpha.GetType()

	switch alphaType {
	case common.AddTrieNodeAlpha:
		addTrieNodeAlpha := alpha.(AddTrieNodeAlpha)
		node, error := t.RetrieveValue(addTrieNodeAlpha.FileName)
		if error != nil && node == "" {
			return false
		}

		return true
	case common.RemoveTrieNodeAlpha:
		addTrieNodeAlpha := alpha.(AddTrieNodeAlpha)
		node, error := t.RetrieveValue(addTrieNodeAlpha.FileName)
		if error != nil && node == "" {
			return true
		}

		return false
	}

	return false
}

type AddTrieNodeAlpha struct {
	FileName     string
	FileLocation string
	Hash         string
}

type RemoveTrieNodeAlpha struct {
	FileName     string
	FileLocation string
	Hash         string
}

func (atna AddTrieNodeAlpha) GetType() byte {
	return common.AddTrieNodeAlpha
}

func (atna AddTrieNodeAlpha) GetId() string {
	return atna.FileName + atna.FileLocation + string(common.AddTrieNodeAlpha)
}

func (atna AddTrieNodeAlpha) GetHash() string {
	return atna.Hash
}

// This is mean to capture the state that the alpha was used to update
// the underlying datastructure
func (atna AddTrieNodeAlpha) SetHash(lastAlpha common.Alpha) {
	prevAlphaHash := lastAlpha.GetHash()
	currentHash := sha1.Sum([]byte(atna.GetId() + prevAlphaHash))
	currentHashStr := fmt.Sprintf("%x", currentHash[:])
	atna.Hash = currentHashStr
}

func (rtna RemoveTrieNodeAlpha) GetType() byte {
	return common.RemoveTrieNodeAlpha
}

func (rtna RemoveTrieNodeAlpha) GetId() string {
	return rtna.FileName + string(common.AddTrieNodeAlpha)
}

func (rtna RemoveTrieNodeAlpha) GetHash() string {
	return rtna.Hash
}

// This is mean to capture the state that the alpha was used to update
// the underlying datastructure
func (rtna RemoveTrieNodeAlpha) SetHash(lastAlpha common.Alpha) {
	prevAlphaHash := lastAlpha.GetHash()
	currentHash := sha1.Sum([]byte(rtna.GetId() + prevAlphaHash))
	currentHashStr := fmt.Sprintf("%x", currentHash[:])
	rtna.Hash = currentHashStr
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
		log.Println("Cannot remove non-existant word")
		return
	}

	if baseTrie.IsEnd {
		baseTrie.IsEnd = false
	} else {
		log.Println("Word does not exist")
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
		log.Println("No words from prefix")
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

func LoadReconcilableTrie(filePath string) *common.Reconcilable {
	file, fileErr := os.Open(filePath)

	if fileErr != nil {
		log.Println("Error opening binary file")
		return nil
	}
	defer file.Close()

	decoder := gob.NewDecoder(file)
	gob.Register(&Trie{})
	var loadedReconcilable *common.Reconcilable
	decodingErr := decoder.Decode(&loadedReconcilable)
	if decodingErr != nil {
		log.Println("Error decoding", decodingErr.Error())
		return nil
	}

	return loadedReconcilable
}
