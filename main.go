package main

import (
	// "github/pm/cmd"
	"fmt"
)

func main() {
	// cmd.Execute()
	pmDag := newDag()
	vertex1 := newVertex("Hi")
	pmDag.addVertex(vertex1)

	search := NewTrie()

	word1 := "hello"
	word2 := "world"
	word3 := "help"
	word4 := "worlly"
	word5 := "worlda"

	search.addWord(word1)
	search.addWord(word2)
	search.addWord(word3)
	search.addWord(word4)
	search.addWord(word5)

	// search.removeWord("world")
	suggestions := search.loadWordsFromPrefix("wo")
	for _, value := range suggestions {
		fmt.Println(value)
	}
	
	Save(search, "epic")

	loadedTrie := Load("epic")
	if loadedTrie == nil {
		fmt.Println("Error loading trie")
		return
	}

	loadedSuggestions := loadedTrie.loadWordsFromPrefix("h")
	for _, value := range loadedSuggestions {
		fmt.Println(value)
	}
}
