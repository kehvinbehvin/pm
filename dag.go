package main

import (
	"encoding/gob"
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

func (d *Dag) addEdge(from *Vertex, to *Vertex) {
	parent, hasParent := d.Vertices[from.ID]
	_, hasChild := d.Vertices[to.ID]

	if !hasParent || !hasChild {
		fmt.Println("From or to Vertex does not exist")
		return
	}

	_, hasChildEdge := parent.Children[to.ID]
	if hasChildEdge {
		fmt.Println("Child Edge already exist")
		return
	}

	hasCycle := dfs(from, to)
	if hasCycle {
		fmt.Println("Cannot add edge as it will create a cycle")
		return
	}

	from.Children[to.ID] = to
}

func (d *Dag) removeEdge(from *Vertex, to *Vertex) {
	_, hasToEdge := from.Children[to.ID]
	if !hasToEdge {
		fmt.Println("Vertex does not exist")
		return
	}

	delete(from.Children, to.ID)
}

func (d *Dag) addVertex(in *Vertex) {
	_, exists := d.Vertices[in.ID]
	if exists {
		fmt.Println("Vertex already exists")
		return
	}

	d.Vertices[in.ID] = in
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
		fmt.Println("Non existent vertex")
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
		fmt.Printf(encodingErr.Error())
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
