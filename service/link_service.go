package service

import (
	"database/sql"
	"radproject/repository"

	"github.com/go-playground/validator/v10"
)

type LinkService struct {
	validator      *validator.Validate
	db             *sql.DB
	linkRepository *repository.LinkRepository
}

func NewLinkService(validator *validator.Validate, db *sql.DB, repository *repository.LinkRepository) *LinkService {
	return &LinkService{
		validator:      validator,
		db:             db,
		linkRepository: repository,
	}
}
