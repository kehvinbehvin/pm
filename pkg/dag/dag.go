package dag

import (
	"encoding/gob"
	"errors"
	"fmt"
	"os"

	"github/pm/pkg/common"
)


func NewReconcilableDag(storageKey string) common.Reconcilable {
	dagAlphaList := common.NewAlphaList();
	dagStorage := NewDag(storageKey);

	return common.Reconcilable{
		AlphaList: dagAlphaList,
		DataStructure: dagStorage,
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

func (d *Dag) Update(alpha common.Alpha) (error) {
	alphaType := alpha.GetType();
	var error error
	switch alphaType {
	case common.AddVertexAlpha:
		addVertexAlpha := alpha.(*AddVertexAlpha)
		error = d.AddVertex(addVertexAlpha.target)
	case common.RemoveVertexAlpha:
		removeVertexAlpha := alpha.(*RemoveVertexAlpha)
		error = d.RemoveVertex(removeVertexAlpha.target)
	case common.AddEdgeAlpha:
		addEdgeAlpha := alpha.(*AddEdgeAlpha)
		error = d.AddEdge(addEdgeAlpha.to, addEdgeAlpha.from)
	case common.RemoveEdgeAlpha:
		removeEdgeAlpha := alpha.(*RemoveEdgeAlpha)
		error = d.RemoveEdge(removeEdgeAlpha.to, removeEdgeAlpha.from)
	}

	return error
}


func (d *Dag) Rewind(alpha common.Alpha) (error) {
	alphaType := alpha.GetType();
	var error error
	switch alphaType {
	case common.AddVertexAlpha:
		removeVertexAlpha := alpha.(*RemoveVertexAlpha)
		error = d.RemoveVertex(removeVertexAlpha.target)
	case common.RemoveVertexAlpha:
		addVertexAlpha := alpha.(*AddVertexAlpha)
		error = d.AddVertex(addVertexAlpha.target)
	case common.AddEdgeAlpha:
		removeEdgeAlpha := alpha.(*RemoveEdgeAlpha)
		error = d.RemoveEdge(removeEdgeAlpha.to, removeEdgeAlpha.from)
	case common.RemoveEdgeAlpha:
		addEdgeAlpha := alpha.(*AddEdgeAlpha)
		error = d.AddEdge(addEdgeAlpha.to, addEdgeAlpha.from)
	}

	return error
}

func (d *Dag) Validate(alpha common.Alpha) (bool) {
	alphaType := alpha.GetType();
	var valid bool

	switch alphaType {
	case common.AddVertexAlpha:
		addVertexAlpha := alpha.(*AddVertexAlpha)
		valid = d.HasVertex(addVertexAlpha.target.ID)
	case common.RemoveVertexAlpha:
		removeVertexAlpha := alpha.(*RemoveVertexAlpha)
		valid = !d.HasVertex(removeVertexAlpha.target.ID)
	case common.AddEdgeAlpha:
		addEdgeAlpha := alpha.(*AddEdgeAlpha)
		valid = d.HasEdge(addEdgeAlpha.from, addEdgeAlpha.to)
	case common.RemoveEdgeAlpha:
		removeEdgeAlpha := alpha.(*RemoveEdgeAlpha)
		valid = d.HasEdge(removeEdgeAlpha.from, removeEdgeAlpha.to)
	}

	return valid
}


func NewDag(fileName string) *Dag {
	return &Dag{
		Id:       fileName,
		Vertices: make(map[string]*Vertex),
	}
}

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
	from *Vertex
	to *Vertex
}

type RemoveEdgeAlpha struct {
	from *Vertex
	to *Vertex
}


func (aea *AddEdgeAlpha) GetType() byte {
	return common.AddEdgeAlpha
}

func (aea *AddEdgeAlpha) GetId() string {
	return aea.to.ID + aea.from.ID + string(common.AddEdgeAlpha)
}

func (rea *RemoveEdgeAlpha) GetType() byte {
	return common.RemoveEdgeAlpha
}

func (rea *RemoveEdgeAlpha) GetId() string {
	return rea.to.ID + rea.from.ID + string(common.AddEdgeAlpha)
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
	target *Vertex
}

type RemoveVertexAlpha struct {
	target *Vertex
}

func (ava *AddVertexAlpha) GetType() byte {
	return common.AddVertexAlpha;
}

func (ava *AddVertexAlpha) GetId() string {
	return ava.target.ID
}

func (rva *RemoveVertexAlpha) GetType() byte {
	return common.RemoveVertexAlpha
}

func (rvd *RemoveVertexAlpha) GetId() string {
	return rvd.target.ID
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

func (d *Dag) SaveDag() {
	file, err := os.Create("./.pm/dag/" + d.Id)
	if err != nil {
		fmt.Printf(err.Error())
		fmt.Println("Error creating file")
		return
	}
	defer file.Close()

	encoder := gob.NewEncoder(file)
	encodingErr := encoder.Encode(d)
	if encodingErr != nil {
		fmt.Println("Error encoding dag")
		return
	}
}

func LoadDag(fileName string) *Dag {
	file, fileErr := os.Open("./.pm/dag/" + fileName)

	if fileErr != nil {
		fmt.Println("Error opening binary file")
		return nil
	}
	defer file.Close()

	decoder := gob.NewDecoder(file)

	var loadedDag *Dag
	decodingErr := decoder.Decode(&loadedDag)
	if decodingErr != nil {
		fmt.Println("Error decoding dag")
		return nil
	}

	return loadedDag
}
