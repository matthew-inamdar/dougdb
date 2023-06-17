package node

import (
	"github.com/matthew-inamdar/dougdb/db"
	"github.com/matthew-inamdar/dougdb/state"
)

type Node struct {
	config *Config
	state  state.State
	db     *db.DB
}

func NewNode(config *Config) *Node {
	return &Node{
		config: config,
		state:  &state.FollowerState{},
		db:     db.NewDB(),
	}
}

func (n *Node) Start() error {
	// TODO: Implement.

	// TODO: Attempt to recover from WAL if present.

	return nil
}
