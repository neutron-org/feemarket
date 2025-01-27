package v2

import (
	"errors"

	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/skip-mev/feemarket/x/feemarket/types"
)

// MigrateStore performs in-place store migrations.
// The migration adds new feemarket param -- SendTipToProposer.
func MigrateStore(ctx sdk.Context, cdc codec.BinaryCodec, storeKey storetypes.StoreKey) error {
	if err := migrateParams(ctx, cdc, storeKey); err != nil {
		return err
	}

	return nil
}

func migrateParams(ctx sdk.Context, cdc codec.BinaryCodec, storeKey storetypes.StoreKey) error {
	ctx.Logger().Info("Migrating feemarket params...")

	// fetch old params
	store := ctx.KVStore(storeKey)
	bz := store.Get(types.KeyParams)
	if bz == nil {
		return errors.New("cannot fetch feemarket params from KV store")
	}
	var oldParams types.Params
	cdc.MustUnmarshal(bz, &oldParams)

	// add new param values
	newParams := types.Params{
		Alpha:               oldParams.Alpha,
		Beta:                oldParams.Beta,
		Gamma:               oldParams.Gamma,
		Delta:               oldParams.Delta,
		MinBaseGasPrice:     oldParams.MinBaseGasPrice,
		MinLearningRate:     oldParams.MinLearningRate,
		MaxLearningRate:     oldParams.MaxLearningRate,
		MaxBlockUtilization: oldParams.MaxBlockUtilization,
		Window:              oldParams.Window,
		FeeDenom:            oldParams.FeeDenom,
		Enabled:             oldParams.Enabled,
		DistributeFees:      oldParams.DistributeFees,
		SendTipToProposer:   true,
	}

	// set params
	bz, err := cdc.Marshal(&newParams)
	if err != nil {
		return err
	}
	store.Set(types.KeyParams, bz)

	ctx.Logger().Info("Finished migrating feemarket params")

	return nil
}
