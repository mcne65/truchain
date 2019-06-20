package trubank

import (
	"encoding/json"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// ModuleName is the name of this module
const ModuleName = "trubank"

// Defines module constants
const (
	RouterKey    = ModuleName
	QuerierRoute = ModuleName
)

var (
	_ sdk.AppModule      = AppModule{}
	_ sdk.AppModuleBasic = AppModuleBasic{}
)

// AppModuleBasic defines the internal data for the module
// ----------------------------------------------------------------------------
type AppModuleBasic struct{}

var _ sdk.AppModuleBasic = AppModuleBasic{}

// Name define the name of the module
func (AppModuleBasic) Name() string {
	return ModuleName
}

// RegisterCodec registers the types needed for amino encoding/decoding
func (AppModuleBasic) RegisterCodec(cdc *codec.Codec) {
	RegisterCodec(cdc)
}

// DefaultGenesis creates the default genesis state for testing
func (AppModuleBasic) DefaultGenesis() json.RawMessage {
	return moduleCodec.MustMarshalJSON(DefaultGenesisState())
}

// ValidateGenesis validates the genesis state
func (AppModuleBasic) ValidateGenesis(bz json.RawMessage) error {
	var data GenesisState
	err := moduleCodec.UnmarshalJSON(bz, &data)
	if err != nil {
		return err
	}
	return ValidateGenesis(data)
}

type AppModule struct {
	AppModuleBasic
	keeper Keeper
}

// NewAppModule creates a new AppModule object
func NewAppModule(keeper Keeper) AppModule {
	return AppModule{
		AppModuleBasic: AppModuleBasic{},
		keeper:         keeper,
	}
}

// RegisterInvariants enforces registering of invariants
func (AppModule) RegisterInvariants(_ sdk.InvariantRouter) {}

// Route defines the key for the route
func (AppModule) Route() string {
	return RouterKey
}

// NewHandler creates the handler for the module
func (am AppModule) NewHandler() sdk.Handler {
	return NewHandler(am.keeper)
}

// QuerierRoute defines the querier route
func (AppModule) QuerierRoute() string {
	return QuerierRoute
}

// NewQuerierHandler creates a new querier handler
func (am AppModule) NewQuerierHandler() sdk.Querier {
	return NewQuerier(am.keeper)
}

// InitGenesis enforces the creation of the genesis state for this module
func (am AppModule) InitGenesis(ctx sdk.Context, data json.RawMessage) []abci.ValidatorUpdate {
	var genesisState GenesisState
	moduleCodec.MustUnmarshalJSON(data, &genesisState)
	InitGenesis(ctx, am.keeper, genesisState)
	return []abci.ValidatorUpdate{}
}

// ExportGenesis enforces exporting this module's data to a genesis file
func (am AppModule) ExportGenesis(ctx sdk.Context) json.RawMessage {
	gs := ExportGenesis(ctx, am.keeper)
	return moduleCodec.MustMarshalJSON(gs)
}

// BeginBlock runs before a block is processed
func (AppModule) BeginBlock(_ sdk.Context, _ abci.RequestBeginBlock) sdk.Tags {
	return sdk.EmptyTags()
}

// EndBlock runs at the end of each block
func (AppModule) EndBlock(_ sdk.Context, _ abci.RequestEndBlock) ([]abci.ValidatorUpdate, sdk.Tags) {
	return []abci.ValidatorUpdate{}, sdk.EmptyTags()
}