package saga

type saga struct{}

func New() *saga {
	return &saga{}
}

func (s *saga) Start() {}
func (s *saga) Stop()  {}
