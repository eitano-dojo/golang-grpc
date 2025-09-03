package main

import (
	"context"
	"log"
	"time"

	pb "github.com/eitano-dojo/golang-grpc/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewGreeterClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.SayHello(ctx, &pb.HelloRequest{Name: "eitan"})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.GetMessage())

	s, err := c.SayGoodbye(ctx, &pb.ByeRequest{Name: "eitan"})
	if err != nil {
		log.Fatalf("could not depart: %v", err)
	}
	log.Printf("Departure: %s", s.GetMessage())

}
