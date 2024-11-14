package plan

import "github.com/eviltomorrow/king/apps/king-cron/domain"

const (
	NameWithTrendFollowing  = "CronWithTrendFollowing"
	AliasWithTrendFollowing = "趋势追踪"
)

func CronWithTrendFollowing() *domain.Plan {
	return &domain.Plan{
		Precondition: func() (domain.StatusCode, error) {
			return 0, nil
		},
		Todo: func(string) error {
			return nil
		},
		WriteToDB: func(string, error) error {
			return nil
		},

		Status: domain.Ready,
		Name:   NameWithTrendFollowing,
		Alias:  AliasWithTrendFollowing,
	}
}
