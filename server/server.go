
package main

import (
	"flag"
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
)

var (
	tls        = flag.Bool("tls", false, "Connection uses TLS if true, else plain TCP")
	//certFile   = flag.String("cert_file", "", "The TLS cert file")
	//keyFile    = flag.String("key_file", "", "The TLS key file")
	port       = flag.Int("port", 10000, "The server port")
)

type swapResolverServer struct {
	mu         sync.Mutex // protects data structure
}

// GetFeature returns the feature at the given point.
func (s *swapResolverServer) ResolveHash(ctx context.Context, point *pb.ResolveReq) (*pb.ResolveResp, error) {

	//if pd.Amount == 3000000{
	//	preimage := [32]byte{
	//		0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11,
	//		0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11,
	//		0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11,
	//		0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11,
	//	}
	//
	//} else{

		return &pb.ResolveResp{
		Preimage: "1111111111111111111111111111111111111111111111111111111111111111",
	}, nil
}


func newServer() *swapResolverServer {
	s := &swapResolverServer{}
	return s
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption
	if *tls {
		//if *certFile == "" {
		//	*certFile = testdata.Path("server1.pem")
		//}
		//if *keyFile == "" {
		//	*keyFile = testdata.Path("server1.key")
		//}
		//creds, err := credentials.NewServerTLSFromFile(*certFile, *keyFile)
		//if err != nil {
		//	log.Fatalf("Failed to generate credentials %v", err)
		//}
		//opts = []grpc.ServerOption{grpc.Creds(creds)}
	}
	grpcServer := grpc.NewServer(opts...)
	pb.RegisterSwapResolverServer(grpcServer, newServer())
	grpcServer.Serve(lis)
}
