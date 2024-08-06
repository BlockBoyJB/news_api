package repo

import (
	"news_api/internal/model/dbmodel"
	"news_api/internal/repo/pgdb"
	"news_api/pkg/postgres"
)

type News interface {
	CreateNews(news *dbmodel.News, categories []int) error
	UpdateNews(news *dbmodel.News, categories []int) error
	FindNews(limit, offset int) ([]dbmodel.News, error)
}

type Categories interface {
	FindNewsCategories(newsId int) ([]int, error)
}

type Repositories struct {
	News
	Categories
}

func NewRepositories(pg postgres.Postgres) *Repositories {
	return &Repositories{
		News:       pgdb.NewNewsRepo(pg),
		Categories: pgdb.NewCategoriesRepo(pg),
	}
}
