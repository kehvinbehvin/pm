package reconciler

import (
	"errors"
	"fmt"
	"math"

	"github/pm/pkg/common"
)

const (
	IncumbentAhead byte = 1
	InBoundAhead   byte = 2
	Identical      byte = 3
	Conflict       byte = 4
	Irreconcilable byte = 5
)

type Reconciler struct {
	reconcilable    common.Reconcilable
	lastCommonAlpha common.Alpha
	inBoundDelta    []common.Alpha
	incumbentDelta  []common.Alpha
}

func (rc *Reconciler) Reconcile(other common.AlphaList) byte {
	incumbentAlphaList := rc.reconcilable.AlphaList
	inBoundAlphaList := other

	// Valid Merges
	// InBound ahead of Incumbent
	// InBound conflic with Invcumbent

	// Invalid Merges
	// Incumbent ahead of Inbound
	// Inbound and Incumbent have no common alpha

	// Find common alpha
	alphaErr := rc.findLastCommonAlpha(incumbentAlphaList, inBoundAlphaList)
	if alphaErr != nil {
		// No common alpha found, cannot merge
		return Irreconcilable
	}

	lastIncumbentAlpha := incumbentAlphaList.Alphas[len(incumbentAlphaList.Alphas)-1].GetId()
	lastInBoundAlpha := inBoundAlphaList.Alphas[len(inBoundAlphaList.Alphas)-1].GetId()

	identical := rc.areAlphaListsIdentical(lastIncumbentAlpha, lastInBoundAlpha)
	if identical {
		// Both lists are identical, nothing to merge
		return Identical
	}

	hasConflict := rc.hasConflict(lastIncumbentAlpha, lastInBoundAlpha)
	if hasConflict {
		// Perform default rebase algo
		return Conflict
	}

	incumbentAhead := rc.isIncumbentAhead(lastIncumbentAlpha, lastInBoundAlpha)
	if incumbentAhead {
		// Incumbent ahead, nothing to merge into incumbent
		return IncumbentAhead
	}

	inBoundAhead := rc.isInBoundAhead(lastIncumbentAlpha, lastInBoundAlpha)
	if inBoundAhead {
		// Merge inBound in
		return InBoundAhead
	}

	// Should not reach this path of execution as all the cases have been accounted for
	// Log unexpected error here
	// Return as irreconcilabe as the state is unrecognisable
	return Irreconcilable
}

func (rc *Reconciler) MergeIn(other common.AlphaList) error {
	reconcilingStrategy := rc.Reconcile(other)
	error := rc.resolveDeltas(other)
	if error != nil {
		return error
	}

	switch reconcilingStrategy {
	case Conflict:
		// Perform Conflict strategy
		rc.MergeSilently(other)
	case InBoundAhead:
		// Perform Fast forward strategy
		rc.FastForward(other)
	case Identical:
	case IncumbentAhead:
	case Irreconcilable:
		// Log and return for (Identical, IncumbentAhead and Irreconcilable)
		return errors.New("There is nothing to merge")
	}

	return nil
}

func (rc *Reconciler) findLastCommonAlpha(incumbent common.AlphaList, inBound common.AlphaList) error {
	incumbentRootAlphaId := incumbent.Alphas[0].GetId()
	inBoundRootAlphaId := inBound.Alphas[0].GetId()

	if inBoundRootAlphaId != incumbentRootAlphaId {
		return errors.New("[findLastCommonAlpha] No common alphas found")
	}

	commonAlpha, searchErr := rc.optimisticSearch(incumbent, inBound)
	if searchErr != nil {
		return searchErr
	}

	rc.lastCommonAlpha = commonAlpha

	return nil
}

func (rc *Reconciler) optimisticSearch(incumbent common.AlphaList, inBound common.AlphaList) (common.Alpha, error) {
	// Optimistic search, alwaays assume that there are common alphas than uncommon ones
	var longerList []common.Alpha
	var shorterList []common.Alpha

	// Doesnt matter if they are the same length
	if len(incumbent.Alphas) > len(inBound.Alphas) {
		longerList = incumbent.Alphas
		shorterList = inBound.Alphas
	} else {
		longerList = inBound.Alphas
		shorterList = incumbent.Alphas
	}

	// Check if last Alpha of shorter list is in longer list
	var pointer int
	pointer = len(shorterList) - 1

	longAlpha := longerList[pointer]
	shortAlpha := shorterList[pointer]

	if isCommonAlpha(shortAlpha, longAlpha) {
		// If 1 last alpha is found in another list, it will be the last common alpha for both
		return shorterList[pointer], nil
	}

	// At this point, the 2 lists have diverged. The common alpha must be between the 1st and last alpha of the shorter list
	// We can perform a binary optimisticSearch
	alpha, searchErr := binarySearch(pointer, shorterList, longerList)
	if searchErr != nil {
		return nil, searchErr
	}

	return alpha, nil
}

func binarySearch(pointer int, shorter []common.Alpha, longer []common.Alpha) (common.Alpha, error) {
	if len(shorter) == 0 || len(longer) == 0 {
		// Technically should not be able to reach here
		return nil, errors.New("Binary search failed: one of the lists is empty")
	}

	if pointer+1 >= len(shorter) || pointer+1 >= len(longer) {
		// Technically should not be able to reach here
		return nil, errors.New("Index out of bounds in binarySearch")
	}

	if pointer == 0 {
		// No more to search, this is the last common alpha
		return shorter[0], nil
	}

	shorterAlpha := shorter[pointer]
	longerAlpha := longer[pointer]

	if isCommonAlpha(shorterAlpha, longerAlpha) {
		// Common alpha must be earlier
		floatPointer := float64(pointer)
		nextFloatPointer := math.Round(floatPointer / 2)
		next := int(nextFloatPointer)

		return binarySearch(next, shorter, longer)
	} else {
		// Test if we found latest alpha
		// Safe to assume that this is not the last alpha in the shorterList
		// Because that case was handled in optimisticSearch and exited early
		// The only case this can occur is for alphas before the last alpha of the shorterlist
		nextShorterAlpha := shorter[pointer+1]
		nextLongerAlpha := longer[pointer+1]

		if isCommonAlpha(nextShorterAlpha, nextLongerAlpha) {
			// Not the last commona alpha
			shortLength := len(shorter)
			difference := shortLength - pointer

			floatDifference := float64(difference)
			nextFloatDifference := math.Round(floatDifference / 2)
			nextDifference := int(nextFloatDifference)

			diff := shortLength + nextDifference

			return binarySearch(diff, shorter, longer)
		} else {
			// Last common alpha found
			return shorter[pointer], nil
		}
	}
}

func (rc *Reconciler) areAlphaListsIdentical(lastIncumbentAlpha string, lastInBoundAlpha string) bool {
	return (rc.lastCommonAlpha.GetId() == lastIncumbentAlpha) && (lastInBoundAlpha == rc.lastCommonAlpha.GetId())
}

func (rc *Reconciler) hasConflict(lastIncumbentAlpha string, lastInBoundAlpha string) bool {
	return (rc.lastCommonAlpha.GetId() != lastInBoundAlpha) && (rc.lastCommonAlpha.GetId() != lastIncumbentAlpha)
}

func (rc *Reconciler) isIncumbentAhead(lastIncumbentAlpha string, lastInBoundAlpha string) bool {
	return (rc.lastCommonAlpha.GetId() == lastInBoundAlpha) && (rc.lastCommonAlpha.GetId() != lastIncumbentAlpha)
}

func (rc *Reconciler) isInBoundAhead(lastIncumbentAlpha string, lastInBoundAlpha string) bool {
	return (rc.lastCommonAlpha.GetId() != lastInBoundAlpha) && (rc.lastCommonAlpha.GetId() == lastIncumbentAlpha)
}

func isCommonAlpha(shorter common.Alpha, longer common.Alpha) bool {
	return shorter.GetId() == longer.GetId()
}

func (rc *Reconciler) resolveDeltas(inBound common.AlphaList) error {
	var inBoundDelta []common.Alpha
	var incumbentDelta []common.Alpha

	// TODO: May include logic to calculate the lastCommonAlpha here
	if rc.lastCommonAlpha == nil {
		return errors.New("Calculate the last common alpha first before resolving the deltas")
	}

	lastCommonAlphaId := rc.lastCommonAlpha.GetId()

	inBoundLength := len(inBound.Alphas)
	incumbentLength := len(rc.reconcilable.AlphaList.Alphas)
	for i := inBoundLength - 1; i >= 0; i-- {
		tmpAlpha := inBound.Alphas[i]
		if tmpAlpha.GetId() != lastCommonAlphaId {
			inBoundDelta = append(inBoundDelta, tmpAlpha)
		} else {
			break
		}

	}

	for i := incumbentLength - 1; i >= 0; i-- {
		tmpAlpha := rc.reconcilable.AlphaList.Alphas[i]
		if tmpAlpha.GetId() != lastCommonAlphaId {
			incumbentDelta = append(incumbentDelta, tmpAlpha)
		} else {
			break
		}

	}

	rc.inBoundDelta = incumbentDelta
	rc.incumbentDelta = incumbentDelta

	return nil
}

func (rc *Reconciler) MergeSilently(inBound common.AlphaList) error {
	// Reset incumbent reconcilable
	error := rc.reconcilable.Reset(rc.lastCommonAlpha)
	if error != nil {
		return error
	}

	// Set all of inBoundDeltas
	error = rc.reconcilable.FastForward(rc.inBoundDelta)
	if error != nil {
		return error
	}

	// Since we take inBound as the source of truth, we will accept all inBound Alphas first
	// Next, we need to append all the incumbentAlphas. however incumbentAlphas depend on
	// the state of the underlying structure at lastCommonAlpha. Hence depedencies might have been changed
	// This requires us to check for every alpha in the incumbentDelta for changed dependencies
	for incumbentIndex := 0; incumbentIndex < len(rc.incumbentDelta); incumbentIndex++ {
		incumbentAlpha := rc.incumbentDelta[incumbentIndex]

		validationError := rc.reconcilable.DataStructure.Validate(incumbentAlpha)

		if validationError {
			// log error but continue with merge first and skip this change
			// this may have a cascading effect on later alphas in the delta
			fmt.Print("There was a validation error with merging")
			continue
		}

		error := rc.reconcilable.Commit(incumbentAlpha)
		if error != nil {
			return error
		}
	}

	return nil
}

func (rc *Reconciler) FastForward(inBound common.AlphaList) {
	error := rc.reconcilable.FastForward(inBound.Alphas)
	if error != nil {
		// Log error
		fmt.Print("There was an error with FastForwarding")
	}
}
