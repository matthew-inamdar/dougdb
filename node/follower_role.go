package node

import (
	"errors"
)

type FollowerRole struct{}

func (r *FollowerRole) HandleAppendEntriesRPC(rpc AppendEntriesRPCRequest) (AppendEntriesRPCResponse, error) {
	var res AppendEntriesRPCResponse
	return res, errors.New("not implemented")
}

func (r *FollowerRole) HandleRequestVoteRPC(rpc RequestVoteRPCRequest) (RequestVoteRPCResponse, error) {
	var res RequestVoteRPCResponse
	return res, errors.New("not implemented")
}

func (r *FollowerRole) HandleExists(n *Node, key string) (bool, error) {
	return false, RedirectToLeaderError{LeaderMember: n.config.}
}

func (r *FollowerRole) HandleGet(n *Node, key string) ([]byte, error) {
	v, ok := n.db.Get(key)
	if !ok {
		return nil, nil
	}
	return v, nil
}

func (r *FollowerRole) HandleSet(n *Node, key string, value []byte) error {
	n.db.Set(key, value)
	return nil
}
