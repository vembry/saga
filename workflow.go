package main

import "log"

// workflow is to contain and orchestrate activities
type workflow[T any] struct {
	name string

	activities []*activity

	// NOTE: directly accessing saga client for now
	// TODO: this should be changed
	sagaclient *saga
}

func NewWorkflow[T any](
	name string,
	activities ...*activity,
) *workflow[T] {
	return &workflow[T]{
		name:       name,
		activities: activities,
	}
}

func (w *workflow[T]) Start(s *saga) {
	log.Printf("starting workflow...")

	w.sagaclient = s
}
func (w *workflow[T]) Stop() {
	log.Printf("stopping workflow...")
}

// Execute starts a process of a workflow. A workflow journey starts here
func (w *workflow[T]) Execute(param T) {
	// create entry into mongo
	w.sagaclient.CreateEntry(w.name, param)

	log.Printf("executing '%s' workflow. payload=%+v", w.name, param)
}
