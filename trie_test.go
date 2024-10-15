package main

import (
	"reflect"
	"testing"
	"os"
	"crypto/sha1"
	"fmt"
)

// Test adding a file to the Trie
func TestAddFile(t *testing.T) {
	trie := NewTrie()
	content := "file content for cat"
	hash := sha1.Sum([]byte(content)) // Generate the SHA-1 hash
	hashStr := fmt.Sprintf("%x", hash[:])
	trie.addFile("cat", hashStr)

	node := trie.walkWord("cat")
	if node == nil || !node.IsEnd {
		t.Errorf("Expected to find the file 'cat' in the Trie, but it was not found")
	}

	if !reflect.DeepEqual(node.Value, hashStr) {
		t.Errorf("Expected to retrieve hash %x for 'cat', but got %x", hashStr, node.Value)
	}
}

// Test adding multiple files to the Trie
func TestAddMultipleFiles(t *testing.T) {
	trie := NewTrie()
	contentCat := "file content for cat"
	contentCar := "file content for car"

	cat := sha1.Sum([]byte(contentCat)) // Generate SHA-1 hash for cat
	hashCat := fmt.Sprintf("%x", cat[:])
	car := sha1.Sum([]byte(contentCar)) // Generate SHA-1 hash for car
	hashCar := fmt.Sprintf("%x", car[:])

	trie.addFile("cat", hashCat)
	trie.addFile("car", hashCar)

	nodeCat := trie.walkWord("cat")
	if nodeCat == nil || !nodeCat.IsEnd {
		t.Errorf("Expected to find the file 'cat' in the Trie, but it was not found")
	}

	nodeCar := trie.walkWord("car")
	if nodeCar == nil || !nodeCar.IsEnd {
		t.Errorf("Expected to find the file 'car' in the Trie, but it was not found")
	}

	if !reflect.DeepEqual(nodeCat.Value, hashCat) {
		t.Errorf("Expected to retrieve hash %x for 'cat', but got %x", hashCat, nodeCat.Value)
	}

	if !reflect.DeepEqual(nodeCar.Value, hashCar) {
		t.Errorf("Expected to retrieve hash %x for 'car', but got %x", hashCar, nodeCar.Value)
	}
}

// Test retrieving a file's hash from the Trie
func TestRetrieveFile(t *testing.T) {
	trie := NewTrie()
	content := "file content for dog"
	hash := sha1.Sum([]byte(content)) // Generate the SHA-1 hash
	hashHex := fmt.Sprintf("%x", hash[:])
	trie.addFile("dog", hashHex)

	retrievedHash, err := trie.retrieveValue("dog")
	if err != nil {
		t.Errorf("Error retrieving hash for 'dog': %v", err)
	}

	if !reflect.DeepEqual(retrievedHash, hashHex) {
		t.Errorf("Expected to retrieve hash %x for 'dog', but got %x", hashHex, retrievedHash)
	}
}


// Test adding a duplicate file to the Trie
func TestAddDuplicateFile(t *testing.T) {
	trie := NewTrie()
	content := "file content for dog"
	hash := sha1.Sum([]byte(content)) // Generate the SHA-1 hash
	hashHex := fmt.Sprintf("%x", hash[:])

	trie.addFile("dog", hashHex)
	trie.addFile("dog", hashHex) // Duplicate addition

	node := trie.walkWord("dog")
	if node == nil || !node.IsEnd {
		t.Errorf("Expected to find the file 'dog' in the Trie, but it was not found")
	}

	if !reflect.DeepEqual(node.Value, hashHex) {
		t.Errorf("Expected to retrieve hash %x for 'dog', but got %x", hashHex, node.Value)
	}
}

// Test removing a file from the Trie
func TestRemoveFile(t *testing.T) {
	trie := NewTrie()
	content := "file content for cat"
	hash := sha1.Sum([]byte(content)) // Generate the SHA-1 hash
	hashHex := fmt.Sprintf("%x", hash[:])
	trie.addFile("cat", hashHex)

	trie.removeWord("cat")

	node := trie.walkWord("cat")
	if node != nil && node.IsEnd {
		t.Errorf("Expected 'cat' to be removed from the Trie, but it still exists")
	}
}

// Test removing a word that doesn't exist
func TestRemoveNonExistentWord(t *testing.T) {
	trie := NewTrie()
  content := "file content for cat"
  hash := sha1.Sum([]byte(content)) // Generate the SHA-1 hash
  hashHex := fmt.Sprintf("%x", hash[:])
	trie.addFile("cat", hashHex)
	trie.removeWord("dog") // Try removing a non-existent word

	node := trie.walkWord("cat")
	if node == nil || !node.IsEnd {
		t.Errorf("Expected the word 'cat' to still exist, but it was removed")
	}
}

// Test removing a word that is a prefix of another word
func TestRemovePrefixWord(t *testing.T) {
	trie := NewTrie()
	contentCat := "file content for cat"
  contentCar := "file content for car"

  cat := sha1.Sum([]byte(contentCat)) // Generate SHA-1 hash for cat
  hashCat := fmt.Sprintf("%x", cat[:])
  car := sha1.Sum([]byte(contentCar)) // Generate SHA-1 hash for car
  hashCar := fmt.Sprintf("%x", car[:])

	trie.addFile("cat", hashCat)
	trie.addFile("car", hashCar)
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
	contentCat := "file content for cat"
  contentCar := "file content for car"
  contentCart := "file content for cart"

  cat := sha1.Sum([]byte(contentCat)) // Generate SHA-1 hash for cat
  hashCat := fmt.Sprintf("%x", cat[:])
  car := sha1.Sum([]byte(contentCar)) // Generate SHA-1 hash for car
  hashCar := fmt.Sprintf("%x", car[:])
  cart := sha1.Sum([]byte(contentCart)) // Generate SHA-1 hash for car
  hashCart := fmt.Sprintf("%x", cart[:])
	trie.addFile("cat", hashCat)
	trie.addFile("car", hashCar)
	trie.addFile("cart", hashCart)

	words := trie.loadWordsFromPrefix("ca")
	expected := []string{"cat", "car", "cart"}

	if !reflect.DeepEqual(words, expected) {
		t.Errorf("Expected words %v, but got %v", expected, words)
	}
}

// Test retrieving words with a non-existent prefix
func TestLoadWordsFromNonExistentPrefix(t *testing.T) {
	trie := NewTrie()
	contentCat := "file content for cat"
  contentCar := "file content for car"
  contentCart := "file content for cart"

  cat := sha1.Sum([]byte(contentCat)) // Generate SHA-1 hash for cat
  hashCat := fmt.Sprintf("%x", cat[:])
  car := sha1.Sum([]byte(contentCar)) // Generate SHA-1 hash for car
  hashCar := fmt.Sprintf("%x", car[:])
  cart := sha1.Sum([]byte(contentCart)) // Generate SHA-1 hash for car
  hashCart := fmt.Sprintf("%x", cart[:])
	trie.addFile("cat", hashCat)
	trie.addFile("car", hashCar)
	trie.addFile("cart", hashCart)

	words := trie.loadWordsFromPrefix("dog")
	if len(words) != 0 {
		t.Errorf("Expected no words for prefix 'dog', but got %v", words)
	}
}

// Test walking through a word in the Trie
func TestWalkWord(t *testing.T) {
	trie := NewTrie()
	contentBat := "file content for bat"
	bat := sha1.Sum([]byte(contentBat)) // Generate SHA-1 hash for car
  hashBat := fmt.Sprintf("%x", bat[:])
	trie.addFile("bat", hashBat)

	node := trie.walkWord("bat")
	if node == nil || !node.IsEnd {
		t.Errorf("Expected to walk to the word 'bat', but it was not found")
	}
}

// Test walking a non-existent word in the Trie
func TestWalkNonExistentWord(t *testing.T) {
	trie := NewTrie()
	contentBat := "file content for bat"
  	bat := sha1.Sum([]byte(contentBat)) // Generate SHA-1 hash for car
    hashBat := fmt.Sprintf("%x", bat[:])
  	trie.addFile("bat", hashBat)

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
	contentBat := "file content for bat"
  bat := sha1.Sum([]byte(contentBat)) // Generate SHA-1 hash for car
  hashBat := fmt.Sprintf("%x", bat[:])
  contentDog := "file content for dog"
  dog := sha1.Sum([]byte(contentDog)) // Generate SHA-1 hash for car
  hashDog := fmt.Sprintf("%x", dog[:])
  trie.addFile("bat", hashBat)
  trie.addFile("dog", hashDog)

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
	contentBat := "file content for bat"
  bat := sha1.Sum([]byte(contentBat)) // Generate SHA-1 hash for car
  hashBat := fmt.Sprintf("%x", bat[:])
  contentDog := "file content for dog"
  dog := sha1.Sum([]byte(contentDog)) // Generate SHA-1 hash for car
  hashDog := fmt.Sprintf("%x", dog[:])
  trie.addFile("bat", hashBat)
  trie.addFile("dog", hashDog)

	// Save the Trie first
	Save(trie, "test_trie.gob")

	// Load the Trie from the file
	loadedTrie := Load("test_trie.gob")
	if loadedTrie == nil {
		t.Errorf("Failed to load the Trie from the file")
	}

	// Check that the loaded Trie contains the words
	if !loadedTrie.walkWord("bat").IsEnd {
		t.Errorf("Expected the word 'cat' to be in the loaded Trie, but it was not found")
	}

	if !loadedTrie.walkWord("dog").IsEnd {
		t.Errorf("Expected the word 'dog' to be in the loaded Trie, but it was not found")
	}
}

// Test saving a Trie to disk and loading it
func TestSaveAndLoadTrie(t *testing.T) {
	trie := NewTrie()
	contentCat := "file content for cat"
	contentDog := "file content for dog"

	cat := sha1.Sum([]byte(contentCat)) // Generate SHA-1 hash for cat
	hashCat := fmt.Sprintf("%x", cat[:])
	dog := sha1.Sum([]byte(contentDog)) // Generate SHA-1 hash for dog
  hashDog := fmt.Sprintf("%x", dog[:])

	trie.addFile("cat", hashCat)
	trie.addFile("dog", hashDog)

	// Save the Trie to disk
	Save(trie, "test_trie.gob")

	// Check if the file was created
	if _, err := os.Stat("./.pm/trie/test_trie.gob"); os.IsNotExist(err) {
		t.Errorf("Expected the file 'test_trie.gob' to exist, but it does not")
	}

	// Load the Trie from the file
	loadedTrie := Load("test_trie.gob")
	if loadedTrie == nil {
		t.Errorf("Failed to load the Trie from the file")
	}

	// Check that the loaded Trie contains the files with the correct hashes
	if !reflect.DeepEqual(loadedTrie.walkWord("cat").Value, hashCat) {
		t.Errorf("Expected the hash for 'cat' in the loaded Trie to be %x, but got %x", hashCat, loadedTrie.walkWord("cat").Value)
	}

	if !reflect.DeepEqual(loadedTrie.walkWord("dog").Value, hashDog) {
		t.Errorf("Expected the hash for 'dog' in the loaded Trie to be %x, but got %x", hashDog, loadedTrie.walkWord("dog").Value)
	}
}

// Clean up test files after running tests
func TestCleanupFiles(t *testing.T) {
	err := os.Remove("./.pm/trie/test_trie.gob")
	if err != nil {
		t.Errorf("Error while cleaning up test file: %v", err)
	}
}
