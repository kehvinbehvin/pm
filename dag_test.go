package main

import (
	"testing"
)

// Test creating a new DAG
func TestCreateDag(t *testing.T) {
	dag := newDag()
	if dag == nil {
		t.Errorf("Expected DAG to be created, but got nil")
	}
}

// Test creating a new Vertex
func TestCreateVertex(t *testing.T) {
	v := newVertex("A")
	if v.id != "A" {
		t.Errorf("Expected vertex ID to be 'A', but got '%s'", v.id)
	}
	if len(v.children) != 0 || len(v.parents) != 0 {
		t.Errorf("Expected vertex to have no children or parents")
	}
}

// Test adding a vertex to the DAG
func TestAddVertex(t *testing.T) {
	dag := newDag()
	v := newVertex("A")
	dag.addVertex(v)

	if _, exists := dag.vertices["A"]; !exists {
		t.Errorf("Expected vertex 'A' to be added to the DAG")
	}
}

// Test trying to add the same vertex twice
func TestAddSameVertexTwice(t *testing.T) {
	dag := newDag()
	v := newVertex("A")
	dag.addVertex(v)
	dag.addVertex(v) // Attempting to add the same vertex

	if len(dag.vertices) != 1 {
		t.Errorf("Expected only one vertex in the DAG, but found %d", len(dag.vertices))
	}
}

// Test removing a vertex from the DAG
func TestRemoveVertex(t *testing.T) {
	dag := newDag()
	v := newVertex("A")
	dag.addVertex(v)

	dag.removeVertex(v)

	if _, exists := dag.vertices["A"]; exists {
		t.Errorf("Expected vertex 'A' to be removed from the DAG")
	}
}

// Test removing a non-existent vertex
func TestRemoveNonExistentVertex(t *testing.T) {
	dag := newDag()
	v := newVertex("A")
	dag.removeVertex(v)

	if len(dag.vertices) != 0 {
		t.Errorf("Expected DAG to have no vertices, but found %d", len(dag.vertices))
	}
}

// Test adding an edge between two vertices
func TestAddEdge(t *testing.T) {
	dag := newDag()
	v1 := newVertex("A")
	v2 := newVertex("B")

	dag.addVertex(v1)
	dag.addVertex(v2)

	dag.addEdge(v1, v2)

	if _, exists := v1.children[v2.id]; !exists {
		t.Errorf("Expected vertex 'B' to be a child of vertex 'A'")
	}

	if _, exists := v2.parents[v1.id]; !exists {
		t.Errorf("Expected vertex 'A' to be a parent of vertex 'B'")
	}
}

// Test adding an edge that would create a cycle
func TestAddEdgeWithCycle(t *testing.T) {
	dag := newDag()
	v1 := newVertex("A")
	v2 := newVertex("B")

	dag.addVertex(v1)
	dag.addVertex(v2)

	dag.addEdge(v1, v2)

	// Try to add an edge from v2 to v1, which would create a cycle
	dag.addEdge(v2, v1)

	if _, exists := v2.children[v1.id]; exists {
		t.Errorf("Expected DAG to prevent cycle when adding an edge from 'B' to 'A'")
	}
}

// Test adding an edge for non-existent vertices
func TestAddEdgeForNonExistentVertex(t *testing.T) {
	dag := newDag()
	v1 := newVertex("A")
	v2 := newVertex("B")

	dag.addVertex(v1)

	// Try to add an edge where vertex B doesn't exist in the DAG
	dag.addEdge(v1, v2)

	if _, exists := v1.children[v2.id]; exists {
		t.Errorf("Expected no edge to be added because vertex 'B' does not exist")
	}
}

// Test removing an edge between two vertices
func TestRemoveEdge(t *testing.T) {
	dag := newDag()
	v1 := newVertex("A")
	v2 := newVertex("B")

	dag.addVertex(v1)
	dag.addVertex(v2)
	dag.addEdge(v1, v2)

	dag.removeEdge(v1, v2)

	if _, exists := v1.children[v2.id]; exists {
		t.Errorf("Expected edge between 'A' and 'B' to be removed")
	}

	if _, exists := v2.parents[v1.id]; exists {
		t.Errorf("Expected edge between 'B' and 'A' to be removed")
	}
}

// Test removing a non-existent edge
func TestRemoveNonExistentEdge(t *testing.T) {
	dag := newDag()
	v1 := newVertex("A")
	v2 := newVertex("B")

	dag.addVertex(v1)
	dag.addVertex(v2)

	// Try to remove an edge that does not exist
	dag.removeEdge(v1, v2)

	if len(v1.children) != 0 || len(v2.parents) != 0 {
		t.Errorf("Expected no edges to exist between 'A' and 'B'")
	}
}

// Test DFS to check if a path exists between two vertices
func TestDfs(t *testing.T) {
	dag := newDag()
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
	dag := newDag()
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
