package rpc

import (
	"Jight/dagchain"
	"Jight/transaction"
)

type Jightd struct{
	DC *dagchain.DagChain
}

type GetBalanceCMD struct {
	Address string
}

type GetTopicsCMD struct {
	Address string
}

type GetAccTopicsReply struct {
	Topics []string
}


type GetTransactionCMD struct {
	Hash string
}

type GetTransactionReply struct {
	Tx transaction.Transaction
}

type CompactTxCMD struct {

}

type CompactTxReply struct {
	Info string
}

type ListAddressCMD struct {
}

type ListAddressReply struct {
	Addresses []string
}

type SendCmd struct {
	From string
	To string
	Info string
}

type SendReply struct {
	SrField string
}

type CreateDagChainCMD struct {
	Address string
}

type CreateDagChainReply struct {
	CdcrField string
}

type CreateWalletCMD struct {
	Address string
}

type CreateWalletReply struct {
	CwrField string
}

// The following variables are just for test
var balancesList = map[string]string{
	"fsdjklfjskljfkldsjklfdjsklf": "123243",
	"fjiosdjfioqrewrewrewtrtrggg": "896969",
}

var addressList = []string {
	"fsdjklfjskljfkldsjklfdjsklf",
	"fjiosdjfioqrewrewrewtrtrggg",
	"r8i9w0jfsfjsdkfjdskfjh8r9ww",
}



