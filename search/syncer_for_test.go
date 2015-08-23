package search_test

type Syncer struct {
	goroutines int
	ready, run chan bool
}

func NewSyncer() *Syncer {
	return &Syncer{
		ready: make(chan bool),
		run:   make(chan bool),
	}
}

func (s *Syncer) Register() { s.goroutines++ }

func (s *Syncer) Sync() {
	select {
	case s.ready <- true:
		<-s.run
	case <-s.run:
	}
}

func (s *Syncer) WaitForAllReady() *Syncer {
	for i := 0; i < s.goroutines; i++ {
		<-s.ready
	}
	return s
}

func (s *Syncer) LetAllRun() {
	for i := 0; i < s.goroutines; i++ {
		s.run <- true
	}
}
