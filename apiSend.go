package goVsysSdk

import (
	"time"
	"errors"
	"encoding/json"
)

// transfer vsys, do not wait the transaction go to chain
func (sdk *VsysApi) SendPaymentSimpleAsync(senderAccount *Account,receiverAddress string,amount int64) (err error){
	tx:=senderAccount.BuildPayment(receiverAddress,amount,[]byte{})
	_,err=sdk.SendPaymentTx(tx)
	if err!=nil{
		return err
	}
	return nil
}

func (sdk *VsysApi) SendPaymentSimpleAsync2(senderAccount *Account,receiverAddress string,amount int64) (resp TransactionResponse,err error){
	tx:=senderAccount.BuildPayment(receiverAddress,amount,[]byte{})
	return sdk.SendPaymentTx(tx)
}

func (sdk *VsysApi) MustSendPaymentSimpleAsync2(senderAccount *Account,receiverAddress string,amount int64) (resp TransactionResponse){
	tx:=senderAccount.BuildPayment(receiverAddress,amount,[]byte{})
	resp,err:= sdk.SendPaymentTx(tx)
	if err!=nil{
		panic(err)
	}
	return resp
}

// transfer vsys, wait the transaction go to chain,(it may be refused by the chain)
func (sdk *VsysApi) SendPaymentSimpleSync(senderAccount *Account,receiverAddress string,amount int64) (err error){
	tx:=senderAccount.BuildPayment(receiverAddress,amount,[]byte{})
	resp,err:=sdk.SendPaymentTx(tx)
	if err!=nil{
		return err
	}
	return sdk.WaitPaymentConfirmedByTranId(resp.Id,time.Now().Add(time.Minute))
}

func (sdk *VsysApi) MustSendPaymentSimpleSync(senderAccount *Account,receiverAddress string,amount int64) {
	err:=sdk.SendPaymentSimpleSync(senderAccount,receiverAddress,amount)
	if err!=nil{
		panic(err)
	}
	return
}

func (sdk *VsysApi) WaitPaymentConfirmedByTranId(id string, deadline time.Time) (err error){
	for {
		list,err:=sdk.GetUnconfirmedHistoryTransactionList()
		if err!=nil{
			return nil
		}
		hasFound:=false
		for _,tran:=range list{
			if tran.Id ==id{
				hasFound = true
				break
			}
		}
		if hasFound==false{
			break
		}
		if time.Now().After(deadline){
			return errors.New("8wdc5ckzrb")
		}
		time.Sleep(time.Second)
	}
	for {
		ht,err:=sdk.GetHistoryTransactionById(id)
		if err==nil && ht.Id == id{
			if ht.Status!="Success"{
				return errors.New("rqztpbezkh ["+ht.Status+"]")
			}
			return nil
		}
		if time.Now().After(deadline){
			return errors.New("8wdc5ckzrb")
		}
		time.Sleep(time.Second)
	}
}



func (api *VsysApi)postSendTx(path string, tx *Transaction) (resp TransactionResponse, err error) {
	data, err := api.httpPost(path, tx)
	if err != nil {
		return resp,err
	}
	//fmt.Println("postSendTx",string(data))
	err = json.Unmarshal(data, &resp)
	if err != nil {
		return resp,err
	}
	return resp,nil
}

func (api *VsysApi)SendPaymentTx(tx *Transaction) (resp TransactionResponse, err error) {
	return api.postSendTx("/vsys/broadcast/payment", tx)
}

func (api *VsysApi)SendLeasingTx(tx *Transaction) (resp TransactionResponse, err error) {
	return api.postSendTx("/leasing/broadcast/lease", tx)
}

func (api *VsysApi)SendCancelLeasingTx(tx *Transaction) (resp TransactionResponse, err error) {
	return api.postSendTx("/leasing/broadcast/cancel", tx)
}

func (api *VsysApi)SendRegisterContractTx(tx *Transaction) (resp TransactionResponse, err error) {
	return api.postSendTx("/contract/broadcast/register", tx)
}

func (api *VsysApi)SendExecuteContractTx(tx *Transaction) (resp TransactionResponse, err error) {
	return api.postSendTx("/contract/broadcast/execute", tx)
}

func (sdk *VsysApi) MustWaitPaymentConfirmedByTranId(id string, deadline time.Time) {
	err:=sdk.WaitPaymentConfirmedByTranId(id,deadline)
	if err!=nil{
		panic(err)
	}
}

func (api *VsysApi) MustWaitPaymentOkByTransactionResponse(resp TransactionResponse) {
	api.MustWaitPaymentConfirmedByTranId(resp.Id,time.Now().Add(time.Minute*2))
}