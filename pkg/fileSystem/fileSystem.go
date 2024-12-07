package filesystem

import (
	"github/pm/pkg/common"
	"github/pm/pkg/dag"
	pmfile "github/pm/pkg/file"
	"github/pm/pkg/blob"

	"os"
	"path/filepath"
	"errors"
	"fmt"
)

type FileSystem struct {
	fileRelationShips common.Reconcilable
	fileTypeIndex common.Reconcilable
}

func checkFileExists(filePath string) bool {
	_, error := os.Stat(filePath)
	//return !os.IsNotExist(err)
	return !errors.Is(error, os.ErrNotExist)
}

func checkDirExists(dirPath string) bool {
	_, err := os.Stat(dirPath)
	return !os.IsNotExist(err)
}

func (fs *FileSystem) BootDag() (error) {
	dagDirectory := filepath.Join(".", ".pm", "dag")
	dagFile := filepath.Join(".", ".pm", "dag", "dag")

	if !checkDirExists(dagDirectory) {
		err := os.MkdirAll(dagDirectory, os.ModePerm)
		if err != nil {
			return errors.New("Error creating directory for file relationships")
		}
	}

	if !checkFileExists(dagFile) {
		file, fileErr := os.Create(dagFile)
		if fileErr != nil {
			fmt.Printf("Error creating tmp file: %v\n", fileErr)
		}

		defer file.Close()

		// TODO: Refactor to pass in path file
		fs.fileRelationShips = dag.NewReconcilableDag("dag")
	} else {
		fs.fileRelationShips = *dag.LoadReconcilableDag(dagFile)
	}

	return nil
}

func (fs *FileSystem) BootFileTypes() (error) {
	fileTypeDirectory := filepath.Join(".", ".pm", "fileTypes")
	fileTypeFile := filepath.Join(".", ".pm", "fileTypes", "types")

	if !checkDirExists(fileTypeDirectory) {
		err := os.MkdirAll(fileTypeDirectory, os.ModePerm)
		if err != nil {
			return errors.New("Error creating directory for file relationships")
		}
	}

	if !checkFileExists(fileTypeFile) {
		fileType := filepath.Join(fileTypeFile)
		file, fileErr := os.Create(fileType)
		if fileErr != nil {
			fmt.Printf("Error creating tmp file: %v\n", fileErr)
		}

		defer file.Close()

		fs.fileTypeIndex = pmfile.NewReconcilableFileTypeIndex("types");
	} else {
		fs.fileTypeIndex = *pmfile.LoadReconcilableFileTypeIndex(fileTypeFile);
	}

	return nil
}

func (fs *FileSystem) Boot() (error) {
	// Load/Create File fileRelationShips
	bootDagErr := fs.BootDag();
	if bootDagErr != nil {
		return bootDagErr
	}

	// Load/CreateFile Type index
	bootIndexErr := fs.BootFileTypes()
	if bootIndexErr != nil {
		return bootIndexErr
	}

	return nil
}

func (fs *FileSystem) getFileIndex() (*pmfile.FileTypeIndex) {
	return fs.fileTypeIndex.DataStructure.(*pmfile.FileTypeIndex)
}

func (fs *FileSystem) getFileTree() (*dag.Dag) {
	return fs.fileTypeIndex.DataStructure.(*dag.Dag)
}

func (fs *FileSystem) CreateFile(fileName string) (error) {
	// Create Blob using fileName
	blobErr := blob.CreateBlob(fileName, "");
	if blobErr != nil {
		return blobErr
	}

	// Add name to FileIndex
	addFileIndexAlpha := pmfile.AddFileTypeIndexAlpha{
		FileName: fileName,
		FileType: "",
	};

	updateErr := fs.fileTypeIndex.DataStructure.Update(&addFileIndexAlpha)
	if updateErr != nil {
		return updateErr
	}

	vertex := dag.NewVertex(fileName);

	// Add Vertex in Dag
	addVertexAlpha := dag.AddVertexAlpha{
		Target: vertex,
	}

	updateErr = fs.fileRelationShips.DataStructure.Update(&addVertexAlpha)
	if updateErr != nil {
		return updateErr
	}

	return nil
}

func (fs *FileSystem) DeleteFile(fileId string) (error) {
	// Remove name from fileTypeIndex

	// Remove Blob using name

	// Remove Vertex from Dag
	return nil
}

func (fs *FileSystem) LinkFile(parentId string, childId string) (error) {
	// Check if Parent and Child exist in File index 

	// Check if Parent and Child blobs exist

	// Check if Parent and Child vertex exist in the FileTree

	// Check if Edge exist for Parent and child

	// Add Edge between parent and child vertex
	return nil
}

func (fs *FileSystem) UnLinkFile(parentId string, childId string) (error) {
	// Check if Parent and Child exist in File index

	// Check if Parent and Child blobs exist

	// Check if Parent and Child vertex exist

	// Check if Parent and Child edge exist

	// Remove edge between Parent and Child
	return nil
}

func (fs *FileSystem) ListFileNamesByType(fileType string) ([]string, error) {
	fileIndex := fs.getFileIndex();
	files, fileErr := fileIndex.RetrieveFilesFromType(fileType)
	if fileErr != nil {
		return nil, fileErr
	}

	return files, nil
}
