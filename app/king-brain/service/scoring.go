package service

type Model interface {
	Desc() string
	Mark(*DataWrapper) (int, error)
}

var repository = func() []Model {
	return make([]Model, 0, 16)
}()

func RegisterModel(m Model) {
	repository = append(repository, m)
}

func Scoring(data *DataWrapper) (int, error) {
	return 0, nil
}
