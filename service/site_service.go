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
}

func NewSiteService(validator *validator.Validate, db *sql.DB, repository *repository.SiteRepository) *SiteService {
	return &SiteService{
		validator:      validator,
		db:             db,
		siteRepository: repository,
	}
}

func (s *SiteService) GetAllSites(ctx context.Context, claims jwt.MapClaims) ([]entity.Sites, error) {
	return s.siteRepository.SelectAllById(ctx, s.db, claims["id"].(string))
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
