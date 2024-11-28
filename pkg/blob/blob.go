package blob

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

func CreateBlob(content string) error {
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

func CompressContent(content string) (string, error) {
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

func DecompressContent(content []byte) (string, error) {
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

func RemoveIfEmpty(dirPath string) error {
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
