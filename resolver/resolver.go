package resolver

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	"github/pm/dag"
)

// primaryTree applies to secondaryTree
func MergeTrees(primaryTree *dag.DeltaTree, secondaryTree *dag.DeltaTree, dag *dag.Dag) error {
	isPrimaryTreeEmpty := len(primaryTree.Seq) == 0
	isSecondaryTreeEmpty := len(secondaryTree.Seq) == 0
	if isPrimaryTreeEmpty && !isSecondaryTreeEmpty {
		deltasToApply := secondaryTree.Seq
		applyErr := applyDeltasToLocal(deltasToApply, primaryTree, dag)
		if applyErr != nil {
			fmt.Println("Error when applying remote to local when local is empty")
			return applyErr
		}
		return nil
	}

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

	lastCommonHash, longerTree, shorterTree, lcsType, lcsErr := CalculateLcs(primaryTree, secondaryTree)
	if lcsErr != nil {
		return lcsErr
	}

	hasDeviation, deviationError := CheckForDeviation(primaryTree, secondaryTree, lastCommonHash)
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
func CalculateLcs(primaryTree *dag.DeltaTree, secondaryTree *dag.DeltaTree) (string, *dag.DeltaTree, *dag.DeltaTree, byte, error) {
	if len(primaryTree.Seq) > 0 && len(secondaryTree.Seq) == 0 {
		primaryDelta := *primaryTree.Seq[0]
		primaryId := primaryDelta.GetId()
		return primaryId, primaryTree, secondaryTree, dag.LocalAhead, nil
	}

	for i := primaryTree.Pointer; i >= 0; i-- {
		primaryDelta := *primaryTree.Seq[i]
		primaryId := primaryDelta.GetId()

		_, ok := secondaryTree.IdTree[primaryId]
		if ok {
			if (i < primaryTree.Pointer) && i == secondaryTree.Pointer {
				fmt.Println("Local is ahead of Remote", primaryId)
				return primaryId, primaryTree, secondaryTree, dag.LocalAhead, nil
			} else if secondaryTree.Pointer > i {
				fmt.Println("Remote is ahead of Local", primaryId)
				return primaryId, secondaryTree, primaryTree, dag.RemoteAhead, nil
			} else {
				fmt.Println("Both are the same", i, primaryId)

				return primaryId, nil, nil, dag.Identical, nil
			}
		}
	}

	return "", nil, nil, 0, nil
}

func getPositionFromHash(hash string, tree *dag.DeltaTree) (int, error) {
	for key, value := range tree.Seq {
		delta := *value
		id := delta.GetId()
		if hash == id {
			return key, nil
		}

	}
	return 0, errors.New("No delta found")
}

func CheckForDeviation(primaryTree *dag.DeltaTree, secondaryTree *dag.DeltaTree, lastCommonHash string) (bool, error) {
	primaryLatestSeq := primaryTree.Pointer
	commonPrimarySeq, priErr := getPositionFromHash(lastCommonHash, primaryTree)
	if priErr != nil {
		return true, priErr
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

func GetDeltasAhead(longTree *dag.DeltaTree, lastCommonSeq string) ([]*dag.Delta, error) {
	longTreeLastCommonSeq, err := getPositionFromHash(lastCommonSeq, longTree)
	if err != nil {
		fmt.Println("Could not get position")
		return nil, err
	}

	nextSeq := longTreeLastCommonSeq + 1
	if nextSeq >= len(longTree.Seq) {
		fmt.Println("Out of index")
		return nil, errors.New("No more commits to add")
	}

	deltasToApply := longTree.Seq[nextSeq:]
	return deltasToApply, nil
}

func applyDeltasToLocal(deltas []*dag.Delta, tree *dag.DeltaTree, dag *dag.Dag) (error) {
	for _, deltaPtr := range deltas {
		delta := *deltaPtr
		// Only fastforward the local tree.
		// If remoteAhead, then local wil be the shortTree
		err := delta.SetDag(dag, tree, false)
		if err != nil {
			fmt.Println(err.Error())
			return err
		}

		err = delta.SetDeltaTree(tree)
		if err != nil {
			fmt.Println(err.Error())
			return err
		}
	}
	
	return nil
}

func fastForward(shortTree *dag.DeltaTree, longTree *dag.DeltaTree, lastCommonSeq string, localDag *dag.Dag, lcsType byte) error {
	deltasToApply, deltaErr := GetDeltasAhead(longTree, lastCommonSeq)
	if deltaErr != nil {
		return deltaErr
	}

	if lcsType == dag.RemoteAhead {
		applyErr := applyDeltasToLocal(deltasToApply, shortTree, localDag)
		if applyErr != nil {
			fmt.Println("Error when apply deltas to local tree and dag")
			return applyErr
		}
	}
	
	return nil
}

func manualMerge(primaryTree *dag.DeltaTree, secondaryTree *dag.DeltaTree, lastCommonHash string, localDag *dag.Dag) error {
	var primaryDeviatedDeltas []*dag.Delta
	var secondaryDeviatedDeltas []*dag.Delta

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
	fmt.Println("Local State")
	fmt.Print(localTreeState)
	fmt.Println("Remote State")
	fmt.Print(remoteTreeState)

	_, deconflictErr := deconflictStates(localTreeState, remoteTreeState, primaryLastCommonSeq, localDag, primaryTree, secondaryDeviatedDeltas)
	if deconflictErr != nil {
		return errors.New("Error deconflicting states")
	}

	return nil
}

func deconflictStates(primaryState *State, secondaryState *State, localLastCommonSeq int, localDag *dag.Dag, localDeltaTree *dag.DeltaTree, remoteDeviatedDeltas []*dag.Delta) ([]*dag.Delta, error) {
	var localDeltasToAskUser []*dag.Delta
	var remoteDeltasToAskUser []*dag.Delta

	var deltaSetToApply []*dag.Delta

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

	var localDeltasToKeep []dag.Delta
	var remoteDeltasToDiscard []*dag.Delta

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

	var localEdgeDeltasToAskUser []*dag.Delta
	var remoteEdgeDeltasToAskUser []*dag.Delta

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
				case dag.AddVertex:
					// Invalid (Cannot add if vertex does not exist)
					// Skip when invalid
				case dag.RemoveVertex:
					// Valid (Vertex exists)
					localEdgeDeltasToAskUser = append(localEdgeDeltasToAskUser, primaryValue)
				}

			} else {
				if op, ok := vertexToDiscard[childGid]; ok {
					switch op {
					case dag.AddVertex:
						// Invalid (Cannot add if vertex does not exist)
						// Skip when invalid
					case dag.RemoveVertex:
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
				case dag.AddVertex:
					// Invalid (Cannot add if vertex does not exist)
					// Skip when invalid
				case dag.RemoveVertex:
					// Valid (Vertex exists)
					remoteEdgeDeltasToAskUser = append(remoteEdgeDeltasToAskUser, secondaryValue)
				}

			} else {
				if op, ok := vertexToDiscard[childGid]; ok {
					switch op {
					case dag.AddVertex:
						// Invalid (Cannot add if vertex does not exist)
						// Skip when invalid
					case dag.RemoveVertex:
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
	var reverseDeltas []*dag.Delta
	for i := localLastCommonSeq; i < count-1; i++ {
		poppedDelta, popErr := localDeltaTree.Pop()
		if popErr != nil {
			return nil, popErr
		}
		reverseDeltas = append(reverseDeltas, poppedDelta)
	}

	fmt.Println("Updating local dag")
	for i := len(reverseDeltas) - 1; i >= 0; i-- {
		delta := *reverseDeltas[i]
		copy := delta.GetCopy()
		copy.SetDag(localDag, localDeltaTree, true)
	}

	fmt.Println("Replay all deltas after commonHash for remote")
	// Replay all deltas after commonHash for remote (This will include common deltas)
	for _, deltaValue := range remoteDeviatedDeltas {
		delta := *deltaValue
		delta.SetDag(localDag, localDeltaTree, false)
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
		delta.SetDag(localDag, localDeltaTree, false)
		delta.SetDeltaTree(localDeltaTree)

	}

	fmt.Println("Applying localDeltasToKeep")
	// Apply all localDeltasToKeep
	for _, delta := range localDeltasToKeep {
		delta.InvertOp()
		delta.SetDag(localDag, localDeltaTree, false)
		delta.SetDeltaTree(localDeltaTree)
	}

	return nil, nil
}

func invertDeltaOp(deltas []*dag.Delta) ([]*dag.Delta, error) {
	var copies []*dag.Delta

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
	Vertexes map[string]*dag.Delta
	Edges    map[string]*dag.Delta
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

func squashIntoState(deltas []*dag.Delta) (*State, error) {
	state := &State{
		Vertexes: make(map[string]*dag.Delta),
		Edges:    make(map[string]*dag.Delta),
	}
	for _, deltaPtr := range deltas {
		delta := *deltaPtr
		op := delta.GetOp()
		gid := delta.GetGid()

		switch op {
		case dag.AddVertex:
			if existingDeltaPtr, ok := state.Vertexes[gid]; !ok {
				state.Vertexes[gid] = &delta
			} else {
				// Merge
				existingDelta := *existingDeltaPtr
				existingOp := existingDelta.GetOp()
				switch existingOp {
				case dag.AddVertex:
					continue
				case dag.RemoveVertex:
					delete(state.Vertexes, gid)
				}
			}
		case dag.RemoveVertex:
			if existingDeltaPtr, ok := state.Vertexes[gid]; !ok {
				state.Vertexes[gid] = &delta
			} else {
				// Merge
				existingDelta := *existingDeltaPtr
				existingOp := existingDelta.GetOp()
				switch existingOp {
				case dag.AddVertex:
					delete(state.Vertexes, gid)
				case dag.RemoveVertex:
					continue
				}
			}
		case dag.AddEdge:
			if existingEdgePtr, ok := state.Edges[gid]; !ok {
				state.Edges[gid] = &delta
			} else {
				// Merge
				existingEdge := *existingEdgePtr
				existingOp := existingEdge.GetOp()
				switch existingOp {
				case dag.AddEdge:
					continue
				case dag.RemoveEdge:
					delete(state.Edges, gid)
				}

			}
		case dag.RemoveEdge:
			if existingEdgePtr, ok := state.Edges[gid]; !ok {
				state.Edges[gid] = &delta
			} else {
				// Merge
				existingEdge := *existingEdgePtr
				existingOp := existingEdge.GetOp()
				switch existingOp {
				case dag.AddEdge:
					delete(state.Edges, gid)
				case dag.RemoveEdge:
					continue
				}
			}
		}
	}
	return state, nil
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
