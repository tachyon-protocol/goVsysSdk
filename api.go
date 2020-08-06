package goVsysSdk

import (
	"sync"
)

type VsysApi struct {
	req NewVsysApiReq
	name string
	nodeAddress string
	network     NetType
	locker sync.Mutex
	walletMap   map[string]*Account
}

type TransactionResponse struct {
	HistoryTransaction
	CommonResp
}

type Proof struct {
	ProofType string `json:"proofType"`
	PublicKey string `json:"publicKey"`
	Address   string `json:"address"` // from read
	Signature string `json:"signature"`
}

type CommonResp struct {
	Error   int    `json:"error"`
	Message string `json:"message"`
}

func (sdk *VsysApi) MustGetNodeVersion() (version string){
	respB,err:=sdk.httpGet("/node/version")
	if err!=nil{
		panic(err)
	}
	var versionResp struct{
		Version string `json:"version"`
	}
	mustJsonUnmarshal(respB,&versionResp)
	return versionResp.Version
}

