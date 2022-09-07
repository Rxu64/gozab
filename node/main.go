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

const nodeNum = 5 // <--------------- NUMBER OF NODES

var (
	vs *grpc.Server
	ls *grpc.Server
	fs *grpc.Server

	// Constants
	userPort    = "198.22.255.14:50056"
	serverPorts = []string{"198.22.255.12:50051", "198.22.255.16:50052", "198.22.255.28:50053", "198.22.255.11:50054", "198.22.255.15:50055"}
	serverMap   = map[string]int32{"198.22.255.12:50051": 0, "198.22.255.16:50052": 1, "198.22.255.28:50053": 2, "198.22.255.11:50054": 3, "198.22.255.15:50055": 4}

	// Global channcels for broadcast-phase leader
	propChannels   [nodeNum]chan *pb.PropTxn
	commitChannels [nodeNum]chan *pb.CommitTxn
	ackChannel     chan Ack // size of nodeNum queue
	beatChannel    chan int // size of nodeNum Channel

	// upFollowersUpdateRoutine is the only code that changes those two variables directly!
	// specifically, leader routine reset this array to true by pushing negative values
	// messenger routiens report follower failure by pushing positive values
	upFollowers              = []bool{true, true, true, true, true}
	upNum                    = nodeNum
	upFollowersUpdateChannel chan int32

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

	//	for leader's use to clean up everything
	universalCleanupHolder chan int32
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
	r := Elect(port, serverMap[port])
	for {
		if r == "elect" {
			r = Elect(port, serverMap[port])
		} else if r == "lead" {
			r = Lead(port)
		} else if r == "follow" {
			r = Follow(port)
		}
	}
}

// Leader's Handler: implementation of user Store handler
func (s *leaderServer) Store(ctx context.Context, in *pb.Vec) (*pb.Empty, error) {
	log.Printf("Leader received user request")
	for i := 0; i < nodeNum; i++ {
		propChannels[i] <- &pb.PropTxn{E: lastEpoch, Transaction: &pb.Txn{V: &pb.Vec{Key: in.GetKey(), Value: in.GetValue()}, Z: &pb.Zxid{Epoch: lastEpoch, Counter: lastCount}}}
	}
	lastCount++
	return &pb.Empty{Content: "Leader recieved your request"}, nil
}

// Leader's Handler: implementation of user Retrieve handler
func (s *leaderServer) Retrieve(ctx context.Context, in *pb.GetTxn) (*pb.ResultTxn, error) {
	log.Printf("Leader received user retrieve")
	return &pb.ResultTxn{Value: dStruct[in.GetKey()]}, nil
}

func Lead(port string) string {

	for i := 0; i < nodeNum; i++ {
		propChannels[i] = make(chan *pb.PropTxn, 1)
		commitChannels[i] = make(chan *pb.CommitTxn, 1)
	}
	ackChannel = make(chan Ack, nodeNum)

	// setup upFollowers stats
	upFollowersUpdateChannel = make(chan int32, nodeNum)
	go upFollowersUpdateRoutine()

	// setup cleanup holder
	universalCleanupHolder = make(chan int32, 1)

	// initialize leader-attached local follower
	go Follow(port)

	go MessengerRoutine(serverPorts[0], 0)
	go MessengerRoutine(serverPorts[1], 1)
	go MessengerRoutine(serverPorts[2], 2)
	go MessengerRoutine(serverPorts[3], 3)
	go MessengerRoutine(serverPorts[4], 4)

	go AckToCmtRoutine()

	go serveL()

	// wait for quorum dead and cleanup signal
	<-universalCleanupHolder

	ls.Stop()
	return "elect"
}

// Register leader's grpc server exposed to user
func serveL() {
	// listen user
	lis, err := net.Listen("tcp", userPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Printf("leader listening at %v", lis.Addr())
	ls = grpc.NewServer()
	pb.RegisterLeaderUserServer(ls, &leaderServer{})
	if err := ls.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

// Voter's Handler: implementation of AskVote handler
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

// Voter's Handler: implementation of NEWEPOCH handler
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

// Voter's Handler: implementation of NEWLEADER handler
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

// Voter's Handler: implementation of CommitNewLeader handler
func (s *voterServer) CommitNewLeader(ctx context.Context, in *pb.Epoch) (*pb.Empty, error) {

	if in.GetEpoch() != lastLeaderProp {
		return &pb.Empty{}, nil // this is simply for prevent delayed delivery of messeges from previous leader
		// do not need to abort here, simply ignore it
	}
	// update dStruct
	log.Printf("updating dStruct")
	for _, v := range pStorage {
		dStruct[v.Transaction.V.Key] = v.Transaction.V.Value
	}
	followerHolder <- true

	return &pb.Empty{}, nil
}

func Elect(port string, serial int32) string {
	voted = make(chan int32, 1)
	followerHolder = make(chan bool, 1)
	voted <- 0

	sleepHolder := make(chan int32, 1)

	go serveV(port, sleepHolder)

	<-sleepHolder // make sure service established before go to sleep

	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(10)*100 + 600 // n will be between 600 and 1500
	log.Printf("entering random sleep...")
	time.Sleep(time.Duration(n) * time.Millisecond)
	log.Printf("woke up: %v", time.Duration(n)*time.Millisecond)
	// time to check
	select {
	case <-voted:
		// didn't receive request, try to be leader
		log.Printf("voted to my self")

		// initialize synchronization channels
		electionHolder := make(chan stateEpoch, 4)
		voteRequestResultChannel := make(chan stateEpoch, 4)
		synchronizationHolder := make(chan stateHist, 4)
		ackeResultChannel := make(chan stateHist, 4)
		commitldHolder := make(chan bool, 4)
		ackldResultChannel := make(chan bool, 4)
		voteHolder := make(chan int32, 4)

		// // initialize election messenger routines
		for i := 0; i < nodeNum; i++ {
			if int32(i) != serial {
				go ElectionMessengerRoutine(serverPorts[i], voteHolder, electionHolder, voteRequestResultChannel, synchronizationHolder, ackeResultChannel, ackldResultChannel, commitldHolder)
			}
		}

		for i := 0; i < 4; i++ {
			<-voteHolder // make sure election messenger routines started bofore checking results below
		}

		noreplyCount := 0

		log.Printf("checking vote request results...")
		var latestE int32 = -1
		voteNum := 4 - noreplyCount
		for i := 0; i < voteNum; i++ {
			result := <-voteRequestResultChannel
			if result.state {
				log.Printf("valid vote")
				if result.epoch > latestE {
					latestE = result.epoch
				}
			} else {
				log.Printf("invalid vote")
				noreplyCount++
			}
		}
		lastEpoch = latestE + 1
		if noreplyCount > nodeNum/2 {
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
		latestHist := make([]*pb.PropTxn, 0)
		ackeNum := 4 - noreplyCount
		for i := 0; i < ackeNum; i++ {
			result := <-ackeResultChannel
			if result.state {
				log.Printf("valid ACK-E")
				if len(result.hist) == 0 { // system startup, no history yet
					continue
				} else if len(latestHist) == 0 { // leader has history, I have no history
					latestHist = result.hist
				} else if result.hist[len(result.hist)-1].E > latestHist[len(latestHist)-1].E || result.hist[len(result.hist)-1].E == latestHist[len(latestHist)-1].E && result.hist[len(result.hist)-1].Transaction.Z.Counter >= latestHist[len(latestHist)-1].Transaction.Z.Counter {
					latestHist = result.hist
				}
			} else {
				noreplyCount++
			}
		}
		if noreplyCount > nodeNum/2 {
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
		ackldNum := 4 - noreplyCount
		for i := 0; i < ackldNum; i++ {
			result := <-ackldResultChannel
			if result {
				log.Printf("valid ACK-LD")
			} else {
				noreplyCount++
			}
		}
		if noreplyCount > nodeNum/2 {
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
		case <-time.After(6 * time.Second):
			log.Printf("time out, restarting election...")
			vs.Stop()
			return "elect"
		}
	}
}

func ElectionMessengerRoutine(port string, voteHolder chan int32, electionHolder chan stateEpoch, voteRequestResultChannel chan stateEpoch, synchronizationHolder chan stateHist, ackeResultChannel chan stateHist, ackldResultChannel chan bool, commitldHolder chan bool) {
	// dial and askvote
	conn, err := grpc.Dial(port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("could not connect to server %s: %v", port, err)
	}
	defer conn.Close()

	voteHolder <- 0

	client := pb.NewVoterCandidateClient(conn)

	if !askvoteHelper(port, voteRequestResultChannel, client) {
		return
	}

	// wait for quorum check and send new epoch
	checkElection := <-electionHolder
	if !checkElection.state {
		return
	}
	if !newepochHelper(port, checkElection.epoch, ackeResultChannel, client) {
		return
	}

	// wait for quorum check and propose new leader
	checkSynchronization := <-synchronizationHolder
	if !checkSynchronization.state {
		return
	}
	if !newleaderHelper(port, checkSynchronization.hist, ackldResultChannel, client) {
		return
	}

	// wait for quorum check and commit new leader
	if !<-commitldHolder {
		return
	}
	commitldHelper(port, client)
}

func newleaderHelper(port string, latestHist []*pb.PropTxn, ackldResultChannel chan bool, client pb.VoterCandidateClient) bool {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// for local follower
	lastLeaderProp = lastEpoch
	pStorage = latestHist
	for _, v := range pStorage {
		v.E = lastEpoch
	}

	r, err := client.NewLeader(ctx, &pb.EpochHist{Epoch: lastEpoch, Hist: latestHist})
	if err != nil {
		// this messenger is dead
		ackldResultChannel <- false
		log.Printf("failed to ask NEWLEADER ack from %s: %v", port, err)
		return false
	}
	if !r.GetVoted() {
		ackldResultChannel <- false
		log.Printf("%s refused to give NEWLEADER ack", port)
		return false
	}
	ackldResultChannel <- true
	return true
}

func commitldHelper(port string, client pb.VoterCandidateClient) bool {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// for local follower
	for _, v := range pStorage {
		dStruct[v.Transaction.V.Key] = v.Transaction.V.Value
	}

	_, err := client.CommitNewLeader(ctx, &pb.Epoch{Epoch: lastEpoch})
	if err != nil {
		log.Printf("failed to send Commit-LD to %s: %v", port, err)
		return false
	}
	return true
}

func askvoteHelper(port string, voteRequestResultChannel chan stateEpoch, client pb.VoterCandidateClient) bool {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := client.AskVote(ctx, &pb.Epoch{Epoch: lastEpochProp})
	if err != nil {
		voteRequestResultChannel <- stateEpoch{false, -1}
		log.Printf("failed to ask vote from %s: %v", port, err)
		return false
	}
	if !r.GetVoted() {
		voteRequestResultChannel <- stateEpoch{false, -1}
		log.Printf("%s refused to vote for me", port)
		return false
	}
	voteRequestResultChannel <- stateEpoch{true, lastEpochProp}
	log.Printf("%s voted me", port)
	return true
}

func newepochHelper(port string, epoch int32, ackeResultChannel chan stateHist, client pb.VoterCandidateClient) bool {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// for local follower
	lastEpochProp = epoch

	hist, err := client.NewEpoch(ctx, &pb.Epoch{Epoch: epoch})
	if err != nil {
		// this messenger is dead
		ackeResultChannel <- stateHist{false, nil}
		log.Printf("messenger quit: failed to get ACK-E from %s: %v", port, err)
		return false
	}
	ackeResultChannel <- stateHist{true, hist.GetHist()}
	return true
}

// Register voter's grpc server exposed to potential leader
func serveV(port string, sleepHolder chan int32) {
	// listen port
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Printf("voter listening at %v", lis.Addr())
	vs = grpc.NewServer()
	pb.RegisterVoterCandidateServer(vs, &voterServer{})
	sleepHolder <- 0
	if err := vs.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

// Follower's Handler: implementation of Broadcast handler
func (s *followerServer) Broadcast(ctx context.Context, in *pb.PropTxn) (*pb.AckTxn, error) {
	if in.GetTransaction().GetZ().Epoch != lastLeaderProp {
		// simply for preventing delayed dilivery from old leader
		return &pb.AckTxn{Content: "Ignore stale messeges"}, nil
	}
	log.Printf("Follower received proposal")
	// writes the proposal to local stable storage
	pStorage = append(pStorage, in)
	log.Printf("local stable storage: %+v", pStorage)
	return &pb.AckTxn{Content: "I Acknowledged"}, nil
}

// Follower's Handler: implementation of Commit handler
func (s *followerServer) Commit(ctx context.Context, in *pb.CommitTxn) (*pb.Empty, error) {
	if in.GetEpoch() != lastLeaderProp {
		// simply for preventing delayed delivery from old leaders
		return &pb.Empty{Content: "prevent stale messeges"}, nil
	}
	log.Printf("Follower received commit")
	// writes the transaction from local stable storage to local data structure
	dStruct[pStorage[len(pStorage)-1].Transaction.V.Key] = pStorage[len(pStorage)-1].Transaction.V.Value
	log.Printf("local data structure: %+v", dStruct)
	return &pb.Empty{Content: "Commit message recieved"}, nil
}

// Follower's Handler: implementation of HeartBeat handler
func (s *followerServer) HeartBeat(ctx context.Context, in *pb.Empty) (*pb.Empty, error) {
	beatChannel <- 1
	return &pb.Empty{Content: "bump"}, nil
}

func Follow(port string) string {
	// follower receiving leader's heartbeat
	beatChannel = make(chan int, nodeNum)
	leaderStat := make(chan string, 1)
	go BeatReceiveRoutine(leaderStat)

	go serveF(port)

	<-leaderStat
	log.Printf("leader dead")
	fs.Stop()
	return "elect"
}

// Register follower's grpc server exposed to the leader
func serveF(port string) {
	// listen leader on port
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Printf("server listening at %v", lis.Addr())
	fs = grpc.NewServer()
	pb.RegisterFollowerLeaderServer(fs, &followerServer{})
	if err := fs.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

// NOTE:	now this routine do very minimal work, this does not even check if the ack is valid.
//
//	most work are handled by the PropAndCmtRoutine
func AckToCmtRoutine() {
	for {
		for i := 0; i < nodeNum; i++ {
			select {
			case <-universalCleanupHolder:
				// this is to prevent PropAndCmtRoutine stucked at waiting commits
				// if quorum is already dead, it's okay to send this commit since
				// the PropAndCmt will check this
				for i := 0; i < nodeNum; i++ {
					commitChannels[i] <- &pb.CommitTxn{Content: "Please commit"}
				}
				return
			default:
				<-ackChannel
			}
		}

		for i := 0; i < nodeNum; i++ {
			commitChannels[i] <- &pb.CommitTxn{Content: "Please commit", Epoch: lastEpoch}
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

	// leader receiving follower's heartbeat
	go BeatSendRoutine(port, serial, client)

	go PropAndCmtRoutine(port, serial, client)

	messengerHolder := make(chan int32, 1)
	<-messengerHolder
}

func PropAndCmtRoutine(port string, serial int32, client pb.FollowerLeaderClient) {
	for {

		var proposal *pb.PropTxn

		select {
		case <-universalCleanupHolder:
			return
		default:
			proposal = <-propChannels[serial]
		}

		if !upFollowers[serial] {
			ackChannel <- Ack{false, serial, -1, -1}
		} else {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			r, err := client.Broadcast(ctx, proposal)
			defer cancel()
			if err != nil {
				log.Printf("could not broadcast to server %s: %v", port, err)
				status, ok := status.FromError(err)
				if ok {
					if status.Code() == codes.DeadlineExceeded {
						log.Printf("Server %s timeout", port)
					}
				}
				upFollowersUpdateChannel <- serial
				ackChannel <- Ack{false, serial, -1, -1}
			}
			if r.GetContent() == "I Acknowledged" {
				ackChannel <- Ack{true, serial, proposal.GetTransaction().GetZ().Epoch, proposal.GetTransaction().GetZ().Counter}
			} else {
				log.Printf("Server %s failed to acknowledge: %s", port, r.GetContent())
				ackChannel <- Ack{false, serial, -1, -1}
			}
		}

		commit := <-commitChannels[serial]
		if upFollowers[serial] && upNum > nodeNum/2 {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			r, err := client.Commit(ctx, commit)
			defer cancel()
			if err != nil {
				log.Printf("could not issue commit to server %s: %v", port, err)
				status, ok := status.FromError(err)
				if ok {
					if status.Code() == codes.DeadlineExceeded {
						log.Printf("Server %s timeout", port)
					}
				}
				upFollowersUpdateChannel <- serial
			}
			if r.GetContent() == "Commit message recieved" {
				log.Printf("Commit feedback recieved from %s", port)
			} else {
				log.Printf("Server %s failed to commit: %s", port, r.GetContent())
			}
		}
	}
}

// for leader
func BeatSendRoutine(port string, serial int32, client pb.FollowerLeaderClient) {
	time.Sleep(2 * time.Second) // wait for followers to start service
	for {
		select {
		case <-universalCleanupHolder:
			return
		default:
			time.Sleep(100 * time.Millisecond)
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			_, err := client.HeartBeat(ctx, &pb.Empty{Content: "Beat"})
			defer cancel()
			if err != nil {
				log.Printf("Server %s heart beat stopped: %v", port, err)
				upFollowersUpdateChannel <- serial
				return
			}
		}
	}
}

// for follower
func BeatReceiveRoutine(leaderStat chan string) {
	for {
		select {
		case <-beatChannel:
			//log.Printf("beat")
		case <-time.After(5 * time.Second):
			leaderStat <- "dead"
			return
		}
	}
}

func upFollowersUpdateRoutine() {

	// election routine reset upFollowers
	for i := 0; i < nodeNum; i++ {
		upFollowers[i] = true
		upNum = nodeNum
	}

	for {
		// handle messenger routine's report
		serial := <-upFollowersUpdateChannel
		if serial >= 0 && upFollowers[serial] {
			// messenger routine report follower failure
			upFollowers[serial] = false
			upNum--
			// Check if quorum dead
			if upNum <= nodeNum/2 {
				log.Printf("quorum dead")
				for i := 0; i < nodeNum*2+2; i++ {
					universalCleanupHolder <- 0
				}
				return
			}
		}
	}
}
