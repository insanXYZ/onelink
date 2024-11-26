create table sites 
(
  id varchar(255) primary key,
  domain varchar(255) not null,
  title varchar(255) not null,
  image tinytext not null,
  user_id varchar(255) not null,
  created_at timestamp default current_timestamp,
  updated_at timestamp default current_timestamp on update CURRENT_TIMESTAMP,
  FOREIGN KEY(user_id) REFERENCES users(id)
);
