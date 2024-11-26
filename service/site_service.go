package service

import (
	"context"
	"database/sql"
	"io"
	"os"
	"radproject/entity"
	"radproject/model"
	"radproject/model/message"
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
	linkRepository *repository.LinkRepository
}

func NewSiteService(validator *validator.Validate, db *sql.DB, siteRepository *repository.SiteRepository, linkRepository *repository.LinkRepository) *SiteService {
	return &SiteService{
		validator:      validator,
		db:             db,
		siteRepository: siteRepository,
		linkRepository: linkRepository,
	}
}

func (s *SiteService) GetAllSites(ctx context.Context, claims jwt.MapClaims) ([]model.SiteResponse, error) {
	res := make([]model.SiteResponse, 0)

	sites, err := s.siteRepository.SelectAllByUserID(ctx, s.db, claims["id"].(string))
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

	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	domain := strings.ReplaceAll(request.Domain, " ", "_")
	_, err = s.siteRepository.SelectWhere(ctx, tx, "domain = ?", domain)
	if err == nil {
		tx.Rollback()
		return message.ERR_CREATE_SITE_DOMAIN_USED
	}

	site := &entity.Sites{
		Id:      uuid.New().String(),
		Domain:  domain,
		Title:   request.Title,
		User_Id: claims["id"].(string),
	}

	split := strings.Split(request.Image.Filename, ".")
	ext := split[len(split)-1]
	filename := "site-" + site.Id + "-" + strconv.Itoa(int(time.Now().Unix())) + "." + ext

	site.Image = filename

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

	site, err := s.siteRepository.SelectWithIdAndUser(ctx, s.db, req.Id, claims["id"].(string))
	if err != nil {
		return err
	}

	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	err = s.linkRepository.DeleteBySite(ctx, tx, site.Id)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = s.siteRepository.Delete(ctx, tx, req.Id, claims["id"].(string))
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (s *SiteService) GetSite(ctx context.Context, claims jwt.MapClaims, id_site string) (*model.SiteResponse, error) {
	err := s.validator.Var(id_site, "required")
	if err != nil {
		return nil, err
	}

	site, err := s.siteRepository.SelectWithJoinLinkByUser(ctx, s.db, id_site, claims["id"].(string))
	if err != nil {
		return nil, err
	}

	res := model.EntityToSiteResponse(site)

	return res, err
}

func (s *SiteService) UpdateSite(ctx context.Context, claims jwt.MapClaims, req *model.UpdateSiteRequest) error {
	err := s.validator.Struct(req)
	if err != nil {
		return err
	}

	site, err := s.siteRepository.SelectWithIdAndUser(ctx, s.db, req.Id, claims["id"].(string))
	if err != nil {
		return err
	}

	if req.Domain != "" {
		site.Domain = strings.ReplaceAll(req.Domain, " ", "_")
	}
	if req.Title != "" {
		site.Title = req.Title
	}
	if req.Image != nil {
		split := strings.Split(req.Image.Filename, ".")
		ext := split[len(split)-1]
		filename := "site-" + site.Id + "-" + strconv.Itoa(int(time.Now().Unix())) + "." + ext

		site.Image = filename
	}

	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	err = s.siteRepository.Update(ctx, tx, site)
	if err != nil {
		tx.Rollback()
		return err
	}
	if req.Image != nil {
		open, err := req.Image.Open()
		if err != nil {
			tx.Rollback()
			return err
		}
		defer open.Close()
		dsn, err := os.Create("storage/image/site/" + site.Image)
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
	}
	return tx.Commit()
}

func (s *SiteService) GetSiteWithDomain(ctx context.Context, req *model.ViewPublishSite) (*model.SiteResponse, error) {
	err := s.validator.Struct(req)
	if err != nil {
		return nil, err
	}

	split := strings.Split(req.DomainSite, ".")
	domain := strings.Join(split[:len(split)-1], ".")
	if split[len(split)-1] != "site" {
		return nil, message.ERR_PUBLISH_SITE_DOMAIN
	}

	site, err := s.siteRepository.SelectWithJoinLinkByDomain(ctx, s.db, domain)
	if err != nil {
		return nil, err
	}
	return model.EntityToSiteResponse(site), nil
}
