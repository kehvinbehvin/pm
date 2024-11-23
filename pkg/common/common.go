package common

type AlphaList struct {
	Alphas []Alpha
	Branches map[string]int
}

func NewAlphaList() AlphaList {
	return AlphaList{
		Alphas: []Alpha{},
		Branches: make(map[string]int),
	}
}

func (al AlphaList) MergeIn() {
}

func (al AlphaList) Diff() {
}
