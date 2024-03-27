package node

import (
	"context"
	"errors"
	"sync"
	"time"
)

type Role int

const (
	RoleFollower Role = iota
	RoleCandidate
	RoleLeader
)

var (
	ErrEntryIndexOutOfRange = errors.New("the index of the entry is out of range")
	ErrUnsupportedOperation = errors.New("the operation is unsupported")
)

type ID string

const NilID ID = ""

type config struct {
	// The ID of this node.
	id ID
	// The IDs of the other nodes in the cluster. This must not contain the ID
	// of this node.
	clusterIDs []ID
	// The timeout to start an election.
	electionTimeout time.Duration
}

type persistentState struct {
	// The current term of this node.
	currentTerm uint64
	// The ID of the node that this node has voted for in the current term.
	votedFor ID
	// The log of entries this node has.
	log []Entry
}

type volatileState struct {
	// The index of the highest known entry to be committed within the cluster,
	// according to this node's knowledge.
	commitIndex uint64
	// The index of the highest committed entry in this node's log.
	lastApplied uint64

	// The ID of the leader node in the cluster. This may be empty if this node
	// is not aware which node is leader.
	leaderID ID
	// The current role of this node.
	role Role
}

type volatileLeaderState struct {
	// For each other node in the cluster, the index of the next entry in the
	// log to attempt to append entries from.
	nextIndex map[ID]uint64
	// For each other node in the cluster, the index of the highest known log
	// entry to have been replicated.
	matchIndex map[ID]uint64
}

type Node struct {
	config

	sync.RWMutex
	persistentState
	volatileState
	volatileLeaderState

	dbState State
}

func NewNode(id ID, clusterIDs []ID, electionTimeout time.Duration) *Node {
	nextIndex, matchIndex := make(map[ID]uint64), make(map[ID]uint64)
	for _, id := range clusterIDs {
		nextIndex[id], matchIndex[id] = 0, 0
	}

	return &Node{
		config: config{
			id:              id,
			clusterIDs:      clusterIDs,
			electionTimeout: electionTimeout,
		},
		RWMutex: sync.RWMutex{},
		volatileLeaderState: volatileLeaderState{
			nextIndex:  nextIndex,
			matchIndex: matchIndex,
		},
		dbState: newMemoryState(),
	}
}

func (n *Node) LeaderID() ID {
	return n.leaderID
}

func (n *Node) SetLeaderID(leaderID ID) {
	n.leaderID = leaderID
}

func (n *Node) CurrentTerm() uint64 {
	return n.currentTerm
}

func (n *Node) HasEntry(index, term uint64) bool {
	if index >= uint64(len(n.log)) {
		return false
	}

	return n.log[index].Term == term
}

func (n *Node) AddEntry(index uint64, entry Entry) error {
	// Node is missing prior entries.
	if index > uint64(len(n.log)) {
		return ErrEntryIndexOutOfRange
	}

	// Needs to overwrite prior entries so remove entries from log up to and
	// excluding new entry index.
	log := make([]Entry, index)
	copy(log, n.log)
	n.log = log

	// Next entry node is expecting.
	n.log = append(n.log, entry)
	return nil
}

func (n *Node) AppendEntry(entry Entry) error {
	return n.AddEntry(uint64(len(n.log)), entry)
}

func (n *Node) LastApplied() uint64 {
	return n.lastApplied
}

func (n *Node) Apply(ctx context.Context, index uint64) error {
	if index >= uint64(len(n.log)) {
		return ErrEntryIndexOutOfRange
	}

	switch n.log[index].Operation {
	case OperationPut:
		err := n.dbState.Put(ctx, n.log[index].Key, n.log[index].Value)
		if err != nil {
			return err
		}
	case OperationDelete:
		err := n.dbState.Delete(ctx, n.log[index].Key)
		if err != nil {
			return err
		}
	default:
		return ErrUnsupportedOperation
	}

	n.lastApplied = index
	return nil
}

func (n *Node) SetCurrentTerm(term uint64) {
	n.currentTerm = term

	// Voted-for must be cleared when the term is changed.
	n.votedFor = NilID
}

func (n *Node) ConditionallySetCommitIndex(leaderCommit uint64) {
	if leaderCommit > n.commitIndex {
		n.commitIndex = leaderCommit
	}
}

func (n *Node) Role() Role {
	return n.role
}

func (n *Node) SetRole(role Role) {
	n.role = role

	if role != RoleLeader {
		return
	}

	// Reset volatile leader state when promoted to leader.
	for id, _ := range n.nextIndex {
		n.nextIndex[id] = n.LastLogIndex() + 1
	}
	for id, _ := range n.matchIndex {
		n.matchIndex[id] = 0
	}
}

func (n *Node) VotedFor() ID {
	return n.votedFor
}

func (n *Node) DBState() State {
	return n.dbState
}

func (n *Node) LastLogTerm() uint64 {
	if len(n.log) == 0 {
		return 0
	}

	return n.log[len(n.log)-1].Term
}

func (n *Node) LastLogIndex() uint64 {
	return uint64(len(n.log)) - 1
}

func (n *Node) LastCommittedLogTerm() uint64 {
	if len(n.log) == 0 {
		return 0
	}

	return n.log[n.commitIndex].Term
}
