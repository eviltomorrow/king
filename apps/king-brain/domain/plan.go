package domain

import (
	"encoding/json"

	"github.com/eviltomorrow/king/apps/king-brain/domain/chart"
)

type Plan struct {
	K *chart.K
}

func (p *Plan) String() string {
	buf, _ := json.Marshal(p)
	return string(buf)
}
