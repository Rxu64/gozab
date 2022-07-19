package main

import (
	"context"
	"log"
	"net"
	"time"

	pb "gozab/gozab"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const serverNum = 4

type Vec struct {
	key   string
	value int32
}

var (
	serverPorts = []string{"", "localhost:50051", "localhost:50052", "localhost:50053", "localhost:50054", "localhost:50055"}
	vecBuff     chan Vec // size 1 queue
)

type leaderServer struct {
	pb.UnimplementedClientConsoleServer
}

// implementation of user request handler
func (s *leaderServer) SendRequest(ctx context.Context, in *pb.Vec) (*pb.Empty, error) {
	log.Printf("Leader received user request\n")
	vecBuff <- Vec{in.GetKey(), in.GetValue()}
	return &pb.Empty{Content: "Leader recieved your request"}, nil
}

func main() {
	var clients [serverNum + 1]pb.SimulationClient
	vecBuff = make(chan Vec)

	go CallFollowers(clients[:]) // propose and commit

	// build 4 client connections
	for i := 1; i <= serverNum; i++ {
		conn, err := grpc.Dial(serverPorts[i], grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatalf("did not connect server%d: %v", i, err)
		}
		defer conn.Close()
		clients[i] = pb.NewSimulationClient(conn)
	}
	log.Printf("Connection to followers established\n")

	// listen user
	lis, err := net.Listen("tcp", serverPorts[5])
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterClientConsoleServer(s, &leaderServer{})
	log.Printf("leader listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func CallFollowers(clients []pb.SimulationClient) {
	var lastEpoch int32 = 1
	var lastCount int32 = 0

	for {
		vec := <-vecBuff // block if empty!
		proposal := &pb.PropTxn{E: lastEpoch, Transaction: &pb.Txn{V: &pb.Vec{Key: vec.key, Value: vec.value}, Z: &pb.Zxid{Epoch: lastEpoch, Counter: lastCount}}}

		SendProposal(proposal, clients[:])
		lastCount++

		SendCommit(clients[:])
	}
}

func SendProposal(proposal *pb.PropTxn, clients []pb.SimulationClient) {
	akgCount := 0

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	for i := 1; i <= serverNum; i++ {
		r, err := clients[i].Broadcast(ctx, proposal)
		if err != nil {
			log.Fatalf("could not broadcast to server%d: %v", i, err)
		}
		if r.GetContent() == "I Acknowledged" {
			akgCount++
		}
	}
	log.Printf("Result: %d acknowledgement received", akgCount)

}

func SendCommit(clients []pb.SimulationClient) {
	cmtCount := 0

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	for i := 1; i <= serverNum; i++ {
		r, err := clients[i].Commit(ctx, &pb.CommitTxn{Content: "Please commit"})
		if err != nil {
			log.Fatalf("could not issue commit to server%d: %v", i, err)
		}
		if r.GetContent() == "Commit message recieved" {
			cmtCount++
		}
	}
	log.Printf("Result: %d commit feedback received", cmtCount)
}
