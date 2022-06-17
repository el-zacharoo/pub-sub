package handler

import (
	"context"

	dapr "github.com/dapr/go-sdk/client"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	pb "github.com/el-zacharoo/pubsub/gen/proto/go/person/v1"
	"github.com/el-zacharoo/pubsub/store"
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
		"pubsubsrv", "create", person,
		dapr.PublishEventWithContentType("application/json"),
	); err != nil {
		return &pb.CreateResponse{}, status.Errorf(codes.Aborted, "%s", "error publishing event")
	}

	if err := s.Store.CreatePerson(person, md); err != nil {
		return &pb.CreateResponse{}, status.Errorf(codes.Aborted, "%v", err)
	}
	return &pb.CreateResponse{Message: message, Person: person}, nil
}

func (s Server) Update(ctx context.Context, req *pb.UpdateRequest) (*pb.UpdateResponse, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return &pb.UpdateResponse{}, status.Errorf(codes.Aborted, "%s", "no incoming context")
	}

	person := req.Person
	id := person.Id
	message := "Update Submission for " + person.Name + " successful"

	// publish event
	if err := s.Dapr.PublishEvent(
		context.Background(),
		"pubsubsrv", "create", person,
		dapr.PublishEventWithContentType("application/json"),
	); err != nil {
		return &pb.UpdateResponse{}, status.Errorf(codes.Aborted, "%s", "error publishing event")
	}

	if err := s.Store.UpdatePerson(id, md, person); err != nil {
		return &pb.UpdateResponse{}, status.Errorf(codes.Aborted, "%v", err)
	}

	return &pb.UpdateResponse{Person: person, Message: message}, nil
}
