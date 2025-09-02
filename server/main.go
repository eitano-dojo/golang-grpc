package main

import (
	"context"
	"log"
	"net"
	"time"

	pb "github.com/eitano-dojo/golang-grpc/proto"
	"github.com/gofiber/fiber/v2"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedGreeterServer
}

func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Message: "Hello " + in.Name}, nil
}

func main() {
	go startGRPCServer()

	app := fiber.New()

	app.Get("/hello/:name", func(c *fiber.Ctx) error {
		name := c.Params("name")
		conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure(), grpc.WithBlock())
		if err != nil {
			return err
		}
		defer conn.Close()
		client := pb.NewGreeterClient(conn)

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		resp, err := client.SayHello(ctx, &pb.HelloRequest{Name: name})
		if err != nil {
			return err
		}

		return c.SendString(resp.GetMessage())
	})

	log.Fatal(app.Listen(":3000"))
}

func startGRPCServer() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
