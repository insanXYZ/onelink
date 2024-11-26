create table clicks 
(
  destination_id varchar(255) not null,
  clicked_at timestamp default current_timestamp
);
