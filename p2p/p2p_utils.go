package p2p

import (
	"Jight/utils"
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-crypto"
	"github.com/libp2p/go-libp2p-host"
	"github.com/libp2p/go-libp2p-peerstore"
	"github.com/libp2p/go-libp2p-protocol"
	"github.com/multiformats/go-multiaddr"
	"log"
	"os"

	"Jight/transaction"

	pnet "github.com/libp2p/go-libp2p-net"
)

// maintain a list of all the P2P connection to write data to
var writerList []*bufio.Writer

// a function to deal with the transaction received from the P2P connection
type TxDealerHandler func (tx *transaction.Transaction, local bool) error

// create a LibP2P host with a private key listening on the given listenPort
// listenPort: port to be listened on
// privateKey: RSA private key
// fullAddrsPath: path of a file storing all the full addresses for p2p connection from other nodes
func MakeBasicHost(listenPort int, privateKey crypto.PrivKey, fullAddrsPath string) (host.Host, error) {
	// create the basicHost
	basicHost, err := libp2p.New(
		context.Background(),
		libp2p.ListenAddrStrings(fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", listenPort)),
		libp2p.Identity(privateKey),
	)
	if err != nil {
		return nil, err
	}

	// generate full addresses
	hostAddr, err := multiaddr.NewMultiaddr(fmt.Sprintf("/ipfs/%s", basicHost.ID().Pretty()))
	if err != nil {
		return nil, err
	}
	var fullAddrs []string
	var fullAddr string
	for i, addr := range basicHost.Addrs() {
		fullAddr = addr.Encapsulate(hostAddr).String()
		log.Printf("The %d th fullAddr: %s\n", i, fullAddr)
		fullAddrs = append(fullAddrs, fullAddr)
	}

	//store the full addresses to the file
	if utils.CheckFileExisted(fullAddrsPath) {
		//log.Fatalf("File %s already exists!", fullAddrsPath)
	} else {
		f, err := os.Create(fullAddrsPath)
		if err != nil {
			return nil, err
		}
		defer f.Close()
		for _, addr := range fullAddrs {
			f.WriteString(addr+"\n")
		}
	}
	return basicHost, nil
}


// open a listening for connection
// txDealer: a function to deal with the transaction received from the P2P connection
func OpenP2PListen(pid string, host host.Host, txDealer TxDealerHandler ) {
	h := func (s pnet.Stream) {
		log.Println("Receive a connection")
		w := bufio.NewWriter(s)
		writerList = append(writerList, w)
		go readTxFromP2P(bufio.NewReader(s), txDealer)
	}
	host.SetStreamHandler(protocol.ID(pid), h)
}

// read transaction from P2P connection and deal with the transaction with the handler txDealer
func readTxFromP2P(r *bufio.Reader, txDealer TxDealerHandler) {
	var tx transaction.Transaction
	for {
		str, err := r.ReadString('\n')
		if err!=nil {
			log.Fatal(err)
		}
		if err := json.Unmarshal([]byte(str), &tx); err!=nil {
			log.Fatal(err)
		}
		if err := txDealer(&tx, false); err!=nil {
			log.Fatal(err)
		}
	}
}

// connect to a specific p2p node
func ConnectP2PNode(dest *string, pid string, host host.Host, txDealer TxDealerHandler) error {
	maddr, err := multiaddr.NewMultiaddr(*dest)
	if err != nil {
		return err
	}
	info, err := peerstore.InfoFromP2pAddr(maddr)
	if err != nil {
		return err
	}
	host.Peerstore().AddAddrs(info.ID, info.Addrs, peerstore.PermanentAddrTTL)

	s, err := host.NewStream(
		context.Background(),
		info.ID,
		protocol.ID(pid),
	)
	if err != nil {
		return err
	}

	// add the P2P stream to the static variable
	w := bufio.NewWriter(s)
	writerList = append(writerList, w)
	go readTxFromP2P(bufio.NewReader(s), txDealer)

	return nil
}

// sync a transaction to all the p2p connection
func SyncTx(tx *transaction.Transaction) error {
	bytes, err := json.Marshal(*tx)
	if err!=nil {
		return err
	}
	for k, w := range writerList{
		w.WriteString(fmt.Sprintf("%s\n", string(bytes)))
		w.Flush()
		log.Println(fmt.Sprintf("Finish sending the tx to %d the node", k))
	}
	return nil
}