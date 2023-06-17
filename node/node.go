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
	return nil
}
