create table links
(
  id varchar(50) primary key not null,
  title varchar(50) not null,
  href longtext not null,
  site_id varchar(50) not null,
  created_at timestamp default current_timestamp,
  updated_at timestamp default current_timestamp on update CURRENT_TIMESTAMP,
  FOREIGN KEY(site_id) REFERENCES sites(id)
);
