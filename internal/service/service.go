package service

import "news_api/internal/repo"

type (
	NewsInput struct {
		Title      string
		Content    string
		Categories []int
	}
	NewsOutput struct {
		Id         int    `json:"Id"`
		Title      string `json:"Title"`
		Content    string `json:"Content"`
		Categories []int  `json:"Categories"`
	}
)

type Auth interface {
	ValidateToken(token string) bool
	CreateToken() (string, error)
}

type News interface {
	CreateNews(input NewsInput) (int, error)
	UpdateNews(id int, input NewsInput) error
	FindNews(limit, offset int) ([]NewsOutput, error)
}

type (
	Services struct {
		Auth Auth
		News News
	}
	ServicesDependencies struct {
		Repos      *repo.Repositories
		PrivateKey string
		PublicKey  string
	}
)

func NewServices(d *ServicesDependencies) *Services {
	return &Services{
		Auth: newAuthService(d.PrivateKey, d.PublicKey),
		News: newNewsService(d.Repos.News, d.Repos.Categories),
	}
}
