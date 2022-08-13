package main

import (
	"context"
	"log"
	"math/rand"
	"net"
	"os"
	"time"

	pb "gozab/gozab"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

type Ack struct {
	valid   bool
	serial  int32
	epoch   int32
	counter int32
}

var (
	vs *grpc.Server
	ls *grpc.Server
	fs *grpc.Server

	serverNum   = 5
	userPort    = "localhost:50056"
	serverPorts = []string{"localhost:50051", "localhost:50052", "localhost:50053", "localhost:50054", "localhost:50055"}
	serverMap   = map[string]int32{"localhost:50051": 0, "localhost:50052": 1, "localhost:50053": 2, "localhost:50054": 3, "localhost:50055": 4}

	propBuffers   [5]chan *pb.PropTxn
	commitBuffers [5]chan *pb.CommitTxn
	ackBuffer     chan Ack // size 5 queue
	beatBuffer    chan int // size 5 buffer

	upFollowers = []bool{true, true, true, true, true}
	upNum       = serverNum

	// for leader convenience only
	lastEpoch int32 = 0
	lastCount int32 = 0

	// for follower use
	lastEpochProp  int32 = 0
	lastLeaderProp int32 = 0

	pStorage []*pb.PropTxn
	dStruct  map[string]int32

	voted          chan int32
	followerHolder chan bool
)

type leaderServer struct {
	pb.UnimplementedLeaderUserServer
}

type followerServer struct {
	pb.UnimplementedFollowerLeaderServer
}

type voterServer struct {
	pb.UnimplementedVoterCandidateServer
}

type stateEpoch struct {
	state bool
	epoch int32
}

type stateHist struct {
	state bool
	hist  []*pb.PropTxn
}

func main() {
	pStorage = make([]*pb.PropTxn, 0)
	dStruct = make(map[string]int32)
	port := os.Args[1]
	r := ElectionRoutine(port, serverMap[port])
	for {
		if r == "elect" {
			r = ElectionRoutine(port, serverMap[port])
		} else if r == "lead" {
			r = LeaderRoutine(port)
		} else if r == "follow" {
			r = FollowerRoutine(port)
		}
	}
}

// Leader: implementation of user Store handler
func (s *leaderServer) Store(ctx context.Context, in *pb.Vec) (*pb.Empty, error) {
	log.Printf("Leader received user request")
	for i := 0; i < serverNum; i++ {
		propBuffers[i] <- &pb.PropTxn{E: lastEpoch, Transaction: &pb.Txn{V: &pb.Vec{Key: in.GetKey(), Value: in.GetValue()}, Z: &pb.Zxid{Epoch: lastEpoch, Counter: lastCount}}}
	}
	lastCount++
	return &pb.Empty{Content: "Leader recieved your request"}, nil
}

// Leader: implementation of user Retrieve handler
func (s *leaderServer) Retrieve(ctx context.Context, in *pb.GetTxn) (*pb.ResultTxn, error) {
	log.Printf("Leader received user retrieve")
	return &pb.ResultTxn{Value: dStruct[in.GetKey()]}, nil
}

func LeaderRoutine(port string) string {
	for i := 0; i < serverNum; i++ {
		propBuffers[i] = make(chan *pb.PropTxn)
		commitBuffers[i] = make(chan *pb.CommitTxn)
	}
	ackBuffer = make(chan Ack, 5)

	go FollowerRoutine(port)

	messengerStat := make(chan string)
	go MessengerRoutine(serverPorts[0], 0, messengerStat)
	go MessengerRoutine(serverPorts[1], 1, messengerStat)
	go MessengerRoutine(serverPorts[2], 2, messengerStat)
	go MessengerRoutine(serverPorts[3], 3, messengerStat)
	go MessengerRoutine(serverPorts[4], 4, messengerStat)

	go AckToCmtRoutine()

	go registerL()

	<-messengerStat
	log.Printf("quorum dead")
	ls.Stop()
	return "elect"
}

func registerL() {
	// listen user
	lis, err := net.Listen("tcp", userPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	ls = grpc.NewServer()
	pb.RegisterLeaderUserServer(ls, &leaderServer{})
	log.Printf("leader listening at %v", lis.Addr())
	if err := ls.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

// Voter: implementation of AskVote handler
func (s *voterServer) AskVote(ctx context.Context, in *pb.Epoch) (*pb.Vote, error) {
	log.Printf("Voter received candidate vote request")
	select {
	case <-voted:
		log.Printf("not voted yet, now vote")
		return &pb.Vote{Voted: true}, nil
	default:
		log.Printf("voted, do not vote")
		return &pb.Vote{Voted: false}, nil
	}

}

// Voter: implementation of NEWEPOCH handler
func (s *voterServer) NewEpoch(ctx context.Context, in *pb.Epoch) (*pb.EpochHist, error) {
	// check new epoch
	if lastEpochProp >= in.GetEpoch() {
		return nil, status.Errorf(codes.InvalidArgument,
			"new epoch not new enough")
	}
	lastEpochProp = in.GetEpoch()
	// acknowledge new epoch proposal
	return &pb.EpochHist{Epoch: lastLeaderProp, Hist: pStorage}, nil
}

// Voter: implementation of NEWLEADER handler
func (s *voterServer) NewLeader(ctx context.Context, in *pb.EpochHist) (*pb.Vote, error) {
	// check new leader's epoch
	if in.GetEpoch() != lastEpochProp {
		followerHolder <- false
		return &pb.Vote{Voted: false}, nil
	}
	// update last leader and pStorage
	lastLeaderProp = in.GetEpoch()
	pStorage = in.GetHist()
	for _, v := range pStorage {
		v.E = in.GetEpoch()
	}
	// acknowledge NEWLEADER proposal
	return &pb.Vote{Voted: true}, nil
}

// Voter: implementation of CommitNewLeader handler
func (s *voterServer) CommitNewLeader(ctx context.Context, in *pb.Empty) (*pb.Empty, error) {
	// update dStruct
	for _, v := range pStorage {
		dStruct[v.Transaction.V.Key] = v.Transaction.V.Value
	}
	// TODO: add epoch num to commit messages, if not match lift holder with false
	followerHolder <- true

	return &pb.Empty{}, nil
}

func ElectionRoutine(port string, serial int32) string {
	voted = make(chan int32, 1)
	voted <- 0
	signal := make(chan int32, 1)
	go registerV(port, signal)
	<-signal // wait for server registration
	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(10)*100 + 600 // n will be between 600 and 1500
	log.Printf("entering random sleep...")
	time.Sleep(time.Duration(n) * time.Millisecond)
	log.Printf("woke up: %v", time.Duration(n)*time.Millisecond)
	// time to check
	select {
	// didn't receive request, try to be leader
	case <-voted:
		log.Printf("voted to my self")
		// initialize routines
		electionHolder := make(chan stateEpoch, 4)
		voteRequestResultBuffer := make(chan stateEpoch, 4)
		synchronizationHolder := make(chan stateHist, 4)
		ackeResultBuffer := make(chan stateHist, 4)
		commitldHolder := make(chan bool, 4)
		ackldResultBuffer := make(chan bool, 4)
		for i := 0; i < 5; i++ {
			if int32(i) != serial {
				go ElectionMessengerRoutine(serverPorts[i], electionHolder, voteRequestResultBuffer, synchronizationHolder, ackeResultBuffer, ackldResultBuffer, commitldHolder)
			}
		}

		log.Printf("checking vote request results...")
		voteCount := 0
		var latestE int32 = -1
		for i := 0; i < 4; i++ {
			result := <-voteRequestResultBuffer
			if result.state {
				log.Printf("valid vote")
				voteCount++
				if result.epoch > latestE {
					latestE = result.epoch
				}
			}
		}
		lastEpoch = latestE + 1
		if voteCount <= serverNum/2 {
			log.Printf("insufficient vote, restarting election...")
			for i := 0; i < 4; i++ {
				electionHolder <- stateEpoch{false, -1}
			}
			vs.Stop()
			return "elect"
		} else {
			log.Printf("lifting election holder...")
			for i := 0; i < 4; i++ {
				electionHolder <- stateEpoch{true, lastEpoch}
			}
		}

		log.Printf("checking ACK-E results...")
		ackeCount := 0
		var latestHist = []*pb.PropTxn{{E: -1, Transaction: &pb.Txn{V: &pb.Vec{Key: "", Value: -1}, Z: &pb.Zxid{Epoch: -1, Counter: -1}}}}
		for i := 0; i < 4; i++ {
			result := <-ackeResultBuffer
			if result.state {
				log.Printf("valid ACK-E")
				ackeCount++
				if len(result.hist) == 0 { // system startup, no history yet
					latestHist = make([]*pb.PropTxn, 0)
				} else if result.hist[len(result.hist)-1].E > latestHist[len(latestHist)-1].E || result.hist[len(result.hist)-1].E == latestHist[len(latestHist)-1].E && result.hist[len(result.hist)-1].Transaction.Z.Counter >= latestHist[len(latestHist)-1].Transaction.Z.Counter {
					latestHist = result.hist
				}
			}
		}
		if ackeCount <= serverNum/2 {
			log.Printf("insufficient acke, restarting election...")
			for i := 0; i < 4; i++ {
				synchronizationHolder <- stateHist{false, []*pb.PropTxn{{E: -1, Transaction: &pb.Txn{V: &pb.Vec{Key: "", Value: -1}, Z: &pb.Zxid{Epoch: -1, Counter: -1}}}}}
			}
			vs.Stop()
			return "elect"
		} else {
			log.Printf("lifting synchronization holder...")
			for i := 0; i < 4; i++ {
				synchronizationHolder <- stateHist{true, latestHist}
			}
		}

		log.Printf("proceeding to phase 2 as a prospective leader...")

		log.Printf("checking NEWLEADER ack results...")
		ackldCount := 0
		for i := 0; i < 4; i++ {
			result := <-ackldResultBuffer
			if result {
				log.Printf("valid ACK-LD")
				ackldCount++
			}
		}
		if ackldCount <= serverNum/2 {
			log.Printf("insufficient ackld, restarting election...")
			for i := 0; i < 4; i++ {
				commitldHolder <- false
			}
			vs.Stop()
			return "elect"
		} else {
			log.Printf("lifting commit holder...")
			for i := 0; i < 4; i++ {
				commitldHolder <- true
			}
		}

		log.Printf("proceeding to phase 3 as an estabilished leader...")
		vs.Stop()
		return "lead"
	default:
		select {
		case r := <-followerHolder:
			if !r {
				log.Printf("followerHolder returned false, restarting election...")
				vs.Stop()
				return "elect"
			}

			log.Printf("proceeding to phase 3 as a follower...")
			vs.Stop()
			return "follow"
		case <-time.After(5 * time.Second):
			log.Printf("time out, restarting election...")
			vs.Stop()
			return "elect"
		}
	}
}

func ElectionMessengerRoutine(port string, electionHolder chan stateEpoch, voteRequestResultBuffer chan stateEpoch, synchronizationHolder chan stateHist, ackeResultBuffer chan stateHist, ackldResultBuffer chan bool, commitldHolder chan bool) {
	// dial and askvote
	conn, err := grpc.Dial(port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("could not connect to server %s: %v", port, err)
	}
	defer conn.Close()

	client := pb.NewVoterCandidateClient(conn)

	if !askvoteHelper(port, voteRequestResultBuffer, client) {
		return
	}

	// wait for quorum check and send new epoch
	checkElection := <-electionHolder
	if !checkElection.state {
		return
	}
	if !newepochHelper(port, checkElection.epoch, ackeResultBuffer, client) {
		return
	}

	// wait for quorum check and propose new leader
	checkSynchronization := <-synchronizationHolder
	if !checkSynchronization.state {
		return
	}
	if !newleaderHelper(port, checkSynchronization.hist, ackldResultBuffer, client) {
		return
	}

	// wait for quorum check and commit new leader
	if !<-commitldHolder {
		return
	}
	commitldHelper(port, client)
}

func newleaderHelper(port string, latestHist []*pb.PropTxn, ackldResultBuffer chan bool, client pb.VoterCandidateClient) bool {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := client.NewLeader(ctx, &pb.EpochHist{Epoch: lastEpoch, Hist: latestHist})
	if err != nil {
		// this messenger is dead
		ackldResultBuffer <- false
		log.Printf("failed to ask NEWLEADER ack from %s: %v", port, err)
		return false
	}
	if !r.GetVoted() {
		ackldResultBuffer <- false
		log.Printf("%s refused to give NEWLEADER ack", port)
		return false
	}
	ackldResultBuffer <- true
	return true
}

func commitldHelper(port string, client pb.VoterCandidateClient) bool {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	_, err := client.CommitNewLeader(ctx, &pb.Empty{})
	if err != nil {
		log.Printf("failed to send Commit-LD to %s: %v", port, err)
		return false
	}
	return true
}

func askvoteHelper(port string, voteRequestResultBuffer chan stateEpoch, client pb.VoterCandidateClient) bool {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := client.AskVote(ctx, &pb.Epoch{Epoch: lastEpochProp})
	if err != nil {
		voteRequestResultBuffer <- stateEpoch{false, -1}
		log.Printf("failed to ask vote from %s: %v", port, err)
		return false
	}
	if !r.GetVoted() {
		voteRequestResultBuffer <- stateEpoch{false, -1}
		log.Printf("%s refused to vote for me", port)
		return false
	}
	voteRequestResultBuffer <- stateEpoch{true, lastEpochProp}
	return true
}

func newepochHelper(port string, epoch int32, ackeResultBuffer chan stateHist, client pb.VoterCandidateClient) bool {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	hist, err := client.NewEpoch(ctx, &pb.Epoch{Epoch: epoch})
	if err != nil {
		// this messenger is dead
		ackeResultBuffer <- stateHist{false, nil}
		log.Printf("messenger quit: failed to get ACK-E from %s: %v", port, err)
		return false
	}
	ackeResultBuffer <- stateHist{true, hist.GetHist()}
	return true
}

func registerV(port string, signal chan int32) {
	// listen port
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	vs = grpc.NewServer()
	pb.RegisterVoterCandidateServer(vs, &voterServer{})
	signal <- 0
	log.Printf("voter listening at %v", lis.Addr())
	if err := vs.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

// Follower: implementation of Broadcast handler
func (s *followerServer) Broadcast(ctx context.Context, in *pb.PropTxn) (*pb.AckTxn, error) {
	log.Printf("Follower received proposa")
	// writes the proposal to local stable storage
	pStorage = append(pStorage, in)
	log.Printf("local stable storage: %+v", pStorage)
	return &pb.AckTxn{Content: "I Acknowledged"}, nil
}

// Follower: implementation of Commit handler
func (s *followerServer) Commit(ctx context.Context, in *pb.CommitTxn) (*pb.Empty, error) {
	log.Printf("Follower received commit")
	// writes the transaction from local stable storage to local data structure
	dStruct[pStorage[len(pStorage)-1].Transaction.V.Key] = pStorage[len(pStorage)-1].Transaction.V.Value
	log.Printf("local data structure: %+v", dStruct)
	return &pb.Empty{Content: "Commit message recieved"}, nil
}

// Follower: implementation of HeartBeat handler
func (s *followerServer) HeartBeat(ctx context.Context, in *pb.Empty) (*pb.Empty, error) {
	beatBuffer <- 1
	return &pb.Empty{Content: "bump"}, nil
}

func FollowerRoutine(port string) string {
	// follower receiving leader's heartbeat
	beatBuffer = make(chan int, 5)
	leaderStat := make(chan string)
	go BeatReceiver(leaderStat)

	go registerF(port)

	<-leaderStat
	log.Printf("leader dead")
	fs.Stop()
	return "elect"
}

func registerF(port string) {
	// listen leader on port
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	fs = grpc.NewServer()
	pb.RegisterFollowerLeaderServer(fs, &followerServer{})
	log.Printf("server listening at %v", lis.Addr())
	if err := fs.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func AckToCmtRoutine() {
	for {
		// Collect acknowledgements from Messengers
		// Update upFollowers statistics
		for i := 0; i < serverNum; i++ {
			ack := <-ackBuffer
			if !ack.valid && upFollowers[ack.serial] {
				upFollowers[ack.serial] = false
				upNum--
			}
		}

		// Check if quorum dead
		if upNum <= serverNum/2 {
			log.Printf("quorum dead")
			return
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
func MessengerRoutine(port string, serial int32, messengerStat chan string) {
	// dial follower
	conn, err := grpc.Dial(serverPorts[serial], grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("could not connect to server %s: %v", port, err)
	}
	defer conn.Close()

	client := pb.NewFollowerLeaderClient(conn)

	// leader receiving follower's heartbeat
	followerStat := make(chan string)
	go BeatSender(client, followerStat)

	go propcmt(client, port, serial)

	<-followerStat
	upFollowers[serial] = false
	upNum--
	log.Printf("follower %s down, messenger quit", port)
	// check quorum dead
	if upNum <= serverNum/2 {
		messengerStat <- "returned"
	}
}

func propcmt(client pb.FollowerLeaderClient, port string, serial int32) {
	for {
		prop(client, port, serial)
		cmt(client, port, serial)
	}
}

func prop(client pb.FollowerLeaderClient, port string, serial int32) {
	// send proposal
	proposal := <-propBuffers[serial]
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	r, err := client.Broadcast(ctx, proposal)
	cancel() // defer?
	if err != nil {
		log.Printf("could not broadcast to server %s: %v", port, err)
		status, ok := status.FromError(err)
		if ok {
			if status.Code() == codes.DeadlineExceeded {
				log.Printf("Server %s timeout, Messenger exit", port)
			}
		}
		// dead Messenger
		log.Printf("messenger on %s initiated death mode", port)
		ackBuffer <- Ack{false, serial, -1, -1}
		for {
			<-propBuffers[serial]
			ackBuffer <- Ack{false, serial, -1, -1}
		}
	}
	if r.GetContent() == "I Acknowledged" {
		ackBuffer <- Ack{true, serial, proposal.GetTransaction().GetZ().Epoch, proposal.GetTransaction().GetZ().Counter}
	}
}

func cmt(client pb.FollowerLeaderClient, port string, serial int32) {
	// send commit
	commit := <-commitBuffers[serial]
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	r, err := client.Commit(ctx, commit)
	cancel() // defer?
	if err != nil {
		log.Printf("could not issue commit to server %s: %v", port, err)
		status, ok := status.FromError(err)
		if ok {
			if status.Code() == codes.DeadlineExceeded {
				log.Printf("Server %s timeout, Messenger exit", port)
			}
		}
		// dead Messenger
		log.Printf("messenger on %s initiated death mode", port)
		for {
			<-propBuffers[serial]
			ackBuffer <- Ack{false, serial, -1, -1}
		}
	}
	if r.GetContent() == "Commit message recieved" {
		log.Printf("Commit feedback recieved from %s", port)
	}
}

func BeatSender(client pb.FollowerLeaderClient, followerStat chan string) {
	for {
		time.Sleep(100 * time.Millisecond)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		_, err := client.HeartBeat(ctx, &pb.Empty{Content: "Beat"})
		cancel() // defer?
		if err != nil {
			followerStat <- "dead"
			return
		}
	}
}

func BeatReceiver(leaderStat chan string) {
	for {
		select {
		case <-beatBuffer:
			log.Printf("beat")
		case <-time.After(5 * time.Second):
			leaderStat <- "dead"
			return
		}
	}
}
