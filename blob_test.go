package main

import (
	"crypto/sha1"
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

// Test committing a file and creating a blob
func TestCommit(t *testing.T) {
	// Setup
	trie := NewTrie("test")
	content := "This is a test content"
	filename := "testfile.txt"

	// Commit the file
	err := commit(filename, content, trie)
	if err != nil {
		t.Fatalf("Error during commit: %v", err)
	}

	// Check if the blob was created
	hash := sha1.Sum([]byte(content))
	hashStr := fmt.Sprintf("%x", hash[:])
	blobPath := filepath.Join(".pm/blobs", hashStr[:2], hashStr)

	if _, err := os.Stat(blobPath); os.IsNotExist(err) {
		t.Fatalf("Blob file %s not found after commit", blobPath)
	}

	// Check if the file is indexed in the Trie
	node := trie.walkWord(filename)
	if node == nil || !node.IsEnd {
		t.Errorf("File %s not indexed in the Trie", filename)
	}
}

// // Test creating a blob with compression for large files
func TestCreateBlobWithCompression(t *testing.T) {
	content := make([]byte, compressionThreshold+1) // Exceed threshold
	for i := range content {
		content[i] = 'A'
	}
	hashStr := fmt.Sprintf("%x", content[:])
	err := createBlob(string(content))
	if err != nil {
		t.Fatalf("Error creating blob with compression: %v", err)
	}

	// Check if the blob was created
	blobPath := filepath.Join(".pm/blobs", hashStr[:2], hashStr)
	if _, err := os.Stat(blobPath); os.IsNotExist(err) {
		t.Fatalf("Blob file %s not found after creation", blobPath)
	}
}

// Test creating a blob without compression for small files
func TestCreateBlobWithoutCompression(t *testing.T) {
	content := "This is a small file content"
	hash := sha1.Sum([]byte(content))
	hashStr := fmt.Sprintf("%x", hash[:])
	err := createBlob(content)
	if err != nil {
		t.Fatalf("Error creating blob without compression: %v", err)
	}

	// Check if the blob was created
	blobPath := filepath.Join(".pm/blobs", hashStr[:2], hashStr)
	if _, err := os.Stat(blobPath); os.IsNotExist(err) {
		t.Fatalf("Blob file %s not found after creation", blobPath)
	}
}

// Test compressing and decompressing content
func TestCompressAndDecompress(t *testing.T) {
	content := "This is a content to compress"

	// Compress content
	compressedContent, err := compressContent(content)
	if err != nil {
		t.Fatalf("Error compressing content: %v", err)
	}

	// Decompress content
	decompressedContent, err := decompressContent([]byte(compressedContent))
	if err != nil {
		t.Fatalf("Error decompressing content: %v", err)
	}

	// Compare the original and decompressed content
	if content != decompressedContent {
		t.Errorf("Expected decompressed content to be '%s', but got '%s'", content, decompressedContent)
	}
}

// Test deleting a blob and removing its index from the Trie
func TestDeleteBlob(t *testing.T) {
	// Setup
	trie := NewTrie("test")
	content := "This is another test content"
	filename := "deletefile.txt"

	// Commit the file
	err := commit(filename, content, trie)
	if err != nil {
		t.Fatalf("Error during commit: %v", err)
	}

	// Delete the blob
	err = deleteBlob(filename, trie)
	if err != nil {
		t.Fatalf("Error during deleteBlob: %v", err)
	}

	// Check if the blob was deleted
	hash := sha1.Sum([]byte(content))
	hashStr := fmt.Sprintf("%x", hash[:])
	blobPath := filepath.Join(".pm/blobs", hashStr[:2], hashStr)

	if _, err := os.Stat(blobPath); !os.IsNotExist(err) {
		t.Errorf("Expected blob file %s to be deleted, but it still exists", blobPath)
	}

	// Check if the file was removed from the Trie
	node := trie.walkWord(filename)
	if node != nil {
		if node.IsEnd {
			t.Errorf("Expected file %s to be removed from the Trie, but it still exists", filename)
		}
	}
}

// Test retrieving content from blob
func TestRetrieveContent(t *testing.T) {
	// Setup
	trie := NewTrie("test")
	content := "This is the content to be retrieved"
	filename := "retrievefile.txt"

	// Commit the file
	err := commit(filename, content, trie)
	if err != nil {
		t.Fatalf("Error during commit: %v", err)
	}

	// Retrieve the content
	retrievedContent, err := retrieveContent(filename, trie)
	if err != nil {
		t.Fatalf("Error retrieving content: %v", err)
	}

	// Check if the retrieved content matches the original
	if string(retrievedContent) != content {
		t.Errorf("Expected retrieved content to be '%s', but got '%s'", content, retrievedContent)
	}
}

// Test cleanup of blobs directory after tests
// func TestCleanupBlobFiles(t *testing.T) {
// 	blobDir := ".pm/blobs"
// 	err := os.RemoveAll(blobDir)
// 	if err != nil {
// 		t.Fatalf("Error cleaning up blob directory: %v", err)
// 	}
// }
