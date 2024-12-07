package saga

// workflow is to contain and orchestrate activities
type workflow struct {
	activities []activity
}

func NewWorkflow() *workflow {
	return &workflow{}
}

func (w *workflow) Start() {}
func (w *workflow) Stop()  {}

func (s *workflow) RegisterActivities(activities ...activity) {
	s.activities = activities
}
