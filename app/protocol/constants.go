package protocol

import sdk "github.com/orientwalt/htdf/types"

const (
	// all store name
	AccountStore         = "acc"
	StakeStore           = "staking"
	StakeTransientStore  = "transient_stake"
	MintStore            = "mint"
	DistrStore           = "distr"
	DistrTransientStore  = "transient_distr"
	SlashingStore        = "slashing"
	GovStore             = "gov"
	FeeStore             = "fee"
	ParamsStore          = "params"
	ParamsTransientStore = "transient_params"
	ServiceStore         = "service"
	GuardianStore        = "guardian"
	UpgradeStore         = "upgrade"

	// all route for query and handler
	BankRoute     = "bank"
	AccountRoute  = AccountStore
	StakeRoute    = StakeStore
	DistrRoute    = DistrStore
	SlashingRoute = SlashingStore
	GovRoute      = GovStore
	ParamsRoute   = ParamsStore
	ServiceRoute  = ServiceStore
	GuardianRoute = GuardianStore
	UpgradeRoute  = UpgradeStore
)

var (
	KeyMain     = sdk.NewKVStoreKey(sdk.MainStore)
	KeyAccount  = sdk.NewKVStoreKey(AccountStore)
	KeyStake    = sdk.NewKVStoreKey(StakeStore)
	TkeyStake   = sdk.NewTransientStoreKey(StakeTransientStore)
	KeyMint     = sdk.NewKVStoreKey(MintStore)
	KeyDistr    = sdk.NewKVStoreKey(DistrStore)
	TkeyDistr   = sdk.NewTransientStoreKey(DistrTransientStore)
	KeySlashing = sdk.NewKVStoreKey(SlashingStore)
	KeyGov      = sdk.NewKVStoreKey(GovStore)
	KeyFee      = sdk.NewKVStoreKey(FeeStore)
	KeyParams   = sdk.NewKVStoreKey(ParamsStore)
	TkeyParams  = sdk.NewTransientStoreKey(ParamsTransientStore)
	KeyService  = sdk.NewKVStoreKey(ServiceStore)
	KeyGuardian = sdk.NewKVStoreKey(GuardianStore)
	KeyUpgrade  = sdk.NewKVStoreKey(UpgradeStore)
)
