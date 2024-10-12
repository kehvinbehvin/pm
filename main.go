package main

import (
	// "github/pm/cmd"
	"fmt"
)

func main() {
  // cmd.Execute()
  pmDag := newDag(); 

  vertex1 := newVertex("Test 1")
  vertex2 := newVertex("Test 2")
  vertex3 := newVertex("Test 3")

  pmDag.addVertex(vertex1);
  pmDag.addVertex(vertex2);
  pmDag.addVertex(vertex3);
  pmDag.addEdge(vertex1, vertex2)
  pmDag.addEdge(vertex2, vertex3)
  pmDag.addEdge(vertex3, vertex3)

  for _, value := range pmDag.vertices {
    fmt.Println(value);
  }

  pmDag.removeEdge(vertex1, vertex3)
  for _, value := range pmDag.vertices {
    fmt.Println(value);
  }
}
