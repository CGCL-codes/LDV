package api

import (
	"Jight/config"
	jrpc "Jight/rpc"
	"Jight/utils"
	"Jight/wallet"
	"errors"
	"fmt"
	"github.com/urfave/cli"
	"log"
	"net/rpc"
	"os"
	"strconv"
)

type CMD struct {
	rpcPort int
}


func (cmd *CMD) Run() {
	// TODO nodeID should be mac or ip address, nodeID specialize a node
	nodeID := "self"

	app := cli.NewApp()

	app.Flags = []cli.Flag {
		cli.IntFlag{
			Name: "rpcport, r",
			Value: 9525,
			Usage: "Daemon RPC port to connect to",
			Destination: &cmd.rpcPort,
		},
	}

	app.Commands = []cli.Command{
		// createdagchain CMD
		{
			Name:        "createdagchain",
			Aliases:     []string{"cdc"},
			Usage:       "Create a dagchain and send genesis tx reward to ADDRESS",
			UsageText:   "createdagchain -address ADDRESS",
			Description: "Create a dagchain",
			Flags: []cli.Flag{
				cli.StringFlag{Name: "address, a"},
			},
			Action: func(c *cli.Context) error {
				address := c.String("address")
				cmd.createDagChain(address)
				return nil
			},
		},


		// getbalance CMD
		{
			Name:        "getAccountTopics",
			Aliases:     []string{"gat"},
			Usage:       "Get topics of ADDRESS",
			UsageText:   "getAccountTopics -address ADDRESS",
			Description: "Get topics of ADDRESS",
			Flags: []cli.Flag{
				cli.StringFlag{Name: "address, a"},
			},
			Action: func(c *cli.Context) error {
				address := c.String("address")
				cmd.getTopics(address)
				return nil
			},
		},

		// get tx CMD
		{
			Name:        "getTransaction",
			Aliases:     []string{"gtx"},
			Usage:       "Get transaction by tx hash",
			UsageText:   "getTransaction -id ADDRESS",
			Description: "Get transaction by tx hash",
			Flags: []cli.Flag{
				cli.StringFlag{Name: "id, i"},
			},
			Action: func(c *cli.Context) error {
				hash := c.String("id")
				cmd.getTransaction(hash)
				return nil
			},
		},

		// list all address of wallet CMD
		{
			Name:        "listaddress",
			Aliases:     []string{"la"},
			Usage:       "Lists all addresses from the wallet file",
			UsageText:   "listaddress",
			Description: "Lists all addresses from the wallet file",
			Action: func(c *cli.Context) error {
				cmd.listAddress(nodeID)
				return nil
			},
		},

		// compact tx in log to db files
		{
			Name:        "compact",
			Aliases:     []string{"cmp"},
			Usage:       "Compact tx from log file to db files",
			UsageText:   "compact",
			Description: "Compact tx from log file to db files",
			Action: func(c *cli.Context) error {
				cmd.compact()
				return nil
			},
		},

		// send money from f to t CMD
		{
			Name:        "send",
			Aliases:     []string{"sd"},
			Usage:       "Send information from FROM address to topic TO",
			UsageText:   "send -from FROM -to TO -info INFO",
			Description: "Send information from FROM address to topic TO",
			Flags: []cli.Flag{
				cli.StringFlag{Name: "from, f"},
				cli.StringFlag{Name: "to, t"},
				cli.StringFlag{Name: "info, i"},
			},
			Action: func(c *cli.Context) error {
				from := c.String("from")
				to := c.String("to")
				info := c.String("info")
				cmd.send(from, info, to)
				return nil
			},
		},

		{
			Name: "initializewallets",
			Aliases: []string{"iw"},
			Usage: "Initialize the wallets for test according to the 1000 genesis addresses",
			Action: func(c *cli.Context) error {
				if err := cmd.initializeWallets(nodeID); err != nil {
					return  err
				}
				return nil
			},
		},

		{
			Name: "createwallets",
			Aliases: []string{"cw"},
			Usage: "Create the wallets",
			Action: func(c *cli.Context) error {
				cmd.createWallet(nodeID)
				return nil
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Panic(err)
	}

}

// TODO this func should implemented in daemon
func (cmd *CMD) createDagChain(address string) {
	fmt.Println("create DAGCHAIN: ", address)
	client, err := rpc.DialHTTP("tcp", "localhost:"+ strconv.Itoa(cmd.rpcPort))
	if err != nil {
		log.Fatal("Connection error: ", err)
	}
	cdcc := jrpc.CreateDagChainCMD{
		Address: "tJight123344566765765756",
	}
	var cdcr jrpc.CreateDagChainReply

	err = client.Call("Jightd.CreateDagChain", cdcc, &cdcr)
	if err != nil {
		log.Fatal("CreateDagChain error: ", err)
	} else {
		fmt.Println("CreateDagChain CMD: ", cdcr.CdcrField)
	}

}

func (cmd *CMD) initializeWallets(nodeID string) error {
	if utils.CheckFileExisted(config.WALLET_FILE) {
		return errors.New("wallet file is already existed")
	}
	//for _, addr := range(config.GenesisAddresses) {
		cmd.createWallet(nodeID)
	//}
	return nil
}


func (cmd *CMD) createWallet(nodeID string) {
	fmt.Println("create wallet")

	// write the new address to local wallet_self.dat
	wallets, _ := wallet.LoadWallets(nodeID)
	address := wallets.CreateWallet()
	wallets.SaveToFile(nodeID)
	//address := "test"

	// write the new address to DB via rpc
	log.Println("cmd.rpcPort: ", cmd.rpcPort)
	client, err := rpc.DialHTTP("tcp", "localhost:"+ strconv.Itoa(cmd.rpcPort))
	if err != nil {
		log.Fatal("Connection error: ", err)
	}
	cwc := jrpc.CreateWalletCMD{
		Address: address,
	}

	var cwr jrpc.CreateWalletReply

	err = client.Call("Jightd.CreateWallet", cwc, &cwr)
	if err != nil {
		log.Fatal("CreateWallet error: ", err)
	}

	fmt.Printf("Generated address: %s\n", address)

}

func (cmd *CMD) getTopics(address string) {
	fmt.Println("get topic of ", address)
	client, err := rpc.DialHTTP("tcp", "localhost:"+ strconv.Itoa(cmd.rpcPort))
	if err != nil {
		log.Fatal("Connection error: ", err)
	}

	gbc := jrpc.GetTopicsCMD{
		Address: address,
	}

	var gbr jrpc.GetAccTopicsReply

	err = client.Call("Jightd.GetTopics", gbc, &gbr)
	if err != nil {
		log.Fatal("GetTopics error: ", err)
	} else {
		fmt.Printf("Topics of %s: %s\n", address, gbr.Topics)
	}
}

/*
// Merge conflict
func (cmd *CMD) listAddress() {
	fmt.Println("list address")
	client, err := rpc.DialHTTP("tcp", "localhost:"+ strconv.Itoa(cmd.rpcPort))
	if err != nil {
		log.Fatal("Connection error: ", err)
	}
	lac := jrpc.ListAddressCMD{
	}

	var lar jrpc.ListAddressReply

	err = client.Call("Jightd.ListAddress", lac, &lar)
	if err != nil {
		log.Fatal("ListAddress error: ", err)
	} else {
		fmt.Println("ListAddress CMD: ", lar.Addresses)
*/

func (cmd *CMD) listAddress(nodeID string) {
	wallets, err := wallet.LoadWallets(nodeID)
	if err != nil {
		log.Panic(err)
	}
	addresses := wallets.GetAddresses()

	for i, address := range addresses {
		if i >= 5 {
			break
		}
		fmt.Println(address)
	}
}

func (cmd *CMD) send(from string, info string, to string) {
	fmt.Println("send from ", from, "to ", to, "info ", info)
	client, err := rpc.DialHTTP("tcp", "localhost:"+ strconv.Itoa(cmd.rpcPort))
	if err != nil {
		log.Fatal("Connection error: ", err)
	}

	sc := jrpc.SendCmd{
		From: from,
		To: to,
		Info: info,
	}

	var sr jrpc.SendReply

	err = client.Call("Jightd.Send", sc, &sr)
	if err != nil {
		log.Fatal("Send error: ", err)
	} else {
		fmt.Println("Send CMD: ", sr.SrField)
	}
}

func (cmd *CMD) getTransaction(hash string) {
	fmt.Println("get transaction of ", hash)
	client, err := rpc.DialHTTP("tcp", "localhost:"+ strconv.Itoa(cmd.rpcPort))
	if err != nil {
		log.Fatal("Connection error: ", err)
	}

	gbc := jrpc.GetTransactionCMD{
		Hash: hash,
	}

	var gbr jrpc.GetTransactionReply

	err = client.Call("Jightd.GetTransaction", gbc, &gbr)
	if err != nil {
		log.Fatal("GetTransaction error: ", err)
	} else {
		fmt.Printf("Transaction Topic: %s, Info: %s\n", gbr.Tx.Topic, string(gbr.Tx.Info))
	}
}

func (cmd *CMD) compact() {
	fmt.Println("compact tx")
	client, err := rpc.DialHTTP("tcp", "localhost:"+ strconv.Itoa(cmd.rpcPort))
	if err != nil {
		log.Fatal("Connection error: ", err)
	}

	gbc := jrpc.CompactTxCMD{

	}

	var gbr jrpc.CompactTxReply

	err = client.Call("Jightd.CompactTx", gbc, &gbr)
	if err != nil {
		log.Fatal("Compact tx error: ", err)
	} else {
		fmt.Printf("Compact tx success, Info: %s\n", gbr.Info)
	}
}
