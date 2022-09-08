package main

import (
	"context"
	"fmt"
	"log"

	pb "gozab/gozab"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	serverPorts = []string{"198.22.255.12:50051", "198.22.255.16:50052", "198.22.255.28:50053", "198.22.255.11:50054", "198.22.255.15:50055"}
	serverMap   = map[string]int32{"198.22.255.12:50051": 0, "198.22.255.16:50052": 1, "198.22.255.28:50053": 2, "198.22.255.11:50054": 3, "198.22.255.15:50055": 4}
)

func main() {
	conn, err := grpc.Dial("localhost:50056", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect to leader: %v", err)
	}
	defer conn.Close()
	c := pb.NewLeaderUserClient(conn)

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
				log.Fatalf("could not send request to leader: %v", err)
			}
			log.Printf("Greeting: %s", r.GetContent())
		} else if command == "Get" {
			fmt.Printf("Enter a key: ")
			fmt.Scanf("%s", &key)
			r, err := c.Retrieve(context.Background(), &pb.GetTxn{Key: key})
			if err != nil {
				log.Fatalf("could not send Get request to leader: %v", err)
			}
			log.Printf("Greeting: %d", r.GetValue())
		} else {
			continue
		}
	}
}
