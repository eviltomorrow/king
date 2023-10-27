package scoring

import "github.com/eviltomorrow/king/apps/king-brain/service"

type Barycentric struct {
	Name      string
	ScoreStep []int
}

func (s *Barycentric) Desc() string {
	return s.Name
}

func (s *Barycentric) Mark(data *service.DataWrapper) (int, error) {
	return 0, nil
}

func init() {
	service.RegisterModel(&Barycentric{
		Name: "重心移动模型",
		ScoreStep: []int{
			0,
			10,
		},
	})
}
