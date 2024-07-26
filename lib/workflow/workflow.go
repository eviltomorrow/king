package workflow

type Job struct {
	Name string
	F    func() error
}

var cache []Job

func Register(name string, f func() error) {
	job := Job{
		Name: name,
		F:    f,
	}
	cache = append(cache, job)
}

func Finish() []Job {
	return cache
}
