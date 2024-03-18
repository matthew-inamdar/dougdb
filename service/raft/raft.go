package raft

import (
	"context"

	raftgrpc "github.com/matthew-inamdar/dougdb/gen/grpc/raft"
)

type Service struct {
	raftgrpc.UnimplementedRaftServer
}

func (s *Service) AppendEntries(ctx context.Context, req *raftgrpc.AppendEntriesRequest) (*raftgrpc.AppendEntriesResponse, error) {
	// TODO
	return nil, nil
}

func (s *Service) RequestVote(ctx context.Context, req *raftgrpc.RequestVoteRequest) (*raftgrpc.RequestVoteResponse, error) {
	// TODO
	return nil, nil
}

var _ raftgrpc.RaftServer = (*Service)(nil)
