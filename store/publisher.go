package store

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"google.golang.org/grpc/metadata"

	pb "github.com/el-zacharoo/pubsub/gen/proto/go/person/v1"
)

type Storer interface {
	CreatePerson(msg *pb.Person, md metadata.MD) error
	UpdatePerson(id string, md metadata.MD, u *pb.Person) error
}

func (s Store) CreatePerson(msg *pb.Person, md metadata.MD) error {
	_, err := s.locaColl.InsertOne(context.Background(), msg)
	if err != nil {
		log.Fatal(err)
	}
	return err
}

func (s Store) UpdatePerson(id string, md metadata.MD, u *pb.Person) error {
	_, err := s.locaColl.ReplaceOne(context.Background(), bson.M{"id": id}, u)
	if err != nil {
		log.Fatal(err)
	}
	return err
}
