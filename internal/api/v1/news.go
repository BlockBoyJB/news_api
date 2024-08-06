package v1

import (
	"errors"
	"github.com/gofiber/fiber/v3"
	"net/http"
	"news_api/internal/service"
	"news_api/pkg/validator"
)

type newsRouter struct {
	news      service.News
	validator validator.Validator
}

func newNewsRouter(g fiber.Router, news service.News, validator validator.Validator) {
	r := &newsRouter{
		news:      news,
		validator: validator,
	}

	g.Post("/create", r.create)
	g.Post("/edit", r.edit) // я бы использовал PATCH, но задание есть задание
	g.Get("/list", r.list)
}

type newsCreateInput struct {
	Title      string `json:"Title" validate:"required"`
	Content    string `json:"Content" validate:"required"`
	Categories []int  `json:"Categories" validate:"required,category"`
}

func (r *newsRouter) create(c fiber.Ctx) error {
	var input newsCreateInput

	if err := c.Bind().Body(&input); err != nil {
		errorResponse(c, http.StatusBadRequest, fiber.ErrBadRequest)
		return nil
	}

	if err := r.validator.Validate(input); err != nil {
		errorResponse(c, http.StatusBadRequest, err)
		return nil
	}

	id, err := r.news.CreateNews(service.NewsInput{
		Title:      input.Title,
		Content:    input.Content,
		Categories: input.Categories,
	})
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, fiber.ErrInternalServerError)
		return err
	}

	type response struct {
		Id int `json:"id"`
	}
	return c.JSON(response{Id: id})
}

type newsEditInput struct {
	Id         int    `json:"Id"  validate:"required"`
	Title      string `json:"Title"`
	Content    string `json:"Content"`
	Categories []int  `json:"Categories" validate:"category"`
}

func (r *newsRouter) edit(c fiber.Ctx) error {
	var input newsEditInput

	if err := c.Bind().Body(&input); err != nil {
		errorResponse(c, http.StatusBadRequest, fiber.ErrBadRequest)
		return nil
	}

	if err := r.validator.Validate(input); err != nil {
		errorResponse(c, http.StatusBadRequest, err)
		return nil
	}

	err := r.news.UpdateNews(input.Id, service.NewsInput{
		Title:      input.Title,
		Content:    input.Content,
		Categories: input.Categories,
	})
	if err != nil {
		if errors.Is(err, service.ErrNewsNotFound) {
			errorResponse(c, http.StatusBadRequest, err)
			return nil
		}
		errorResponse(c, http.StatusInternalServerError, fiber.ErrInternalServerError)
		return err
	}
	return c.SendStatus(http.StatusOK)
}

type newsListInput struct {
	Limit  int `json:"Limit" validate:"limit"`
	Offset int `json:"Offset" validate:"value"`
}

func (r *newsRouter) list(c fiber.Ctx) error {
	var input newsListInput

	if err := c.Bind().Body(&input); err != nil {
		errorResponse(c, http.StatusBadRequest, fiber.ErrBadRequest)
		return err
	}

	if err := r.validator.Validate(input); err != nil {
		errorResponse(c, http.StatusBadRequest, err)
		return nil
	}

	news, err := r.news.FindNews(input.Limit, input.Offset)
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, fiber.ErrInternalServerError)
		return err
	}

	type response struct {
		Success bool                 `json:"Success"`
		News    []service.NewsOutput `json:"News"`
	}
	return c.JSON(response{
		Success: true,
		News:    news,
	})
}
