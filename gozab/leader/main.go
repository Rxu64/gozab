package main

import (
	"context"
	"gozab/follower"
	"log"
	"net"
	"time"

	pb "gozab/gozab"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	serverNum = 5
	userPort  = "localhost:50056"
)

type Ack struct {
	serial  int32
	epoch   int32
	counter int32
}

var (
	serverPorts   = []string{"localhost:50051", "localhost:50052", "localhost:50053", "localhost:50054", "localhost:50055"}
	propBuffers   [5]chan *pb.PropTxn
	commitBuffers [5]chan *pb.CommitTxn
	ackBuffer     chan Ack // size 4 queue
	activity      chan *pb.Vec
	upFollowers         = []bool{true, true, true, true, true}
	lastEpoch     int32 = 1
	lastCount     int32 = 0
)

type leaderServer struct {
	pb.UnimplementedClientConsoleServer
}

// implementation of user request handler
func (s *leaderServer) SendRequest(ctx context.Context, in *pb.Vec) (*pb.Empty, error) {
	log.Printf("Leader received user request\n")
	activity <- in
	for i, up := range upFollowers {
		if up {
			propBuffers[i] <- &pb.PropTxn{E: lastEpoch, Transaction: &pb.Txn{V: &pb.Vec{Key: in.GetKey(), Value: in.GetValue()}, Z: &pb.Zxid{Epoch: lastEpoch, Counter: lastCount}}}
		}
	}
	lastCount++
	return &pb.Empty{Content: "Leader recieved your request"}, nil
}

func main() {
	for i := 0; i < serverNum; i++ {
		propBuffers[i] = make(chan *pb.PropTxn)
		commitBuffers[i] = make(chan *pb.CommitTxn)
	}
	ackBuffer = make(chan Ack, 5)
	activity = make(chan *pb.Vec)

	// Launch Routines
	go follower.FollowerRoutine([]string{"", serverPorts[4]}) // start follower on 50055

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
	pb.RegisterClientConsoleServer(s, &leaderServer{})
	log.Printf("leader listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func AckToCmtRoutine() {
	upNum := serverNum
	const timeout = time.Second
	ackFollowers := make([]bool, 5)
	for {
		// Collect acknowledgements with timeout
		<-activity     // wait for user activity
		downCount := 0 // reset count
		for i := range ackFollowers {
			ackFollowers[i] = false // reset marker
		}

		for i := 0; i < upNum; i++ {
			select {
			case ack := <-ackBuffer:
				ackFollowers[ack.serial] = true // mark this follower as acknowledged
			case <-time.After(timeout):
				downCount++
				continue
			}
		}

		for i, acknowledged := range ackFollowers {
			if !acknowledged {
				upFollowers[i] = false // mark followers who did not acknowledge as down
			}
		}

		upNum -= downCount
		if upNum < 2 {
			log.Fatalf("quorum dead")
		}

		// Tell the acknowledged Messengers to send commits
		for i, acked := range upFollowers {
			if acked {
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
		log.Printf("did not connect server %s: %v", port, err)
	}
	defer conn.Close()

	client := pb.NewSimulationClient(conn)

	for {
		// send proposal
		proposal := <-propBuffers[serial] // should add timeout
		rb, errb := client.Broadcast(context.Background(), proposal)
		if errb != nil {
			log.Printf("could not broadcast to server %s: %v", port, errb)
		}
		if rb.GetContent() == "I Acknowledged" {
			ackBuffer <- Ack{serial, proposal.GetTransaction().GetZ().Epoch, proposal.GetTransaction().GetZ().Counter}
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
