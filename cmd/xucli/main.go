package main

import (
	"log"
	"time"
	"github.com/urfave/cli"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	pb "github.com/offerm/swap-resolver/swapresolver"
	"os"
	"fmt"
)

const (
	defaultRpcPort          = "10000"
	defaultRpcHostPort      = "localhost:" + defaultRpcPort
)

var (
	//Commit stores the current commit hash of this build. This should be
	//set using -ldflags during compilation.
	Commit string

)

func main() {
	app := cli.NewApp()
	app.Name = "xucli"
	app.Version = fmt.Sprintf("%s commit=%s", "0.0.1", Commit)
	app.Usage = "Use me to simulate order taking by the taker"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "rpcserver",
			Value: defaultRpcHostPort,
			Usage: "host:port of xu daemon",
		},
	}
	app.Commands = []cli.Command{
		takeordercommand,
	}

	if err := app.Run(os.Args); err != nil {
		fatal(err)
	}
}

func fatal(err error) {
	fmt.Fprintf(os.Stderr, "[lncli] %v\n", err)
	os.Exit(1)
}

func getClient(ctx *cli.Context) (pb.P2PClient, func()) {
	conn := getClientConn(ctx, false)

	cleanUp := func() {
		conn.Close()
	}

	return pb.NewP2PClient(conn), cleanUp
}

func getClientConn(ctx *cli.Context, skipMacaroons bool) *grpc.ClientConn {
	var opts []grpc.DialOption
	//if *tls {
	//if *caFile == "" {
	//	*caFile = testdata.Path("ca.pem")
	//}
	//creds, err := credentials.NewClientTLSFromFile(*caFile, *serverHostOverride)
	//if err != nil {
	//	log.Fatalf("Failed to create TLS credentials %v", err)
	//}
	//opts = append(opts, grpc.WithTransportCredentials(creds))
	//} else {
	opts = append(opts, grpc.WithInsecure())
	//}
	conn, err := grpc.Dial(ctx.GlobalString("rpcserver"), opts...)

	if err != nil {
		fatal(err)
	}

	return conn
}


var takeordercommand = cli.Command{
	Name:     "takeorder",
	Category: "Order",
	Usage:    "Instruct XUD to take an order.",
	Description: `
	`,
	//ArgsUsage: "order_id maker_amount maker_coin taker_amount taker_coin",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "order_id",
			Usage: "the chain to generate an address for",
		},
		cli.Int64Flag{
			Name:  "maker_amount",
			Usage: "the number of coins denominated in {l/s}atoshis the maker should send",
		},
		cli.StringFlag{
			Name:  "maker_coin",
			Usage: "the coins which the maker should send to the taker",
		},
		cli.Int64Flag{
			Name:  "taker_amount",
			Usage: "the number of coins denominated in {l/s}atoshis the taker should send",
		},
		cli.StringFlag{
			Name:  "taker_coin",
			Usage: "the coins which the taker should send to the maker",
		},

	},

	Action: takeOrder,
}

func takeOrder(ctx *cli.Context) error{

	client, cleanUp := getClient(ctx)
	defer cleanUp()

	//args := ctx.Args()
	//var err error

	// Show command help if no arguments provided
	if ctx.NArg() == 0 && ctx.NumFlags() == 0 {
		cli.ShowCommandHelp(ctx, "takeOrder")
		return nil
	}

	req := &pb.TakeOrderReq{
	}

	if !ctx.IsSet("order_id"){
		return fmt.Errorf("order_id argument missing")
	}
	req.Orderid = ctx.String("order_id")

	if !ctx.IsSet("maker_amount"){
		return fmt.Errorf("maker_amount argument missing")
	}
	req.MakerAmount = int64(ctx.Int("maker_amount"))

	if !ctx.IsSet("maker_coin"){
		return fmt.Errorf("maker_coin argument missing")
	}
	req.MakerCoin = ctx.String("maker_coin")

	if !ctx.IsSet("taker_amount"){
		return fmt.Errorf("taker_amount argument missing")
	}
	req.TakerAmount = int64(ctx.Int("taker_amount"))

	if !ctx.IsSet("taker_coin"){
		return fmt.Errorf("taker_coin argument missing")
	}
	req.TakerCoin = ctx.String("taker_coin")


	log.Printf("Starting takeOrder command: %v",req)
	ctxt, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	preimage, err := client.TakeOrder(ctxt, req)
	if err != nil {
		log.Fatalf("%v.ResolveHash(_) = _, %v: ", client, err)
	}
	log.Printf("Got response from Resolver: %v \n",preimage)


	return nil;

}



