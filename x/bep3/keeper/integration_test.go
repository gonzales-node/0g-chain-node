package keeper_test

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	authexported "github.com/cosmos/cosmos-sdk/x/auth/exported"

	"github.com/tendermint/tendermint/crypto"
	tmtime "github.com/tendermint/tendermint/types/time"

	"github.com/kava-labs/kava/app"
	"github.com/kava-labs/kava/x/bep3/types"
)

const (
	TestSenderOtherChain    = "bnb1uky3me9ggqypmrsvxk7ur6hqkzq7zmv4ed4ng7"
	TestRecipientOtherChain = "bnb1urfermcg92dwq36572cx4xg84wpk3lfpksr5g7"
	TestDeputy              = "kava1xy7hrjy9r0algz9w3gzm8u6mrpq97kwta747gj"
)

var (
	DenomMap  = map[int]string{0: "btc", 1: "eth", 2: "bnb", 3: "xrp", 4: "dai"}
	TestUser1 = sdk.AccAddress(crypto.AddressHash([]byte("KavaTestUser1")))
	TestUser2 = sdk.AccAddress(crypto.AddressHash([]byte("KavaTestUser2")))
)

func i(in int64) sdk.Int                    { return sdk.NewInt(in) }
func c(denom string, amount int64) sdk.Coin { return sdk.NewInt64Coin(denom, amount) }
func cs(coins ...sdk.Coin) sdk.Coins        { return sdk.NewCoins(coins...) }
func ts(minOffset int) int64                { return tmtime.Now().Add(time.Duration(minOffset) * time.Minute).Unix() }

func NewAuthGenStateFromAccs(accounts ...authexported.GenesisAccount) app.GenesisState {
	authGenesis := auth.NewGenesisState(auth.DefaultParams(), accounts)
	return app.GenesisState{auth.ModuleName: auth.ModuleCdc.MustMarshalJSON(authGenesis)}
}

func NewBep3GenStateMulti(deputyAddress sdk.AccAddress) app.GenesisState {
	bep3Genesis := types.GenesisState{
		Params: types.Params{
			AssetParams: types.AssetParams{
				types.AssetParam{
					Denom:  "bnb",
					CoinID: 714,
					SupplyLimit: types.SupplyLimit{
						Limit:          sdk.NewInt(350000000000000),
						TimeLimited:    false,
						TimeBasedLimit: sdk.ZeroInt(),
						TimePeriod:     time.Hour,
					},
					Active:        true,
					DeputyAddress: deputyAddress,
					FixedFee:      sdk.NewInt(1000),
					MinSwapAmount: sdk.OneInt(),
					MaxSwapAmount: sdk.NewInt(1000000000000),
					MinBlockLock:  types.DefaultMinBlockLock,
					MaxBlockLock:  types.DefaultMaxBlockLock,
				},
				types.AssetParam{
					Denom:  "inc",
					CoinID: 9999,
					SupplyLimit: types.SupplyLimit{
						Limit:          sdk.NewInt(100000000000000),
						TimeLimited:    true,
						TimeBasedLimit: sdk.NewInt(50000000000),
						TimePeriod:     time.Hour,
					},
					Active:        false,
					DeputyAddress: deputyAddress,
					FixedFee:      sdk.NewInt(1000),
					MinSwapAmount: sdk.OneInt(),
					MaxSwapAmount: sdk.NewInt(100000000000),
					MinBlockLock:  types.DefaultMinBlockLock,
					MaxBlockLock:  types.DefaultMaxBlockLock,
				},
			},
		},
		Supplies: types.AssetSupplies{
			types.NewAssetSupply(
				sdk.NewCoin("bnb", sdk.ZeroInt()),
				sdk.NewCoin("bnb", sdk.ZeroInt()),
				sdk.NewCoin("bnb", sdk.ZeroInt()),
				sdk.NewCoin("bnb", sdk.ZeroInt()),
				time.Duration(0),
			),
			types.NewAssetSupply(
				sdk.NewCoin("inc", sdk.ZeroInt()),
				sdk.NewCoin("inc", sdk.ZeroInt()),
				sdk.NewCoin("inc", sdk.ZeroInt()),
				sdk.NewCoin("inc", sdk.ZeroInt()),
				time.Duration(0),
			),
		},
		PreviousBlockTime: types.DefaultPreviousBlockTime,
	}
	return app.GenesisState{types.ModuleName: types.ModuleCdc.MustMarshalJSON(bep3Genesis)}
}

func atomicSwaps(ctx sdk.Context, count int) types.AtomicSwaps {
	var swaps types.AtomicSwaps
	for i := 0; i < count; i++ {
		swap := atomicSwap(ctx, i)
		swaps = append(swaps, swap)
	}
	return swaps
}

func atomicSwap(ctx sdk.Context, index int) types.AtomicSwap {
	expireOffset := uint64(200) // Default expire height + offet to match timestamp
	timestamp := ts(index)      // One minute apart
	randomNumber, _ := types.GenerateSecureRandomNumber()
	randomNumberHash := types.CalculateRandomHash(randomNumber[:], timestamp)

	return types.NewAtomicSwap(cs(c("bnb", 50000)), randomNumberHash,
		uint64(ctx.BlockHeight())+expireOffset, timestamp, TestUser1, TestUser2,
		TestSenderOtherChain, TestRecipientOtherChain, 0, types.Open, true,
		types.Incoming)
}

func assetSupplies(count int) types.AssetSupplies {
	if count > 5 { // Max 5 asset supplies
		return types.AssetSupplies{}
	}

	var supplies types.AssetSupplies

	for i := 0; i < count; i++ {
		supply := assetSupply(DenomMap[i])
		supplies = append(supplies, supply)
	}
	return supplies
}

func assetSupply(denom string) types.AssetSupply {
	return types.NewAssetSupply(c(denom, 0), c(denom, 0), c(denom, 0), c(denom, 0), time.Duration(0))
}
