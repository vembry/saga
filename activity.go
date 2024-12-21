package main

// activity a unit of execution that contain both commit and rollback
type activity struct {
	name string

	commit   func()
	rollback func()
}

func NewActivity(name string, commit func(), rollback func()) *activity {
	return &activity{
		name:     name,
		commit:   commit,
		rollback: rollback,
	}
}

func (a *activity) Start() {}
func (a *activity) Stop()  {}
