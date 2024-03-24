package raft

import (
	"context"
	"errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

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

	// Check if leader's term is less than the current term (stale term).
	if req.GetTerm() < s.node.CurrentTerm() {
		return &raftgrpc.AppendEntriesResponse{
			Term:    s.node.CurrentTerm(),
			Success: false,
		}, nil
	}

	// If leader's term is greater or equal, then update to match and convert to
	// follower (if not already).
	if req.GetTerm() >= s.node.CurrentTerm() {
		if req.GetTerm() > s.node.CurrentTerm() {
			s.node.SetCurrentTerm(req.GetTerm())
		}
		if s.node.Role() != node.RoleFollower {
			s.node.SetRole(node.RoleFollower)
		}
	}

	// Check if the log doesn't contain an entry at the prevLogIndex with the
	// same term as the prevLogTerm.
	if !s.node.HasEntry(req.GetPrevLogIndex(), req.GetPrevLogTerm()) {
		// The leader should retry with prevLogIndex-1.
		return &raftgrpc.AppendEntriesResponse{
			Term:    s.node.CurrentTerm(),
			Success: false,
		}, nil
	}

	// Replace and append entries starting from the prevLogIndex+1.
	for i, e := range req.GetEntries() {
		op, err := mapRPCOperation(e.GetOperation())
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "operation %d not supported", op)
		}

		entry := node.Entry{
			Term:      e.GetTerm(),
			Operation: op,
			Key:       e.GetKey(),
			Value:     e.GetValue(),
		}
		entryIdx := uint64(i+1) + req.GetPrevLogIndex()
		err = s.node.AddEntry(entryIdx, entry)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed adding entry with index %d to log", entryIdx)
		}
	}

	// Apply committed entries to the state machine.
	for i := s.node.LastApplied() + 1; i <= req.GetLeaderCommit(); i++ {
		err := s.node.Commit(ctx, i)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed commiting entry with index %d", i)
		}
	}

	// Update commitIndex.
	s.node.ConditionallySetCommitIndex(req.GetLeaderCommit())

	return &raftgrpc.AppendEntriesResponse{
		Term:    s.node.CurrentTerm(),
		Success: true,
	}, nil
}

func (s *Service) RequestVote(ctx context.Context, req *raftgrpc.RequestVoteRequest) (*raftgrpc.RequestVoteResponse, error) {
	s.node.Lock()
	defer s.node.Unlock()

	select {
	case <-ctx.Done():
		// TODO: gRPC error response
		return nil, ctx.Err()
	default:
	}

	// Check if candidate's term is less than current term (stale term).
	if req.GetTerm() < s.node.CurrentTerm() {
		return &raftgrpc.RequestVoteResponse{
			Term:        s.node.CurrentTerm(),
			VoteGranted: false,
		}, nil
	}

	// If leader's term is greater or equal, then update to match and convert to
	// follower (if not already).
	if req.GetTerm() >= s.node.CurrentTerm() {
		if req.GetTerm() > s.node.CurrentTerm() {
			s.node.SetCurrentTerm(req.GetTerm())
		}
		if s.node.Role() != node.RoleFollower {
			s.node.SetRole(node.RoleFollower)
		}
	}

	// If candidate has already voted, then do not grant vote.
	if s.node.VotedFor() != node.NilID {
		return &raftgrpc.RequestVoteResponse{
			Term:        s.node.CurrentTerm(),
			VoteGranted: false,
		}, nil
	}

	// TODO: If votedFor is null or candidateID, and candidate's log is at least
	// TODO: as up-to-date as receiver's log, grant vote (§5.2, §5.4)

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
