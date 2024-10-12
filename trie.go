package main

import "fmt"

type TrieNode struct {
  character rune
  isEnd bool
  children map[rune]*TrieNode
  parents []*TrieNode
}

type Trie struct {
  root *TrieNode
}

func newTrie() *Trie {
  return &Trie{
    root: &TrieNode{
      character: 0,
      isEnd: false,
      children: make(map[rune]*TrieNode),
      parents: []*TrieNode{},
    },
  }
}

func (t *Trie) addWord(word string) {
  currentNode := t.root

  for _ , charac := range word {
    if _ , exist := currentNode.children[charac]; !exist {
      currentNode.children[charac] = &TrieNode{
        character: charac,
        isEnd: false,
        children: make(map[rune]*TrieNode),
      }

      childNode := currentNode.children[charac]
      _ = append(childNode.parents, currentNode)
    }

    currentNode = currentNode.children[charac]
  }
  
  currentNode.isEnd = true
}

func (t *Trie) removeWord(word string) {
  baseTrie := t.walkWord(word);
  if baseTrie == nil {
    fmt.Println("Cannot remove non-existant word")
    return
  }

  if baseTrie.isEnd {
    baseTrie.isEnd = false
  } else {
    fmt.Println("Word does not exist")
    return
  }

  if len(baseTrie.children) != 0 {
    return
  }

  for _, parent := range baseTrie.parents {
    t.removeParents(parent, baseTrie)
  }
}

func (t *Trie) removeParents(parent *TrieNode, child *TrieNode) {
  _, childExist := parent.children[child.character]
  if childExist {
    delete(parent.children, child.character)
  }

  if (parent.isEnd) {
    return
  }

  if len(parent.children) > 0 {
    return
  }

  for _, grandParent := range parent.parents {
    t.removeParents(grandParent, parent)
  }
}

func (t *Trie) isBarren(node *TrieNode) bool {
  return (len(node.children) == 0) && !node.isEnd
}

// TODO: Remove the duplicated code
func (t *Trie) loadWordsFromPrefix(prefix string) []string {
  var words []string
  baseTrie := t.walkWord(prefix);
  
  if baseTrie == nil {
    fmt.Println("No words from prefix")
    return words
  }

  if baseTrie.isEnd {
    words = append(words, prefix)
  }

  if len(baseTrie.children) == 0 {
    return words
  }

  for _, child := range baseTrie.children {
    buildWordsFromChildren(prefix, child, &words);
  }

  return words
}

func buildWordsFromChildren(base string, node *TrieNode, words *[]string) {
  newBase := base + string(node.character)
  if node.isEnd {
    *words = append(*words, newBase)
  }

  if len(node.children) == 0 {
    return
  }
  
  for _ , child := range node.children {
    buildWordsFromChildren(newBase, child, words)
  }
}

func (t *Trie) walkWord(word string) *TrieNode {
  currentNode := t.root
  for _, charac := range word {
    if _ , characExist := currentNode.children[charac]; !characExist {
      return nil
    }

    currentNode = currentNode.children[charac]
  }

  return currentNode
}
