package main

import (
	"bufio"
	"crypto/sha1"
	"encoding/gob"
	"errors"
	"fmt"
	"os"
	"strings"
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
	invertErr := delta.InvertOp();
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
		Id:        uuidStr,
		Operation: addEdge,
		Parent:    parent,
		Child:     child,
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
		Id:        uuidStr,
		Operation: removeEdge,
		Parent:    parent,
		Child:     child,
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
		Id:        uuidStr,
		Operation: addVertex,
		Vertex:    vertex,
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
		Id:        uuidStr,
		Operation: removeVertex,
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

const (
	localAhead  byte = 1
	remoteAhead byte = 2
	deviated    byte = 3
	conflicted  byte = 4
	identical   byte = 5
)

// primaryTree applies to secondaryTree
func MergeTrees(primaryTree *DeltaTree, secondaryTree *DeltaTree, dag *Dag) error {
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

	lastCommonHash, longerTree, shorterTree, lcsType, lcsErr := calculateLcs(primaryTree, secondaryTree)
	if lcsErr != nil {
		return lcsErr
	}

	hasDeviation, deviationError := checkForDeviation(primaryTree, secondaryTree, lastCommonHash)
	if deviationError != nil {
		return deviationError
	}

	// No conflicts, return all deltas required to bring both trees to the lastest of 2 trees
	if !hasDeviation {
		mergeErr := fastForward(shorterTree, longerTree, lastCommonHash, dag, lcsType)
		if mergeErr != nil {
			return errors.New("Error fast forwarding tree")
		}

	} else {
		conflictingErr := manualMerge(primaryTree, secondaryTree, lastCommonHash, dag)
		if conflictingErr != nil {
			return errors.New("Error performing manual merge")
		}
	}
	return nil
}

// Primary should be lcoal
// Secondary should be remote
func calculateLcs(primaryTree *DeltaTree, secondaryTree *DeltaTree) (string, *DeltaTree, *DeltaTree, byte, error) {
	for i := primaryTree.Pointer; i >= 0; i-- {
		primaryDelta := *primaryTree.Seq[i]
		primaryId := primaryDelta.GetId()

		_, ok := secondaryTree.IdTree[primaryId]
		if ok {
			if (i < primaryTree.Pointer) && i == secondaryTree.Pointer {
				fmt.Println("Local is ahead of Remote", primaryId)
				return primaryId, primaryTree, secondaryTree, localAhead, nil
			} else if secondaryTree.Pointer > i {
				fmt.Println("Remote is ahead of Local", primaryId)
				return primaryId, secondaryTree, primaryTree, remoteAhead, nil
			} else {
				fmt.Println("Both are the same", i, primaryId)

				return primaryId, nil, nil, identical, nil
			}
		}
	}

	return "", nil, nil, 0, nil
}

func getPositionFromHash(hash string, tree *DeltaTree) (int, error) {
	for key, value := range tree.Seq {
		delta := *value
		id := delta.GetId()
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

func fastForward(shortTree *DeltaTree, longTree *DeltaTree, lastCommonSeq string, dag *Dag, lcsType byte) error {
	longTreeLastCommonSeq, err := getPositionFromHash(lastCommonSeq, longTree)
	if err != nil {
		fmt.Println("Could not get position")
		return err
	}

	nextSeq := longTreeLastCommonSeq + 1
	fmt.Println(nextSeq)
	if nextSeq >= len(longTree.Seq) {
		fmt.Println("Out of index")
		return errors.New("No more commits to add")
	}

	deltasToApply := longTree.Seq[nextSeq:]
	for _, deltaPtr := range deltasToApply {
		fmt.Println(deltaPtr)
		delta := *deltaPtr
		if lcsType == remoteAhead {
			err = delta.SetDag(dag, shortTree, false)
			if err != nil {
				fmt.Println(err.Error())
				return err
			}

		}

		err = delta.SetDeltaTree(shortTree)
		if err != nil {
			fmt.Println(err.Error())
			return err
		}
	}

	return nil
}

func manualMerge(primaryTree *DeltaTree, secondaryTree *DeltaTree, lastCommonHash string, dag *Dag) error {
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

	localTreeState, priErr := squashIntoState(primaryDeviatedDeltas)
	if priErr != nil {
		return priErr
	}
	remoteTreeState, secErr := squashIntoState(secondaryDeviatedDeltas)
	if secErr != nil {
		return secErr
	}
	fmt.Print(localTreeState)
	fmt.Print(remoteTreeState)

	_, deconflictErr := deconflictStates(localTreeState, remoteTreeState, primaryLastCommonSeq, dag, primaryTree, secondaryDeviatedDeltas)
	if deconflictErr != nil {
		return errors.New("Error deconflicting states")
	}

	fmt.Println("Merge successful")
	return nil
}

// YesNoPrompt asks yes/no questions using the label.
func YesNoPrompt(label string, def bool) bool {
	choices := "Y/n"
	if !def {
		choices = "y/N"
	}

	r := bufio.NewReader(os.Stdin)
	var s string

	for {
		fmt.Fprintf(os.Stderr, "%s (%s) ", label, choices)
		s, _ = r.ReadString('\n')
		s = strings.TrimSpace(s)
		if s == "" {
			return def
		}
		s = strings.ToLower(s)
		if s == "y" || s == "yes" {
			return true
		}
		if s == "n" || s == "no" {
			return false
		}
	}
}

func deconflictStates(primaryState *State, secondaryState *State, localLastCommonSeq int, dag *Dag, localDeltaTree *DeltaTree, remoteDeviatedDeltas []*Delta) ([]*Delta, error) {
	var localDeltasToAskUser []*Delta
	var remoteDeltasToAskUser []*Delta

	var deltaSetToApply []*Delta

	// Compare 2 states
	for primaryKey, primaryValue := range primaryState.Vertexes {
		if _, ok := secondaryState.Vertexes[primaryKey]; !ok {
				localDeltasToAskUser = append(localDeltasToAskUser, primaryValue)
		} else {
			deltaSetToApply = append(deltaSetToApply, primaryValue)
		}
	}

	for secondaryKey, secondaryValue := range secondaryState.Vertexes {
		if _, ok := primaryState.Vertexes[secondaryKey]; !ok {
				remoteDeltasToAskUser = append(remoteDeltasToAskUser, secondaryValue)
		}
	}

	var localDeltasToKeep []Delta
	var remoteDeltasToDiscard []*Delta
	
	vertexToDiscard := make(map[string]byte)

	for _, deltaValue := range localDeltasToAskUser {
		delta := *deltaValue
		gid := delta.GetGid()
		ok := YesNoPrompt(gid+"\n"+"Do you want to keep this local change?", true)
		if ok {
			localDeltasToKeep = append(localDeltasToKeep, *deltaValue)
		} else {
			vertexToDiscard[gid] = delta.GetOp()
		}
	}

	for _, deltaValue := range remoteDeltasToAskUser {
		delta := *deltaValue
		gid := delta.GetGid()
		ok := YesNoPrompt(gid+"\n"+"Do you want to keep this remote change?", true)
		if !ok {
			remoteDeltasToDiscard = append(remoteDeltasToDiscard, deltaValue)
			vertexToDiscard[gid] = delta.GetOp()
		}
	}

	var localEdgeDeltasToAskUser []*Delta
	var remoteEdgeDeltasToAskUser []*Delta

	edgesToDiscard := make(map[string]byte)

	for primaryKey, primaryValue := range primaryState.Edges {
		if _, ok := secondaryState.Edges[primaryKey]; !ok {
			delta := *primaryValue
			combinedGid := delta.GetGid()
			gids := strings.Split(combinedGid, "|")
			parentGid := gids[0]
			childGid := gids[1]
			// We just need to see if either the parent or child vertex was removed to not consider this edge anymore
			if op, ok := vertexToDiscard[parentGid]; ok {
				switch op {
				case addVertex:
					// Invalid (Cannot add if vertex does not exist)
					// Skip when invalid
				case removeVertex:
					// Valid (Vertex exists)
					localEdgeDeltasToAskUser = append(localEdgeDeltasToAskUser, primaryValue)
				}

			} else {
				if op, ok := vertexToDiscard[childGid]; ok {
					switch op {
					case addVertex:
						// Invalid (Cannot add if vertex does not exist)
						// Skip when invalid
					case removeVertex:
						// Valid (Vertex exists)
						localEdgeDeltasToAskUser = append(localEdgeDeltasToAskUser, primaryValue)
					}
				} else {
					localEdgeDeltasToAskUser = append(localEdgeDeltasToAskUser, primaryValue)
				}
			}


		}
	}

	for secondaryKey, secondaryValue := range secondaryState.Edges {
		if _, ok := primaryState.Edges[secondaryKey]; !ok {
			delta := *secondaryValue
			combinedGid := delta.GetGid()
			gids := strings.Split(combinedGid, "|")
			parentGid := gids[0]
			childGid := gids[1]
			// TODO: Refactor into function
			// We just need to see if either the parent or child vertex was removed to not consider this edge anymore
			if op, ok := vertexToDiscard[parentGid]; ok {
				switch op {
				case addVertex:
					// Invalid (Cannot add if vertex does not exist)
					// Skip when invalid
				case removeVertex:
					// Valid (Vertex exists)
					remoteEdgeDeltasToAskUser = append(remoteEdgeDeltasToAskUser, secondaryValue)
				}

			} else {
				if op, ok := vertexToDiscard[childGid]; ok {
					switch op {
					case addVertex:
						// Invalid (Cannot add if vertex does not exist)
						// Skip when invalid
					case removeVertex:
						// Valid (Vertex exists)
						remoteEdgeDeltasToAskUser = append(remoteEdgeDeltasToAskUser, secondaryValue)
					}
				} else {
					remoteEdgeDeltasToAskUser = append(remoteEdgeDeltasToAskUser, secondaryValue)
				}
			}

		}
	}

	for _, deltaValue := range localEdgeDeltasToAskUser {
		delta := *deltaValue
		gid := delta.GetGid()
		ok := YesNoPrompt(gid+"\n"+"Do you want to keep this local change?", true)
		if ok {
			localDeltasToKeep = append(localDeltasToKeep, *deltaValue)
		} else {
			edgesToDiscard[gid] = delta.GetOp()
		}
	}

	for _, deltaValue := range remoteEdgeDeltasToAskUser {
		delta := *deltaValue
		gid := delta.GetGid()
		ok := YesNoPrompt(gid+"\n"+"Do you want to keep this remote change?", true)
		if !ok {
			remoteDeltasToDiscard = append(remoteDeltasToDiscard, deltaValue)
			edgesToDiscard[gid] = delta.GetOp()
		}
	}

	// Rewind local dag and deltaTree to last common hash
	fmt.Println("Rewinding existing local delta tree")
	count := len(localDeltaTree.Seq)
	var reverseDeltas []*Delta
	for i := localLastCommonSeq; i < count - 1; i++ {
		poppedDelta, popErr:= localDeltaTree.Pop();
		if popErr != nil {
			return nil, popErr
		}
		reverseDeltas = append(reverseDeltas, poppedDelta)
	}

	fmt.Println("Updating local dag")
	for i := len(reverseDeltas) - 1; i >= 0; i-- {
		delta := *reverseDeltas[i]
		copy := delta.GetCopy()
		copy.SetDag(dag, localDeltaTree, true)
	}

	fmt.Println("Replay all deltas after commonHash for remote")
	// Replay all deltas after commonHash for remote (This will include common deltas)
	for _, deltaValue := range remoteDeviatedDeltas {
		delta := *deltaValue
		delta.SetDag(dag, localDeltaTree, false)
		delta.SetDeltaTree(localDeltaTree)
	}

	// Apply negative commits of to attain desired state for remoteDeltasToKeep
	mergeCommitDeltas, invertErr := invertDeltaOp(remoteDeltasToDiscard)
	if invertErr != nil {
		return nil, invertErr
	}

	fmt.Println("Applying negative commits to retain desired state for remote")
	for i := len(mergeCommitDeltas) - 1; i >= 0; i-- {
		delta := *mergeCommitDeltas[i]
		delta.SetDag(dag, localDeltaTree, false)
		delta.SetDeltaTree(localDeltaTree)

	}

	fmt.Println("Applying localDeltasToKeep")
	// Apply all localDeltasToKeep
	for _, delta := range localDeltasToKeep {
		delta.InvertOp()
		delta.SetDag(dag, localDeltaTree, false)
		delta.SetDeltaTree(localDeltaTree)
	}

	return nil, nil
}

func invertDeltaOp(deltas []*Delta) ([]*Delta, error) {
	var copies []*Delta

	for _, ptr := range deltas {
		delta := *ptr
		copied := delta.GetCopy()
		err := copied.InvertOp()
		if err != nil {
			return nil, err
		}

		copies = append(copies, &copied)
	}
	return copies, nil
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
