package main

import (
	"context"
	"fmt"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func main() {
	// setup db
	db, err := setupDB()
	if err != nil {
		log.Fatalf("error on setting up db. err=%v", err)
	}
	defer db.Client().Disconnect(context.Background())

	// setup rabbit
	rabbit, err := setupRabbit()
	if err != nil {
		log.Fatalf("error on setting up rabbit. err=%v", err)
	}
	defer rabbit.Close()

	// setup saga
	// ==========

	// setup sample activities
	activity1 := NewActivity("activity-1", samplecommit, samplerollback)
	activity2 := NewActivity("activity-2", samplecommit, samplerollback)
	activity3 := NewActivity("activity-3", samplecommit, samplerollback)

	// setup sample workflow
	workflow1 := NewWorkflow[sampleparameter](
		"workflow-1",
		activity1, activity2, activity3,
	)

	// setup saga client
	saga := New(
		db,
		rabbit,
		workflow1,
	)

	saga.Start()

	for range 5 {
		workflow1.Execute(sampleparameter{})
	}
}

type sampleparameter struct {
	Field1 string `bson:"field_1"`
	Field2 string `bson:"field_2"`
	Field3 string `bson:"field_3"`
}

func samplecommit() {
	log.Printf("this is a commit :D")
}
func samplerollback() {
	log.Printf("this is a rollback :D")
}

// setupRabbit setups rabbit connection for saga purposes
func setupRabbit() (*amqp.Connection, error) {
	connection, err := amqp.Dial("amqp://guest:guest@host.docker.internal:5672/")
	if err != nil {
		return nil, fmt.Errorf("error on dialing to rabbit. err=%w", err)
	}

	return connection, nil
}

// setupDB setups db connection to mongo db for saga purposes
func setupDB() (*mongo.Database, error) {
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
