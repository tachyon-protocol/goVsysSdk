package goVsysSdk

import (
	"encoding/binary"
	"strconv"
)

type DataEncoder struct {
	result []byte
}

func (de *DataEncoder) WriteContractAccount(contractId string){
	bytes := MustBase58Decode(contractId)
	if len(bytes)!=26{
		panic("WriteContractAccount fail1 "+strconv.Itoa(len(bytes)))
	}
	de.result = append(de.result, DeTypeContractAccount)
	de.result = append(de.result, bytes...)
}

func (de *DataEncoder) WriteShortByte(bytes []byte){
	de.result = append(de.result, DeTypeShortByte)
	de.result = append(de.result, bytesToByteArrayWithSize(bytes)...)
}

func (de *DataEncoder) WriteAddress(address string){
	bytes := MustBase58Decode(address)
	if len(bytes)!=26{
		panic("WriteAccount fail1 "+strconv.Itoa(len(bytes)))
	}
	de.result = append(de.result, DeTypeAddress)
	de.result = append(de.result, bytes...)
}

func (de *DataEncoder) WriteAmount(amount int64){
	bytes := uint64ToByte(uint64(amount))
	de.result = append(de.result, DeTypeAmount)
	de.result = append(de.result, bytes...)
}

func (de *DataEncoder) EncodeArgAmount(amount int16) {
	de.result = append(de.result, uint16ToByte(uint16(amount))...)
}

func (de *DataEncoder) Encode(data interface{}, dataEntryType byte) {
	switch dataEntryType {
	case DeTypeContractAccount:
		de.WriteContractAccount(data.(string))
		return
	//case DeTypeAccount:
	//	de.WriteAccount(data.(string))
	//	return
	case DeTypePublicKey, DeTypeAddress, DeTypeTokenId:
		bytes := MustBase58Decode(data.(string))
		de.result = append(de.result, dataEntryType)
		de.result = append(de.result, bytes...)
	case DeTypeAmount:
		de.WriteAmount(data.(int64))
		return
	case DeTypeInt32:
		bytes := uint32ToByte(data.(int32))
		de.result = append(de.result, dataEntryType)
		de.result = append(de.result, bytes...)
	case DeTypeShortText:
		bytes := []byte(data.(string))
		de.result = append(de.result, dataEntryType)
		de.result = append(de.result, bytesToByteArrayWithSize(bytes)...)
	case DeTypeShortByte:
		de.WriteShortByte(data.([]byte))
		return
	case DeTypeTimeStamp:
		bytes := uint64ToByte(uint64(data.(int64)))
		de.result = append(de.result, dataEntryType)
		de.result = append(de.result, bytes...)
	default:
	}
}

func (de *DataEncoder) Decode(data []byte) (list []DataEntry) {
	for i := 2; i < len(data); {
		deType := data[i]
		i++
		switch deType {
		case DeTypePublicKey:
			list = append(list, DataEntry{
				Type:  int8(deType),
				Value: Base58Encode(data[i : i+32]),
			})
			i = i + 32
		case DeTypeTokenId:
			list = append(list, DataEntry{
				Type:  int8(deType),
				Value: Base58Encode(data[i : i+30]),
			})
			i = i + 30
		case DeTypeAddress, DeTypeContractAccount:
			list = append(list, DataEntry{
				Type:  int8(deType),
				Value: Base58Encode(data[i : i+26]),
			})
			i = i + 26
		case DeTypeAmount:
			list = append(list, DataEntry{
				Type:  int8(deType),
				Value: int64(binary.BigEndian.Uint64(data[i : i+8])),
			})
			i = i + 8
		case DeTypeInt32:
			list = append(list, DataEntry{
				Type:  int8(deType),
				Value: int64(binary.BigEndian.Uint32(data[i : i+4])),
			})
			i = i + 4
		case DeTypeShortText:
			length := int(binary.BigEndian.Uint16(data[i : i+2]))
			i = i + 2
			list = append(list, DataEntry{
				Type:  int8(deType),
				Value: string(data[i : i+length]),
			})
			i = i + length
		case DeTypeShortByte:
			length := int(binary.BigEndian.Uint16(data[i : i+2]))
			i = i + 2
			list = append(list, DataEntry{
				Type:  int8(deType),
				Value: []byte(data[i : i+length]),
			})
			i = i + length
		}
	}
	return list
}

func (de *DataEncoder) Result() []byte {
	return de.result
}