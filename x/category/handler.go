package category

import (
	t "github.com/TruStory/truchain/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewHandler creates a new handler for all category messages
func NewHandler(k WriteKeeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case CreateCategoryMsg:
			return handleCreateCategoryMsg(ctx, k, msg)
		default:
			return t.ErrMsgHandler(msg)
		}
	}
}

// ============================================================================

func handleCreateCategoryMsg(ctx sdk.Context, k WriteKeeper, msg CreateCategoryMsg) sdk.Result {
	if err := msg.ValidateBasic(); err != nil {
		return err.Result()
	}

	id, err := k.NewCategory(ctx, msg.Title, msg.Slug, msg.Description)
	if err != nil {
		return err.Result()
	}

	return t.Result(id)
}