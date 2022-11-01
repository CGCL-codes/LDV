package main

import (
    "Jight/config"
    "Jight/transaction"
    "Jight/utils"
    "Jight/wallet"
    "errors"
    "fmt"
    "log"
    "math/rand"
    "time"
)

func init() {
    rand.Seed(time.Now().UnixNano())
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringRunes(n int) string {
    b := make([]rune, n)
    for i := range b {
        b[i] = letterRunes[rand.Intn(len(letterRunes))]
    }
    return string(b)
}

func main() {
    //for i:=3; i<16; i++ {
    //    fmt.Println("---------------------------")
    //    var n = int(math.Pow(2, float64(i)))
    //    fmt.Println("len of info: ", n)
    //    tx := createTx(RandStringRunes(n))
    //    fmt.Println(tx)
    //    txRaw := tx.Serialize()
    //    originLen := len(txRaw)
    //    fmt.Println(originLen)
    //
    //    txHash := tx.ToTxHash()
    //    txHashRaw := txHash.Serialize()
    //    txHashLen := len(txHashRaw)
    //    fmt.Println(txHashLen)
    //
    //}
    fmt.Println(RandStringRunes(1000))
}

func createTx(info string) *transaction.Transaction {
    txTmp := transaction.NewTx(
        []byte("444444333asdf"),
        [][]byte{[]byte("sdfkjlwke"), []byte("sdfkjlwke")},
        []byte("12E9tH5tL2td8jgkrK8WCtkEJpCRgHDXpX"),
        []byte("This road is fine"),
        []byte("luoyu"),
    )

    ref := txTmp.HashTx()
    var validateRef [][]byte
    validateRef = append(validateRef, ref, ref, ref, ref, ref, ref)

    nodeID := "self"

    addr, _ := initializeWallets(nodeID)

    ws, err := wallet.LoadWallets(nodeID)
    if err != nil {
        log.Panic(err)
    }

    wallet := ws.GetWallet(addr)

    tx := transaction.NewTx(ref, validateRef, []byte(addr), []byte(info), []byte("topic1"))

    tx.Sign(wallet.PrivateKey)
    tx.Nonce = tx.Pow()
    tx.Hash = tx.HashTx()

    return tx
}

func initializeWallets(nodeID string) (string, error) {
    if utils.CheckFileExisted(config.WALLET_FILE) {
        return "", errors.New("wallet file is already existed")
    }
    //for _, addr := range(config.GenesisAddresses) {
    addr := createWallet(nodeID)
    //}
    return addr, nil
}


func createWallet(nodeID string) string {
    fmt.Println("create wallet")

    // write the new address to local wallet_self.dat
    wallets, _ := wallet.LoadWallets(nodeID)
    address := wallets.CreateWallet()
    wallets.SaveToFile(nodeID)
    //address := "test"

    fmt.Printf("Generated address: %s\n", address)
    return address
}

