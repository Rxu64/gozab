package main

import (
	"context"
	"log"
	"net"
	"os"
	"time"

	pb "gozab/gozab"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

const (
	serverNum = 5
	userPort  = "localhost:50056"
)

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

type Ack struct {
	valid   bool
	serial  int32
	epoch   int32
	counter int32
}

var (
	serverPorts   = []string{"localhost:50051", "localhost:50052", "localhost:50053", "localhost:50054", "localhost:50055"}
	propBuffers   [5]chan *pb.PropTxn
	commitBuffers [5]chan *pb.CommitTxn
	ackBuffer     chan Ack // size 5 queue
	upFollowers            = []bool{true, true, true, true, true}
	lastEpoch     int32    = 1
	lastCount     int32    = 0

	pStorage []Proposal
	dStruct  map[string]int32
)

type leaderServer struct {
	pb.UnimplementedLeaderUserServer
}

type followerServer struct {
	pb.UnimplementedFollowerLeaderServer
}

func main() {
	pStorage = make([]Proposal, 0)
	dStruct = make(map[string]int32)
	if os.Args[1] == "lead" {
		LeaderRoutine(os.Args[2])
	} else if os.Args[1] == "follow" {
		FollowerRoutine(os.Args[2])
	} else {
		log.Fatalf("unknown command")
	}
}

// Leader: implementation of user Store handler
func (s *leaderServer) Store(ctx context.Context, in *pb.Vec) (*pb.Empty, error) {
	log.Printf("Leader received user request\n")
	for i, up := range upFollowers {
		if up {
			propBuffers[i] <- &pb.PropTxn{E: lastEpoch, Transaction: &pb.Txn{V: &pb.Vec{Key: in.GetKey(), Value: in.GetValue()}, Z: &pb.Zxid{Epoch: lastEpoch, Counter: lastCount}}}
		}
	}
	lastCount++
	return &pb.Empty{Content: "Leader recieved your request"}, nil
}

// Leader: implementation of user Retrieve handler
func (s *leaderServer) Retrieve(ctx context.Context, in *pb.GetTxn) (*pb.ResultTxn, error) {
	log.Printf("Leader received user retrieve\n")
	return &pb.ResultTxn{Value: dStruct[in.GetKey()]}, nil
}

func LeaderRoutine(port string) {
	for i := 0; i < serverNum; i++ {
		propBuffers[i] = make(chan *pb.PropTxn)
		commitBuffers[i] = make(chan *pb.CommitTxn)
	}
	ackBuffer = make(chan Ack, 5)

	go FollowerRoutine(port)

	go MessengerRoutine(serverPorts[0], 0)
	go MessengerRoutine(serverPorts[1], 1)
	go MessengerRoutine(serverPorts[2], 2)
	go MessengerRoutine(serverPorts[3], 3)
	go MessengerRoutine(serverPorts[4], 4)

	go AckToCmtRoutine()

	// listen user
	lis, err := net.Listen("tcp", userPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterLeaderUserServer(s, &leaderServer{})
	log.Printf("leader listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

// Follower: implementation of Broadcast handler
func (s *followerServer) Broadcast(ctx context.Context, in *pb.PropTxn) (*pb.AckTxn, error) {
	log.Printf("Follower received proposal\n")
	// writes the proposal to local stable storage
	var prop = Proposal{in.GetE(), Transaction{Vec{in.GetTransaction().GetV().GetKey(), in.GetTransaction().GetV().GetValue()}, Zxid{in.GetTransaction().GetZ().GetEpoch(), in.GetTransaction().GetZ().GetCounter()}}}
	pStorage = append(pStorage, prop)
	log.Printf("local stable storage: %+v\n", pStorage)
	return &pb.AckTxn{Content: "I Acknowledged"}, nil
}

// Follower: implementation of Commit handler
func (s *followerServer) Commit(ctx context.Context, in *pb.CommitTxn) (*pb.Empty, error) {
	log.Printf("Follower received commit\n")
	// writes the transaction from local stable storage to local data structure
	dStruct[pStorage[len(pStorage)-1].txn.v.key] = pStorage[len(pStorage)-1].txn.v.value
	log.Printf("local data structure: %+v\n", dStruct)
	return &pb.Empty{Content: "Commit message recieved"}, nil
}

func FollowerRoutine(port string) {
	// listen on port
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// register server
	s := grpc.NewServer()
	pb.RegisterFollowerLeaderServer(s, &followerServer{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func AckToCmtRoutine() {
	upNum := serverNum
	for {
		// Collect acknowledgements from Messengers
		// Update upFollowers statistics
		for i := 0; i < serverNum; i++ {
			ack := <-ackBuffer
			if !ack.valid {
				upFollowers[ack.serial] = false
				upNum--
			}
		}

		// Check if quorum dead
		if upNum <= serverNum/2 {
			log.Fatalf("quorum dead")
		}

		// Send commits to the acknowledged followers
		for i, acknowledged := range upFollowers {
			if acknowledged {
				commitBuffers[i] <- &pb.CommitTxn{Content: "Please commit"}
			}
		}
	}
}

// CORE BROACAST FUNCTION!
func MessengerRoutine(port string, serial int32) {
	// dial follower
	conn, err := grpc.Dial(serverPorts[serial], grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("could not connect to server %s: %v", port, err)
	}
	defer conn.Close()

	client := pb.NewFollowerLeaderClient(conn)

	for {
		// send proposal
		proposal := <-propBuffers[serial]
		clientDeadline := time.Now().Add(time.Duration(time.Second)) // 1 second deadline
		ctx, cancel := context.WithDeadline(context.Background(), clientDeadline)
		defer cancel()
		rb, errb := client.Broadcast(ctx, proposal)
		if errb != nil {
			log.Printf("could not broadcast to server %s: %v", port, errb)
			status, ok := status.FromError(err)
			if ok {
				if status.Code() == codes.DeadlineExceeded {
					log.Printf("Server %s timeout, Messenger exit", port)
				}
			}
			// dead Messenger
			log.Printf("messenger on %s initiated death mode", port)
			ackBuffer <- Ack{false, serial, proposal.GetTransaction().GetZ().Epoch, proposal.GetTransaction().GetZ().Counter}
			for {
				proposal = <-propBuffers[serial]
				ackBuffer <- Ack{false, serial, proposal.GetTransaction().GetZ().Epoch, proposal.GetTransaction().GetZ().Counter}
			}
		}
		if rb.GetContent() == "I Acknowledged" {
			ackBuffer <- Ack{true, serial, proposal.GetTransaction().GetZ().Epoch, proposal.GetTransaction().GetZ().Counter}
		}

		// send commit
		commit := <-commitBuffers[serial]
		rc, errc := client.Commit(context.Background(), commit)
		if errc != nil {
			log.Printf("could not issue commit to server %s: %v", port, errc)
		}
		if rc.GetContent() == "Commit message recieved" {
			log.Printf("Commit feedback recieved from %s", port)
		}
	}
}