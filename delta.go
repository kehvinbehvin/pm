package main

import (
	"fmt"
)

const (
	addVertex    byte = 1
	removeVertex byte = 2
	addEdge      byte = 3
	removeEdge   byte = 4
)

type Delta interface {
	SetDag(*Dag, *DeltaTree, bool) error
	SetDeltaTree(*DeltaTree) error
	UnSet()
	GetId() string
	String() string
	GetOp() byte
	GetGid() string
	InvertOp() error
	GetCopy() Delta
}

type VertexDelta struct {
	Id        string
	Operation byte
	Vertex    *Vertex
}

func (vd *VertexDelta) GetGid() string {
	vertex := *vd.Vertex
	gid := vertex.ID
	return gid
}

func (vd *VertexDelta) GetOp() byte {
	return vd.Operation
}

func (vd *VertexDelta) GetCopy() Delta {
	originalDelta := *vd
	originalVertex := *originalDelta.Vertex
	vertexCopy := originalVertex
	copy := &VertexDelta{
		Id: originalDelta.Id,
		Operation: originalDelta.Operation,
		Vertex: &vertexCopy,
	}

	return copy
}

func (vd *VertexDelta) InvertOp() error {
	switch vd.Operation {
	case addVertex:
		vd.Operation = removeVertex
	case removeVertex:
		vd.Operation = addVertex
	}

	return nil
}

func (vd *VertexDelta) SetDag(dag *Dag, tree *DeltaTree, silent bool) error {
	switch vd.Operation {
	case addVertex:
		err := dag.addVertex(vd.Vertex)
		if err != nil {
			return err
		}
	case removeVertex:
		err := dag.removeVertex(vd.Vertex, tree, silent)
		if err != nil {
			return err
		}
	}

	return nil
}

func (vd *VertexDelta) SetDeltaTree(tree *DeltaTree) error {
	switch vd.Operation {
	case addVertex:
		err := tree.addVertexEvent(vd.Vertex)
		if err != nil {
			return nil
		}
	case removeVertex:
		err := tree.removeVertexEvent(vd.Vertex)
		if err != nil {
			return nil
		}
	}

	return nil
}

func (vd *VertexDelta) UnSet() {

}

func (vd *VertexDelta) GetId() string {
	return vd.Id
}

func (vd *VertexDelta) String() string {
	return fmt.Sprintf("VertexDelta(Id: %s, Operation: %d, Vertex: %s)", vd.Id, vd.Operation, vd.Vertex.ID)
}

type EdgeDelta struct {
	Id        string
	Operation byte
	Parent    *Vertex
	Child     *Vertex
}

func (ed *EdgeDelta) GetGid() string {
	parent := *ed.Parent
	child := *ed.Child
	gid := parent.ID + "|" + child.ID
	return gid
}

func (ed *EdgeDelta) GetOp() byte {
	return ed.Operation
}

func (ed *EdgeDelta) GetCopy() Delta {
	originalDelta := *ed
	originalParent := *originalDelta.Parent
	originalChild := *originalDelta.Child
	parentCopy := originalParent
	childCopy := originalChild

	copy := &EdgeDelta{
		Id: originalDelta.Id,
		Operation: originalDelta.Operation,
		Parent: &parentCopy,
		Child: &childCopy,
	}

	return copy
}

func (ed *EdgeDelta) InvertOp() error {
	switch ed.Operation {
	case addEdge:
		ed.Operation = removeEdge
	case removeEdge:
		ed.Operation = addEdge
	}

	return nil
}

func (ed *EdgeDelta) SetDag(dag *Dag, tree *DeltaTree, silent bool) error {
	switch ed.Operation {
	case addEdge:
		err := dag.addEdge(ed.Parent, ed.Child)
		if err != nil {
			return err
		}
	case removeEdge:
		err := dag.removeEdge(ed.Parent, ed.Child)
		if err != nil {
			return err
		}
	}

	return nil
}

func (ed *EdgeDelta) SetDeltaTree(tree *DeltaTree) error {
	switch ed.Operation {
	case addEdge:
		err := tree.addEdgeEvent(ed.Parent, ed.Child)
		if err != nil {
			return nil
		}
	case removeEdge:
		err := tree.removeEdgeEvent(ed.Parent, ed.Child)
		if err != nil {
			return nil
		}
	}

	return nil
}

func (ed *EdgeDelta) UnSet() {
}

func (ed *EdgeDelta) GetId() string {
	return ed.Id
}

func (ed *EdgeDelta) String() string {
	return fmt.Sprintf("EdgeDelta(Id: %s, Operation: %d, Parent: %s, Child: %s)", ed.Id, ed.Operation, ed.Parent.ID, ed.Child.ID)
}
