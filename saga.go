package main

import (
	"context"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type IProcess interface {
	Start(s *saga)
	Stop()
}

// saga is the client to orchestrate workflows
type saga struct {
	db     *mongo.Database
	rabbit *amqp.Connection

	workflows []IProcess
}

func New(db *mongo.Database, rabbit *amqp.Connection) *saga {
	return &saga{
		db:     db,
		rabbit: rabbit,
	}
}

func (s *saga) Start() {
	log.Printf("starting saga...")

	for _, workflow := range s.workflows {
		workflow.Start(s)
	}
}
func (s *saga) Stop() {
	log.Printf("stopping saga...")
}

func (s *saga) RegisterWorkflows(workflows ...IProcess) {
	s.workflows = workflows
}

type workflowEntry struct {
	Name      string      `bson:"name"`
	CreatedAt time.Time   `json:"created_at"`
	Payload   interface{} `bson:"payload"`
}

func (s *saga) CreateEntry(workflowName string, payload interface{}) {
	res, err := s.db.Collection("workflows").InsertOne(
		context.Background(),
		workflowEntry{
			Name:      workflowName,
			CreatedAt: time.Now(),
			Payload:   payload,
		},
	)

	if err != nil {
		log.Fatalf("error on creating workflow entry to db. err=%v", err)
	}

	log.Printf("created entry for '%s' workflow. id=%+v", workflowName, res.InsertedID.(bson.ObjectID).Hex())
}
