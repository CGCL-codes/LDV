package config

import (
	"container/list"
	"fmt"
	"math"
)

// TODO how to set a pow target which need few computational
const POW_TARGET_BITS int = 24
const MAX_NONCE uint32 = math.MaxUint32


// TODO genesis tx value
//const GENESIS_VALUE uint32 = 10
const GENESIS_VALUE string = "genesis transaction"
var GENESIS_ADDR_COUNT uint = 10

// Default data storage location
const LOC string = "E:\\Jight\\"

// TODO db location storing dagchain
const DB_LOCATION string = LOC + "db"

var Interested_Topics []string

var SlideWindow *list.List

var Is_Prune = false
var Prune_Range = 9000

var Keep_Hash = false

// Wallet const
const (
	VERSION = byte(0x00)
	ADDRESS_CHECK_SUM_LEN = 4
	WALLET_FILE = "wallet_%s.dat"
)

//var GenesisAddresses = make([]string, GENESIS_ADDR_COUNT)
var GenesisAddresses []string

// convert a integer (< 1000) to a 4 byte string
func convertIntTo4Byte(i int) string {
	if i < 10 {
		return fmt.Sprintf("000%d", i)
	} else if i < 100 {
		return fmt.Sprintf("00%d", i)
	} else if i < 1000 {
		return fmt.Sprintf("0%d", i)
	} else {
		return fmt.Sprintf("%d", i)
	}
}

/*func init() {
	for i:= 0; i< int(GENESIS_ADDR_COUNT); i++ {
		GenesisAddresses[i] = "158ewJ1itTAAE8gy1Hk45JdAmqvbpd" + convertIntTo4Byte(i)
	}
}*/
