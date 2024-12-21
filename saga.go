package main

import "log"

type IProcess interface {
	Start()
	Stop()
}

// saga is the client to orchestrate workflows
type saga struct {
	workflows []IProcess
}

func New() *saga {
	return &saga{}
}

func (s *saga) Start() {
	log.Printf("starting saga...")
}
func (s *saga) Stop() {
	log.Printf("stopping saga...")
}

func (s *saga) RegisterWorkflows(workflows ...IProcess) {
	s.workflows = workflows
}
