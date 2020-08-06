package goVsysSdk

import (
	"testing"
	"bytes"
)

func TestBase58(ot *testing.T){
	for i:=0;i<256;i++{
		in:=[]byte{byte(i)}
		outS:=Base58Encode(in)
		outB:= MustBase58Decode(outS)
		testOk(bytes.Equal(in,outB))
	}
	_,ok:= Base58Decode("\x00")
	testOk(ok==false)
}

func testOk(expectTrue bool){
	if expectTrue==true{
		return
	}
	panic("ok fail")
}