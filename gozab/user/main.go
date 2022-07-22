package main

import (
	"context"
	"fmt"
	"log"

	pb "gozab/gozab"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.Dial("localhost:50056", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect to leader: %v", err)
	}
	defer conn.Close()
	c := pb.NewClientConsoleClient(conn)

	// stdin key value
	for {
		var key string
		var val int32
		var command string

		fmt.Printf("Enter 'Send' or 'Get' (^C to quit): ")
		fmt.Scanf("%s", &command)

		if command == "Send" {
			fmt.Printf("Enter key-value pair: ")
			fmt.Scanf("%s %d", &key, &val)
			r, err := c.SendRequest(context.Background(), &pb.Vec{Key: key, Value: val})
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
