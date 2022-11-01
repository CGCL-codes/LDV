package api

import (
	"github.com/urfave/cli"
	"os"
)

type CMDDeamon struct {
	HelpFlag bool
	P2PPort int
	RPCPort int
	TargetPath string
	Pid string
	FullAddrsPath string
	CreatePK bool
	PrivKeyPath string
	GenesisNode bool
	LikeTopics cli.StringSlice
	GenesisTxNum uint
	IsPrune bool
	PruneRange int
	KeepHash bool
}

func (cmd *CMDDeamon) Run() error {
	app := cli.NewApp()
	app.Name = "jightd"
	app.Flags = []cli.Flag {
		cli.IntFlag{
			Name: "p2pport, l",
			Usage: "P2P Port to listen to",
			Value: 8525,
			Destination: &cmd.P2PPort,
		},
		cli.IntFlag{
			Name: "rpcport, r",
			Usage: "RPC port to listen to",
			Value: 9525,
			Destination: &cmd.RPCPort,
		},
		cli.StringFlag{
			Name: "targetPath, t",
			Usage: "Path of a file storing the destination p2p nodes to connect to",
			Value: "targetaddrs.txt",
			Destination: &cmd.TargetPath,
		},
		cli.StringFlag{
			Name: "pid, p",
			Value: "/jightd/1.0.0",
			Usage: "pid to identify a p2p network",
			Destination: &cmd.Pid,
		},
		cli.StringFlag{
			Name: "fulladdrspath, f",
			Value: "fulladdrs.txt",
			Usage: "Path of a file to store the full p2p addresses for listening",
			Destination: &cmd.FullAddrsPath,
		},
		cli.BoolFlag{
			Name: "createpk, c",
			Hidden: false,
			Usage: "Whether to create a new private key",
			Destination: &cmd.CreatePK,
		},
		cli.StringFlag{
			Name: "privkeypath, k",
			Value: "jight.pk",
			Usage: "Path of a file to store the private key",
			Destination: &cmd.PrivKeyPath,
		},
		cli.BoolFlag{
			Name: "genesisnode, g",
			Hidden: false,
			Usage: "Whether is the node a genesis node",
			Destination: &cmd.GenesisNode,
		},
		cli.StringSliceFlag{
			Name: "follow, fo",
			//Value: &cli.StringSlice{"topic1", "topic2"},
			Usage: "Topics of interest",
		},
		cli.UintFlag{
			Name: "amount, a",
			Value: 100,
			Usage: "Amounts of genesis tx",
			Destination: &cmd.GenesisTxNum,
		},
		cli.IntFlag{
			Name: "prunerange, pr",
			Value: 9000,
			Usage: "Range of pruning transactions",
			Destination: &cmd.PruneRange,
		},
		cli.BoolFlag{
			Name: "prune, u",
			Hidden: false,
			Usage: "Whether to prune history",
			Destination: &cmd.IsPrune,
		},
		cli.BoolFlag{
			Name: "keephash, kh",
			Hidden: false,
			Usage: "Whether to keep hash of transaction",
			Destination: &cmd.KeepHash,
		},
	}

	app.Action = func(c *cli.Context) error {
		cmd.LikeTopics = c.StringSlice("follow")
		return nil
	}

	cli.HelpFlag = cli.BoolFlag{
		Name: "help, h",
		Destination: &cmd.HelpFlag,
	}

	err := app.Run(os.Args)

	if err!=nil {
		return err
	}
	return nil
}
