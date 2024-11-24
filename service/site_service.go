package service

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"os"
	"radproject/entity"
	"radproject/model"
	"radproject/repository"
	"strconv"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/insanXYZ/sage"
)

type SiteService struct {
	validator      *validator.Validate
	db             *sql.DB
	siteRepository *repository.SiteRepository
	LinkRepository *repository.LinkRepository
}

func NewSiteService(validator *validator.Validate, db *sql.DB, siteRepository *repository.SiteRepository, linkRepository *repository.LinkRepository) *SiteService {
	return &SiteService{
		validator:      validator,
		db:             db,
		siteRepository: siteRepository,
		LinkRepository: linkRepository,
	}
}

func (s *SiteService) GetAllSites(ctx context.Context, claims jwt.MapClaims) ([]model.SiteResponse, error) {
	res := make([]model.SiteResponse, 0)

	sites, err := s.siteRepository.SelectAllById(ctx, s.db, claims["id"].(string))
	if err != nil {
		return res, err
	}

	for _, v := range sites {
		res = append(res, *model.EntityToSiteResponse(&v))
	}

	return res, err
}

func (s *SiteService) CreateSite(ctx context.Context, claims jwt.MapClaims, request *model.CreateSiteRequest) error {
	err := s.validator.Struct(request)
	if err != nil {
		return err
	}

	err = sage.Validate(request.Image)
	if err != nil {
		return err
	}

	site := &entity.Sites{
		Id:      uuid.New().String(),
		Domain:  strings.ReplaceAll(request.Domain, " ", "_"),
		Title:   request.Title,
		User_Id: claims["id"].(string),
	}

	split := strings.Split(request.Image.Filename, ".")
	ext := split[len(split)-1]
	filename := "site-" + site.Id + "-" + strconv.Itoa(int(time.Now().Unix())) + "." + ext

	site.Image = filename

	fmt.Println(*site)

	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	err = s.siteRepository.Save(ctx, tx, site)
	if err != nil {
		return err
	}

	open, err := request.Image.Open()
	if err != nil {
		tx.Rollback()
		return err
	}
	defer open.Close()

	dsn, err := os.Create("storage/image/site/" + filename)
	if err != nil {
		tx.Rollback()
		return err
	}
	defer dsn.Close()

	_, err = io.Copy(dsn, open)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (s *SiteService) DeleteSite(ctx context.Context, claims jwt.MapClaims, req *model.DeleteSiteRequest) error {
	err := s.validator.Struct(req)
	if err != nil {
		return err
	}
	return s.siteRepository.Delete(ctx, s.db, req.Id, claims["id"].(string))
}
