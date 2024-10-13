package main

import "fmt"

type TrieNode struct {
  Character rune
  IsEnd bool
  Children map[rune]*TrieNode
  Parents []*TrieNode
}

type Trie struct {
  Root *TrieNode
}

func NewTrie() *Trie {
  return &Trie{
    Root: &TrieNode{
      Character: 0,
      IsEnd: false,
      Children: make(map[rune]*TrieNode),
      Parents: []*TrieNode{},
    },
  }
}

func (t *Trie) addWord(word string) {
  currentNode := t.Root

  for _ , charac := range word {
    if _ , exist := currentNode.Children[charac]; !exist {
      currentNode.Children[charac] = &TrieNode{
        Character: charac,
        IsEnd: false,
        Children: make(map[rune]*TrieNode),
      }

      childNode := currentNode.Children[charac]
      _ = append(childNode.Parents, currentNode)
    }

    currentNode = currentNode.Children[charac]
  }
  
  currentNode.IsEnd = true
}

func (t *Trie) removeWord(word string) {
  baseTrie := t.walkWord(word);
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

  if (parent.IsEnd) {
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
  baseTrie := t.walkWord(prefix);
  
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
    buildWordsFromChildren(prefix, child, &words);
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
  
  for _ , child := range node.Children {
    buildWordsFromChildren(newBase, child, words)
  }
}

func (t *Trie) walkWord(word string) *TrieNode {
  currentNode := t.Root
  for _, charac := range word {
    if _ , characExist := currentNode.Children[charac]; !characExist {
      return nil
    }

    currentNode = currentNode.Children[charac]
  }

  return currentNode
}
