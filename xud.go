
package main

import (
	"fmt"
	//"io"
	//"io/ioutil"
	"log"
	//"math"
	"net"
	"sync"
	//"time"

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
)

type swapResolverServer struct {
	mu         sync.Mutex // protects data structure
}

type role int

const (
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

	deals []deal

)


// GetFeature returns the feature at the given point.
func (s *swapResolverServer) ResolveHash(ctx context.Context, point *pb.ResolveReq) (*pb.ResolveResp, error) {

		preimage := [32]byte{
			0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11,
			0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11,
			0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11,
			0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11,
		}

		return &pb.ResolveResp{
		Preimage: string(preimage[:]),
	}, nil
}


func newServer() *swapResolverServer {
	s := &swapResolverServer{}
	return s
}


type P2PServer struct {
	peer 	   pb.P2PClient
	mu         sync.Mutex // protects data structure
}

// TakeOrder is called to initiate a swap between maker and taker
// it is a temporary service needed until the integration with XUD
// intended to be called from CLI to simulate order taking by taker
func (s *P2PServer) TakeOrder(ctx context.Context, req *pb.TakeOrderReq) (*pb.TakeOrderResp, error){

	log.Printf("TakeOrder stating with %v: ",spew.Sdump(req))
	newDeal := deal{
		myRole:			Taker,
		orderId:		req.Orderid,
		takerDealId:	uniuri.New(),
		takerAmount:	req.TakerAmount,
		takerCoin:		req.TakerCoin,
		makerAmount:	req.MakerAmount,
		makerCoin:		req.MakerCoin,
		createTime:     time.Now(),
	}

	log.Printf("Suggesting deal to peer. deal data %v: ",spew.Sdump(newDeal))

	// TODO: takerPubKey

	suggestDealResp, err := s.peer.SuggestDeal(ctx, &pb.SuggestDealReq{
		Orderid: 		newDeal.orderId,
		TakerCoin: 		newDeal.takerCoin,
		TakerAmount:	req.TakerAmount,
		TakerDealId:	newDeal.takerDealId,
		TakerPubkey:	newDeal.takerPubKey,
		MakerCoin:		req.MakerCoin,
		MakerAmount: 	req.MakerAmount,
	})
	if err != nil{
		err = fmt.Errorf("SuggestDeal failed with %v", err)
		return nil, err
	}

	newDeal.makerDealId = suggestDealResp.MakerDealId
	newDeal.makerPubKey = suggestDealResp.MakerPubkey
	copy(newDeal.hash[:],suggestDealResp.RHash[:32])

	log.Printf("Deal Agreed with maker. deal data %v: ",spew.Sdump(newDeal))

	deals = append(deals, newDeal)


	swapResp, err := s.peer.Swap(ctx, &pb.SwapReq{
		MakerDealId:  suggestDealResp.MakerDealId,
	})
	if err != nil{
		log.Fatal("Swap failed with %v", err)
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

	newDeal := deal{
		myRole:			Maker,
		orderId:		req.Orderid,
		takerDealId:	req.TakerDealId,
		takerAmount:	req.TakerAmount,
		takerCoin:		req.TakerCoin,
		takerPubKey: 	req.TakerPubkey,
		makerPubKey: 	"TODO",
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

	//TODO: makerPubKey,hash, preimage

	deals = append(deals, newDeal)

	ret := pb.SuggestDealResp{
		Orderid: 	newDeal.orderId,
		RHash:		newDeal.hash[:],
		MakerDealId:newDeal.makerDealId,
	}
	return &ret, nil
}
// Swap initiates the swap. It is called by the taker to confirm that
// he has the hash and confirm the deal.
func (s *P2PServer) Swap(ctx context.Context, req *pb.SwapReq) (*pb.SwapResp, error) {
	deal := deals[0]
	// TODO - find the deal base on deal ID
	ret := &pb.SwapResp{
		RPreimage: 	deal.preImage[:],
	}
	return ret, nil
}

func newP2PServer(peer pb.P2PClient) *P2PServer {
	s := &P2PServer{
		peer: peer,
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
	}

	app.Action = func(c *cli.Context) error {
		log.Printf("Server starting")

		peer, close := getClient(c)
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
		pb.RegisterSwapResolverServer(grpcServer, newServer())
		pb.RegisterP2PServer(grpcServer, newP2PServer(peer))
		log.Printf("Server ready")
		grpcServer.Serve(lis)

		return nil

	}

	if err := app.Run(os.Args); err != nil {
		fatal(err)
	}

}
