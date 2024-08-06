package pgdb

import (
	"errors"
	log "github.com/sirupsen/logrus"
	"gopkg.in/reform.v1"
	"news_api/internal/model/dbmodel"
	"news_api/pkg/postgres"
)

const (
	categoriesPrefixLog = "/pgdb/categories"
)

type CategoriesRepo struct {
	postgres.Postgres
}

func NewCategoriesRepo(pg postgres.Postgres) *CategoriesRepo {
	return &CategoriesRepo{pg}
}

func (r *CategoriesRepo) FindNewsCategories(newsId int) ([]int, error) {
	rows, err := r.SelectRows(dbmodel.CategoriesTable, "WHERE news_id = $1", newsId)
	if err != nil {
		if errors.Is(err, reform.ErrNoRows) { // это конечно невозможно, но пусть будет
			return []int{}, nil
		}
		log.Errorf("%s/FindAll error find news categories: %s", categoriesPrefixLog, err)
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var result []int
	for rows.Next() {
		var category dbmodel.Categories

		err = rows.Scan(
			&category.Id,
			&category.NewsId,
			&category.CategoryId,
		)
		if err != nil {
			log.Errorf("%s/FindAll error scan category row: %s", categoriesPrefixLog, err)
			return nil, err
		}
		result = append(result, category.CategoryId)
	}
	return result, nil
}
