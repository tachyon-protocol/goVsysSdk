package goVsysSdk

import (
	"golang.org/x/crypto/curve25519"
)

type Account struct {
	publicKey  []byte
	privateKey []byte
	network    NetType
	accSeed    string
}

// get account address string
func (acc *Account) GetAddress() string {
	return publicKeyToAddress(acc.publicKey, acc.network)
}

func (acc *Account) SetNetwork(network NetType){
	acc.network = network
}

func publicKeyToAddress(publicKey []byte, network NetType) string {
	uAddr := make([]byte, 0)
	uAddr = append(uAddr, addrVersion, byte(network))
	uAddr = append(uAddr, HashChain(publicKey)[:20]...)
	return Base58Encode(append(uAddr, HashChain(uAddr)[:4]...))
}

// get account privateKey string
func (acc *Account) PrivateKey() string {
	return Base58Encode(acc.privateKey)
}

// get account publicKey string
func (acc *Account) PublicKey() string {
	return Base58Encode(acc.publicKey)
}

func (acc *Account) AccountSeed() string {
	return acc.accSeed
}

// SignData sign data bytes and
// the output is base58 encoded data
func (acc *Account) SignData(data []byte) string {
	return Base58Encode(Sign(acc.privateKey, data, genRandomBytes(64)))
}

// VerifySignature check if signature is correct
func (acc *Account) VerifySignature(data, signature []byte) bool {
	return Verify(acc.publicKey, data, signature) == 1
}

// InitAccount return account with network initiated
func InitAccount(network NetType) *Account {
	return &Account{network: network}
}

// BuildFromPrivateKey build account using privateKey
func (acc *Account) BuildFromPrivateKey(privateKey string) {
	var bPrivateKey [32]byte
	var originPublicKey = new([32]byte)
	copy(bPrivateKey[:], MustBase58Decode(privateKey)[:])
	curve25519.ScalarBaseMult(originPublicKey, &bPrivateKey)
	acc.publicKey = originPublicKey[:]
	acc.privateKey = bPrivateKey[:]
}

// BuildFromPrivateKey build account using seed and nonce
func (acc *Account) BuildFromSeed(seed string, nonce int) {
	seedHash := BuildSeedHash(seed, nonce)
	keyPair := GenerateKeyPair(seedHash)
	acc.publicKey = keyPair.publicKey
	acc.privateKey = keyPair.privateKey
	acc.accSeed = seed
}

// BuildPayment build payment transaction
// recipient should be address
// amount is in minimum unit
// attachment can be empty
func (acc *Account) BuildPayment(recipient string, amount int64, attachment []byte) *Transaction {
	transaction := NewPaymentTransaction(recipient, amount, attachment)
	transaction.SenderPublicKey = acc.PublicKey()
	transaction.Signature = acc.SignData(transaction.BuildTxData())
	return transaction
}

// BuildLeasing build leasing transaction
// recipient should be address
// amount is in minimum unit
func (acc *Account) BuildLeasing(recipient string, amount int64) *Transaction {
	transaction := NewLeaseTransaction(recipient, amount)
	transaction.SenderPublicKey = acc.PublicKey()
	transaction.Signature = acc.SignData(transaction.BuildTxData())
	return transaction
}

// BuildCancelLeasing build Cancel transaction
func (acc *Account) BuildCancelLeasing(txId string) *Transaction {
	transaction := NewCancelLeaseTransaction(txId)
	transaction.SenderPublicKey = acc.PublicKey()
	transaction.Signature = acc.SignData(transaction.BuildTxData())
	return transaction
}



// BuildExecuteContract build ExecuteContract transaction
func (acc *Account) BuildExecuteContract(contractId string, funcIdx int16, funcData []byte, attachment []byte) *Transaction {
	transaction := NewExecuteTransaction(contractId, funcIdx, Base58Encode(funcData), attachment)
	transaction.SenderPublicKey = acc.PublicKey()
	transaction.Signature = acc.SignData(transaction.BuildTxData())
	return transaction
}

// BuildExecuteContract build SendToken transaction
func (acc *Account) BuildSendTokenTransaction(tokenId string, recipient string, amount int64, isSplitSupported bool, attachment []byte) *Transaction {
	a := &Contract{
		ContractId: TokenId2ContractId(tokenId),
		Amount:     amount,
		Recipient:  recipient,
	}
	funcData := a.BuildSendData()
	funcIdx := FuncidxSend
	if isSplitSupported {
		funcIdx = FuncidxSendSplit
	}
	transaction := NewExecuteTransaction(a.ContractId, int16(funcIdx), Base58Encode(funcData), attachment)
	transaction.SenderPublicKey = acc.PublicKey()
	transaction.Signature = acc.SignData(transaction.BuildTxData())
	return transaction
}

type BuildPayment2Req struct{
	Recipient string
	Amount int64
}
func (acc *Account) BuildPayment2(transaction *Transaction) *Transaction {
	transaction.txType = TxTypePayment
	transaction.SenderPublicKey = acc.PublicKey()
	transaction.Signature = acc.SignData(transaction.BuildTxData())
	return transaction
}