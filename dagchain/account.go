package dagchain

import (
	"bytes"
	"encoding/gob"
	"log"
)

type Account struct {
	Account []byte // address of accounts
	LastId  []byte // hash id of Account's last tx
	//Balance uint32 // Balance of this Account
	Topics [][]byte
}


// serialize Account
func (acc Account) Serialize() []byte {
	var encode bytes.Buffer

	enc := gob.NewEncoder(&encode)
	err := enc.Encode(acc)

	if err != nil {
		log.Panic("Account encode fail:", err)
	}

	return encode.Bytes()
}

// Deserialize Account
func DeserializeAcc(data []byte) Account {
	var acc Account

	decode := gob.NewDecoder(bytes.NewReader(data))

	err := decode.Decode(&acc)
	if err != nil {
		log.Panic("Account decode fail:", err)
	}

	return acc
}