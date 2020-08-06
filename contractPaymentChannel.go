package goVsysSdk

type RegisterContractPaymentChannelReq struct{
	Sender              *Account
	Vsys_token_id string
}
func (api *VsysApi) MustRegisterContractPaymentChannel(req RegisterContractPaymentChannelReq) (resp TransactionResponse){
	data := DataEncoder{}
	data.EncodeArgAmount(1)
	data.Encode(req.Vsys_token_id, DeTypeTokenId)
	dataB:=data.Result()
	return api.MustApiCallRegisterContract(ApiCallRegisterContractReq{
		Sender: req.Sender,
		ContractCode: ContractCodePaymentChannel,
		DataB: dataB,
	})
}

/*
recipient_address = Account(chain=chain, seed='<recipient-seed>', nonce=0)
create_recipient_data_entry = DataEntry(recipient_address.address, Type.address)
create_amount_data_entry = DataEntry(50, Type.amount)
create_expiration_time_data_entry = DataEntry(response2["timestamp"] + 10000000000, Type.timestamp)

create_data_stack = [create_recipient_data_entry, create_amount_data_entry, create_expiration_time_data_entry]
 */
type PaymentChannelCreateAndLoadReq struct{
	Sender *Account
	ContractId string
	RecipientAddress string
	Amount int64
	TimeStamp int64
}
func (api *VsysApi) MustPaymentChannelCreateAndLoad(req PaymentChannelCreateAndLoadReq)(resp TransactionResponse){
	data := DataEncoder{}
	data.EncodeArgAmount(3)
	data.Encode(req.RecipientAddress,DeTypeAddress)
	data.Encode(req.Amount,DeTypeAmount)
	data.Encode(req.TimeStamp,DeTypeTimeStamp)
	dataB:=data.Result()
	return api.MustApiCallExecuteContract(ApiCallExecuteContractReq{
		Sender: req.Sender,
		ContractId: req.ContractId,
		FuncIdx: 0,
		FuncData: dataB,
	})
}

type PaymentChannelExtendExpirationTimeReq struct{
	Sender *Account
	ContractId string
	ChannelId string
	TimeStamp int64
}
func (api *VsysApi) MustPaymentChannelExtendExpirationTime(req PaymentChannelExtendExpirationTimeReq)(resp TransactionResponse){
	data := DataEncoder{}
	data.EncodeArgAmount(2)
	data.Encode(MustBase58Decode(req.ChannelId),DeTypeShortByte)
	data.Encode(req.TimeStamp,DeTypeTimeStamp)
	dataB:=data.Result()
	return api.MustApiCallExecuteContract(ApiCallExecuteContractReq{
		Sender: req.Sender,
		ContractId: req.ContractId,
		FuncIdx: 1,
		FuncData: dataB,
	})
}

type PaymentChannelLoadReq struct{
	Sender *Account
	ContractId string
	ChannelId string
	Amount int64
}

func (api *VsysApi) MustPaymentChannelLoad(req PaymentChannelLoadReq)(resp TransactionResponse){
	data := DataEncoder{}
	data.EncodeArgAmount(2)
	data.Encode(MustBase58Decode(req.ChannelId),DeTypeShortByte)
	data.Encode(req.Amount,DeTypeAmount)
	dataB:=data.Result()
	return api.MustApiCallExecuteContract(ApiCallExecuteContractReq{
		Sender: req.Sender,
		ContractId: req.ContractId,
		FuncIdx: 2,
		FuncData: dataB,
	})
}

type PaymentChannelAbortReq struct{
	Sender *Account
	ContractId string
	ChannelId string
}
func (api *VsysApi) MustPaymentChannelAbort(req PaymentChannelAbortReq)(resp TransactionResponse){
	data := DataEncoder{}
	data.EncodeArgAmount(1)
	data.Encode(MustBase58Decode(req.ChannelId),DeTypeShortByte)
	dataB:=data.Result()
	return api.MustApiCallExecuteContract(ApiCallExecuteContractReq{
		Sender: req.Sender,
		ContractId: req.ContractId,
		FuncIdx: 3,
		FuncData: dataB,
	})
}

type PaymentChannelUnloadReq struct{
	Sender *Account
	ContractId string
	ChannelId string
}
func (api *VsysApi) MustPaymentChannelUnload(req PaymentChannelUnloadReq)(resp TransactionResponse){
	data := DataEncoder{}
	data.EncodeArgAmount(1)
	data.Encode(MustBase58Decode(req.ChannelId),DeTypeShortByte)
	dataB:=data.Result()
	return api.MustApiCallExecuteContract(ApiCallExecuteContractReq{
		Sender: req.Sender,
		ContractId: req.ContractId,
		FuncIdx: 4,
		FuncData: dataB,
	})
}

func PaymentChannelGenerateSenderPaymentSignature(Sender *Account,ChannelId string,Amount int64) string{
	buf:=[]byte{}
	channelByte:=MustBase58Decode(ChannelId)
	buf = append(buf,uint16ToByte(uint16(len(channelByte)))...)
	buf = append(buf,channelByte...)
	buf = append(buf,uint64ToByte(uint64(Amount))...)
	return Sender.SignData(buf)
}

type PaymentChannelCollectReq struct{
	Receiver *Account
	ContractId string
	ChannelId string
	Amount int64
	Payment_signature_str string
}
func (api *VsysApi) MustPaymentChannelCollect(req PaymentChannelCollectReq)(resp TransactionResponse){
	data := DataEncoder{}
	data.EncodeArgAmount(3)
	data.WriteShortByte(MustBase58Decode(req.ChannelId))
	data.WriteAmount(req.Amount)
	data.WriteShortByte(MustBase58Decode(req.Payment_signature_str))
	dataB:=data.Result()
	return api.MustApiCallExecuteContract(ApiCallExecuteContractReq{
		Sender: req.Receiver,
		ContractId: req.ContractId,
		FuncIdx: 5,
		FuncData: dataB,
	})
}

func (api *VsysApi) MustPaymentChannelGetMasterContractBalance(ContractId string,address string) int64{
	data := DataEncoder{}
	data.Encode(address,DeTypeAddress)
	dataB:=data.Result()
	content:=append([]byte{0},dataB...)
	key:=Base58Encode(content)
	//fmt.Println(hex.Dump(content),key)
	resp:=api.MustApiCallGetContractData(ContractId,key)
	return resp.ValueI
}

func (api *VsysApi) MustPaymentChannelGetCreatorAddressOfChannel(ContractId string,ChannelId string) string{
	resp:=api.mustPaymentChannelContractGetByChannelId(ContractId,ChannelId,1)
	if resp.DataType!="Address"{
		panic(resp.DataType)
	}
	return resp.ValueS
}

func (api *VsysApi) MustPaymentChannelGetCreatorPublicKeyOfChannel(ContractId string,ChannelId string) string{
	resp:=api.mustPaymentChannelContractGetByChannelId(ContractId,ChannelId,2)
	if resp.DataType!="PublicKey"{
		panic(resp.DataType)
	}
	return resp.ValueS
}

func (api *VsysApi) MustPaymentChannelGetRecipientAddressOfChannel(ContractId string,ChannelId string) string{
	resp:=api.mustPaymentChannelContractGetByChannelId(ContractId,ChannelId,3)
	if resp.DataType!="Address"{
		panic(resp.DataType)
	}
	return resp.ValueS
}

func (api *VsysApi) MustPaymentChannelGetChannelCapacity(ContractId string,ChannelId string) int64{
	resp:=api.mustPaymentChannelContractGetByChannelId(ContractId,ChannelId,4)
	if resp.DataType!=""{
		panic(resp.DataType)
	}
	return resp.ValueI
}

func (api *VsysApi) MustPaymentChannelGetChannelCollectedBalance(ContractId string,ChannelId string) int64{
	resp:=api.mustPaymentChannelContractGetByChannelId(ContractId,ChannelId,5)
	if resp.DataType!=""{
		panic(resp.DataType)
	}
	return resp.ValueI
}

func (api *VsysApi) MustPaymentChannelGetChannelExpirationTime(ContractId string,ChannelId string) int64{
	resp:=api.mustPaymentChannelContractGetByChannelId(ContractId,ChannelId,6)
	if resp.DataType!="Timestamp"{
		panic(resp.DataType)
	}
	return resp.ValueI
}

// true mean channel is open, false mean sender abort it.
// false can not call ExtendExpirationTime on it.
func (api *VsysApi) MustPaymentChannelGetChannelIsOpen(ContractId string,ChannelId string) bool{
	resp:=api.mustPaymentChannelContractGetByChannelId(ContractId,ChannelId,7)
	if resp.DataType!="Boolean"{
		panic(resp.DataType)
	}
	return resp.ValueS=="True"
}

func (api *VsysApi) mustPaymentChannelContractGetByChannelId(ContractId string,ChannelId string,index byte) (resp GetContractDataResp){
	data := DataEncoder{}
	data.Encode(MustBase58Decode(ChannelId),DeTypeShortByte)
	dataB:=data.Result()
	content:=append([]byte{index},dataB...)
	key:=Base58Encode(content)
	return api.MustApiCallGetContractData(ContractId,key)
}