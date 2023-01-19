package main

import (
	"context"
	"log"
	"math/rand"
	"net"
	"os"
	"strconv"
	"strings"
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

const nodeNum = 5 // <--------------- NUMBER OF SERVER NODES

const (
	EMPTY = 0
	EPOCH = 1

	PROP_TXN   = 2
	ACK_TXN    = 3
	COMMIT_TXN = 4

	VEC   = 5
	KEY   = 6
	VALUE = 7

	VOTE = 8

	NEW_EPOCH = 9
	ACK_E     = 10

	NEW_LEADER = 11
	ACK_NEW    = 12

	COMMIT_NEW_LEADER = 13
)

var (
	vs *grpc.Server
	ls *grpc.Server
	fs *grpc.Server

	// Constants
	portForUser string
	serverPorts = []string{"localhost:50051", "localhost:50052", "localhost:50053", "localhost:50054", "localhost:50055"}
	// Global channcels for broadcast-phase leader
	propChannels   [nodeNum]chan *pb.Message
	commitChannels [nodeNum]chan *pb.Message
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
	lastCount int32

	// for follower use
	lastEpochProp  int32 = 0
	lastLeaderProp int32 = 0

	pStorage []*pb.PropTxn
	dStruct  map[string]int32

	voted          chan int32
	votedBuffer    chan int32
	followerHolder chan bool

	//	for leader's use to clean up everything
	universalCleanupHolder chan int32

	// uses to pause and resume voter service handlers
	voterPaused       = true
	voterPauseChannel chan bool
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
	serial, _ := strconv.Atoi(os.Args[1])
	portForUser = (serverPorts[serial])[0:strings.Index(serverPorts[serial], ":")] + ":50056"
	voterPauseChannel = make(chan bool, 2)
	go voterPauseUpdateRoutine()
	voterPauseChannel <- true
	r := Elect(serverPorts[serial], int32(serial))
	for {
		if r == "elect" {
			r = Elect(serverPorts[serial], int32(serial))
		} else if r == "lead" {
			r = Lead(serverPorts[serial])
		} else if r == "follow" {
			r = Follow(serverPorts[serial])
		}
	}
}

// Leader's Handler: implementation of user Store handler
func (s *leaderServer) Store(ctx context.Context, in *pb.Message) (*pb.Message, error) {
	log.Printf("Leader received user request")
	for i := 0; i < nodeNum; i++ {
		propChannels[i] <- &pb.Message{Type: PROP_TXN, Epoch: lastEpoch, Transaction: &pb.Txn{V: &pb.Vec{Key: in.GetKey(), Value: in.GetValue()}, Z: &pb.Zxid{Epoch: lastEpoch, Counter: lastCount}}}
	}
	lastCount++
	return &pb.Message{Type: EMPTY, Content: "Leader recieved your request"}, nil
}

// Leader's Handler: implementation of user Retrieve handler
func (s *leaderServer) Retrieve(ctx context.Context, in *pb.Message) (*pb.Message, error) {
	log.Printf("Leader received user retrieve")
	return &pb.Message{Type: VALUE, Value: dStruct[in.GetKey()]}, nil
}

// Leader's Handler: implementation of user Identify handler
func (s *leaderServer) Identify(ctx context.Context, in *pb.Message) (*pb.Message, error) {
	log.Printf("Leader received user identify")
	return &pb.Message{Type: VOTE, Voted: true}, nil
}

func Lead(port string) string {
	lastCount = 0

	for i := 0; i < nodeNum; i++ {
		propChannels[i] = make(chan *pb.Message, 1)
		commitChannels[i] = make(chan *pb.Message, 1)
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

	fs.Stop()
	ls.Stop()
	return "elect"
}

// Register leader's grpc server exposed to user
func serveL() {
	// listen user
	lis, err := net.Listen("tcp", portForUser)
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
func (s *voterServer) AskVote(ctx context.Context, in *pb.Message) (*pb.Message, error) {
	if voterPaused {
		return nil, status.Error(codes.PermissionDenied, "voter server is paused")
	}
	log.Printf("Voter received candidate vote request")
	select {
	case <-voted:
		log.Printf("not voted yet, now vote")
		votedBuffer <- 1
		return &pb.Message{Type: VOTE, Voted: true}, nil
	default:
		log.Printf("voted, do not vote")
		return &pb.Message{Type: VOTE, Voted: false}, nil
	}
}

// Voter's Handler: implementation of NEWEPOCH handler
func (s *voterServer) NewEpoch(ctx context.Context, in *pb.Message) (*pb.Message, error) {
	if voterPaused {
		return nil, status.Error(codes.PermissionDenied, "voter server is paused")
	}
	// check new epoch
	if lastEpochProp >= in.GetEpoch() {
		return nil, status.Errorf(codes.InvalidArgument,
			"new epoch not new enough")
	}
	lastEpochProp = in.GetEpoch()
	// acknowledge new epoch proposal
	return &pb.Message{Type: ACK_E, Epoch: lastLeaderProp, Hist: pStorage}, nil
}

// Voter's Handler: implementation of NEWLEADER handler
func (s *voterServer) NewLeader(ctx context.Context, in *pb.Message) (*pb.Message, error) {
	if voterPaused {
		return nil, status.Error(codes.PermissionDenied, "voter server is paused")
	}
	// check new leader's epoch
	if in.GetEpoch() != lastEpochProp {
		followerHolder <- false
		return &pb.Message{Type: ACK_NEW, Voted: false}, nil
	}
	// update last leader and pStorage
	lastLeaderProp = in.GetEpoch()
	pStorage = in.GetHist()
	for _, v := range pStorage {
		v.E = in.GetEpoch()
	}
	// acknowledge NEWLEADER proposal
	return &pb.Message{Type: ACK_NEW, Voted: true}, nil
}

// Voter's Handler: implementation of CommitNewLeader handler
func (s *voterServer) CommitNewLeader(ctx context.Context, in *pb.Message) (*pb.Message, error) {
	if voterPaused {
		return nil, status.Error(codes.PermissionDenied, "voter server is paused")
	}
	if in.GetEpoch() != lastLeaderProp {
		return &pb.Message{Type: EMPTY}, nil // this is simply for prevent delayed delivery of messeges from previous leader
		// do not need to abort here, simply ignore it
	}
	// update dStruct
	log.Printf("updating dStruct")
	for _, v := range pStorage {
		dStruct[v.T.V.Key] = v.T.V.Value
	}
	followerHolder <- true

	return &pb.Message{Type: EMPTY}, nil
}

func Elect(port string, serial int32) string {
	voted = make(chan int32, 1)
	votedBuffer = make(chan int32, 1)
	followerHolder = make(chan bool, 1)
	voted <- 0
	voterPauseChannel <- false
	go serveV(port)
	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(10)*300 + 600 // n will be between 600 and 1500
	// time to check
	select {
	case <-time.After(time.Duration(n) * time.Millisecond):
		// didn't receive request, try to be leader
		<-voted
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
			voterPauseChannel <- true
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
				} else if result.hist[len(result.hist)-1].E > latestHist[len(latestHist)-1].E || result.hist[len(result.hist)-1].E == latestHist[len(latestHist)-1].E && result.hist[len(result.hist)-1].T.Z.Counter >= latestHist[len(latestHist)-1].T.Z.Counter {
					latestHist = result.hist
				}
			} else {
				noreplyCount++
			}
		}
		if noreplyCount > nodeNum/2 {
			log.Printf("insufficient acke, restarting election...")
			for i := 0; i < 4; i++ {
				synchronizationHolder <- stateHist{false, []*pb.PropTxn{{E: -1, T: &pb.Txn{V: &pb.Vec{Key: "", Value: -1}, Z: &pb.Zxid{Epoch: -1, Counter: -1}}}}}
			}
			voterPauseChannel <- true
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
			voterPauseChannel <- true
			vs.Stop()
			return "elect"
		} else {
			log.Printf("lifting commit holder...")
			for i := 0; i < 4; i++ {
				commitldHolder <- true
			}
		}

		log.Printf("proceeding to phase 3 as an estabilished leader...")
		voterPauseChannel <- true
		vs.Stop()
		return "lead"
	case <-votedBuffer:
		select {
		case r := <-followerHolder:
			if !r {
				log.Printf("followerHolder returned false, restarting election...")
				voterPauseChannel <- true
				vs.Stop()
				return "elect"
			}

			log.Printf("proceeding to phase 3 as a follower...")
			voterPauseChannel <- true
			vs.Stop()
			return "follow"
		case <-time.After(6 * time.Second):
			log.Printf("time out, restarting election...")
			voterPauseChannel <- true
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

	r, err := client.NewLeader(ctx, &pb.Message{Type: NEW_LEADER, Epoch: lastEpoch, Hist: latestHist})
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
		dStruct[v.T.V.Key] = v.T.V.Value
	}

	_, err := client.CommitNewLeader(ctx, &pb.Message{Type: COMMIT_NEW_LEADER, Epoch: lastEpoch})
	if err != nil {
		log.Printf("failed to send Commit-LD to %s: %v", port, err)
		return false
	}
	return true
}

func askvoteHelper(port string, voteRequestResultChannel chan stateEpoch, client pb.VoterCandidateClient) bool {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := client.AskVote(ctx, &pb.Message{Type: EPOCH, Epoch: lastEpochProp})
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

	hist, err := client.NewEpoch(ctx, &pb.Message{Type: NEW_EPOCH, Epoch: epoch})
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
func serveV(port string) {
	// listen port
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Printf("voter server start listening at %v", lis.Addr())
	vs = grpc.NewServer()
	pb.RegisterVoterCandidateServer(vs, &voterServer{})
	if err := vs.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

// Follower's Handler: implementation of Broadcast handler
func (s *followerServer) Broadcast(ctx context.Context, in *pb.Message) (*pb.Message, error) {
	if in.GetTransaction().GetZ().Epoch != lastLeaderProp {
		// simply for preventing delayed dilivery from old leader
		return &pb.Message{Type: ACK_TXN, Content: "Ignore stale messeges"}, nil
	}
	log.Printf("Follower received proposal")
	// writes the proposal to local stable storage
	pStorage = append(pStorage, &pb.PropTxn{E: in.GetEpoch(), T: in.GetTransaction()})
	log.Printf("local stable storage: %+v", pStorage)
	return &pb.Message{Type: ACK_TXN, Content: "I Acknowledged"}, nil
}

// Follower's Handler: implementation of Commit handler
func (s *followerServer) Commit(ctx context.Context, in *pb.Message) (*pb.Message, error) {
	if in.GetEpoch() != lastLeaderProp {
		// simply for preventing delayed delivery from old leaders
		return &pb.Message{Type: EMPTY}, nil
	}
	log.Printf("Follower received commit")
	// writes the transaction from local stable storage to local data structure
	dStruct[pStorage[len(pStorage)-1].T.V.Key] = pStorage[len(pStorage)-1].T.V.Value
	log.Printf("local data structure: %+v", dStruct)
	return &pb.Message{Type: EMPTY, Content: "Commit message recieved"}, nil
}

// Follower's Handler: implementation of HeartBeat handler
func (s *followerServer) HeartBeat(ctx context.Context, in *pb.Message) (*pb.Message, error) {
	beatChannel <- 1
	return &pb.Message{Type: EMPTY}, nil
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
					commitChannels[i] <- &pb.Message{Type: COMMIT_TXN}
				}
				return
			default:
				<-ackChannel
			}
		}

		for i := 0; i < nodeNum; i++ {
			commitChannels[i] <- &pb.Message{Type: COMMIT_TXN, Epoch: lastEpoch}
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

		var proposal *pb.Message

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
			_, err := client.HeartBeat(ctx, &pb.Message{Type: EMPTY})
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

func voterPauseUpdateRoutine() {
	for {
		newState := <-voterPauseChannel
		voterPaused = newState
	}
}
