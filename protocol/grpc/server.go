package grpc

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"

	"github.com/rehandwi03/guser/model"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func UserRunServer(ctx context.Context, service model.UserServiceServer, port string) error {
	listen, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}
	server := grpc.NewServer()
	model.RegisterUserServiceServer(server, service)
	reflection.Register(server)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			log.Println("Shutting down gRPC Server")
			server.GracefulStop()
			<-ctx.Done()
		}
	}()
	log.Println("starting gRPC User Server at...", port)
	return server.Serve(listen)
}
