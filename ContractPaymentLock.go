package goVsysSdk

const ContractCodeLock = `4Qgfi31k6qfLxTguJg8AeYzmmgaCTJCEPQyAdoRUUSrFDc91PhkdU6C8QQSsNCFc2xEud2XnuQ4YNJ51HgdNtBdnxZcU5Rnqdzyop41Ck81v4nRKkHpTdTrfD8vTur2w4mTFeTFKVzGvGjpHXUVvT47vZiKLBHSB7FHHpGf69bu8DQGXWu6xnZZkn9v2Rfc9mByhwVLSNghNdRhrQwRWPFJ9Qt7Yb8N8WdmcUCAC6PrC3Ha3Z9w7dyf6CsKcCMS6JmB2gvNQitm9jqAfjRxDdqPBUR6TtyjSdmHP9BZRGgiVCaQH7X8fbJZVWSib4RXvFoSrqY4SfVftDY3PU4hXASaRWbaheB8m4VgM4mA8nKDbZvRWZtZ4cHdWeNFyVPs6HxHQZHrQ3GZGNPjmBSyAkGRFS7i5dK8aYWQDEYu1Xijk63UFAWuf6tRdR44ZgRjWGUZJtdQBDFB38XaU8LSFEj2eaC1yNqZ6nnGeRXDzS1q3YKsGyJTqaDDMHvPHiHonGn76JQHAZN7eGU7biaSLxoikW4MaTPSfmcTmDyPGJyJNHjc8MrpV8aQSaGGyDkf1a9MpoJcyEjsPFQbxYzSJVqFEFg2oUL7Z8VUtJK2kYcWDz7w8UiiQqe3uuQnKDGb1nJ5Ad3W8ZPfVP6YHbJrnBKZXMMypNoveokVvxZMCkSNYDsoBxJzrwFvm5DcDJbePQU6VbeZ5SzQw9XTAw4DZpxkQm9RwRE9PXPqogpp9P6LhaiUa6ZD1cWUAHypjWLJ2Rds96oap3biBp5aESunuh99HByoXg5Aa7EQ3FrEvmeq9TLVFYpJraZyW`

func (api *VsysApi) MustRegisterContractLock(Sender *Account,Vsys_token_id string) (resp TransactionResponse){
	data := DataEncoder{}
	data.EncodeArgAmount(1)
	data.Encode(Vsys_token_id, DeTypeTokenId)
	dataB:=data.Result()
	return api.MustApiCallRegisterContract(ApiCallRegisterContractReq{
		Sender: Sender,
		ContractCode: ContractCodeLock,
		DataB: dataB,
	})
}

type ContractLockReq struct{
	Sender *Account
	ContractId string
	Timestamp int64
}

func (api *VsysApi) MustContractLock(req ContractLockReq)(resp TransactionResponse){
	data := DataEncoder{}
	data.EncodeArgAmount(1)
	data.Encode(req.Timestamp,DeTypeTimeStamp)
	dataB:=data.Result()
	return api.MustApiCallExecuteContract(ApiCallExecuteContractReq{
		Sender: req.Sender,
		ContractId: req.ContractId,
		FuncIdx: 0,
		FuncData: dataB,
	})
}