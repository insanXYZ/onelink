package repository

import (
	"context"
)

type ClickRepository struct{}

func NewClikRepository() *ClickRepository {
	return &ClickRepository{}
}

func (r *ClickRepository) Save(ctx context.Context, db SqlMetdod, id string) error {
	query := "insert into clicks(destination_id) values(?)"
	_, err := db.ExecContext(ctx, query, id)
	return err
}
