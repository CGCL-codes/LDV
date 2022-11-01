package transaction

import (
	"Jight/config"
	"bytes"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"encoding/gob"
	"fmt"
	"log"

	"math/big"
	"time"
)

// Transaction struct
type Transaction struct {
	Hash []byte			// Hash of the tx
	Parent []byte		// parent hash ref
	Validate [][]byte	// validate hash ref
	Sender []byte		// sender public key not address
	//Value uint32		// the money of transformation
	Info []byte		// the money of transformation
	Topic []byte		// receiver address
	Nonce uint32		// nonce of mining
	Timestamp int64		// Timestamp of tx
	Signature []byte	// Signature of tx
}

type TxHash struct {
	Hash []byte
	Parent []byte
	Validate [][]byte
}

func (th TxHash) Serialize() []byte {
	var encode bytes.Buffer

	enc := gob.NewEncoder(&encode)
	err := enc.Encode(th)

	if err != nil {
		log.Panic("tx encode fail:", err)
	}

	return encode.Bytes()
}

func (tx Transaction) ToTxHash() *TxHash {
	txHash := &TxHash{Hash: []byte(tx.Hash), Parent: []byte(tx.Parent), Validate: [][]byte(tx.Validate)}

	return txHash
}

// TODO mining
func (tx Transaction) Pow() uint32 {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-config.POW_TARGET_BITS))

	/*var hashInt big.Int
	var hash [32]byte
	var nonce uint32 = 0

	for nonce < config.MAX_NONCE {
		data := getRawData(tx, nonce)

		hash = sha256.Sum256(data)

		hashInt.SetBytes(hash[:])

		if hashInt.Cmp(target) == -1 {
			break
		} else {
			nonce++
		}
	}

	return nonce*/
	return 0
}

// return the raw data of tx need in pow
func getRawData(tx Transaction, nonce uint32) []byte {
	return bytes.Join(
		[][]byte{
			tx.Parent,
			bytes.Join(tx.Validate, []byte{}),
			tx.Sender,
			tx.Info,
			tx.Topic,
			IntToBytes(int64(nonce)),
			IntToBytes(tx.Timestamp),
			tx.Signature,
		},
		[]byte{})
}

func IntToBytes(data int64) []byte {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.BigEndian, data)
	if err != nil {
		log.Panic(err)
	}

	return buf.Bytes()
}

// calculate hash of tx
func (tx Transaction) HashTx() []byte {

	txCopy := tx
	txCopy.Hash = []byte{}

	//serRaw := txCopy.Serialize()

	// fix deep copy
	//var bt []byte = make([]byte, len(serRaw))

	//copy(bt, serRaw)

	hash := sha256.Sum256(txCopy.Serialize())

	return hash[:]
}

// Sign a tx with private key
func (tx *Transaction) Sign(privKey ecdsa.PrivateKey) {
	if tx.IsGenesisTx() {
		return
	}

	// TODO not complete copy
	txCopy := *tx

	txCopy.Hash = nil
	txCopy.Signature = nil

	dataToSign := fmt.Sprintf("%x\n", txCopy)

	r, s, err := ecdsa.Sign(rand.Reader, &privKey, []byte(dataToSign))
	if err != nil {
		log.Panic(err)
	}

	signature := append(r.Bytes(), s.Bytes()...)
	tx.Signature = signature
}

// verify a signature of tx
func (tx *Transaction) Verify() bool {
	/*
	// It just return true in the test
	if tx.IsGenesisTx() {
		return true
	}

	// TODO not complete copy
	txCopy := *tx

	r := big.Int{}
	s := big.Int{}
	sigLen := len(tx.Signature)
	r.SetBytes(tx.Signature[:(sigLen/2)])
	s.SetBytes(tx.Signature[(sigLen/2):])

	x := big.Int{}
	y := big.Int{}
	keyLen := len(tx.Sender)
	x.SetBytes(tx.Sender[:(keyLen/2)])
	y.SetBytes(tx.Sender[(keyLen/2):])

	txCopy.Hash = nil
	txCopy.Signature = nil

	dataToVerify := fmt.Sprintf("%x\n", txCopy)

	curve := elliptic.P256()

	rawPubKey := ecdsa.PublicKey{Curve: curve, X: &x, Y: &y}
	if ecdsa.Verify(&rawPubKey, []byte(dataToVerify), &r, &s) == false {
		return false
	}*/

	return true
}

// serialize tx
func (tx Transaction) Serialize() []byte {
	var encode bytes.Buffer

	enc := gob.NewEncoder(&encode)
	err := enc.Encode(tx)

	if err != nil {
		log.Panic("tx encode fail:", err)
	}

	return encode.Bytes()
}

// Deserialize tx
func DeserializeTx(data []byte) Transaction {
	var tx Transaction

	decode := gob.NewDecoder(bytes.NewReader(data))

	err := decode.Decode(&tx)
	if err != nil {
		log.Panic("tx decode fail:", err)
	}

	return tx
}


// create a new tx including mining
func NewTx(par []byte, validate [][]byte, sender []byte, info []byte, receiver []byte) *Transaction {
	tx := &Transaction{[]byte{}, par, validate, sender, info, receiver, 0, time.Now().Unix(), nil}


	return tx
}

// create a genesis tx
func NewGenesisTx(info []byte, receiver []byte) *Transaction {

	return NewTx([]byte{}, [][]byte{}, []byte{}, info, receiver)

}

// determine if a tx is a genesis tx
func (tx Transaction) IsGenesisTx() bool {
	return 0 == len(tx.Parent) && 0 == len(tx.Validate) && 0 == len(tx.Sender)
}


