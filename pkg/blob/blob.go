package blob

import (
	"bytes"
	"compress/zlib"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"log"
)

const compressionThreshold = 10 * 1024 // 10 KB threshold for compression

// Use Human readable fileNames for easier reading and potability
func CreateBlob(fileName string, content string) error {
	fileName = fileName + ".md"
	blobDirectory := filepath.Join(".", ".pm", "blobs")
	blobFile := filepath.Join(".", ".pm", "blobs", fileName)

	err := os.MkdirAll(blobDirectory, os.ModePerm)
	if err != nil {
		log.Println("Error creating blob dir")
		return err
	}

	contentSize := len(content)
	if contentSize > compressionThreshold {
		content, err = CompressContent(content)
	}

	if err != nil {
		log.Println("Error compressing content")
		return err
	}

	err = os.WriteFile(blobFile, []byte(content), 0644)
	if err != nil {
		log.Println("Error writing to file")
		return err
	}

	return nil
}

func checkFileExists(filePath string) bool {
	_, error := os.Stat(filePath)
	//return !os.IsNotExist(err)
	return !errors.Is(error, os.ErrNotExist)
}

func Exists(fileName string) bool {
	fileName = fileName + ".md"
	path := filepath.Join(".", ".pm", "./blobs", fileName)
	return checkFileExists(path)
}

func DeleteBlob(fileName string) error {
	fileName = fileName + ".md"
	path := filepath.Join(".", ".pm", "./blobs", fileName)
	exists := checkFileExists(path)
	if !exists {
		return errors.New("Cannot remove non existent file")
	}

	deleteErr := os.Remove(path)
	if deleteErr != nil {
		return deleteErr
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
		log.Printf("Directory %s was empty and has been deleted.\n", dirPath)
	}

	return nil
}
