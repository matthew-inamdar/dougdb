package node

import (
	"errors"
	"fmt"
)

type FollowerRole struct{}

func (r FollowerRole) HandleAppendEntriesRPC(n *Node, rpc AppendEntriesRPCRequest) (AppendEntriesRPCResponse, error) {
	//
	prevLogIndex, ok := n.log.At(rpc.PrevLogIndex)
	if !ok {
		return AppendEntriesRPCResponse{}, fmt.Errorf("log is smaller than prevLogIndex %d", rpc.PrevLogIndex)
	}
	if prevLogIndex.Term != rpc.PrevLogTerm {
		return AppendEntriesRPCResponse{}, fmt.Errorf(
			"log term at prevLogIndex %d is %d, RPC expected %d",
			rpc.PrevLogIndex,
			prevLogIndex.Term,
			rpc.PrevLogIndex,
		)
	}

	// TODO: Remove.
	var res AppendEntriesRPCResponse
	return res, errors.New("not implemented")
}

func (r FollowerRole) HandleRequestVoteRPC(n *Node, rpc RequestVoteRPCRequest) (RequestVoteRPCResponse, error) {
	var res RequestVoteRPCResponse
	return res, errors.New("not implemented")
}

func (r FollowerRole) HandleExists(n *Node, key string) (bool, error) {
	return false, RedirectToLeaderError{LeaderMember: n.currentLeader}
}

func (r FollowerRole) HandleGet(n *Node, key string) ([]byte, error) {
	return nil, RedirectToLeaderError{LeaderMember: n.currentLeader}
}

func (r FollowerRole) HandleSet(n *Node, key string, value []byte) error {
	return RedirectToLeaderError{LeaderMember: n.currentLeader}
}
