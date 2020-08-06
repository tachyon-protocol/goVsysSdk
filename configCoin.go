package goVsysSdk

import (
	"strconv"
)

const TokenSendFeeVsys = 0.3
const IpxContractId = "CC8Jx8aLkKVQmzuHBWNnhCSkn1GBLcjZ32k"
const IpxTokenId = "TWZZfKFqcaNVe5TrphLRNEm5DQFnBRJMjDDByqv84"
const IpxAmountRate = 1e9
const VsysAmountRate = 1e8

func FormatIpx(ipx int64)string{
	return strconv.FormatFloat(float64(ipx)/float64(IpxAmountRate), 'f', -1, 64)+" ipx"
}

func FormatVsys(v int64) string{
	return strconv.FormatFloat(float64(v)/1e8, 'f', -1, 64)+"v"
}