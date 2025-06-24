package keeper

import (
	"context"
	"errors"
	"strconv"

	errorsmod "cosmossdk.io/errors"
	"github.com/achelabov/cosmossdk-example/x/notes/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) CreateNote(ctx context.Context, msg *types.MsgCreateNote) (*types.MsgCreateNoteResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Creator); err != nil {
		return nil, errorsmod.Wrap(err, "invalid authority address")
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)

	// 1. Проверяем, что текст не пустой
	if len(msg.Text) == 0 {
		return nil, errors.New("text cannot be empty")
	}

	// 2. Создаём объект Note
	note := types.Note{
		Text:      msg.Text,
		Creator:   msg.Creator, // Автор берётся из подписи транзакции
		CreatedAt: sdkCtx.BlockTime().Unix(),
	}

	// 3. Сохраняем в хранилище
	k.Note.Set(ctx, msg.Creator, note)

	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent("note_created",
			sdk.NewAttribute("creator", msg.Creator),
			sdk.NewAttribute("text_length", strconv.Itoa(len(msg.Text))),
		),
	)

	return &types.MsgCreateNoteResponse{Id: msg.Creator}, nil
}
