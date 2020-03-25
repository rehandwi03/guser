package cmd

import (
	"context"
	"database/sql"
	"log"

	"github.com/rehandwi03/guser/protocol/grpc"
	"github.com/rehandwi03/guser/protocol/rest"
	"github.com/rehandwi03/guser/service"
)

const (
	grpcPortUser     = "7777"
	grpcPortKaryawan = "7776"
	httpPortUser     = "7778"
	httpPortKaryawan = "7779"
)

func RunServer() error {
	ctx := context.Background()
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/grpcuser")
	if err != nil {
		log.Fatalf("could not connect to db: %v", err)
	}
	defer db.Close()
	UserAPI := service.NewUserServiceServer(db)
	go func() {
		_ = rest.UserRunServer(ctx, grpcPortUser, httpPortUser)
	}()
	return grpc.UserRunServer(ctx, UserAPI, grpcPortUser)
}

// func RunServer() error {
// 	ctx := context.Background()
// 	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/grpcuser")
// 	if err != nil {
// 		log.Fatalf("could not connect to db: %v", err)
// 	}
// 	defer db.Close()
// 	API := service.NewUserServiceServer(db)
// 	go func() {
// 		_ = rest.RunServer(ctx, grpcPort, httpPort)
// 	}()
// 	return grpc.RunServer(ctx, API, grpcPort)
// }
