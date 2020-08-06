package goVsysSdk

import (
	"encoding/json"
)

type RegisterContractTokenReq struct{
	Sender              *Account
	Max                 int64
	Unity               int64
	TokenDescription    string
	ContractDescription string
}
func (api *VsysApi) MustRegisterContractToken(req RegisterContractTokenReq) (resp TransactionResponse) {
	ContractCode:=ContractCodeToken
	data := DataEncoder{}
	data.EncodeArgAmount(3)
	data.Encode(req.Max, DeTypeAmount)
	data.Encode(req.Unity, DeTypeAmount)
	data.Encode(req.TokenDescription, DeTypeShortText)

	dataB:= data.Result()
	return api.MustApiCallRegisterContract(ApiCallRegisterContractReq{
		Sender:              req.Sender,
		ContractCode:        ContractCode,
		ContractDescription: req.ContractDescription,
		DataB:               dataB,
	})
}

func (sdk *VsysApi) SendTokenSimpleAsync(tokenId string,senderAccount *Account,receiverAddress string,amount int64) (err error){
	tx:=senderAccount.BuildSendTokenTransaction(tokenId,receiverAddress,amount,false,[]byte{})
	_,err=sdk.SendExecuteContractTx(tx)
	if err!=nil{
		return err
	}
	return nil
}

func (sdk *VsysApi) SendTokenSimpleAsync2(tokenId string,senderAccount *Account,receiverAddress string,amount int64) (resp TransactionResponse,err error){
	tx:=senderAccount.BuildSendTokenTransaction(tokenId,receiverAddress,amount,false,[]byte{})
	return sdk.SendExecuteContractTx(tx)
}

func (sdk *VsysApi) MustSendTokenSimpleSync(tokenId string,senderAccount *Account,receiverAddress string,amount int64) (resp TransactionResponse){
	tx:=senderAccount.BuildSendTokenTransaction(tokenId,receiverAddress,amount,false,[]byte{})
	resp,err:=sdk.SendExecuteContractTx(tx)
	if err!=nil{
		panic(err)
	}
	sdk.MustWaitPaymentOkByTransactionResponse(resp)
	return resp
}

type TokenIssueReq struct{
	Sender *Account
	TokenId string
	ContractId string
	Amount int64
}
func (api *VsysApi) MustTokenIssue(req TokenIssueReq) (resp TransactionResponse){
	data := DataEncoder{}
	data.EncodeArgAmount(1)
	data.Encode(req.Amount, DeTypeAmount)
	dataB:= data.Result()
	if req.ContractId==""{
		req.ContractId = TokenId2ContractId(req.TokenId)
	}
	return api.MustApiCallExecuteContract(ApiCallExecuteContractReq{
		Sender:     req.Sender,
		ContractId: req.ContractId,
		FuncIdx:    FuncidxIssue,
		FuncData:   dataB,
	})
}

type TokenInfoResp struct{
	TokenID     string `json:"tokenId"`
	ContractID  string `json:"contractId"`
	Max         int64  `json:"max"`
	Total       int64  `json:"total"`
	Unity       int    `json:"unity"`
	Description string `json:"description"`
}

/*
{
  "tokenId" : "TWZZfKFqcaNVe5TrphLRNEm5DQFnBRJMjDDByqv84",
  "contractId" : "CC8Jx8aLkKVQmzuHBWNnhCSkn1GBLcjZ32k",
  "max" : 1500000000000000000,
  "total" : 1000000000000000000,
  "unity" : 1000000000,
  "description" : "15sb7d"
}
*/
func (sdk *VsysApi)MustGetTokenInfo(tokenId string) (TokenInfoResp){
	b,err:= sdk.httpGet("/contract/tokenInfo/"+ urlv(tokenId))
	if err!=nil{
		panic(err)
	}
	var resp TokenInfoResp
	err = json.Unmarshal(b,&resp)
	if err!=nil{
		panic(err)
	}
	return resp
}

type GetContractTokenBalanceResp struct {
	Address string `json:"address"`
	TokenId string `json:"tokenId"`
	Balance int64  `json:"balance"`
	Unity   int    `json:"unity"`
}
/*
{
  "address" : "AR695aEbZPsDQzVjHBLvDYxadrpe21zdfHf",
  "tokenId" : "TWZZfKFqcaNVe5TrphLRNEm5DQFnBRJMjDDByqv84",
  "balance" : 35141117913,
  "unity" : 1000000000
}
*/
func (sdk *VsysApi) MustGetContractTokenBalanceObj(address string,tokenId string) (GetContractTokenBalanceResp){
	b,err:= sdk.httpGet("/contract/balance/"+ urlv(address)+"/"+ urlv(tokenId))
	if err!=nil{
		panic(err)
	}
	var resp GetContractTokenBalanceResp
	err = json.Unmarshal(b,&resp)
	if err!=nil{
		panic(err)
	}
	return resp
}

func (sdk *VsysApi) MustGetContractTokenBalance(address string,tokenId string) (v int64){
	obj:=sdk.MustGetContractTokenBalanceObj(address,tokenId)
	return obj.Balance
}

type TokenDepositReq struct{
	Sender *Account
	TokenId string
	ReceiveAddress string
	TokenContractId string
	ReceiveContractId string
	Amount int64
}
func (api *VsysApi) MustTokenDeposit(req TokenDepositReq) (resp TransactionResponse){
	data := DataEncoder{}
	data.EncodeArgAmount(3)
	data.WriteAddress(req.ReceiveAddress)
	data.WriteContractAccount(req.ReceiveContractId)
	data.WriteAmount(int64(req.Amount))
	dataB:= data.Result()
	if req.TokenContractId==""{
		req.TokenContractId = TokenId2ContractId(req.TokenId)
	}
	return api.MustApiCallExecuteContract(ApiCallExecuteContractReq{
		Sender:     req.Sender,
		ContractId: req.TokenContractId,
		FuncIdx:    FuncidxDeposit,
		FuncData:   dataB,
	})
}

type TokenWithdrawReq struct{
	Sender *Account
	TokenId string
	ReceiveAddress string
	TokenContractId string
	ReceiveContractId string
	Amount int64
}
func (api *VsysApi) MustTokenWithdraw(req TokenWithdrawReq) (resp TransactionResponse){
	data := DataEncoder{}
	data.EncodeArgAmount(3)
	data.WriteContractAccount(req.ReceiveContractId)
	data.WriteAddress(req.ReceiveAddress)
	data.WriteAmount(int64(req.Amount))
	dataB:= data.Result()
	if req.TokenContractId==""{
		req.TokenContractId = TokenId2ContractId(req.TokenId)
	}
	return api.MustApiCallExecuteContract(ApiCallExecuteContractReq{
		Sender:     req.Sender,
		ContractId: req.TokenContractId,
		FuncIdx:    FuncidxWithdraw,
		FuncData:   dataB,
	})
}