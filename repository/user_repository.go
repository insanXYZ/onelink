package repository

import (
	"context"
	"database/sql"

	"radproject/entity"
)

type UserRepository struct{}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

func (r *UserRepository) Save(ctx context.Context, db *sql.DB, model *entity.User) error {
	query := "insert into users(id,name,email,password,image) values(?,?,?,?,?)"
	_, err := db.ExecContext(ctx, query, model.ID, model.Name, model.Email, model.Password, model.Image)
	return err
}

func (r *UserRepository) Select(ctx context.Context, db *sql.DB, with, value string) (*entity.User, error) {
	user := new(entity.User)

	query := "select * from users where " + with + " = ?"
	err := db.QueryRowContext(ctx, query, value).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.Image,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	return user, err
}

func (r *UserRepository) UpdateAccount(ctx context.Context, db *sql.DB, ent *entity.User) error {
	query := "update users set name = ? , email = ?, image = ? where id = ?"
	_, err := db.ExecContext(ctx, query, ent.Name, ent.Email, ent.Image, ent.ID)
	return err
}
