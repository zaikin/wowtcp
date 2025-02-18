package service

//go:generate mockery --name=Repository --output=./mocks --outpkg=mocks

type Repository interface {
	GetWoWQuote() string
}

type Service struct {
	r Repository
}

func NewService(r Repository) *Service {
	return &Service{
		r: r,
	}
}
