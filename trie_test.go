package main

import (
	"reflect"
	"testing"
	"os"
)

// Test adding a word to the Trie
func TestAddWord(t *testing.T) {
	trie := NewTrie()
	trie.addWord("cat")

	node := trie.walkWord("cat")
	if node == nil || !node.IsEnd {
		t.Errorf("Expected to find the word 'cat' in the Trie, but it was not found")
	}
}

// Test adding multiple words to the Trie
func TestAddMultipleWords(t *testing.T) {
	trie := NewTrie()
	trie.addWord("cat")
	trie.addWord("car")

	nodeCat := trie.walkWord("cat")
	if nodeCat == nil || !nodeCat.IsEnd {
		t.Errorf("Expected to find the word 'cat' in the Trie, but it was not found")
	}

	nodeCar := trie.walkWord("car")
	if nodeCar == nil || !nodeCar.IsEnd {
		t.Errorf("Expected to find the word 'car' in the Trie, but it was not found")
	}
}

// Test adding a duplicate word to the Trie
func TestAddDuplicateWord(t *testing.T) {
	trie := NewTrie()
	trie.addWord("dog")
	trie.addWord("dog") // Duplicate addition

	node := trie.walkWord("dog")
	if node == nil || !node.IsEnd {
		t.Errorf("Expected to find the word 'dog' in the Trie, but it was not found")
	}
}

// Test removing a word from the Trie
func TestRemoveWord(t *testing.T) {
	trie := NewTrie()
	trie.addWord("cat")
	trie.removeWord("cat")

	node := trie.walkWord("cat")
	if node != nil && node.IsEnd {
		t.Errorf("Expected 'cat' to be removed from the Trie, but it still exists")
	}
}

// Test removing a word that doesn't exist
func TestRemoveNonExistentWord(t *testing.T) {
	trie := NewTrie()
	trie.addWord("cat")
	trie.removeWord("dog") // Try removing a non-existent word

	node := trie.walkWord("cat")
	if node == nil || !node.IsEnd {
		t.Errorf("Expected the word 'cat' to still exist, but it was removed")
	}
}

// Test removing a word that is a prefix of another word
func TestRemovePrefixWord(t *testing.T) {
	trie := NewTrie()
	trie.addWord("cat")
	trie.addWord("car")
	trie.removeWord("cat") // "car" should still exist

	nodeCar := trie.walkWord("car")
	if nodeCar == nil || !nodeCar.IsEnd {
		t.Errorf("Expected the word 'car' to still exist, but it was removed")
	}

	nodeCat := trie.walkWord("cat")
	if nodeCat != nil && nodeCat.IsEnd {
		t.Errorf("Expected 'cat' to be removed from the Trie, but it still exists")
	}
}

// Test retrieving words with a common prefix
func TestLoadWordsFromPrefix(t *testing.T) {
	trie := NewTrie()
	trie.addWord("cat")
	trie.addWord("car")
	trie.addWord("cart")

	words := trie.loadWordsFromPrefix("ca")
	expected := []string{"cat", "car", "cart"}

	if !reflect.DeepEqual(words, expected) {
		t.Errorf("Expected words %v, but got %v", expected, words)
	}
}

// Test retrieving words with a non-existent prefix
func TestLoadWordsFromNonExistentPrefix(t *testing.T) {
	trie := NewTrie()
	trie.addWord("cat")
	trie.addWord("car")

	words := trie.loadWordsFromPrefix("dog")
	if len(words) != 0 {
		t.Errorf("Expected no words for prefix 'dog', but got %v", words)
	}
}

// Test walking through a word in the Trie
func TestWalkWord(t *testing.T) {
	trie := NewTrie()
	trie.addWord("bat")

	node := trie.walkWord("bat")
	if node == nil || !node.IsEnd {
		t.Errorf("Expected to walk to the word 'bat', but it was not found")
	}
}

// Test walking a non-existent word in the Trie
func TestWalkNonExistentWord(t *testing.T) {
	trie := NewTrie()
	trie.addWord("bat")

	node := trie.walkWord("cat")
	if node != nil {
		t.Errorf("Expected to not find the word 'cat', but it was found")
	}
}

// Test if the node is barren (no Children, not an end of word)
func TestIsBarren(t *testing.T) {
	trie := NewTrie()
	node := &TrieNode{Children: make(map[rune]*TrieNode), IsEnd: false}

	if !trie.isBarren(node) {
		t.Errorf("Expected node to be barren, but it was not")
	}
}

// Test if a node is not barren (has Children or is the end of a word)
func TestIsNotBarren(t *testing.T) {
	trie := NewTrie()
	node := &TrieNode{Children: make(map[rune]*TrieNode), IsEnd: true}

	if trie.isBarren(node) {
		t.Errorf("Expected node to not be barren, but it was barren")
	}
}

// Test saving a Trie to disk
func TestSaveTrie(t *testing.T) {
	trie := NewTrie()
	trie.addWord("cat")
	trie.addWord("dog")

	// Save the Trie to disk
	Save(trie, "test_trie.gob")

	// Check if the file was created
	if _, err := os.Stat("./.pm/trie/test_trie.gob"); os.IsNotExist(err) {
		t.Errorf("Expected the file 'test_trie.gob' to exist, but it does not")
	}
}

// Test loading a Trie from disk
func TestLoadTrie(t *testing.T) {
	trie := NewTrie()
	trie.addWord("cat")
	trie.addWord("dog")

	// Save the Trie first
	Save(trie, "test_trie.gob")

	// Load the Trie from the file
	loadedTrie := Load("test_trie.gob")
	if loadedTrie == nil {
		t.Errorf("Failed to load the Trie from the file")
	}

	// Check that the loaded Trie contains the words
	if !loadedTrie.walkWord("cat").IsEnd {
		t.Errorf("Expected the word 'cat' to be in the loaded Trie, but it was not found")
	}

	if !loadedTrie.walkWord("dog").IsEnd {
		t.Errorf("Expected the word 'dog' to be in the loaded Trie, but it was not found")
	}
}

// Clean up test files after running tests
func TestCleanupFiles(t *testing.T) {
	err := os.Remove("./.pm/trie/test_trie.gob")
	if err != nil {
		t.Errorf("Error while cleaning up test file: %v", err)
	}
}
