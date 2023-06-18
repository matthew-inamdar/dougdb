package node

import (
	"github.com/matthew-inamdar/dougdb/db"
	"github.com/matthew-inamdar/dougdb/log"
)

type Node struct {
	config *Config
	role   Role
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
		db:     db.NewDB(),
		role:   &FollowerRole{},
		log:    l,
	}, nil
}

func (n *Node) Start() error {
	// TODO: Implement.

	// TODO: Attempt to recover from WAL if present.
	// TODO: Listen for RPC requests.
	// TODO: Listen for client requests.

	return nil
}
