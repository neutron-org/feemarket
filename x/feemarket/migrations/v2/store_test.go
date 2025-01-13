package v2_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"cosmossdk.io/math"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/testutil"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"

	"github.com/skip-mev/feemarket/x/feemarket"
	v2 "github.com/skip-mev/feemarket/x/feemarket/migrations/v2"
	"github.com/skip-mev/feemarket/x/feemarket/types"
)

func TestParamsUpgrade(t *testing.T) {
	var (
		encCfg = moduletestutil.MakeTestEncodingConfig(feemarket.AppModuleBasic{})
		cdc    = encCfg.Codec

		storeKey = storetypes.NewKVStoreKey(types.StoreKey)
		tKey     = storetypes.NewTransientStoreKey("transient_test")
		ctx      = testutil.DefaultContext(storeKey, tKey)
	)

	// Write old params
	oldParams := types.Params{
		Alpha:               math.LegacyMustNewDecFromStr("0.0"),
		Beta:                math.LegacyMustNewDecFromStr("1.0"),
		Gamma:               math.LegacyMustNewDecFromStr("0.0"),
		Delta:               math.LegacyMustNewDecFromStr("0.0"),
		MinBaseGasPrice:     math.LegacyOneDec(),
		MinLearningRate:     math.LegacyMustNewDecFromStr("0.125"),
		MaxLearningRate:     math.LegacyMustNewDecFromStr("0.125"),
		MaxBlockUtilization: 30_000_000,
		Window:              1,
		FeeDenom:            types.DefaultFeeDenom,
		Enabled:             true,
		DistributeFees:      true,
		SendTipToProposer:   false,
	}

	store := ctx.KVStore(storeKey)
	bz, err := cdc.Marshal(&oldParams)
	require.NoError(t, err)

	store.Set(types.KeyParams, bz)

	// Run migration
	require.NoError(t, v2.MigrateStore(ctx, cdc, storeKey))

	bz = store.Get(types.KeyParams)
	require.NotNil(t, bz)

	var newParams types.Params
	cdc.MustUnmarshal(bz, &newParams)

	// Check params are correct
	require.EqualValues(t, oldParams.Alpha, newParams.Alpha)
	require.EqualValues(t, oldParams.Beta, newParams.Beta)
	require.EqualValues(t, oldParams.Gamma, newParams.Gamma)
	require.EqualValues(t, oldParams.Delta, newParams.Delta)
	require.EqualValues(t, oldParams.MinBaseGasPrice, newParams.MinBaseGasPrice)
	require.EqualValues(t, oldParams.MinLearningRate, newParams.MinLearningRate)
	require.EqualValues(t, oldParams.MaxLearningRate, newParams.MaxLearningRate)
	require.EqualValues(t, oldParams.MaxBlockUtilization, newParams.MaxBlockUtilization)
	require.EqualValues(t, oldParams.Window, newParams.Window)
	require.EqualValues(t, oldParams.FeeDenom, newParams.FeeDenom)
	require.EqualValues(t, oldParams.Enabled, newParams.Enabled)
	require.EqualValues(t, oldParams.DistributeFees, newParams.DistributeFees)
	require.EqualValues(t, true, newParams.SendTipToProposer)
}
