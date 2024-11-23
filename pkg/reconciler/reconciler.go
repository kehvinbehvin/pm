package reconciler

import (
	"errors"
	"math"

	"github/pm/pkg/common"
)

const (
	IncumbentAhead byte = 1;
	InBoundAhead byte = 2;
	Identical byte = 3;
	Conflict byte = 4;
	Irreconcilable byte = 5;
)

type Reconciler struct {
	reconcilable common.Reconcilable
}

func (rc *Reconciler) Reconcile(other *common.Reconcilable) byte {
	incumbentAlphaList := rc.reconcilable.AlphaList;
	inBoundAlphaList := other.AlphaList;

	// Valid Merges 
	// InBound ahead of Incumbent
	// InBound conflic with Invcumbent

	// Invalid Merges 
	// Incumbent ahead of Inbound 
	// Inbound and Incumbent have no common alpha
	
	// Find common alpha
	latestCommonAlpha, alphaErr := rc.findLastCommonAlpha(incumbentAlphaList, inBoundAlphaList);
	if alphaErr != nil {
		// No common alpha found, cannot merge
		return Irreconcilable
	}

	lastIncumbentAlpha := incumbentAlphaList.Alphas[len(incumbentAlphaList.Alphas) - 1].GetId()
	lastInBoundAlpha := inBoundAlphaList.Alphas[len(inBoundAlphaList.Alphas) - 1].GetId()
	commonAlpha := latestCommonAlpha.GetId()


	identical := rc.areAlphaListsIdentical(lastIncumbentAlpha, lastInBoundAlpha, commonAlpha)
	if identical {
		// Both lists are identical, nothing to merge
		return Identical
	}

	hasConflict := rc.hasConflict(lastIncumbentAlpha, lastInBoundAlpha, commonAlpha)
	if hasConflict {
		// Perform default rebase algo
		return Conflict
	}

	incumbentAhead := rc.isIncumbentAhead(lastIncumbentAlpha, lastInBoundAlpha, commonAlpha)
	if incumbentAhead {
		// Incumbent ahead, nothing to merge into incumbent
		return IncumbentAhead
	}

	inBoundAhead := rc.isInBoundAhead(lastIncumbentAlpha, lastInBoundAlpha, commonAlpha)
	if inBoundAhead {
		// Merge inBound in
		return InBoundAhead
	}

	// Should not reach this path of execution as all the cases have been accounted for
	// Log unexpected error here
	// Return as irreconcilabe as the state is unrecognisable
	return Irreconcilable
}

func (rc *Reconciler) MergeIn(other *common.Reconcilable) {
	if other == nil {
		return;
	}

	reconcilingStrategy := rc.Reconcile(other);
	switch reconcilingStrategy {
	case Conflict:
		// Perform Conflict strategy
	case InBoundAhead:
		// Perform Fast forward strategy
	case Identical:
	case IncumbentAhead:
	case Irreconcilable:
		// Log and return for (Identical, IncumbentAhead and Irreconcilable)
	}
}

func (rc *Reconciler) findLastCommonAlpha(incumbent common.AlphaList, inBound common.AlphaList) (common.Alpha, error) {
	incumbentRootAlphaId := incumbent.Alphas[0].GetId();
	inBoundRootAlphaId := inBound.Alphas[0].GetId();

	if inBoundRootAlphaId != incumbentRootAlphaId {
		return nil, errors.New("[findLastCommonAlpha] No common alphas found");
	}
	
	commonAlpha, searchErr := rc.optimisticSearch(incumbent, inBound)
	if searchErr != nil {
		return nil, searchErr
	}

	return commonAlpha, nil
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

	if isCommonAlpha(shortAlpha, longAlpha ) {
		// If 1 last alpha is found in another list, it will be the last common alpha for both
		return shorterList[pointer], nil
	}

	// At this point, the 2 lists have diverged. The common alpha must be between the 1st and last alpha of the shorter list
	// We can perform a binary optimisticSearch
	alpha, searchErr := binarySearch(pointer, shorterList, longerList)
	if searchErr != nil {
		return nil, searchErr;
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
		next := int(nextFloatPointer);

		return binarySearch(next, shorter, longer)
	} else {
		// Test if we found latest alpha
		// Safe to assume that this is not the last alpha in the shorterList
		// Because that case was handled in optimisticSearch and exited early
		// The only case this can occur is for alphas before the last alpha of the shorterlist
		nextShorterAlpha := shorter[pointer + 1]
		nextLongerAlpha := longer[pointer + 1]

		if isCommonAlpha(nextShorterAlpha, nextLongerAlpha) {
			// Not the last commona alpha
			shortLength := len(shorter)
			difference := shortLength - pointer;
			
			floatDifference := float64(difference)
			nextFloatDifference := math.Round(floatDifference / 2)
			nextDifference := int(nextFloatDifference);

			diff := shortLength + nextDifference;

			return binarySearch(diff, shorter, longer)
		} else {
			// Last common alpha found
			return shorter[pointer], nil
		}
	}
}

func (rc *Reconciler) areAlphaListsIdentical(lastIncumbentAlpha string, lastInBoundAlpha string, commonAlpha string) bool {	
	return (commonAlpha == lastIncumbentAlpha) && (lastInBoundAlpha == commonAlpha)
}

func (rc *Reconciler) hasConflict(lastIncumbentAlpha string, lastInBoundAlpha string, commonAlpha string) bool {
	return (commonAlpha != lastInBoundAlpha) && (commonAlpha != lastIncumbentAlpha)
}

func (rc *Reconciler) isIncumbentAhead(lastIncumbentAlpha string, lastInBoundAlpha string, commonAlpha string) bool {
	return (commonAlpha == lastInBoundAlpha) && (commonAlpha != lastIncumbentAlpha)
}

func (rc *Reconciler) isInBoundAhead(lastIncumbentAlpha string, lastInBoundAlpha string, commonAlpha string) bool {
	return (commonAlpha != lastInBoundAlpha) && (commonAlpha == lastIncumbentAlpha)
}

func isCommonAlpha(shorter common.Alpha, longer common.Alpha) bool {
    return shorter.GetId() == longer.GetId()
}
