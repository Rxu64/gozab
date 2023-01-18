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

func BroadcastInference(message *pb.Message) LeaderPartialState {
	return LeaderPartialState{pStorage: nil, dStruct: nil, epoch: message.Epoch, counter: message.Counter}
}

func CommitInference(message *pb.Message) LeaderPartialState {
	return LeaderPartialState{pStorage: nil, dStruct: nil, epoch: message.Epoch, counter: -1}
}

type CandidatePartialState struct {
}
