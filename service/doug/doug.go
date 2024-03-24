package doug

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	douggrpc "github.com/matthew-inamdar/dougdb/gen/grpc/doug"
	"github.com/matthew-inamdar/dougdb/node"
)

type Service struct {
	douggrpc.UnimplementedDougServer

	node *node.Node
}

func (s *Service) Get(ctx context.Context, req *douggrpc.GetRequest) (*douggrpc.GetResponse, error) {
	s.node.RLock()
	defer s.node.RUnlock()

	select {
	case <-ctx.Done():
		return nil, status.Error(codes.Canceled, "client canceled request")
	default:
	}

	// Redirect if not leader.
	if s.node.Role() != node.RoleLeader {
		return nil, s.redirectToLeaderError()
	}

	// Read value from DB state.
	v, err := s.node.DBState().Get(ctx, req.GetKey())
	if err != nil {
		return nil, status.Error(codes.Internal, "failed reading record")
	}

	// Return value.
	return &douggrpc.GetResponse{
		Value: v,
	}, nil
}

func (s *Service) Put(ctx context.Context, req *douggrpc.PutRequest) (*douggrpc.PutResponse, error) {
	s.node.Lock()
	defer s.node.Unlock()

	select {
	case <-ctx.Done():
		return nil, status.Error(codes.Canceled, "client canceled request")
	default:
	}

	// Redirect if not leader.
	if s.node.Role() != node.RoleLeader {
		return nil, s.redirectToLeaderError()
	}

	// Append entry to internal log.
	err := s.node.AppendEntry(node.Entry{
		Term:      s.node.CurrentTerm(),
		Operation: node.OperationPut,
		Key:       req.GetKey(),
		Value:     req.GetValue(),
	})
	if err != nil {
		return nil, status.Error(codes.Internal, "failed appending entry to log")
	}

	// TODO: Issue AppendEntries RPC to append entry to all nodes (except itself).
	// TODO: Wait for majority to accept entry.
	// TODO: -> Concurrently, continuously retry to nodes that haven't responded.
	// TODO: Once majority accepted, commit locally
	// TODO: ???????? Issue AppendEntries RPC to commit entry all nodes.
	// TODO: Wait until majority has committed entry.

	// Write value to DB state.
	err = s.node.DBState().Put(ctx, req.GetKey(), req.GetValue())
	if err != nil {
		return nil, status.Error(codes.Internal, "failed writing record")
	}

	// Return.
	return &douggrpc.PutResponse{}, nil
}

func (s *Service) Delete(ctx context.Context, req *douggrpc.DeleteRequest) (*douggrpc.DeleteResponse, error) {
	s.node.Lock()
	defer s.node.Unlock()

	select {
	case <-ctx.Done():
		return nil, status.Error(codes.Canceled, "client canceled request")
	default:
	}

	// Redirect if not leader.
	if s.node.Role() != node.RoleLeader {
		return nil, s.redirectToLeaderError()
	}

	// Append entry to internal log.
	err := s.node.AppendEntry(node.Entry{
		Term:      s.node.CurrentTerm(),
		Operation: node.OperationDelete,
		Key:       req.GetKey(),
	})
	if err != nil {
		return nil, status.Error(codes.Internal, "failed appending entry to log")
	}

	// TODO: Issue AppendEntries RPC to append entry to all nodes (except itself).
	// TODO: Wait for majority to accept entry.
	// TODO: -> Concurrently, continuously retry to nodes that haven't responded.
	// TODO: Once majority accepted, commit locally
	// TODO: ???????? Issue AppendEntries RPC to commit entry all nodes.
	// TODO: Wait until majority has committed entry.

	// Remove value to DB state.
	err = s.node.DBState().Delete(ctx, req.GetKey())
	if err != nil {
		return nil, status.Error(codes.Internal, "failed deleting record")
	}

	// Return.
	return &douggrpc.DeleteResponse{}, nil
}

var _ douggrpc.DougServer = (*Service)(nil)

func (s *Service) redirectToLeaderError() error {
	st := status.New(codes.FailedPrecondition, "node is not leader")
	st, err := st.WithDetails(&douggrpc.ErrorDetails{
		IsLeader: false,
		LeaderID: string(s.node.LeaderID()),
	})
	if err != nil {
		return status.Error(codes.Internal, "failed adding gRPC richer error details")
	}
	return st.Err()
}
