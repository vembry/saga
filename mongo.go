package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// setupMongo setups db connection to mongo db for saga purposes
func setupMongo() (*mongo.Database, error) {
	// saga uses mongodb and expects:
	// - "workflows" collection
	// ....

	// establish db connection
	client, err := mongo.Connect(options.Client().ApplyURI("mongodb://root:root@host.docker.internal:27017"))
	if err != nil {
		return nil, fmt.Errorf("error on connecting to db. err=%w", err)
	}

	// select db
	db := client.Database("saga")

	// test insert
	res, err := db.Collection("tests").InsertOne(context.Background(), map[string]interface{}{
		"name": "setup",
		"time": time.Now().String(),
	})
	if err != nil {
		return nil, fmt.Errorf("error on inserting entry into collection. err=%v", err)
	}

	// log out insertions output
	val := res.InsertedID.(bson.ObjectID)
	log.Printf("test inserted=%s", val.Hex())

	return db, nil
}
