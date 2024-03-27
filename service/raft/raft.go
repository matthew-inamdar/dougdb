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
		return nil, status.Error(codes.Canceled, "client canceled request")
	default:
	}

	// Check if leader's term is less than the current term (stale term).
	if req.GetTerm() < s.node.CurrentTerm() {
		return &raftgrpc.AppendEntriesResponse{
			Term:    s.node.CurrentTerm(),
			Success: false,
		}, nil
	}

	// If the leader's term is greater than the current term, then update
	// internal state.
	if req.GetTerm() > s.node.CurrentTerm() {
		s.node.SetCurrentTerm(req.GetTerm())
	}

	// If this node isn't a follower, then convert to follower.
	if s.node.Role() != node.RoleFollower {
		s.node.SetRole(node.RoleFollower)
	}

	// If this is a new leader, then update internal state.
	if leaderID := node.ID(req.GetLeaderID()); s.node.LeaderID() != leaderID {
		s.node.SetLeaderID(leaderID)
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
		err := s.node.Apply(ctx, i)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed applying entry with index %d", i)
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
		return nil, status.Error(codes.Canceled, "client canceled request")
	default:
	}

	// Check if candidate's term is less than current term (stale term).
	if req.GetTerm() < s.node.CurrentTerm() {
		return &raftgrpc.RequestVoteResponse{
			Term:        s.node.CurrentTerm(),
			VoteGranted: false,
		}, nil
	}

	// If the leader's term is greater than the current term, then update
	// internal state.
	if req.GetTerm() > s.node.CurrentTerm() {
		s.node.SetCurrentTerm(req.GetTerm())
	}

	// If this node isn't a follower, then convert to follower.
	if s.node.Role() != node.RoleFollower {
		s.node.SetRole(node.RoleFollower)
	}

	// If the candidate has already been granted the vote for the current term,
	// then idempotently grant again.
	if node.ID(req.GetCandidateID()) == s.node.VotedFor() {
		return &raftgrpc.RequestVoteResponse{
			Term:        s.node.CurrentTerm(),
			VoteGranted: true,
		}, nil
	}

	// If vote has already been granted for current term, then do not grant
	// another vote.
	if s.node.VotedFor() != node.NilID {
		return &raftgrpc.RequestVoteResponse{
			Term:        s.node.CurrentTerm(),
			VoteGranted: false,
		}, nil
	}

	// If current node has larger term in last log entry, then do not grant
	// vote.
	if req.GetLastLogTerm() < s.node.LastLogTerm() {
		return &raftgrpc.RequestVoteResponse{
			Term:        s.node.CurrentTerm(),
			VoteGranted: false,
		}, nil
	}

	// If current node has same term in last log entry, and current node has
	// more entries (higher index), then do not grant vote.
	if req.GetLastLogTerm() == s.node.LastLogTerm() && req.GetLastLogIndex() < s.node.LastLogIndex() {
		return &raftgrpc.RequestVoteResponse{
			Term:        s.node.CurrentTerm(),
			VoteGranted: false,
		}, nil
	}

	// Candidate is suitable for leadership, so grant vote.
	return &raftgrpc.RequestVoteResponse{
		Term:        s.node.CurrentTerm(),
		VoteGranted: true,
	}, nil
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
