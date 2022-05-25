package store

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Store struct {
	locaColl *mongo.Collection
}

func Connect() Store {
	clientOptions := options.Client().ApplyURI("mongodb+srv://geoloaction:e2Fyk5w2ZJnV6uzN@cluster0.u4qeu.mongodb.net")
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	db := client.Database("info")

	return Store{
		locaColl: db.Collection("subscriber"),
	}
}
