package raft

import (
	"context"

	raftgrpc "github.com/matthew-inamdar/dougdb/gen/grpc/raft"
	"github.com/matthew-inamdar/dougdb/node"
)

type Service struct {
	raftgrpc.UnimplementedRaftServer

	node *node.Node
}

func (s *Service) AppendEntries(ctx context.Context, req *raftgrpc.AppendEntriesRequest) (*raftgrpc.AppendEntriesResponse, error) {
	s.node.Lock()
	defer s.node.Unlock()

	select {
	case <-ctx.Done():
		// TODO: gRPC error response
		return nil, ctx.Err()
	default:
	}

	// Check if leader's term is less than the current term.
	ct := s.node.CurrentTerm()

	if req.GetTerm() < ct {
		return &raftgrpc.AppendEntriesResponse{
			Term:    ct,
			Success: false,
		}, nil
	}

	// Check if the log doesn't contain an entry at the prevLogIndex with the
	// same term as the prevLogTerm.
	if !s.node.HasEntry(req.GetPrevLogIndex(), req.GetPrevLogTerm()) {
		// TODO: What happens here???
		return &raftgrpc.AppendEntriesResponse{
			Term:    ct,
			Success: false,
		}, nil
	}

	// TODO: Replace and append entries starting from the prevLogIndex+1.

	// TODO: Update commitIndex and apply committed entries to the state machine.

	// TODO
	return nil, nil
}

func (s *Service) RequestVote(ctx context.Context, req *raftgrpc.RequestVoteRequest) (*raftgrpc.RequestVoteResponse, error) {
	// TODO
	return nil, nil
}

var _ raftgrpc.RaftServer = (*Service)(nil)
