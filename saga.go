package main

import (
	"context"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type IWorkflow interface {
	Start(s *saga)
	Stop()
}

// saga is the client to orchestrate workflows
type saga struct {
	db     *mongo.Database
	rabbit *amqp.Connection

	workflows []IWorkflow
}

func New(
	db *mongo.Database,
	rabbit *amqp.Connection,
	workflows ...IWorkflow,
) *saga {
	return &saga{
		db:        db,
		rabbit:    rabbit,
		workflows: workflows,
	}
}

func (s *saga) Start() {
	log.Printf("starting saga...")

	go func() {

	}()

	for _, workflow := range s.workflows {
		workflow.Start(s)
	}
}
func (s *saga) Stop() {
	log.Printf("stopping saga...")
}

type workflowStatus string

var (
	workflowStatusPending    workflowStatus = "pending"
	workflowStatusInProgress workflowStatus = "in_progress"
	workflowStatusCompleted  workflowStatus = "pending"
)

type workflowEntry struct {
	Name       string         `bson:"name"`
	Status     workflowStatus `bson:"status"`
	activities []string       `bson:"activities"`
	CreatedAt  time.Time      `bson:"created_at"`
	UpdatedAt  time.Time      `bson:"updated_at"`
	Payload    interface{}    `bson:"payload"`
}

func (s *saga) CreateEntry(workflowName string, payload interface{}) {
	res, err := s.db.Collection("workflows").InsertOne(
		context.Background(),
		workflowEntry{
			Name:       workflowName,
			Status:     workflowStatusPending,
			activities: []string{},
			CreatedAt:  time.Now(),
			Payload:    payload,
		},
	)

	if err != nil {
		log.Fatalf("error on creating workflow entry to db. err=%v", err)
	}

	log.Printf("created entry for '%s' workflow. id=%+v", workflowName, res.InsertedID.(bson.ObjectID).Hex())
}

func (s *saga) fetchReadyWorkflows() {
	var out []workflowEntry
	cursor, err := s.db.Collection("workflows").Find(context.Background(), bson.D{{Key: "status", Value: workflowStatusPending}})
	if err != nil {

	}
	cursor.Decode(&out)
}
