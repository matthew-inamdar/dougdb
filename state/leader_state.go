package state

import "errors"

type LeaderState struct{}

func (s *LeaderState) HandleAppendEntriesRPC(rpc AppendEntriesRPCRequest) (AppendEntriesRPCResponse, error) {
	var res AppendEntriesRPCResponse
	return res, errors.New("not implemented")
}

func (s *LeaderState) HandleRequestVoteRPC(rpc RequestVoteRPCRequest) (RequestVoteRPCResponse, error) {
	var res RequestVoteRPCResponse
	return res, errors.New("not implemented")
}
