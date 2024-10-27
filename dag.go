package main

import (
	"encoding/gob"
	"errors"
	"fmt"
	"os"
)

type Vertex struct {
	ID       string
	Children map[string]*Vertex
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

func (d *Dag) removeEdge(from *Vertex, to *Vertex) error {
	_, hasToEdge := from.Children[to.ID]
	if !hasToEdge {
		fmt.Println("Vertex does not exist")
		return errors.New("Vertex does not exist")
	}

	delete(from.Children, to.ID)
	return nil
}

func (d *Dag) addVertex(in *Vertex) error {
	_, exists := d.Vertices[in.ID]
	if exists {
		return errors.New("Vertex already exists")
	}

	d.Vertices[in.ID] = in
	return nil
}

func (d *Dag) removeVertex(out *Vertex, deltaTree *DeltaTree) error {
	_, exists := d.Vertices[out.ID]
	if !exists {
		fmt.Println("Deleting non existent vertex")
		return errors.New("Deleting non existent vertex")
	}

	for _, value := range out.Children {
		removeErr := d.removeEdge(out, value)
		if removeErr != nil {
			return removeErr
		}
		deltaTree.removeEdgeEvent(out, value)
	}

	delete(d.Vertices, out.ID)
	return nil
}

func (d *Dag) retrieveVertex(vertexID string) *Vertex {
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
