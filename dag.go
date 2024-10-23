package main

import (
	"crypto/sha1"
	"encoding/gob"
	"errors"
	"fmt"
	"os"
	"strings"
)

const (
	addVertex    byte = 1
	removeVertex byte = 2
	addEdge      byte = 3
	removeEdge   byte = 4
)

type Vertex struct {
	ID       string
	Children map[string]*Vertex
}

type Delta interface {
	Set()
	UnSet()
	GetId() string
	String() string
	GetParent(*DeltaTree) *Delta
}

type VertexDelta struct {
	Id          string
	Operation   byte
	Vertex      *Vertex
	ParentDelta string
}

func (vd *VertexDelta) Set() {

}

func (vd *VertexDelta) UnSet() {

}

func (vd *VertexDelta) GetId() string {
	return vd.Id
}

func (vd *VertexDelta) GetParent(dt *DeltaTree) *Delta {
	parentId := vd.ParentDelta
	delta, ok := dt.Tree[parentId]
	if !ok {
		fmt.Println("Cannot find parent")
		return nil
	}

	return delta
}

func (vd *VertexDelta) String() string {
	return fmt.Sprintf("VertexDelta(Id: %s, Operation: %d, Vertex: %s, ParentDelta: %s)", vd.Id, vd.Operation, vd.Vertex.ID, vd.ParentDelta)
}

type EdgeDelta struct {
	Id          string
	Operation   byte
	Parent      *Vertex
	Child       *Vertex
	ParentDelta string
}

func (ed *EdgeDelta) Set() {
}

func (ed *EdgeDelta) UnSet() {
}

func (ed *EdgeDelta) GetId() string {
	return ed.Id
}

func (ed *EdgeDelta) GetParent(dt *DeltaTree) *Delta {
	parentId := ed.ParentDelta
	delta, ok := dt.Tree[parentId]
	if !ok {
		fmt.Println("Cannot find parent")
		return nil
	}

	return delta
}

func (ed *EdgeDelta) String() string {
	return fmt.Sprintf("EdgeDelta(Id: %s, Operation: %d, Parent: %s, Child: %s)", ed.Id, ed.Operation, ed.Parent.ID, ed.Child.ID)
}

type DeltaTree struct {
	Tree          map[string]*Delta
	RemotePointer string
	LocalPointer  string
	ParentDelta   string
}

func NewDeltaTree() *DeltaTree {
	return &DeltaTree{
		Tree:          make(map[string]*Delta),
		RemotePointer: "",
		LocalPointer:  "",
	}
}

func (dt *DeltaTree) Checkout(nodeHash string) error {
	_, ok := dt.Tree[nodeHash]
	if !ok {
		return errors.New("No hash found")
	}

	return nil
}

func (dt *DeltaTree) Push(delta Delta) error {
	deltaId := delta.GetId()
	_, ok := dt.Tree[deltaId]
	if ok {
		return errors.New("Existing delta found")
	}

	dt.Tree[deltaId] = &delta
	dt.LocalPointer = deltaId

	return nil
}

func (dt *DeltaTree) String() string {
	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("DeltaTree (remotePointer: %s, localPointer: %s)\n", dt.RemotePointer, dt.LocalPointer))
	builder.WriteString("Deltas:\n")

	for id, delta := range dt.Tree {
		builder.WriteString(fmt.Sprintf("ID: %s, Delta: %s\n", id, *delta))
	}

	return builder.String()
}

func newVertex(id string) *Vertex {
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

type Dag struct {
	Id       string
	Vertices map[string]*Vertex
}

func newDag(fileName string) *Dag {
	return &Dag{
		Id:       fileName,
		Vertices: make(map[string]*Vertex),
	}
}

func (d *Dag) addEdge(from *Vertex, to *Vertex) error {
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

func (dt *DeltaTree) addEdgeEvent(parent *Vertex, child *Vertex) error {
	parentHash := sha1.Sum([]byte(parent.ID))
	parentStr := fmt.Sprintf("%x", parentHash[:])

	childHash := sha1.Sum([]byte(child.ID))
	childStr := fmt.Sprintf("%x", childHash[:])

	opHash := sha1.Sum([]byte{addEdge})
	opStr := fmt.Sprintf("%x", opHash[:])

	prevHash := dt.LocalPointer

	uuid := parentStr + childStr + opStr + prevHash
	uuidHash := sha1.Sum([]byte(uuid))
	uuidStr := fmt.Sprintf("%x", uuidHash[:])

	delta := &EdgeDelta{
		Id:          uuidStr,
		Operation:   addEdge,
		Parent:      parent,
		Child:       child,
		ParentDelta: prevHash,
	}

	dt.Push(delta)

	return nil
}

func (dt *DeltaTree) removeEdgeEvent(parent *Vertex, child *Vertex) error {
	parentHash := sha1.Sum([]byte(parent.ID))
	parentStr := fmt.Sprintf("%x", parentHash[:])

	childHash := sha1.Sum([]byte(child.ID))
	childStr := fmt.Sprintf("%x", childHash[:])

	opHash := sha1.Sum([]byte{removeEdge})
	opStr := fmt.Sprintf("%x", opHash[:])

	prevHash := dt.LocalPointer

	uuid := parentStr + childStr + opStr + prevHash
	uuidHash := sha1.Sum([]byte(uuid))
	uuidStr := fmt.Sprintf("%x", uuidHash[:])

	delta := &EdgeDelta{
		Id:          uuidStr,
		Operation:   removeEdge,
		Parent:      parent,
		Child:       child,
		ParentDelta: prevHash,
	}

	dt.Push(delta)

	return nil
}

func (d *Dag) removeEdge(from *Vertex, to *Vertex) {
	_, hasToEdge := from.Children[to.ID]
	if !hasToEdge {
		fmt.Println("Vertex does not exist")
		return
	}

	delete(from.Children, to.ID)
}

func (d *Dag) addVertex(in *Vertex) (error) {
	_, exists := d.Vertices[in.ID]
	if exists {
		return errors.New("Vertex already exists")
	}

	d.Vertices[in.ID] = in
	return nil
}

func (dt *DeltaTree) addVertexEvent(vertex *Vertex) error {
	idHash := sha1.Sum([]byte(vertex.ID))
	hashStr := fmt.Sprintf("%x", idHash[:])

	opHash := sha1.Sum([]byte{addVertex})
	opStr := fmt.Sprintf("%x", opHash[:])

	prevHash := dt.LocalPointer

	uuid := hashStr + opStr + prevHash
	uuidHash := sha1.Sum([]byte(uuid))
	uuidStr := fmt.Sprintf("%x", uuidHash[:])

	delta := &VertexDelta{
		Id:          uuidStr,
		Operation:   addVertex,
		Vertex:      vertex,
		ParentDelta: prevHash,
	}

	dt.Push(delta)

	return nil
}

func (dt *DeltaTree) removeVertexEvent(vertex *Vertex) error {
	idHash := sha1.Sum([]byte(vertex.ID))
	hashStr := fmt.Sprintf("%x", idHash[:])

	opHash := sha1.Sum([]byte{addVertex})
	opStr := fmt.Sprintf("%x", opHash[:])

	prevHash := dt.LocalPointer

	uuid := hashStr + opStr + prevHash
	uuidHash := sha1.Sum([]byte(uuid))
	uuidStr := fmt.Sprintf("%x", uuidHash[:])

	delta := &VertexDelta{
		Id:          uuidStr,
		Operation:   removeVertex,
		Vertex:      vertex,
		ParentDelta: prevHash,
	}

	dt.Push(delta)

	return nil
}

func (d *Dag) removeVertex(out *Vertex) {
	_, exists := d.Vertices[out.ID]
	if !exists {
		fmt.Println("Deleting non existent vertex")
		return
	}

	for _, value := range out.Children {
		d.removeEdge(out, value)
	}

	delete(d.Vertices, out.ID)
}

func (d *Dag) retrieveVertex(vertexID string) *Vertex {
	vertex, exists := d.Vertices[vertexID]
	if !exists {
		fmt.Println("Non existent vertex", vertexID)
		return nil
	}

	return vertex
}

func (dt *DeltaTree) SaveDelta() {
	file, err := os.Create("./.pm/delta")
	if err != nil {
		fmt.Printf(err.Error())
		fmt.Println("Error creating file")
		return
	}
	defer file.Close()

	// Register the types with gob
	gob.Register(&VertexDelta{})
	gob.Register(&EdgeDelta{})

	encoder := gob.NewEncoder(file)
	encodingErr := encoder.Encode(dt)
	if encodingErr != nil {
		fmt.Printf(encodingErr.Error())
		fmt.Println("Error encoding delta")
		return
	}
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

func LoadDelta() *DeltaTree {
	file, fileErr := os.Open("./.pm/delta")

	if fileErr != nil {
		fmt.Println("Error opening binary file")
		return nil
	}
	defer file.Close()

	gob.Register(&VertexDelta{})
	gob.Register(&EdgeDelta{})

	decoder := gob.NewDecoder(file)

	var deltaTree *DeltaTree
	decodingErr := decoder.Decode(&deltaTree)
	if decodingErr != nil {
		fmt.Println("Error decoding delta tree")
		return nil
	}

	return deltaTree

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
