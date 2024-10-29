package dag

import (
	"crypto/sha1"
	"encoding/gob"
	"fmt"
	"os"
	"strings"
)

const (
	LocalAhead  byte = 1
	RemoteAhead byte = 2
	Deviated    byte = 3
	Conflicted  byte = 4
	Identical   byte = 5
)

type DeltaTree struct {
	Seq     []*Delta
	IdTree  map[string]*Delta
	Pointer int
}

func NewDeltaTree() *DeltaTree {
	return &DeltaTree{
		Seq:     []*Delta{},
		IdTree:  make(map[string]*Delta),
		Pointer: 0,
	}
}

func (dt *DeltaTree) Push(delta Delta) error {
	if len(dt.Seq) > 0 {
		dt.Pointer += 1
	}
	dt.Seq = append(dt.Seq, &delta)

	deltaId := delta.GetId()
	dt.IdTree[deltaId] = &delta

	return nil
}

func (dt *DeltaTree) Pop() (*Delta, error) {
	if len(dt.Seq) == 0 {
		return nil, fmt.Errorf("cannot pop from an empty DeltaTree")
	}

	// Get the last element in Seq
	lastDelta := dt.Seq[len(dt.Seq)-1]
	delta := *lastDelta

	// Remove the last element from Seq
	dt.Seq = dt.Seq[:len(dt.Seq)-1]

	// Update Pointer to the new last element (or -1 if Seq is empty)
	if len(dt.Seq) == 0 {
		dt.Pointer = 0
	} else {
		dt.Pointer = len(dt.Seq) - 1
	}

	// Remove the last delta from IdTree
	deltaId := delta.GetId()
	invertErr := delta.InvertOp()
	if invertErr != nil {
		return nil, invertErr
	}
	delete(dt.IdTree, deltaId)

	return &delta, nil
}

func (dt *DeltaTree) String() string {
	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("DeltaTree (Pointer: %d)\n", dt.Pointer))
	builder.WriteString("Deltas:\n")
	//
	// for id, delta := range dt.SeqTree {
	// 	builder.WriteString(fmt.Sprintf("Seq: %d, Delta: %s\n", id, *delta))
	// }

	for id, delta := range dt.IdTree {
		builder.WriteString(fmt.Sprintf("ID: %s, Delta: %s\n", id, *delta))
	}

	return builder.String()
}

func (dt *DeltaTree) AddEdgeEvent(parent *Vertex, child *Vertex) error {
	parentHash := sha1.Sum([]byte(parent.ID))
	parentStr := fmt.Sprintf("%x", parentHash[:])

	childHash := sha1.Sum([]byte(child.ID))
	childStr := fmt.Sprintf("%x", childHash[:])

	opHash := sha1.Sum([]byte{AddEdge})
	opStr := fmt.Sprintf("%x", opHash[:])

	var prevHash string
	deltaTreeLength := len(dt.Seq)

	if deltaTreeLength == 0 {
		prevHash = ""
	} else {
		parentDeltaPtr := dt.Seq[dt.Pointer]
		parentDelta := *parentDeltaPtr
		prevHash = parentDelta.GetId()
	}

	uuid := parentStr + childStr + opStr + prevHash
	uuidHash := sha1.Sum([]byte(uuid))
	uuidStr := fmt.Sprintf("%x", uuidHash[:])

	delta := &EdgeDelta{
		Id:        uuidStr,
		Operation: AddEdge,
		Parent:    parent,
		Child:     child,
	}

	dt.Push(delta)

	return nil
}

func (dt *DeltaTree) RemoveEdgeEvent(parent *Vertex, child *Vertex) error {
	parentHash := sha1.Sum([]byte(parent.ID))
	parentStr := fmt.Sprintf("%x", parentHash[:])

	childHash := sha1.Sum([]byte(child.ID))
	childStr := fmt.Sprintf("%x", childHash[:])

	opHash := sha1.Sum([]byte{RemoveEdge})
	opStr := fmt.Sprintf("%x", opHash[:])

	var prevHash string
	deltaTreeLength := len(dt.Seq)

	if deltaTreeLength == 0 {
		prevHash = ""
	} else {
		parentDeltaPtr := dt.Seq[dt.Pointer]
		parentDelta := *parentDeltaPtr
		prevHash = parentDelta.GetId()
	}

	uuid := parentStr + childStr + opStr + prevHash
	uuidHash := sha1.Sum([]byte(uuid))
	uuidStr := fmt.Sprintf("%x", uuidHash[:])

	delta := &EdgeDelta{
		Id:        uuidStr,
		Operation: RemoveEdge,
		Parent:    parent,
		Child:     child,
	}

	dt.Push(delta)

	return nil
}

func (dt *DeltaTree) AddVertexEvent(vertex *Vertex) error {
	idHash := sha1.Sum([]byte(vertex.ID))
	hashStr := fmt.Sprintf("%x", idHash[:])

	opHash := sha1.Sum([]byte{AddVertex})
	opStr := fmt.Sprintf("%x", opHash[:])

	var prevHash string

	deltaTreeLength := len(dt.Seq)
	if deltaTreeLength == 0 {
		prevHash = ""
	} else {
		parentDeltaPtr := dt.Seq[dt.Pointer]
		parentDelta := *parentDeltaPtr
		prevHash = parentDelta.GetId()
	}

	uuid := hashStr + opStr + prevHash
	uuidHash := sha1.Sum([]byte(uuid))
	uuidStr := fmt.Sprintf("%x", uuidHash[:])

	delta := &VertexDelta{
		Id:        uuidStr,
		Operation: AddVertex,
		Vertex:    vertex,
	}

	dt.Push(delta)

	return nil
}

func (dt *DeltaTree) RemoveVertexEvent(vertex *Vertex) error {
	idHash := sha1.Sum([]byte(vertex.ID))
	hashStr := fmt.Sprintf("%x", idHash[:])

	opHash := sha1.Sum([]byte{AddVertex})
	opStr := fmt.Sprintf("%x", opHash[:])

	var prevHash string

	deltaTreeLength := len(dt.Seq)
	if deltaTreeLength == 0 {
		prevHash = ""
	} else {
		parentDeltaPtr := dt.Seq[dt.Pointer]
		parentDelta := *parentDeltaPtr
		prevHash = parentDelta.GetId()
	}

	uuid := hashStr + opStr + prevHash
	uuidHash := sha1.Sum([]byte(uuid))
	uuidStr := fmt.Sprintf("%x", uuidHash[:])

	delta := &VertexDelta{
		Id:        uuidStr,
		Operation: RemoveVertex,
		Vertex:    vertex,
	}

	dt.Push(delta)

	return nil
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

func (dt *DeltaTree) SaveDelta(fileName string) {
	file, err := os.Create(fileName)
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

func LoadRemoteDelta() *DeltaTree {
	file, fileErr := os.Open("./.pm/remote/delta")

	if fileErr != nil {
		fmt.Println("Error opening remote delta file")
		return nil
	}
	defer file.Close()

	gob.Register(&VertexDelta{})
	gob.Register(&EdgeDelta{})

	decoder := gob.NewDecoder(file)

	var deltaTree *DeltaTree
	decodingErr := decoder.Decode(&deltaTree)
	if decodingErr != nil {
		fmt.Println("Error decoding remote delta tree")
		return nil
	}

	return deltaTree
}
