package node

import (
	"github.com/matthew-inamdar/dougdb/db"
	"github.com/matthew-inamdar/dougdb/log"
	"github.com/matthew-inamdar/dougdb/store"
)

type Node struct {
	config      *Config
	role        Role
	db          *db.DB
	currentTerm *store.Store
	votedFor    *store.Store
	log         *log.Log
}

func NewNode(config *Config) (*Node, error) {
	curTerm, err := store.NewStore("current_term", config.Node.ID)
	if err != nil {
		return nil, err
	}

	votedFor, err := store.NewStore("voted_for", config.Node.ID)
	if err != nil {
		return nil, err
	}

	l, err := log.NewLog(config.Node.ID)
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
	}, nil
}

func (n *Node) Start() error {
	// TODO: Implement.

	// TODO: Attempt to recover from WAL if present.
	// TODO: Listen for RPC requests.
	// TODO: Listen for client requests.

	return nil
}
