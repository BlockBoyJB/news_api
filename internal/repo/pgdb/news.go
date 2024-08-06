package pgdb

import (
	"errors"
	"github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"gopkg.in/reform.v1"
	"news_api/internal/model/dbmodel"
	"news_api/internal/repo/pgerrs"
	"news_api/pkg/postgres"
)

const (
	newsPrefixLog = "/pgdb/news"
)

type NewsRepo struct {
	postgres.Postgres
}

func NewNewsRepo(pg postgres.Postgres) *NewsRepo {
	return &NewsRepo{pg}
}

func (r *NewsRepo) CreateNews(news *dbmodel.News, categories []int) error {
	tx, err := r.Begin()
	if err != nil {
		log.Errorf("%s/CreateNews error init tx: %s", newsPrefixLog, err)
		return err
	}
	defer func() { _ = tx.Rollback() }()

	if err = tx.Save(news); err != nil {
		log.Errorf("%s/CreateNews error save news: %s", newsPrefixLog, err)
		return err
	}

	for _, c := range categories {
		category := &dbmodel.Categories{
			NewsId:     news.Id,
			CategoryId: c,
		}
		if err = tx.Save(category); err != nil {
			log.Errorf("%s/CreateNews error save news category: %s", newsPrefixLog, err)
			return err
		}
	}

	if err = tx.Commit(); err != nil {
		log.Errorf("%s/CreateNews error commit: %s", newsPrefixLog, err)
		return err
	}

	return nil
}

func (r *NewsRepo) UpdateNews(news *dbmodel.News, categories []int) error {
	// если в категории приходит не пустой массив, то нужно удалить все прошлые категории из категорий и записать новые
	tx, err := r.Begin()
	if err != nil {
		log.Errorf("%s/UpdateNews error init tx: %s", newsPrefixLog, err)
		return err
	}
	defer func() { _ = tx.Rollback() }()

	var columns []string
	if news.Title != "" {
		columns = append(columns, "title")
	}
	if news.Content != "" {
		columns = append(columns, "content")
	}

	if columns != nil {
		if err = tx.UpdateColumns(news, columns...); err != nil {
			if errors.Is(err, reform.ErrNoRows) {
				return pgerrs.ErrNotFound
			}
			log.Errorf("%s/UpdateNews error update news: %s", newsPrefixLog, err)
			return err
		}
	}

	if len(categories) != 0 {
		if _, err = tx.DeleteFrom(dbmodel.CategoriesTable, "WHERE news_id = $1", news.Id); err != nil {
			if errors.Is(err, reform.ErrNoRows) {
				return pgerrs.ErrNotFound
			}
			log.Errorf("%s/UpdateNews error delete old categories: %s", newsPrefixLog, err)
			return err
		}

		for _, c := range categories {
			category := &dbmodel.Categories{
				NewsId:     news.Id,
				CategoryId: c,
			}
			if err = tx.Save(category); err != nil {
				// интересно, что если на вход дать неправильный id и только новый список категорий,
				// то ошибка отсутствия новости проявится только здесь
				var pgErr *pq.Error

				if ok := errors.As(err, &pgErr); ok {
					if pgErr.Code == "23503" {
						return pgerrs.ErrNotFound
					}
				}
				log.Errorf("%s/UpdateNews error save news category: %s", newsPrefixLog, err)
				return err
			}
		}
	}

	if err = tx.Commit(); err != nil {
		log.Errorf("%s/UpdateNews error commit: %s", newsPrefixLog, err)
		return err
	}

	return nil
}

func (r *NewsRepo) FindNews(limit, offset int) ([]dbmodel.News, error) {
	rows, err := r.SelectRows(dbmodel.NewsTable, "ORDER BY id LIMIT $1 OFFSET $2", limit, offset)
	if err != nil {
		log.Errorf("%s/FindNews error find news: %s", newsPrefixLog, err)
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var result []dbmodel.News
	for rows.Next() {
		var news dbmodel.News

		err = rows.Scan(
			&news.Id,
			&news.Title,
			&news.Content,
		)
		if err != nil {
			log.Errorf("%s/FindNews error scan news row: %s", newsPrefixLog, err)
			return nil, err
		}
		result = append(result, news)
	}

	return result, nil
}
