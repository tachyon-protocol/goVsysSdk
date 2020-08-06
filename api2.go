package goVsysSdk

import (
	"encoding/json"
	"errors"
	"sync"
)

func (sdk *VsysApi) NewAccountFromSeedAndNonce(seed string,nonce int) (acc *Account){
	sender:=InitAccount(sdk.network)
	sender.BuildFromSeed(seed,nonce)
	return sender
}

func (sdk *VsysApi) NewAccountFromSeedAndNonceV2(seed string,nonce int) (acc *Account){
	sender:=InitAccount(sdk.network)
	sender.BuildFromSeed(seed,nonce)
	return sender
}

func (sdk *VsysApi) NewAccountFromSeedHash(seedHash []byte) (acc *Account){
	sender:=GenerateKeyPair(seedHash)
	sender.SetNetwork(sdk.network)
	return sender
}

func (sdk *VsysApi)GetAccountBalance(address string) (v int64,err error){
	b,err:= sdk.httpGet("/addresses/balance/details/"+ urlv(address))
	if err!=nil{
		return 0,err
	}
	type respT struct {
		Address        string `json:"address"`
		Regular        int64  `json:"regular"`
		MintingAverage int    `json:"mintingAverage"`
		Available      int64  `json:"available"`
		Effective      int64  `json:"effective"`
		Height         int    `json:"height"`
	}
	var respO respT
	err = json.Unmarshal(b,&respO)
	if err!=nil{
		return 0,err
	}
	if respO.Address!=address{
		return 0,errors.New("dnv4ub5faz")
	}
	return respO.Available,nil
}

func (sdk *VsysApi)MustGetAccountBalance(address string) (v int64) {
	v,err:=sdk.GetAccountBalance(address)
	if err!=nil{
		panic(err)
	}
	return v
}
func (sdk *VsysApi)MustGetAccountBalanceString(address string) (s string) {
	v:=sdk.MustGetAccountBalance(address)
	return FormatVsys(v)
}

func (tran *HistoryTransaction) GetFirstSenderAddress() string{
	if len(tran.Proofs)==0{
		return ""
	}
	return tran.Proofs[0].Address
}



func (sdk *VsysApi) AddToWallet(account *Account){
	addr:=account.GetAddress()
	sdk.locker.Lock()
	sdk.walletMap[addr] = account
	sdk.locker.Unlock()
}

func (sdk *VsysApi) GetAccountFromWallet(addr string) *Account{
	sdk.locker.Lock()
	acc:=sdk.walletMap[addr]
	sdk.locker.Unlock()
	return acc
}

func (sdk *VsysApi) DeleteAccountFromWallet(addr string) {
	sdk.locker.Lock()

	sdk.locker.Unlock()
}

func (sdk *VsysApi) IsValidAddress(address string)bool{
	return IsValidateAddress(address,sdk.network)
}

type NewVsysApiReq struct{
	NodeAddress string
	Network NetType
	Name string // can be get by GetSdkName()
	ApiKey string
}

func NewVsysApi(req NewVsysApiReq) *VsysApi{
	return &VsysApi{
		nodeAddress: req.NodeAddress,
		network:     req.Network,
		name: req.Name,
		walletMap:   map[string]*Account{},
		req: req,
	}
}

func NewPublicTestNetApi() *VsysApi{
	return NewVsysApi(NewVsysApiReq{
		NodeAddress: "http://test.v.systems:9922",
		Network: Testnet,
		Name: "PubTestNet",
		ApiKey: "vsystest2018",
	})
}

func NewPublicMainNetApi() *VsysApi{
	return NewVsysApi(NewVsysApiReq{
		NodeAddress: "https://wallet.v.systems/api",
		Network: Mainnet,
		Name: "PubMainNet",
	})
}

func (sdk *VsysApi) GetSdkName() string{
	return sdk.name
}

type VsysApiList struct{
	m map[string]*VsysApi
	locker sync.RWMutex
}

func (list *VsysApiList) AddVsysApi(api *VsysApi){
	if api==nil{
		return
	}
	list.locker.Lock()
	if list.m==nil{
		list.m = map[string]*VsysApi{}
	}
	list.m[api.GetSdkName()] = api
	list.locker.Unlock()
}

func (list *VsysApiList) GetByName(name string) *VsysApi{
	list.locker.RLock()
	if list.m==nil {
		list.locker.RUnlock()
		return nil
	}
	api:=list.m[name]
	list.locker.RUnlock()
	return api
}

func (list *VsysApiList) GetList() []*VsysApi{
	var output []*VsysApi
	list.locker.RLock()
	for _,api:=range list.m{
		output = append(output,api)
	}
	list.locker.RUnlock()
	return output
}