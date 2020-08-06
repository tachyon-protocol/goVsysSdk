package goVsysSdk

import (
	"strconv"
	"sort"
	"encoding/json"
	"strings"
)

func (sdk *VsysApi) GetBlockHeight() (h int,errMsg string){
	b,err:= sdk.httpGet("/blocks/height")
	if err!=nil{
		return 0,"qkw9bg6ufq ["+err.Error()+"]"
	}
	var respObj struct{
		Height int `json:"height"`
	}
	err = json.Unmarshal(b, &respObj)
	if err!=nil{
		return 0,"69dszw2fx4 ["+err.Error()+"]"
	}
	return respObj.Height,""
}

func (sdk *VsysApi) MustGetBlockHeight() int{
	b,err:= sdk.httpGet("/blocks/height")
	if err!=nil{
		panic("qkw9bg6ufq ["+err.Error()+"]")
	}
	var respObj struct{
		Height int `json:"height"`
	}
	mustJsonUnmarshal(b,&respObj)
	return respObj.Height
}

func (sdk *VsysApi) MustApiCallBlockSeq(from int,to int) (output []Block){
	b,err:= sdk.httpGet("/blocks/seq/"+urlv(strconv.Itoa(from))+"/"+urlv(strconv.Itoa(to)))
	if err!=nil{
		panic("qkw9bg6ufq ["+err.Error()+"]")
	}
	mustJsonUnmarshal(b,&output)
	return output
}

func (sdk *VsysApi) MustGetBlockCallbackDesc(cb func(block Block) bool) {
	thisHeight:=sdk.MustGetBlockHeight()
	const OneBatchNum = 99
	for{
		low:=thisHeight-OneBatchNum
		if low<=0{
			low = 0
		}
		blockList:=sdk.MustApiCallBlockSeq(low,thisHeight)
		sort.Slice(blockList,func(i int,j int)bool{
			return blockList[i].Height>blockList[j].Height
		})
		for _,block:=range blockList{
			isNext:=cb(block)
			if isNext==false{
				return
			}
		}
		thisHeight = low
		if thisHeight<=0{
			break
		}
	}
	return
}

type Block struct {
	Version       int    `json:"version"`
	Timestamp     int64  `json:"timestamp"`
	Reference     string `json:"reference"`
	SPOSConsensus struct {
		MintTime    int64 `json:"mintTime"`
		MintBalance int64 `json:"mintBalance"`
	} `json:"SPOSConsensus"`
	ResourcePricingData struct {
		Computation  int64 `json:"computation"`
		Storage      int64 `json:"storage"`
		Memory       int64 `json:"memory"`
		RandomIO     int64 `json:"randomIO"`
		SequentialIO int64 `json:"sequentialIO"`
	} `json:"resourcePricingData"`
	TransactionMerkleRoot string `json:"TransactionMerkleRoot"`
	Transactions          []HistoryTransaction `json:"transactions"`
	Generator        string `json:"generator"`
	Signature        string `json:"signature"`
	Fee              int64    `json:"fee"`
	Blocksize        int64    `json:"blocksize"`
	Height           int64    `json:"height"`
	TransactionCount int64    `json:"transaction count"`
}

func (sdk *VsysApi) MustApiGetBlockLast() Block{
	b,err:= sdk.httpGet("/blocks/last")
	if err!=nil{
		panic("g9ggn6pamr ["+err.Error()+"]")
	}
	var respObj Block
	mustJsonUnmarshal(b,&respObj)
	return respObj
}

func (sdk *VsysApi) MustApiGetBlockHeightBySignature(signature string) int64{
	b,err:= sdk.httpGet("/blocks/height/"+urlv(signature))
	if err!=nil{
		errMsg:=err.Error()
		if strings.Contains(errMsg,"block does not exist"){
			return 0
		}
		panic("8jby83bsa3 ["+errMsg+"]")
	}
	var respObj struct{
		Height int64 `json:"height"`
	}
	err = json.Unmarshal(b, &respObj)
	if err!=nil{
		panic("s6mqh3s93w ["+err.Error()+"]")
	}
	return respObj.Height
	//var respObj Block
	//mustJsonUnmarshal(b,&respObj)
	//return respObj
}