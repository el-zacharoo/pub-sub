package handler

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	pb "github.com/el-zacharoo/publisher/gen/proto/go/person/v1"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/encoding/protojson"
)

type Server struct {
	pb.UnimplementedPersonServiceServer
}

func (s Server) Person(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodOptions {
		return
	}
	defer r.Body.Close()

	byt, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("err %v", err)))
	}

	req := &pb.Request{}
	err = protojson.Unmarshal(byt, req)

	rsp, err := s.SendProto(req)

	if err != nil {
		w.Write([]byte(fmt.Sprintf("err %v", err)))
	} else {
		w.Write([]byte(fmt.Sprintf("%s", rsp.GetMessage())))
	}

}

func (s Server) SendProto(req *pb.Request) (*pb.Response, error) {
	conn, err := grpc.Dial(":8082", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewPersonServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	sub := req.Person
	sub.Id = uuid.NewString()

	per, err := c.Create(ctx, &pb.Request{Person: sub})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Subscribing service response: %s", per.GetMessage())

	return per, nil
}
