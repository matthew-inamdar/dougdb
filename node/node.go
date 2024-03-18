package node

import "sync"

type role int

const (
	roleFollower role = iota
	roleCandidate
	roleLeader
)

type ID string

type persistentState struct {
	role        role
	currentTerm uint64
	votedFor    ID
	log         []Entry
}

type volatileState struct {
	commitIndex uint64
	lastApplied uint64
}

type volatileLeaderState struct {
	nextIndex  map[ID]uint64
	matchIndex map[ID]uint64
}

type Node struct {
	id ID

	// Raft state.
	mu sync.Mutex
	persistentState
	volatileState
	volatileLeaderState

	// DB state.
	DBState state
}

func NewNode(id ID) *Node {
	return &Node{
		id: id,
		mu: sync.Mutex{},
		volatileLeaderState: volatileLeaderState{
			nextIndex:  make(map[ID]uint64),
			matchIndex: make(map[ID]uint64),
		},
		DBState: newMemoryState(),
	}
}
