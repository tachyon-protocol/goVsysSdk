package goVsysSdk

import (
	"encoding/json"
	"fmt"
)

func mustJsonUnmarshal(b []byte,obj interface{}){
	err := json.Unmarshal(b, obj)
	if err != nil {
		panic(err)
	}
}

const hexTable = "0123456789ABCDEF"
func urlv(s string) string {
	afterLen:=len(s)
	for i:=0;i<len(s);i++{
		c := s[i]
		if ('A' <= c && c <= 'Z') || ('a' <= c && c <= 'z') || ('0' <= c && c <= '9') || c == '-' || c == '.' || c == '_' {
		} else {
			afterLen+=2
		}
	}
	if afterLen==len(s){
		return s
	}
	out := make([]byte, afterLen)
	outPos:=0
	for i := 0; i < len(s); i++ {
		c := s[i]
		if ('A' <= c && c <= 'Z') || ('a' <= c && c <= 'z') || ('0' <= c && c <= '9') || c == '-' || c == '.' || c == '_' {
			out[outPos] = c
			outPos++
		} else {
			out[outPos] = '%'
			out[outPos+1] = hexTable[c>>4]
			out[outPos+2] = hexTable[c&15]
			outPos+=3
		}
	}
	return string(out)
}

func panicIfError(err error) {
	if err != nil {
		panic(err)
	}
}

func panicToErrorMsg(f func()) (errMsg string) {
	hasFinish := false
	defer func() {
		if hasFinish {
			return
		}
		out := recover()
		errMsg = interfaceToStringNotEmpty(out)
	}()
	f()
	hasFinish = true
	return errMsg
}

func interfaceToStringNotEmpty(i interface{}) (outS string) {
	switch out := i.(type) {
	case error:
		outS = out.Error()
	case string:
		outS = out
	default:
		return fmt.Sprintf("%#v",out)
	}
	if outS==""{
		outS = `<"">`
	}
	return outS
}