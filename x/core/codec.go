package htdfservice

import (
	"github.com/deep2chain/htdf/codec"
)

// RegisterCodec registers concrete types on the Amino codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgSend{}, "htdfservice/send", nil)
	// cdc.RegisterConcrete(MsgAdd{}, "htdfservice/add", nil)
}
