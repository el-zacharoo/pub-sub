package store

import (
	"context"
	"log"

	"google.golang.org/grpc/metadata"

	pb "github.com/el-zacharoo/publisher/gen/proto/go/person/v1"
)

type Storer interface {
	CreatePerson(msg *pb.Person, md metadata.MD) error
}

func (s Store) CreatePerson(msg *pb.Person, md metadata.MD) error {
	_, err := s.locaColl.InsertOne(context.Background(), msg)
	if err != nil {
		log.Fatal(err)
	}
	return err
}
