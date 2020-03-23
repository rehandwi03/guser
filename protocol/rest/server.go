package rest

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/handlers"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/rehandwi03/guser/model"
	"google.golang.org/grpc"
)

func UserRunServer(ctx context.Context, grpcPort, httpPort string) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	if err := model.RegisterUserServiceHandlerFromEndpoint(ctx, mux, "localhost:"+grpcPort, opts); err != nil {
		log.Fatalf("failed to start HTTP gateway: %v", err)
	}
	// newMux := handlers.CORS(
	// 	handlers.AllowedOrigins([]string{"GET", "POST", "PUT", "DELETE", "OPTION"}),
	// 	handlers.AllowedOrigins([]string{"http://localhost:7888"}),
	// 	handlers.AllowedOrigins([]string{"Content-Type"}),
	// )(mux)
	srv := &http.Server{
		Addr:    ":" + httpPort,
		Handler: handlers.CORS()(mux),
	}
	// graceful shutdown
	c := make(chan os.Signal, 1)

	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {

		}
		_, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		_ = srv.Shutdown(ctx)
	}()

	log.Println("starting User HTTP/REST gateway on port ", httpPort)
	return srv.ListenAndServe()
}
