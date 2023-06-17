package node

import (
	"github.com/matthew-inamdar/dougdb/db"
	"github.com/matthew-inamdar/dougdb/log"
	"github.com/matthew-inamdar/dougdb/state"
)

type Node struct {
	config *Config
	state  state.State
	db     *db.DB
	log    *log.Log
}

func NewNode(config *Config) (*Node, error) {
	l, err := log.NewLog(config.Node.ID)
	if err != nil {
		return nil, err
	}

	return &Node{
		config: config,
		state:  &state.FollowerState{},
		db:     db.NewDB(),
		log:    l,
	}, nil
}

func (n *Node) Start() error {
	// TODO: Implement.

	// TODO: Attempt to recover from WAL if present.
	// TODO: Listen for RPC requests.

	return nil
}
