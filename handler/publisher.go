package handler

import (
	"context"

	dapr "github.com/dapr/go-sdk/client"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	pb "github.com/el-zacharoo/publisher/gen/proto/go/person/v1"
	"github.com/el-zacharoo/publisher/store"
)

type Server struct {
	Dapr  dapr.Client
	Store store.Storer
	pb.UnimplementedPersonServiceServer
}

func (s Server) Create(ctx context.Context, req *pb.CreateRequest) (*pb.CreateResponse, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return &pb.CreateResponse{}, status.Errorf(codes.Aborted, "%s", "no incoming context")
	}

	person := req.Person
	person.Id = uuid.NewString()
	message := "Submission for " + person.GetName() + " posted successfully."

	if err := s.Dapr.PublishEvent(
		context.Background(),
		"pubsubsrv", "zacharysPubSub", person,
		dapr.PublishEventWithContentType("application/json"),
	); err != nil {
		return &pb.CreateResponse{}, status.Errorf(codes.Aborted, "%s", "error publishing event")
	}

	if err := s.Store.CreatePerson(person, md); err != nil {
		return &pb.CreateResponse{}, status.Errorf(codes.Aborted, "%v", err)
	}
	return &pb.CreateResponse{Message: message, Person: person}, nil
}
