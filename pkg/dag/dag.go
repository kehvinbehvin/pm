package dag

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
	Children []*DirectedEdge // A Vertex can have multiple types of edges connected to it
}

// Cannot store parent address because encoding/gob does not allow recursive data structures
type DirectedEdge struct {
	Label string // Allows us to add meaning to the edge
	To    *Vertex
}

func (d *Dag) Update(alpha common.Alpha) error {
	alphaType := alpha.GetType()
	var error error
	switch alphaType {
	case common.AddVertexAlpha:
		addVertexAlpha := alpha.(*AddVertexAlpha)
		error = d.AddVertex(addVertexAlpha.Target)
		if error == nil {
			log.Println("addVertexAlpha: " + addVertexAlpha.Target.String())
		}
	case common.RemoveVertexAlpha:
		log.Println("removing vertex")
		removeVertexAlpha := alpha.(*RemoveVertexAlpha)
		error = d.RemoveVertex(removeVertexAlpha.Target)
		if error == nil {
			log.Println("removeVertexAlpha: " + removeVertexAlpha.Target.String())
		} else {
			log.Println(error)
		}

	case common.AddEdgeAlpha:
		addEdgeAlpha := alpha.(*AddEdgeAlpha)
		log.Println("addEdgeAlpha: " + addEdgeAlpha.From.String())
		log.Println("addEdgeAlpha: " + addEdgeAlpha.To.String())
		error = d.AddEdge(addEdgeAlpha.From, addEdgeAlpha.To, addEdgeAlpha.Label)
	case common.RemoveEdgeAlpha:
		removeEdgeAlpha := alpha.(*RemoveEdgeAlpha)
		log.Println("removeEdgeAlpha: " + removeEdgeAlpha.From.String())
		log.Println("removeEdgeAlpha: " + removeEdgeAlpha.To.String())
		error = d.RemoveEdge(removeEdgeAlpha.From, removeEdgeAlpha.To, removeEdgeAlpha.Label)
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
		error = d.RemoveEdge(removeEdgeAlpha.From, removeEdgeAlpha.To, removeEdgeAlpha.Label)
	case common.RemoveEdgeAlpha:
		addEdgeAlpha := alpha.(*AddEdgeAlpha)
		error = d.AddEdge(addEdgeAlpha.From, addEdgeAlpha.To, addEdgeAlpha.Label)
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
		valid = d.canAddEdge(addEdgeAlpha.From, addEdgeAlpha.To)
	case common.RemoveEdgeAlpha:
		removeEdgeAlpha := alpha.(*RemoveEdgeAlpha)
		valid = d.canAddEdge(removeEdgeAlpha.From, removeEdgeAlpha.To)
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
		Children: make([]*DirectedEdge, 0),
	}
}

func NewDirectedEdge(to *Vertex, label string) *DirectedEdge {
	return &DirectedEdge{
		To:    to,
		Label: label,
	}
}

func (v *Vertex) String() string {
	var childrenIDs []string
	var childrenRelationships []string

	for _, edgePointer := range v.Children {
		relationship := edgePointer.Label
		childrenRelationships = append(childrenRelationships, relationship)

		child := edgePointer.To.ID
		childrenIDs = append(childrenIDs, child)
	}

	return fmt.Sprintf("Vertex(id: %s, children: %v, relationships: %v)", v.ID, childrenIDs, childrenRelationships)
}

func dfs(from *Vertex, to *Vertex, label string) bool {
	if from.ID == to.ID {
		return true
	}

	for _, value := range to.Children {
		childLabel := value.Label
		if childLabel != label {
			continue
		}

		next := value.To
		if dfs(from, next, label) {
			return true
		}
	}

	return false
}

type AddEdgeAlpha struct {
	From  *Vertex
	To    *Vertex
	Hash  string
	Label string
}

type RemoveEdgeAlpha struct {
	From  *Vertex
	To    *Vertex
	Hash  string
	Label string
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

func (d *Dag) canAddEdge(from *Vertex, to *Vertex) bool {
	_, hasParent := d.Vertices[from.ID]
	_, hasChild := d.Vertices[to.ID]

	return hasParent && hasChild
}

func (d *Dag) isExistingEdge(from *Vertex, to *Vertex, label string) bool {
	parent := d.Vertices[from.ID]
	for _, edgePointer := range parent.Children {
		if edgePointer.Label != label {
			continue
		}

		child := edgePointer.To.ID
		if child == to.ID {
			return true
		}
	}

	return false
}

func (d *Dag) AddEdge(from *Vertex, to *Vertex, label string) error {
	if !d.canAddEdge(from, to) {
		return errors.New("From or to Vertex does not exist")
	}

	if d.isExistingEdge(from, to, label) {
		return errors.New("Edge for this label already exists")
	}

	hasCycle := dfs(from, to, label)
	if hasCycle {
		return errors.New("Cannot add edge as it will create a cycle")
	}

	directedEdge := NewDirectedEdge(to, label)

	from.Children = append(from.Children, directedEdge)
	return nil
}

func (d *Dag) deleteEdge(from *Vertex, to *Vertex, label string) error {
	for index, child := range from.Children {
		if child.Label != label {
			continue
		}

		if child.To.ID != to.ID {
			continue
		}

		from.Children = append(from.Children[:index], from.Children[index+1:]...)
		return nil
	}

	return errors.New("Cannot delete non-existent edge")
}

func (d *Dag) RemoveEdge(from *Vertex, to *Vertex, label string) error {
	if !d.isExistingEdge(from, to, label) {
		return errors.New("Cannot delete non-existent edge")
	}

	deleteErr := d.deleteEdge(from, to, label)
	if deleteErr != nil {
		return deleteErr
	}

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
		log.Println("Deleting non existent vertex")
		return errors.New("Deleting non existent vertex")
	}

	for _, value := range out.Children {
		log.Println("Removing edges")
		valueLabel := value.Label
		value := value.To

		removeErr := d.RemoveEdge(out, value, valueLabel)
		if removeErr != nil {
			log.Println("Error removing related edges " + removeErr.Error())
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
		log.Println("Non existent vertex", vertexID)
		return nil
	}

	return vertex
}

func LoadReconcilableDag(filePath string) common.Reconcilable {
	file, fileErr := os.Open(filePath)

	if fileErr != nil {
		log.Println("Error opening binary file")
		return common.Reconcilable{}
	}
	defer file.Close()

	gob.Register(&Dag{})
	decoder := gob.NewDecoder(file)
	var loadedReconcilable common.Reconcilable
	decodingErr := decoder.Decode(&loadedReconcilable)
	if decodingErr != nil {
		log.Println("Error decoding", decodingErr.Error())
		return common.Reconcilable{}
	}

	return loadedReconcilable
}
