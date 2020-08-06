goVsysSdk
=============

[![GoDoc](https://godoc.org/github.com/tachyon-protocol/goVsysSdk?status.svg)](https://godoc.org/github.com/tachyon-protocol/goVsysSdk)

The golang library for V Systems Blockchain.

## Warning

The full functionality of SDK is still under development. The API may not be stable, please use it at your own risk.

## Installing

Use `go get` to retrieve the SDK sto add it to your `GOPATH` workspace, or
project's Go module dependencies.

	go get github.com/tachyon-protocol/goVsysSdk

### Dependencies

The metadata of the SDK's dependencies can be found in the Go module file `go.mod`

### usage example
* Notice you need 202.0v in TestNetApi to success run this example.
    * Change the seed to your seed in TestNetApi. or make the balance of AUCvntTpQo39XBaGrx26459RZJsragYBCSj(sender address in the example) more than 202.0v.
```go
package main

import (
	"github.com/tachyon-protocol/goVsysSdk"
	"fmt"
	"time"
)

func main(){
	const seed = "test_qpbw2c2dmru97k8fptrp4kyj768cg7y99jn8rzk5dbydgxchfebnpkayjxqybmhq"
	api:=goVsysSdk.NewPublicTestNetApi()
	sender :=api.NewAccountFromSeedAndNonceV2(seed,0)
	receiver :=api.NewAccountFromSeedAndNonceV2(seed,1)
	senderVsys:=api.MustGetAccountBalance(sender.GetAddress())
	if senderVsys<202.0*goVsysSdk.VsysAmountRate{
		fmt.Println("WARNGING: you need at least 202.0v to finish this example.")
		fmt.Println("=========================================================")
	}
	fmt.Println("example of send vsys to another account.")
	fmt.Println("=========================================================")
	fmt.Println("sender", sender.GetAddress(),api.MustGetAccountBalanceString(sender.GetAddress()))
	fmt.Println("receiver", receiver.GetAddress(),api.MustGetAccountBalanceString(receiver.GetAddress()))
	fmt.Println("height",api.MustGetBlockHeight())
	fmt.Println("start send 1 vsys from sender to receiver.")
	api.MustSendPaymentSimpleSync(sender, receiver.GetAddress(),1*goVsysSdk.VsysAmountRate)
	fmt.Println("after 1 vsys from sender to receiver.")
	fmt.Println("sender", sender.GetAddress(),api.MustGetAccountBalanceString(sender.GetAddress()))
	fmt.Println("receiver", receiver.GetAddress(),api.MustGetAccountBalanceString(receiver.GetAddress()))
	fmt.Println("=========================================================")

	fmt.Println("example of create a token contract")
	fmt.Println("=========================================================")
	resp1 := api.MustRegisterContractToken(goVsysSdk.RegisterContractTokenReq{
		Sender:              sender,
		Max:                 1e9,
		Unity:               1e9,
		TokenDescription:    "td_2f8w2zwzj6",
		ContractDescription: "cd_svyh2c9ax4",
	})
	api.MustWaitPaymentOkByTransactionResponse(resp1)
	TokenContractId := resp1.ContractId
	tokenId:= goVsysSdk.ContractId2TokenId(resp1.ContractId,0)
	fmt.Println("after create a token contract, tokenId",tokenId,"TokenContractId:",TokenContractId)
	fmt.Println("sender vsys", sender.GetAddress(),api.MustGetAccountBalanceString(sender.GetAddress()),"token",api.MustGetContractTokenBalance(sender.GetAddress(),tokenId))
	fmt.Println("=========================================================")

	fmt.Println("example of token issue for tokenId",tokenId)
	fmt.Println("=========================================================")
	resp2:=api.MustTokenIssue(goVsysSdk.TokenIssueReq{
		Sender:  sender,
		TokenId: tokenId,
		Amount:  int64(1e9),
	})
	api.MustWaitPaymentOkByTransactionResponse(resp2)
	fmt.Println("after token issue for tokenId",tokenId)
	fmt.Println("sender vsys", sender.GetAddress(),api.MustGetAccountBalanceString(sender.GetAddress()),"token",api.MustGetContractTokenBalance(sender.GetAddress(),tokenId))
	fmt.Println("=========================================================")

	fmt.Println("example of create a PaymentChannel contract for tokenId",tokenId)
	fmt.Println("=========================================================")
	resp3 := api.MustRegisterContractPaymentChannel(goVsysSdk.RegisterContractPaymentChannelReq{
		Sender:        sender,
		Vsys_token_id: tokenId,
	})
	api.MustWaitPaymentOkByTransactionResponse(resp3)
	PaymentChannelContractId := resp3.ContractId
	fmt.Println("after create a PaymentChannel contract, PaymentChannelContractId",PaymentChannelContractId)
	fmt.Println("sender vsys", sender.GetAddress(),api.MustGetAccountBalanceString(sender.GetAddress()),"token",api.MustGetContractTokenBalance(sender.GetAddress(),tokenId))
	fmt.Println("=========================================================")

	fmt.Println("example of token deposit from ",tokenId," to PaymentChannelContractId",PaymentChannelContractId)
	fmt.Println("=========================================================")
	resp4 := api.MustTokenDeposit(goVsysSdk.TokenDepositReq{
		Sender:            sender,
		TokenId:           tokenId,
		ReceiveAddress:    sender.GetAddress(),
		ReceiveContractId: PaymentChannelContractId,
		Amount:            3,
	})
	api.MustWaitPaymentOkByTransactionResponse(resp4)
	fmt.Println("after token deposit")
	fmt.Println("sender vsys", sender.GetAddress(),api.MustGetAccountBalanceString(sender.GetAddress()),"token",api.MustGetContractTokenBalance(sender.GetAddress(),tokenId))
	fmt.Println("=========================================================")

	fmt.Println("example of PaymentChannelCreateAndLoad PaymentChannelContractId:",PaymentChannelContractId,"receiver address",receiver.GetAddress())
	fmt.Println("=========================================================")
	resp5 := api.MustPaymentChannelCreateAndLoad(goVsysSdk.PaymentChannelCreateAndLoadReq{
		Sender:           sender,
		ContractId:       PaymentChannelContractId,
		RecipientAddress: receiver.GetAddress(),
		Amount:           3,
		TimeStamp:        time.Now().Add(time.Hour * 24).UnixNano(),
	})
	api.MustWaitPaymentOkByTransactionResponse(resp5)
	channelId := resp5.Id
	fmt.Println("after PaymentChannelCreateAndLoad channelId",channelId)
	fmt.Println("sender vsys", sender.GetAddress(),api.MustGetAccountBalanceString(sender.GetAddress()),"token",api.MustGetContractTokenBalance(sender.GetAddress(),tokenId))
	fmt.Println("=========================================================")

	fmt.Println("example of PaymentChannelGenerateSenderPaymentSignature (no api call) channelId",channelId)
	fmt.Println("=========================================================")
	paymentS := goVsysSdk.PaymentChannelGenerateSenderPaymentSignature(sender, channelId, 1)
	fmt.Println("after PaymentChannelGenerateSenderPaymentSignature channelId",channelId,"paymentS",paymentS)
	fmt.Println("=========================================================")

	fmt.Println("example of MustPaymentChannelCollect channelId",channelId)
	fmt.Println("=========================================================")
	resp6 := api.MustPaymentChannelCollect(goVsysSdk.PaymentChannelCollectReq{
		Receiver:              receiver,
		ContractId:           PaymentChannelContractId,
		ChannelId:             channelId,
		Amount:                1,
		Payment_signature_str: paymentS,
	})
	api.MustWaitPaymentOkByTransactionResponse(resp6)
	fmt.Println("after MustPaymentChannelCollect channelId",channelId,"paymentS",paymentS)
	fmt.Println("receiver vsys", receiver.GetAddress(),api.MustGetAccountBalanceString(receiver.GetAddress()),"token",api.MustGetContractTokenBalance(receiver.GetAddress(),tokenId))
	fmt.Println("=========================================================")

	fmt.Println("example of MustTokenWithdraw tokenId",tokenId,"receiver",receiver.GetAddress(),"PaymentChannelContractId",PaymentChannelContractId)
	fmt.Println("=========================================================")
	resp7 := api.MustTokenWithdraw(goVsysSdk.TokenWithdrawReq{
		Sender:            receiver,
		TokenId:           tokenId,
		ReceiveAddress:    receiver.GetAddress(),
		ReceiveContractId: PaymentChannelContractId,
		Amount:            1,
	})
	api.MustWaitPaymentOkByTransactionResponse(resp7)
	fmt.Println("after MustTokenWithdraw channelId",channelId,"paymentS",paymentS)
	fmt.Println("receiver vsys", receiver.GetAddress(),api.MustGetAccountBalanceString(receiver.GetAddress()),"token",api.MustGetContractTokenBalance(receiver.GetAddress(),tokenId))
	fmt.Println("=========================================================")
	fmt.Println("sender vsys", sender.GetAddress(),api.MustGetAccountBalanceString(sender.GetAddress()),"token",api.MustGetContractTokenBalance(sender.GetAddress(),tokenId))
}
```

output of above example
```
example of send vsys to another account.
=========================================================
sender AUCvntTpQo39XBaGrx26459RZJsragYBCSj 238.7v
receiver AU3CJrTHUEL386t25tq3AzL2r9EqNnRKmZM 3.4v
height 13682383
start send 1 vsys from sender to receiver.
after 1 vsys from sender to receiver.
sender AUCvntTpQo39XBaGrx26459RZJsragYBCSj 237.6v
receiver AU3CJrTHUEL386t25tq3AzL2r9EqNnRKmZM 4.4v
=========================================================
example of create a token contract
=========================================================
after create a token contract, tokenId TWu5D99q1cnZS1FWSrTsUtgcEaawoHXu8QqHeBLHq TokenContractId: CF7EVuokPpqYy6JEH9JR2BzhUAaGZv2CeKw
sender vsys AUCvntTpQo39XBaGrx26459RZJsragYBCSj 137.6v token 0
=========================================================
example of token issue for tokenId TWu5D99q1cnZS1FWSrTsUtgcEaawoHXu8QqHeBLHq
=========================================================
after token issue for tokenId TWu5D99q1cnZS1FWSrTsUtgcEaawoHXu8QqHeBLHq
sender vsys AUCvntTpQo39XBaGrx26459RZJsragYBCSj 137.3v token 1000000000
=========================================================
example of create a PaymentChannel contract for tokenId TWu5D99q1cnZS1FWSrTsUtgcEaawoHXu8QqHeBLHq
=========================================================
after create a PaymentChannel contract, PaymentChannelContractId CF4Fya2CyAdaMrxtBMzydTKD2dX2kJ9f7Jv
sender vsys AUCvntTpQo39XBaGrx26459RZJsragYBCSj 37.3v token 1000000000
=========================================================
example of token deposit from  TWu5D99q1cnZS1FWSrTsUtgcEaawoHXu8QqHeBLHq  to PaymentChannelContractId CF4Fya2CyAdaMrxtBMzydTKD2dX2kJ9f7Jv
=========================================================
after token deposit
sender vsys AUCvntTpQo39XBaGrx26459RZJsragYBCSj 37v token 999999997
=========================================================
example of PaymentChannelCreateAndLoad PaymentChannelContractId: CF4Fya2CyAdaMrxtBMzydTKD2dX2kJ9f7Jv receiver address AU3CJrTHUEL386t25tq3AzL2r9EqNnRKmZM
=========================================================
after PaymentChannelCreateAndLoad channelId DA1CLjfJunKwqW6AhDrWJHZbVoiu7eAEyVAFtnpdZCwQ
sender vsys AUCvntTpQo39XBaGrx26459RZJsragYBCSj 36.7v token 999999997
=========================================================
example of PaymentChannelGenerateSenderPaymentSignature (no api call) channelId DA1CLjfJunKwqW6AhDrWJHZbVoiu7eAEyVAFtnpdZCwQ
=========================================================
after PaymentChannelGenerateSenderPaymentSignature channelId DA1CLjfJunKwqW6AhDrWJHZbVoiu7eAEyVAFtnpdZCwQ paymentS 7dxKicDVz11hSy77gGxWBWfGYKQhHyXmbKt4wX36yvYzBpwLHrhj8ze4N6wZRrN69byBnCZgHsCFEE8tHVD6bto
=========================================================
example of MustPaymentChannelCollect channelId DA1CLjfJunKwqW6AhDrWJHZbVoiu7eAEyVAFtnpdZCwQ
=========================================================
after MustPaymentChannelCollect channelId DA1CLjfJunKwqW6AhDrWJHZbVoiu7eAEyVAFtnpdZCwQ paymentS 7dxKicDVz11hSy77gGxWBWfGYKQhHyXmbKt4wX36yvYzBpwLHrhj8ze4N6wZRrN69byBnCZgHsCFEE8tHVD6bto
receiver vsys AU3CJrTHUEL386t25tq3AzL2r9EqNnRKmZM 4.1v token 0
=========================================================
example of MustTokenWithdraw tokenId TWu5D99q1cnZS1FWSrTsUtgcEaawoHXu8QqHeBLHq receiver AU3CJrTHUEL386t25tq3AzL2r9EqNnRKmZM PaymentChannelContractId CF4Fya2CyAdaMrxtBMzydTKD2dX2kJ9f7Jv
=========================================================
after MustTokenWithdraw channelId DA1CLjfJunKwqW6AhDrWJHZbVoiu7eAEyVAFtnpdZCwQ paymentS 7dxKicDVz11hSy77gGxWBWfGYKQhHyXmbKt4wX36yvYzBpwLHrhj8ze4N6wZRrN69byBnCZgHsCFEE8tHVD6bto
receiver vsys AU3CJrTHUEL386t25tq3AzL2r9EqNnRKmZM 3.8v token 1
=========================================================
sender vsys AUCvntTpQo39XBaGrx26459RZJsragYBCSj 36.7v token 999999997
```

### License
  This package is licensed under the Unlicense. See LICENSE for details.

### Notice
* public main net explorer https://explorer.v.systems/
* public test net explorer https://testexplorer.v.systems/
* fork from https://github.com/walkbean/vsys-sdk-go , change a lot api from that version.