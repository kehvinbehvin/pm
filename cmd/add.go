package cmd

import (
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Initialize a new .pm project",
	Run: func(cmd *cobra.Command, args []string) {
		createEpic("Loyalty", "hello")
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}

func createEpic(name string, content string) {
	storeFileName(name)

	hash := sha1.Sum([]byte(content))
	hashStr := fmt.Sprintf("%x", hash[:])
	blobDir := ".pm/blobs/"
	subDir := blobDir + hashStr[:2]
	err := os.MkdirAll(subDir, os.ModePerm)
	if err != nil {
		fmt.Println("error 1")
	}

	blobPath := filepath.Join(subDir, hashStr)

	compressContent, err := compressContent(content)
	if err != nil {
		fmt.Println("error 2")
	}
	err = os.WriteFile(blobPath, compressContent, 0644)
	if err != nil {
		fmt.Println("error 3")
	}

	decompressed, err := decompressContent(compressContent)
	if err != nil {
		fmt.Println("Error decompressing string:", err)
		return
	}
	fmt.Println("Decompressed string:", decompressed)
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

func storeFileName(filename string) error {
	// Step 1: Hash the filename using SHA-256
	hash := sha1.Sum([]byte(filename))
	hashStr := fmt.Sprintf("%x", hash[:]) // Convert hash to a hex string

	// Step 2: Create the 4-level directory structure
	dir1 := hashStr[:2]
	dir2 := hashStr[2:4]
	dir3 := hashStr[4:6]
	fileName := hashStr[6:]

	// Step 3: Create the directories
	baseDir := ".pm/names/"
	fullDir := filepath.Join(baseDir, dir1, dir2, dir3)
	err := os.MkdirAll(fullDir, os.ModePerm) // Create all necessary directories
	if err != nil {
		return err
	}

	// Step 4: Store the file in the last directory
	filePath := filepath.Join(fullDir, fileName+".txt")
	err = os.WriteFile(filePath, []byte(filename), 0644) // Write file content
	if err != nil {
		return err
	}

	fmt.Println("File stored at:", filePath)
	return nil
}
