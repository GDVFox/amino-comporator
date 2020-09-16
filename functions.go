package main

// MaxInt returns max int
func MaxInt(a, b int) int {
	if a > b {
		return a
	}

	return b
}

// MinInt returns min int
func MinInt(a, b int) int {
	if a < b {
		return a
	}

	return b
}

type Semaphore struct {
	ch chan struct{}
}

func NewSemaphore(cap uint) *Semaphore {
	return &Semaphore{
		ch: make(chan struct{}, cap),
	}
}

func (s *Semaphore) Accuire() {
	s.ch <- struct{}{}
}

func (s *Semaphore) Release() {
	<-s.ch
}

func (s *Semaphore) Len() int {
	return len(s.ch)
}
