package sscqservice

import (
	"github.com/deep2chain/htdf/codec"
)

// RegisterCodec registers concrete types on the Amino codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgSend{}, "sscqservice/send", nil)
	// cdc.RegisterConcrete(MsgAdd{}, "sscqservice/add", nil)
}
