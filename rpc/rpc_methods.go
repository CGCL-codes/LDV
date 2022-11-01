package rpc

import (
	"Jight/dagchain"
	"Jight/p2p"
	"Jight/wallet"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
)

func (jd *Jightd) GetTopics (cmd GetTopicsCMD, reply *GetAccTopicsReply) error {
	account, err := dagchain.FetchAccount(jd.DC.DB, []byte(cmd.Address))
	if err != nil {
		return errors.New("no match address: " + cmd.Address)
	}

	//reply.Balance = strconv.FormatUint(uint64(balance), 10)
	for _, topic := range account.Topics {
		reply.Topics = append(reply.Topics, string(topic[:]))
	}

	return nil
}

func (jd *Jightd) Send (cmd SendCmd, reply *SendReply) error {

	// TODO nodeID should be mac or ip address, nodeID specialize a node
	nodeID := "self"
	ws, err := wallet.LoadWallets(nodeID)
	if err != nil {
		log.Panic(err)
	}

	wallet := ws.GetWallet(cmd.From)

	tx := jd.DC.CreateTx(wallet.PrivateKey, []byte(cmd.From), []byte(cmd.To), []byte(cmd.Info))

	reply.SrField = fmt.Sprintf("Send %d from %s to topic %s, hash of tx: %x", cmd.Info, cmd.From, cmd.To, tx.Hash)

	jd.DC.AddTx(tx, true)
	p2p.SyncTx(tx)
	return nil
}

func (jd *Jightd) CreateWallet (cmd CreateWalletCMD, reply *CreateWalletReply) error {
	address := cmd.Address
	// lastid of new account is nil
	acc := dagchain.Account{[]byte(address), nil, [][]byte{}}
	log.Println("func CreateWallet is called")
	return dagchain.StoreAccount(jd.DC.DB, acc)
}

func (jd *Jightd) GetTransaction(cmd GetTransactionCMD, reply *GetTransactionReply) error {
	hashByte, _ := hex.DecodeString(cmd.Hash)

	tx, err := dagchain.FetchTx(jd.DC.DB, hashByte)

	if err != nil {
		return errors.New("no tx of " + cmd.Hash)
	}

	reply.Tx = *tx

	return nil
}

func (jd *Jightd) CompactTx(cmd CompactTxCMD, reply *CompactTxReply) error {

	err := dagchain.CompactTx(jd.DC.DB)

	if err != nil {
		return errors.New("compact tx fails")
	}

	reply.Info = "get request"

	return nil
}
