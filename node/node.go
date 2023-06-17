package node

import "github.com/matthew-inamdar/dougdb/state"

type Node struct {
	config *Config
	state  state.State
}

func NewNode(config *Config) *Node {
	return &Node{
		config: config,
		state:  &state.FollowerState{},
	}
}

func (n *Node) Start() error {
	// TODO: Implement.

	// TODO: Attempt to recover from WAL if present.

	return nil
}
