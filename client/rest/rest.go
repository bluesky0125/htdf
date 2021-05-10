package rest

import (
	"log"
	"net/http"

	"github.com/deep2chain/htdf/client"
	"github.com/deep2chain/htdf/client/context"
	"github.com/deep2chain/htdf/client/utils"
	"github.com/deep2chain/htdf/codec"
	sdk "github.com/deep2chain/htdf/types"
	"github.com/deep2chain/htdf/types/rest"
	"github.com/deep2chain/htdf/x/auth"
	authtxb "github.com/deep2chain/htdf/x/auth/client/txbuilder"
)

//-----------------------------------------------------------------------------
// Building / Sending utilities

// WriteGenerateStdTxResponse writes response for the generate only mode.
func WriteGenerateStdTxResponse(w http.ResponseWriter, cdc *codec.Codec,
	cliCtx context.CLIContext, br rest.BaseReq, msgs []sdk.Msg) {

	gasAdj, ok := rest.ParseFloat64OrReturnBadRequest(w, br.GasAdjustment, client.DefaultGasAdjustment)
	if !ok {
		return
	}

	simAndExec, gasWanted, err := client.ParseGas(br.GasWanted)
	if err != nil {
		rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	var gasPrice uint64
	gasPrice, err = client.ParseGasPrice(br.GasPrice)
	if err != nil {
		rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
	}
	txBldr := authtxb.NewTxBuilder(
		utils.GetTxEncoder(cdc), br.AccountNumber, br.Sequence, gasWanted, gasAdj,
		br.Simulate, br.ChainID, br.Memo, gasPrice,
	)

	if br.Simulate || simAndExec {
		if gasAdj < 0 {
			rest.WriteErrorResponse(w, http.StatusBadRequest, client.ErrInvalidGasAdjustment.Error())
			return
		}

		txBldr, err = utils.EnrichWithGas(txBldr, cliCtx, msgs)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		if br.Simulate {
			rest.WriteSimulationResponse(w, cdc, txBldr.GasWanted())
			return
		}
	}

	stdMsg, err := txBldr.BuildSignMsg(msgs)
	if err != nil {
		rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	output, err := cdc.MarshalJSON(auth.NewStdTx(stdMsg.Msgs, stdMsg.Fee, nil, stdMsg.Memo))
	if err != nil {
		rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(output); err != nil {
		log.Printf("could not write response: %v", err)
	}
	return
}
