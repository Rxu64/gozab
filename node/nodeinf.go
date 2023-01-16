package main

import (
	pb "gozab/gozab"
)

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
