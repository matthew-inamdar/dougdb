package node

type Role int

const (
	RoleFollower Role = iota
	RoleCandidate
	RoleLeader
)

type ID string

type Node struct {
	// Persistent state.
	role        Role
	currentTerm uint64
	votedFor    ID
	log         []Entry

	// Volatile state.
	commitIndex uint64
	lastApplied uint64

	// Volatile leader state.
	nextIndex  map[ID]uint64
	matchIndex map[ID]uint64
}
