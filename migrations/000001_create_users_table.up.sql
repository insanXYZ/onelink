CREATE TABLE users
(
    id VARCHAR(255) PRIMARY KEY ,
    name VARCHAR(255) NOT NULL ,
    email VARCHAR(255) NOT NULL ,
    password VARCHAR(255) NOT NULL ,
    image VARCHAR(255) NOT NULL,
    created_at timestamp default current_timestamp,
    updated_at timestamp default current_timestamp on update CURRENT_TIMESTAMP
) ;
