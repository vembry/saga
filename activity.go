package saga

// activity a unit of execution that contain both commit and rollback
type activity struct{}

func NewActivity() *activity {
	return &activity{}
}

func (a *activity) Start() {}
func (a *activity) Stop()  {}
