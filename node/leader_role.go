package node

import (
	"errors"
)

type LeaderRole struct{}

func (r *LeaderRole) HandleAppendEntriesRPC(rpc AppendEntriesRPCRequest) (AppendEntriesRPCResponse, error) {
	var res AppendEntriesRPCResponse
	return res, errors.New("not implemented")
}

func (r *LeaderRole) HandleRequestVoteRPC(rpc RequestVoteRPCRequest) (RequestVoteRPCResponse, error) {
	var res RequestVoteRPCResponse
	return res, errors.New("not implemented")
}

func (r *LeaderRole) HandleExists(n *Node, key string) (bool, error) {
	return n.db.Has(key), nil
}

func (r *LeaderRole) HandleGet(n *Node, key string) ([]byte, error) {
	v, ok := n.db.Get(key)
	if !ok {
		return nil, nil
	}
	return v, nil
}

func (r *LeaderRole) HandleSet(n *Node, key string, value []byte) error {
	n.db.Set(key, value)
	return nil
}
