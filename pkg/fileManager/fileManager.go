package fileManager

import (
	"crypto/sha1"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github/pm/pkg/blob"
	"github/pm/pkg/common"
	"github/pm/pkg/trie"
)

func UpdateBlobContent(name string, content string, index *common.Reconcilable) error {
	defer index.SaveReconcilable()
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

func Commit(name string, content string, index *common.Reconcilable) error {
	defer index.SaveReconcilable()
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

func IndexName(fileName string, hashContent string, index *common.Reconcilable) error {
	trie := index.DataStructure.(*trie.Trie)
	indexErr := trie.AddFile(fileName, hashContent)
	if indexErr != nil {
		return indexErr
	}

	return nil
}

func ReIndex(fileName string, hashContent string, index *common.Reconcilable) error {
	trie := index.DataStructure.(*trie.Trie)
	indexErr := trie.UpdateValue(fileName, hashContent)
	if indexErr != nil {
		return indexErr
	}

	return nil
}

// Delete a blob and remove its index in the Trie
func DeleteBlob(fileName string, index *common.Reconcilable) error {
	trie := index.DataStructure.(*trie.Trie)
	// Retrieve the hash from the Trie for the given file name
	hash, err := trie.RetrieveValue(fileName)
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
	trie.RemoveWord(fileName)

	return nil
}

func RetrieveContent(fileName string, index *common.Reconcilable) ([]byte, error) {
	trie := index.DataStructure.(*trie.Trie)
	hash, err := trie.RetrieveValue(fileName)
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
