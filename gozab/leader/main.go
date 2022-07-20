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

type Ack struct {
	serial  int32
	epoch   int32
	counter int32
}

var (
	serverPorts   = []string{"", "localhost:50051", "localhost:50052", "localhost:50053", "localhost:50054", "localhost:50055"}
	propBuffers   []chan *pb.PropTxn // size 1 queue
	commitBuffers []chan *pb.CommitTxn
	ackBuffer     chan Ack
	lastEpoch     int32 = 1
	lastCount     int32 = 0
)

type leaderServer struct {
	pb.UnimplementedClientConsoleServer
}

// implementation of user request handler
func (s *leaderServer) SendRequest(ctx context.Context, in *pb.Vec) (*pb.Empty, error) {
	log.Printf("Leader received user request\n")
	for i := 1; i < serverNum; i++ {
		propBuffers[i] <- &pb.PropTxn{E: lastEpoch, Transaction: &pb.Txn{V: &pb.Vec{Key: in.GetKey(), Value: in.GetValue()}, Z: &pb.Zxid{Epoch: lastEpoch, Counter: lastCount}}}
	}
	lastCount++
	return &pb.Empty{Content: "Leader recieved your request"}, nil
}

func main() {
	for i := 0; i < serverNum; i++ {
		propBuffers[i] = make(chan *pb.PropTxn)
		commitBuffers[i] = make(chan *pb.CommitTxn)
	}

	LaunchMessengerRoutines()

	go AckToCommitRoutine() // TODO: finish this!

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

// simply launch 4 Messenger routines
func LaunchMessengerRoutines() {
	for i := 0; i < serverNum; i++ {
		go MessengerRoutine(serverPorts[i], int32(i))
	}
}

func AckToCommitRoutine() {
	for {

		// Collect acknowledgements
		ackCount := 0
		for ackCount <= 2 {
			<-ackBuffer
			ackCount++
		}

		// Tell the Messengers to send commits
		for i := 0; i < serverNum; i++ {
			commitBuffers[i] <- &pb.CommitTxn{Content: "Please commit"}
		}
	}
}

// CORE BROACAST FUNCTION!
func MessengerRoutine(port string, serial int32) {

	// Dial the port
	conn, err := grpc.Dial(serverPorts[serial], grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect server %s: %v", port, err)
	}
	defer conn.Close()

	client := pb.NewSimulationClient(conn)

	// Coulson: what does this line do?
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	for {
		// send proposal
		proposal := <-propBuffers[serial]
		rb, errb := client.Broadcast(ctx, proposal)
		if errb != nil {
			log.Fatalf("could not broadcast to server %s: %v", port, errb)
		}

		// receive acknowledgement
		if rb.GetContent() == "I Acknowledged" {
			ackBuffer <- Ack{serial, proposal.GetTransaction().GetZ().Epoch, proposal.GetTransaction().GetZ().Counter}
		}

		commit := <-commitBuffers[serial]
		rc, errc := client.Commit(ctx, commit)
		if errc != nil {
			log.Fatalf("could not issue commit to server %s: %v", port, errc)
		}
		if rc.GetContent() == "Commit message recieved" {
			log.Printf("Commit feedback recieved from %s", port)
		}
	}
}
