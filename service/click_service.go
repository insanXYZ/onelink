package service

import (
	"context"
	"database/sql"
	"radproject/model"
	"radproject/model/message"
	"radproject/repository"
	"strings"

	"github.com/go-playground/validator/v10"
)

type ClickService struct {
	validator       *validator.Validate
	db              *sql.DB
	linkRepository  *repository.LinkRepository
	siteRepository  *repository.SiteRepository
	clickRepository *repository.ClickRepository
}

func NewClickService(validator *validator.Validate, db *sql.DB, linkRepository *repository.LinkRepository, siteRepository *repository.SiteRepository, clickRepository *repository.ClickRepository) *ClickService {
	return &ClickService{
		validator:       validator,
		db:              db,
		linkRepository:  linkRepository,
		siteRepository:  siteRepository,
		clickRepository: clickRepository,
	}
}

func (s *ClickService) Visit(ctx context.Context, req *model.VisitDestination) error {
	var id string

	if req.Link_Id != "" {
		link, err := s.linkRepository.SelectWhere(ctx, s.db, "id = ?", req.Link_Id)
		if err != nil {
			return err
		}
		id = link.Id
	}

	if req.Domain != "" {
		split := strings.Split(req.Domain, ".")
		domain := strings.Join(split[:len(split)-1], ".")
		site, err := s.siteRepository.SelectWhere(ctx, s.db, "domain = ?", domain)
		if err != nil {
			return err
		}
		id = site.Id
	}

	if id != "" {
		err := s.clickRepository.Save(ctx, s.db, id)
		if err != nil {
			return err
		}
		return nil
	}

	return message.ERR_CLICK
}
