package goVsysSdk

import (
	"time"
)

type ApiRespNodeStatus struct{
	BlockchainHeight int64  `json:"blockchainHeight"`
	StateHeight      int64  `json:"stateHeight"`
	UpdatedDate      string `json:"updatedDate"`
	UpdatedTimestamp int64  `json:"updatedTimestamp"`
}
func (sdk *VsysApi) MustApiGetNodeStatus() (resp ApiRespNodeStatus){
	b,err:= sdk.httpGet("/node/status")
	if err!=nil{
		panic("ms34x43qy4 ["+err.Error()+"]")
	}
	mustJsonUnmarshal(b,&resp)
	return resp
}

type NodePeerConnected struct{
	Address            string `json:"address"`
	ApplicationName    string `json:"applicationName"`
	ApplicationVersion string `json:"applicationVersion"`
	DeclaredAddress    string `json:"declaredAddress"`
	PeerName           string `json:"peerName"`
	PeerNonce          int64  `json:"peerNonce"`
}


func (sdk *VsysApi) MustApiGetPeersConnected() []NodePeerConnected {
	b,err:= sdk.httpGet("/peers/connected")
	if err!=nil{
		panic("6npygzdqyk ["+err.Error()+"]")
	}
	var resp struct{
		Peers []NodePeerConnected `json:"peers"`
	}
	mustJsonUnmarshal(b,&resp)
	return resp.Peers
}

type NodePeerAll struct{
	Address            string `json:"address"`
	LastSeen int64  `json:"lastSeen"`
}

func (sdk *VsysApi) MustApiGetPeersAll() []NodePeerAll {
	b,err:= sdk.httpGet("/peers/all")
	if err!=nil{
		panic("6npygzdqyk ["+err.Error()+"]")
	}
	//fmt.Println("MustApiGetPeersAll",string(b))
	var resp struct{
		Peers []NodePeerAll `json:"peers"`
	}
	mustJsonUnmarshal(b,&resp)
	return resp.Peers
}

type NodePeerSuspended struct{
	Hostname            string `json:"hostname"`
	Timestamp int64  `json:"timestamp"`
}

func (sdk *VsysApi) MustApiGetPeersSuspended() []NodePeerSuspended {
	b,err:= sdk.httpGet("/peers/suspended")
	if err!=nil{
		panic("6npygzdqyk ["+err.Error()+"]")
	}
	//fmt.Println("MustApiGetPeersSuspended",string(b))
	var resp []NodePeerSuspended
	mustJsonUnmarshal(b,&resp)
	return resp
}

type NodePeerBlacklisted struct{
	Hostname            string `json:"hostname"`
	Timestamp int64  `json:"timestamp"`
	Reason    string `json:"reason"`
}

func (sdk *VsysApi) MustApiGetPeersBlacklisted() []NodePeerBlacklisted {
	b,err:= sdk.httpGet("/peers/blacklisted")
	if err!=nil{
		panic("zpxt95sj6n ["+err.Error()+"]")
	}
	//fmt.Println("MustApiGetPeersBlacklisted",string(b))
	var resp []NodePeerBlacklisted
	mustJsonUnmarshal(b,&resp)
	return resp
}

type NodeHeath struct{
	SelfAddr            string    `json:",omitempty"`
	AllAddrList   []string  `json:",omitempty"`
	SuspendedAddrList []string `json:",omitempty"`
	ConnectedAddrList   []string  `json:",omitempty"`
	BlacklistedAddrList []string  `json:",omitempty"`
	LastBlockHeight     int64     `json:",omitempty"`
	LastBlockSignature  string    `json:",omitempty"`
	LastBlockTimestamp  time.Time `json:",omitempty"`
	UnconfirmedSize int64 // 0 is valid value
	ErrMsg string `json:",omitempty"`
}

func (api *VsysApi) GetNodeHeath() NodeHeath{
	nh:=NodeHeath{}
	nh.SelfAddr = api.nodeAddress
	errMsg:=panicToErrorMsg(func(){
		lastBlock:=api.MustApiGetBlockLast()
		nh.LastBlockHeight = lastBlock.Height
		nh.LastBlockSignature = lastBlock.Signature
		nh.LastBlockTimestamp = time.Unix(0,lastBlock.Timestamp).UTC()

		pcList:=api.MustApiGetPeersAll()
		for _,pc:=range pcList{
			nh.AllAddrList = append(nh.AllAddrList,pc.Address)
		}
		pcList2:=api.MustApiGetPeersConnected()
		for _,pc:=range pcList2{
			nh.ConnectedAddrList = append(nh.ConnectedAddrList,pc.Address)
		}
		pbList2:=api.MustApiGetPeersSuspended()
		for _,pc:=range pbList2{
			nh.SuspendedAddrList = append(nh.SuspendedAddrList,pc.Hostname)
		}
		pbList:=api.MustApiGetPeersBlacklisted()
		for _,pc:=range pbList{
			nh.BlacklistedAddrList = append(nh.BlacklistedAddrList,pc.Hostname)
		}
		nh.UnconfirmedSize = api.MustGetTransactionUnconfirmedSize()
	})
	nh.ErrMsg = errMsg
	return nh
}

type NodeSyncRelation struct{
	LastNodeAddr string `json:",omitempty"`
	CheckNodeAddr string `json:",omitempty"`
	FindHeight int64 `json:",omitempty"`
	ErrMsg string `json:",omitempty"`
}
type NodeListHeathCheckResp struct{
	HealthList []NodeHeath
	SyncRelationList []NodeSyncRelation
	ProblemList []string
}
func NodeListHeathCheck(nodeAddrList []string,NetType NetType) (resp NodeListHeathCheckResp){
	SyncRelationM1:=map[string]NodeSyncRelation{}
	canNotConnectMap:=map[string]struct{}{}
	for _, addr1 :=range nodeAddrList {
		api:=NewVsysApi(NewVsysApiReq{
			NodeAddress: addr1,
			Network:     NetType,
		})
		nh:=api.GetNodeHeath()
		resp.HealthList = append(resp.HealthList,nh)
		if nh.ErrMsg!=""{
			canNotConnectMap[nh.SelfAddr] = struct{}{}
			resp.ProblemList = append(resp.ProblemList,"can not connect to ["+nh.SelfAddr+"]")
			continue
		}
		for _, addr2 :=range nodeAddrList{
			if addr2 == addr1 {
				continue
			}
			api2:=NewVsysApi(NewVsysApiReq{
				NodeAddress: addr2,
				Network:     NetType,
			})
			findHeight:=int64(0)
			errMsg:=panicToErrorMsg(func(){
				findHeight=api2.MustApiGetBlockHeightBySignature(nh.LastBlockSignature)
			})
			nsr:=NodeSyncRelation{
				LastNodeAddr: addr1,
				CheckNodeAddr: addr2,
				FindHeight: findHeight,
				ErrMsg: errMsg,
			}
			resp.SyncRelationList = append(resp.SyncRelationList,nsr)
			SyncRelationM1[nsr.LastNodeAddr+"_"+nsr.CheckNodeAddr] = nsr
		}
	}
	for _, addr1 :=range nodeAddrList {
		_,ok:=canNotConnectMap[addr1]
		if ok==true{
			continue
		}
		for _, addr2 :=range nodeAddrList {
			if addr2 >= addr1 {
				continue
			}
			_,ok:=canNotConnectMap[addr2]
			if ok==true{
				continue
			}
			nsr1:=SyncRelationM1[addr1+"_"+addr2]
			nsr2:=SyncRelationM1[addr2+"_"+addr1]
			if nsr1.FindHeight==0 && nsr2.FindHeight==0{
				resp.ProblemList = append(resp.ProblemList,"["+addr1+"]_["+addr2+"] is not synced")
			}
		}
	}
	return resp
}