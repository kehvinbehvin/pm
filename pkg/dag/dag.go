package dag

import (
	"crypto/sha1"
	"encoding/gob"
	"errors"
	"fmt"
	"os"

	"github/pm/pkg/common"
)

/**
How to use:
1. Create a file save in the file system to store the dag
2. The name of the file is the storage key
3. Use NewReconcilableDag to create new Dag
4. Create Add/Remove Vertex/Edge Alphas and call Update on Reconcilable
5. Save Reconcilable once Update is completed


Rewind and Validate are future features for merging strategies
Not need now
*/

// Provide the storage key to identify instances of dags
func NewReconcilableDag(storageKey string) common.Reconcilable {
	dagAlphaList := common.NewAlphaList()
	dagStorage := NewDag(storageKey)
	filePath := "./.pm/dag/" + storageKey

	return common.Reconcilable{
		AlphaList:     dagAlphaList,
		DataStructure: dagStorage,
		FilePath:      filePath,
	}
}

type Dag struct {
	Id       string
	Vertices map[string]*Vertex
}

type Vertex struct {
	ID       string
	Children map[string]*Vertex
}

func (d *Dag) Update(alpha common.Alpha) error {
	alphaType := alpha.GetType()
	var error error
	switch alphaType {
	case common.AddVertexAlpha:
		addVertexAlpha := alpha.(*AddVertexAlpha)
		error = d.AddVertex(addVertexAlpha.Target)
	case common.RemoveVertexAlpha:
		removeVertexAlpha := alpha.(*RemoveVertexAlpha)
		error = d.RemoveVertex(removeVertexAlpha.Target)
	case common.AddEdgeAlpha:
		addEdgeAlpha := alpha.(*AddEdgeAlpha)
		error = d.AddEdge(addEdgeAlpha.To, addEdgeAlpha.From)
	case common.RemoveEdgeAlpha:
		removeEdgeAlpha := alpha.(*RemoveEdgeAlpha)
		error = d.RemoveEdge(removeEdgeAlpha.To, removeEdgeAlpha.From)
	}

	return error
}

func (d *Dag) Rewind(alpha common.Alpha) error {
	alphaType := alpha.GetType()
	var error error
	switch alphaType {
	case common.AddVertexAlpha:
		removeVertexAlpha := alpha.(*RemoveVertexAlpha)
		error = d.RemoveVertex(removeVertexAlpha.Target)
	case common.RemoveVertexAlpha:
		addVertexAlpha := alpha.(*AddVertexAlpha)
		error = d.AddVertex(addVertexAlpha.Target)
	case common.AddEdgeAlpha:
		removeEdgeAlpha := alpha.(*RemoveEdgeAlpha)
		error = d.RemoveEdge(removeEdgeAlpha.To, removeEdgeAlpha.From)
	case common.RemoveEdgeAlpha:
		addEdgeAlpha := alpha.(*AddEdgeAlpha)
		error = d.AddEdge(addEdgeAlpha.To, addEdgeAlpha.From)
	}

	return error
}

func (d *Dag) Validate(alpha common.Alpha) bool {
	alphaType := alpha.GetType()
	var valid bool

	switch alphaType {
	case common.AddVertexAlpha:
		addVertexAlpha := alpha.(*AddVertexAlpha)
		valid = d.HasVertex(addVertexAlpha.Target.ID)
	case common.RemoveVertexAlpha:
		removeVertexAlpha := alpha.(*RemoveVertexAlpha)
		valid = !d.HasVertex(removeVertexAlpha.Target.ID)
	case common.AddEdgeAlpha:
		addEdgeAlpha := alpha.(*AddEdgeAlpha)
		valid = d.HasEdge(addEdgeAlpha.From, addEdgeAlpha.To)
	case common.RemoveEdgeAlpha:
		removeEdgeAlpha := alpha.(*RemoveEdgeAlpha)
		valid = d.HasEdge(removeEdgeAlpha.From, removeEdgeAlpha.To)
	}

	return valid
}

func NewDag(fileName string) *Dag {
	return &Dag{
		Id:       fileName,
		Vertices: make(map[string]*Vertex),
	}
}

// Id provided must be unique in nature.
func NewVertex(id string) *Vertex {
	return &Vertex{
		ID:       id,
		Children: make(map[string]*Vertex),
	}
}

func (v *Vertex) String() string {
	var childrenIDs []string
	for childID := range v.Children {
		childrenIDs = append(childrenIDs, childID)
	}

	return fmt.Sprintf("Vertex(id: %s, children: %v)", v.ID, childrenIDs)
}

func dfs(from *Vertex, to *Vertex) bool {
	if from.ID == to.ID {
		return true
	}

	for _, value := range to.Children {
		if dfs(from, value) {
			return true
		}
	}

	return false
}

type AddEdgeAlpha struct {
	From *Vertex
	To   *Vertex
	Hash string
}

type RemoveEdgeAlpha struct {
	From *Vertex
	To   *Vertex
	Hash string
}

func (aea *AddEdgeAlpha) GetType() byte {
	return common.AddEdgeAlpha
}

func (aea *AddEdgeAlpha) GetId() string {
	return aea.To.ID + aea.From.ID + string(common.AddEdgeAlpha)
}

func (aea *AddEdgeAlpha) GetHash() string {
	return aea.Hash
}

// This is mean to capture the state that the alpha was used to update
// the underlying datastructure
func (aea *AddEdgeAlpha) SetHash(lastAlpha common.Alpha) {
	prevAlphaHash := lastAlpha.GetHash()
	currentHash := sha1.Sum([]byte(aea.GetId() + prevAlphaHash))
	currentHashStr := fmt.Sprintf("%x", currentHash[:])
	aea.Hash = currentHashStr
}

func (rea *RemoveEdgeAlpha) GetType() byte {
	return common.RemoveEdgeAlpha
}

func (rea *RemoveEdgeAlpha) GetId() string {
	return rea.To.ID + rea.From.ID + string(common.AddEdgeAlpha)
}

func (rea *RemoveEdgeAlpha) GetHash() string {
	return rea.Hash
}

// This is mean to capture the state that the alpha was used to update
// the underlying datastructure
func (rea *RemoveEdgeAlpha) SetHash(lastAlpha common.Alpha) {
	prevAlphaHash := lastAlpha.GetHash()
	currentHash := sha1.Sum([]byte(rea.GetId() + prevAlphaHash))
	currentHashStr := fmt.Sprintf("%x", currentHash[:])
	rea.Hash = currentHashStr
}

func (d *Dag) HasEdge(from *Vertex, to *Vertex) bool {
	_, hasParent := d.Vertices[from.ID]
	_, hasChild := d.Vertices[to.ID]

	return hasParent && hasChild
}

func (d *Dag) AddEdge(from *Vertex, to *Vertex) error {
	parent, hasParent := d.Vertices[from.ID]
	_, hasChild := d.Vertices[to.ID]

	if !hasParent || !hasChild {
		return errors.New("From or to Vertex does not exist")
	}

	_, hasChildEdge := parent.Children[to.ID]
	if hasChildEdge {
		return errors.New("Child Edge already exist")
	}

	hasCycle := dfs(from, to)
	if hasCycle {
		return errors.New("Cannot add edge as it will create a cycle")
	}

	from.Children[to.ID] = to
	return nil
}

func (d *Dag) RemoveEdge(from *Vertex, to *Vertex) error {
	_, hasToEdge := from.Children[to.ID]
	if !hasToEdge {
		fmt.Println("Vertex does not exist")
		return errors.New("Vertex does not exist")
	}

	delete(from.Children, to.ID)
	return nil
}

type AddVertexAlpha struct {
	Target *Vertex
	Hash   string
}

type RemoveVertexAlpha struct {
	Target *Vertex
	Hash   string
}

func (ava *AddVertexAlpha) GetType() byte {
	return common.AddVertexAlpha
}

func (ava *AddVertexAlpha) GetId() string {
	return ava.Target.ID
}

func (ava *AddVertexAlpha) GetHash() string {
	return ava.Hash
}

func (ava *AddVertexAlpha) SetHash(lastAlpha common.Alpha) {
	prevAlphaHash := lastAlpha.GetHash()
	currentHash := sha1.Sum([]byte(ava.GetId() + prevAlphaHash))
	currentHashStr := fmt.Sprintf("%x", currentHash[:])
	ava.Hash = currentHashStr
}

func (rva *RemoveVertexAlpha) GetType() byte {
	return common.RemoveVertexAlpha
}

func (rvd *RemoveVertexAlpha) GetId() string {
	return rvd.Target.ID
}

func (rvd *RemoveVertexAlpha) GetHash() string {
	return rvd.Hash
}

func (rvd *RemoveVertexAlpha) SetHash(lastAlpha common.Alpha) {
	prevAlphaHash := lastAlpha.GetHash()
	currentHash := sha1.Sum([]byte(rvd.GetId() + prevAlphaHash))
	currentHashStr := fmt.Sprintf("%x", currentHash[:])
	rvd.Hash = currentHashStr

}

func (d *Dag) AddVertex(in *Vertex) error {
	_, exists := d.Vertices[in.ID]
	if exists {
		return errors.New("Vertex already exists")
	}

	d.Vertices[in.ID] = in
	return nil
}

func (d *Dag) RemoveVertex(out *Vertex) error {
	_, exists := d.Vertices[out.ID]
	if !exists {
		fmt.Println("Deleting non existent vertex")
		return errors.New("Deleting non existent vertex")
	}

	for _, value := range out.Children {
		removeErr := d.RemoveEdge(out, value)
		if removeErr != nil {
			return removeErr
		}
	}

	delete(d.Vertices, out.ID)
	return nil
}

func (d *Dag) HasVertex(vertexID string) bool {
	_, exists := d.Vertices[vertexID]
	return exists
}

func (d *Dag) RetrieveVertex(vertexID string) *Vertex {
	vertex, exists := d.Vertices[vertexID]
	if !exists {
		fmt.Println("Non existent vertex", vertexID)
		return nil
	}

	return vertex
}

func LoadReconcilableDag(filePath string) common.Reconcilable {
	file, fileErr := os.Open(filePath)

	if fileErr != nil {
		fmt.Println("Error opening binary file")
		return common.Reconcilable{}
	}
	defer file.Close()

	gob.Register(&Dag{})
	decoder := gob.NewDecoder(file)
	var loadedReconcilable common.Reconcilable
	decodingErr := decoder.Decode(&loadedReconcilable)
	if decodingErr != nil {
		fmt.Println("Error decoding", decodingErr.Error())
		return common.Reconcilable{}
	}

	return loadedReconcilable
}
