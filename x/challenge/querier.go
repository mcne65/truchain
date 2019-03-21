package challenge

import (
	app "github.com/TruStory/truchain/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// query endpoints supported by the truchain Querier
const (
	QueryPath                     = "challenges"
	QueryByStoryID                = "storyID"
	QueryByStoryIDAndCreator      = "storyIDAndCreator"
	QueryChallengeAmountByStoryID = "totalAmountByStoryID"
)

// NewQuerier returns a function that handles queries on the KVStore
func NewQuerier(k ReadKeeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case QueryByStoryID:
			return queryByStoryID(ctx, req, k)
		case QueryByStoryIDAndCreator:
			return queryByStoryIDAndCreator(ctx, req, k)
		case QueryChallengeAmountByStoryID:
			return queryChallengeAmountByStoryID(ctx, req, k)
		default:
			return nil, sdk.ErrUnknownRequest("Unknown query endpoint")
		}
	}
}

// ============================================================================

func queryByStoryID(
	ctx sdk.Context,
	req abci.RequestQuery,
	k ReadKeeper) (res []byte, sdkErr sdk.Error) {

	params := app.QueryByIDParams{}

	sdkErr = app.UnmarshalQueryParams(req, &params)
	if sdkErr != nil {
		return
	}

	challenges, sdkErr := k.ChallengesByStoryID(ctx, params.ID)
	if sdkErr != nil {
		return
	}

	return app.MustMarshal(challenges), nil
}

func queryChallengeAmountByStoryID(ctx sdk.Context, req abci.RequestQuery, k ReadKeeper) (res []byte, err sdk.Error) {
	params := app.QueryByIDParams{}

	if err = app.UnmarshalQueryParams(req, &params); err != nil {
		return
	}

	challengePool, err := k.TotalChallengeAmount(ctx, params.ID)
	if err != nil {
		return
	}

	return app.MustMarshal(challengePool), nil
}

func queryByStoryIDAndCreator(
	ctx sdk.Context,
	req abci.RequestQuery,
	k ReadKeeper) (res []byte, sdkErr sdk.Error) {

	params := app.QueryByStoryIDAndCreatorParams{}

	sdkErr = app.UnmarshalQueryParams(req, &params)
	if sdkErr != nil {
		return
	}

	// convert address bech32 string to bytes
	addr, err := sdk.AccAddressFromBech32(params.Creator)
	if err != nil {
		return res, sdk.ErrInvalidAddress("Cannot decode address")
	}

	challenge, sdkErr := k.ChallengeByStoryIDAndCreator(ctx, params.StoryID, addr)
	if sdkErr != nil {
		return
	}

	return app.MustMarshal(challenge), nil
}
