package main

import (
	"log"
	"net"
	"time"

	daprpb "github.com/dapr/dapr/pkg/proto/runtime/v1"
	dapr "github.com/dapr/go-sdk/client"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	pb "github.com/el-zacharoo/pubsub/gen/proto/go/person/v1"
	"github.com/el-zacharoo/pubsub/handler"
	"github.com/el-zacharoo/pubsub/store"
)

func main() {
	time.Sleep(2 * time.Second)
	client, err := dapr.NewClient()
	if err != nil {
		log.Fatalf("failed to initialise Dapr client: %v", err)
	}
	defer client.Close()

	grpcSrv := grpc.NewServer()
	defer grpcSrv.Stop()         // stop server on exit
	reflection.Register(grpcSrv) // for postman

	h := &handler.Server{
		Store: store.Connect(),
		Dapr:  client,
	}
	pb.RegisterPersonServiceServer(grpcSrv, h)

	ch := handler.CallbackServer{}
	daprpb.RegisterAppCallbackServer(grpcSrv, ch)

	lis, err := net.Listen("tcp", ":8081")
	if err != nil {
		log.Fatalln("Failed to listen:", err)
	}
	if err := grpcSrv.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
