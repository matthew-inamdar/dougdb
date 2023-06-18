package node

import "errors"

type CandidateRole struct{}

func (r CandidateRole) HandleAppendEntriesRPC(n *Node, rpc AppendEntriesRPCRequest) (AppendEntriesRPCResponse, error) {
	var res AppendEntriesRPCResponse
	return res, errors.New("not implemented")
}

func (r CandidateRole) HandleRequestVoteRPC(n *Node, rpc RequestVoteRPCRequest) (RequestVoteRPCResponse, error) {
	var res RequestVoteRPCResponse
	return res, errors.New("not implemented")
}

func (r CandidateRole) HandleExists(n *Node, key string) (bool, error) {
	return false, RedirectToLeaderError{LeaderMember: n.currentLeader}
}

func (r CandidateRole) HandleGet(n *Node, key string) ([]byte, error) {
	return nil, RedirectToLeaderError{LeaderMember: n.currentLeader}
}

func (r CandidateRole) HandleSet(n *Node, key string, value []byte) error {
	return RedirectToLeaderError{LeaderMember: n.currentLeader}
}
