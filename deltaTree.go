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
	Seq []*Delta
	IdTree  map[string]*Delta
	Pointer int
}

func NewDeltaTree() *DeltaTree {
	return &DeltaTree{
		Seq: []*Delta{},
		IdTree:  make(map[string]*Delta),
		Pointer: 0,
	}
}

func (dt *DeltaTree) Push(delta Delta) error {
	// deltaSeq := delta.GetSeq()
	// _, ok := dt.SeqTree[deltaSeq]
	// if ok {
	// 	return errors.New("Existing delta found")
	// }
	//
	// dt.Pointer = delta.GetSeq()
	// dt.SeqTree[dt.Pointer] = &delta
	dt.Pointer += 1
	dt.Seq = append(dt.Seq, &delta)

	deltaId := delta.GetId()
	dt.IdTree[deltaId] = &delta

	return nil
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

func (dt *DeltaTree) addEdgeEvent(parent *Vertex, child *Vertex) error {
	parentHash := sha1.Sum([]byte(parent.ID))
	parentStr := fmt.Sprintf("%x", parentHash[:])

	childHash := sha1.Sum([]byte(child.ID))
	childStr := fmt.Sprintf("%x", childHash[:])

	opHash := sha1.Sum([]byte{addEdge})
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
		Id:          uuidStr,
		Operation:   addEdge,
		Parent:      parent,
		Child:       child,
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
		Id:          uuidStr,
		Operation:   removeEdge,
		Parent:      parent,
		Child:       child,
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
		Id:          uuidStr,
		Operation:   addVertex,
		Vertex:      vertex,
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
		Id:          uuidStr,
		Operation:   removeVertex,
		Vertex:      vertex,
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
	primaryTreeLatestDelta := *primaryTree.Seq[primaryTree.Pointer]
	primaryDeltaId := primaryTreeLatestDelta.GetId()

	secondaryTreeLatestDelta := *secondaryTree.Seq[secondaryTree.Pointer]
	secondaryDeltaId := secondaryTreeLatestDelta.GetId()

	// Early return for identical trees
	isIdentical := primaryDeltaId == secondaryDeltaId
	if isIdentical {
		fmt.Println("Both trees are identical")
		return nil
	}

	lastCommonHash, lcsErr := calculateLcs(primaryTree, secondaryTree)
	if lcsErr != nil {
		return lcsErr
	}
	fmt.Println("LCS: ", lastCommonHash)

	hasDeviation, deviationError := checkForDeviation(primaryTree, secondaryTree, lastCommonHash)
	if deviationError != nil {
		return deviationError
	}

	// No conflicts, return all deltas required to bring both trees to the lastest of 2 trees
	if !hasDeviation {
		mergeErr := fastForward(primaryTree, secondaryTree, lastCommonHash)
		if mergeErr != nil {
			return errors.New("Error fast forwarding tree")
		}

	} else {
		conflictingErr := manualMerge(primaryTree, secondaryTree, lastCommonHash)
		if conflictingErr != nil {
			return errors.New("Error performing manual merge")
		}
	}
	return nil
}

// Primary should be lcoal
// Secondary should be remote
func calculateLcs(primaryTree *DeltaTree, secondaryTree *DeltaTree) (string, error) {
	for i := primaryTree.Pointer; i >= 0; i-- {
		primaryDelta := *primaryTree.Seq[i]
		primaryId := primaryDelta.GetId()

		_, ok := secondaryTree.IdTree[primaryId]
		if ok {
			if (i < primaryTree.Pointer) && i == secondaryTree.Pointer {
				fmt.Println("Local is ahead of Remote", primaryId)
				return primaryId, nil
			} else if secondaryTree.Pointer > i {
				fmt.Println("Remote is ahead of Local", primaryId)
				return primaryId, nil
			} else {
				fmt.Println("Both are the same", i, primaryId)

				return primaryId, nil
			}
		}
	}

	return "", nil
}

func getPositionFromHash(hash string, tree *DeltaTree) (int, error) {
	for key, value := range tree.Seq {
		delta := *value;
		id := delta.GetId();;
		if hash == id {
			return key, nil
		}

	}
	return 0, errors.New("No delta found")
}

func checkForDeviation(primaryTree *DeltaTree, secondaryTree *DeltaTree, lastCommonHash string) (bool, error) {
	primaryLatestSeq := primaryTree.Pointer
	commonPrimarySeq, priErr := getPositionFromHash(lastCommonHash, primaryTree)
	if priErr != nil {
		return false, priErr
	}

	secondaryLatestSeq := secondaryTree.Pointer
	commonSecondarySeq, secErr := getPositionFromHash(lastCommonHash, secondaryTree)
	if secErr != nil {
		return false, secErr
	}

	if (primaryLatestSeq > commonPrimarySeq) && (secondaryLatestSeq > commonSecondarySeq) {
		return true, nil
	}

	return false, nil
}

func fastForward(primaryTree *DeltaTree, secondaryTree *DeltaTree, lastCommonSeq string) error {
	fmt.Println("Can fast forward")
	return nil
}

func manualMerge(primaryTree *DeltaTree, secondaryTree *DeltaTree, lastCommonHash string) error {
	fmt.Println("Need to manually merge")

	var primaryDeviatedDeltas []*Delta
	var secondaryDeviatedDeltas []*Delta

	primaryLastCommonSeq, priErr := getPositionFromHash(lastCommonHash, primaryTree)
	if priErr != nil {
		return priErr
	}

	secondaryLastCommonSeq, secErr := getPositionFromHash(lastCommonHash, secondaryTree)
	if secErr != nil {
		return secErr
	}

	for i := primaryLastCommonSeq + 1; i <= primaryTree.Pointer; i++ {
		primaryDeviatedDeltas = append(primaryDeviatedDeltas, primaryTree.Seq[i])
	}

	for i := secondaryLastCommonSeq + 1; i <= secondaryTree.Pointer; i++ {
		secondaryDeviatedDeltas = append(secondaryDeviatedDeltas, secondaryTree.Seq[i])
	}

	primaryTreeState, priErr := squashIntoState(primaryDeviatedDeltas)
	if priErr != nil {
		return priErr
	}
	secondaryTreeState, secErr := squashIntoState(secondaryDeviatedDeltas)
	if secErr != nil {
		return secErr
	}
	fmt.Print(primaryTreeState)
	fmt.Print(secondaryTreeState)

	// deltasToApply, deconflictErr := deconflictStates(primaryTreeState, secondaryTreeState)
	// if deconflictErr != nil {
	// 	return errors.New("Error deconflicting states")
	// }

	fmt.Println("Merge successful")
	return nil
}

func deconflictStates(primaryState *State, secondaryState *State) ([]*Delta, error) {

	return nil, nil
}

type State struct {
	Vertexes map[string]*Delta
	Edges    map[string]*Delta
}

func (s *State) String() string {
	var sb strings.Builder

	// Print Vertexes
	sb.WriteString("Vertexes:\n")
	for id, deltaPtr := range s.Vertexes {
		delta := *deltaPtr
		sb.WriteString(fmt.Sprintf("  %s: %s\n", id, delta))
	}

	// Print Edges
	sb.WriteString("Edges:\n")
	for id, deltaPtr := range s.Edges {
		delta := *deltaPtr
		sb.WriteString(fmt.Sprintf("  %s: %s\n", id, delta))
	}

	return sb.String()
}

func squashIntoState(deltas []*Delta) (*State, error) {
	state := &State{
		Vertexes: make(map[string]*Delta),
		Edges:    make(map[string]*Delta),
	}
	for _, deltaPtr := range deltas {
		delta := *deltaPtr
		op := delta.GetOp()
		gid := delta.GetGid()

		switch op {
		case addVertex:
			fmt.Println("addVertex")
			if existingDeltaPtr, ok := state.Vertexes[gid]; !ok {
				state.Vertexes[gid] = &delta
			} else {
				// Merge
				existingDelta := *existingDeltaPtr
				existingOp := existingDelta.GetOp()
				switch existingOp {
				case addVertex:
					continue
				case removeVertex:
					delete(state.Vertexes, gid)
				}
			}
		case removeVertex:
			fmt.Println("removeVertex")
			if existingDeltaPtr, ok := state.Vertexes[gid]; !ok {
				state.Vertexes[gid] = &delta
			} else {
				// Merge
				existingDelta := *existingDeltaPtr
				existingOp := existingDelta.GetOp()
				switch existingOp {
				case addVertex:
					delete(state.Vertexes, gid)
				case removeVertex:
					continue
				}
			}
		case addEdge:
			fmt.Println("addEdge")
			if existingEdgePtr, ok := state.Edges[gid]; !ok {
				state.Edges[gid] = &delta
			} else {
				// Merge
				existingEdge := *existingEdgePtr
				existingOp := existingEdge.GetOp()
				switch existingOp {
				case addEdge:
					continue
				case removeEdge:
					delete(state.Edges, gid)
				}

			}
		case removeEdge:
			fmt.Println("removeEdge")
			if existingEdgePtr, ok := state.Edges[gid]; !ok {
				state.Edges[gid] = &delta
			} else {
				// Merge
				existingEdge := *existingEdgePtr
				existingOp := existingEdge.GetOp()
				switch existingOp {
				case addEdge:
					delete(state.Edges, gid)
				case removeEdge:
					continue
				}
			}
		}
	}
	return state, nil
}
