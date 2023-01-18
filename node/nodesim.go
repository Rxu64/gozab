package main

import (
	pb "gozab/gozab"
)

type State struct {
	// TODO: 0 for candidate; 1 for leader; 2 for follower
	// currRole int32

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

func Apply(s State, msg *pb.Message) {

}

// LEADER receive from client
func (s *State) StoreTransformation(msg *pb.Vec) {
	s.proposals = append(s.proposals, &pb.PropTxn{E: s.lastEpoch, T: &pb.Txn{V: msg, Z: &pb.Zxid{Epoch: s.lastEpoch, Counter: -1}}})
}

// LEADER receive from follower
func (s *State) AckTransformation(msg *pb.Message) {
	if msg.Content == "I Acknowledged" {
		s.ackcnt++
	}
}

// FOLLOWER receive from leader
func (s *State) BroadcastTransformation(msg *pb.PropTxn) {
	if s.lastLeaderProp == msg.E {
		s.pStorage = append(s.pStorage, msg)
	}
}

func (s *State) CommitTransformation(msg *pb.Message) {
	if s.lastLeaderProp == msg.Epoch {
		s.dStruct[s.pStorage[len(s.pStorage)-1].T.V.Key] = s.pStorage[len(s.pStorage)-1].T.V.Value
	}
}

// CANDIDATE receive from candidate
func (s *State) NewEpochTransformation(msg *pb.Message) {
	if s.lastEpochProp < msg.Epoch {
		s.lastEpochProp = msg.Epoch
	}
}

func (s *State) NewLeaderTransformation(msg *pb.Message) {
	if s.lastEpochProp == msg.Epoch {
		s.lastLeaderProp = msg.Epoch
		s.pStorage = msg.Hist
		for _, v := range s.pStorage {
			v.E = msg.Epoch
		}
	}
}

func (s *State) CommitNewLeaderTransformation(msg *pb.Message) {
	if s.lastLeaderProp == msg.Epoch {
		for _, v := range s.pStorage {
			s.dStruct[v.T.V.Key] = v.T.V.Value
		}
	}
}
