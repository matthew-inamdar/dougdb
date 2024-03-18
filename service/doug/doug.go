package doug

import (
	"context"

	douggrpc "github.com/matthew-inamdar/dougdb/gen/grpc/doug"
	"github.com/matthew-inamdar/dougdb/node"
)

type Service struct {
	douggrpc.UnimplementedDougServer

	dbState node.State
}

func (s *Service) Get(ctx context.Context, req *douggrpc.GetRequest) (*douggrpc.GetResponse, error) {
	v, err := s.dbState.Get(ctx, req.GetKey())
	if err != nil {
		// TODO: gRPC error response
		return nil, err
	}
	return &douggrpc.GetResponse{
		Value: v,
	}, nil
}

func (s *Service) Put(ctx context.Context, req *douggrpc.PutRequest) (*douggrpc.PutResponse, error) {
	err := s.dbState.Put(ctx, req.GetKey(), req.GetValue())
	if err != nil {
		// TODO: gRPC error response
		return nil, err
	}
	return &douggrpc.PutResponse{}, nil
}

func (s *Service) Delete(ctx context.Context, req *douggrpc.DeleteRequest) (*douggrpc.DeleteResponse, error) {
	err := s.dbState.Delete(ctx, req.GetKey())
	if err != nil {
		// TODO: gRPC error response
		return nil, err
	}
	return &douggrpc.DeleteResponse{}, nil
}

var _ douggrpc.DougServer = (*Service)(nil)
