package service

import (
	"context"
	"database/sql"
	"radproject/entity"
	"radproject/model"
	"radproject/repository"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type LinkService struct {
	validator      *validator.Validate
	db             *sql.DB
	linkRepository *repository.LinkRepository
	siteRepository *repository.SiteRepository
}

func NewLinkService(validator *validator.Validate, db *sql.DB, linkRepository *repository.LinkRepository, siteRepository *repository.SiteRepository) *LinkService {
	return &LinkService{
		validator:      validator,
		db:             db,
		linkRepository: linkRepository,
		siteRepository: siteRepository,
	}
}

func (s *LinkService) CreateLink(ctx context.Context, claims jwt.MapClaims, req *model.CreateLinkRequest) error {
	err := s.validator.Struct(req)
	if err != nil {
		return err
	}

	site, err := s.siteRepository.SelectWithIdAndUser(ctx, s.db, req.Site_Id, claims["id"].(string))
	if err != nil {
		return err
	}

	link := &entity.Links{
		Id:      uuid.New().String(),
		Site_Id: site.Id,
		Title:   req.Title,
		Href:    req.Href,
	}

	return s.linkRepository.Save(ctx, s.db, link)
}

func (s *LinkService) Delete(ctx context.Context, claims jwt.MapClaims, req *model.DeleteLinkRequest) error {
	err := s.validator.Struct(req)
	if err != nil {
		return err
	}

	site, err := s.siteRepository.SelectWithIdAndUser(ctx, s.db, req.Site_Id, claims["id"].(string))
	if err != nil {
		return err
	}

	return s.linkRepository.Delete(ctx, s.db, req.Link_Id, site.Id)
}

func (s *LinkService) Visit(ctx context.Context, req *model.VisitLinkRequest) (*model.LinkResponse, error) {
	err := s.validator.Struct(req)
	if err != nil {
		return nil, err
	}

	link, err := s.linkRepository.SelectWhere(ctx, s.db, "id = ?", req.Link_Id)
	if err != nil {
		return nil, err
	}

	return model.EntitytoLinkResponse(link), nil
}
