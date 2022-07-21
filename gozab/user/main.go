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

	//ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	//defer cancel()

	// stdin key value
	for {
		var key string
		var val int32
		safe := 0
		fmt.Printf("Enter key-value pair: ")
		fmt.Scanf("%s %d", &key, &val)
		safe = 1
		if safe == 1 {
			r, err := c.SendRequest(context.Background(), &pb.Vec{Key: key, Value: val})
			if err != nil {
				log.Fatalf("could not send request to leader: %v", err)
			}
			log.Printf("Greeting: %s", r.GetContent())
			safe = 0
		}
	}
}
