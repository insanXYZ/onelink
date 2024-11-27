package repository

import (
	"context"
	"database/sql"
	"radproject/entity"
	"radproject/model"
)

type SiteRepository struct{}

func NewSiteRepository() *SiteRepository {
	return &SiteRepository{}
}

func (r *SiteRepository) SelectWhere(ctx context.Context, db SqlMetdod, field string, arguments ...any) (*entity.Sites, error) {
	site := new(entity.Sites)
	query := "select * from sites where " + field + " limit 1"
	err := db.QueryRowContext(ctx, query, arguments...).Scan(&site.Id, &site.Domain, &site.Title, &site.Image, &site.User_Id, &site.Created_At, &site.Updated_At)
	return site, err
}

// you shoud use s variable for represent sites table on field parameter
// for example
// // s.id -> sites.id
func (r *SiteRepository) SelectWhereWithJoinLink(ctx context.Context, db SqlMetdod, field string, arguments ...string) (*entity.Sites, error) {
	site := new(entity.Sites)
	query := "select * from sites s left join links l on s.id = l.site_id where " + field

	rows, err := db.QueryContext(ctx, query, arguments)
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

func (r *SiteRepository) SelectWithIdAndUser(ctx context.Context, db SqlMetdod, id, user_id string) (*entity.Sites, error) {
	site := new(entity.Sites)
	query := "select * from sites where id = ? and user_id = ?"
	err := db.QueryRowContext(ctx, query, id, user_id).Scan(&site.Id, &site.Domain, &site.Title, &site.Image, &site.User_Id, &site.Created_At, &site.Updated_At)
	return site, err
}

func (r *SiteRepository) SelectWithJoinLinkByDomain(ctx context.Context, db SqlMetdod, domain string) (*entity.Sites, error) {
	site := new(entity.Sites)
	query := "select * from sites s left join links l on s.id = l.site_id where s.domain = ?"

	rows, err := db.QueryContext(ctx, query, domain)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		link := entity.Links{}
		rows.Scan(&site.Id, &site.Domain, &site.Title, &site.Image, &site.User_Id, &site.Created_At, &site.Updated_At, &link.Id, &link.Title, &link.Href, &link.Site_Id, &link.CreatedAt, &link.UpdatedAt)
		site.Links = append(site.Links, link)
	}

	if site.Id == "" {
		return nil, sql.ErrNoRows
	}
	return site, nil
}

func (r *SiteRepository) SelectWithJoinLinkByUser(ctx context.Context, db SqlMetdod, id, user_id string) (*entity.Sites, error) {
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

func (r *SiteRepository) GetAllWithNumberClickByUser(ctx context.Context, db SqlMetdod, user_id string, req *model.DashboardRequest) ([]entity.Sites, []any, error) {
	sites := make([]entity.Sites, 0)
	ids := make([]any, 0)
	query := "select s.id, s.domain , s.title, s.image,s.created_at, count(c.clicked_at) from sites s left join clicks c on s.id = c.destination_id where s.user_id = ? and c.clicked_at between ? and ? group by s.id order by count(c.clicked_at)"
	rows, err := db.QueryContext(ctx, query, user_id, req.From, req.To)
	if err != nil {
		return nil, nil, err
	}

	defer rows.Close()

	for rows.Next() {
		site := entity.Sites{}
		err = rows.Scan(&site.Id, &site.Domain, &site.Title, &site.Image, &site.Created_At, &site.Clicks)
		if err != nil {
			return nil, nil, err
		}
		ids = append(ids, site.Id)
		sites = append(sites, site)
	}

	return sites, ids, nil
}
