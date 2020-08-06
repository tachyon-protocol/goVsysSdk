package goVsysSdk

import (
	"time"
	"sync"
)

type Transaction struct {
	TxId            string `json:"txId,omitempty"`
	Timestamp       int64  `json:"timestamp"`
	Fee             int64  `json:"fee"`
	FeeScale        int16  `json:"feeScale"`
	Amount          int64  `json:"amount,omitempty"`
	SenderPublicKey string `json:"senderPublicKey"`
	Attachment      []byte `json:"attachment,omitempty"`
	Recipient       string `json:"recipient,omitempty"`
	Signature       string `json:"signature,omitempty"`
	// contract
	Contract      string `json:"contract,omitempty"`
	InitData      string `json:"initData,omitempty"`
	ContractId    string `json:"contractId,omitempty"`
	TokenIdx      int32  `json:"tokenIdx,omitempty"`
	Description   string `json:"description,omitempty"`
	FunctionIndex int16  `json:"functionIndex"`
	FunctionData  string  `json:"functionData,omitempty"`
	txType        uint8
}

func NewPaymentTransaction(recipient string, amount int64, attachment []byte) *Transaction {
	return &Transaction{
		Timestamp:  newTimestampNow(),
		Fee:        DefaultTxFee,
		FeeScale:   DefaultFeeScale,
		Recipient:  recipient,
		Amount:     amount,
		Attachment: attachment,
		txType:     TxTypePayment,
	}
}

func NewLeaseTransaction(recipient string, amount int64) *Transaction {
	return &Transaction{
		Timestamp: newTimestampNow(),
		Fee:       DefaultTxFee,
		FeeScale:  DefaultFeeScale,
		Recipient: recipient,
		Amount:    amount,
		txType:    TxTypeLease,
	}
}

func NewCancelLeaseTransaction(txId string) *Transaction {
	return &Transaction{
		Timestamp: newTimestampNow(),
		Fee:       DefaultTxFee,
		FeeScale:  DefaultFeeScale,
		TxId:      txId,
		txType:    TxTypeCancelLease,
	}
}

func NewRegisterTransaction(contract string, data string, contractDescription string) *Transaction {
	return &Transaction{
		txType:      TxTypeContractRegister,
		Contract:    contract,
		InitData:    data,
		Description: contractDescription,
		Fee:         100 * VsysPrecision,
		FeeScale:    DefaultFeeScale,
		Timestamp:   newTimestampNow(),
	}
}

func NewExecuteTransaction(contractId string, funcIdx int16, funcData string, attachment []byte) *Transaction {
	return &Transaction{
		txType:        TxTypeContractExecute,
		ContractId:    contractId,
		FunctionIndex: funcIdx,
		FunctionData:  funcData,
		Attachment:    attachment,
		Fee:           ContractExecFee,
		FeeScale:      DefaultFeeScale,
		Timestamp:     newTimestampNow(),
	}
}

func (tx *Transaction) GetTxType() int {
	return int(tx.txType)
}

// BuildTxData generate data which is used to be signed
func (tx *Transaction) BuildTxData() []byte {
	data := make([]byte, 0)
	data = append(data, tx.txType)
	if tx.Timestamp <= 0 {
		tx.Timestamp = newTimestampNow()
	}
	if tx.Fee <= 0 {
		tx.Fee = DefaultTxFee
	}
	if tx.FeeScale <= 0 {
		tx.FeeScale = DefaultFeeScale
	}
	switch tx.txType {
	case TxTypePayment:
		return tx.buildPaymentData(data)
	case TxTypeLease:
		return tx.buildLeaseData(data)
	case TxTypeCancelLease:
		return tx.buildLeaseCancelData(data)
	case TxTypeContractRegister:
		return tx.buildRegisterContractData(data)
	case TxTypeContractExecute:
		return tx.buildExecuteContractData(data)
	}
	return data
}

func (tx *Transaction) buildPaymentData(data []byte) []byte {
	data = append(data, uint64ToByte(uint64(tx.Timestamp))...)
	data = append(data, uint64ToByte(uint64(tx.Amount))...)
	data = append(data, uint64ToByte(uint64(tx.Fee))...)
	data = append(data, uint16ToByte(uint16(tx.FeeScale))...)
	data = append(data, MustBase58Decode(tx.Recipient)...)
	data = append(data, uint16ToByte(uint16(len(tx.Attachment)))...)
	data = append(data, tx.Attachment...)
	return data
}

func (tx *Transaction) buildLeaseData(data []byte) []byte {
	data = append(data, MustBase58Decode(tx.Recipient)...)
	data = append(data, uint64ToByte(uint64(tx.Amount))...)
	data = append(data, uint64ToByte(uint64(tx.Fee))...)
	data = append(data, uint16ToByte(uint16(tx.FeeScale))...)
	data = append(data, uint64ToByte(uint64(tx.Timestamp))...)
	return data
}

func (tx *Transaction) buildLeaseCancelData(data []byte) []byte {
	data = append(data, uint64ToByte(uint64(tx.Fee))...)
	data = append(data, uint16ToByte(uint16(tx.FeeScale))...)
	data = append(data, uint64ToByte(uint64(tx.Timestamp))...)
	data = append(data, MustBase58Decode(tx.TxId)...)
	return data
}

func (tx *Transaction) buildRegisterContractData(data []byte) []byte {
	data = append(data, bytesToByteArrayWithSize(MustBase58Decode(tx.Contract))...)
	data = append(data, bytesToByteArrayWithSize(MustBase58Decode(tx.InitData))...)
	data = append(data, bytesToByteArrayWithSize([]byte(tx.Description))...)
	data = append(data, uint64ToByte(uint64(tx.Fee))...)
	data = append(data, uint16ToByte(uint16(tx.FeeScale))...)
	data = append(data, uint64ToByte(uint64(tx.Timestamp))...)
	return data
}

func (tx *Transaction) buildExecuteContractData(data []byte) []byte {
	data = append(data, MustBase58Decode(tx.ContractId)...)
	data = append(data, uint16ToByte(uint16(tx.FunctionIndex))...)
	data = append(data, bytesToByteArrayWithSize(MustBase58Decode(tx.FunctionData))...)
	data = append(data, uint16ToByte(uint16(len(tx.Attachment)))...)
	data = append(data, tx.Attachment...)
	data = append(data, uint64ToByte(uint64(tx.Fee))...)
	data = append(data, uint16ToByte(uint16(tx.FeeScale))...)
	data = append(data, uint64ToByte(uint64(tx.Timestamp))...)
	return data
}

var gLastTimeStamp int64
var glastTimeStampLocker sync.Mutex
func newTimestampNow() int64{
	v:=time.Now().UnixNano()
	glastTimeStampLocker.Lock()
	if v<=gLastTimeStamp{
		v = gLastTimeStamp+1
	}
	gLastTimeStamp = v
	glastTimeStampLocker.Unlock()
	return v
}