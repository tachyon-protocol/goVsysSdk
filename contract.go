package goVsysSdk

import (
	"encoding/binary"
	"encoding/json"
)

type Contract struct {
	ContractId       string
	Contract         string
	Max              int64
	Unity            int64
	TokenDescription string
	ContractDescription string
	Amount           int64
	TokenIdx         int32 // const 0
	Recipient        string
	SenderPublicKey  string
	NewUnity         int64  // split newUnity
	NewIssuer        string // supersede newIssuer

	Textual   Textual // [init func, user defined func, stateVar]
	Functions []Func
}

type Textual struct {
	Triggers       string
	Descriptors    string
	StateVariables string
}

type Func struct {
	Name    string
	Args    []string
	RetArgs []string
}

const (
	DeTypePublicKey       = 0x01
	DeTypeAddress         = 0x02
	DeTypeAmount          = 0x03
	DeTypeInt32           = 0x04
	DeTypeShortText       = 0x05
	DeTypeContractAccount = 0x06
	//DeTypeAccount         = 0x07 // Account is not a data entry, please use DeTypeAddress or DeTypeContractAccount
	DeTypeTokenId         = 0x08
	DeTypeTimeStamp       = 0x09
	DeTypeBool            = 10
	DeTypeShortByte       = 11
)

func (c *Contract) BuildRegisterData() []byte {
	data := DataEncoder{}
	data.EncodeArgAmount(3)
	data.Encode(c.Max, DeTypeAmount)
	data.Encode(c.Unity, DeTypeAmount)
	data.Encode(c.TokenDescription, DeTypeShortText)

	return data.Result()
}

func (c *Contract) BuildIssueData() []byte {
	data := DataEncoder{}
	data.EncodeArgAmount(1)
	data.Encode(c.Amount, DeTypeAmount)
	return data.Result()
}

func (c *Contract) BuildSendData() []byte {
	data := DataEncoder{}
	data.EncodeArgAmount(2)
	data.Encode(c.Recipient, DeTypeAddress)
	data.Encode(c.Amount, DeTypeAmount)

	return data.Result()
}

func (c *Contract) BuildSplitData() []byte {
	data := DataEncoder{}
	data.EncodeArgAmount(1)
	data.Encode(c.NewUnity, DeTypeAmount)

	return data.Result()
}

func (c *Contract) BuildDestroyData() []byte {
	data := DataEncoder{}
	data.EncodeArgAmount(1)
	data.Encode(c.Amount, DeTypeAmount)

	return data.Result()
}

func (c *Contract) DecodeRegister(data []byte) {
	de := DataEncoder{}
	list := de.Decode(data)
	c.Max = list[0].Value.(int64)
	c.Unity = list[1].Value.(int64)
	c.TokenDescription = list[2].Value.(string)
}

func (c *Contract) DecodeIssue(data []byte) {
	de := DataEncoder{}
	list := de.Decode(data)
	c.Amount = list[0].Value.(int64)
}

func (c *Contract) DecodeDestroy(data []byte) {
	de := DataEncoder{}
	list := de.Decode(data)
	c.Amount = list[0].Value.(int64)
}

func (c *Contract) DecodeSend(data []byte) {
	de := DataEncoder{}
	list := de.Decode(data)
	c.Recipient = list[0].Value.(string)
	c.Amount = list[1].Value.(int64)
}

func (c *Contract) DecodeSplit(data []byte) {
	de := DataEncoder{}
	list := de.Decode(data)
	c.NewUnity = list[0].Value.(int64)
}

func (c *Contract) DecodeSupersede(data []byte) {
	de := DataEncoder{}
	list := de.Decode(data)
	c.NewIssuer = list[0].Value.(string)
}

func (c *Contract) DecodeTexture() {
	c.Functions = append(c.Functions, decodeFunc(c.Textual.Descriptors)...)
}

func decodeFunc(data string) []Func {
	bytes := MustBase58Decode(data)
	if len(bytes) < 2 {
		return []Func{}
	}
	numFunc := int(binary.BigEndian.Uint16(bytes[0:2]))
	p := 2
	var funcs []Func
	for i := 0; i < numFunc; i++ {
		p = p + 2
		funcNameLen := int(binary.BigEndian.Uint16(bytes[p : p+2]))
		p = p + 2
		funcName := string(bytes[p : p+funcNameLen])
		p = p + funcNameLen
		p = p + 2
		retArgCount := int(binary.BigEndian.Uint16(bytes[p : p+2]))
		p = p + 2
		var retArgs []string
		for j := 0; j < retArgCount; j++ {
			retArgLen := int(binary.BigEndian.Uint16(bytes[p : p+2]))
			retArgs = append(retArgs, string(bytes[p+2:p+2+retArgLen]))
			p = p + 2 + retArgLen
		}
		argsCount := int(binary.BigEndian.Uint16(bytes[p : p+2]))
		p = p + 2
		var args []string
		for argIndex := 0; argIndex < argsCount; argIndex++ {
			argNameLen := int(binary.BigEndian.Uint16(bytes[p : p+2]))
			args = append(args, string(bytes[p+2:p+2+argNameLen]))
			p = p + 2 + argNameLen
		}
		funcs = append(funcs, Func{
			Name:    funcName,
			Args:    args,
			RetArgs: retArgs,
		})
	}
	return funcs
}

type DataEntry struct {
	Type  int8
	Value interface{}
}



func DecodeContractTexture(data string) string {
	funcs := decodeFunc(data)
	res, _ := json.Marshal(funcs)
	return string(res)
}

func ContractId2TokenId(contractId string, tokenIndex int) string {
	if contractId == "" {
		return ""
	}
	contractIdBytes := MustBase58Decode(contractId)
	bytes := append([]byte{132}, contractIdBytes[1:len(contractIdBytes)-4]...)
	bytes = append(bytes, uint32ToByte(int32(tokenIndex))...)
	checksum := HashChain(bytes)[:4]
	return Base58Encode(append(bytes, checksum...))
}

func TokenId2ContractId(tokenId string) string {
	if tokenId == "" {
		return ""
	}
	tokenIdBytes := MustBase58Decode(tokenId)
	bytes := append([]byte{6}, tokenIdBytes[1:len(tokenIdBytes)-8]...)
	checksum := HashChain(bytes)[:4]
	return Base58Encode(append(bytes, checksum...))
}

