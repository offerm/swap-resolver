package main

import (
	"flag"
	//"io"
	"log"
	//"math/rand"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	//"google.golang.org/grpc/credentials"
	pb "github.com/offerm/swap-resolver/swapresolver"
	//"google.golang.org/grpc/testdata"
)

var (
	tls                = flag.Bool("tls", false, "Connection uses TLS if true, else plain TCP")
	//caFile             = flag.String("ca_file", "", "The file containning the CA root cert file")
	serverAddr         = flag.String("server_addr", "127.0.0.1:10000", "The server address in the format of host:port")
	//serverHostOverride = flag.String("server_host_override", "x.test.youtube.com", "The server name use to verify the hostname returned by TLS handshake")
)

// printFeature gets the feature for the given point.
func printPreImage(client pb.SwapResolverClient, req *pb.ResolveReq) {
	log.Printf("Getting pre-image for hash: %v",req.Hash)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	preimage, err := client.ResolveHash(ctx, req)
	if err != nil {
		log.Fatalf("%v.ResolveHash(_) = _, %v: ", client, err)
	}
	log.Printf("Got response from Resolver: %v \n",preimage)
}

func main() {
	flag.Parse()
	var opts []grpc.DialOption
	if *tls {
		//if *caFile == "" {
		//	*caFile = testdata.Path("ca.pem")
		//}
		//creds, err := credentials.NewClientTLSFromFile(*caFile, *serverHostOverride)
		//if err != nil {
		//	log.Fatalf("Failed to create TLS credentials %v", err)
		//}
		//opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		opts = append(opts, grpc.WithInsecure())
	}
	conn, err := grpc.Dial(*serverAddr, opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()
	client := pb.NewSwapResolverClient(conn)

	printPreImage(client, &pb.ResolveReq{
		Hash: "this is my hash",
	})

}
