package filesystem

import (
	"github/pm/pkg/blob"
	"github/pm/pkg/common"
	dag "github/pm/pkg/dag"
	pmfile "github/pm/pkg/file"

	"sort"
	"errors"
	"os"
	"path/filepath"
	"log"
)

type FileSystem struct {
	fileRelationShips common.Reconcilable
	fileTypeIndex     common.Reconcilable
}

func NewFileSystem() (*FileSystem) {
	return &FileSystem{}
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

func (fs *FileSystem) ShutDown() error {
	fs.fileRelationShips.SaveReconcilable()
	fs.fileTypeIndex.SaveReconcilable()
	return nil
}

func (fs *FileSystem) BootDag() error {
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
			log.Printf("Error creating tmp file: %v\n", fileErr)
		}

		defer file.Close()

		// TODO: Refactor to pass in path file
		fs.fileRelationShips = dag.NewReconcilableDag("dag")
		defer fs.fileRelationShips.SaveReconcilable()
	} else {
		reconcilable := dag.LoadReconcilableDag(dagFile)
		fs.fileRelationShips = reconcilable
		defer fs.fileRelationShips.SaveReconcilable()
	}

	return nil
}

func (fs *FileSystem) BootFileTypes() error {
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
			log.Printf("Error creating tmp file: %v\n", fileErr)
		}

		defer file.Close()

		fs.fileTypeIndex = pmfile.NewReconcilableFileTypeIndex("types")
		defer fs.fileTypeIndex.SaveReconcilable()
	} else {
		fileTypeIndex := pmfile.LoadReconcilableFileTypeIndex(fileTypeFile)
		fs.fileTypeIndex = fileTypeIndex
		defer fs.fileTypeIndex.SaveReconcilable()
	}

	return nil
}

func (fs *FileSystem) Boot() error {
	// Load/Create File fileRelationShips
	bootDagErr := fs.BootDag()
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

func (fs *FileSystem) getFileIndex() *pmfile.FileTypeIndex {
	return fs.fileTypeIndex.DataStructure.(*pmfile.FileTypeIndex)
}

func (fs *FileSystem) getFileTree() *dag.Dag {
	return fs.fileRelationShips.DataStructure.(*dag.Dag)
}

func (fs *FileSystem) CreateFile(fileName string, fileType string) error {
  log.Println("Filename: " + fileName + " created");
	// Create Blob using fileName
	// TODO: refactor to use reconcilable data structure
	blobErr := blob.CreateBlob(fileName, "")
	if blobErr != nil {
		return blobErr
	}

	// Add name to FileIndex
	addFileIndexAlpha := pmfile.AddFileTypeIndexAlpha{
		FileName: fileName,
		FileType: fileType,
	}

	updateErr := fs.fileTypeIndex.DataStructure.Update(&addFileIndexAlpha)
	if updateErr != nil {
		return updateErr
	}

	vertex := dag.NewVertex(fileName)

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

func (fs *FileSystem) DeleteFile(fileName string) error {
	// Remove name from fileTypeInde
	removeFileIndexAlpha := pmfile.RemoveFileTypeIndexAlpha{
		FileName: fileName,
		FileType: "",
	}

	updateErr := fs.fileTypeIndex.DataStructure.Update(&removeFileIndexAlpha)
	if updateErr != nil {
		return updateErr
	}

	// Already handles non-existent blobs
	deleteErr := blob.DeleteBlob(fileName)
	if deleteErr != nil {
		return deleteErr
	}

	fileTree := fs.getFileTree()
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

func (fs *FileSystem) validateFileExists(fileName string) error {
	// Check if Parent and Child exist in File index
	fileIndex := fs.getFileIndex()
	_, fileErr := fileIndex.RetrieveFileType(fileName)
	if fileErr != nil {
		log.Println("File not in index")
		return errors.New("File not in index: File: " + fileName)
	}

	// Check if Parent and Child blobs exist
	fileBlob := blob.Exists(fileName)
	if !fileBlob {
		log.Println("File not in blobs")
		return errors.New("File blob not found: File: " + fileName)
	}

	// Check if Parent and Child vertex exist in the FileTree
	fileTree := fs.getFileTree()
	fileVertex := fileTree.RetrieveVertex(fileName)
	if fileVertex == nil {
		log.Println("File not in tree")
		return errors.New("File Vertex not found: File " + fileName)
	}

	return nil
}

func (fs *FileSystem) LinkFile(parentName string, childName string) error {

	if parentName == "" || childName == "" {
		return nil
	}

	log.Println("Parent: " + parentName + "; Child: " + childName);
	parentErr := fs.validateFileExists(parentName)
	if parentErr != nil {
		log.Println("Parent cannot be found");
		return parentErr
	}

	childErr := fs.validateFileExists(childName)
	if childErr != nil {
		log.Println("Child cannot be found");
		return childErr
	}

	fileTree := fs.getFileTree()
	parentVertex := fileTree.RetrieveVertex(parentName)
	childVertex := fileTree.RetrieveVertex(childName)

	// Add Edge between parent and child vertex
	addEdgeAlpha := dag.AddEdgeAlpha{
		From: childVertex,
		To:   parentVertex,
	}

	updateErr := fs.fileRelationShips.DataStructure.Update(&addEdgeAlpha)
	if updateErr != nil {
		return updateErr
	}

	return nil
}

func (fs *FileSystem) UnLinkFile(parentName string, childName string) error {
	parentErr := fs.validateFileExists(parentName)
	if parentErr != nil {
		return parentErr
	}

	childErr := fs.validateFileExists(childName)
	if childErr != nil {
		return childErr
	}

	fileTree := fs.getFileTree()
	parentVertex := fileTree.RetrieveVertex(parentName)
	childVertex := fileTree.RetrieveVertex(childName)

	// Add Edge between parent and child vertex
	removeEdgeAlpha := dag.RemoveEdgeAlpha{
		From: parentVertex,
		To:   childVertex,
	}

	updateErr := fs.fileRelationShips.DataStructure.Update(&removeEdgeAlpha)
	if updateErr != nil {
		return updateErr
	}

	return nil
}

func (fs *FileSystem) ListFileNamesByType(fileType string) ([]string, error) {
	fileIndex := fs.getFileIndex()
	files, fileErr := fileIndex.RetrieveFilesFromType(fileType)
	if fileErr != nil {
		return nil, fileErr
	}

	sort.Strings(files)

	return files, nil
}

func (fs *FileSystem) ListChildIssues(fileName string) ([]string, error) {
	dag := fs.fileRelationShips.DataStructure.(*dag.Dag)
	vertex := dag.RetrieveVertex(fileName)
	children := vertex.Children
	
	var childIssues []string
	for _, child := range children {
		childIssues = append(childIssues, child.ID)
	}

	return childIssues, nil
}
