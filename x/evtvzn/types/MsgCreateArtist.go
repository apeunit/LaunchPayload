package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/google/uuid"
)

var _ sdk.Msg = &MsgCreateArtist{}

type MsgCreateArtist struct {
  ID      string
  Creator sdk.AccAddress `json:"creator" yaml:"creator"`
  Name string `json:"name" yaml:"name"`
  Dropaddress string `json:"dropaddress" yaml:"dropaddress"`
}

func NewMsgCreateArtist(creator sdk.AccAddress, name string, dropaddress string) MsgCreateArtist {
  return MsgCreateArtist{
    ID: uuid.New().String(),
		Creator: creator,
    Name: name,
    Dropaddress: dropaddress,
	}
}

func (msg MsgCreateArtist) Route() string {
  return RouterKey
}

func (msg MsgCreateArtist) Type() string {
  return "CreateArtist"
}

func (msg MsgCreateArtist) GetSigners() []sdk.AccAddress {
  return []sdk.AccAddress{sdk.AccAddress(msg.Creator)}
}

func (msg MsgCreateArtist) GetSignBytes() []byte {
  bz := ModuleCdc.MustMarshalJSON(msg)
  return sdk.MustSortJSON(bz)
}

func (msg MsgCreateArtist) ValidateBasic() error {
  if msg.Creator.Empty() {
    return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "creator can't be empty")
  }
  return nil
}