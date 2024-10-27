package main

import (
	"crypto/sha1"
	"encoding/gob"
	"errors"
	"fmt"
	"os"
	"strings"
)

type DeltaTree struct {
	SeqTree map[int]*Delta
	IdTree  map[string]*Delta
	Pointer int
}

func NewDeltaTree() *DeltaTree {
	return &DeltaTree{
		SeqTree: make(map[int]*Delta),
		IdTree:  make(map[string]*Delta),
		Pointer: 0,
	}
}

// Not in use now
func (dt *DeltaTree) Checkout(nodeHash int) error {
	_, ok := dt.SeqTree[nodeHash]
	if !ok {
		return errors.New("No hash found")
	}

	return nil
}

func (dt *DeltaTree) Push(delta Delta) error {
	deltaId := delta.GetSeq()
	_, ok := dt.SeqTree[deltaId]
	if ok {
		return errors.New("Existing delta found")
	}

	dt.Pointer = delta.GetSeq()
	dt.SeqTree[dt.Pointer] = &delta
	return nil
}

func (dt *DeltaTree) String() string {
	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("DeltaTree (Pointer: %d)\n", dt.Pointer))
	builder.WriteString("Deltas:\n")

	for id, delta := range dt.SeqTree {
		builder.WriteString(fmt.Sprintf("ID: %d, Delta: %s\n", id, *delta))
	}

	return builder.String()
}

func (dt *DeltaTree) addEdgeEvent(parent *Vertex, child *Vertex) error {
	parentHash := sha1.Sum([]byte(parent.ID))
	parentStr := fmt.Sprintf("%x", parentHash[:])

	childHash := sha1.Sum([]byte(child.ID))
	childStr := fmt.Sprintf("%x", childHash[:])

	opHash := sha1.Sum([]byte{addEdge})
	opStr := fmt.Sprintf("%x", opHash[:])

	var prevHash string

	deltaParentSeq := dt.Pointer
	if deltaParentSeq == 0 {
		prevHash = ""
	} else {
		deltaParentPtr, ok := dt.SeqTree[deltaParentSeq]
		if !ok {
			return errors.New("Parent seq not found")
		}

		deltaParent := *deltaParentPtr
		prevHash = deltaParent.GetId()
	}

	uuid := parentStr + childStr + opStr + prevHash
	uuidHash := sha1.Sum([]byte(uuid))
	uuidStr := fmt.Sprintf("%x", uuidHash[:])

	delta := &EdgeDelta{
		Id:          uuidStr,
		Operation:   addEdge,
		Parent:      parent,
		Child:       child,
		ParentDelta: deltaParentSeq,
		Seq:         deltaParentSeq + 1,
	}

	dt.Push(delta)

	return nil
}

func (dt *DeltaTree) removeEdgeEvent(parent *Vertex, child *Vertex) error {
	parentHash := sha1.Sum([]byte(parent.ID))
	parentStr := fmt.Sprintf("%x", parentHash[:])

	childHash := sha1.Sum([]byte(child.ID))
	childStr := fmt.Sprintf("%x", childHash[:])

	opHash := sha1.Sum([]byte{removeEdge})
	opStr := fmt.Sprintf("%x", opHash[:])

	var prevHash string

	deltaParentSeq := dt.Pointer
	if deltaParentSeq == 0 {
		prevHash = ""
	} else {
		deltaParentPtr, ok := dt.SeqTree[deltaParentSeq]
		if !ok {
			return errors.New("Parent seq not found")
		}

		deltaParent := *deltaParentPtr
		prevHash = deltaParent.GetId()
	}

	uuid := parentStr + childStr + opStr + prevHash
	uuidHash := sha1.Sum([]byte(uuid))
	uuidStr := fmt.Sprintf("%x", uuidHash[:])

	delta := &EdgeDelta{
		Id:          uuidStr,
		Operation:   removeEdge,
		Parent:      parent,
		Child:       child,
		ParentDelta: deltaParentSeq,
		Seq:         deltaParentSeq + 1,
	}

	dt.Push(delta)

	return nil
}

func (dt *DeltaTree) addVertexEvent(vertex *Vertex) error {
	idHash := sha1.Sum([]byte(vertex.ID))
	hashStr := fmt.Sprintf("%x", idHash[:])

	opHash := sha1.Sum([]byte{addVertex})
	opStr := fmt.Sprintf("%x", opHash[:])

	var prevHash string

	deltaParentSeq := dt.Pointer
	if deltaParentSeq == 0 {
		prevHash = ""
	} else {
		deltaParentPtr, ok := dt.SeqTree[deltaParentSeq]
		if !ok {
			return errors.New("Parent seq not found")
		}

		deltaParent := *deltaParentPtr
		prevHash = deltaParent.GetId()
	}

	uuid := hashStr + opStr + prevHash
	uuidHash := sha1.Sum([]byte(uuid))
	uuidStr := fmt.Sprintf("%x", uuidHash[:])

	delta := &VertexDelta{
		Id:          uuidStr,
		Operation:   addVertex,
		Vertex:      vertex,
		ParentDelta: deltaParentSeq,
		Seq:         deltaParentSeq + 1,
	}

	dt.Push(delta)

	return nil
}

func (dt *DeltaTree) removeVertexEvent(vertex *Vertex) error {
	idHash := sha1.Sum([]byte(vertex.ID))
	hashStr := fmt.Sprintf("%x", idHash[:])

	opHash := sha1.Sum([]byte{addVertex})
	opStr := fmt.Sprintf("%x", opHash[:])

	var prevHash string

	deltaParentSeq := dt.Pointer
	if deltaParentSeq == 0 {
		prevHash = ""
	} else {
		deltaParentPtr, ok := dt.SeqTree[deltaParentSeq]
		if !ok {
			return errors.New("Parent seq not found")
		}

		deltaParent := *deltaParentPtr
		prevHash = deltaParent.GetId()
	}

	uuid := hashStr + opStr + prevHash
	uuidHash := sha1.Sum([]byte(uuid))
	uuidStr := fmt.Sprintf("%x", uuidHash[:])

	delta := &VertexDelta{
		Id:          uuidStr,
		Operation:   removeVertex,
		Vertex:      vertex,
		ParentDelta: deltaParentSeq,
		Seq:         deltaParentSeq + 1,
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

func (dt *DeltaTree) SaveDelta() {
	file, err := os.Create("./.pm/delta")
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

type DeltaPacket struct {
	DeltaOperations []*Delta
	Tree            *DeltaTree
}

type ConflictingDeltaPacket struct {
	Primary   *Delta
	Secondary *Delta
	Selection *Delta
	Tree      *DeltaTree
}

// primaryTree applies to secondaryTree
func MergeTrees(primaryTree *DeltaTree, secondaryTree *DeltaTree) error {
	primaryTreeLatestDelta := *primaryTree.SeqTree[primaryTree.Pointer]
	primaryDeltaId := primaryTreeLatestDelta.GetId()

	secondaryTreeLatestDelta := *secondaryTree.SeqTree[secondaryTree.Pointer]
	secondaryDeltaId := secondaryTreeLatestDelta.GetId()

	// Early return for identical trees
	isIdentical := primaryDeltaId == secondaryDeltaId
	if isIdentical {
		return nil
	}

	lastCommonSeq, lcsErr := calculateLcs(primaryTree, secondaryTree)
	if lcsErr != nil {
		return lcsErr
	}
	hasDeviation, deviationError := checkForDeviation(primaryTree, secondaryTree, lastCommonSeq)
	if deviationError != nil {
		return deviationError
	}

	// No conflicts, return all deltas required to bring both trees to the lastest of 2 trees
	if !hasDeviation {
		mergeErr := fastForward(primaryTree, secondaryTree, lastCommonSeq)
		if mergeErr != nil {
			return errors.New("Error fast forwarding tree")
		}

	} else {
		conflictingErr := manualMerge(primaryTree, secondaryTree, lastCommonSeq)
		if conflictingErr != nil {
			return errors.New("Error performing manual merge")
		}
	}
	return nil
}

func calculateLcs(primaryTree *DeltaTree, secondaryTree *DeltaTree) (int, error) {
	return 0, nil
}

func checkForDeviation(primaryTree *DeltaTree, secondaryTree *DeltaTree, lastCommonSeq int) (bool, error) {
	return true, nil
}

func fastForward(primaryTree *DeltaTree, secondaryTree *DeltaTree, lastCommonSeq int) error {
	return nil
}

func manualMerge(primaryTree *DeltaTree, secondaryTree *DeltaTree, lastCommonSeq int) error {
	return nil
}
