package dagchain

import (
	"Jight/transaction"
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
	"log"
)

func StoreTxHash(db *leveldb.DB, th transaction.TxHash) error {
	blockdata := th.Serialize()
	_ = db.Put(bytes.Join([][]byte{[]byte("th"),th.Hash},[]byte{}),[]byte(blockdata),nil)
	log.Println("Store TxHash: ", hex.EncodeToString(th.Hash))
	return nil
}

// nq: store tx to db
func StoreTx(db *leveldb.DB, tx transaction.Transaction) error {
	blockdata := tx.Serialize()
	_ = db.Put(bytes.Join([][]byte{[]byte("t"),tx.Hash},[]byte{}),[]byte(blockdata),nil)
	log.Println("Store Tx: ", hex.EncodeToString(tx.Hash))
	return nil
}

// nq: get tx via hash from db
func FetchTx(db *leveldb.DB, hash []byte) (*transaction.Transaction, error) {
	data,err :=db.Get(bytes.Join([][]byte{[]byte("t"),hash},[]byte{}),nil)
	if err != nil {
		return nil, err
	}
	tx := transaction.DeserializeTx(data)
	return &tx, nil
}

// Delete a specified tip with tip tx hash
func RemoveTx(db *leveldb.DB, hash []byte) error {
	err := db.Delete(bytes.Join([][]byte{[]byte("t"),hash},[]byte{}),nil)
	if err != nil {
		return errors.New("remove tx failed")
	}
	return nil
}

// nq: store account to db
func StoreAccount(db *leveldb.DB, acc Account) error {
	Accountdata := acc.Serialize()
	key := bytes.Join([][]byte{[]byte("a"),acc.Account},[]byte{})
	value := []byte(Accountdata)
	log.Println("Store key: ", key)
	if err := db.Put(key ,value,nil); err!=nil {
		log.Fatal(err)
		return err
	}
	return nil
}

// nq: get account via address from db
func FetchAccount(db *leveldb.DB, addr []byte) (*Account, error) {
	log.Println("Before db.get in FetchAccount")
	key := bytes.Join([][]byte{[]byte("a"),addr},[]byte{})
	fmt.Println("Fetch key: ", key)
	data, err := db.Get(key,nil)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	log.Println("After db.get in FetchAccount")
	acc := DeserializeAcc(data)
	return &acc, nil
}

// nq: store all tips to db
func StoreTips(db *leveldb.DB, tips []*Tip) error {
	for i := 0; i < len(tips); i++ {
		log.Println("Store tip: ", hex.EncodeToString(tips[i].TipID))
		data := tips[i].Serialize()
		_ = db.Put(bytes.Join([][]byte{[]byte("tip"),tips[i].TipID},[]byte{}),[]byte(data),nil)
	}
	return nil
}

// Add a new tip
func PutTip(db *leveldb.DB, tip Tip) error {
	Tipdata := tip.Serialize()
	err := db.Put(bytes.Join([][]byte{[]byte("tip"),tip.TipID},[]byte{}),[]byte(Tipdata),nil)
	if err != nil {
		return errors.New("add Tip failed")
	}
	return nil
}

// Delete a specified tip with tip tx hash
func RemoveTip(db *leveldb.DB, tipID []byte) error {
	err := db.Delete(bytes.Join([][]byte{[]byte("tip"),tipID},[]byte{}),nil)
	if err != nil {
		return errors.New("remove tip failed")
	}
	return nil
}

// nq: get all tips from db
func FetchTips(db *leveldb.DB) ([]*Tip, error) {
	var temptips []*Tip
	iter := db.NewIterator(util.BytesPrefix([]byte("tip")),nil)
	for iter.Next() {
		data := DeserializeTip(iter.Value())
		temptips = append(temptips, &data)
	}
	iter.Release()
	tips := temptips[:]
	return tips,nil
}

func CompactTx(db *leveldb.DB) error {
	txRange := util.BytesPrefix([]byte("t"))
	return db.CompactRange(*txRange)
}


// open db instance or create a db if not exist
func LoadDB() (*leveldb.DB, error) {
	db, err := leveldb.OpenFile("Jightdb",nil)
	if err != nil {
		return nil, err
	}
	return db, nil
}