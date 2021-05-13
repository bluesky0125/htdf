package main

import (
	"fmt"
	"os"
	"path"

	"github.com/deep2chain/htdf/client/bech32"

	"github.com/deep2chain/htdf/params"
	svrConfig "github.com/deep2chain/htdf/server/config"

	"github.com/deep2chain/htdf/client"
	"github.com/deep2chain/htdf/client/lcd"
	"github.com/deep2chain/htdf/client/rpc"
	"github.com/deep2chain/htdf/client/tx"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/libs/cli"

	sdk "github.com/deep2chain/htdf/types"
	authcmd "github.com/deep2chain/htdf/x/auth/client/cli"
	sscqservicecmd "github.com/deep2chain/htdf/x/core/client/cli"

	accounts "github.com/deep2chain/htdf/accounts/cli"
	accrest "github.com/deep2chain/htdf/accounts/rest"
	"github.com/deep2chain/htdf/app"
	hsrest "github.com/deep2chain/htdf/x/core/client/rest"

	dist "github.com/deep2chain/htdf/x/distribution/client/rest"
	gv "github.com/deep2chain/htdf/x/gov"
	gov "github.com/deep2chain/htdf/x/gov/client/rest"
	mint "github.com/deep2chain/htdf/x/mint/client/rest"
	sl "github.com/deep2chain/htdf/x/slashing"
	slashing "github.com/deep2chain/htdf/x/slashing/client/rest"
	st "github.com/deep2chain/htdf/x/staking"
	staking "github.com/deep2chain/htdf/x/staking/client/rest"

	hscliversion "github.com/deep2chain/htdf/server"
	distcmd "github.com/deep2chain/htdf/x/distribution"
	hsdistClient "github.com/deep2chain/htdf/x/distribution/client"
	hsgovClient "github.com/deep2chain/htdf/x/gov/client"
	hsmintClient "github.com/deep2chain/htdf/x/mint/client/cli"
	hslashingClient "github.com/deep2chain/htdf/x/slashing/client"
	hstakingClient "github.com/deep2chain/htdf/x/staking/client"
	upgradecmd "github.com/deep2chain/htdf/x/upgrade/client/cli"
	upgraderest "github.com/deep2chain/htdf/x/upgrade/client/rest"
)

const (
	storeAcc = "acc"
	storeHS  = "hs"
)

var (
	DEBUGAPI  = "OFF"
	GitCommit = ""
	GitBranch = ""
)

func main() {
	cobra.EnableCommandSorting = false

	if DEBUGAPI == svrConfig.ValueDebugApi_On {
		svrConfig.ApiSecurityLevel = svrConfig.ValueSecurityLevel_Low
	}

	cdc := app.MakeLatestCodec()

	// set address prefix
	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount(sdk.Bech32PrefixAccAddr, sdk.Bech32PrefixAccPub)
	config.SetBech32PrefixForValidator(sdk.Bech32PrefixValAddr, sdk.Bech32PrefixValPub)
	config.SetBech32PrefixForConsensusNode(sdk.Bech32PrefixConsAddr, sdk.Bech32PrefixConsPub)
	config.Seal()

	mc := []sdk.ModuleClients{
		hsgovClient.NewModuleClient(gv.StoreKey, cdc),
		hsdistClient.NewModuleClient(distcmd.StoreKey, cdc),
		hstakingClient.NewModuleClient(st.StoreKey, cdc),
		hslashingClient.NewModuleClient(sl.StoreKey, cdc),
	}

	rootCmd := &cobra.Command{
		Use:   "hscli",
		Short: "sscqservice Client",
	}

	// Add --chain-id to persistent flags and mark it required
	rootCmd.PersistentFlags().String(client.FlagChainID, "", "Chain ID of tendermint node")
	rootCmd.PersistentPreRunE = func(_ *cobra.Command, _ []string) error {
		return initConfig(rootCmd)
	}

	// Construct Root Command
	rootCmd.AddCommand(
		rpc.StatusCommand(),
		client.ConfigCmd(app.DefaultCLIHome),
		queryCmd(cdc, mc), // check the below
		txCmd(cdc, mc),    // check the below
		versionCmd(cdc, mc),
		client.LineBreak,
		lcd.ServeCommand(cdc, registerRoutes),
		client.LineBreak,
		accounts.Commands(),
		client.LineBreak,
		hscliversion.VersionHscliCmd,
		bech32.Bech32Commands(),
	)

	executor := cli.PrepareMainCmd(rootCmd, "HS", app.DefaultCLIHome)
	err := executor.Execute()
	if err != nil {
		panic(err)
	}
}

func registerRoutes(rs *lcd.RestServer) {
	rs.CliCtx = rs.CliCtx.WithAccountDecoder(rs.Cdc)
	rpc.RegisterRoutes(rs.CliCtx, rs.Mux, rs.Cdc)
	tx.RegisterRoutes(rs.CliCtx, rs.Mux, rs.Cdc)
	hsrest.RegisterRoutes(rs.CliCtx, rs.Mux, rs.Cdc, storeHS)
	accrest.RegisterRoutes(rs.CliCtx, rs.Mux, rs.Cdc)
	accrest.RegisterRoute(rs.CliCtx, rs.Mux, rs.Cdc, storeAcc)
	dist.RegisterRoutes(rs.CliCtx, rs.Mux, rs.Cdc, distcmd.StoreKey)
	staking.RegisterRoutes(rs.CliCtx, rs.Mux, rs.Cdc, rs.KeyBase)
	slashing.RegisterRoutes(rs.CliCtx, rs.Mux, rs.Cdc, rs.KeyBase)
	gov.RegisterRoutes(rs.CliCtx, rs.Mux, rs.Cdc)
	mint.RegisterRoutes(rs.CliCtx, rs.Mux, rs.Cdc)
	upgraderest.RegisterRoutes(rs.CliCtx, rs.Mux, rs.Cdc)
}

func versionCmd(cdc *amino.Codec, mc []sdk.ModuleClients) *cobra.Command {
	cbCmd := &cobra.Command{
		Use:   "version",
		Short: "print version, api security level",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("GitCommit=%s|version=%s|GitBranch=%s|DEBUGAPI=%s|ApiSecurityLevel=%s\n", GitCommit, params.Version, GitBranch, DEBUGAPI, svrConfig.ApiSecurityLevel)
		},
	}

	return cbCmd
}

func queryCmd(cdc *amino.Codec, mc []sdk.ModuleClients) *cobra.Command {
	queryCmd := &cobra.Command{
		Use:     "query",
		Aliases: []string{"q"},
		Short:   "Querying subcommands",
	}

	queryCmd.AddCommand(
		rpc.ValidatorCommand(cdc),
		rpc.BlockCommand(),
		tx.SearchTxCmd(cdc),
		tx.QueryTxCmd(cdc),
		client.LineBreak,
		authcmd.GetAccountCmd(storeAcc, cdc),
		sscqservicecmd.GetCmdCall(cdc),
		hsmintClient.GetCmdQueryBlockRewards(cdc),
		hsmintClient.GetCmdQueryTotalProvisions(cdc),
		upgradecmd.GetInfoCmd("upgrade", cdc),
		upgradecmd.GetCmdQuerySignals("upgrade", cdc),
	)

	for _, m := range mc {
		queryCmd.AddCommand(m.GetQueryCmd())
	}

	return queryCmd
}

func txCmd(cdc *amino.Codec, mc []sdk.ModuleClients) *cobra.Command {
	txCmd := &cobra.Command{
		Use:   "tx",
		Short: "Transactions subcommands",
	}

	if svrConfig.ApiSecurityLevel == svrConfig.ValueSecurityLevel_Low {
		txCmd.AddCommand(
			sscqservicecmd.GetCmdBurn(cdc),
			sscqservicecmd.GetCmdCreate(cdc),
			sscqservicecmd.GetCmdSend(cdc),
			sscqservicecmd.GetCmdSign(cdc),
		)
	}

	txCmd.AddCommand(
		sscqservicecmd.GetCmdBroadCast(cdc),
		client.LineBreak,
	)

	for _, m := range mc {
		txCmd.AddCommand(m.GetTxCmd())
	}

	return txCmd
}

func initConfig(cmd *cobra.Command) error {
	home, err := cmd.PersistentFlags().GetString(cli.HomeFlag)
	if err != nil {
		return err
	}

	cfgFile := path.Join(home, "config", "config.toml")
	if _, err := os.Stat(cfgFile); err == nil {
		viper.SetConfigFile(cfgFile)

		if err := viper.ReadInConfig(); err != nil {
			return err
		}
	}
	if err := viper.BindPFlag(client.FlagChainID, cmd.PersistentFlags().Lookup(client.FlagChainID)); err != nil {
		return err
	}
	if err := viper.BindPFlag(cli.EncodingFlag, cmd.PersistentFlags().Lookup(cli.EncodingFlag)); err != nil {
		return err
	}
	return viper.BindPFlag(cli.OutputFlag, cmd.PersistentFlags().Lookup(cli.OutputFlag))
}
