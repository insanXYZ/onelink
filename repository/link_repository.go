package repository

import (
	"context"
	"database/sql"
	"fmt"
	"radproject/entity"
	"radproject/model"
	"strings"
)

type LinkRepository struct{}

func NewLinkRepository() *LinkRepository {
	return &LinkRepository{}
}

func (r *LinkRepository) SelectWhere(ctx context.Context, db SqlMetdod, field string, arguments ...any) (*entity.Links, error) {
	ent := new(entity.Links)
	query := "select * from links where " + field + " limit 1"
	err := db.QueryRowContext(ctx, query, arguments...).Scan(&ent.Id, &ent.Title, &ent.Href, &ent.Site_Id, &ent.CreatedAt, &ent.UpdatedAt)
	return ent, err
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

func (r *LinkRepository) GetAllWithNumberClickBySiteId(ctx context.Context, db SqlMetdod, site_ids []any, req *model.DashboardRequest) ([]entity.Links, error) {
	links := make([]entity.Links, 0)

	placeholders := strings.Repeat("?,", len(site_ids)-1) + "?"

	query := fmt.Sprintf("select l.title, l.site_id,l.created_at, count(c.clicked_at) from links l left join clicks c on l.id = c.destination_id where l.site_id in (%s) and c.clicked_at between ? and ? group by l.id", placeholders)

	args := []any{}
	args = append(args, site_ids...)
	args = append(args, req.From, req.To)
	rows, err := db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		link := entity.Links{}

		err := rows.Scan(&link.Title, &link.Site_Id, &link.CreatedAt, &link.Clicks)
		if err != nil {
			return nil, err
		}

		links = append(links, link)
	}

	return links, nil
}
