package dag

type Graph struct {}

type GraphAlpha struct {}

func (G Graph) Merge(other Graph) (bool, error) {
	return true, nil
}

func (G Graph) Compare(other Graph) ([]GraphAlpha, error) {
	var diffs []GraphAlpha;

	return diffs, nil
}

func (G Graph) Alpha() Alpha {
	return GraphAlpha{}
}

func (GA GraphAlpha) Apply() (bool, error) {
	return true, nil
}

type LinkedList struct {}

type LinkedListAlpha struct {}

func (LL LinkedList) Merge(other Graph) (bool, error) {
	return true, nil
}

func (LL LinkedList) Compare(other Graph) ([]LinkedListAlpha, error) {
	var diffs []LinkedListAlpha;

	return diffs, nil
}

func (LL LinkedList) Alpha() Alpha {
	return LinkedListAlpha{}
}

func (LLA LinkedListAlpha) Apply() (bool, error) {
	return true, nil
}

type DataStruct interface {
	Graph | LinkedList
}

type CompositeStruct[D DataStruct] struct {
	Structure D
	History AlphaLog[D]
}

type Alpha interface {
	Apply() (bool, error)
}

type AlphaLog[D DataStruct] struct {
	log []D
}

type Mergable[D DataStruct] interface {
	Compare(D) ([]Alpha, error)
	Merge(D) (bool, error)
}

type Trackable interface {
	Apply() Alpha
}
/**
Alpha: A unit of change. Generalisable across all data structures.
DataStruct: A structure that is
- Mergable: Able to merge in new Alphas into its own structure
- Comparable: Able to compare itself against another of its type and return a list of alphas
- Loggable: Able to log changes.

Mergable:
- Merge(T DataStruct)() T

Comparable:
- Compare(T DataStruct)() T



*/
