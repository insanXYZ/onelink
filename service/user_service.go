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
	b, _ := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)

	user := &entity.User{
		ID:       uuid.New().String(),
		Name:     request.Name,
		Email:    request.Email,
		Image:    "default_user.jpeg",
		Password: string(b),
	}

	return s.userRepository.Save(ctx, s.db, user)
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
