package repository

import (
	"context"
	"database/sql"
	"radproject/entity"
)

type LinkRepository struct{}

func NewLinkRepository() *LinkRepository {
	return &LinkRepository{}
}

func (r *LinkRepository) SelectAllWithSiteId(ctx context.Context, db *sql.DB, site_id string) ([]entity.Links, error) {
	links := make([]entity.Links, 0)
	query := "select * from links where site_id = ?"
	rows, err := db.QueryContext(ctx, query, site_id)
	if err != nil {
		return links, err
	}
	defer rows.Close()
	for rows.Next() {
		link := entity.Links{}
		rows.Scan(&link.Id, &link.Title, &link.Href, &link.Site_Id, &link.CreatedAt, &link.UpdatedAt)
		links = append(links, link)
	}

	return links, err
}

func (r *LinkRepository) Save(ctx context.Context, db *sql.DB, ent *entity.Links) error {
	query := "insert into links(id,title,href,site_id) values(?,?,?,?)"
	_, err := db.ExecContext(ctx, query, ent.Id, ent.Title, ent.Href, ent.Site_Id)
	return err
}

func (r *LinkRepository) DeleteBySite(ctx context.Context, db SqlMetdod, site_id string) error {
	query := "delete from links where site_id = ?"
	_, err := db.ExecContext(ctx, query, site_id)
	return err
}

func (r *LinkRepository) Delete(ctx context.Context, db *sql.DB, id, site_id string) error {
	query := "delete from links where id = ? and site_id = ? "
	_, err := db.ExecContext(ctx, query, id, site_id)
	return err
}
