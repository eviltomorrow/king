package service

var repository = map[string]*Model{}

type Model struct {
	Name string
	F    func() (int, bool)
}

func RegisterModel(name string, model *Model) {
	repository[name] = model
}

func FindPossibleChance() {
}
