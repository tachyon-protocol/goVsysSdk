package goVsysSdk

import (
	"strconv"
	"encoding/json"
	"errors"
	"strings"
)

type HistoryTransaction struct{
	Type       int    `json:"type,omitempty"` // TX_TYPE
	Id         string `json:"id,omitempty"`
	Fee        int64    `json:"fee,omitempty"`
	FeeScale   int64    `json:"feeScale,omitempty"`
	Timestamp  int64  `json:"timestamp,omitempty"`
	Proofs     []Proof `json:"proofs,omitempty"`

	Recipient  string `json:"recipient,omitempty"`
	Amount     int64    `json:"amount,omitempty"`
	Attachment string `json:"attachmen,omitempty"`

	ContractId    string `json:"contractId,omitempty"`
	FunctionIndex int    `json:"functionIndex,omitempty"` // FuncidxSplit,...
	FunctionData  string `json:"functionData,omitempty"`
	Contract   *HtContract `json:"contract,omitempty"`
	InitData    string `json:"initData,omitempty"`
	Description string `json:"description,omitempty"`

	Status     string `json:"status,omitempty"`
	FeeCharged int    `json:"feeCharged,omitempty"`
	Height     int    `json:"height,omitempty"`

	CurrentBlockHeight int    `json:"currentBlockHeight,omitempty"`
}

type HtContract struct{
	LanguageCode    string   `json:"languageCode,omitempty"`
	LanguageVersion int      `json:"languageVersion,omitempty"`
	Triggers        []string `json:"triggers,omitempty"`
	Descriptors     []string `json:"descriptors,omitempty"`
	StateVariables  []string `json:"stateVariables,omitempty"`
	StateMaps       []string `json:"stateMaps,omitempty"`
	Textual         struct {
		Triggers       string `json:"triggers,omitempty"`
		Descriptors    string `json:"descriptors,omitempty"`
		StateVariables string `json:"stateVariables,omitempty"`
		StateMaps      string `json:"stateMaps,omitempty"`
	} `json:"textual,omitempty"`
}

func (ht *HistoryTransaction) IsContractSend() bool{
	return ht.Type==TxTypeContractExecute && ht.ContractId!="" && ht.FunctionIndex==3
}
func (ht *HistoryTransaction) MustDecodeContractSend() HtContractSend{
	if ht.IsContractSend()==false{
		panic("rvxzrh2rqz")
	}
	de := DataEncoder{}
	bytes := MustBase58Decode(ht.FunctionData)
	list := de.Decode(bytes)
	if len(list)!=2{
		panic("5xrrhxtuc6")
	}
	return HtContractSend{
		AddressSend: ht.Proofs[0].Address,
		AddressRecv: list[0].Value.(string),
		Amount: list[1].Value.(int64),
	}
}
type HtContractSend struct{
	AddressSend string
	AddressRecv string
	Amount      int64
}

func (sdk *VsysApi)GetAccountHistoryTransactionList(address string,limit int) (list []HistoryTransaction,err error){
	b,err:= sdk.httpGet("/transactions/address/"+ urlv(address)+"/limit/"+ urlv(strconv.Itoa(limit)))
	if err!=nil{
		return nil,err
	}
	//fmt.Println(string(b))
	type respT [][]HistoryTransaction
	var respO respT
	err = json.Unmarshal(b,&respO)
	if err!=nil{
		return nil,err
	}
	if len(respO)!=1{
		return nil,errors.New("dnv4ub5faz ["+strconv.Itoa(int(len(respO))) )
	}
	return respO[0],nil
}
// TxType 0 to return all transaction
func (sdk *VsysApi)GetAccountAllHistoryTransactionCallback(address string,TxType int,visitor func(tran HistoryTransaction)) (err error){
	const itemPerPage = 10000
	offset:=0
	for {
		list2,err:=sdk.GetAccountHistoryTransactionList2(TransactionList2Req{
			Address: address,
			TxType: TxType,
			Limit: itemPerPage,
			Offset: offset,
		})
		if err!=nil{
			return err
		}
		for _,ht:=range list2{
			visitor(ht)
		}
		offset+=len(list2)
		if len(list2)<itemPerPage{
			return nil
		}
	}
}

type TransactionList2Req struct{
	Address string
	Limit int
	Offset int
	TxType int
}

func (sdk *VsysApi)GetAccountHistoryTransactionList2(req TransactionList2Req) (list []HistoryTransaction,err error){
	type respT struct {
		TotalCount   int `json:"totalCount"`
		Size         int `json:"size"`
		Transactions []HistoryTransaction `json:"transactions"`
	}
	uS:="/transactions/list?address="+ urlv(req.Address)+"&limit="+ urlv(strconv.Itoa(req.Limit))
	if req.Offset!=0{
		uS+="&offset="+ urlv(strconv.Itoa(req.Offset))
	}
	if req.TxType!=0{
		uS+="&txType="+ urlv(strconv.Itoa(req.TxType))
	}
	b,err:= sdk.httpGet(uS)
	if err!=nil{
		return nil,errors.New("rs2u4qna66 "+err.Error())
	}
	var respO respT
	err = json.Unmarshal(b,&respO)
	if err!=nil{
		return nil,errors.New("9mtxu65dyt "+string(b)+" "+err.Error())
	}
	return respO.Transactions,nil
}

func (sdk *VsysApi)MustGetAccountHistoryTransactionList(address string,limit int) (list []HistoryTransaction) {
	list,err:=sdk.GetAccountHistoryTransactionList(address,limit)
	if err!=nil{
		panic(err)
	}
	return list
}
func (sdk *VsysApi)GetUnconfirmedHistoryTransactionList() (list []HistoryTransaction,err error){
	b,err:= sdk.httpGet("/transactions/unconfirmed")
	if err!=nil{
		return nil,err
	}
	var respO []HistoryTransaction
	err = json.Unmarshal(b,&respO)
	if err!=nil{
		return nil,err
	}
	return respO,nil
}
func (sdk *VsysApi)GetHistoryTransactionById(txId string) (ht HistoryTransaction,err error){
	b,err:= sdk.httpGet("/transactions/info/"+ urlv(txId))
	if err!=nil{
		return ht,errors.New("qkw9bg6ufq ["+err.Error()+"]")
	}
	var respO HistoryTransaction
	err = json.Unmarshal(b,&respO)
	if err!=nil{
		return ht,errors.New("w4wm8yw2wb ["+err.Error()+"]")
	}
	if respO.Id !=txId{
		return ht,errors.New("fbhs64quz4 ["+respO.Id +"]")
	}
	return respO,nil
}

// "Transaction is not in blockchain"
func IsErrTransactionNotInBlockChain(err error) bool{
	if err==nil{
		return false
	}
	errS:=err.Error()
	return strings.Contains(errS,`"Transaction is not in blockchain"`)
}

func (sdk *VsysApi) MustGetHistoryTransactionById(txId string)(ht HistoryTransaction){
	ht,err:=sdk.GetHistoryTransactionById(txId)
	panicIfError(err)
	return ht
}

func (sdk *VsysApi)MustGetTransactionUnconfirmedSize() int64{
	b,err:= sdk.httpGet("/transactions/unconfirmed/size")
	if err!=nil{
		panic(err)
	}
	var respO struct{
		Size int64 `json:"size"`
	}
	err = json.Unmarshal(b,&respO)
	if err!=nil{
		panic(err)
	}
	return respO.Size
}