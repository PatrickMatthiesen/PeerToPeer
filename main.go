package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"time"

	gRPC "github.com/PatrickMatthiesen/PeerToPeer/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var Name = flag.String("name", "default", "Senders name")
var Port = flag.String("port", "5400", "Own Tcp server port")
var Address = flag.String("address", "localhost", "Tcp server address")
var peerPort = flag.String("peerPort", "5401", "Peer Tcp server port")
var peerAddress = flag.String("peerAddress", "localhost", "Peer Tcp server address")

type Server struct {
	gRPC.UnimplementedHelloServiceServer // You need this line if you have a server
}

func (s *Server) Hello(ctx context.Context, in *gRPC.HelloMessage) (*gRPC.HelloMessage, error) {
	fmt.Printf("%s send: %s\n", in.Sender, in.Message)
	return &gRPC.HelloMessage{Message: "Hello " + in.Sender}, nil
}

func main() {
	flag.Parse()
	fmt.Println("Starting server...")
	// log.SetOutput(os.Stdout)

	// Create a listener on TCP port
	lis, err := net.Listen("tcp", *Address+":"+*Port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Create a gRPC server object
	s := grpc.NewServer()

	// Attach the Hello service to the server
	gRPC.RegisterHelloServiceServer(s, &Server{})

	// start sending requests to peer
	go runClient()

	// Serve gRPC server
	fmt.Println("Serving gRPC on port " + *Port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
	// code here is unreachable because grpcServer.Serve occupies the current thread.
}

func runClient() {
	fmt.Println("starting client...")
	// Set up a connection to the server.
	opts := []grpc.DialOption{
		grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	conn, err := grpc.Dial(*peerAddress+":"+*peerPort, opts...)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := gRPC.NewHelloServiceClient(conn)

	for {
		c.Hello(context.Background(), &gRPC.HelloMessage{Sender: *Name, Message: "Hope you are doing well!"})
		time.Sleep(time.Second)
	}
}
