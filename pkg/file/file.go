package file

import (
	"crypto/sha1"
	"encoding/gob"
	"errors"
	"fmt"
	"log"
	"os"

	"github/pm/pkg/common"
)

/**
Caller can create a new file which
Creates a blob and associates the blob to a vertex

How can the caller indicate the type of file. Options
Bucket them or tag them.

Conclusion: Bucket them using a map that writes

*/

type FileTypeIndex struct {
	TypeToFile map[string]map[string]string
	FileToType map[string]string
}

func NewReconcilableFileTypeIndex(storageKey string) common.Reconcilable {
	fileTypeIndexAlphaList := common.NewAlphaList()
	indexStorage := NewFileTypeIndex()
	filePath := "./.pm/fileTypes/" + storageKey

	return common.Reconcilable{
		AlphaList:     fileTypeIndexAlphaList,
		DataStructure: indexStorage,
		FilePath:      filePath,
	}

}

// For mvp, system declares file types
// TODO: refactor into constants
func NewFileTypeIndex() *FileTypeIndex {
	fileTypes := map[string]map[string]string{
		"prd":   {},
		"epic":  {},
		"story": {},
		"task":  {},
	}

	files := map[string]string{}

	return &FileTypeIndex{
		TypeToFile: fileTypes,
		FileToType: files,
	}
}

func (ft *FileTypeIndex) AddFileToIndex(fileName string, fileType string) error {
	value, ok := ft.TypeToFile[fileType]
	if !ok {
		return errors.New("File type not found, file not indexed")
	}

	_, ok = value[fileName]
	if ok {
		return errors.New("File already exists, cannot create new file")
	}

	value[fileName] = ""

	_, ok = ft.FileToType[fileName]
	if ok {
		return errors.New("File already exists, cannot create new file")
	}

	ft.FileToType[fileName] = fileType

	return nil
}

func (ft *FileTypeIndex) RemoveFileFromIndex(fileName string, fileType string) error {
	value, ok := ft.TypeToFile[fileType]
	if !ok {
		return errors.New("File type not found, cannot remove from non existent file type. Type: " + fileType)
	}

	_, ok = value[fileName]
	if !ok {
		return errors.New("File not found in index, cannot remove non existent file. Name: " + fileName)
	}

	delete(value, fileName)

	_, ok = ft.FileToType[fileName]
	if !ok {
		return errors.New("File not found in index, cannot remove non existent file. Name: " + fileName)
	}

	delete(ft.FileToType, fileName)

	return nil
}

func (ft *FileTypeIndex) RetrieveFilesFromType(fileType string) ([]string, error) {
	value, ok := ft.TypeToFile[fileType]
	if !ok {
		return nil, errors.New("File type not found, cannot remove from non existent file type. Type: " + fileType)
	}

	index := 0
	files := make([]string, len(value))
	for i := range value {
		files[index] = i
		index++
	}

	return files, nil
}

func (ft *FileTypeIndex) RetrieveFileType(fileName string) (string, error) {
	value, ok := ft.FileToType[fileName]
	if !ok {
		log.Println("Error finding this file: " + fileName)
		return "", errors.New("File not found in index. FileName: " + fileName)
	}

	return value, nil
}

type AddFileTypeIndexAlpha struct {
	Hash     string
	FileName string
	FileType string
}

func (aft *AddFileTypeIndexAlpha) GetType() byte {
	return common.AddFileAlpha
}

func (aft *AddFileTypeIndexAlpha) GetId() string {
	return aft.FileName + aft.FileType + string(common.AddFileAlpha)
}

func (aft *AddFileTypeIndexAlpha) GetHash() string {
	return aft.Hash
}

func (aft *AddFileTypeIndexAlpha) SetHash(lastAlpha common.Alpha) {
	prevAlphaHash := lastAlpha.GetHash()
	currentHash := sha1.Sum([]byte(aft.GetId() + prevAlphaHash))
	currentHashStr := fmt.Sprintf("%x", currentHash[:])
	aft.Hash = currentHashStr
}

type RemoveFileTypeIndexAlpha struct {
	Hash     string
	FileName string
	FileType string
}

func (rft *RemoveFileTypeIndexAlpha) GetType() byte {
	return common.RemoveFileAlpha
}

func (rft *RemoveFileTypeIndexAlpha) GetId() string {
	return rft.FileName + rft.FileType + string(common.RemoveFileAlpha)
}

func (rft *RemoveFileTypeIndexAlpha) GetHash() string {
	return rft.Hash
}

func (rft *RemoveFileTypeIndexAlpha) SetHash(lastAlpha common.Alpha) {
	prevAlphaHash := lastAlpha.GetHash()
	currentHash := sha1.Sum([]byte(rft.GetId() + prevAlphaHash))
	currentHashStr := fmt.Sprintf("%x", currentHash[:])
	rft.Hash = currentHashStr
}

func (ft *FileTypeIndex) Update(alpha common.Alpha) error {
	alphaType := alpha.GetType()
	var error error

	switch alphaType {
	case common.AddFileAlpha:
		addFileAlpha := alpha.(*AddFileTypeIndexAlpha)
		error = ft.AddFileToIndex(addFileAlpha.FileName, addFileAlpha.FileType)
	case common.RemoveFileAlpha:
		removeFileAlpha := alpha.(*RemoveFileTypeIndexAlpha)
		error = ft.RemoveFileFromIndex(removeFileAlpha.FileName, removeFileAlpha.FileType)
	}

	if error != nil {
		return error
	}

	return nil
}

// Dont want to wrok on this now as it is not used
func (ft *FileTypeIndex) Rewind(alpha common.Alpha) error {
	return nil
}

// Dont want to work on this now
func (ft *FileTypeIndex) Validate(alpha common.Alpha) bool {
	return true
}

// TODO: refactor into common utility function in the future
// TODO: refacotr into a parser that can read plain text file
func LoadReconcilableFileTypeIndex(filePath string) common.Reconcilable {
	file, fileErr := os.Open(filePath)

	if fileErr != nil {
		log.Println("Error opening binary file")
		return common.Reconcilable{}
	}
	defer file.Close()

	gob.Register(&FileTypeIndex{})
	decoder := gob.NewDecoder(file)
	var loadedReconcilable common.Reconcilable
	decodingErr := decoder.Decode(&loadedReconcilable)
	if decodingErr != nil {
		log.Println("Error decoding", decodingErr.Error())
		return common.Reconcilable{}
	}

	return loadedReconcilable
}
