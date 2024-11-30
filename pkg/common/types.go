package common

const (
	AddVertexAlpha      byte = 1
	RemoveVertexAlpha   byte = 2
	AddEdgeAlpha        byte = 3
	RemoveEdgeAlpha     byte = 4
	AddTrieNodeAlpha    byte = 5
	RemoveTrieNodeAlpha byte = 6
)

type Alpha interface {
	GetType() byte
	GetId() string
	GetHash() string
	SetHash(Alpha)
}

type AlphaHistory interface {
	MergeIn(AlphaList)
	Diff(AlphaList)
}

type DataStructure interface {
	Update(Alpha) error
	Rewind(Alpha) error
	Validate(Alpha) bool
}
