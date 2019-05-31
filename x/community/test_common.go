package community

import (
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	amino "github.com/tendermint/go-amino"
	abci "github.com/tendermint/tendermint/abci/types"
	cryptoAmino "github.com/tendermint/tendermint/crypto/encoding/amino"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
)

func mockDB() (sdk.Context, Keeper) {
	db := dbm.NewMemDB()

	communityKey := sdk.NewKVStoreKey(StoreKey)

	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(communityKey, sdk.StoreTypeIAVL, db)
	ms.LoadLatestVersion()

	ctx := sdk.NewContext(ms, abci.Header{}, false, log.NewNopLogger())

	cdc := amino.NewCodec()
	cryptoAmino.RegisterAmino(cdc)
	RegisterAmino(cdc)

	ck := NewKeeper(communityKey, cdc)

	return ctx, ck
}

func getFakeCommunityParams() (name string, slug string, description string) {
	name, slug, description = "Randomness", "randomness", "All the random quantum things happen in this community."

	return
}
