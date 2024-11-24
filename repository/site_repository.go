package repository

import (
	"context"
	"database/sql"
	"radproject/entity"
)

type ExecSql interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
}
type SiteRepository struct{}

func NewSiteRepository() *SiteRepository {
	return &SiteRepository{}
}

func (r *SiteRepository) SelectAllById(ctx context.Context, db *sql.DB, id string) ([]entity.Sites, error) {
	res := make([]entity.Sites, 0)

	query := "select * from sites where user_id = ?"
	rows, err := db.QueryContext(ctx, query, id)
	if err != nil {
		return res, err
	}

	defer rows.Close()

	for rows.Next() {
		site := entity.Sites{}
		rows.Scan(&site.Id, &site.Domain, &site.Title, &site.Image, &site.User_Id, &site.Created_At, &site.Updated_At)
		res = append(res, site)
	}

	return res, err
}

func (r *SiteRepository) Save(ctx context.Context, db ExecSql, ent *entity.Sites) error {
	query := "insert into sites(id,domain,title,image,user_id) values(?,?,?,?,?)"
	_, err := db.ExecContext(ctx, query, ent.Id, ent.Domain, ent.Title, ent.Image, ent.User_Id)
	return err
}

func (r *SiteRepository) Delete(ctx context.Context, db ExecSql, id, user_id string) error {
	query := "delete from sites where id = ? and user_id = ?"
	_, err := db.ExecContext(ctx, query, id, user_id)
	return err
}
