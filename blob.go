package main;

import (
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

const compressionThreshold = 1024 // 1 KB threshold for compression

func commit(name string, content string, index *Trie) error {
  defer Save(index, "test")
  hash := sha1.Sum([]byte(content))
	hashStr := fmt.Sprintf("%x", hash[:])

  blobErr := createBlob(content);
  if blobErr != nil {
    return blobErr;
  }

  indexErr := indexName(name, hashStr, index);
  if indexErr != nil {
    return indexErr;
  }

  return nil
}

func createBlob(content string) (error) {
  hash := sha1.Sum([]byte(content))
	hashStr := fmt.Sprintf("%x", hash[:])
	blobDir := ".pm/blobs/"
	subDir := blobDir + hashStr[:2]
	err := os.MkdirAll(subDir, os.ModePerm)
	if err != nil {
		fmt.Println("Error creating blob file")
		return err
	}

	blobPath := filepath.Join(subDir, hashStr)

	contentSize := len(content)
	if contentSize > compressionThreshold {
		content, err = compressContent(content)
	}

	if err != nil {
		fmt.Println("Error compressing content")
		return err
	}

	err = os.WriteFile(blobPath, []byte(content), 0644)
	if err != nil {
		fmt.Println("Error writing to file")
		return err
	}

	return nil
}

func compressContent(content string) (string, error) {
	var b bytes.Buffer

	w := zlib.NewWriter(&b)
	_, err := w.Write([]byte(content))
	if err != nil {
		return "", err
	}

	err = w.Close()
	if err != nil {
		return "", err
	}

	return string(b.Bytes()), nil
}

func decompressContent(content []byte) (string, error) {
	b := bytes.NewReader(content)
	r, err := zlib.NewReader(b)
	if err != nil {
		return "", err
	}
	defer r.Close()

	var out bytes.Buffer
	_, err = io.Copy(&out, r)
	if err != nil {
		return "", err
	}

	return out.String(), nil
}

func indexName(fileName string, hashContent string, index *Trie) error {
	indexErr := index.addFile(fileName, hashContent)
	if indexErr != nil {
		return indexErr
	}

	return nil
}

// Delete a blob and remove its index in the Trie
func deleteBlob(fileName string, index *Trie) error {
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
	err = removeIfEmpty(blobDir);

	if err != nil {
    return fmt.Errorf("failed to delete blob dir: %v", err)
  }

	// Remove the file name from the Trie index
	index.removeWord(fileName)

	return nil
}

func removeIfEmpty(dirPath string) error {
	// Open the directory
	dir, err := os.Open(dirPath)
	if err != nil {
		return fmt.Errorf("failed to open directory: %v", err)
	}
	defer dir.Close()

	// Read directory contents
	files, err := dir.Readdir(1) // Read one file to check if it's empty
	if err != nil && err != io.EOF {
		return fmt.Errorf("failed to read directory: %v", err)
	}

	// Check if the directory is empty
	if len(files) == 0 {
		// Directory is empty, proceed to delete
		err = os.Remove(dirPath)
		if err != nil {
			return fmt.Errorf("failed to remove directory: %v", err)
		}
		fmt.Printf("Directory %s was empty and has been deleted.\n", dirPath)
	}

	return nil
}

func retrieveContent(fileName string, index *Trie) ([]byte, error) {
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
