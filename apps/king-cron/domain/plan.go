package domain

type Plan struct {
	Precondition func() (bool, error)
	Todo         func() error
}

const (
	StatusProcessing = "processing"
	StatusCompleted  = "completed"
)
