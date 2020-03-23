// package main

// import (
// 	"context"
// 	"log"
// 	"time"

// 	"github.com/golang/protobuf/ptypes/empty"
// 	"github.com/rehandwi03/guser/model"
// 	"google.golang.org/grpc"
// )

// func main() {
// 	conn, err := grpc.Dial("localhost:7777", grpc.WithInsecure())
// 	if err != nil {
// 		log.Fatalf("did not connect: %v", err)
// 	}
// 	defer conn.Close()

// 	c := model.NewUserServiceClient(conn)

// 	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 	defer cancel()

// 	res, err := c.GetUsers(ctx, new(empty.Empty))
// 	if err != nil {
// 		log.Fatalf("GetUser failed: %v", err)
// 	}
// 	log.Printf("GetUser result: <%+v>\n\n", res)

// 	req1 := model.CreateRequest{
// 		User.User.User: "manjiw",
// 		Password: "manjiw",
// 	}
// 	res1, err := c.Create(ctx, &req1)
// 	if err != nil {
// 		log.Fatalf("Create failed: %v", err)
// 	}
// 	log.Printf("Create result: <%+v>\n\n", res1)

// 	id := res1.User.Id

// 	req2 := model.ReadRequest{
// 		Id: id,
// 	}
// 	res2, err := c.Read(ctx, &req2)
// 	if err != nil {
// 		log.Fatalf("Read failed: %v", err)
// 	}
// 	log.Printf("Read result: %v", res2)

// }
