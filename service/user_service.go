package service

import (
	"context"
	"database/sql"
	"os"
	"time"

	"radproject/entity"
	"radproject/model"
	"radproject/repository"
	"radproject/util"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	validator      *validator.Validate
	db             *sql.DB
	userRepository *repository.UserRepository
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
	user, err := s.userRepository.GetWithEmail(ctx, s.db, request.Email)
	if err != nil {
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password))
	if err != nil {
		return "", err
	}

	claim := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"name": user.Name,
		"exp":  time.Now().Add(5 * time.Minute).Unix(),
	})

	return claim.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))
}

func (s *UserService) Register(ctx context.Context, request *model.RegisterRequest) error {
	err := s.validator.Struct(request)
	if err != nil {
		return util.HandleValidatorStruct(err)
	}
	b, _ := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)

	return s.userRepository.Save(ctx, s.db, &entity.User{
		ID:       uuid.New().String(),
		Name:     request.Name,
		Email:    request.Email,
		Password: string(b),
	})
}
