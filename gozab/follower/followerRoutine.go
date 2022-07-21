package follower

import (
	"context"
	"log"
	"net"

	pb "gozab/gozab"

	"google.golang.org/grpc"
)

var pStorage []Proposal
var dStruct map[string]int32

type Vec struct {
	key   string
	value int32
}

type Zxid struct {
	epoch   int32
	counter int32
}

type Transaction struct {
	v Vec
	z Zxid
}

type Proposal struct {
	e   int32
	txn Transaction
}

type server struct {
	pb.UnimplementedSimulationServer
}

// implementation of Broadcast handler
func (s *server) Broadcast(ctx context.Context, in *pb.PropTxn) (*pb.AckTxn, error) {
	log.Printf("Follower received proposal\n")
	// writes the proposal to local stable storage
	var prop = Proposal{in.GetE(), Transaction{Vec{in.GetTransaction().GetV().GetKey(), in.GetTransaction().GetV().GetValue()}, Zxid{in.GetTransaction().GetZ().GetEpoch(), in.GetTransaction().GetZ().GetCounter()}}}
	pStorage = append(pStorage, prop)
	log.Printf("local stable storage: %+v\n", pStorage)
	return &pb.AckTxn{Content: "I Acknowledged"}, nil
}

// implementation of Commit handler
func (s *server) Commit(ctx context.Context, in *pb.CommitTxn) (*pb.Empty, error) {
	log.Printf("Follower received commit\n")
	// writes the transaction from local stable storage to local data structure
	dStruct[pStorage[len(pStorage)-1].txn.v.key] = pStorage[len(pStorage)-1].txn.v.value
	log.Printf("local data structure: %+v\n", dStruct)
	return &pb.Empty{Content: "Commit message recieved"}, nil
}

func FollowerRoutine(args []string) {
	pStorage = make([]Proposal, 0)
	dStruct = make(map[string]int32)
	lis, err := net.Listen("tcp", args[1])
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterSimulationServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
