package main

import (
	// "github/pm/cmd"
	"fmt"
)

type Vertex struct {
  id string
  children map[string]*Vertex
  parents map[string]*Vertex
}

func newVertex(id string) *Vertex {
  return &Vertex{
    id: id,
    children: make(map[string]*Vertex),
    parents: make(map[string]*Vertex),
  }
}

func (d Dag) addEdge(from *Vertex, to *Vertex) {
  parent, hasParent := d.vertices[from.id]
  child, hasChild := d.vertices[to.id]

  if !hasParent || !hasChild {
    fmt.Printf("From or to Vertex does not exist")  
    return 
  }

  _ , hasChildEdge := parent.children[to.id]
  if hasChildEdge {
    fmt.Printf("Child Edge already exist")
    return
  }

  _ , hasParentEdge := child.parents[from.id]
  if hasParentEdge {
    fmt.Printf("Parent Edge already exist")
    return
  }

  parent.children[to.id] = to
  child.parents[from.id] = from
}

func (v *Vertex) String() string {
    var childrenIDs, parentsIDs []string
    for childID := range v.children {
        childrenIDs = append(childrenIDs, childID)
    }
    for parentID := range v.parents {
        parentsIDs = append(parentsIDs, parentID)
    }

    return fmt.Sprintf("Vertex(id: %s, children: %v, parents: %v)", v.id, childrenIDs, parentsIDs)
}

type Dag struct {
  vertices map[string]*Vertex
}

func newDag() *Dag {
  return &Dag{
    vertices: make(map[string]*Vertex),
  }
}

func (d Dag) addVertex(in *Vertex) {
  _ , exists := d.vertices[in.id]
  if exists {
    fmt.Println("Vertex already exists")
    return
  }

  d.vertices[in.id] = in
}

func main() {
  // cmd.Execute()
  pmDag := newDag(); 

  vertex1 := newVertex("Test 1")
  vertex2 := newVertex("Test 2")
  vertex3 := newVertex("Test 2")

  pmDag.addVertex(vertex1);
  pmDag.addVertex(vertex2);
  pmDag.addVertex(vertex3);
  pmDag.addEdge(vertex1, vertex2)

  for _, value := range pmDag.vertices {
    fmt.Println(value);
  }
}
