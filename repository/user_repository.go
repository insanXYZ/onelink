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

func (r *UserRepository) GetWithEmail(ctx context.Context, db *sql.DB, email string) (*entity.User, error) {
	query := "select * from users where email = ?"
	rows, err := db.QueryContext(ctx, query, email)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	user := entity.User{}

	for rows.Next() {
		err = rows.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Image, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, err
		}
	}

	return &user, nil
}

func (r *UserRepository) SearchById(ctx context.Context, db *sql.DB, id string) (*entity.User, error) {
	user := new(entity.User)

	query := "select * from users where id = ?"
	err := db.QueryRowContext(ctx, query, id).Scan(
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
