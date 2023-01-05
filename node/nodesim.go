package main

import (
	pb "gozab/gozab"
)

type State struct {
	pStorage []*pb.PropTxn
	dStruct  map[string]int32

	lastEpochProp  int32
	lastLeaderProp int32
}

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
