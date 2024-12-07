package saga

// saga is the client to orchestrate workflows
type saga struct {
	workflows []workflow
}

func New() *saga {
	return &saga{}
}

func (s *saga) Start() {}
func (s *saga) Stop()  {}

func (s *saga) RegisterWorkflows(workflows ...workflow) {
	s.workflows = workflows
}
