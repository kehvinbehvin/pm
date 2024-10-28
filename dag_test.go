package main

import (
	"os"
	"testing"
)

// Test creating a new DAG
func TestCreateDag(t *testing.T) {
	dag := newDag("testDag")
	if dag == nil {
		t.Errorf("Expected DAG to be created, but got nil")
	}
}

// Test creating a new Vertex
func TestCreateVertex(t *testing.T) {
	v := newVertex("A")
	if v.ID != "A" {
		t.Errorf("Expected vertex ID to be 'A', but got '%s'", v.ID)
	}
	if len(v.Children) != 0 {
		t.Errorf("Expected vertex to have no children")
	}
}

// Test adding a vertex to the DAG
func TestAddVertex(t *testing.T) {
	dag := newDag("testDag")
	v := newVertex("A")
	dag.addVertex(v)

	if _, exists := dag.Vertices["A"]; !exists {
		t.Errorf("Expected vertex 'A' to be added to the DAG")
	}
}

// Test trying to add the same vertex twice
func TestAddSameVertexTwice(t *testing.T) {
	dag := newDag("testDag")
	v := newVertex("A")
	dag.addVertex(v)
	dag.addVertex(v) // Attempting to add the same vertex

	if len(dag.Vertices) != 1 {
		t.Errorf("Expected only one vertex in the DAG, but found %d", len(dag.Vertices))
	}
}

// Test removing a vertex from the DAG
func TestRemoveVertex(t *testing.T) {
	dag := newDag("testDag")
	v := newVertex("A")
	dag.addVertex(v)

	dag.removeVertex(v, false)

	if _, exists := dag.Vertices["A"]; exists {
		t.Errorf("Expected vertex 'A' to be removed from the DAG")
	}
}

// Test removing a non-existent vertex
func TestRemoveNonExistentVertex(t *testing.T) {
	dag := newDag("testDag")
	v := newVertex("A")
	dag.removeVertex(v, false)

	if len(dag.Vertices) != 0 {
		t.Errorf("Expected DAG to have no vertices, but found %d", len(dag.Vertices))
	}
}

// Test adding an edge between two vertices
func TestAddEdge(t *testing.T) {
	dag := newDag("testDag")
	v1 := newVertex("A")
	v2 := newVertex("B")

	dag.addVertex(v1)
	dag.addVertex(v2)

	dag.addEdge(v1, v2)

	if _, exists := v1.Children[v2.ID]; !exists {
		t.Errorf("Expected vertex 'B' to be a child of vertex 'A'")
	}
}

// Test adding an edge that would create a cycle
func TestAddEdgeWithCycle(t *testing.T) {
	dag := newDag("testDag")
	v1 := newVertex("A")
	v2 := newVertex("B")

	dag.addVertex(v1)
	dag.addVertex(v2)

	dag.addEdge(v1, v2)

	// Try to add an edge from v2 to v1, which would create a cycle
	dag.addEdge(v2, v1)

	if _, exists := v2.Children[v1.ID]; exists {
		t.Errorf("Expected DAG to prevent cycle when adding an edge from 'B' to 'A'")
	}
}

// Test adding an edge for non-existent vertices
func TestAddEdgeForNonExistentVertex(t *testing.T) {
	dag := newDag("testDag")
	v1 := newVertex("A")
	v2 := newVertex("B")

	dag.addVertex(v1)

	// Try to add an edge where vertex B doesn't exist in the DAG
	dag.addEdge(v1, v2)

	if _, exists := v1.Children[v2.ID]; exists {
		t.Errorf("Expected no edge to be added because vertex 'B' does not exist")
	}
}

// Test removing an edge between two vertices
func TestRemoveEdge(t *testing.T) {
	dag := newDag("testDag")
	v1 := newVertex("A")
	v2 := newVertex("B")

	dag.addVertex(v1)
	dag.addVertex(v2)
	dag.addEdge(v1, v2)

	dag.removeEdge(v1, v2)

	if _, exists := v1.Children[v2.ID]; exists {
		t.Errorf("Expected edge between 'A' and 'B' to be removed")
	}
}

// Test removing a non-existent edge
func TestRemoveNonExistentEdge(t *testing.T) {
	dag := newDag("testDag")
	v1 := newVertex("A")
	v2 := newVertex("B")

	dag.addVertex(v1)
	dag.addVertex(v2)

	// Try to remove an edge that does not exist
	dag.removeEdge(v1, v2)

	if len(v1.Children) != 0 {
		t.Errorf("Expected no edges to exist between 'A' and 'B'")
	}
}

// Test DFS to check if a path exists between two vertices
func TestDfs(t *testing.T) {
	dag := newDag("testDag")
	v1 := newVertex("A")
	v2 := newVertex("B")
	v3 := newVertex("C")

	dag.addVertex(v1)
	dag.addVertex(v2)
	dag.addVertex(v3)

	dag.addEdge(v1, v2)
	dag.addEdge(v2, v3)

	found := dfs(v3, v1)

	if !found {
		t.Errorf("Expected DFS to find a path from 'A' to 'C'")
	}
}

// Test DFS where no path exists
func TestDfsNoPath(t *testing.T) {
	dag := newDag("testDag")
	v1 := newVertex("A")
	v2 := newVertex("B")
	v3 := newVertex("C")

	dag.addVertex(v1)
	dag.addVertex(v2)
	dag.addVertex(v3)

	// No edge between v1 and v3
	found := dfs(v1, v3)

	if found {
		t.Errorf("Expected DFS not to find a path from 'A' to 'C'")
	}
}

// Test saving a DAG to disk
func TestSaveDag(t *testing.T) {
	dag := newDag("testDag")
	v1 := newVertex("A")
	v2 := newVertex("B")

	dag.addVertex(v1)
	dag.addVertex(v2)
	dag.addEdge(v1, v2)

	// Save the DAG to disk
	dag.SaveDag()

	// Check if the file was created
	if _, err := os.Stat("./.pm/dag/testDag"); os.IsNotExist(err) {
		t.Errorf("Expected the file 'testDag' to exist, but it does not")
	}
}

// Test loading a DAG from disk
func TestLoadDag(t *testing.T) {
	dag := newDag("testDag")
	v1 := newVertex("A")
	v2 := newVertex("B")
	v3 := newVertex("C")

	dag.addVertex(v1)
	dag.addVertex(v2)
	dag.addVertex(v3)
	dag.addEdge(v1, v2)
	dag.addEdge(v2, v3)

	// Save the DAG first
	dag.SaveDag()

	// Load the DAG from the file
	loadedDag := LoadDag("testDag")
	if loadedDag == nil {
		t.Errorf("Failed to load the DAG from the file")
	}

	// Check that the loaded DAG contains the vertices and edges
	if _, exists := loadedDag.Vertices["A"]; !exists {
		t.Errorf("Expected vertex 'A' to exist in the loaded DAG, but it was not found")
	}
	if _, exists := loadedDag.Vertices["B"]; !exists {
		t.Errorf("Expected vertex 'B' to exist in the loaded DAG, but it was not found")
	}
	if _, exists := loadedDag.Vertices["C"]; !exists {
		t.Errorf("Expected vertex 'C' to exist in the loaded DAG, but it was not found")
	}

	// Check that the loaded DAG contains the edges
	if _, exists := loadedDag.Vertices["A"].Children["B"]; !exists {
		t.Errorf("Expected vertex 'B' to be a child of vertex 'A' in the loaded DAG")
	}
	if _, exists := loadedDag.Vertices["B"].Children["C"]; !exists {
		t.Errorf("Expected vertex 'C' to be a child of vertex 'B' in the loaded DAG")
	}
}
