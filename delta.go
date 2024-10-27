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
	Set()
	UnSet()
	GetId() string
	String() string
	GetParent(*DeltaTree) *Delta
	GetSeq() int
	GetOp() byte
	GetGid() string
}

type VertexDelta struct {
	Id          string
	Operation   byte
	Vertex      *Vertex
	ParentDelta int
	Seq         int
}

func (vd *VertexDelta) GetGid() string {
	vertex := *vd.Vertex
	gid := vertex.ID
	return gid
}

func (vd *VertexDelta) GetOp() byte {
	return vd.Operation
}

func (vd *VertexDelta) Set() {

}

func (vd *VertexDelta) UnSet() {

}

func (vd *VertexDelta) GetId() string {
	return vd.Id
}

func (vd *VertexDelta) GetSeq() int {
	return vd.Seq
}

func (vd *VertexDelta) GetParent(dt *DeltaTree) *Delta {
	parentId := vd.ParentDelta
	delta, ok := dt.SeqTree[parentId]
	if !ok {
		fmt.Println("Cannot find parent")
		return nil
	}

	return delta
}

func (vd *VertexDelta) String() string {
	return fmt.Sprintf("VertexDelta(Id: %s, Operation: %d, Vertex: %s, ParentDelta: %d)", vd.Id, vd.Operation, vd.Vertex.ID, vd.ParentDelta)
}

type EdgeDelta struct {
	Id          string
	Operation   byte
	Parent      *Vertex
	Child       *Vertex
	ParentDelta int
	Seq         int
}

func (ed *EdgeDelta) GetGid() string {
	parent := *ed.Parent
	child := *ed.Child
	gid := parent.ID + child.ID
	return gid
}

func (ed *EdgeDelta) GetOp() byte {
	return ed.Operation
}

func (ed *EdgeDelta) Set() {
}

func (ed *EdgeDelta) UnSet() {
}

func (ed *EdgeDelta) GetId() string {
	return ed.Id
}

func (ed *EdgeDelta) GetSeq() int {
	return ed.Seq
}

func (ed *EdgeDelta) GetParent(dt *DeltaTree) *Delta {
	parentId := ed.ParentDelta
	delta, ok := dt.SeqTree[parentId]
	if !ok {
		fmt.Println("Cannot find parent")
		return nil
	}

	return delta
}

func (ed *EdgeDelta) String() string {
	return fmt.Sprintf("EdgeDelta(Id: %s, Operation: %d, Parent: %s, Child: %s, ParentDelta: %d)", ed.Id, ed.Operation, ed.Parent.ID, ed.Child.ID, ed.ParentDelta)
}
