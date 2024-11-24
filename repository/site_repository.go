package repository

import (
	"context"
	"database/sql"
	"radproject/entity"
)

type SiteRepository struct{}

func NewSiteRepository() *SiteRepository {
	return &SiteRepository{}
}

func (r *SiteRepository) SelectWithIdAndUser(ctx context.Context, db SqlMetdod, id, user_id string) (*entity.Sites, error) {
	site := new(entity.Sites)
	query := "select * from sites where id = ? and user_id = ?"
	err := db.QueryRowContext(ctx, query, id, user_id).Scan(&site.Id, &site.Domain, &site.Title, &site.Image, &site.User_Id, &site.Created_At, &site.Updated_At)
	return site, err
}

func (r *SiteRepository) SelectWithJoinLink(ctx context.Context, db SqlMetdod, id, user_id string) (*entity.Sites, error) {
	site := new(entity.Sites)
	query := "select * from sites s left join links l on s.id = l.site_id where s.id = ? and s.user_id = ?"

	rows, err := db.QueryContext(ctx, query, id, user_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		link := entity.Links{}
		rows.Scan(&site.Id, &site.Domain, &site.Title, &site.Image, &site.User_Id, &site.Created_At, &site.Updated_At, &link.Id, &link.Title, &link.Href, &link.Site_Id, &link.CreatedAt, &link.UpdatedAt)
		site.Links = append(site.Links, link)
	}
	return site, nil

}

func (r *SiteRepository) SelectAllByUserID(ctx context.Context, db *sql.DB, user_id string) ([]entity.Sites, error) {
	res := make([]entity.Sites, 0)

	query := "select * from sites where user_id = ?"
	rows, err := db.QueryContext(ctx, query, user_id)
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

func (r *SiteRepository) Save(ctx context.Context, db SqlMetdod, ent *entity.Sites) error {
	query := "insert into sites(id,domain,title,image,user_id) values(?,?,?,?,?)"
	_, err := db.ExecContext(ctx, query, ent.Id, ent.Domain, ent.Title, ent.Image, ent.User_Id)
	return err
}

func (r *SiteRepository) Delete(ctx context.Context, db SqlMetdod, id, user_id string) error {
	query := "delete from sites where id = ? and user_id = ?"
	_, err := db.ExecContext(ctx, query, id, user_id)
	return err
}

func (r *SiteRepository) Update(ctx context.Context, db SqlMetdod, ent *entity.Sites) error {
	query := "update sites set domain = ? ,title = ?,image =? where id = ?"
	_, err := db.ExecContext(ctx, query, ent.Domain, ent.Title, ent.Image, ent.Id)
	return err
}
