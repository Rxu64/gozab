package main

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "gozab/gozab"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

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
	serverPorts = []string{"localhost:50056", "localhost:50056", "localhost:50056", "localhost:50056", "localhost:50056"}
)

func main() {
	found := make(chan int, 1)
	lost := make(chan int, 1)
	for {
		go connectionRoutine(1, found, lost)
		go connectionRoutine(2, found, lost)
		go connectionRoutine(3, found, lost)
		go connectionRoutine(4, found, lost)
		go connectionRoutine(0, found, lost)
		select {
		case <-found:
			<-lost
		case <-time.After(2 * time.Second):
			continue
		}
	}
}

func connectionRoutine(serial int, found chan int, lost chan int) {
	conn, err := grpc.Dial(serverPorts[serial], grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("failed to connect: %v", err)
		return
	}
	defer conn.Close()

	c := pb.NewLeaderUserClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.Identify(ctx, &pb.Message{Type: EMPTY})
	if err != nil {
		log.Printf("%s is not the leader currently: %v", serverPorts[serial], err)
		return
	}
	if r.Voted { // connected to leader
		found <- 1
		var key string
		var val int32
		var command string
		// stdin key value
		for {
			fmt.Printf("Enter 'Send' or 'Get' (^C to quit): ")
			fmt.Scanf("%s", &command)

			if command == "Send" {
				fmt.Printf("Enter key-value pair: ")
				fmt.Scanf("%s %d", &key, &val)
				r, err := c.Store(context.Background(), &pb.Message{Type: VEC, Key: key, Value: val})
				if err != nil {
					log.Printf("could not send request to leader, try again: %v", err)
					lost <- 1
					return
				}
				log.Printf("Greeting: %s", r.GetContent())
			} else if command == "Get" {
				fmt.Printf("Enter a key: ")
				fmt.Scanf("%s", &key)
				r, err := c.Retrieve(context.Background(), &pb.Message{Type: KEY, Key: key})
				if err != nil {
					log.Printf("could not send Get request to leader, try again: %v", err)
					lost <- 1
					return
				}
				log.Printf("Greeting: %d", r.GetValue())
			} else {
				continue
			}
		}
	}
}
