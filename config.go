package goVsysSdk

type NetType byte

const (
	//Protocol          = "v.systems"
	//Api               = 3
	addrVersion uint8 = 5

	// Fee
	VsysPrecision   int64 = 1e8
	ContractExecFee int64 = 3e7
	DefaultTxFee    int64 = 1e7
	DefaultFeeScale int16 = 100

	// Network
	Testnet NetType = 'T'
	Mainnet NetType = 'M'

	// TX_TYPE
	TxTypePayment          = 2
	TxTypeLease            = 3
	TxTypeCancelLease      = 4
	TxTypeMinting          = 5
	TxTypeContractRegister = 8
	TxTypeContractExecute  = 9

	//contract funcIdx variable
	ActionInit      = "init"
	ActionSupersede = "supersede"
	ActionIssue     = "issue"
	ActionDestroy   = "destroy"
	ActionSplit     = "split"
	ActionSend      = "send"
	ActionTransfer  = "transfer"
	ActionDeposit   = "deposit"
	ActionWithdraw  = "withdraw"

	// function index
	FuncidxSupersede     = 0
	FuncidxIssue         = 1
	FuncidxDestroy       = 2
	FuncidxSplit         = 3
	FuncidxSend          = 3
	FuncidxSendSplit     = 4
	FuncidxWithdraw      = 6
	FuncidxWithdrawSplit = 7
	FuncidxDeposit       = 5
	FuncidxDepositSplit  = 6
)

