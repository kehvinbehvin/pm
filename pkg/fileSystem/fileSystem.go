package fileSystem

import (
	"github/pm/pkg/blob"
	"github/pm/pkg/common"
	"github/pm/pkg/dag"
	pmfile "github/pm/pkg/file"

	"errors"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
)

const FILE_RELATIONSHIP_DEPENDENCY = "DEPENDENCY"
const FILE_RELATIONSHIPS_HIERARCHY = "HIERARCHY"

type FileSystem struct {
	fileRelationShips       common.Reconcilable
	fileTypeIndex           common.Reconcilable
	fileParentRelationships common.Reconcilable
}

func NewFileSystem() *FileSystem {
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
	fs.fileParentRelationships.SaveReconcilable()
	fs.fileTypeIndex.SaveReconcilable()
	return nil
}

func (fs *FileSystem) BootDag(key string) (common.Reconcilable, error) {
	dagDirectory := filepath.Join(".", ".pm", "dag")
	dagFile := filepath.Join(".", ".pm", "dag", key)

	if !checkDirExists(dagDirectory) {
		err := os.MkdirAll(dagDirectory, os.ModePerm)
		if err != nil {
			return common.Reconcilable{}, errors.New("Error creating directory for file relationships")
		}
	}

	if !checkFileExists(dagFile) {
		file, fileErr := os.Create(dagFile)
		if fileErr != nil {
			log.Printf("Error creating tmp file: %v\n", fileErr)
		}

		defer file.Close()

		// TODO: Refactor to pass in path file
		return dag.NewReconcilableDag(key), nil
	} else {
		return dag.LoadReconcilableDag(dagFile), nil
	}
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
	childDag, bootChildrenDagErr := fs.BootDag("children")
	if bootChildrenDagErr != nil {
		return bootChildrenDagErr
	}

	fs.fileRelationShips = childDag
	defer fs.fileRelationShips.SaveReconcilable()

	parentDag, bootParentDagErr := fs.BootDag("parent")
	if bootParentDagErr != nil {
		return bootParentDagErr
	}

	fs.fileParentRelationships = parentDag
	defer fs.fileRelationShips.SaveReconcilable()

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

func (fs *FileSystem) getParentFileTree() *dag.Dag {
	return fs.fileParentRelationships.DataStructure.(*dag.Dag)
}

func (fs *FileSystem) CreateFile(fileName string, fileType string) error {
	log.Println("Filename: " + fileName + " created")
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

	parentVertex := dag.NewVertex(fileName)
	addParentVertexAlpha := dag.AddVertexAlpha{
		Target: parentVertex,
	}

	updateErr = fs.fileParentRelationships.DataStructure.Update(&addParentVertexAlpha)
	if updateErr != nil {
		return updateErr
	}

	return nil
}

func (FileSystem) EditFile(fileName string) {
	filePath := filepath.Join(".", ".pm", "./blobs", fileName+".md")
	editor := os.Getenv("EDITOR")
	if editor == "" {
		// Fallback to a default editor if $EDITOR is not set
		editor = "vim"
	}

	// Open the file in the editor
	err := openEditor(editor, filePath)
	if err != nil {
		return
	}
}

func openEditor(editor string, filePath string) error {
	// Create an exec command to open the file in the editor
	cmd := exec.Command(editor, filePath)

	// Set the command to use the same standard input, output, and error streams as the Go process
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Run the command and wait for it to finish
	err := cmd.Run()
	if err != nil {
		log.Printf("Error creating tmp file: %v\n", err)
		return nil
	}

	return nil
}

func (fs *FileSystem) DeleteFile(fileName string, fileType string) error {
	// Remove name from fileTypeInde
	removeFileIndexAlpha := pmfile.RemoveFileTypeIndexAlpha{
		FileName: fileName,
		FileType: fileType,
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

	parentFileTree := fs.getParentFileTree()
	parentVertex := parentFileTree.RetrieveVertex(fileName)
	if parentVertex != nil {
		return errors.New("File not found in file system")
	}

	// Remove Vertex from Dag
	removeParentVertexAlpha := dag.RemoveVertexAlpha{
		Target: parentVertex,
	}

	updateErr = fs.fileParentRelationships.DataStructure.Update(&removeParentVertexAlpha)
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

func (fs *FileSystem) RetrieveFileContents(fileName string) (string, error) {
	return blob.ReturnBlobContent(fileName)
}

func (fs *FileSystem) LinkHierarchy(parentName string, childName string) error {
	return fs.linkFile(parentName, childName, FILE_RELATIONSHIPS_HIERARCHY)
}

func (fs *FileSystem) LinkDependency(parentName string, childName string) error {
	return fs.linkFile(parentName, childName, FILE_RELATIONSHIP_DEPENDENCY)
}

func (fs *FileSystem) linkFile(parentName string, childName string, relationship string) error {
	if parentName == "" || childName == "" || relationship == "" {
		return nil
	}

	log.Println("Parent: " + parentName + "; Child: " + childName + "; Relationship: " + relationship)
	parentErr := fs.validateFileExists(parentName)
	if parentErr != nil {
		log.Println("Parent cannot be found")
		return parentErr
	}

	childErr := fs.validateFileExists(childName)
	if childErr != nil {
		log.Println("Child cannot be found")
		return childErr
	}

	fileTree := fs.getFileTree()
	parentVertex := fileTree.RetrieveVertex(parentName)
	childVertex := fileTree.RetrieveVertex(childName)

	// Add Edge between parent and child vertex
	addEdgeAlpha := dag.AddEdgeAlpha{
		From:  parentVertex,
		To:    childVertex,
		Label: relationship,
	}

	updateErr := fs.fileRelationShips.DataStructure.Update(&addEdgeAlpha)
	if updateErr != nil {
		log.Println("Error Linking file")
		return updateErr
	}

	parentFileTree := fs.getParentFileTree()
	parentVertex = parentFileTree.RetrieveVertex(parentName)
	childVertex = parentFileTree.RetrieveVertex(childName)

	// Add opposite direction edge
	addOppEdgeAlpha := dag.AddEdgeAlpha{
		To:    parentVertex,
		From:  childVertex,
		Label: relationship,
	}

	updateErr = fs.fileParentRelationships.DataStructure.Update(&addOppEdgeAlpha)
	if updateErr != nil {
		log.Println("Error Parent Linking file" + updateErr.Error())
		return updateErr
	}

	return nil
}

func (fs *FileSystem) UnLinkHierarchy(parentName string, childName string) error {
	return fs.unLinkFile(parentName, childName, FILE_RELATIONSHIPS_HIERARCHY)
}

func (fs *FileSystem) UnLinkDependency(parentName string, childName string) error {
	return fs.unLinkFile(parentName, childName, FILE_RELATIONSHIP_DEPENDENCY)
}

func (fs *FileSystem) unLinkFile(parentName string, childName string, relationship string) error {
	if parentName == "" || childName == "" || relationship == "" {
		return nil
	}

	log.Println("Parent: " + parentName + "; Child: " + childName + "; Relationship: " + relationship)
	parentErr := fs.validateFileExists(parentName)
	if parentErr != nil {
		return parentErr
	}

	childErr := fs.validateFileExists(childName)
	if childErr != nil {
		return childErr
	}

	fileTree := fs.getFileTree()
	parentVertex := fileTree.RetrieveVertex(childName)
	childVertex := fileTree.RetrieveVertex(parentName)

	// Add Edge between parent and child vertex
	removeEdgeAlpha := dag.RemoveEdgeAlpha{
		From:  parentVertex,
		To:    childVertex,
		Label: relationship,
	}

	updateErr := fs.fileRelationShips.DataStructure.Update(&removeEdgeAlpha)
	if updateErr != nil {
		log.Println("Error unlinking file")
		return updateErr
	}

	parentFileTree := fs.getParentFileTree()
	parentVertex = parentFileTree.RetrieveVertex(childName)
	childVertex = parentFileTree.RetrieveVertex(parentName)

	// Add Edge between parent and child vertex
	removeOppEdgeAlpha := dag.RemoveEdgeAlpha{
		To:    parentVertex,
		From:  childVertex,
		Label: relationship,
	}

	updateErr = fs.fileRelationShips.DataStructure.Update(&removeOppEdgeAlpha)
	if updateErr != nil {
		log.Println("Error unlinking parent file")
		return updateErr
	}

	return nil
}

func (fs *FileSystem) ListAllFilesWithTypes() (map[string][]string, error) {
	fileIndex := fs.getFileIndex()
	files, fileErr := fileIndex.RetrieveAllFilesWithTypes()
	if fileErr != nil {
		return nil, fileErr
	}

	for _, v := range files {
		sort.Strings(v)
	}

	return files, nil
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

func (fs *FileSystem) ListRelatedHierarchy(fileName string) ([]string, error) {
	return fs.ListRelatedIssues(fileName, FILE_RELATIONSHIPS_HIERARCHY)
}

func (fs *FileSystem) ListRelatedDependency(fileName string) ([]string, error) {
	return fs.ListRelatedIssues(fileName, FILE_RELATIONSHIP_DEPENDENCY)
}

func (fs *FileSystem) ListRelatedParentDependency(fileName string) ([]string, error) {
	return fs.ListRelatedParents(fileName, FILE_RELATIONSHIP_DEPENDENCY)
}

func (fs *FileSystem) ListRelatedParents(fileName string, fileRelationship string) ([]string, error) {
	childDag := fs.fileParentRelationships.DataStructure.(*dag.Dag)
	vertex := childDag.RetrieveVertex(fileName)
	children := vertex.Children

	var issues []string
	for _, child := range children {
		if fileRelationship != child.Label {
			continue
		}

		issues = append(issues, child.To.ID)
	}

	return issues, nil
}

func (fs *FileSystem) ListRelatedIssues(fileName string, fileRelationship string) ([]string, error) {
	childDag := fs.fileRelationShips.DataStructure.(*dag.Dag)
	vertex := childDag.RetrieveVertex(fileName)
	children := vertex.Children

	var issues []string
	for _, child := range children {
		if fileRelationship != child.Label {
			continue
		}

		issues = append(issues, child.To.ID)
	}

	return issues, nil
}

func (fs *FileSystem) GetFileType(fileName string) (string, error) {
	fileIndex := fs.getFileIndex()
	fileType, typeErr := fileIndex.RetrieveFileType(fileName)
	if typeErr != nil {
		return "", typeErr
	}

	return fileType, nil
}

func (fs *FileSystem) GetFileChildMeta(fileName string) *dag.Vertex {
	return fs.getFileTree().RetrieveVertex(fileName)
}

type FileGraphRenderer struct {
	Fs FileSystem
}

func (fgr FileGraphRenderer) Build(fileName string) (string, error) {
	output := fgr.BuildIssue(fileName, 0)

	return output, nil
}

func (fgr FileGraphRenderer) BuildIssue(vertexID string, lane int) string {
	childrenDag := fgr.Fs.getFileTree()
	issueWithChildren := childrenDag.RetrieveVertex(vertexID)

	head := fgr.AddHeadInLane(lane) + " " + issueWithChildren.ID
	body := fgr.AddBodyInLane(lane)
	dependencies := fgr.BuildDepedencies(issueWithChildren, lane)

	return head + body + dependencies
}

// func (fgr FileGraphRenderer) BuildHierarchy(issueWithChildren *dag.Vertex, lane int) (string) {
// 	if len(issueWithChildren.Children) == 0 {
// 		return ""
// 	}
//
// 	var children string;
//
// 	for _, value := range issueWithChildren.Children {
// 		if value.Label == FILE_RELATIONSHIPS_HIERARCHY {
// 			children += fgr.BuildIssue(value.To.ID, lane)
// 		}
// 	}
//
// 	return children
// }

func (fgr FileGraphRenderer) BuildDepedencies(issueWithChildren *dag.Vertex, lane int) string {
	if len(issueWithChildren.Children) == 0 {
		return ""
	}

	var children string

	for index, value := range issueWithChildren.Children {
		if value.Label == FILE_RELATIONSHIP_DEPENDENCY {
			position := fgr.LanePosition(lane + index)
			children += position + "\\" + fgr.BuildIssue(value.To.ID, lane+index+1)
		}
	}

	return "\n" + children
}

func (fgr FileGraphRenderer) AddBodyInLane(lane int) string {
	laneString := fgr.LanePosition(lane)
	return "\n" + laneString + "|"
}

func (fgr FileGraphRenderer) AddHeadInLane(lane int) string {
	laneString := fgr.LanePosition(lane)
	return "\n" + laneString + "*"
}

func (fgr FileGraphRenderer) LanePosition(lane int) string {
	return strings.Repeat(" ", lane)
}
