package keeper

import (
	"context"
	"errors"

	"cosmossdk.io/collections"
	"github.com/achelabov/cosmossdk-example/x/notes/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (q queryServer) ListNote(ctx context.Context, req *types.QueryAllNoteRequest) (*types.QueryAllNoteResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	notes, pageRes, err := query.CollectionPaginate(
		ctx,
		q.k.Note,
		req.Pagination,
		func(_ string, value types.Note) (types.Note, error) {
			return value, nil
		},
	)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllNoteResponse{Note: notes, Pagination: pageRes}, nil
}

func (q queryServer) GetNote(ctx context.Context, req *types.QueryGetNoteRequest) (*types.QueryGetNoteResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	val, err := q.k.Note.Get(ctx, req.Index)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "not found")
		}

		return nil, status.Error(codes.Internal, "internal error")
	}

	return &types.QueryGetNoteResponse{Note: val}, nil
}
