package keeper_test

import (
	"context"
	"strconv"
	"testing"

	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/achelabov/cosmossdk-example/x/notes/keeper"
	"github.com/achelabov/cosmossdk-example/x/notes/types"
)

func createNNote(keeper keeper.Keeper, ctx context.Context, n int) []types.Note {
	items := make([]types.Note, n)
	for i := range items {
		items[i].Index = strconv.Itoa(i)
		items[i].Text = strconv.Itoa(i)
		items[i].Creator = strconv.Itoa(i)
		_ = keeper.Note.Set(ctx, items[i].Index, items[i])
	}
	return items
}

func TestNoteQuerySingle(t *testing.T) {
	f := initFixture(t)
	qs := keeper.NewQueryServerImpl(f.keeper)
	msgs := createNNote(f.keeper, f.ctx, 2)
	tests := []struct {
		desc     string
		request  *types.QueryGetNoteRequest
		response *types.QueryGetNoteResponse
		err      error
	}{
		{
			desc: "First",
			request: &types.QueryGetNoteRequest{
				Index: msgs[0].Index,
			},
			response: &types.QueryGetNoteResponse{Note: msgs[0]},
		},
		{
			desc: "Second",
			request: &types.QueryGetNoteRequest{
				Index: msgs[1].Index,
			},
			response: &types.QueryGetNoteResponse{Note: msgs[1]},
		},
		{
			desc: "KeyNotFound",
			request: &types.QueryGetNoteRequest{
				Index: strconv.Itoa(100000),
			},
			err: status.Error(codes.NotFound, "not found"),
		},
		{
			desc: "InvalidRequest",
			err:  status.Error(codes.InvalidArgument, "invalid request"),
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			response, err := qs.GetNote(f.ctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				require.EqualExportedValues(t, tc.response, response)
			}
		})
	}
}

func TestNoteQueryPaginated(t *testing.T) {
	f := initFixture(t)
	qs := keeper.NewQueryServerImpl(f.keeper)
	msgs := createNNote(f.keeper, f.ctx, 5)

	request := func(next []byte, offset, limit uint64, total bool) *types.QueryAllNoteRequest {
		return &types.QueryAllNoteRequest{
			Pagination: &query.PageRequest{
				Key:        next,
				Offset:     offset,
				Limit:      limit,
				CountTotal: total,
			},
		}
	}
	t.Run("ByOffset", func(t *testing.T) {
		step := 2
		for i := 0; i < len(msgs); i += step {
			resp, err := qs.ListNote(f.ctx, request(nil, uint64(i), uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.Note), step)
			require.Subset(t, msgs, resp.Note)
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(msgs); i += step {
			resp, err := qs.ListNote(f.ctx, request(next, 0, uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.Note), step)
			require.Subset(t, msgs, resp.Note)
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		resp, err := qs.ListNote(f.ctx, request(nil, 0, 0, true))
		require.NoError(t, err)
		require.Equal(t, len(msgs), int(resp.Pagination.Total))
		require.EqualExportedValues(t, msgs, resp.Note)
	})
	t.Run("InvalidRequest", func(t *testing.T) {
		_, err := qs.ListNote(f.ctx, nil)
		require.ErrorIs(t, err, status.Error(codes.InvalidArgument, "invalid request"))
	})
}
