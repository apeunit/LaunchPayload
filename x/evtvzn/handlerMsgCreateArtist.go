package evtvzn

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/apeunit/evtvzn/x/evtvzn/types"
	"github.com/apeunit/evtvzn/x/evtvzn/keeper"
)

func handleMsgCreateArtist(ctx sdk.Context, k keeper.Keeper, msg types.MsgCreateArtist) (*sdk.Result, error) {
	var artist = types.Artist{
		Creator: msg.Creator,
		ID:      msg.ID,
    Name: msg.Name,
    Dropaddress: msg.Dropaddress,
	}
	k.CreateArtist(ctx, artist)

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}
