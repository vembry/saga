package main

import "log"

// workflow is to contain and orchestrate activities
type workflow struct {
	activities []activity
}

func NewWorkflow() *workflow {
	return &workflow{}
}

func (w *workflow) Start() {
	log.Printf("starting workflow...")
}
func (w *workflow) Stop() {
	log.Printf("stopping workflow...")
}

func (s *workflow) RegisterActivities(activities ...activity) {
	s.activities = activities
}
