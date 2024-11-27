package service

import (
	"context"
	"database/sql"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"radproject/entity"
	"radproject/model"
	"radproject/model/message"
	"radproject/repository"
	"radproject/util"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/insanXYZ/sage"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	validator      *validator.Validate
	db             *sql.DB
	userRepository *repository.UserRepository
	siteRepository *repository.SiteRepository
	linkRepository *repository.LinkRepository
}

func NewUserService(validator *validator.Validate, db *sql.DB, repository *repository.UserRepository) *UserService {
	return &UserService{
		validator:      validator,
		db:             db,
		userRepository: repository,
	}
}

func (s *UserService) Login(ctx context.Context, request *model.LoginRequest) (string, error) {
	err := s.validator.Struct(request)
	if err != nil {
		return "", err
	}
	user, err := s.userRepository.Select(ctx, s.db, "email", request.Email)
	if err != nil {
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password))
	if err != nil {
		return "", err
	}

	exp, _ := strconv.Atoi(os.Getenv("JWT_EXP"))

	claim := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"name": user.Name,
		"id":   user.ID,
		"exp":  time.Now().Add(time.Duration(exp) * time.Minute).Unix(),
	})

	return claim.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))
}

func (s *UserService) Register(ctx context.Context, request *model.RegisterRequest) error {
	err := s.validator.Struct(request)
	if err != nil {
		return util.HandleValidatorStruct(err)
	}

	_, err = s.userRepository.Select(ctx, s.db, "email", request.Email)
	if err == nil {
		return message.ERR_REGISTER_EMAIL_USED
	}
	b, _ := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)

	userCreate := &entity.User{
		ID:       uuid.New().String(),
		Name:     request.Name,
		Email:    request.Email,
		Image:    "default_user.jpeg",
		Password: string(b),
	}

	return s.userRepository.Save(ctx, s.db, userCreate)
}

func (s *UserService) GetAccount(ctx context.Context, claims jwt.MapClaims) (*model.UserResponse, error) {
	user, err := s.userRepository.Select(ctx, s.db, "id", claims["id"].(string))
	if err != nil {
		return nil, err
	}

	resp := model.EntityToUserResponse(user)
	return resp, nil
}

func (s *UserService) UpdateUser(ctx context.Context, claims jwt.MapClaims, request *model.UpdateUserRequest) error {
	err := s.validator.Struct(request)
	if err != nil {
		return err
	}

	if request.Image != nil {
		err = sage.Validate(request.Image)
		if err != nil {
			return err
		}
	}

	user, err := s.userRepository.Select(ctx, s.db, "id", claims["id"].(string))
	if err != nil {
		return err
	}

	if request.Name != "" {
		user.Name = request.Name
	}
	if request.Email != "" {
		user.Email = request.Email
	}
	if request.Image != nil {

		open, err := request.Image.Open()
		split := strings.Split(request.Image.Filename, ".")
		ext := split[len(split)-1]

		filename := "user-" + user.ID + "-" + strconv.Itoa(int(time.Now().Unix())) + "." + ext
		dsn, err := os.Create("storage/image/user/" + filename)
		if err != nil {
			return err
		}
		defer dsn.Close()

		_, err = io.Copy(dsn, open)
		if err != nil {
			return err
		}

		user.Image = filename
	}

	return s.userRepository.UpdateAccount(ctx, s.db, user)

}

func (s *UserService) Dashboard(ctx context.Context, claims jwt.MapClaims, req *model.DashboardRequest) (*model.UserResponse, int, int, error) {
	var t time.Time
	var sumSite, sumLink int
	if req.From == "" {
		t = time.Now()
		req.From = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location()).Format(time.DateTime)
	} else {
		tt, err := time.Parse(time.DateOnly, req.From)
		if err != nil {
			return nil, 0, 0, err
		}
		t = tt
		req.From = time.Date(tt.Year(), tt.Month(), tt.Day(), 0, 0, 0, 0, tt.Location()).Format(time.DateTime)
	}

	if req.To == "" {
		req.To = time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 59, t.Location()).Format(time.DateTime)
	} else {
		t, err := time.Parse(time.DateOnly, req.To)
		if err != nil {
			return nil, 0, 0, err
		}
		req.To = time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 59, t.Location()).Format(time.DateTime)
	}

	user, err := s.userRepository.Select(ctx, s.db, "id", claims["id"].(string))
	if err != nil {
		return nil, 0, 0, err
	}

	userResponse := model.EntityToUserResponse(user)

	sites, ids, err := s.siteRepository.GetAllWithNumberClickByUser(ctx, s.db, user.ID, req)
	if err != nil {
		return nil, 0, 0, err
	}

	if len(sites) != 0 {

		for _, v := range sites {
			sumSite += v.Clicks
			userResponse.SiteResponse = append(userResponse.SiteResponse, *model.EntityToSiteResponse(&v))
		}

		links, err := s.linkRepository.GetAllWithNumberClickBySiteId(ctx, s.db, ids, req)
		if err != nil {
			return nil, 0, 0, err
		}

		if len(links) != 0 {
			for _, vv := range links {
				sumLink += vv.Clicks
				userResponse.LinkResponse = append(userResponse.LinkResponse, *model.EntitytoLinkResponse(&vv))
			}
		}
	}

	return userResponse, sumSite, sumLink, nil

}
