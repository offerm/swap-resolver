
package main
// TODO: seperate resolver RPC into a standalone package
// TODO: properly handle the time locks - Important!!!!

import (
	"fmt"
	//"io"
	//"io/ioutil"
	"log"
	//"math"
	"net"
	"sync"
	//"time"
	"crypto/rand"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	//"google.golang.org/grpc/credentials"
	//"google.golang.org/grpc/testdata"

	//"github.com/golang/protobuf/proto"

	pb "github.com/offerm/swap-resolver/swapresolver"
	"github.com/urfave/cli"
	"os"
	"github.com/davecgh/go-spew/spew"
	"time"
	"github.com/dchest/uniuri"
	//"github.com/lightningnetwork/lnd/macaroons"
	//"github.com/lightningnetwork/lnd/lncfg"
	"github.com/lightningnetwork/lnd/lnrpc"
	"github.com/btcsuite/btcutil"
	"path/filepath"
	"google.golang.org/grpc/credentials"
	//"io/ioutil"
	//"gopkg.in/macaroon.v2"
	//"strings"
	//"os/user"
	//"github.com/lightningnetwork/lnd/macaroons"
	//"github.com/lightningnetwork/lnd/lncfg"
	"crypto/sha256"
)

type swapResolverServer struct {
	p2pServer *P2PServer
	mu         sync.Mutex // protects data structure
}

type role int

const (
	defaultTLSCertFilename  = "tls.cert"
	defaultMacaroonFilename = "admin.macaroon"
	defaultRpcPort          = "10009"
	defaultRpcHostPort      = "localhost:" + defaultRpcPort


	Maker role = iota + 1
	Taker
)

type deal struct {
	// Maker or Taker ?
	myRole		role
	// global order it in XU network
	orderId     string
	takerDealId string
	takerAmount int64
	// takerCoin is the name of the coin the taker is expecting to get
	takerCoin   pb.CoinType
	takerPubKey string
	makerDealId string
	makerAmount int64
	// makerCoin is the name of the coin the maker is expecting to get
	makerCoin   pb.CoinType
	makerPubKey string
	hash		[32]byte
	preImage	[32]byte
	createTime	time.Time
	executeTime	time.Time
	competionTime	time.Time
}


var (
	//tls        = flag.Bool("tls", false, "Connection uses TLS if true, else plain TCP")
	//certFile   = flag.String("cert_file", "", "The TLS cert file")
	//keyFile    = flag.String("key_file", "", "The TLS key file")
	//port       = flag.Int("port", 10000, "The server port")

	//Commit stores the current commit hash of this build. This should be
	//set using -ldflags during compilation.
	Commit string

	deals []*deal

	defaultLndDir       = btcutil.AppDataDir("lnd", false)
	defaultTLSCertPath  = filepath.Join(defaultLndDir, defaultTLSCertFilename)
	defaultMacaroonPath = filepath.Join(defaultLndDir, defaultMacaroonFilename)

)


// GetFeature returns the feature at the given point.
func (s *swapResolverServer) ResolveHash(ctx context.Context, req *pb.ResolveReq) (*pb.ResolveResp, error) {

	var deal *deal

	log.Printf("ResolveHash stating with for hash: %v ",req.Hash)

	for _, d := range deals{
		if string(d.hash[:]) == req.Hash{
			deal = d
			break
		}
	}

	if deal == nil{
		log.Printf("Something went wrong. Can't find deal in swapResolverServer: %v ",req.Hash)
		return nil, fmt.Errorf("Something went wrong. Can't find deal in swapResolverServer")
	}

	// If I'm the taker I need to forward the payment to the other chanin
	// TODO: check that I got the right amount before sending out the agreed amount
	if deal.myRole == Taker{

		log.Printf("Taker code")

		cmdLnd := s.p2pServer.lnBTC

		switch deal.makerCoin{
		case pb.CoinType_BTC:

		case pb.CoinType_LTC:
			cmdLnd = s.p2pServer.lnLTC

		}

		lncctx := context.Background()
		resp, err := cmdLnd.SendPaymentSync(lncctx,&lnrpc.SendRequest{
			DestString:deal.makerPubKey,
			Amt:deal.makerAmount,
			PaymentHash:deal.hash[:],
		})
		if err != nil{
			err = fmt.Errorf("Got error sending  %d %v by taker - %v",
				deal.makerAmount,deal.makerCoin.String(),err)
			log.Printf(err.Error())
			return nil, err
		}
		if resp.PaymentError != ""{
			err = fmt.Errorf("Got PaymentError sending %d %v by taker - %v",
				deal.makerAmount,deal.makerCoin.String(), resp.PaymentError)
			log.Printf(err.Error())
			return nil, err
		}

		log.Printf("sendPayment response from maker to taker:%v",spew.Sdump(resp))


		return &pb.ResolveResp{
			Preimage: string(resp.PaymentPreimage[:]),
		}, nil
	}

	// If we are here we are the maker
	log.Printf("Maker code")

	return &pb.ResolveResp{
		Preimage: string(deal.preImage[:]),
	}, nil

}


func newServer(p2pServer *P2PServer) *swapResolverServer {
	s := &swapResolverServer{
		p2pServer: p2pServer,
	}
	return s
}


type P2PServer struct {
	xuPeer 	   	pb.P2PClient
	lnLTC 		lnrpc.LightningClient
	lnBTC		lnrpc.LightningClient
	mu     	    sync.Mutex // protects data structure
}

// TakeOrder is called to initiate a swap between maker and taker
// it is a temporary service needed until the integration with XUD
// intended to be called from CLI to simulate order taking by taker
func (s *P2PServer) TakeOrder(ctx context.Context, req *pb.TakeOrderReq) (*pb.TakeOrderResp, error){

	log.Printf("TakeOrder (maker) stating with %v: ",spew.Sdump(req))

	ctxt, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var (
		info *lnrpc.GetInfoResponse
		err 	error
	)

	switch req.TakerCoin{
	case pb.CoinType_BTC:
		info, err = s.lnBTC.GetInfo(ctxt, &lnrpc.GetInfoRequest{})
		if err != nil{
			return nil,fmt.Errorf("Can't getInfo from BTC LND - %v", err)
		}

	case pb.CoinType_LTC:
		info, err = s.lnLTC.GetInfo(ctxt, &lnrpc.GetInfoRequest{})
		if err != nil{
			return nil,fmt.Errorf("Can't getInfo from LTC LND - %v", err)
		}

	}
	log.Printf("%v getinfo:%v",req.TakerCoin.String(),spew.Sdump(info))
	spew.Sdump(info)

	newDeal := deal{
		myRole:			Taker,
		orderId:		req.Orderid,
		takerDealId:	uniuri.New(),
		takerAmount:	req.TakerAmount,
		takerCoin:		req.TakerCoin,
		takerPubKey:	info.IdentityPubkey,
		makerAmount:	req.MakerAmount,
		makerCoin:		req.MakerCoin,
		createTime:     time.Now(),
	}

	log.Printf("Suggesting deal to peer. deal data %v: ",spew.Sdump(newDeal))

	suggestDealResp, err := s.xuPeer.SuggestDeal(ctx, &pb.SuggestDealReq{
		Orderid: 		newDeal.orderId,
		TakerCoin: 		newDeal.takerCoin,
		TakerAmount:	newDeal.takerAmount,
		TakerDealId:	newDeal.takerDealId,
		TakerPubkey:	newDeal.takerPubKey,
		MakerCoin:		newDeal.makerCoin,
		MakerAmount: 	newDeal.makerAmount,
	})
	if err != nil{
		err = fmt.Errorf("SuggestDeal failed with %v", err)
		return nil, err
	}

	newDeal.makerDealId = suggestDealResp.MakerDealId
	newDeal.makerPubKey = suggestDealResp.MakerPubkey
	copy(newDeal.hash[:],suggestDealResp.RHash[:32])

	log.Printf("Deal Agreed with maker. deal data %v: ",spew.Sdump(newDeal))

	// TODO: protect data structure
	deals = append(deals, &newDeal)
	newDeal.executeTime = time.Now()


	swapResp, err := s.xuPeer.Swap(ctx, &pb.SwapReq{
		MakerDealId:  suggestDealResp.MakerDealId,
	})
	if err != nil{
		err = fmt.Errorf("Swap failed with %v", err)
		return nil, err
	}


	ret := &pb.TakeOrderResp{
		RPreimage:  swapResp.RPreimage,
	}
	return ret, nil
}
// SuggestDeal is called by the taker to inform the maker that he
// would like to execute a swap. The maker may reject the request
// for now, the maker can only accept/reject and can't rediscuss the
// deal or suggest partial amount. If accepted the maker should respond
// with a hash that would be used for teh swap.
func (s *P2PServer) SuggestDeal(ctx context.Context, req *pb.SuggestDealReq) (*pb.SuggestDealResp, error){

	log.Printf("SuggestDeal (taker) stating with %v: ",spew.Sdump(req))

	ctxt, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var (
		info *lnrpc.GetInfoResponse
		err 	error
	)

	switch req.MakerCoin{
	case pb.CoinType_BTC:
		info, err = s.lnBTC.GetInfo(ctxt, &lnrpc.GetInfoRequest{})
		if err != nil{
			return nil,fmt.Errorf("Can't getInfo from BTC LND - %v", err)
		}

	case pb.CoinType_LTC:
		info, err = s.lnLTC.GetInfo(ctxt, &lnrpc.GetInfoRequest{})
		if err != nil{
			return nil,fmt.Errorf("Can't getInfo from LTC LND - %v", err)
		}
	}
	log.Printf("%v getinfo:%v",req.TakerCoin.String(),spew.Sdump(info))
	spew.Sdump(info)



	newDeal := deal{
		myRole:			Maker,
		orderId:		req.Orderid,
		takerDealId:	req.TakerDealId,
		takerAmount:	req.TakerAmount,
		takerCoin:		req.TakerCoin,
		takerPubKey: 	req.TakerPubkey,
		makerPubKey: 	info.IdentityPubkey,
		makerAmount:	req.MakerAmount,
		makerCoin:		req.MakerCoin,
		makerDealId: 	uniuri.New(),
		hash: 			[32]byte{
			0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11,
			0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11,
			0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11,
			0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11,
		},
		createTime:		time.Now(),
	}

	// create preimage and a hash for the deal
	if _, err := rand.Read(newDeal.preImage[:]); err != nil {
		return nil, fmt.Errorf("Cant create preimage - %v", err)
	}

	newDeal.hash = sha256.Sum256(newDeal.preImage[:])


	deals = append(deals, &newDeal)

	ret := pb.SuggestDealResp{
		Orderid: 	newDeal.orderId,
		RHash:		newDeal.hash[:],
		MakerDealId:newDeal.makerDealId,
		MakerPubkey:newDeal.makerPubKey,
	}
	return &ret, nil
}
// Swap initiates the swap. It is called by the taker to confirm that
// he has the hash and confirm the deal.
func (s *P2PServer) Swap(ctx context.Context, req *pb.SwapReq) (*pb.SwapResp, error) {
	var deal *deal

	log.Printf("Swap (maker) stating with: %v ",spew.Sdump(req))

	for _, d := range deals{
		if d.makerDealId == req.MakerDealId{
			deal = d
			break
		}
	}

	if deal == nil{
		return nil, fmt.Errorf("Something went wrong. Can't find maker deal in swap")
	}

	cmdLnd := s.lnLTC

	switch deal.makerCoin{
	case pb.CoinType_BTC:

	case pb.CoinType_LTC:
		cmdLnd = s.lnBTC

	}

	lncctx := context.Background()
	resp, err := cmdLnd.SendPaymentSync(lncctx,&lnrpc.SendRequest{
		DestString:deal.takerPubKey,
		Amt:deal.takerAmount,
		PaymentHash:deal.hash[:],
	})

	if err != nil{
		return nil, fmt.Errorf("Got error sending  %d %v by maker - %v",
			deal.takerAmount,deal.takerCoin.String(),err)
	}

	if resp.PaymentError != ""{
		return nil, fmt.Errorf("Got error sending %d %v by maker - %v",
			deal.takerAmount,deal.takerCoin.String(), resp.PaymentError)
	}

	log.Printf("sendPayment response:%v",spew.Sdump(resp))

	ret := &pb.SwapResp{
		RPreimage: 	deal.preImage[:],
	}
	return ret, nil
}

func newP2PServer(xuPeer pb.P2PClient, lnLTC lnrpc.LightningClient, lnBTC lnrpc.LightningClient) *P2PServer {
	s := &P2PServer{
		xuPeer: xuPeer,
		lnLTC: lnLTC,
		lnBTC: lnBTC,
	}
	return s
}

func fatal(err error) {
	fmt.Fprintf(os.Stderr, "[lncli] %v\n", err)
	os.Exit(1)
}

func getClient(ctx *cli.Context) (pb.P2PClient, func()) {
	conn := getPeerConn(ctx, false)

	cleanUp := func() {
		conn.Close()
	}

	return pb.NewP2PClient(conn), cleanUp
}

func getPeerConn(ctx *cli.Context, skipMacaroons bool) *grpc.ClientConn {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	conn, err := grpc.Dial(ctx.GlobalString("peer"), opts...)

	if err != nil {
		fatal(err)
	}

	return conn
}



func main() {

	app := cli.NewApp()
	app.Name = "xud"
	app.Version = fmt.Sprintf("%s commit=%s", "0.0.1", Commit)
	app.Usage = "Use me to simulate order taking by the taker"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "listen",
			Value: "localhost:10000",
			Usage: "An ip:port to listen for peer connections",
		},
		cli.StringFlag{
			Name:  "peer",
			Value: "a.b.com:10000",
			Usage: "A host:port of peer xu daemon",
		},
		cli.StringFlag{
			Name:  "lnd-rpc-ltc",
			Value: "localhost:10001",
			Usage: "RPC host:port of LND connected to LTC chain",
		},
		cli.StringFlag{
			Name:  "lnd-rpc-btc",
			Value: "localhost:10002",
			Usage: "RPC host:port of LND connected to LTC chain",
		},

	}

	app.Action = func(c *cli.Context) error {
		log.Printf("Server starting")

		lnLTC, close := getLNDClient(c,"lnd-rpc-ltc")
		defer close()

		lnBTC, close := getLNDClient(c,"lnd-rpc-btc")
		defer close()

		xuPeer, close := getClient(c)
		defer close()
		log.Printf("Got peer connection")

		listen := c.GlobalString("listen")
		lis, err := net.Listen("tcp", listen)
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		log.Printf("listening on %s", listen)

		var opts []grpc.ServerOption
		grpcServer := grpc.NewServer(opts...)
		p2pServer := newP2PServer(xuPeer, lnLTC, lnBTC)
		pb.RegisterSwapResolverServer(grpcServer, newServer(p2pServer))
		pb.RegisterP2PServer(grpcServer, p2pServer)
		log.Printf("Server ready")
		grpcServer.Serve(lis)

		return nil

	}

	if err := app.Run(os.Args); err != nil {
		fatal(err)
	}

}

func getLNDClient(ctx *cli.Context, name string) (lnrpc.LightningClient, func()) {
	conn := getLNDClientConn(ctx, false, name)

	cleanUp := func() {
		conn.Close()
	}

	return lnrpc.NewLightningClient(conn), cleanUp
}

func getLNDClientConn(ctx *cli.Context, skipMacaroons bool, name string) *grpc.ClientConn {
	creds, err := credentials.NewClientTLSFromFile(defaultTLSCertPath, "")
	if err != nil {
		fatal(err)
	}

	// Create a dial options array.
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(creds),
	}


	conn, err := grpc.Dial(ctx.String(name), opts...)
	if err != nil {
		fatal(err)
	}

	return conn
}