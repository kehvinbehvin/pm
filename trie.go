package main

import (
	"encoding/gob"
	"errors"
	"fmt"
	"os"
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

func (t *Trie) addFile(fileName string, fileLocation string) error {
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

func (t *Trie) removeWord(word string) {
	baseTrie := t.walkWord(word)
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
		t.removeParents(parent, baseTrie)
	}
}

func (t *Trie) removeParents(parent *TrieNode, child *TrieNode) {
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
		t.removeParents(grandParent, parent)
	}
}

func (t *Trie) isBarren(node *TrieNode) bool {
	return (len(node.Children) == 0) && !node.IsEnd
}

// TODO: Remove the duplicated code
func (t *Trie) loadWordsFromPrefix(prefix string) []string {
	var words []string
	baseTrie := t.walkWord(prefix)

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
		buildWordsFromChildren(prefix, child, &words)
	}

	return words
}

func buildWordsFromChildren(base string, node *TrieNode, words *[]string) {
	newBase := base + string(node.Character)
	if node.IsEnd {
		*words = append(*words, newBase)
	}

	if len(node.Children) == 0 {
		return
	}

	for _, child := range node.Children {
		buildWordsFromChildren(newBase, child, words)
	}
}

func (t *Trie) walkWord(word string) *TrieNode {
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

func (t *Trie) retrieveValue(word string) (string, error) {
	wordNode := t.walkWord(word)
	if wordNode != nil {
		return wordNode.Value, nil
	}

	err := errors.New("Cannot retrieve value")
	return "", err
}

func (t *Trie) loadAllWords() ([]string, error){
	baseTrie := t.Root;
	var words []string
	prefix := ""

	for _, child := range baseTrie.Children {
		buildWordsFromChildren(prefix, child, &words)
	}
	
	return words, nil
}

func (t *Trie) updateValue(word string, fileLocation string) error {
	wordNode := t.walkWord(word)
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

func Load(fileName string) *Trie {
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
		newTrie := NewTrie(fileName)
		newTrie.Save()
		return newTrie
	}

	decoder := gob.NewDecoder(file)

	var loadedTrie Trie
	decodingErr := decoder.Decode(&loadedTrie)
	if decodingErr != nil {
		fmt.Println("Error decoding trie")
		return nil
	}

	return &loadedTrie
}
