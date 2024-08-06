package service

import (
	"errors"
	"news_api/internal/model/dbmodel"
	"news_api/internal/repo"
	"news_api/internal/repo/pgerrs"
)

type newsService struct {
	news       repo.News
	categories repo.Categories
}

func newNewsService(news repo.News, categories repo.Categories) *newsService {
	return &newsService{
		news:       news,
		categories: categories,
	}
}

func (s *newsService) CreateNews(input NewsInput) (int, error) {
	news := &dbmodel.News{
		Title:   input.Title,
		Content: input.Content,
	}
	if err := s.news.CreateNews(news, input.Categories); err != nil {
		return 0, err
	}
	return news.Id, nil
}

func (s *newsService) UpdateNews(id int, input NewsInput) error {
	news := &dbmodel.News{
		Id:      id,
		Title:   input.Title,
		Content: input.Content,
	}
	if err := s.news.UpdateNews(news, input.Categories); err != nil {
		if errors.Is(err, pgerrs.ErrNotFound) {
			return ErrNewsNotFound
		}
		return err
	}
	return nil
}

func (s *newsService) FindNews(limit, offset int) ([]NewsOutput, error) {
	news, err := s.news.FindNews(limit, offset)
	if err != nil {
		return nil, err
	}

	var result []NewsOutput
	for _, n := range news {
		categories, e := s.categories.FindNewsCategories(n.Id)
		if e != nil {
			return nil, e
		}
		result = append(result, NewsOutput{
			Id:         n.Id,
			Title:      n.Title,
			Content:    n.Content,
			Categories: categories,
		})
	}

	return result, nil
}
