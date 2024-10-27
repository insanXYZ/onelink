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
	query := "insert into users(id,name,email,password) values(?,?,?,?)"
	_, err := db.ExecContext(ctx, query, model.ID, model.Name, model.Email, model.Password)
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
		err = rows.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, err
		}
	}

	return &user, nil
}
