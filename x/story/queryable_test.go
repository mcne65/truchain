package story

import (
	"encoding/json"
	"testing"

	app "github.com/TruStory/truchain/types"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
)

func TestQueryStories_ErrNotFound(t *testing.T) {
	ctx, k, _ := mockDB()

	queryParams := QueryCategoryStoriesParams{
		CategoryID: 1,
	}

	bz, errRes := json.Marshal(queryParams)
	require.Nil(t, errRes)

	query := abci.RequestQuery{
		Path: "/custom/stories/category",
		Data: bz,
	}

	_, err := queryStoriesByCategoryID(ctx, query, k)
	require.NotNil(t, err)
	require.Equal(t, ErrStoriesWithCategoryNotFound(1).Code(), err.Code(), "should get error")
}

func TestQueryStoryByID(t *testing.T) {
	ctx, sk, ck := mockDB()

	createFakeStory(ctx, sk, ck)

	queryParams := QueryStoryByIDParams{
		ID: 1,
	}

	bz, errRes := json.Marshal(queryParams)
	require.Nil(t, errRes)

	query := abci.RequestQuery{
		Path: "/custom/stories/id",
		Data: bz,
	}
	_, err := queryStoryByID(ctx, query, sk)

	require.Nil(t, err)
}

func TestQueryStoriesWithCategory(t *testing.T) {
	ctx, sk, ck := mockDB()

	createFakeStory(ctx, sk, ck)

	queryParams := QueryCategoryStoriesParams{
		CategoryID: 1,
	}

	bz, errRes := json.Marshal(queryParams)
	require.Nil(t, errRes)

	query := abci.RequestQuery{
		Path: "/custom/stories/category",
		Data: bz,
	}
	_, err := queryStoriesByCategoryID(ctx, query, sk)

	require.Nil(t, err)
}

func TestQueryStoriesByFeedID(t *testing.T) {
	ctx, sk, ck := mockDB()

	createFakeStory(ctx, sk, ck)

	queryParams := app.QueryByIDParams{
		ID: 1,
	}

	bz, errRes := json.Marshal(queryParams)
	require.Nil(t, errRes)

	query := abci.RequestQuery{
		Path: "/custom/stories/feedID",
		Data: bz,
	}

	_, err := queryStoriesByFeedID(ctx, query, sk)
	require.Nil(t, err)
}
