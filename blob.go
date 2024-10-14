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

func commit(name string, content string, index *Trie) error {
  blobErr := createBlob(content);
  if blobErr != nil {
    return blobErr;
  }

  indexErr := indexName(name, index);
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

	compressContent, err := compressContent(content)
	if err != nil {
		fmt.Println("Error compressing content")
		return err
	}

	err = os.WriteFile(blobPath, compressContent, 0644)
	if err != nil {
		fmt.Println("Error writing to file")
		return err
	}

	return nil
}

func compressContent(content string) ([]byte, error) {
	var b bytes.Buffer

	w := zlib.NewWriter(&b)
	_, err := w.Write([]byte(content))
	if err != nil {
		return nil, err
	}

	err = w.Close()
	if err != nil {
		return nil, err
	}

	return b.Bytes(), nil
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

func indexName(fileName string, index *Trie) error {
	indexErr := index.addWord(fileName)
	if indexErr != nil {
		return indexErr
	}

	return nil
}
