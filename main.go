package main

import (
	// "github/pm/cmd"
	"fmt"
	"encoding/gob"
	"os"
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
  for _ , value := range suggestions {
    fmt.Println(value)
  }

  file, err := os.Create("epic");
  if err != nil {
    fmt.Println("Error creating file")
    return
  }
  defer file.Close()

  encoder := gob.NewEncoder(file)
  encodingErr := encoder.Encode(search)
  if encodingErr != nil {
    fmt.Printf(encodingErr.Error())
    fmt.Println("Error encoding trie")
    return
  }

  fmt.Println("File encoded and saved")

  savedFile, fileErr := os.Open("epic")
  if fileErr != nil {
    fmt.Println("Error opening binary file")
    return
  }
  defer file.Close()

  decoder := gob.NewDecoder(savedFile)
  var loadedTrie Trie
  decodingErr := decoder.Decode(&loadedTrie)
  if decodingErr != nil {
    fmt.Println("Error decoding trie")
    return
  }
  
  fmt.Println("Loaded trie from binary")
  loadedSuggestions := loadedTrie.loadWordsFromPrefix("h")
  for _ , value := range loadedSuggestions  {
    fmt.Println(value)
  }
}
