package node

import (
	"context"
	"errors"
	"sync"
)

type role int

const (
	roleFollower role = iota
	roleCandidate
	roleLeader
)

var (
	ErrEntryIndexOutOfRange = errors.New("the index of the entry is out of range")
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

	// Raft State.
	sync.Mutex
	persistentState
	volatileState
	volatileLeaderState

	// DB State.
	DBState State
}

func NewNode(id ID) *Node {
	return &Node{
		id:    id,
		Mutex: sync.Mutex{},
		volatileLeaderState: volatileLeaderState{
			nextIndex:  make(map[ID]uint64),
			matchIndex: make(map[ID]uint64),
		},
		DBState: newMemoryState(),
	}
}

func (n *Node) CurrentTerm() uint64 {
	return n.persistentState.currentTerm
}

func (n *Node) HasEntry(index, term uint64) bool {
	if index >= uint64(len(n.log)) {
		return false
	}

	return n.log[index].Term == term
}

func (n *Node) AddEntry(index uint64, entry Entry) error {
	if index > uint64(len(n.log)) {
		return ErrEntryIndexOutOfRange
	}

	if index == uint64(len(n.log)) {
		n.log = append(n.log, entry)
		return nil
	}

	n.log[index] = entry
	return nil
}

func (n *Node) CommitIndex() uint64 {
	return n.commitIndex
}

func (n *Node) Commit(ctx context.Context, index uint64) error {
	if index >= uint64(len(n.log)) {
		return ErrEntryIndexOutOfRange
	}

	switch n.log[index].Operation {
	case OperationPut:
		err := n.DBState.Put(ctx, n.log[index].Key, n.log[index].Value)
		if err != nil {
			return err
		}
	case OperationDelete:
		err := n.DBState.Delete(ctx, n.log[index].Key)
		if err != nil {
			return err
		}
	default:
		// We should logically never arrive here.
		panic("the operation is unsupported")
	}

	n.commitIndex = index
	return nil
}
