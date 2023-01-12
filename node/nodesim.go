package main

import (
	pb "gozab/gozab"
)

type State struct {
	// follower states
	pStorage []*pb.PropTxn
	dStruct  map[string]int32

	lastEpochProp  int32
	lastLeaderProp int32

	// leader states
	lastEpoch int32
	ackcnt    int32
	proposals []*pb.PropTxn
}

// leader functions
func (s *State) StoreTransformation(message *pb.Vec) {
	s.proposals = append(s.proposals, &pb.PropTxn{E: s.lastEpoch, Transaction: &pb.Txn{V: message, Z: &pb.Zxid{Epoch: s.lastEpoch, Counter: -1}}})
}

func (s *State) AckTransformation(message *pb.AckTxn) {
	if message.Content == "I Acknowledged" {
		s.ackcnt++
	}
}

// follower functions
func (s *State) BroadcastTransformation(message *pb.PropTxn) {
	if s.lastLeaderProp == message.E {
		s.pStorage = append(s.pStorage, message)
	}
}

func (s *State) CommitTransformation(message *pb.CommitTxn) {
	if s.lastLeaderProp == message.Epoch {
		s.dStruct[s.pStorage[len(s.pStorage)-1].Transaction.V.Key] = s.pStorage[len(s.pStorage)-1].Transaction.V.Value
	}
}

// candidate functions
func (s *State) NewEpochTransformation(message *pb.Epoch) {
	if s.lastEpochProp < message.Epoch {
		s.lastEpochProp = message.Epoch
	}
}

func (s *State) NewLeaderTransformation(message *pb.EpochHist) {
	if s.lastEpochProp == message.Epoch {
		s.lastLeaderProp = message.Epoch
		s.pStorage = message.Hist
		for _, v := range s.pStorage {
			v.E = message.Epoch
		}
	}
}

func (s *State) CommitNewLeaderTransformation(message *pb.Epoch) {
	if s.lastLeaderProp == message.Epoch {
		for _, v := range s.pStorage {
			s.dStruct[v.Transaction.V.Key] = v.Transaction.V.Value
		}
	}
}

type LeaderPartialState struct {
	pStorage []*pb.PropTxn
	dStruct  map[string]int32

	epoch   int32
	counter int32
}

func BroadcastInference(message *pb.PropTxn) LeaderPartialState {
	return LeaderPartialState{pStorage: nil, dStruct: nil, epoch: message.Transaction.Z.Epoch, counter: message.Transaction.Z.Counter}
}

func CommitInference(message *pb.CommitTxn) LeaderPartialState {
	return LeaderPartialState{pStorage: nil, dStruct: nil, epoch: message.Epoch, counter: -1}
}

type CandidatePartialState struct {
}
