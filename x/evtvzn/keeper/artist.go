package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/apeunit/evtvzn/x/evtvzn/types"
  "github.com/cosmos/cosmos-sdk/codec"
)

func (k Keeper) CreateArtist(ctx sdk.Context, artist types.Artist) {
	store := ctx.KVStore(k.storeKey)
	key := []byte(types.ArtistPrefix + artist.ID)
	value := k.cdc.MustMarshalBinaryLengthPrefixed(artist)
	store.Set(key, value)
}

func listArtist(ctx sdk.Context, k Keeper) ([]byte, error) {
  var artistList []types.Artist
  store := ctx.KVStore(k.storeKey)
  iterator := sdk.KVStorePrefixIterator(store, []byte(types.ArtistPrefix))
  for ; iterator.Valid(); iterator.Next() {
    var artist types.Artist
    k.cdc.MustUnmarshalBinaryLengthPrefixed(store.Get(iterator.Key()), &artist)
    artistList = append(artistList, artist)
  }
  res := codec.MustMarshalJSONIndent(k.cdc, artistList)
  return res, nil
}