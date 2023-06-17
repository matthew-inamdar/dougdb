package state

import "errors"

type FollowerState struct{}

func (s *FollowerState) HandleAppendEntriesRPC(rpc AppendEntriesRPCRequest) (AppendEntriesRPCResponse, error) {
	var res AppendEntriesRPCResponse
	return res, errors.New("not implemented")
}

func (s *FollowerState) HandleRequestVoteRPC(rpc RequestVoteRPCRequest) (RequestVoteRPCResponse, error) {
	var res RequestVoteRPCResponse
	return res, errors.New("not implemented")
}
