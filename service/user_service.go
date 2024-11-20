package service

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strconv"
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
	fmt.Println(*request)
	user, err := s.userRepository.GetWithEmail(ctx, s.db, request.Email)
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
	b, _ := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)

	user := &entity.User{
		ID:       uuid.New().String(),
		Name:     request.Name,
		Email:    request.Email,
		Image:    "default_user.jpeg",
		Password: string(b),
	}

	fmt.Println(*user, request.Password)

	return s.userRepository.Save(ctx, s.db, user)
}
