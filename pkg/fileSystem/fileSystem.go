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
	// TODO: refactor to use reconcilable data structure
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

func (fs *FileSystem) DeleteFile(fileName string) (error) {
	// Remove name from fileTypeInde
	removeFileIndexAlpha := pmfile.RemoveFileTypeIndexAlpha{
		FileName: fileName,
		FileType: "",
	};

	updateErr := fs.fileTypeIndex.DataStructure.Update(&removeFileIndexAlpha)
	if updateErr != nil {
		return updateErr
	}


	// Already handles non-existent blobs
	deleteErr := blob.DeleteBlob(fileName)
	if deleteErr != nil {
		return deleteErr
	}

	fileTree := fs.getFileTree();
	vertex := fileTree.RetrieveVertex(fileName)
	if vertex != nil {
		return errors.New("File not found in file system")
	}

	// Remove Vertex from Dag
	removeVertexAlpha := dag.RemoveVertexAlpha{
		Target: vertex,
	}

	updateErr = fs.fileRelationShips.DataStructure.Update(&removeVertexAlpha)
	if updateErr != nil {
		return updateErr
	}

	return nil
}

func (fs *FileSystem) validateFileExists(fileName string) (error) {
	// Check if Parent and Child exist in File index 
	fileIndex := fs.getFileIndex()
	_ , fileErr := fileIndex.RetrieveFileType(fileName)
	if fileErr != nil {
		return errors.New("File not in index: File: " + fileName)
	}

	// Check if Parent and Child blobs exist
	fileBlob := blob.Exists(fileName)
	if !fileBlob {
		return errors.New("File blob not found: File: " + fileName)
	}

	// Check if Parent and Child vertex exist in the FileTree
	fileTree := fs.getFileTree();
	fileVertex := fileTree.RetrieveVertex(fileName)
	if fileVertex != nil {
		return errors.New("File Vertext not found: File " + fileName)
	}

	return nil
}

func (fs *FileSystem) LinkFile(parentName string, childName string) (error) {
	parentErr := fs.validateFileExists(parentName);
	if parentErr != nil {
		return parentErr
	}

	childErr := fs.validateFileExists(childName);
	if childErr != nil {
		return childErr
	}

	fileTree := fs.getFileTree();
	parentVertex := fileTree.RetrieveVertex(parentName)
	childVertex := fileTree.RetrieveVertex(childName)

	// Add Edge between parent and child vertex
	addEdgeAlpha := dag.AddEdgeAlpha{
		From: parentVertex,
		To: childVertex,
	}

	updateErr := fs.fileRelationShips.DataStructure.Update(&addEdgeAlpha)
	if updateErr != nil {
		return updateErr
	}

	return nil
}

func (fs *FileSystem) UnLinkFile(parentName string, childName string) (error) {
	parentErr := fs.validateFileExists(parentName);
	if parentErr != nil {
		return parentErr
	}

	childErr := fs.validateFileExists(childName);
	if childErr != nil {
		return childErr
	}

	fileTree := fs.getFileTree();
	parentVertex := fileTree.RetrieveVertex(parentName)
	childVertex := fileTree.RetrieveVertex(childName)

	// Add Edge between parent and child vertex
	removeEdgeAlpha := dag.RemoveEdgeAlpha{
		From: parentVertex,
		To: childVertex,
	}

	updateErr := fs.fileRelationShips.DataStructure.Update(&removeEdgeAlpha)
	if updateErr != nil {
		return updateErr
	}

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
