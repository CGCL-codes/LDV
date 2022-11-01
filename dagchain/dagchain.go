package dagchain

import (
	"Jight/config"
	"Jight/transaction"
	"Jight/wallet"
	"bytes"
	"crypto/ecdsa"
	"encoding/gob"
	"errors"
	"fmt"
	"github.com/syndtr/goleveldb/leveldb"
	"log"
	"math/rand"
	"time"
)

/*
 Different nodes may see different views of tips due to the network latency
 However, since we cannot deploy the nodes all across the world currently, we have to simulate the different views of tips
 Specifically,
	1) if a node issues a tx citing two tips by itself, it will delete the tips directly;
	2) if a node receives a tx from other nodes which cites two tips, only the 'cited' flags in these two tips are set
	3) tips with 'cited' flags being set will be deleted later, either periodically or triggered by the cli
 */
type Tip struct {
	TipID []byte	// tip txid
	cited bool	// denote if it has been cited by the txs from other nodes
	//Account []byte	// Account address of tip
}

type DagChain struct {
	Tips []*Tip // tips of dagchain, persistence
	//Account map[string][]byte		// no persistence, reindex at CreateDagChain or RecoverDagChain, to reduce storage costs and speed up FindParent
	DB *leveldb.DB
}

/*
	--------------
	  tx related
	--------------
*/
// receive tx from p2p network, then store it to db
// local; denote if the tx is created locally or received from other nodes
func (dc *DagChain) AddTx(tx *transaction.Transaction, local bool) error {
	// Tips view has no effect in data reduction, TODO remove local

	fmt.Printf("Tips count: %d\n", len(dc.Tips))
	// verify tx from network
	if tx.Verify() == false {
		return errors.New("Tx verify failed, invalid Tx")
	}

	local = true

	if local {



		fmt.Println("Receive tx locally: ", tx)

		// remove old tips and add new tip
		var index []int
		oldTips := tx.Validate
		for i, t := range dc.Tips {
			if bytes.Equal(t.TipID, oldTips[0]) || bytes.Equal(t.TipID, oldTips[1]) {
				index = append(index, i)
			}
		}
		newTip := &Tip{tx.Hash, false}
		log.Println("length of oldTips: ", len(index))
		if len(index) == 2 {
			dc.Tips[index[0]] = newTip
			dc.Tips = append(dc.Tips[:index[1]], dc.Tips[index[1]+1:]...)
		} else if len(index) == 1 {
			dc.Tips[index[0]] = newTip
		} else {
			dc.Tips = append(dc.Tips, newTip)
		}
		// write to DB
		if ContainTopic(tx.Topic) {
			err := StoreTx(dc.DB, *tx)
			if err != nil {
				log.Panic("store tx failed: ", err)
			}

			if config.Is_Prune {
				if config.SlideWindow.Len() < config.Prune_Range {
					config.SlideWindow.PushBack(tx.Hash)
				} else {
					first := config.SlideWindow.Front()
					if first == nil {
						return errors.New("the first element of silde window is nil")
					}

					rmHashInf := config.SlideWindow.Remove(first)

					rmHash := rmHashInf.([]byte)

					if config.Keep_Hash {
						// get tx from db
						rmTx, errRmTx := FetchTx(dc.DB, rmHash)
						if errRmTx != nil {
							return errors.New("get the removed tx fails")
						}

						rmTxHash := rmTx.ToTxHash()
						errRmTxHash := StoreTxHash(dc.DB, *rmTxHash)

						if errRmTxHash != nil {
							log.Panic("store rmTxHash failed: ", errRmTxHash)
						}
					}

					errRm := RemoveTx(dc.DB, rmHash)

					if errRm != nil {
						return errors.New("rm tx fails")
					}
					config.SlideWindow.PushBack(tx.Hash)
				}
			}

		} else if config.Keep_Hash {
			// store hash
			txHash := tx.ToTxHash()
			err := StoreTxHash(dc.DB, *txHash)

			if err != nil {
				log.Panic("store txHash failed: ", err)
			}
		}

		for _, v := range oldTips {
			RemoveTip(dc.DB, v)
		}
		PutTip(dc.DB, *newTip)
		//dc.Tips = append(dc.Tips, newTip)
	} else {
		fmt.Println("Receive tx remotely: ", tx)
		// set the cited flag and add new tip
		//oldTips := tx.Validate
		var citedChangeTips []*Tip
		/*for _, t := range dc.Tips {
			if bytes.Equal(t.TipID, oldTips[0]) || bytes.Equal(t.TipID, oldTips[1]) {
				t.cited = true
				citedChangeTips = append(citedChangeTips, t)
			}
		}*/
		newTip := &Tip{tx.Hash, false}
		dc.Tips = append(dc.Tips, newTip)

		// write to DB
		err := StoreTx(dc.DB, *tx)
		if err != nil {
			log.Panic("store tx failed: ", err)
		}
		for _, v := range citedChangeTips {
			PutTip(dc.DB, *v)
		}
		PutTip(dc.DB, *newTip)
		dc.Tips = append(dc.Tips, newTip)
	}

	fmt.Printf("Tips count: %d\n", len(dc.Tips))

	return nil
}
/*
func GetBytes(v interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(v)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}*/

func ContainTopic(topic []byte) bool {
	topicStr := string(topic[:])

	for _, t := range config.Interested_Topics {
		if topicStr == t {
			return true
		}
	}

	return false
}


// create (mining) a tx
func (dc *DagChain) CreateTx(privKey ecdsa.PrivateKey, sender, receiver []byte, info []byte) *transaction.Transaction {
	var validateRef [][]byte

	//index, vTips := dc.SelectTips()
	_, vTips := dc.SelectTips()
	for dc.VerifyTips(vTips) == false {
		_, vTips = dc.SelectTips()
	}
	log.Println("length of selectedTips: ", len(vTips))
	lastTx := dc.FindParent(sender)

	for _, v := range vTips {
		validateRef = append(validateRef, v.TipID)
	}

	tx := transaction.NewTx(lastTx, validateRef, sender, info, receiver)

	// sign a tx
	tx.Sign(privKey)
	tx.Nonce = tx.Pow()
	tx.Hash = tx.HashTx()

	return tx
}

// Signature a tx with privKey
func (dc *DagChain) SignTx(tx *transaction.Transaction, privKey ecdsa.PrivateKey)  {
	if tx.IsGenesisTx() {
		return
	}

	tx.Sign(privKey)
}

// Verify if a tx is valid
func (dc *DagChain) VerifyTx(tx *transaction.Transaction) bool {
	if tx.IsGenesisTx() {
		return true
	}

	return tx.Verify()
}

// Init a dagchain used in daemon
func Init() (*DagChain, error) {
	// db will be closed in the main() function
	db, err := LoadDB()

	if err != nil {
		return nil, err
	}
	tips, err := FetchTips(db)
	if err != nil {
		return nil, err
	}
	//address := []byte(config.GENESIS_TO_ADDRESS)

	// if tips is nil, initialize the new dagchain
	if tips == nil {
		tips, err = InitializeGenesisTxs(db)
		if err!=nil {
			return nil, err
		}
	}

	return &DagChain{tips, db}, nil
}

/*
	--------------------
	  dagchain related
	--------------------
*/

// Generate account for testing
func InitAccount(db *leveldb.DB) error {
	nodeID := "self"

	for i:= 0; i< int(config.GENESIS_ADDR_COUNT); i++ {
		wallets, _ := wallet.LoadWallets(nodeID)
		address := wallets.CreateWallet()
		wallets.SaveToFile(nodeID)

		// lastid of new account is nil
		//acc := Account{[]byte(address), nil, 0}
		log.Println("func InitAccount is called, address:", address)
		//err := StoreAccount(db, acc)

		config.GenesisAddresses[i] = address

		/*if err != nil {
			return err
		}*/
	}

	return nil
}

// Initialize the 1000 genesis transactions and 1000 addresses, each address has 10000 coins
// @parameter address: receiver of genesis tx
func InitializeGenesisTxs(db *leveldb.DB) ([]*Tip, error) {

	_ = InitAccount(db)

	var tips []*Tip
	for i, addr:= range(config.GenesisAddresses) {
		genesisTx := transaction.NewGenesisTx([]byte(config.GENESIS_VALUE), []byte("genesis"+string(i)))
		genesisTx.Nonce = genesisTx.Pow()
		genesisTx.Hash = genesisTx.HashTx()
		genesisID := genesisTx.Hash
		acc := Account{[]byte(addr), genesisID, [][]byte{[]byte("account topic1")}}
		tips = append(tips, &Tip{genesisID, false})

		log.Println("genesis tx is creating, hash: ", genesisID)

		// write to db
		if err := StoreTx(db, *genesisTx); err != nil {
			return nil, err
		}
		if err := StoreAccount(db, acc); err != nil {
			return nil, err
		}

	}
	if err := StoreTips(db, tips); err != nil {
		return nil, err
	}
	return tips, nil
}


/*
	--------------
     tips related
	--------------
*/

//serialize tips
func (tip Tip) Serialize() []byte {
	var encode bytes.Buffer

	enc := gob.NewEncoder(&encode)
	err := enc.Encode(tip)

	if err != nil {
		log.Panic("tip encode fail:", err)
	}

	return encode.Bytes()
}

//deserialize tips
func DeserializeTip(data []byte) Tip {
	var tip Tip

	decode := gob.NewDecoder(bytes.NewReader(data))

	err := decode.Decode(&tip)
	if err != nil {
		log.Panic("tip decode fail:", err)
	}

	return tip
}

// @return tips of DagChain TODO
func (dc *DagChain) FindTips() []*Tip {

	return dc.Tips
}

// @return last tx's hash of special Account address
func (dc *DagChain) FindParent(account []byte) []byte {
	// Todo read db file of Account to find an Account

	acc, err := FetchAccount(dc.DB, account)
	if err != nil {
		log.Panic(err)
	}

	return acc.LastId

	/*lastTx, errTip := database.FetchTx(dc.DB, acc.LastId)
	if errTip != nil {
		log.Panic(errTip)
	}

	return lastTx*/

	/*for _, tip := range tips {
		tx := transaction.FindTxByHash(tip.TipID)
		if bytes.Equal(tx.Sender, Account) {
			return tip, tx
		}
	}*/
}

// select two tips to verify, randomly
func (dc *DagChain) SelectTips() ([]int, []*Tip) {
	allTips := dc.FindTips()

	var twoTips [2]*Tip
	var index [2]int

	rand.Seed(time.Now().UnixNano())
	first := rand.Intn(len(allTips))
	second := rand.Intn(len(allTips))
	/*for first == second {
		log.Println("first == second")
		second = rand.Intn(len(allTips))
	}*/

	index[0] = first
	index[1] = second

	twoTips[0] = allTips[first]
	twoTips[1] = allTips[second]


	/*twoTips := [2]*Tip {
		&Tip{[]byte("tseafooler111111111"), false},
		&Tip{[]byte("tseafooler222222222"), false},
	}
	index := [2]int {2, 3}*/

	return index[:], twoTips[:]
}

// verify selected tips tx
func (dc *DagChain) VerifyTips(tips []*Tip) bool {
	// verify selected tx
	/*for _, t := range tips {
		tipTx, err := FetchTx(dc.DB, t.TipID)
		if err != nil {
			log.Panic(err)
		}
		if tipTx.Verify() == false {
			return false
		}
	}*/

	return true
}
