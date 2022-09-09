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

var (
	serverPorts = []string{"198.22.255.12:50056", "198.22.255.16:50056", "198.22.255.28:50056", "198.22.255.11:50056", "198.22.255.15:50056"}
)

func main() {
	found := make(chan int, 1)
	lost := make(chan int, 1)
	for {
		go finder(1, found, lost)
		go finder(2, found, lost)
		go finder(3, found, lost)
		go finder(4, found, lost)
		go finder(0, found, lost)
		select {
		case <-found:
			<-lost
		case <-time.After(2 * time.Second):
			continue
		}
	}
}

func finder(serial int, found chan int, lost chan int) {
	conn, err := grpc.Dial(serverPorts[serial], grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("failed to connect: %v", err)
		return
	}
	defer conn.Close()

	c := pb.NewLeaderUserClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.Identify(ctx, &pb.Empty{Content: "are you the leader?"})
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
				r, err := c.Store(context.Background(), &pb.Vec{Key: key, Value: val})
				if err != nil {
					log.Printf("could not send request to leader: %v", err)
					lost <- 1
					return
				}
				log.Printf("Greeting: %s", r.GetContent())
			} else if command == "Get" {
				fmt.Printf("Enter a key: ")
				fmt.Scanf("%s", &key)
				r, err := c.Retrieve(context.Background(), &pb.GetTxn{Key: key})
				if err != nil {
					log.Printf("could not send Get request to leader: %v", err)
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
