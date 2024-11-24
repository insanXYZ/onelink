create table sites 
(
  id varchar(50) primary key,
  domain varchar(50) not null,
  title varchar(20) not null,
  image tinytext not null,
  user_id varchar(50) not null,
  created_at timestamp default current_timestamp,
  updated_at timestamp default current_timestamp on update CURRENT_TIMESTAMP,
  FOREIGN KEY(user_id) REFERENCES users(id)
);
