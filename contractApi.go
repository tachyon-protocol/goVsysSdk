package goVsysSdk

type ApiCallRegisterContractReq struct{
	Sender              *Account
	ContractCode        string
	ContractDescription string
	DataB []byte // from DataEncoder{}
}
func (api *VsysApi) MustApiCallRegisterContract(req ApiCallRegisterContractReq)(resp TransactionResponse){
	transaction := NewRegisterTransaction(req.ContractCode, Base58Encode(req.DataB), req.ContractDescription)
	transaction.SenderPublicKey = req.Sender.PublicKey()
	transaction.Signature = req.Sender.SignData(transaction.BuildTxData())

	resp,err:=api.SendRegisterContractTx(transaction)
	panicIfError(err)
	return resp
}

type ApiCallExecuteContractReq struct{
	Sender              *Account
	ContractId string
	FuncIdx int16
	FuncData []byte
	Attachment []byte
}
func (api *VsysApi) MustApiCallExecuteContract(req ApiCallExecuteContractReq)(resp TransactionResponse){
	transaction := NewExecuteTransaction(req.ContractId, req.FuncIdx, Base58Encode(req.FuncData), req.Attachment)
	transaction.SenderPublicKey = req.Sender.PublicKey()
	transaction.Signature = req.Sender.SignData(transaction.BuildTxData())
	resp,err:=api.SendExecuteContractTx(transaction)
	panicIfError(err)
	return resp
}

type GetContractDataResp struct{
	ContractId string `json:"contractId"`
	Key        string `json:"key"`
	Height     int64    `json:"height"`
	DbName     string `json:"dbName"`
	DataType string `json:"dataType"`
	ValueI     int64
	ValueS     string
}
func (api *VsysApi) MustApiCallGetContractData(contractId string,key string) (resp GetContractDataResp){
	respB,err:=api.httpGet("/contract/data/"+urlv(contractId)+"/"+urlv(key))
	panicIfError(err)
	//fmt.Println(string(respB))
	mustJsonUnmarshal(respB,&resp)
	switch resp.DataType {
	case "","Timestamp":
		var resp2 struct {
			Value     int64 `json:"value"`
		}
		mustJsonUnmarshal(respB,&resp2)
		resp.ValueI = resp2.Value
		return resp
	case "Address","PublicKey","Boolean":
		var resp2 struct {
			Value     string `json:"value"`
		}
		mustJsonUnmarshal(respB,&resp2)
		resp.ValueS = resp2.Value
		return resp
	default:
		panic("y8s46u9raq ["+resp.DataType+"]")
	}
	return resp
}