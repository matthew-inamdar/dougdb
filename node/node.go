package node

import (
	"github.com/matthew-inamdar/dougdb/db"
	"github.com/matthew-inamdar/dougdb/log"
	"github.com/matthew-inamdar/dougdb/store"
)

type Node struct {
	config        *Config
	role          Role
	db            *db.DB
	currentLeader *Server

	// Persistent state.
	currentTerm *store.Store
	votedFor    *store.Store
	log         *log.Log

	// Volatile state.
	commitIndex int
	lastApplied int
	nextIndex   map[*Server]int
	matchIndex  map[*Server]int
}

func NewNode(config *Config) (*Node, error) {
	curTerm, err := store.NewStore("current_term", config.ThisServer.ID)
	if err != nil {
		return nil, err
	}

	votedFor, err := store.NewStore("voted_for", config.ThisServer.ID)
	if err != nil {
		return nil, err
	}

	l, err := log.NewLog(config.ThisServer.ID)
	if err != nil {
		return nil, err
	}

	return &Node{
		config:      config,
		db:          db.NewDB(),
		role:        &FollowerRole{},
		currentTerm: curTerm,
		votedFor:    votedFor,
		log:         l,
		nextIndex:   make(map[*Server]int),
		matchIndex:  make(map[*Server]int),
	}, nil
}

func (n *Node) Start() error {
	// TODO: Implement.

	// TODO: Attempt to recover from WAL if present.
	// TODO: Listen for RPC requests.
	// TODO: Listen for client requests.

	return nil
}
