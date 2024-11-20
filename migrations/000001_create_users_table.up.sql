CREATE TABLE users
(
    id VARCHAR(50) PRIMARY KEY ,
    name VARCHAR(50) NOT NULL ,
    email VARCHAR(50) NOT NULL ,
    password VARCHAR(60) NOT NULL ,
    image TINYTEXT NOT NULL,
    created_at timestamp default current_timestamp,
    updated_at timestamp default current_timestamp on update CURRENT_TIMESTAMP
) ;
