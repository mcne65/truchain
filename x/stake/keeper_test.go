package stake

import (
	"testing"
	"time"

	"github.com/TruStory/truchain/x/story"
	"github.com/TruStory/truchain/x/trubank"

	app "github.com/TruStory/truchain/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

type FakeStaker struct {
	*Vote `json:"vote"`
}

// ID implements `Voter.ID`
func (s FakeStaker) ID() int64 {
	return s.Vote.ID
}

// StoryID implements `Voter.StoryID`
func (s FakeStaker) StoryID() int64 {
	return s.Vote.StoryID
}

// Amount implements `Voter.Amount`
func (s FakeStaker) Amount() sdk.Coin {
	return s.Vote.Amount
}

// Creator implements `Voter.Creator`
func (s FakeStaker) Creator() sdk.AccAddress {
	return s.Vote.Creator
}

// VoteChoice implements `Voter.VoteChoice`
func (s FakeStaker) VoteChoice() bool {
	return s.Vote.Vote
}

// Timestamp implements `Voter.Timestamp`
func (s FakeStaker) Timestamp() app.Timestamp {
	return s.Vote.Timestamp
}

func TestRedistributeStakeTrueWins(t *testing.T) {
	ctx, k := mockDB()

	amount := sdk.NewCoin("trusteak", sdk.NewInt(10*app.Shanev))
	creator1 := sdk.AccAddress([]byte{1, 2})
	creator2 := sdk.AccAddress([]byte{2, 4})
	trueVote := Vote{
		ID:         1,
		StoryID:    1,
		Amount:     amount,
		ArgumentID: 0,
		Creator:    creator1,
		Vote:       true,
		Timestamp:  app.NewTimestamp(ctx.BlockHeader()),
	}
	backingVote := FakeStaker{&trueVote}

	falseVote := Vote{
		ID:         2,
		StoryID:    1,
		Amount:     amount,
		ArgumentID: 0,
		Creator:    creator2,
		Vote:       false,
		Timestamp:  app.NewTimestamp(ctx.BlockHeader()),
	}
	challengeVote := FakeStaker{&falseVote}

	votes := []Voter{backingVote, backingVote, backingVote, challengeVote}

	_, err := k.RedistributeStake(ctx, votes)
	assert.NoError(t, err)

	// check stake amounts
	transactions, err := k.truBankKeeper.TransactionsByCreator(ctx, creator1)
	assert.NoError(t, err)
	assert.Len(t, transactions, 3)
	assert.Equal(t, trubank.RewardPool, transactions[0].TransactionType)
	assert.Equal(t, trubank.RewardPool, transactions[1].TransactionType)
	assert.Equal(t, trubank.RewardPool, transactions[2].TransactionType)
	assert.Equal(t, "13333333333trusteak", transactions[0].Amount.String())
	assert.Equal(t, "13333333333trusteak", transactions[1].Amount.String())
	assert.Equal(t, "13333333333trusteak", transactions[2].Amount.String())

	// false voters get nothing back
	transactions, err = k.truBankKeeper.TransactionsByCreator(ctx, creator2)
	assert.NoError(t, err)
	assert.Len(t, transactions, 0)

}

func TestRedistributeStakeFalseWins(t *testing.T) {
	ctx, k := mockDB()

	amount := sdk.NewCoin("trusteak", sdk.NewInt(10*app.Shanev))
	creator1 := sdk.AccAddress([]byte{1, 2})
	creator2 := sdk.AccAddress([]byte{2, 4})
	trueVote := Vote{
		ID:         1,
		StoryID:    1,
		Amount:     amount,
		ArgumentID: 0,
		Creator:    creator1,
		Vote:       true,
		Timestamp:  app.NewTimestamp(ctx.BlockHeader()),
	}
	backingVote := FakeStaker{&trueVote}

	falseVote := Vote{
		ID:         2,
		StoryID:    1,
		Amount:     amount,
		ArgumentID: 0,
		Creator:    creator2,
		Vote:       false,
		Timestamp:  app.NewTimestamp(ctx.BlockHeader()),
	}
	challengeVote := FakeStaker{&falseVote}

	votes := []Voter{backingVote, challengeVote, challengeVote, challengeVote}

	_, err := k.RedistributeStake(ctx, votes)
	assert.NoError(t, err)

	transactions, err := k.truBankKeeper.TransactionsByCreator(ctx, creator1)
	assert.NoError(t, err)
	// true voters don't get anything back
	assert.Len(t, transactions, 0)

	transactions, err = k.truBankKeeper.TransactionsByCreator(ctx, creator2)
	assert.NoError(t, err)
	assert.Len(t, transactions, 3)
	assert.Equal(t, trubank.RewardPool, transactions[0].TransactionType)
	assert.Equal(t, trubank.RewardPool, transactions[1].TransactionType)
	assert.Equal(t, trubank.RewardPool, transactions[2].TransactionType)
	assert.Equal(t, "13333333333trusteak", transactions[0].Amount.String())
	assert.Equal(t, "13333333333trusteak", transactions[1].Amount.String())
	assert.Equal(t, "13333333333trusteak", transactions[2].Amount.String())
}

func TestRedistributeStakeNoMajority(t *testing.T) {
	ctx, k := mockDB()

	amount := sdk.NewCoin("trusteak", sdk.NewInt(10*app.Shanev))
	creator1 := sdk.AccAddress([]byte{1, 2})
	creator2 := sdk.AccAddress([]byte{2, 4})
	trueVote := Vote{
		ID:         1,
		StoryID:    1,
		Amount:     amount,
		ArgumentID: 0,
		Creator:    creator1,
		Vote:       true,
		Timestamp:  app.NewTimestamp(ctx.BlockHeader()),
	}
	backingVote := FakeStaker{&trueVote}

	falseVote := Vote{
		ID:         2,
		StoryID:    1,
		Amount:     amount,
		ArgumentID: 0,
		Creator:    creator2,
		Vote:       false,
		Timestamp:  app.NewTimestamp(ctx.BlockHeader()),
	}
	challengeVote := FakeStaker{&falseVote}

	votes := []Voter{backingVote, challengeVote}

	_, err := k.RedistributeStake(ctx, votes)
	assert.NoError(t, err)

	transactions, err := k.truBankKeeper.TransactionsByCreator(ctx, creator1)
	assert.NoError(t, err)
	assert.Len(t, transactions, 1)
	assert.Equal(t, trubank.BackingReturned, transactions[0].TransactionType)
	assert.Equal(t, "10000000000trusteak", transactions[0].Amount.String())

	transactions, err = k.truBankKeeper.TransactionsByCreator(ctx, creator2)
	assert.NoError(t, err)
	assert.Len(t, transactions, 1)
	assert.Equal(t, trubank.ChallengeReturned, transactions[0].TransactionType)
	assert.Equal(t, "10000000000trusteak", transactions[0].Amount.String())
}

func TestDistributeInterest(t *testing.T) {
	ctx, k := mockDB()

	amount := sdk.NewCoin("trusteak", sdk.NewInt(10*app.Shanev))
	creator1 := sdk.AccAddress([]byte{1, 2})
	creator2 := sdk.AccAddress([]byte{2, 4})
	trueVote := Vote{
		ID:         1,
		StoryID:    1,
		Amount:     amount,
		ArgumentID: 0,
		Creator:    creator1,
		Vote:       true,
		Timestamp:  app.NewTimestamp(ctx.BlockHeader()),
	}
	backingVote := FakeStaker{&trueVote}

	falseVote := Vote{
		ID:         2,
		StoryID:    1,
		Amount:     amount,
		ArgumentID: 0,
		Creator:    creator2,
		Vote:       false,
		Timestamp:  app.NewTimestamp(ctx.BlockHeader()),
	}
	challengeVote := FakeStaker{&falseVote}

	votes := []Voter{backingVote, backingVote, backingVote, challengeVote}

	_, err := k.DistributeInterest(ctx, votes)
	assert.NoError(t, err)

	transactions, err := k.truBankKeeper.TransactionsByCreator(ctx, creator1)
	assert.NoError(t, err)
	assert.Len(t, transactions, 3)
	assert.Equal(t, trubank.Interest, transactions[0].TransactionType)
	assert.Equal(t, trubank.Interest, transactions[1].TransactionType)
	assert.Equal(t, trubank.Interest, transactions[2].TransactionType)
	assert.Equal(t, "3330trusteak", transactions[0].Amount.String())
	assert.Equal(t, "3330trusteak", transactions[1].Amount.String())
	assert.Equal(t, "3330trusteak", transactions[2].Amount.String())

	transactions, err = k.truBankKeeper.TransactionsByCreator(ctx, creator2)
	assert.NoError(t, err)
	assert.Len(t, transactions, 1)
	assert.Equal(t, trubank.Interest, transactions[0].TransactionType)
	assert.Equal(t, "3330trusteak", transactions[0].Amount.String())
}

func Test_interest_MidAmountMidPeriod(t *testing.T) {
	ctx, k := mockDB()

	amount := sdk.NewCoin("crypto", sdk.NewInt(500000000000000))
	period := 12 * time.Hour

	interest := k.interest(ctx, amount, period)
	assert.Equal(t, sdk.NewInt(25000000000000).String(), interest.String())
}

func Test_interest_MaxAmountMinPeriod(t *testing.T) {
	ctx, k := mockDB()
	amount := sdk.NewCoin("crypto", sdk.NewInt(1000000000000000))
	period := 0 * time.Hour

	interest := k.interest(ctx, amount, period)
	assert.Equal(t, sdk.NewInt(33300000000000).String(), interest.String())
}

func Test_interest_MinAmountMaxPeriod(t *testing.T) {
	ctx, k := mockDB()
	amount := sdk.NewCoin("crypto", sdk.NewInt(0))
	period := 24 * time.Hour

	interest := k.interest(ctx, amount, period)
	assert.Equal(t, interest.String(), sdk.NewInt(0).String())
}

func Test_interest_MaxAmountMaxPeriod(t *testing.T) {
	ctx, k := mockDB()
	amount := sdk.NewCoin("crypto", sdk.NewInt(1000000000000000))
	period := 24 * time.Hour
	maxInterestRate := k.GetParams(ctx).MaxInterestRate
	expected := sdk.NewDecFromInt(amount.Amount).Mul(maxInterestRate)

	interest := k.interest(ctx, amount, period)
	assert.Equal(t, expected.RoundInt().String(), interest.String())
}

func Test_interest_ZeroAmountWeight(t *testing.T) {
	ctx, k := mockDB()
	params := Params{
		MaxAmount:       k.GetParams(ctx).MaxAmount,
		MinInterestRate: k.GetParams(ctx).MinInterestRate,
		MaxInterestRate: k.GetParams(ctx).MaxInterestRate,
		MajorityPercent: k.GetParams(ctx).MajorityPercent,
		AmountWeight:    sdk.ZeroDec(),
		PeriodWeight:    sdk.NewDec(1),
	}
	k.SetParams(ctx, params)

	storyParams := story.Params{
		ExpireDuration: 72 * time.Hour,
		MinStoryLength: 25,
		MaxStoryLength: 350,
	}
	storyWriteKeeper := k.storyKeeper.(story.WriteKeeper)
	storyWriteKeeper.SetParams(ctx, storyParams)

	// 10000000000 amount
	//           0 interest
	amount := sdk.NewCoin("crypto", sdk.NewInt(10*app.Shanev))
	period := 0 * time.Hour
	interest := k.interest(ctx, amount, period)
	assert.Equal(t, sdk.NewInt(0).String(), interest.String())

	// 10000000000 amount
	//   333333333 interest
	amount = sdk.NewCoin("crypto", sdk.NewInt(10*app.Shanev))
	period = 24 * time.Hour
	interest = k.interest(ctx, amount, period)
	assert.Equal(t, sdk.NewInt(333333333).String(), interest.String())

	// 10000000000 amount
	//   666666667 interest
	amount = sdk.NewCoin("crypto", sdk.NewInt(10*app.Shanev))
	period = 48 * time.Hour
	interest = k.interest(ctx, amount, period)
	assert.Equal(t, sdk.NewInt(666666667).String(), interest.String())

	// 10000000000 amount
	//  1000000000 interest
	amount = sdk.NewCoin("crypto", sdk.NewInt(10*app.Shanev))
	period = 72 * time.Hour
	interest = k.interest(ctx, amount, period)
	assert.Equal(t, sdk.NewInt(1000000000).String(), interest.String())
}
