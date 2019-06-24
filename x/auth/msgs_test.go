package auth

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

func TestMsgRegisterKey_Success(t *testing.T) {
	ctx, keeper := mockDB()

	_, publicKey, address, coins, _ := getFakeAppAccountParams()

	registrar := keeper.GetParams(ctx).Registrar

	msg := NewMsgRegisterKey(registrar, address, publicKey, "secp256k1", coins)
	err := msg.ValidateBasic()
	assert.Nil(t, err)
	assert.Equal(t, ModuleName, msg.Route())
	assert.Equal(t, TypeMsgRegisterKey, msg.Type())
}

func TestMsgNewCommunity_InvalidAddress(t *testing.T) {
	ctx, keeper := mockDB()

	_, publicKey, _, coins, _ := getFakeAppAccountParams()
	invalidAddress := sdk.AccAddress(nil)

	registrar := keeper.GetParams(ctx).Registrar

	msg := NewMsgRegisterKey(registrar, invalidAddress, publicKey, "secp256k1", coins)
	err := msg.ValidateBasic()
	assert.NotNil(t, err)
	assert.Equal(t, sdk.ErrInvalidAddress("").Code(), err.Code())
}