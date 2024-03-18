package raft

import (
	"context"
	"errors"

	raftgrpc "github.com/matthew-inamdar/dougdb/gen/grpc/raft"
	"github.com/matthew-inamdar/dougdb/node"
)

var (
	ErrUnsupportedOperation = errors.New("the operation is unsupported")
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

	// Replace and append entries starting from the prevLogIndex+1.
	for i, e := range req.GetEntries() {
		op, err := mapRPCOperation(e.GetOperation())
		if err != nil {
			// TODO: Return gRPC error.
			return nil, err
		}

		entry := node.Entry{
			Term:      e.GetTerm(),
			Operation: op,
			Key:       e.GetKey(),
			Value:     e.GetValue(),
		}
		err = s.node.AddEntry(uint64(i+1)+req.GetPrevLogIndex(), entry)
	}

	// TODO: Update commitIndex and apply committed entries to the state machine.

	return &raftgrpc.AppendEntriesResponse{
		Term:    ct,
		Success: true,
	}, nil
}

func (s *Service) RequestVote(ctx context.Context, req *raftgrpc.RequestVoteRequest) (*raftgrpc.RequestVoteResponse, error) {
	// TODO
	return nil, nil
}

var _ raftgrpc.RaftServer = (*Service)(nil)

func mapRPCOperation(operation raftgrpc.AppendEntriesRequest_Entry_Operation) (node.Operation, error) {
	switch operation {
	case raftgrpc.AppendEntriesRequest_Entry_OPERATION_PUT:
		return node.OperationPut, nil
	case raftgrpc.AppendEntriesRequest_Entry_OPERATION_DELETE:
		return node.OperationDelete, nil
	default:
		return 0, ErrUnsupportedOperation
	}
}
