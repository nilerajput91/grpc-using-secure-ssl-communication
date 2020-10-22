package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	"grpc/model"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var (
	address = "expert-Inspiron-3541:8080"
)

func main() {

	flag.StringVar(&address, "a", address, "gRPC server address host:port")
	flag.Parse()
	var opts []grpc.ServerOption

	//configure tls
	creds, _ := credentials.NewServerTLSFromFile("../cert.pem", "../key.pem")

	opts = append(opts, grpc.Creds(creds))

	server := grpc.NewServer(opts...)

	model.RegisterMyMathServiceServer(server, &myMathService{})
	model.RegisterDataServiceServer(server, &myDataService{})

	lis, err := net.Listen("tcp", address)

	if err != nil {
		log.Fatal(fmt.Errorf("Unable to start the gRPC server on address %v: %v", address, err))
	}

	server.Serve(lis)
}
