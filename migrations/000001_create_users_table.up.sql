CREATE TABLE shortened_links (
    hash varchar(255) not null,
    user_id varchar(255) not null,
    original_url varchar(255) unique not null,
    created_at timestamp default current_timestamp
);
