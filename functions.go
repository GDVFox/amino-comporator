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

// Semaphore reprasents sync primitive semaphore
type Semaphore struct {
	ch chan struct{}
}

// NewSemaphore returns new instance of Semaphore
func NewSemaphore(cap uint) *Semaphore {
	return &Semaphore{
		ch: make(chan struct{}, cap),
	}
}

// Acquire quota
func (s *Semaphore) Acquire() {
	s.ch <- struct{}{}
}

// Release quota
func (s *Semaphore) Release() {
	<-s.ch
}

// Len retruns total count of quote
func (s *Semaphore) Len() int {
	return len(s.ch)
}
