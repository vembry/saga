package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func main() {
	// setup db
	mongo, err := setupMongo()
	if err != nil {
		log.Fatalf("error on setting up db. err=%v", err)
	}
	defer mongo.Client().Disconnect(context.Background())

	// setup rabbit
	rabbit, err := setupRabbit()
	if err != nil {
		log.Fatalf("error on setting up rabbit. err=%v", err)
	}
	defer rabbit.Close()

	workflow1 := NewWorkflow("some-workflow", mongo, rabbit)

	workflow1.Start()

	workflow1.Execute(SomeParameter{})

	watchForExitSignal()

	log.Printf("shutting down")
}

// WatchForExitSignal is to awaits incoming interrupt signal
// sent to the service
func watchForExitSignal() os.Signal {
	log.Printf("awaiting sigterm...")
	ch := make(chan os.Signal, 4)
	signal.Notify(
		ch,
		syscall.SIGINT,
		syscall.SIGQUIT,
		syscall.SIGTERM,
		syscall.SIGTSTP,
	)

	return <-ch
}

type SomeParameter struct {
	Field1 string
	Field2 string
	Field3 string
}

type activity struct {
	id string

	commit   func()
	rollback func()
}

func NewActivity(activityId string, commit func(), rollback func()) *activity {
	return &activity{
		id:       activityId,
		commit:   commit,
		rollback: rollback,
	}
}

func (a *activity) Commit() {
	a.commit()
}

func (a *activity) Rollback() {
	a.commit()
}

type workflow struct {
	id     string
	db     *mongo.Database
	worker *amqp.Connection

	activities  []*activity
	activityMap map[string]*activity
}

func NewWorkflow(
	workflowId string,
	db *mongo.Database,
	worker *amqp.Connection,
	activities ...*activity,
) *workflow {
	activityMap := make(map[string]*activity)
	for _, activity := range activities {
		activityMap[activity.id] = activity
	}

	return &workflow{
		id:          workflowId,
		db:          db,
		worker:      worker,
		activities:  activities,
		activityMap: activityMap,
	}
}

func (w *workflow) Start() {
	w.startConsumer()
}

type workflowStatus string

var (
	WorkflowStatusNil        workflowStatus = ""
	WorkflowStatusPending    workflowStatus = "pending"
	WorkflowStatusInProgress workflowStatus = "in_progress"
	WorkflowStatusCompleted  workflowStatus = "completed"
)

type workflowEntry struct {
	Id                  bson.ObjectID  `bson:"_id"`
	WorkflowId          string         `bson:"code"`
	availableActivities []string       `bson:"available_activities"`
	processedActivities []string       `bson:"processed_activities"`
	Status              workflowStatus `bson:"status"`
	State               interface{}    `bson:"state"`
	CreatedAt           time.Time      `bson:"created_at"`
	UpdatedAt           time.Time      `bson:"updated_at"`
}

func (w *workflow) Execute(args interface{}) {
	log.Printf("executing '%s' workflow. args=%+v", w.id, args)

	activities := make([]string, len(w.activities))
	for i, activity := range w.activities {
		activities[i] = activity.id
	}

	// create state
	res, err := w.db.Collection("workflows").InsertOne(context.Background(), workflowEntry{
		Id:                  bson.NewObjectID(),
		WorkflowId:          w.id,
		availableActivities: activities,
		processedActivities: make([]string, 0),
		Status:              WorkflowStatusPending,
		State:               args,
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
	})
	if err != nil {
		log.Fatalf("error on inserting one workflow entry. err=%v", err)
	}

	// establish channel
	ch, err := w.worker.Channel()
	if ch != nil {
		defer ch.Close()
	}
	if err != nil {
		log.Fatalf("error on opening new channel to rabbit. err=%v", err)
	}

	// construct worker's payload
	payload := res.InsertedID.(bson.ObjectID).String()

	// publish message
	err = ch.Publish("", w.constructWorkerKey(), false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte(payload),
	})
	if err != nil {
		log.Fatalf("error on publish message to rabbit. err=%v", err)
	}
	log.Printf("finished executing '%s' workflow", w.id)
}

func (w *workflow) constructWorkerKey() string {
	return fmt.Sprintf("saga.%s", w.id)
}

func (w *workflow) startConsumer() {
	// establish channel
	ch, err := w.worker.Channel()
	if ch != nil {
		defer ch.Close()
	}
	if err != nil {
		log.Fatalf("error on opening new channel to rabbit. err=%v", err)
	}

	// setup consumer to queue
	messageCh, err := ch.Consume(
		w.constructWorkerKey(), // queue
		"",                     // consumer
		false,                  // auto-ack
		false,                  // exclusive
		false,                  // no-local
		false,                  // no-wait
		nil,                    // args
	)
	if err != nil {
		log.Fatalf("error on opening message channel. err=%v", err)
	}

	go func() {
		for message := range messageCh {
			log.Printf("consuming '%s' workflow. payload=%s", w.id, string(message.Body))
		}
	}()
}
