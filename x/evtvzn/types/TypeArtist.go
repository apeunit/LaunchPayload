package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Artist struct {
	Creator sdk.AccAddress `json:"creator" yaml:"creator"`
	ID      string         `json:"id" yaml:"id"`
  Name string `json:"name" yaml:"name"`
  Dropaddress string `json:"dropaddress" yaml:"dropaddress"`
}