package main

import "fmt"

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

func dfs(from *Vertex, to *Vertex) bool {
  if (from.id == to.id) {
    return true
  }

  for _ , value := range to.children {
      if dfs(from, value) {
         return true
      }
  }

  return false
}

type Dag struct {
  vertices map[string]*Vertex
}

func newDag() *Dag {
  return &Dag{
    vertices: make(map[string]*Vertex),
  }
}

func (d Dag) addEdge(from *Vertex, to *Vertex) {
  parent, hasParent := d.vertices[from.id]
  child, hasChild := d.vertices[to.id]

  if !hasParent || !hasChild {
    fmt.Println("From or to Vertex does not exist")
    return
  }

  _ , hasChildEdge := parent.children[to.id]
  if hasChildEdge {
    fmt.Println("Child Edge already exist")
    return
  }

  _ , hasParentEdge := child.parents[from.id]
  if hasParentEdge {
    fmt.Println("Parent Edge already exist")
    return
  }

  hasCycle := dfs(from, to);
  if hasCycle {
     fmt.Println("Cannot add edge as it will create a cycle");
     return
  }

  parent.children[to.id] = to
  child.parents[from.id] = from
}

func (d Dag) removeEdge(from *Vertex, to *Vertex) {
  _, hasToEdge := from.children[to.id]
  if !hasToEdge {
    fmt.Println("Vertex does not exist")
    return
  }

  _, hasFromEdge := to.parents[from.id]
  if !hasFromEdge {
    fmt.Println("Vertex does not exist")
    return
  }

  delete(from.children, to.id)
  delete(to.parents, from.id)
}

func (d Dag) addVertex(in *Vertex) {
  _ , exists := d.vertices[in.id]
  if exists {
    fmt.Println("Vertex already exists")
    return
  }

  d.vertices[in.id] = in
}

func (d Dag) removeVertex(out *Vertex) {
  _ , exists := d.vertices[out.id]
  if !exists {
    fmt.Println("Deleting non existent vertex")
    return
  }

  for _ , value := range out.parents {
    d.removeEdge(value, out)
  }

  for _ , value := range out.children {
    d.removeEdge(out, value)
  }

  delete(d.vertices, out.id)
}