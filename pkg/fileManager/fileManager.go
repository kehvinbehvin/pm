package fileManager

import (
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github/pm/blob"
	"github/pm/trie"
)

func UpdateBlobContent(name string, content string, index *trie.Trie) error {
	defer index.Save()
	hash := sha1.Sum([]byte(content))
	hashStr := fmt.Sprintf("%x", hash[:])

	blobErr := blob.CreateBlob(content)
	if blobErr != nil {
		return blobErr
	}

	indexErr := ReIndex(name, hashStr, index)
	if indexErr != nil {
		return indexErr
	}

	return nil
}

func Commit(name string, content string, index *trie.Trie) error {
	defer index.Save()
	hash := sha1.Sum([]byte(content))
	hashStr := fmt.Sprintf("%x", hash[:])

	blobErr := blob.CreateBlob(content)
	if blobErr != nil {
		return blobErr
	}

	indexErr := IndexName(name, hashStr, index)
	if indexErr != nil {
		return indexErr
	}

	return nil
}

func IndexName(fileName string, hashContent string, index *trie.Trie) error {
	indexErr := index.addFile(fileName, hashContent)
	if indexErr != nil {
		return indexErr
	}

	return nil
}

func ReIndex(fileName string, hashContent string, index *trie.Trie) error {
	indexErr := index.updateValue(fileName, hashContent)
	if indexErr != nil {
		return indexErr
	}

	return nil
}

// Delete a blob and remove its index in the Trie
func DeleteBlob(fileName string, index *trie.Trie) error {
	// Retrieve the hash from the Trie for the given file name
	hash, err := index.retrieveValue(fileName)
	if err != nil {
		return fmt.Errorf("failed to retrieve hash for file '%s': %v", fileName, err)
	}

	// Construct the blob path using the hash
	blobPath := filepath.Join(".pm/blobs", hash[:2], hash)
	blobDir := filepath.Join(".pm/blobs", hash[:2])

	// Remove the blob file
	err = os.Remove(blobPath)
	if err != nil {
		return fmt.Errorf("failed to delete blob: %v", err)
	}
	err = blob.RemoveIfEmpty(blobDir)

	if err != nil {
		return fmt.Errorf("failed to delete blob dir: %v", err)
	}

	// Remove the file name from the Trie index
	index.removeWord(fileName)

	return nil
}

func RetrieveContent(fileName string, index *trie.Trie) ([]byte, error) {
	hash, err := index.retrieveValue(fileName)
	if err != nil {
		fmt.Errorf("failed to retrieve hash for file '%s': %v", fileName, err)
		return []byte(""), err
	}

	// Construct the blob path using the hash
	blobPath := filepath.Join(".pm/blobs", hash[:2], hash)

	fileContent, err := os.Open(blobPath)
	if err != nil {
		return []byte(""), fmt.Errorf("failed to open blob file: %v", err)
	}
	defer fileContent.Close()

	content, err := io.ReadAll(fileContent)
	if err != nil {
		return []byte(""), fmt.Errorf("failed to read blob file: %v", err)
	}

	return content, nil
}